package services

import (
	"fmt"
	"greenride/internal/protocol"
)

// SendEmailMessage 实现EmailMessage接口中的SendEmailMessage方法
func (s *EmailAccountService) SendEmailMessage(message *Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	params := protocol.MapData{}
	params.Copy(DefaultParams)
	params.Copy(message.Params)

	// 获取接收人
	to, ok := params["to"]
	if !ok || to == nil {
		return fmt.Errorf("recipient ('to') is required for email")
	}
	recipient := fmt.Sprintf("%v", to)
	if recipient == "" {
		return fmt.Errorf("empty recipient for email")
	}

	// 获取标题
	titleVal, ok := params["title"]
	if !ok || titleVal == nil {
		return fmt.Errorf("title is required for email")
	}
	title := fmt.Sprintf("%v", titleVal)

	// 获取内容
	contentVal, ok := params["content"]
	if !ok || contentVal == nil {
		return fmt.Errorf("content is required for email")
	}
	content := fmt.Sprintf("%v", contentVal)

	// 发送邮件
	return s.SendEmail(recipient, title, content)
}
