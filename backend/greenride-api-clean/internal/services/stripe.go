package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

// Stripe PaymentIntent status mapping
var StripeStatusMapping = map[stripe.PaymentIntentStatus]string{
	stripe.PaymentIntentStatusRequiresPaymentMethod: protocol.StatusPending,
	stripe.PaymentIntentStatusRequiresConfirmation:  protocol.StatusPending,
	stripe.PaymentIntentStatusRequiresAction:        protocol.StatusPending,
	stripe.PaymentIntentStatusProcessing:            protocol.StatusPending,
	stripe.PaymentIntentStatusRequiresCapture:       protocol.StatusPending,
	stripe.PaymentIntentStatusCanceled:              protocol.StatusCancelled,
	stripe.PaymentIntentStatusSucceeded:             protocol.StatusSuccess,
}

// StripeConfig holds configuration for Stripe API
type StripeConfig struct {
	SecretKey       string `json:"secret_key"`          // sk_test_... or sk_live_...
	PublishableKey  string `json:"publishable_key"`     // pk_test_... or pk_live_...
	WebhookSecret   string `json:"webhook_secret"`      // whsec_...
	Currency        string `json:"currency"`            // Default currency (e.g., "eur", "usd")
	StatementDesc   string `json:"statement_descriptor"` // Appears on customer statement (max 22 chars)
	Timeout         int    `json:"timeout"`             // Request timeout in seconds
}

// Validate validates Stripe configuration
func (c *StripeConfig) Validate() error {
	if c.SecretKey == "" {
		return errors.New("Stripe secret key is required")
	}
	if c.PublishableKey == "" {
		return errors.New("Stripe publishable key is required")
	}
	if c.Currency == "" {
		c.Currency = "eur"
	}
	c.Currency = strings.ToLower(c.Currency)
	if c.Timeout <= 0 {
		c.Timeout = 30
	}
	// Validate statement descriptor length
	if len(c.StatementDesc) > 22 {
		c.StatementDesc = c.StatementDesc[:22]
	}
	return nil
}

// StripeService implements PaymentChannel interface for Stripe
type StripeService struct {
	config *StripeConfig
}

// NewStripeService creates a new Stripe service from PaymentChannels model
func NewStripeService(channelConfig *models.PaymentChannels) *StripeService {
	cfg := channelConfig.Config
	if len(cfg) == 0 {
		return nil
	}
	var stripeConfig StripeConfig
	cfg.ToObject(&stripeConfig)
	return NewStripeServiceWithConfig(&stripeConfig)
}

// NewStripeServiceWithConfig creates a new Stripe service with config
func NewStripeServiceWithConfig(stripeConfig *StripeConfig) *StripeService {
	if stripeConfig == nil {
		return nil
	}
	if err := stripeConfig.Validate(); err != nil {
		fmt.Printf("[Stripe] Config validation failed: %v\n", err)
		return nil
	}
	return &StripeService{
		config: stripeConfig,
	}
}

// GetPublishableKey returns the publishable key for frontend use
func (s *StripeService) GetPublishableKey() string {
	return s.config.PublishableKey
}

