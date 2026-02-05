package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MoMo API Status Mapping
var MoMoStatusMapping = map[string]string{
	"PENDING":    protocol.StatusPending,
	"SUCCESSFUL": protocol.StatusSuccess,
	"FAILED":     protocol.StatusFailed,
	"REJECTED":   protocol.StatusFailed,
	"TIMEOUT":    protocol.StatusFailed,
}

// MoMoConfig holds configuration for MTN MoMo API
type MoMoConfig struct {
	Environment       string `json:"environment"`        // "sandbox" or "production"
	SubscriptionKey   string `json:"subscription_key"`   // Ocp-Apim-Subscription-Key
	APIUserID         string `json:"api_user_id"`        // API user ID (X-Reference-Id)
	APIKey            string `json:"api_key"`            // Generated API key
	CallbackURL       string `json:"callback_url"`       // Webhook callback URL
	TargetEnvironment string `json:"target_environment"` // "sandbox" or country code (e.g., "rwandacollection")
	Currency          string `json:"currency"`           // Default currency (e.g., "RWF", "EUR")
	Timeout           int    `json:"timeout"`            // Request timeout in seconds
}

// Validate validates MoMo configuration
func (c *MoMoConfig) Validate() error {
	if c.SubscriptionKey == "" {
		return errors.New("MoMo subscription key is required")
	}
	if c.APIUserID == "" {
		return errors.New("MoMo API user ID is required")
	}
	if c.APIKey == "" {
		return errors.New("MoMo API key is required")
	}
	if c.Environment == "" {
		c.Environment = "sandbox"
	}
	if c.TargetEnvironment == "" {
		if c.Environment == "production" {
			c.TargetEnvironment = "rwandacollection"
		} else {
			c.TargetEnvironment = "sandbox"
		}
	}
	if c.Currency == "" {
		c.Currency = "RWF"
	}
	if c.Timeout <= 0 {
		c.Timeout = 30
	}
	// Set callback URL from global config if not set
	if c.CallbackURL == "" {
		if _cfg := config.Get(); _cfg != nil {
			if _cfg.Payment != nil && _cfg.Payment.CallbackHost != "" {
				c.CallbackURL = _cfg.Payment.CallbackHost + config.DefaultMoMoCallbackURL
			}
		}
	}
	return nil
}

// MoMoService implements PaymentChannel interface for MTN MoMo API
type MoMoService struct {
	config      *MoMoConfig
	accessToken string
	tokenExpiry int64
	mutex       sync.Mutex
}

// MoMoTokenResponse represents the token response from MoMo API
type MoMoTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// MoMoRequestToPayBody represents the request body for Request to Pay
type MoMoRequestToPayBody struct {
	Amount       string          `json:"amount"`
	Currency     string          `json:"currency"`
	ExternalID   string          `json:"externalId"`
	Payer        MoMoPayerInfo   `json:"payer"`
	PayerMessage string          `json:"payerMessage"`
	PayeeNote    string          `json:"payeeNote"`
}

// MoMoPayerInfo represents payer information
type MoMoPayerInfo struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}

// MoMoStatusResponse represents the status response from MoMo API
type MoMoStatusResponse struct {
	Amount                 string        `json:"amount"`
	Currency               string        `json:"currency"`
	FinancialTransactionID string        `json:"financialTransactionId"`
	ExternalID             string        `json:"externalId"`
	Payer                  MoMoPayerInfo `json:"payer"`
	PayerMessage           string        `json:"payerMessage"`
	PayeeNote              string        `json:"payeeNote"`
	Status                 string        `json:"status"`
	Reason                 *MoMoReason   `json:"reason,omitempty"`
}

// MoMoReason represents error reason
type MoMoReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewMoMoService creates a new MoMo service from PaymentChannels model
func NewMoMoService(channelConfig *models.PaymentChannels) *MoMoService {
	cfg := channelConfig.Config
	if len(cfg) == 0 {
		return nil
	}
	var momoConfig MoMoConfig
	cfg.ToObject(&momoConfig)
	return NewMoMoServiceWithConfig(&momoConfig)
}

// NewMoMoServiceWithConfig creates a new MoMo service with config
func NewMoMoServiceWithConfig(momoConfig *MoMoConfig) *MoMoService {
	if momoConfig == nil {
		return nil
	}
	if err := momoConfig.Validate(); err != nil {
		fmt.Printf("[MoMo] Config validation failed: %v\n", err)
		return nil
	}
	return &MoMoService{
		config: momoConfig,
	}
}

// getBaseURL returns the appropriate MoMo API base URL
func (s *MoMoService) getBaseURL() string {
	return config.GetMoMoBaseURL(s.config.Environment)
}

