package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// InnoPaaSRequest represents the API request body
type InnoPaaSRequest struct {
	To       string `json:"to"`
	Type     string `json:"type"`     // "3" for OTP
	Language string `json:"language"` // "en"
	Code     string `json:"code"`
}

// InnoPaaSResponse represents the API response
type InnoPaaSResponse struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"` // Message ID
}

// SendSmsMessage sends an SMS using InnoPaaS OTP API
func (s *InnoPaaSService) SendSmsMessage(message *Message) error {
	cfg := config.Get().InnoPaaS
	if cfg == nil {
		return fmt.Errorf("InnoPaaS configuration is missing")
	}

	// Extract phone number (to)
	to, ok := message.Params["to"].(string)
	if !ok || to == "" {
		return fmt.Errorf("missing recipient phone number")
	}

	// Extract content. For this specific API, we need the "code"
	// The message content usually comes in as "Your code is 1234"
	// We need to either pass just the code in the 'content' param or extract it.
	// Assumption: The caller (VerifyCodeService) will pass the code as 'content'
	// OR we assume the content IS the code if it's numeric.
	content, ok := message.Params["content"].(string)
	if !ok || content == "" {
		return fmt.Errorf("missing message content")
	}

	// Prepare request body
	// "to" must be digits only (e.g., 250...)
	// In user's example: "to": "250784928786" (no plus sign)
	cleanTo := cleanPhoneNumber(to)

	reqBody := InnoPaaSRequest{
		To:       cleanTo,
		Type:     "3",
		Language: "en",
		Code:     content, // Assuming content is the OTP code
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create Request
	req, err := http.NewRequest("POST", cfg.Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	// "Authorization" header is the "secret"
	req.Header.Set("Authorization", cfg.Authorization)
	req.Header.Set("appKey", cfg.AppKey)

	// Send Request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to InnoPaaS: %v", err)
	}
	defer resp.Body.Close()

	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse Response
	var apiResp InnoPaaSResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Get().Errorf("InnoPaaS raw response: %s", string(body))
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("InnoPaaS API returned error: %s (Code: %s)", apiResp.Message, apiResp.Code)
	}

	log.Get().Infof("SMS sent via InnoPaaS, MessageID: %s", apiResp.Data)
	return nil
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
