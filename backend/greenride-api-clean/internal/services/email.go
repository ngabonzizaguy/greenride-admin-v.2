package services

import (
	"fmt"
	"greenride/internal/config"
	"greenride/internal/protocol"
)

var (
	emailService *EmailService
)

type EmailService struct {
	Handler EmailMessage
}

type EmailMessage interface {
	ServiceName() string
	SendEmailMessage(message *Message) error
}

func SetupEmailService() {
	emailService = &EmailService{
		Handler: GetEmailServiceHandler(),
	}
}

func GetEmailServiceHandler() EmailMessage {
	cfg := config.Get().Email

	// 根据配置选择邮件服务
	switch cfg.ServiceName {
	case "account":
		// 使用自建邮箱账户服务
		return NewEmailAccountService(cfg.Accounts)
	}

	return nil
}

func GetEmailService() *EmailService {
	if emailService == nil {
		SetupEmailService()
	}
	return emailService
}

func (e *EmailService) SendEmail(to string, subject string, body string) error {
	if e.Handler != nil {
		// 创建一个Message对象传递给SendEmailMessage
		message := &Message{
			Type:     protocol.MsgTypeGeneric, // 使用通用消息类型
			Channels: []string{protocol.MsgChannelEmail},
			Params: map[string]any{
				"to":      to,
				"title":   subject,
				"content": body,
			},
		}
		return e.Handler.SendEmailMessage(message)
	}
	return fmt.Errorf("invalid email service")
}

// SendMessage sends an email using a Message object
func (e *EmailService) SendMessage(message *Message) error {
	if e.Handler != nil {
		return e.Handler.SendEmailMessage(message)
	}
	return fmt.Errorf("invalid email service")
}