// Pay implements PaymentChannel interface - creates a PaymentIntent
func (s *StripeService) Pay(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypePayment,
		ChannelCode:   protocol.PaymentChannelStripe,
	}

	// Validate payment amount
	if payment.GetAmount().LessThanOrEqual(decimal.Zero) {
		result.ResCode = protocol.ResCodeInvalidAmount
		result.ResMsg = "Invalid payment amount"
		return result
	}

	// Set Stripe API key
	stripe.Key = s.config.SecretKey

	// Convert amount to cents/smallest currency unit
	// Stripe expects amounts in smallest currency unit (e.g., cents for USD/EUR)
	amountInCents := payment.GetAmount().Mul(decimal.NewFromInt(100)).IntPart()

	// Determine currency
	currency := strings.ToLower(payment.GetCurrency())
	if currency == "" {
		currency = s.config.Currency
	}

	// Create PaymentIntent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"payment_id": payment.PaymentID,
			"order_id":   payment.GetOrderID(),
			"user_id":    payment.GetUserID(),
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Add description
	if payment.GetOrderSku() != "" {
		params.Description = stripe.String(payment.GetOrderSku())
	}

	// Add statement descriptor if configured
	if s.config.StatementDesc != "" {
		params.StatementDescriptor = stripe.String(s.config.StatementDesc)
	}

	// Add customer email if available
	if email := payment.GetEmail(); email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	fmt.Printf("[Stripe] Creating PaymentIntent: amount=%d, currency=%s, payment_id=%s\n",
		amountInCents, currency, payment.PaymentID)

	// Create PaymentIntent
	pi, err := paymentintent.New(params)
	if err != nil {
		stripeErr, ok := err.(*stripe.Error)
		if ok {
			result.ResCode = string(stripeErr.Code)
			result.ResMsg = stripeErr.Msg
		} else {
			result.ResCode = protocol.ResCodeRequestFailed
			result.ResMsg = err.Error()
		}
		fmt.Printf("[Stripe] PaymentIntent creation failed: %v\n", err)
		return result
	}

	// Map status
	result.Status = protocol.StatusPending
	if systemStatus, ok := StripeStatusMapping[pi.Status]; ok {
		result.Status = systemStatus
	}
	result.ChannelStatus = string(pi.Status)
	result.ChannelPaymentID = pi.ID
	result.ResCode = "created"
	result.ResMsg = "PaymentIntent created successfully"

	// Return client_secret and publishable_key for frontend
	result.Metadata = protocol.MapData{
		"client_secret":   pi.ClientSecret,
		"publishable_key": s.config.PublishableKey,
		"payment_intent":  pi.ID,
	}

	fmt.Printf("[Stripe] PaymentIntent created: id=%s, status=%s\n", pi.ID, pi.Status)
	return result
}

// Refund implements PaymentChannel interface - processes refund
func (s *StripeService) Refund(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypeRefund,
		ChannelCode:   protocol.PaymentChannelStripe,
	}

	// Get the PaymentIntent ID
	paymentIntentID := payment.GetChannelPaymentID()
	if paymentIntentID == "" {
		result.ResCode = protocol.ResCodeMissingFields
		result.ResMsg = "Missing PaymentIntent ID"
		return result
	}

	// Set Stripe API key
	stripe.Key = s.config.SecretKey

	// Create refund parameters
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	// If partial refund, set amount
	refundAmount := payment.GetRefundAmount()
	if refundAmount.GreaterThan(decimal.Zero) {
		amountInCents := refundAmount.Mul(decimal.NewFromInt(100)).IntPart()
		params.Amount = stripe.Int64(amountInCents)
	}

	// Add reason if provided
	if payment.GetRefundReason() != "" {
		params.Reason = stripe.String(string(stripe.RefundReasonRequestedByCustomer))
		params.AddMetadata("reason", payment.GetRefundReason())
	}

	fmt.Printf("[Stripe] Creating refund: paymentIntent=%s\n", paymentIntentID)

	// Create refund
	ref, err := refund.New(params)
	if err != nil {
		stripeErr, ok := err.(*stripe.Error)
		if ok {
			result.ResCode = string(stripeErr.Code)
			result.ResMsg = stripeErr.Msg
		} else {
			result.ResCode = protocol.ResCodeRequestFailed
			result.ResMsg = err.Error()
		}
		fmt.Printf("[Stripe] Refund creation failed: %v\n", err)
		return result
	}

	// Map refund status
	switch ref.Status {
	case stripe.RefundStatusSucceeded:
		result.Status = protocol.StatusRefunded
		result.ChannelStatus = "succeeded"
	case stripe.RefundStatusPending:
		result.Status = protocol.StatusPending
		result.ChannelStatus = "pending"
	case stripe.RefundStatusFailed:
		result.Status = protocol.StatusFailed
		result.ChannelStatus = "failed"
	case stripe.RefundStatusCanceled:
		result.Status = protocol.StatusCancelled
		result.ChannelStatus = "canceled"
	}

	result.ChannelPaymentID = ref.ID
	result.ResCode = string(ref.Status)
	result.ResMsg = "Refund processed"

	fmt.Printf("[Stripe] Refund created: id=%s, status=%s\n", ref.ID, ref.Status)
	return result
}

