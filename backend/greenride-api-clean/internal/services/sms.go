package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/protocol"
)

// SMSService represents the SMS service
type SMSService struct {
	Handler SmsMessage
}

// SmsMessage 短信服务接口
type SmsMessage interface {
	ServiceName() string
	SendSmsMessage(message *Message) error
}

var (
	smsServiceInstance *SMSService
	smsServiceOnce     sync.Once
)

// GetSMSService returns the singleton instance of SMSService
func GetSMSService() *SMSService {
	smsServiceOnce.Do(func() {
		SetupSMSService()
	})
	return smsServiceInstance
}

// SetupSMSService initializes the SMS service
func SetupSMSService() {
	smsServiceInstance = &SMSService{
		Handler: GetSMSServiceHandler(),
	}

	// Startup diagnostics: warn about misconfigured providers
	cfg := config.Get()
	if cfg.SMS != nil && cfg.SMS.ServiceName == "twilio" {
		if cfg.InnoPaaS == nil {
			log.Get().Warn("[SMS] InnoPaaS config is nil — Twilio has no fallback provider")
		} else if strings.Contains(cfg.InnoPaaS.AppKey, "PASTE_YOUR") || cfg.InnoPaaS.AppKey == "" {
			log.Get().Warn("[SMS] InnoPaaS has placeholder credentials — fallback will fail if Twilio is down")
		}
	}
	if cfg.Twilio == nil || len(cfg.Twilio.Accounts) == 0 {
		log.Get().Error("[SMS] No Twilio accounts configured — OTP sending will fail")
	} else {
		for _, acc := range cfg.Twilio.Accounts {
			if acc.ServiceSID == "" {
				log.Get().Warnf("[SMS] Twilio account %s has no Verify Service SID — OTP via Verify API will fail", acc.AccountSID)
			}
		}
		log.Get().Infof("[SMS] Twilio initialized with %d account(s), primary service=%s", len(cfg.Twilio.Accounts), cfg.SMS.ServiceName)
	}
}

// GetSMSServiceHandler returns the SMS service handler
func GetSMSServiceHandler() SmsMessage {
	cfg := config.Get().SMS

	// Select SMS service based on configuration
	switch cfg.ServiceName {
	case "twilio":
		return GetTwilioService()
	case "innopaas":
		return GetInnoPaaSService()
	}

	return nil
}

// SendSMS sends an SMS message
// from: optional phone number(s) to send from (comma-separated for priority list)
func (s *SMSService) SendSMS(to string, message string, from string) error {
	if s.Handler == nil {
		return fmt.Errorf("SMS service not properly initialized")
	}

	// Create a Message object to pass to SendSmsMessage
	msg := &Message{
		Type:     protocol.MsgTypeGeneric, // Use generic message type
		Channels: []string{protocol.MsgChannelSms},
		Params: map[string]any{
			"to":      to,
			"content": message,
			"from":    from,
		},
	}

	// Send via SMS service handler
	err := s.Handler.SendSmsMessage(msg)
	if err != nil {
		log.Get().Errorf("Failed to send SMS: %v", err)
		return err
	}

	return nil
}

// ServiceName returns the service name
func (s *SMSService) ServiceName() string {
	return "sms"
}

// SendMessage sends an SMS using a Message object
func (s *SMSService) SendMessage(message *Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Resolve handler defensively in case service init didn't run as expected.
	primary := s.Handler
	if primary == nil {
		primary = GetSMSServiceHandler()
	}
	if primary == nil {
		return fmt.Errorf("SMS service not properly initialized")
	}

	cfg := config.Get().SMS
	// OTP failover policy:
	// - If primary is Twilio and this is a verify_code message, attempt Twilio (fast-timeout client).
	// - If Twilio errors (including timeout), fall back to InnoPaaS.
	if cfg != nil && cfg.ServiceName == "twilio" && message.Type == protocol.MsgTypeVerifyCode {
		start := time.Now()
		err := primary.SendSmsMessage(message)
		elapsed := time.Since(start)
		if err == nil {
			log.Get().Infof("OTP sent via twilio in %s", elapsed)
			return nil
		}

		log.Get().Warnf("[SMS] OTP via Twilio FAILED after %s: %v — attempting InnoPaaS fallback", elapsed, err)
		fallback := GetInnoPaaSService()
		if fallback == nil {
			log.Get().Error("[SMS] InnoPaaS fallback is nil — no fallback available, OTP will NOT be delivered")
			return err
		}
		startFb := time.Now()
		fbErr := fallback.SendSmsMessage(message)
		fbElapsed := time.Since(startFb)
		if fbErr == nil {
			log.Get().Infof("OTP sent via innopaas in %s (after twilio failure)", fbElapsed)
			return nil
		}
		log.Get().Errorf("OTP send via innopaas also failed after %s: %v", fbElapsed, fbErr)
		return fmt.Errorf("twilio otp failed: %v; innopaas otp failed: %v", err, fbErr)
	}

	return primary.SendSmsMessage(message)
}
