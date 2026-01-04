package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// SupportService 支持配置服务
type SupportService struct{}

// GetConfig 获取支持配置
func (s *SupportService) GetConfig() (*protocol.SupportConfigResponse, error) {
	db := models.GetDB()
	if db == nil {
		return nil, nil
	}

	var config models.SupportConfig

	// Get the first (and should be only) config record
	if err := db.First(&config).Error; err != nil {
		// If not found, create default config
		config = *models.NewSupportConfig()
		if err := db.Create(&config).Error; err != nil {
			return nil, err
		}
	}

	return toSupportConfigResponse(&config), nil
}

// UpdateConfig 更新支持配置
func (s *SupportService) UpdateConfig(req *protocol.SupportConfigUpdateRequest, adminID string) error {
	db := models.GetDB()
	if db == nil {
		return nil
	}

	var config models.SupportConfig

	// Get existing config or create new
	if err := db.First(&config).Error; err != nil {
		config = *models.NewSupportConfig()
	}

	// Build updates
	updates := make(map[string]interface{})

	if req.SupportEmail != nil {
		updates["support_email"] = *req.SupportEmail
	}
	if req.SupportPhone != nil {
		updates["support_phone"] = *req.SupportPhone
	}
	if req.SupportHours != nil {
		updates["support_hours"] = *req.SupportHours
	}
	if req.EmergencyPhone != nil {
		updates["emergency_phone"] = *req.EmergencyPhone
	}
	if req.WhatsAppNumber != nil {
		updates["whatsapp_number"] = *req.WhatsAppNumber
	}
	if req.ResponseTimeTarget != nil {
		updates["response_time_target"] = *req.ResponseTimeTarget
	}
	if req.AutoReplyEnabled != nil {
		updates["auto_reply_enabled"] = *req.AutoReplyEnabled
	}
	if req.AutoReplyMessage != nil {
		updates["auto_reply_message"] = *req.AutoReplyMessage
	}
	if req.EscalationEnabled != nil {
		updates["escalation_enabled"] = *req.EscalationEnabled
	}
	if req.EscalationTimeout != nil {
		updates["escalation_timeout"] = *req.EscalationTimeout
	}
	if req.WorkdayStart != nil {
		updates["workday_start"] = *req.WorkdayStart
	}
	if req.WorkdayEnd != nil {
		updates["workday_end"] = *req.WorkdayEnd
	}
	if req.WorkDays != nil {
		updates["work_days"] = *req.WorkDays
	}
	if req.NotifyOnNewFeedback != nil {
		updates["notify_on_new_feedback"] = *req.NotifyOnNewFeedback
	}
	if req.NotifyOnHighPriority != nil {
		updates["notify_on_high_priority"] = *req.NotifyOnHighPriority
	}
	if req.NotifyOnSafetyIssue != nil {
		updates["notify_on_safety_issue"] = *req.NotifyOnSafetyIssue
	}
	if req.NotifyOnEscalation != nil {
		updates["notify_on_escalation"] = *req.NotifyOnEscalation
	}
	if req.FAQUrl != nil {
		updates["faq_url"] = *req.FAQUrl
	}
	if req.HelpCenterUrl != nil {
		updates["help_center_url"] = *req.HelpCenterUrl
	}

	updates["updated_by"] = adminID

	if config.ID == 0 {
		// Create new config
		for k, v := range updates {
			switch k {
			case "support_email":
				config.SupportEmail = v.(*string)
			case "support_phone":
				config.SupportPhone = v.(*string)
			// ... other fields will be set via GORM
			}
		}
		return db.Create(&config).Error
	}

	return db.Model(&config).Updates(updates).Error
}

// toSupportConfigResponse converts model to response
func toSupportConfigResponse(config *models.SupportConfig) *protocol.SupportConfigResponse {
	return &protocol.SupportConfigResponse{
		SupportEmail:         config.GetSupportEmail(),
		SupportPhone:         config.GetSupportPhone(),
		SupportHours:         config.GetSupportHours(),
		EmergencyPhone:       config.GetEmergencyPhone(),
		WhatsAppNumber:       config.GetWhatsAppNumber(),
		ResponseTimeTarget:   config.GetResponseTimeTarget(),
		AutoReplyEnabled:     config.GetAutoReplyEnabled(),
		AutoReplyMessage:     getStringPtrValue(config.AutoReplyMessage),
		EscalationEnabled:    config.GetEscalationEnabled(),
		EscalationTimeout:    config.GetEscalationTimeout(),
		WorkdayStart:         getStringPtrValue(config.WorkdayStart),
		WorkdayEnd:           getStringPtrValue(config.WorkdayEnd),
		WorkDays:             getStringPtrValue(config.WorkDays),
		NotifyOnNewFeedback:  getBoolPtrValue(config.NotifyOnNewFeedback),
		NotifyOnHighPriority: getBoolPtrValue(config.NotifyOnHighPriority),
		NotifyOnSafetyIssue:  getBoolPtrValue(config.NotifyOnSafetyIssue),
		NotifyOnEscalation:   getBoolPtrValue(config.NotifyOnEscalation),
		FAQUrl:               getStringPtrValue(config.FAQUrl),
		HelpCenterUrl:        getStringPtrValue(config.HelpCenterUrl),
		UpdatedAt:            config.UpdatedAt,
	}
}

func getStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getBoolPtrValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// GetSupportService 获取支持服务实例
func GetSupportService() *SupportService {
	return &SupportService{}
}

