package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// Email消息模板定义 - 系统自带模板
var (
	// 英文Email模板
	DefaultVerifyCodeEmailEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangEnglish,
		Title:       "Verification Code for {{.app_name}}",
		Content:     "Your verification code is {{.code}}. It will be valid for 5 minutes.",
		Status:      protocol.StatusActive,
		Description: "Verification code template - English",
	}

	DefaultRegisterSuccessEmailEN = &models.MessageTemplate{
		Type:        protocol.MsgTypeRegisterSuccess,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangEnglish,
		Title:       "Welcome to {{.app_name}}",
		Content:     "Thank you for registering with {{.app_name}}. Your account is now active.",
		Status:      protocol.StatusActive,
		Description: "Registration success email template - English",
	}

	// 中文Email模板
	DefaultVerifyCodeEmailZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangChinese,
		Title:       "{{.app_name}} 验证码",
		Content:     "您的验证码是 {{.code}}，有效期为5分钟。",
		Status:      protocol.StatusActive,
		Description: "验证码邮件模板 - 中文",
	}

	DefaultRegisterSuccessEmailZH = &models.MessageTemplate{
		Type:        protocol.MsgTypeRegisterSuccess,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangChinese,
		Title:       "欢迎加入 {{.app_name}}",
		Content:     "感谢您注册 {{.app_name}}。您的账户现已激活。",
		Status:      protocol.StatusActive,
		Description: "注册成功邮件模板 - 中文",
	}

	// 法语Email模板
	DefaultVerifyCodeEmailFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangFrench,
		Title:       "Code de vérification pour {{.app_name}}",
		Content:     "Votre code de vérification est {{.code}}. Il sera valide pendant 5 minutes.",
		Status:      protocol.StatusActive,
		Description: "Modèle de code de vérification - Français",
	}

	DefaultRegisterSuccessEmailFR = &models.MessageTemplate{
		Type:        protocol.MsgTypeRegisterSuccess,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangFrench,
		Title:       "Bienvenue sur {{.app_name}}",
		Content:     "Merci de vous être inscrit sur {{.app_name}}. Votre compte est maintenant actif.",
		Status:      protocol.StatusActive,
		Description: "Modèle d'e-mail de réussite d'inscription - Français",
	}

	// 卢旺达语Email模板
	DefaultVerifyCodeEmailRW = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangKinyarwanda,
		Title:       "Kode yo kugenzura kuri {{.app_name}}",
		Content:     "Kode yanyu yo kugenzura ni {{.code}}. Izamara iminota 5.",
		Status:      protocol.StatusActive,
		Description: "Verification code template - Kinyarwanda",
	}

	DefaultRegisterSuccessEmailRW = &models.MessageTemplate{
		Type:        protocol.MsgTypeRegisterSuccess,
		Channel:     protocol.MsgChannelEmail,
		Language:    protocol.LangKinyarwanda,
		Title:       "Murakaza neza kuri {{.app_name}}",
		Content:     "Urakoze kwiyandikisha kuri {{.app_name}}. Konti yawe ubu irakora.",
		Status:      protocol.StatusActive,
		Description: "Registration success email template - Kinyarwanda",
	}

	// 默认Email模板集合
	DefaultEmailTemplates = []*models.MessageTemplate{
		// 英文模板
		DefaultVerifyCodeEmailEN,
		DefaultRegisterSuccessEmailEN,

		// 中文模板
		DefaultVerifyCodeEmailZH,
		DefaultRegisterSuccessEmailZH,

		// 法语模板
		DefaultVerifyCodeEmailFR,
		DefaultRegisterSuccessEmailFR,

		// 卢旺达语模板
		DefaultVerifyCodeEmailRW,
		DefaultRegisterSuccessEmailRW,
	}
)