// Status implements PaymentChannel interface - checks payment status
func (s *StripeService) Status(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypePayment,
		ChannelCode:   protocol.PaymentChannelStripe,
	}

	// Get the PaymentIntent ID
	paymentIntentID := payment.GetChannelPaymentID()
	if paymentIntentID == "" {
		result.ResCode = protocol.ResCodeMissingFields
		result.ResMsg = "Missing PaymentIntent ID"
		return result
	}

	// Set Stripe API key
	stripe.Key = s.config.SecretKey

	// Retrieve PaymentIntent
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		stripeErr, ok := err.(*stripe.Error)
		if ok {
			result.ResCode = string(stripeErr.Code)
			result.ResMsg = stripeErr.Msg
		} else {
			result.ResCode = protocol.ResCodeRequestFailed
			result.ResMsg = err.Error()
		}
		fmt.Printf("[Stripe] Status check failed: %v\n", err)
		return result
	}

	// Map status
	result.ChannelStatus = string(pi.Status)
	result.ChannelPaymentID = pi.ID
	if systemStatus, ok := StripeStatusMapping[pi.Status]; ok {
		result.Status = systemStatus
	} else {
		result.Status = protocol.StatusPending
	}
	result.ResCode = string(pi.Status)
	result.ResMsg = "Status retrieved"

	// Serialize to callback data
	if piBytes, err := json.Marshal(pi); err == nil {
		result.CallbackData = string(piBytes)
	}

	fmt.Printf("[Stripe] Status check: id=%s, status=%s\n", pi.ID, pi.Status)
	return result
}

// VerifyWebhookSignature verifies the Stripe webhook signature and returns the event
func (s *StripeService) VerifyWebhookSignature(payload []byte, signature string) (*stripe.Event, error) {
	if s.config.WebhookSecret == "" {
		return nil, errors.New("webhook secret not configured")
	}

	event, err := webhook.ConstructEvent(payload, signature, s.config.WebhookSecret)
	if err != nil {
		return nil, fmt.Errorf("webhook signature verification failed: %w", err)
	}

	return &event, nil
}

// ResolvePaymentIntentEvent processes a PaymentIntent event and returns ChannelResult
func (s *StripeService) ResolvePaymentIntentEvent(pi *stripe.PaymentIntent) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		ChannelCode:      protocol.PaymentChannelStripe,
		OrderType:        protocol.PaymentTypePayment,
		ChannelPaymentID: pi.ID,
		ChannelStatus:    string(pi.Status),
	}

	// Map status
	if systemStatus, ok := StripeStatusMapping[pi.Status]; ok {
		result.Status = systemStatus
	} else {
		result.Status = protocol.StatusPending
	}

	result.ResCode = string(pi.Status)

	// Set appropriate message based on status
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		result.ResMsg = "Payment succeeded"
	case stripe.PaymentIntentStatusCanceled:
		if pi.CancellationReason != "" {
			result.ResMsg = fmt.Sprintf("Payment canceled: %s", pi.CancellationReason)
		} else {
			result.ResMsg = "Payment canceled"
		}
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		if pi.LastPaymentError != nil {
			result.ResMsg = pi.LastPaymentError.Msg
		} else {
			result.ResMsg = "Payment method required"
		}
	default:
		result.ResMsg = string(pi.Status)
	}

	// Serialize to callback data
	if piBytes, err := json.Marshal(pi); err == nil {
		result.CallbackData = string(piBytes)
	}

	return result
}

// ResolveResponse parses webhook data for generic processing
func (s *StripeService) ResolveResponse(resp protocol.MapData) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		ChannelCode: protocol.PaymentChannelStripe,
		OrderType:   protocol.PaymentTypePayment,
	}

	// Extract status
	status := resp.Get("status")
	result.ChannelStatus = status

	// Map status
	if status == "succeeded" {
		result.Status = protocol.StatusSuccess
	} else if status == "canceled" {
		result.Status = protocol.StatusCancelled
	} else if status == "requires_payment_method" || status == "processing" {
		result.Status = protocol.StatusPending
	} else {
		result.Status = protocol.StatusPending
	}

	// Extract IDs
	result.ChannelPaymentID = resp.Get("id")
	result.ResCode = status
	result.ResMsg = status
	result.CallbackData = resp.ToJson()

	return result
}
