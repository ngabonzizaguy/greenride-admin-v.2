package services

import (
	"fmt"
	"sync"

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
		// Add other SMS service providers here
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
	if s.Handler == nil {
		return fmt.Errorf("SMS service not properly initialized")
	}

	return s.Handler.SendSmsMessage(message)
}
