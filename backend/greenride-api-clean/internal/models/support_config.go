package models

import (
	"greenride/internal/utils"
)

// SupportConfig 客服支持配置表
type SupportConfig struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ConfigID string `json:"config_id" gorm:"column:config_id;type:varchar(64);uniqueIndex"`
	*SupportConfigValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type SupportConfigValues struct {
	// 联系方式
	SupportEmail   *string `json:"support_email" gorm:"column:support_email;type:varchar(255)"`
	SupportPhone   *string `json:"support_phone" gorm:"column:support_phone;type:varchar(50)"`
	SupportHours   *string `json:"support_hours" gorm:"column:support_hours;type:varchar(255)"`
	EmergencyPhone *string `json:"emergency_phone" gorm:"column:emergency_phone;type:varchar(50)"`
	WhatsAppNumber *string `json:"whatsapp_number" gorm:"column:whatsapp_number;type:varchar(50)"`

	// 响应时间配置（小时）
	ResponseTimeTarget *int `json:"response_time_target" gorm:"column:response_time_target;type:int;default:24"`

	// 自动回复配置
	AutoReplyEnabled *bool   `json:"auto_reply_enabled" gorm:"column:auto_reply_enabled;default:true"`
	AutoReplyMessage *string `json:"auto_reply_message" gorm:"column:auto_reply_message;type:text"`

	// 升级配置
	EscalationEnabled *bool `json:"escalation_enabled" gorm:"column:escalation_enabled;default:true"`
	EscalationTimeout *int  `json:"escalation_timeout" gorm:"column:escalation_timeout;type:int;default:48"` // 小时

	// 工作时间配置
	WorkdayStart *string `json:"workday_start" gorm:"column:workday_start;type:varchar(10)"` // 如 "08:00"
	WorkdayEnd   *string `json:"workday_end" gorm:"column:workday_end;type:varchar(10)"`     // 如 "18:00"
	WorkDays     *string `json:"work_days" gorm:"column:work_days;type:varchar(50)"`         // 如 "1,2,3,4,5" (周一到周五)

	// 通知配置
	NotifyOnNewFeedback   *bool `json:"notify_on_new_feedback" gorm:"column:notify_on_new_feedback;default:true"`
	NotifyOnHighPriority  *bool `json:"notify_on_high_priority" gorm:"column:notify_on_high_priority;default:true"`
	NotifyOnSafetyIssue   *bool `json:"notify_on_safety_issue" gorm:"column:notify_on_safety_issue;default:true"`
	NotifyOnEscalation    *bool `json:"notify_on_escalation" gorm:"column:notify_on_escalation;default:true"`

	// 管理员通知邮箱列表（JSON数组）
	NotificationEmails *string `json:"notification_emails" gorm:"column:notification_emails;type:json"`

	// FAQ链接
	FAQUrl      *string `json:"faq_url" gorm:"column:faq_url;type:varchar(500)"`
	HelpCenterUrl *string `json:"help_center_url" gorm:"column:help_center_url;type:varchar(500)"`

	// 元数据
	UpdatedBy *string `json:"updated_by" gorm:"column:updated_by;type:varchar(64)"`
	UpdatedAt int64   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (SupportConfig) TableName() string {
	return "t_support_config"
}

// NewSupportConfig 创建新的支持配置
func NewSupportConfig() *SupportConfig {
	return &SupportConfig{
		ConfigID: utils.GenerateID(),
		SupportConfigValues: &SupportConfigValues{
			SupportEmail:          utils.StringPtr("support@greenride.rw"),
			SupportPhone:          utils.StringPtr("+250 788 000 000"),
			SupportHours:          utils.StringPtr("Mon-Fri 8:00 AM - 6:00 PM"),
			EmergencyPhone:        utils.StringPtr("+250 788 000 001"),
			WhatsAppNumber:        utils.StringPtr("+250 788 000 000"),
			ResponseTimeTarget:    utils.IntPtr(24),
			AutoReplyEnabled:      utils.BoolPtr(true),
			EscalationEnabled:     utils.BoolPtr(true),
			EscalationTimeout:     utils.IntPtr(48),
			WorkdayStart:          utils.StringPtr("08:00"),
			WorkdayEnd:            utils.StringPtr("18:00"),
			WorkDays:              utils.StringPtr("1,2,3,4,5"),
			NotifyOnNewFeedback:   utils.BoolPtr(true),
			NotifyOnHighPriority:  utils.BoolPtr(true),
			NotifyOnSafetyIssue:   utils.BoolPtr(true),
			NotifyOnEscalation:    utils.BoolPtr(true),
		},
	}
}

// Getter methods
func (c *SupportConfigValues) GetSupportEmail() string {
	if c.SupportEmail == nil {
		return ""
	}
	return *c.SupportEmail
}

func (c *SupportConfigValues) GetSupportPhone() string {
	if c.SupportPhone == nil {
		return ""
	}
	return *c.SupportPhone
}

func (c *SupportConfigValues) GetSupportHours() string {
	if c.SupportHours == nil {
		return ""
	}
	return *c.SupportHours
}

func (c *SupportConfigValues) GetEmergencyPhone() string {
	if c.EmergencyPhone == nil {
		return ""
	}
	return *c.EmergencyPhone
}

func (c *SupportConfigValues) GetWhatsAppNumber() string {
	if c.WhatsAppNumber == nil {
		return ""
	}
	return *c.WhatsAppNumber
}

func (c *SupportConfigValues) GetResponseTimeTarget() int {
	if c.ResponseTimeTarget == nil {
		return 24
	}
	return *c.ResponseTimeTarget
}

func (c *SupportConfigValues) GetAutoReplyEnabled() bool {
	if c.AutoReplyEnabled == nil {
		return true
	}
	return *c.AutoReplyEnabled
}

func (c *SupportConfigValues) GetEscalationEnabled() bool {
	if c.EscalationEnabled == nil {
		return true
	}
	return *c.EscalationEnabled
}

func (c *SupportConfigValues) GetEscalationTimeout() int {
	if c.EscalationTimeout == nil {
		return 48
	}
	return *c.EscalationTimeout
}

// Setter methods
func (c *SupportConfigValues) SetSupportEmail(email string) *SupportConfigValues {
	c.SupportEmail = &email
	return c
}

func (c *SupportConfigValues) SetSupportPhone(phone string) *SupportConfigValues {
	c.SupportPhone = &phone
	return c
}

func (c *SupportConfigValues) SetSupportHours(hours string) *SupportConfigValues {
	c.SupportHours = &hours
	return c
}

func (c *SupportConfigValues) SetEmergencyPhone(phone string) *SupportConfigValues {
	c.EmergencyPhone = &phone
	return c
}

func (c *SupportConfigValues) SetWhatsAppNumber(number string) *SupportConfigValues {
	c.WhatsAppNumber = &number
	return c
}

func (c *SupportConfigValues) SetResponseTimeTarget(hours int) *SupportConfigValues {
	c.ResponseTimeTarget = &hours
	return c
}

func (c *SupportConfigValues) SetAutoReplyEnabled(enabled bool) *SupportConfigValues {
	c.AutoReplyEnabled = &enabled
	return c
}

func (c *SupportConfigValues) SetEscalationEnabled(enabled bool) *SupportConfigValues {
	c.EscalationEnabled = &enabled
	return c
}

func (c *SupportConfigValues) SetEscalationTimeout(hours int) *SupportConfigValues {
	c.EscalationTimeout = &hours
	return c
}

func (c *SupportConfigValues) SetUpdatedBy(adminID string) *SupportConfigValues {
	c.UpdatedBy = &adminID
	return c
}

