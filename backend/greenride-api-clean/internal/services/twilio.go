package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"

	"github.com/twilio/twilio-go"
	twiApi "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// TwilioService represents the Twilio SMS service implementation
type TwilioService struct {
	// Map of phone numbers to Twilio accounts
	phoneToAccount map[string]*TwilioAccount
	// Default account to use when no specific phone match is found
	defaultAccount *TwilioAccount
}

// TwilioAccount represents a single Twilio account configuration
type TwilioAccount struct {
	AccountSID string
	AuthToken  string
	Client     *twilio.RestClient
	Phones     []string // List of phone numbers associated with this account
	ServiceSID string   // Verify Service SID (optional, for Verify API)
}

var (
	twilioServiceInstance *TwilioService
	twilioServiceOnce     sync.Once
)

// GetTwilioService returns the singleton instance of TwilioService
func GetTwilioService() *TwilioService {
	twilioServiceOnce.Do(func() {
		SetupTwilioService()
	})
	return twilioServiceInstance
}

// SetupTwilioService initializes the Twilio service with configuration
func SetupTwilioService() {
	cfg := config.Get()
	if cfg.Twilio == nil {
		log.Get().Warn("Twilio configuration is nil, skipping Twilio service setup")
		return
	}

	service := &TwilioService{
		phoneToAccount: make(map[string]*TwilioAccount),
	}

	// Load all Twilio accounts from configuration
	for i, accountCfg := range cfg.Twilio.Accounts {
		account := &TwilioAccount{
			AccountSID: accountCfg.AccountSID,
			AuthToken:  accountCfg.AuthToken,
			Phones:     accountCfg.Phones,
			ServiceSID: accountCfg.ServiceSID, // Optional for Verify API
			Client: twilio.NewRestClientWithParams(twilio.ClientParams{
				Username: accountCfg.AccountSID,
				Password: accountCfg.AuthToken,
			}),
		}
		// For OTP flows, we prefer failing fast and falling back to InnoPaaS rather than waiting
		// on slow Twilio Verify requests. This timeout applies to Twilio requests for this client.
		// NOTE: This is set once at init to avoid races from per-request toggling.
		account.Client.SetTimeout(5 * time.Second)

		// If no default account is set and this is the first account, use it as default
		if service.defaultAccount == nil && i == 0 {
			service.defaultAccount = account
		}

		// Map each phone number to this account
		for _, phone := range accountCfg.Phones {
			service.phoneToAccount[phone] = account
		}
	}

	twilioServiceInstance = service
}

// SendSMS sends an SMS message using Twilio Message API
// This method allows complete control over message content
// from: optional sender phone number, will use the first available if not specified
// If a prioritized list of numbers is provided, it will try them in order
func (t *TwilioService) SendSMS(to string, message string, from string) error {
	// Get the appropriate account and phone number
	account, fromNumber := t.getAccountAndPhone(from)
	if account == nil {
		return fmt.Errorf("no valid Twilio account found for sending SMS")
	}

	// Create the message parameters
	params := &twiApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	// Send the message
	resp, err := account.Client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %v", err)
	}

	// Log the message SID for tracking
	if resp.Sid != nil {
		log.Get().Infof("SMS sent via Twilio with SID: %s", *resp.Sid)
	}

	return nil
}

// ServiceName returns the service name
func (t *TwilioService) ServiceName() string {
	return "twilio"
}

// SendCustomVerificationCode sends a custom verification code using Twilio Verify API with custom content
// Uses Twilio Verify API's customization options for better security and compliance
func (t *TwilioService) SendCustomVerificationCode(to, code, locale string) error {
	account := t.defaultAccount
	if account == nil || account.ServiceSID == "" {
		return fmt.Errorf("no Twilio Verify service configured")
	}

	// Create verification parameters
	params := &verify.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	// Set locale if provided
	if locale != "" {
		params.SetLocale(locale)
	}

	// Set custom code if provided
	if code != "" {
		params.SetCustomCode(code)
	}

	// Send verification code using Verify API
	resp, err := account.Client.VerifyV2.CreateVerification(account.ServiceSID, params)
	if err != nil {
		return fmt.Errorf("failed to send custom verification code: %v", err)
	}

	log.Get().Infof("Custom verification code [%v] sent to %s with locale %s, SID: %s", code, to, locale, *resp.Sid)
	return nil
}

// VerifyCode verifies a code using Twilio Verify API
func (t *TwilioService) VerifyCode(to, code string) (bool, error) {
	account := t.defaultAccount
	if account == nil || account.ServiceSID == "" {
		return false, fmt.Errorf("no Twilio Verify service configured")
	}

	// Create verification check parameters
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	// Verify the code
	resp, err := account.Client.VerifyV2.CreateVerificationCheck(account.ServiceSID, params)
	if err != nil {
		return false, fmt.Errorf("failed to verify code: %v", err)
	}

	isValid := resp.Status != nil && *resp.Status == "approved"
	log.Get().Infof("Code verification for %s: %s", to, *resp.Status)
	return isValid, nil
}

// getAccountAndPhone determines which Twilio account and phone number to use
// It handles the prioritized list of "from" numbers if provided
func (t *TwilioService) getAccountAndPhone(from string) (*TwilioAccount, string) {
	// If a specific from number is provided, try to use it
	if from != "" {
		// Check if we have a specific account for this phone
		if account, ok := t.phoneToAccount[from]; ok {
			return account, from
		}

		// If from contains multiple numbers (comma-separated), try them in order
		if phones := splitPhoneNumbers(from); len(phones) > 0 {
			for _, phone := range phones {
				if account, ok := t.phoneToAccount[phone]; ok {
					return account, phone
				}
			}
		}
	}

	// If no specific or valid "from" is provided, use the default account
	if t.defaultAccount != nil {
		// Use the first phone number from the default account
		if len(t.defaultAccount.Phones) > 0 {
			return t.defaultAccount, t.defaultAccount.Phones[0]
		}
	}

	return nil, ""
}

// splitPhoneNumbers splits a comma-separated list of phone numbers
func splitPhoneNumbers(phones string) []string {
	if phones == "" {
		return nil
	}

	// Simple split by comma and trim spaces
	var result []string
	for _, phone := range strings.Split(phones, ",") {
		phone = strings.TrimSpace(phone)
		if phone != "" {
			result = append(result, phone)
		}
	}
	return result
}
