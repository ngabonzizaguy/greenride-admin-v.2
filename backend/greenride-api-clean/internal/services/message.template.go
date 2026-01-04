package services

import (
	"fmt"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"html/template"
	"os"
	"time"
)

var (
	templateService *MessageTemplateService
)
var (
	DefaultParams = map[string]any{
		"app_name": "Greenride",
		"website":  "https://www.greenride.com",
	}

	// 保留默认模板作为兼容性支持
	DefaultVerifyCodeTemplate = &models.MessageTemplate{
		Type:        protocol.MsgTypeVerifyCode,
		Channel:     protocol.MsgChannelEmail,
		Title:       "Verification Code for {{.app_name}}",
		Content:     "Your verification code is {{.code}}",
		Status:      protocol.StatusActive,
		Description: "Default verification code template",
	}
	DefaultRegisterEmailTemplate = &models.MessageTemplate{
		Type:        protocol.MsgTypeRegisterSuccess,
		Channel:     protocol.MsgChannelEmail,
		Title:       "Welcome to {{.app_name}}",
		Content:     "Thank you for registering with {{.app_name}}. Your account is now active.",
		Status:      protocol.StatusActive,
		Description: "Default registration email template",
	}
)

type MessageTemplateService struct {
	TemplateLib map[string]*protocol.MessageTemplate
	lastUpdate  int64
}

func SetupMessageTemplateService() {
	templateService = &MessageTemplateService{
		TemplateLib: make(map[string]*protocol.MessageTemplate),
	}
	if err := templateService.LoadTemplates(); err != nil {
		log.Get().Errorf("Failed to load message templates: %v", err)
	}
}

func GetMessageTemplateService() *MessageTemplateService {
	if templateService == nil {
		SetupMessageTemplateService()
	}
	return templateService
}

func (m *MessageTemplateService) GetTemplateByMessage(message *Message, channel string) *protocol.MessageTemplate {
	var matchTemplate *protocol.MessageTemplate
	var maxScore int
	for _, template := range m.TemplateLib {
		if template.Type != message.Type {
			continue
		}
		currentScore := 0
		if template.Channel != "" && template.Channel != protocol.Default && template.Channel == channel {
			currentScore++
		}
		if template.Language != "" && template.Language != protocol.Default && template.Language == message.Language {
			currentScore++
		}
		if template.Region != "" && template.Region != protocol.Default && template.Region == message.Region {
			currentScore++
		}
		if maxScore < currentScore {
			maxScore = currentScore
			matchTemplate = template
		}
	}
	return matchTemplate
}

func (m *MessageTemplateService) LoadTemplates() error {
	var templates []*models.MessageTemplate
	if err := models.DB.Find(&templates).Error; err != nil {
		log.Get().Errorf("Failed to load message templates from database: %v", err)
	}

	// 添加默认的Email模板
	templates = append(templates, DefaultEmailTemplates...)

	// 添加默认的SMS模板
	templates = append(templates, DefaultSmsTemplates...)

	// 添加默认的FCM模板
	templates = append(templates, DefaultFcmTemplates...)

	newTemplates := map[string]*protocol.MessageTemplate{}
	for _, t := range templates {
		pt, err := m.LoadTemplate(t)
		if err != nil {
			return err
		}

		// 构建模板键
		key := t.Type
		if t.Channel != "" {
			key += "_" + t.Channel
		} else {
			key += "_" + protocol.Default
		}
		if t.Language != "" {
			key += "_" + t.Language
		} else {
			key += "_" + protocol.Default
		}
		if t.Region != "" {
			key += "_" + t.Region
		} else {
			key += "_" + protocol.Default
		}

		newTemplates[key] = pt
	}

	m.TemplateLib = newTemplates
	m.lastUpdate = time.Now().UnixMilli()
	return nil
}

// parseTemplate 封装模板解析逻辑
func parseTemplate(content string) (*template.Template, error) {
	if content == "" {
		return nil, nil
	}
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		log.Get().Errorf("Failed to parse template: %v ,error: %v", content, err)
		return nil, err
	}
	return tmpl, nil
}

func (m *MessageTemplateService) LoadTemplate(mt *models.MessageTemplate) (*protocol.MessageTemplate, error) {
	titleTmpl, _ := parseTemplate(mt.Title)

	var contentTmpl *template.Template
	if mt.Content != "" {
		contentTmpl, _ = parseTemplate(mt.Content)
	} else if mt.Url != "" {
		content, _ := FetchTemplateContent(mt.Url)
		contentTmpl, _ = parseTemplate(content)
	}

	return &protocol.MessageTemplate{
		ID:          mt.ID,
		TemplateID:  mt.TemplateID,
		Type:        mt.Type,
		Channel:     mt.Channel,
		DeviceType:  mt.DeviceType,
		Platform:    mt.Platform,
		Language:    mt.Language,
		Region:      mt.Region,
		Tags:        mt.Tags,
		Title:       titleTmpl,
		Content:     contentTmpl,
		Status:      mt.Status,
		Description: mt.Description,
	}, nil
}

// RefreshIfNeeded 检查是否需要刷新模板缓存
func (m *MessageTemplateService) RefreshIfNeeded() error {
	var lastUpdated models.MessageTemplate
	err := models.DB.Select([]string{"updated_at"}).Order("updated_at desc").First(&lastUpdated).Error
	if err != nil {
		return err
	}

	if m.lastUpdate < lastUpdated.UpdatedAt {
		return m.LoadTemplates()
	}
	return nil
}

// FetchTemplateContent reads template content from a URL or file path
func FetchTemplateContent(url string) (string, error) {
	content, err := os.ReadFile(url)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	return string(content), nil
}
