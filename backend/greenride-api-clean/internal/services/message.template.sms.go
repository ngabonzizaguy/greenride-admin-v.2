package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// SMS消息模板定义 - 系统自带模板
var (
	// 英文SMS模板
	DefaultVerifyCodeSmsEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangEnglish,
		Content:     "Your {{.app_name}} verification code is {{.code}}. Valid for 5 minutes.",
		Status:      protocol.StatusActive,
		Description: "SMS verification code template - English",
	}

	DefaultGenericSmsEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeGeneric,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangEnglish,
		Content:     "{{.content}}",
		Status:      protocol.StatusActive,
		Description: "Generic SMS template - English",
	}

	// 中文SMS模板
	DefaultVerifyCodeSmsZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangChinese,
		Content:     "【{{.app_name}}】您的验证码是{{.code}}，有效期为5分钟，请勿泄露给他人。",
		Status:      protocol.StatusActive,
		Description: "验证码短信模板 - 中文",
	}

	DefaultGenericSmsZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeGeneric,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangChinese,
		Content:     "{{.content}}",
		Status:      protocol.StatusActive,
		Description: "通用短信模板 - 中文",
	}

	// 法语SMS模板
	DefaultVerifyCodeSmsFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangFrench,
		Content:     "Votre code de vérification {{.app_name}} est {{.code}}. Valable pendant 5 minutes.",
		Status:      protocol.StatusActive,
		Description: "Modèle SMS de code de vérification - Français",
	}

	DefaultGenericSmsFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeGeneric,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangFrench,
		Content:     "{{.content}}",
		Status:      protocol.StatusActive,
		Description: "Modèle SMS générique - Français",
	}

	// 卢旺达语SMS模板
	DefaultVerifyCodeSmsRW = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangKinyarwanda,
		Content:     "Kode yanyu yo kugenzura {{.app_name}} ni {{.code}}. Izamara iminota 5.",
		Status:      protocol.StatusActive,
		Description: "SMS verification code template - Kinyarwanda",
	}

	DefaultGenericSmsRW = &models.MessageTemplate{
		Type:        protocol.MsgTypeGeneric,
		Channel:     protocol.MsgChannelSms,
		Language:    protocol.LangKinyarwanda,
		Content:     "{{.content}}",
		Status:      protocol.StatusActive,
		Description: "Generic SMS template - Kinyarwanda",
	}

	// 默认SMS模板集合
	DefaultSmsTemplates = []*models.MessageTemplate{
		// 英文模板
		DefaultVerifyCodeSmsEN,
		DefaultGenericSmsEN,

		// 中文模板
		DefaultVerifyCodeSmsZH,
		DefaultGenericSmsZH,

		// 法语模板
		DefaultVerifyCodeSmsFR,
		DefaultGenericSmsFR,

		// 卢旺达语模板
		DefaultVerifyCodeSmsRW,
		DefaultGenericSmsRW,
	}
)
