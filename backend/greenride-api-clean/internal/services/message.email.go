package services

import (
	"fmt"
	"greenride/internal/protocol"

	"github.com/spf13/cast"
)

// 具体发送邮件的业务逻辑
func (m *MessageService) SendEmailMessage(message *Message) error {
	if m.EmailService == nil {
		return fmt.Errorf("email service not available")
	}

	// 使用辅助函数准备邮件消息
	emailMessage := m.PrepareMessage(message, protocol.MsgChannelEmail)
	if emailMessage == nil {
		return fmt.Errorf("failed to prepare email message: template not found for type %s", message.Type)
	}

	// 检查必要参数
	if to := cast.ToString(emailMessage.Params["to"]); to == "" {
		return fmt.Errorf("recipient ('to') is required for email")
	}

	// 标题和内容是可选的，但如果提供了应该非空
	if title := cast.ToString(emailMessage.Params["title"]); title == "" {
		emailMessage.Params["title"] = "No Subject"
	}

	if content := cast.ToString(emailMessage.Params["content"]); content == "" {
		emailMessage.Params["content"] = "No Content"
	}

	// 调用邮件服务发送
	if err := m.EmailService.SendMessage(emailMessage); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