// getBasicAuth returns the Basic Auth header value
func (s *MoMoService) getBasicAuth() string {
	auth := fmt.Sprintf("%s:%s", s.config.APIUserID, s.config.APIKey)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// refreshTokenIfNeeded refreshes the access token if expired or not set
func (s *MoMoService) refreshTokenIfNeeded() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if token is still valid (with 60 second buffer)
	if s.accessToken != "" && time.Now().Unix() < s.tokenExpiry-60 {
		return nil
	}

	// Request new token
	url := fmt.Sprintf("%s/collection/token/", s.getBaseURL())

	headers := map[string]string{
		"Authorization":             s.getBasicAuth(),
		"Ocp-Apim-Subscription-Key": s.config.SubscriptionKey,
		"Content-Type":              "application/json",
	}

	body, resp, err := utils.PostJsonDataWithHeader(url, []byte(""), headers)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, body)
	}

	var tokenResp MoMoTokenResponse
	if err := json.Unmarshal([]byte(body), &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	s.accessToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Unix() + int64(tokenResp.ExpiresIn)

	fmt.Printf("[MoMo] Token refreshed, expires in %d seconds\n", tokenResp.ExpiresIn)
	return nil
}

// formatPhoneNumber formats phone number for MoMo API (MSISDN format)
func (s *MoMoService) formatPhoneNumber(phone string) string {
	// Remove spaces, dashes, and plus sign
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// MoMo expects MSISDN format (e.g., 250781234567 for Rwanda)
	// If phone starts with 0, replace with country code
	if strings.HasPrefix(phone, "0") && len(phone) == 10 {
		phone = "250" + phone[1:] // Rwanda country code
	}

	return phone
}

// Pay implements PaymentChannel interface - initiates a payment request
func (s *MoMoService) Pay(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypePayment,
		ChannelCode:   protocol.PaymentChannelMoMo,
	}

	// Validate payment amount
	if payment.GetAmount().LessThanOrEqual(decimal.Zero) {
		result.ResCode = protocol.ResCodeInvalidAmount
		result.ResMsg = "Invalid payment amount"
		return result
	}

	// Ensure we have a valid token
	if err := s.refreshTokenIfNeeded(); err != nil {
		result.ResCode = protocol.ResCodeAuthFailed
		result.ResMsg = err.Error()
		return result
	}

	// Generate unique reference ID for this transaction
	referenceID := uuid.New().String()

	// Build request body
	reqBody := MoMoRequestToPayBody{
		Amount:     payment.GetAmount().StringFixed(0), // MoMo expects integer amounts
		Currency:   payment.GetCurrency(),
		ExternalID: payment.PaymentID,
		Payer: MoMoPayerInfo{
			PartyIDType: "MSISDN",
			PartyID:     s.formatPhoneNumber(payment.GetPhone()),
		},
		PayerMessage: payment.GetOrderSku(),
		PayeeNote:    fmt.Sprintf("Order %s", payment.GetOrderID()),
	}

	// Use config currency if payment currency not set
	if reqBody.Currency == "" {
		reqBody.Currency = s.config.Currency
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		result.ResCode = protocol.ResCodeRequestFailed
		result.ResMsg = "Failed to marshal request body"
		return result
	}

	// Build callback URL with payment ID
	callbackURL := s.config.CallbackURL
	if callbackURL != "" && !strings.HasSuffix(callbackURL, "/") {
		callbackURL += "/"
	}
	callbackURL += payment.PaymentID

	// Build headers
	headers := map[string]string{
		"Authorization":             "Bearer " + s.accessToken,
		"X-Reference-Id":            referenceID,
		"X-Target-Environment":      s.config.TargetEnvironment,
		"Ocp-Apim-Subscription-Key": s.config.SubscriptionKey,
		"Content-Type":              "application/json",
		"X-Callback-Url":            callbackURL,
	}

	// Send request to pay
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", s.getBaseURL())

	fmt.Printf("[MoMo] Initiating payment: referenceID=%s, amount=%s, phone=%s\n",
		referenceID, reqBody.Amount, reqBody.Payer.PartyID)

	body, resp, err := utils.PostJsonDataWithHeader(url, bodyBytes, headers)
	if err != nil {
		result.ResCode = protocol.ResCodeRequestFailed
		result.ResMsg = fmt.Sprintf("Request failed: %v", err)
		return result
	}

	// MoMo returns 202 Accepted for successful request
	if resp.StatusCode == 202 {
		result.Status = protocol.StatusPending
		result.ChannelStatus = "PENDING"
		result.ChannelPaymentID = referenceID
		result.ResCode = "202"
		result.ResMsg = "Payment request accepted"
		fmt.Printf("[MoMo] Payment request accepted: referenceID=%s\n", referenceID)
		return result
	}

	// Handle error responses
	if resp.StatusCode == 400 {
		result.ResCode = protocol.ResCodeMissingFields
		result.ResMsg = fmt.Sprintf("Bad request: %s", body)
	} else if resp.StatusCode == 401 {
		result.ResCode = protocol.ResCodeAuthFailed
		result.ResMsg = "Authentication failed"
	} else if resp.StatusCode == 403 {
		result.ResCode = protocol.ResCodeAuthFailed
		result.ResMsg = "Access forbidden"
	} else if resp.StatusCode == 500 {
		result.ResCode = protocol.ResCodeChannelError
		result.ResMsg = "MoMo server error"
	} else {
		result.ResCode = protocol.ResCodeRequestFailed
		result.ResMsg = fmt.Sprintf("Request failed with status %d: %s", resp.StatusCode, body)
	}

	fmt.Printf("[MoMo] Payment request failed: status=%d, body=%s\n", resp.StatusCode, body)
	return result
}

