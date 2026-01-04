package services

import (
	"fmt"
	"greenride/internal/protocol"

	"github.com/spf13/cast"
)

// SendSMSMessage 发送短信消息
func (m *MessageService) SendSMSMessage(message *Message) error {
	if m.SMSService == nil {
		return fmt.Errorf("SMS service not available")
	}

	// 使用辅助函数准备短信消息
	smsMessage := m.PrepareMessage(message, protocol.MsgChannelSms)
	if smsMessage == nil {
		return fmt.Errorf("failed to prepare SMS message: template not found for type %s", message.Type)
	}

	// 检查必要参数
	if to := cast.ToString(smsMessage.Params["to"]); to == "" {
		return fmt.Errorf("recipient ('to') is required for SMS")
	}

	if content := cast.ToString(smsMessage.Params["content"]); content == "" {
		return fmt.Errorf("content is required for SMS")
	}

	// 调用短信服务发送
	if err := m.SMSService.SendMessage(smsMessage); err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	return nil
}
