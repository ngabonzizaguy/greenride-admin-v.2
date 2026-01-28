package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
)

// InnoPaaSService represents the InnoPaaS SMS service implementation
type InnoPaaSService struct {
	client *http.Client
}

var (
	innoPaaSServiceInstance *InnoPaaSService
	innoPaaSServiceOnce     sync.Once
)

// GetInnoPaaSService returns the singleton instance
func GetInnoPaaSService() *InnoPaaSService {
	innoPaaSServiceOnce.Do(func() {
		innoPaaSServiceInstance = &InnoPaaSService{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}
	})
	return innoPaaSServiceInstance
}

// ServiceName returns the service name
func (s *InnoPaaSService) ServiceName() string {
	return "innopaas"
}

// InnoPaaSRequest represents the OTP API v3.0 request body
// type: "1"=WhatsApp, "3"=SMS
type InnoPaaSRequest struct {
	Type     string `json:"type"`     // "3" for SMS
	Language string `json:"language"` // "en"
	To       string `json:"to"`      // E.164 e.g. "+12025551234"
	Code     string `json:"code"`
	Sender   string `json:"sender,omitempty"` // Optional
}

// InnoPaaSResponse represents the OTP API v3.0 response (success when code == "000000")
type InnoPaaSResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"` // Message ID on success
}

// SendSmsMessage sends an OTP using InnoPaaS OTP API v3.0
func (s *InnoPaaSService) SendSmsMessage(message *Message) error {
	cfg := config.Get().InnoPaaS
	if cfg == nil {
		return fmt.Errorf("InnoPaaS configuration is missing")
	}

	to, ok := message.Params["to"].(string)
	if !ok || to == "" {
		return fmt.Errorf("missing recipient phone number")
	}

	var code string
	if codeVal, ok := message.Params["code"]; ok && codeVal != nil {
		code = fmt.Sprintf("%v", codeVal)
	} else {
		content, ok := message.Params["content"].(string)
		if !ok || content == "" {
			return fmt.Errorf("missing OTP code in message params")
		}
		code = extractCodeFromContent(content)
		if code == "" {
			return fmt.Errorf("failed to extract OTP code from message content")
		}
	}
	if code == "" {
		return fmt.Errorf("missing OTP code")
	}

	// OTP v3.0: to in international format e.g. "+12025551234"
	toE164 := "+" + cleanPhoneNumber(to)

	reqBody := InnoPaaSRequest{
		Type:     "3",  // SMS
		Language: "en",
		To:       toE164,
		Code:     code,
	}
	if cfg.SenderID != "" {
		reqBody.Sender = cfg.SenderID
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", cfg.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	if cfg.AppSecret != "" {
		// OTP v3.0: x-appKey, x-nonce, x-timestamp, x-signature (MD5 over sorted params + secret)
		nonce := strconv.FormatInt(time.Now().UnixMilli(), 10)
		timestamp := nonce
		req.Header.Set("x-appKey", cfg.AppKey)
		req.Header.Set("x-nonce", nonce)
		req.Header.Set("x-timestamp", timestamp)
		sign := innopaasSign(cfg.AppSecret, cfg.AppKey, nonce, timestamp, reqBody)
		req.Header.Set("x-signature", sign)
	} else {
		// Legacy: Authorization token + appKey (if no app_secret configured)
		req.Header.Set("Authorization", cfg.Authorization)
		req.Header.Set("appKey", cfg.AppKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to InnoPaaS: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResp InnoPaaSResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Get().Errorf("InnoPaaS raw response: %s", string(body))
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// v3.0 success is code == "000000"
	if apiResp.Code != "000000" {
		return fmt.Errorf("InnoPaaS API error: %s (code: %s)", apiResp.Message, apiResp.Code)
	}

	log.Get().Infof("SMS sent via InnoPaaS, MessageID: %s", apiResp.Data)
	return nil
}

// innopaasSign computes InnoPaaS MD5 signature: sort params, concat key+value for non-blank, append secret, MD5 hex lower
func innopaasSign(secret, appKey, nonce, timestamp string, body InnoPaaSRequest) string {
	params := map[string]string{
		"appKey":    appKey,
		"code":      body.Code,
		"language":  body.Language,
		"nonce":     nonce,
		"timestamp": timestamp,
		"to":        body.To,
		"type":      body.Type,
	}
	if body.Sender != "" {
		params["sender"] = body.Sender
	}
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		v := params[k]
		if v != "" {
			sb.WriteString(k)
			sb.WriteString(v)
		}
	}
	sb.WriteString(secret)
	sum := md5.Sum([]byte(sb.String()))
	return strings.ToLower(hex.EncodeToString(sum[:]))
}

// cleanPhoneNumber removes '+' and other non-digit characters
func cleanPhoneNumber(phone string) string {
	// Implementation to keep only digits
	// User example shows '250...'
	var result []rune
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			result = append(result, r)
		}
	}
	return string(result)
}

// extractCodeFromContent extracts numeric OTP code from formatted message content
// Handles formats like "Your code is 1234", "Code: 1234", "1234", etc.
func extractCodeFromContent(content string) string {
	// First, try to find a sequence of 4-6 digits (typical OTP length)
	re := regexp.MustCompile(`\b\d{4,6}\b`)
	matches := re.FindString(content)
	if matches != "" {
		return matches
	}
	// If no match, return empty string
	return ""
}
