package services

import (
	"errors"
	"fmt"
	"greenride/internal/i18n"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"strings"
)

var (
	messageService *MessageService
)

type MessageService struct {
	TemplateService *MessageTemplateService
	EmailService    *EmailService
	FirebaseService *FirebaseService // FCM服务
	SMSService      *SMSService      // SMS服务
}

func SetupMessageService() {
	messageService = &MessageService{
		TemplateService: GetMessageTemplateService(),
		EmailService:    GetEmailService(),
		FirebaseService: GetFirebaseService(),
		SMSService:      GetSMSService(),
	}
}

func GetMessageService() *MessageService {
	if messageService == nil {
		SetupMessageService()
	}
	return messageService
}

type Message struct {
	To          string         `json:"to,omitempty"`
	Type        string         `json:"type"`
	Channels    []string       `json:"channels"`
	Params      map[string]any `json:"params"`
	Language    string         `json:"language"`
	Region      string         `json:"region"`
	AttachUrls  []string       `json:"attach_urls,omitempty"`
	AttachFiles [][]byte       `json:"attach_files,omitempty"`
}

func (m *MessageService) SendMessage(message *Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	var lastErr error
	success := false

	for _, channel := range message.Channels {
		switch channel {
		case protocol.MsgChannelEmail:
			if m.EmailService != nil {
				if err := m.SendEmailMessage(message); err != nil {
					lastErr = err
				} else {
					success = true
				}
			}
		case protocol.MsgChannelFcm:
			if m.FirebaseService != nil {
				if err := m.SendFcmMessage(message); err != nil {
					lastErr = err
				} else {
					success = true
				}
			}
		case protocol.MsgChannelSms:
			if m.SMSService != nil {
				if err := m.SendSMSMessage(message); err != nil {
					lastErr = err
				} else {
					success = true
				}
			}
		default:
			lastErr = fmt.Errorf("unsupported message channel: %s", channel)
		}
	}

	if !success && lastErr != nil {
		return lastErr
	}

	return nil
}

func (m *MessageService) PrepareMessage(message *Message, channel string) *Message {
	// 获取模板
	template := m.TemplateService.GetTemplateByMessage(message, channel)
	if template == nil {
		return nil
	}

	// 构建参数
	params := protocol.MapData{}
	params.Copy(DefaultParams)
	params.Copy(message.Params)

	// 处理标题
	title := &strings.Builder{}
	if template.Title != nil {
		_ = template.Title.Execute(title, params)
	}
	params.Set("title", title.String())

	// 处理内容
	content := &strings.Builder{}
	if template.Content != nil {
		_ = template.Content.Execute(content, params)
	}
	params.Set("content", content.String())
	// 创建统一的消息对象
	result := &Message{
		To:          message.To,
		Language:    message.Language,
		Type:        message.Type,
		Channels:    []string{channel},
		AttachUrls:  message.AttachUrls,
		AttachFiles: message.AttachFiles,
		Params:      params,
	}
	return result
}

// getUserLanguage 获取用户的语言偏好，实现回退机制
// 优先级：用户偏好语言 -> 卢旺达语(rw) -> 英文(en)
func getUserLanguage(user *models.User) string {
	if user == nil {
		return protocol.LangKinyarwanda // 默认卢旺达语
	}

	userLang := user.GetLanguage()

	// 检查是否为支持的语言
	if i18n.IsValidLanguage(userLang) {
		return userLang
	}

	// 回退到卢旺达语
	return protocol.LangKinyarwanda
}
