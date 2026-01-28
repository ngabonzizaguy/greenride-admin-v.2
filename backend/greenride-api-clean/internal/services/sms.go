package services

import (
	"fmt"
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

		log.Get().Warnf("OTP send via twilio failed after %s; falling back to innopaas: %v", elapsed, err)
		fallback := GetInnoPaaSService()
		if fallback == nil {
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
