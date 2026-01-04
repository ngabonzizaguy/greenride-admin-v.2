package services

import (
	"fmt"
	"greenride/internal/protocol"
)

// SendSmsMessage implements the SmsMessage interface
func (t *TwilioService) SendSmsMessage(message *Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	params := protocol.MapData{}
	params.Copy(DefaultParams)
	params.Copy(message.Params)

	// 获取接收人
	to, ok := params["to"]
	if !ok || to == nil {
		return fmt.Errorf("recipient ('to') is required for SMS")
	}
	recipient := fmt.Sprintf("%v", to)
	if recipient == "" {
		return fmt.Errorf("empty recipient for SMS")
	}

	// 检查是否是验证码消息类型
	if message.Type == protocol.MsgTypeVerifyCode {
		// 获取验证码
		codeVal, ok := params["code"]
		if !ok || codeVal == nil {
			return fmt.Errorf("verification code is required")
		}
		code := fmt.Sprintf("%v", codeVal)

		// 使用 Verify API 发送自定义验证码
		return t.SendCustomVerificationCode(recipient, code, message.Language)
	}
	// 发送普通短信
	contentVal, ok := params["content"]
	if !ok || contentVal == nil {
		return fmt.Errorf("content is required for SMS")
	}
	content := fmt.Sprintf("%v", contentVal)
	if content == "" {
		return fmt.Errorf("empty SMS content")
	}

	// 获取发送号码（如果指定）
	from := ""
	if fromVal, ok := params["from"]; ok && fromVal != nil {
		from = fmt.Sprintf("%v", fromVal)
	}

	return t.SendSMS(recipient, content, from)
}