// Refund implements PaymentChannel interface - processes refund
func (s *MoMoService) Refund(payment *models.Payment) *protocol.ChannelResult {
	// MoMo Collection API does not support direct refunds
	// Refunds must be processed via Disbursement API or manually
	return &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		ResCode:       protocol.ResCodeUnsupportedPaymentMethod,
		ResMsg:        "MoMo refunds require Disbursement API - please process manually",
		OrderType:     protocol.PaymentTypeRefund,
		ChannelCode:   protocol.PaymentChannelMoMo,
	}
}

// Status implements PaymentChannel interface - checks payment status
func (s *MoMoService) Status(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypePayment,
		ChannelCode:   protocol.PaymentChannelMoMo,
	}

	// Get the reference ID (stored as ChannelPaymentID)
	referenceID := payment.GetChannelPaymentID()
	if referenceID == "" {
		result.ResCode = protocol.ResCodeMissingFields
		result.ResMsg = "Missing reference ID"
		return result
	}

	// Ensure we have a valid token
	if err := s.refreshTokenIfNeeded(); err != nil {
		result.ResCode = protocol.ResCodeAuthFailed
		result.ResMsg = err.Error()
		return result
	}

	// Build headers
	headers := map[string]string{
		"Authorization":             "Bearer " + s.accessToken,
		"X-Target-Environment":      s.config.TargetEnvironment,
		"Ocp-Apim-Subscription-Key": s.config.SubscriptionKey,
	}

	// Send status request
	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", s.getBaseURL(), referenceID)

	body, resp, err := utils.GetWithHeader(url, headers)
	if err != nil {
		result.ResCode = protocol.ResCodeRequestFailed
		result.ResMsg = fmt.Sprintf("Status request failed: %v", err)
		return result
	}

	if resp.StatusCode != 200 {
		result.ResCode = protocol.ResCodeRequestFailed
		result.ResMsg = fmt.Sprintf("Status request failed with status %d: %s", resp.StatusCode, body)
		return result
	}

	// Parse status response
	var statusResp MoMoStatusResponse
	if err := json.Unmarshal([]byte(body), &statusResp); err != nil {
		result.ResCode = protocol.ResCodeResponseParseFailed
		result.ResMsg = fmt.Sprintf("Failed to parse status response: %v", err)
		return result
	}

	// Map MoMo status to system status
	result.ChannelStatus = statusResp.Status
	result.ChannelPaymentID = referenceID

	if systemStatus, ok := MoMoStatusMapping[statusResp.Status]; ok {
		result.Status = systemStatus
	} else {
		result.Status = protocol.StatusPending
	}

	if statusResp.FinancialTransactionID != "" {
		result.ResCode = statusResp.FinancialTransactionID
	} else {
		result.ResCode = statusResp.Status
	}

	if statusResp.Reason != nil {
		result.ResMsg = fmt.Sprintf("%s: %s", statusResp.Reason.Code, statusResp.Reason.Message)
	} else {
		result.ResMsg = statusResp.Status
	}

	result.CallbackData = body

	fmt.Printf("[MoMo] Status check: referenceID=%s, status=%s\n", referenceID, statusResp.Status)
	return result
}

// ResolveResponse parses MoMo webhook callback data
func (s *MoMoService) ResolveResponse(resp protocol.MapData) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		ChannelCode: protocol.PaymentChannelMoMo,
		OrderType:   protocol.PaymentTypePayment,
	}

	// Extract status from callback
	status := resp.Get("status")
	if status == "" {
		// Try alternate field names
		status = resp.Get("Status")
	}

	// Map status
	result.ChannelStatus = status
	if systemStatus, ok := MoMoStatusMapping[status]; ok {
		result.Status = systemStatus
	} else {
		result.Status = protocol.StatusPending
	}

	// Extract transaction ID
	financialTxID := resp.Get("financialTransactionId")
	if financialTxID == "" {
		financialTxID = resp.Get("FinancialTransactionId")
	}
	result.ResCode = financialTxID
	if result.ResCode == "" {
		result.ResCode = status
	}

	// Extract reference ID
	referenceID := resp.Get("referenceId")
	if referenceID == "" {
		referenceID = resp.Get("externalId")
	}
	result.ChannelPaymentID = referenceID

	// Extract reason if failed
	if reason := resp.Get("reason"); reason != "" {
		result.ResMsg = reason
	} else if reasonMap := resp.GetMapData("reason"); reasonMap != nil {
		code := reasonMap.Get("code")
		message := reasonMap.Get("message")
		result.ResMsg = fmt.Sprintf("%s: %s", code, message)
	} else {
		result.ResMsg = status
	}

	result.CallbackData = resp.ToJson()

	fmt.Printf("[MoMo] Webhook resolved: status=%s, channelStatus=%s\n", result.Status, result.ChannelStatus)
	return result
}
