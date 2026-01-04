package models

import (
	"greenride/internal/utils"
	"strings"
)

// Announcement 公告表 - 系统公告和通知管理
type Announcement struct {
	ID             int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AnnouncementID string `json:"announcement_id" gorm:"column:announcement_id;type:varchar(64);uniqueIndex"`
	Salt           string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*AnnouncementValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type AnnouncementValues struct {
	// 基本信息
	Title    *string `json:"title" gorm:"column:title;type:varchar(255)"`            // 公告标题
	Content  *string `json:"content" gorm:"column:content;type:text"`                // 公告内容
	Summary  *string `json:"summary" gorm:"column:summary;type:varchar(500)"`        // 摘要
	Type     *string `json:"type" gorm:"column:type;type:varchar(30);index"`         // system, promotion, maintenance, emergency, feature, security
	Category *string `json:"category" gorm:"column:category;type:varchar(50);index"` // 分类
	Level    *string `json:"level" gorm:"column:level;type:varchar(20);index"`       // info, warning, critical, urgent

	// 目标受众
	TargetAudience    *string `json:"target_audience" gorm:"column:target_audience;type:varchar(30)"`          // all, users, drivers, admins, vip, new_users
	TargetUserSegment *string `json:"target_user_segment" gorm:"column:target_user_segment;type:varchar(100)"` // 用户细分
	TargetUserLevel   *int    `json:"target_user_level" gorm:"column:target_user_level;type:int"`              // 目标用户等级
	TargetUserIDs     *string `json:"target_user_ids" gorm:"column:target_user_ids;type:text"`                 // 特定用户ID列表
	ExcludedUserIDs   *string `json:"excluded_user_ids" gorm:"column:excluded_user_ids;type:text"`             // 排除用户ID列表

	// 地理限制
	TargetCities       *string `json:"target_cities" gorm:"column:target_cities;type:varchar(500)"`               // 目标城市
	TargetRegions      *string `json:"target_regions" gorm:"column:target_regions;type:varchar(500)"`             // 目标地区
	TargetServiceAreas *string `json:"target_service_areas" gorm:"column:target_service_areas;type:varchar(500)"` // 目标服务区域
	ExcludedAreas      *string `json:"excluded_areas" gorm:"column:excluded_areas;type:varchar(500)"`             // 排除区域

	// 时间设置
	PublishTime  *int64 `json:"publish_time" gorm:"column:publish_time"`   // 发布时间
	StartTime    *int64 `json:"start_time" gorm:"column:start_time"`       // 生效时间
	EndTime      *int64 `json:"end_time" gorm:"column:end_time"`           // 失效时间
	ExpiresAt    *int64 `json:"expires_at" gorm:"column:expires_at"`       // 过期时间
	ReminderTime *int64 `json:"reminder_time" gorm:"column:reminder_time"` // 提醒时间

	// 显示设置
	IsSticky       *bool `json:"is_sticky" gorm:"column:is_sticky;default:false"`             // 是否置顶
	IsBanner       *bool `json:"is_banner" gorm:"column:is_banner;default:false"`             // 是否横幅显示
	IsPopup        *bool `json:"is_popup" gorm:"column:is_popup;default:false"`               // 是否弹窗显示
	IsFullScreen   *bool `json:"is_fullscreen" gorm:"column:is_fullscreen;default:false"`     // 是否全屏显示
	ShowOnce       *bool `json:"show_once" gorm:"column:show_once;default:false"`             // 是否只显示一次
	RequireConfirm *bool `json:"require_confirm" gorm:"column:require_confirm;default:false"` // 是否需要确认
	ForceRead      *bool `json:"force_read" gorm:"column:force_read;default:false"`           // 是否强制阅读

	// 显示位置
	DisplayPosition *string `json:"display_position" gorm:"column:display_position;type:varchar(50)"` // home, profile, ride, payment, notification
	DisplayOrder    *int    `json:"display_order" gorm:"column:display_order;type:int;default:0"`     // 显示顺序
	Priority        *int    `json:"priority" gorm:"column:priority;type:int;default:0"`               // 优先级

	// 样式设置
	BackgroundColor *string `json:"background_color" gorm:"column:background_color;type:varchar(20)"` // 背景颜色
	TextColor       *string `json:"text_color" gorm:"column:text_color;type:varchar(20)"`             // 文字颜色
	IconURL         *string `json:"icon_url" gorm:"column:icon_url;type:varchar(500)"`                // 图标URL
	ImageURL        *string `json:"image_url" gorm:"column:image_url;type:varchar(500)"`              // 图片URL
	BannerURL       *string `json:"banner_url" gorm:"column:banner_url;type:varchar(500)"`            // 横幅URL
	VideoURL        *string `json:"video_url" gorm:"column:video_url;type:varchar(500)"`              // 视频URL

	// 交互设置
	ActionType      *string `json:"action_type" gorm:"column:action_type;type:varchar(30)"`        // none, url, deeplink, share, feedback
	ActionURL       *string `json:"action_url" gorm:"column:action_url;type:varchar(500)"`         // 动作链接
	ActionText      *string `json:"action_text" gorm:"column:action_text;type:varchar(100)"`       // 动作按钮文本
	ShareURL        *string `json:"share_url" gorm:"column:share_url;type:varchar(500)"`           // 分享链接
	ShareText       *string `json:"share_text" gorm:"column:share_text;type:text"`                 // 分享文案
	FeedbackEnabled *bool   `json:"feedback_enabled" gorm:"column:feedback_enabled;default:false"` // 是否启用反馈

	// 推送设置
	SendPush      *bool   `json:"send_push" gorm:"column:send_push;default:false"`                     // 是否发送推送
	PushTitle     *string `json:"push_title" gorm:"column:push_title;type:varchar(255)"`               // 推送标题
	PushContent   *string `json:"push_content" gorm:"column:push_content;type:text"`                   // 推送内容
	PushTime      *int64  `json:"push_time" gorm:"column:push_time"`                                   // 推送时间
	PushDelay     *int    `json:"push_delay" gorm:"column:push_delay;type:int;default:0"`              // 推送延迟（分钟）
	PushBatchSize *int    `json:"push_batch_size" gorm:"column:push_batch_size;type:int;default:1000"` // 推送批次大小

	// 邮件设置
	SendEmail     *bool   `json:"send_email" gorm:"column:send_email;default:false"`             // 是否发送邮件
	EmailSubject  *string `json:"email_subject" gorm:"column:email_subject;type:varchar(255)"`   // 邮件主题
	EmailTemplate *string `json:"email_template" gorm:"column:email_template;type:varchar(100)"` // 邮件模板
	EmailTime     *int64  `json:"email_time" gorm:"column:email_time"`                           // 邮件发送时间

	// 短信设置
	SendSMS     *bool   `json:"send_sms" gorm:"column:send_sms;default:false"`             // 是否发送短信
	SMSContent  *string `json:"sms_content" gorm:"column:sms_content;type:varchar(500)"`   // 短信内容
	SMSTemplate *string `json:"sms_template" gorm:"column:sms_template;type:varchar(100)"` // 短信模板
	SMSTime     *int64  `json:"sms_time" gorm:"column:sms_time"`                           // 短信发送时间

	// 状态管理
	Status      *string `json:"status" gorm:"column:status;type:varchar(30);index;default:'draft'"` // draft, published, scheduled, paused, expired, deleted
	IsActive    *bool   `json:"is_active" gorm:"column:is_active;default:false"`                    // 是否激活
	IsPublished *bool   `json:"is_published" gorm:"column:is_published;default:false"`              // 是否已发布
	IsScheduled *bool   `json:"is_scheduled" gorm:"column:is_scheduled;default:false"`              // 是否定时发布
	IsUrgent    *bool   `json:"is_urgent" gorm:"column:is_urgent;default:false"`                    // 是否紧急

	// 渠道限制
	ValidChannels    *string `json:"valid_channels" gorm:"column:valid_channels;type:varchar(200)"`         // app, web, sms, email, push
	ValidPlatforms   *string `json:"valid_platforms" gorm:"column:valid_platforms;type:varchar(100)"`       // ios, android, web
	ValidAppVersions *string `json:"valid_app_versions" gorm:"column:valid_app_versions;type:varchar(200)"` // 有效应用版本
	MinAppVersion    *string `json:"min_app_version" gorm:"column:min_app_version;type:varchar(20)"`        // 最低应用版本

	// 统计信息
	ViewCount      *int `json:"view_count" gorm:"column:view_count;type:int;default:0"`             // 查看次数
	ClickCount     *int `json:"click_count" gorm:"column:click_count;type:int;default:0"`           // 点击次数
	ShareCount     *int `json:"share_count" gorm:"column:share_count;type:int;default:0"`           // 分享次数
	FeedbackCount  *int `json:"feedback_count" gorm:"column:feedback_count;type:int;default:0"`     // 反馈次数
	PushSentCount  *int `json:"push_sent_count" gorm:"column:push_sent_count;type:int;default:0"`   // 推送发送数量
	PushOpenCount  *int `json:"push_open_count" gorm:"column:push_open_count;type:int;default:0"`   // 推送打开数量
	EmailSentCount *int `json:"email_sent_count" gorm:"column:email_sent_count;type:int;default:0"` // 邮件发送数量
	EmailOpenCount *int `json:"email_open_count" gorm:"column:email_open_count;type:int;default:0"` // 邮件打开数量
	SMSSentCount   *int `json:"sms_sent_count" gorm:"column:sms_sent_count;type:int;default:0"`     // 短信发送数量

	// 成功率统计
	ClickRate      *float64 `json:"click_rate" gorm:"column:click_rate;type:decimal(5,2);default:0.00"`           // 点击率
	ShareRate      *float64 `json:"share_rate" gorm:"column:share_rate;type:decimal(5,2);default:0.00"`           // 分享率
	PushOpenRate   *float64 `json:"push_open_rate" gorm:"column:push_open_rate;type:decimal(5,2);default:0.00"`   // 推送打开率
	EmailOpenRate  *float64 `json:"email_open_rate" gorm:"column:email_open_rate;type:decimal(5,2);default:0.00"` // 邮件打开率
	EngagementRate *float64 `json:"engagement_rate" gorm:"column:engagement_rate;type:decimal(5,2);default:0.00"` // 参与率

	// A/B测试
	ExperimentID *string `json:"experiment_id" gorm:"column:experiment_id;type:varchar(100)"` // 实验ID
	VariantID    *string `json:"variant_id" gorm:"column:variant_id;type:varchar(100)"`       // 变体ID
	TestGroup    *string `json:"test_group" gorm:"column:test_group;type:varchar(50)"`        // 测试组
	ControlGroup *bool   `json:"control_group" gorm:"column:control_group;default:false"`     // 是否对照组

	// 审批流程
	ApprovalStatus *string `json:"approval_status" gorm:"column:approval_status;type:varchar(30);default:'pending'"` // pending, approved, rejected
	ApprovedBy     *string `json:"approved_by" gorm:"column:approved_by;type:varchar(100)"`                          // 审批人
	ApprovedAt     *int64  `json:"approved_at" gorm:"column:approved_at"`                                            // 审批时间
	ApprovalNotes  *string `json:"approval_notes" gorm:"column:approval_notes;type:text"`                            // 审批备注

	// 创建者信息
	CreatedBy    *string `json:"created_by" gorm:"column:created_by;type:varchar(100)"`       // 创建者
	CreatorType  *string `json:"creator_type" gorm:"column:creator_type;type:varchar(30)"`    // admin, system, auto
	CreatorID    *string `json:"creator_id" gorm:"column:creator_id;type:varchar(100)"`       // 创建者ID
	DepartmentID *string `json:"department_id" gorm:"column:department_id;type:varchar(100)"` // 部门ID

	// 多语言支持
	Language     *string `json:"language" gorm:"column:language;type:varchar(10);default:'en'"` // 语言
	Translations *string `json:"translations" gorm:"column:translations;type:json"`             // 翻译内容
	IsTranslated *bool   `json:"is_translated" gorm:"column:is_translated;default:false"`       // 是否已翻译

	// 时间记录
	PublishedAt   *int64 `json:"published_at" gorm:"column:published_at"`       // 发布时间
	LastViewedAt  *int64 `json:"last_viewed_at" gorm:"column:last_viewed_at"`   // 最后查看时间
	LastClickedAt *int64 `json:"last_clicked_at" gorm:"column:last_clicked_at"` // 最后点击时间
	LastSharedAt  *int64 `json:"last_shared_at" gorm:"column:last_shared_at"`   // 最后分享时间
	ActivatedAt   *int64 `json:"activated_at" gorm:"column:activated_at"`       // 激活时间
	PausedAt      *int64 `json:"paused_at" gorm:"column:paused_at"`             // 暂停时间

	// 关联信息
	RelatedPromoID   *string `json:"related_promo_id" gorm:"column:related_promo_id;type:varchar(100)"`     // 关联促销ID
	RelatedFeatureID *string `json:"related_feature_id" gorm:"column:related_feature_id;type:varchar(100)"` // 关联功能ID
	RelatedURL       *string `json:"related_url" gorm:"column:related_url;type:varchar(500)"`               // 关联链接

	// 扩展信息
	Tags          *string `json:"tags" gorm:"column:tags;type:varchar(500)"`             // 标签
	Keywords      *string `json:"keywords" gorm:"column:keywords;type:varchar(500)"`     // 关键词
	CustomFields  *string `json:"custom_fields" gorm:"column:custom_fields;type:json"`   // 自定义字段
	InternalNotes *string `json:"internal_notes" gorm:"column:internal_notes;type:text"` // 内部备注
	ExternalNotes *string `json:"external_notes" gorm:"column:external_notes;type:text"` // 外部备注
	Metadata      *string `json:"metadata" gorm:"column:metadata;type:json"`             // 元数据

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Announcement) TableName() string {
	return "t_announcements"
}

// 公告类型常量
const (
	AnnouncementTypeSystem      = "system"      // 系统公告
	AnnouncementTypePromotion   = "promotion"   // 促销公告
	AnnouncementTypeMaintenance = "maintenance" // 维护公告
	AnnouncementTypeEmergency   = "emergency"   // 紧急公告
	AnnouncementTypeFeature     = "feature"     // 功能公告
	AnnouncementTypeSecurity    = "security"    // 安全公告
)

// 公告级别常量
const (
	AnnouncementLevelInfo     = "info"     // 信息
	AnnouncementLevelWarning  = "warning"  // 警告
	AnnouncementLevelCritical = "critical" // 严重
	AnnouncementLevelUrgent   = "urgent"   // 紧急
)

// 目标受众常量
const (
	TargetAudienceAll      = "all"       // 所有用户
	TargetAudienceUsers    = "users"     // 普通用户
	TargetAudienceDrivers  = "drivers"   // 司机
	TargetAudienceAdmins   = "admins"    // 管理员
	TargetAudienceVIP      = "vip"       // VIP用户
	TargetAudienceNewUsers = "new_users" // 新用户
)

// 状态常量
const (
	AnnouncementStatusDraft     = "draft"     // 草稿
	AnnouncementStatusPublished = "published" // 已发布
	AnnouncementStatusScheduled = "scheduled" // 定时发布
	AnnouncementStatusPaused    = "paused"    // 暂停
	AnnouncementStatusExpired   = "expired"   // 已过期
	AnnouncementStatusDeleted   = "deleted"   // 已删除
)

// 动作类型常量
const (
	ActionTypeNone     = "none"     // 无动作
	ActionTypeURL      = "url"      // 链接
	ActionTypeDeeplink = "deeplink" // 深度链接
	ActionTypeShare    = "share"    // 分享
	ActionTypeFeedback = "feedback" // 反馈
)

// 审批状态常量
const (
	AnnouncementApprovalStatusPending  = "pending"  // 待审批
	AnnouncementApprovalStatusApproved = "approved" // 已审批
	AnnouncementApprovalStatusRejected = "rejected" // 已拒绝
)

// 创建新的公告对象
func NewAnnouncementV2() *Announcement {
	return &Announcement{
		AnnouncementID: utils.GenerateAnnouncementID(),
		Salt:           utils.GenerateSalt(),
		AnnouncementValues: &AnnouncementValues{
			Type:            utils.StringPtr(AnnouncementTypeSystem),
			Level:           utils.StringPtr(AnnouncementLevelInfo),
			TargetAudience:  utils.StringPtr(TargetAudienceAll),
			IsSticky:        utils.BoolPtr(false),
			IsBanner:        utils.BoolPtr(false),
			IsPopup:         utils.BoolPtr(false),
			IsFullScreen:    utils.BoolPtr(false),
			ShowOnce:        utils.BoolPtr(false),
			RequireConfirm:  utils.BoolPtr(false),
			ForceRead:       utils.BoolPtr(false),
			DisplayOrder:    utils.IntPtr(0),
			Priority:        utils.IntPtr(0),
			ActionType:      utils.StringPtr(ActionTypeNone),
			FeedbackEnabled: utils.BoolPtr(false),
			SendPush:        utils.BoolPtr(false),
			PushDelay:       utils.IntPtr(0),
			PushBatchSize:   utils.IntPtr(1000),
			SendEmail:       utils.BoolPtr(false),
			SendSMS:         utils.BoolPtr(false),
			Status:          utils.StringPtr(AnnouncementStatusDraft),
			IsActive:        utils.BoolPtr(false),
			IsPublished:     utils.BoolPtr(false),
			IsScheduled:     utils.BoolPtr(false),
			IsUrgent:        utils.BoolPtr(false),
			ViewCount:       utils.IntPtr(0),
			ClickCount:      utils.IntPtr(0),
			ShareCount:      utils.IntPtr(0),
			FeedbackCount:   utils.IntPtr(0),
			PushSentCount:   utils.IntPtr(0),
			PushOpenCount:   utils.IntPtr(0),
			EmailSentCount:  utils.IntPtr(0),
			EmailOpenCount:  utils.IntPtr(0),
			SMSSentCount:    utils.IntPtr(0),
			ClickRate:       utils.Float64Ptr(0.00),
			ShareRate:       utils.Float64Ptr(0.00),
			PushOpenRate:    utils.Float64Ptr(0.00),
			EmailOpenRate:   utils.Float64Ptr(0.00),
			EngagementRate:  utils.Float64Ptr(0.00),
			ControlGroup:    utils.BoolPtr(false),
			ApprovalStatus:  utils.StringPtr(AnnouncementApprovalStatusPending),
			Language:        utils.StringPtr("en"),
			IsTranslated:    utils.BoolPtr(false),
		},
	}
}

// SetValues 更新AnnouncementV2Values中的非nil值
func (a *AnnouncementValues) SetValues(values *AnnouncementValues) {
	if values == nil {
		return
	}

	if values.Title != nil {
		a.Title = values.Title
	}
	if values.Content != nil {
		a.Content = values.Content
	}
	if values.Type != nil {
		a.Type = values.Type
	}
	if values.Level != nil {
		a.Level = values.Level
	}
	if values.TargetAudience != nil {
		a.TargetAudience = values.TargetAudience
	}
	if values.Status != nil {
		a.Status = values.Status
	}
	if values.PublishTime != nil {
		a.PublishTime = values.PublishTime
	}
	if values.StartTime != nil {
		a.StartTime = values.StartTime
	}
	if values.EndTime != nil {
		a.EndTime = values.EndTime
	}
	if values.IsActive != nil {
		a.IsActive = values.IsActive
	}
	if values.UpdatedAt > 0 {
		a.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (a *AnnouncementValues) GetTitle() string {
	if a.Title == nil {
		return ""
	}
	return *a.Title
}

func (a *AnnouncementValues) GetContent() string {
	if a.Content == nil {
		return ""
	}
	return *a.Content
}

func (a *AnnouncementValues) GetSummary() string {
	if a.Summary == nil {
		return ""
	}
	return *a.Summary
}

func (a *AnnouncementValues) GetType() string {
	if a.Type == nil {
		return AnnouncementTypeSystem
	}
	return *a.Type
}

func (a *AnnouncementValues) GetLevel() string {
	if a.Level == nil {
		return AnnouncementLevelInfo
	}
	return *a.Level
}

func (a *AnnouncementValues) GetTargetAudience() string {
	if a.TargetAudience == nil {
		return TargetAudienceAll
	}
	return *a.TargetAudience
}

func (a *AnnouncementValues) GetStatus() string {
	if a.Status == nil {
		return AnnouncementStatusDraft
	}
	return *a.Status
}

func (a *AnnouncementValues) GetIsActive() bool {
	if a.IsActive == nil {
		return false
	}
	return *a.IsActive
}

func (a *AnnouncementValues) GetIsPublished() bool {
	if a.IsPublished == nil {
		return false
	}
	return *a.IsPublished
}

func (a *AnnouncementValues) GetIsSticky() bool {
	if a.IsSticky == nil {
		return false
	}
	return *a.IsSticky
}

func (a *AnnouncementValues) GetIsBanner() bool {
	if a.IsBanner == nil {
		return false
	}
	return *a.IsBanner
}

func (a *AnnouncementValues) GetIsPopup() bool {
	if a.IsPopup == nil {
		return false
	}
	return *a.IsPopup
}

func (a *AnnouncementValues) GetRequireConfirm() bool {
	if a.RequireConfirm == nil {
		return false
	}
	return *a.RequireConfirm
}

func (a *AnnouncementValues) GetIsUrgent() bool {
	if a.IsUrgent == nil {
		return false
	}
	return *a.IsUrgent
}

func (a *AnnouncementValues) GetPriority() int {
	if a.Priority == nil {
		return 0
	}
	return *a.Priority
}

func (a *AnnouncementValues) GetViewCount() int {
	if a.ViewCount == nil {
		return 0
	}
	return *a.ViewCount
}

func (a *AnnouncementValues) GetClickCount() int {
	if a.ClickCount == nil {
		return 0
	}
	return *a.ClickCount
}

func (a *AnnouncementValues) GetLanguage() string {
	if a.Language == nil {
		return "en"
	}
	return *a.Language
}

func (a *AnnouncementValues) GetApprovalStatus() string {
	if a.ApprovalStatus == nil {
		return AnnouncementApprovalStatusPending
	}
	return *a.ApprovalStatus
}

// Setter 方法
func (a *AnnouncementValues) SetTitle(title string) *AnnouncementValues {
	a.Title = &title
	return a
}

func (a *AnnouncementValues) SetContent(content string) *AnnouncementValues {
	a.Content = &content
	return a
}

func (a *AnnouncementValues) SetSummary(summary string) *AnnouncementValues {
	a.Summary = &summary
	return a
}

func (a *AnnouncementValues) SetType(announcementType string) *AnnouncementValues {
	a.Type = &announcementType
	return a
}

func (a *AnnouncementValues) SetLevel(level string) *AnnouncementValues {
	a.Level = &level
	return a
}

func (a *AnnouncementValues) SetTargetAudience(audience string) *AnnouncementValues {
	a.TargetAudience = &audience
	return a
}

func (a *AnnouncementValues) SetStatus(status string) *AnnouncementValues {
	a.Status = &status
	return a
}

func (a *AnnouncementValues) SetTimeRange(publishTime, startTime, endTime int64) *AnnouncementValues {
	a.PublishTime = &publishTime
	a.StartTime = &startTime
	a.EndTime = &endTime
	return a
}

func (a *AnnouncementValues) SetDisplaySettings(isSticky, isBanner, isPopup bool) *AnnouncementValues {
	a.IsSticky = &isSticky
	a.IsBanner = &isBanner
	a.IsPopup = &isPopup
	return a
}

func (a *AnnouncementValues) SetInteractionSettings(requireConfirm, forceRead bool) *AnnouncementValues {
	a.RequireConfirm = &requireConfirm
	a.ForceRead = &forceRead
	return a
}

func (a *AnnouncementValues) SetPushSettings(sendPush bool, pushTitle, pushContent string, pushTime int64) *AnnouncementValues {
	a.SendPush = &sendPush
	a.PushTitle = &pushTitle
	a.PushContent = &pushContent
	a.PushTime = &pushTime
	return a
}

func (a *AnnouncementValues) SetEmailSettings(sendEmail bool, emailSubject, emailTemplate string) *AnnouncementValues {
	a.SendEmail = &sendEmail
	a.EmailSubject = &emailSubject
	a.EmailTemplate = &emailTemplate
	return a
}

func (a *AnnouncementValues) SetAction(actionType, actionURL, actionText string) *AnnouncementValues {
	a.ActionType = &actionType
	a.ActionURL = &actionURL
	a.ActionText = &actionText
	return a
}

func (a *AnnouncementValues) SetTargetUsers(audience, segment string, level int) *AnnouncementValues {
	a.TargetAudience = &audience
	a.TargetUserSegment = &segment
	a.TargetUserLevel = &level
	return a
}

func (a *AnnouncementValues) SetGeographicTarget(cities, regions, serviceAreas string) *AnnouncementValues {
	a.TargetCities = &cities
	a.TargetRegions = &regions
	a.TargetServiceAreas = &serviceAreas
	return a
}

func (a *AnnouncementValues) SetChannelLimits(channels, platforms string) *AnnouncementValues {
	a.ValidChannels = &channels
	a.ValidPlatforms = &platforms
	return a
}

func (a *AnnouncementValues) SetCreator(createdBy, creatorType, creatorID string) *AnnouncementValues {
	a.CreatedBy = &createdBy
	a.CreatorType = &creatorType
	a.CreatorID = &creatorID
	return a
}

func (a *AnnouncementValues) SetLanguage(language string) *AnnouncementValues {
	a.Language = &language
	return a
}

func (a *AnnouncementValues) SetUrgent(isUrgent bool) *AnnouncementValues {
	a.IsUrgent = &isUrgent
	return a
}

func (a *AnnouncementValues) SetPriority(priority int) *AnnouncementValues {
	a.Priority = &priority
	return a
}

// 业务方法
func (a *Announcement) IsActive() bool {
	return a.GetStatus() == AnnouncementStatusPublished && a.GetIsActive()
}

func (a *Announcement) IsPublished() bool {
	return a.GetIsPublished() && a.GetStatus() == AnnouncementStatusPublished
}

func (a *Announcement) IsScheduled() bool {
	return a.GetStatus() == AnnouncementStatusScheduled
}

func (a *Announcement) IsDraft() bool {
	return a.GetStatus() == AnnouncementStatusDraft
}

func (a *Announcement) IsExpired() bool {
	if a.AnnouncementValues.EndTime == nil {
		return false
	}
	return *a.AnnouncementValues.EndTime < utils.TimeNowMilli()
}

func (a *Announcement) IsTimeValid() bool {
	now := utils.TimeNowMilli()

	// 检查开始时间
	if a.AnnouncementValues.StartTime != nil && now < *a.AnnouncementValues.StartTime {
		return false
	}

	// 检查结束时间
	if a.AnnouncementValues.EndTime != nil && now > *a.AnnouncementValues.EndTime {
		return false
	}

	return true
}

func (a *Announcement) IsVisible() bool {
	return a.IsActive() && a.IsTimeValid() && !a.IsExpired()
}

func (a *Announcement) IsSystemAnnouncement() bool {
	return a.GetType() == AnnouncementTypeSystem
}

func (a *Announcement) IsPromotionAnnouncement() bool {
	return a.GetType() == AnnouncementTypePromotion
}

func (a *Announcement) IsMaintenanceAnnouncement() bool {
	return a.GetType() == AnnouncementTypeMaintenance
}

func (a *Announcement) IsEmergencyAnnouncement() bool {
	return a.GetType() == AnnouncementTypeEmergency
}

func (a *Announcement) IsFeatureAnnouncement() bool {
	return a.GetType() == AnnouncementTypeFeature
}

func (a *Announcement) IsSecurityAnnouncement() bool {
	return a.GetType() == AnnouncementTypeSecurity
}

func (a *Announcement) IsInfoLevel() bool {
	return a.GetLevel() == AnnouncementLevelInfo
}

func (a *Announcement) IsWarningLevel() bool {
	return a.GetLevel() == AnnouncementLevelWarning
}

func (a *Announcement) IsCriticalLevel() bool {
	return a.GetLevel() == AnnouncementLevelCritical
}

func (a *Announcement) IsUrgentLevel() bool {
	return a.GetLevel() == AnnouncementLevelUrgent
}

func (a *Announcement) IsForAllUsers() bool {
	return a.GetTargetAudience() == TargetAudienceAll
}

func (a *Announcement) IsForUsers() bool {
	return a.GetTargetAudience() == TargetAudienceUsers
}

func (a *Announcement) IsForDrivers() bool {
	return a.GetTargetAudience() == TargetAudienceDrivers
}

func (a *Announcement) IsForVIPUsers() bool {
	return a.GetTargetAudience() == TargetAudienceVIP
}

func (a *Announcement) IsForNewUsers() bool {
	return a.GetTargetAudience() == TargetAudienceNewUsers
}

func (a *Announcement) RequiresApproval() bool {
	return a.GetApprovalStatus() == AnnouncementApprovalStatusPending
}

func (a *Announcement) IsApproved() bool {
	return a.GetApprovalStatus() == AnnouncementApprovalStatusApproved
}

func (a *Announcement) IsRejected() bool {
	return a.GetApprovalStatus() == AnnouncementApprovalStatusRejected
}

func (a *Announcement) ShouldSendPush() bool {
	if a.AnnouncementValues.SendPush == nil {
		return false
	}
	return *a.AnnouncementValues.SendPush && a.IsVisible()
}

func (a *Announcement) ShouldSendEmail() bool {
	if a.AnnouncementValues.SendEmail == nil {
		return false
	}
	return *a.AnnouncementValues.SendEmail && a.IsVisible()
}

func (a *Announcement) ShouldSendSMS() bool {
	if a.AnnouncementValues.SendSMS == nil {
		return false
	}
	return *a.AnnouncementValues.SendSMS && a.IsVisible()
}

// 统计更新方法
func (a *AnnouncementValues) IncrementView() *AnnouncementValues {
	count := a.GetViewCount()
	count++
	a.ViewCount = &count

	now := utils.TimeNowMilli()
	a.LastViewedAt = &now

	return a
}

func (a *AnnouncementValues) IncrementClick() *AnnouncementValues {
	count := a.GetClickCount()
	count++
	a.ClickCount = &count

	now := utils.TimeNowMilli()
	a.LastClickedAt = &now

	return a
}

func (a *AnnouncementValues) IncrementShare() *AnnouncementValues {
	count := 0
	if a.ShareCount != nil {
		count = *a.ShareCount
	}
	count++
	a.ShareCount = &count

	now := utils.TimeNowMilli()
	a.LastSharedAt = &now

	return a
}

func (a *AnnouncementValues) IncrementFeedback() *AnnouncementValues {
	count := 0
	if a.FeedbackCount != nil {
		count = *a.FeedbackCount
	}
	count++
	a.FeedbackCount = &count

	return a
}

func (a *AnnouncementValues) UpdatePushStats(sent, opened int) *AnnouncementValues {
	if a.PushSentCount != nil {
		sentCount := *a.PushSentCount + sent
		a.PushSentCount = &sentCount
	} else {
		a.PushSentCount = &sent
	}

	if a.PushOpenCount != nil {
		openCount := *a.PushOpenCount + opened
		a.PushOpenCount = &openCount
	} else {
		a.PushOpenCount = &opened
	}

	// 计算打开率
	a.CalculatePushOpenRate()

	return a
}

func (a *AnnouncementValues) UpdateEmailStats(sent, opened int) *AnnouncementValues {
	if a.EmailSentCount != nil {
		sentCount := *a.EmailSentCount + sent
		a.EmailSentCount = &sentCount
	} else {
		a.EmailSentCount = &sent
	}

	if a.EmailOpenCount != nil {
		openCount := *a.EmailOpenCount + opened
		a.EmailOpenCount = &openCount
	} else {
		a.EmailOpenCount = &opened
	}

	// 计算打开率
	a.CalculateEmailOpenRate()

	return a
}

func (a *AnnouncementValues) CalculateClickRate() *AnnouncementValues {
	viewCount := a.GetViewCount()
	if viewCount == 0 {
		return a
	}

	clickCount := a.GetClickCount()
	rate := float64(clickCount) / float64(viewCount) * 100.0
	a.ClickRate = &rate

	return a
}

func (a *AnnouncementValues) CalculateShareRate() *AnnouncementValues {
	viewCount := a.GetViewCount()
	if viewCount == 0 {
		return a
	}

	shareCount := 0
	if a.ShareCount != nil {
		shareCount = *a.ShareCount
	}

	rate := float64(shareCount) / float64(viewCount) * 100.0
	a.ShareRate = &rate

	return a
}

func (a *AnnouncementValues) CalculatePushOpenRate() *AnnouncementValues {
	sentCount := 0
	if a.PushSentCount != nil {
		sentCount = *a.PushSentCount
	}

	if sentCount == 0 {
		return a
	}

	openCount := 0
	if a.PushOpenCount != nil {
		openCount = *a.PushOpenCount
	}

	rate := float64(openCount) / float64(sentCount) * 100.0
	a.PushOpenRate = &rate

	return a
}

func (a *AnnouncementValues) CalculateEmailOpenRate() *AnnouncementValues {
	sentCount := 0
	if a.EmailSentCount != nil {
		sentCount = *a.EmailSentCount
	}

	if sentCount == 0 {
		return a
	}

	openCount := 0
	if a.EmailOpenCount != nil {
		openCount = *a.EmailOpenCount
	}

	rate := float64(openCount) / float64(sentCount) * 100.0
	a.EmailOpenRate = &rate

	return a
}

func (a *AnnouncementValues) CalculateEngagementRate() *AnnouncementValues {
	viewCount := a.GetViewCount()
	if viewCount == 0 {
		return a
	}

	clickCount := a.GetClickCount()
	shareCount := 0
	if a.ShareCount != nil {
		shareCount = *a.ShareCount
	}
	feedbackCount := 0
	if a.FeedbackCount != nil {
		feedbackCount = *a.FeedbackCount
	}

	engagements := clickCount + shareCount + feedbackCount
	rate := float64(engagements) / float64(viewCount) * 100.0
	a.EngagementRate = &rate

	return a
}

// 状态管理
func (a *AnnouncementValues) Publish() *AnnouncementValues {
	a.SetStatus(AnnouncementStatusPublished)
	a.IsActive = utils.BoolPtr(true)
	a.IsPublished = utils.BoolPtr(true)
	now := utils.TimeNowMilli()
	a.PublishedAt = &now
	a.ActivatedAt = &now
	return a
}

func (a *AnnouncementValues) Schedule(publishTime int64) *AnnouncementValues {
	a.SetStatus(AnnouncementStatusScheduled)
	a.IsScheduled = utils.BoolPtr(true)
	a.PublishTime = &publishTime
	return a
}

func (a *AnnouncementValues) Pause() *AnnouncementValues {
	a.SetStatus(AnnouncementStatusPaused)
	a.IsActive = utils.BoolPtr(false)
	now := utils.TimeNowMilli()
	a.PausedAt = &now
	return a
}

func (a *AnnouncementValues) MarkExpired() *AnnouncementValues {
	a.SetStatus(AnnouncementStatusExpired)
	a.IsActive = utils.BoolPtr(false)
	return a
}

func (a *AnnouncementValues) Delete() *AnnouncementValues {
	a.SetStatus(AnnouncementStatusDeleted)
	a.IsActive = utils.BoolPtr(false)
	return a
}

// 审批管理
func (a *AnnouncementValues) Approve(approverID, notes string) *AnnouncementValues {
	a.ApprovalStatus = utils.StringPtr(AnnouncementApprovalStatusApproved)
	a.ApprovedBy = &approverID
	a.ApprovalNotes = &notes
	now := utils.TimeNowMilli()
	a.ApprovedAt = &now

	// 自动发布已审批的公告
	a.Publish()

	return a
}

func (a *AnnouncementValues) Reject(approverID, reason string) *AnnouncementValues {
	a.ApprovalStatus = utils.StringPtr(AnnouncementApprovalStatusRejected)
	a.ApprovedBy = &approverID
	a.ApprovalNotes = &reason
	now := utils.TimeNowMilli()
	a.ApprovedAt = &now

	a.SetStatus(AnnouncementStatusDraft)

	return a
}

// 标签管理
func (a *AnnouncementValues) AddTag(tag string) *AnnouncementValues {
	var tags []string
	if a.Tags != nil && *a.Tags != "" {
		tags = strings.Split(*a.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return a
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	a.Tags = &tagsStr
	return a
}

func (a *AnnouncementValues) HasTag(tag string) bool {
	if a.Tags == nil || *a.Tags == "" {
		return false
	}

	tags := strings.Split(*a.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 便捷创建方法
func NewSystemAnnouncement(title, content string) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypeSystem).
		SetLevel(AnnouncementLevelInfo).
		SetTargetAudience(TargetAudienceAll).
		SetDisplaySettings(false, false, false).
		SetCreator("system", "system", "system")

	announcement.AddTag("system")
	announcement.AddTag("auto-generated")

	return announcement
}

func NewPromotionAnnouncement(title, content string, promoID string) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypePromotion).
		SetLevel(AnnouncementLevelInfo).
		SetTargetAudience(TargetAudienceUsers).
		SetDisplaySettings(true, true, false).
		SetAction(ActionTypeURL, "", "查看详情")

	announcement.AnnouncementValues.RelatedPromoID = &promoID
	announcement.AddTag("promotion")
	announcement.AddTag("marketing")

	return announcement
}

func NewMaintenanceAnnouncement(title, content string, startTime, endTime int64) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypeMaintenance).
		SetLevel(AnnouncementLevelWarning).
		SetTargetAudience(TargetAudienceAll).
		SetTimeRange(utils.TimeNowMilli(), startTime, endTime).
		SetDisplaySettings(true, true, true).
		SetInteractionSettings(true, true).
		SetUrgent(true)

	announcement.AddTag("maintenance")
	announcement.AddTag("urgent")
	announcement.AddTag("system")

	return announcement
}

func NewEmergencyAnnouncement(title, content string) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypeEmergency).
		SetLevel(AnnouncementLevelCritical).
		SetTargetAudience(TargetAudienceAll).
		SetDisplaySettings(true, true, true).
		SetInteractionSettings(true, true).
		SetUrgent(true).
		SetPriority(100).
		SetPushSettings(true, title, content, utils.TimeNowMilli())

	announcement.AddTag("emergency")
	announcement.AddTag("critical")
	announcement.AddTag("urgent")

	// 紧急公告自动审批
	announcement.Approve("system", "Emergency announcement auto-approved")

	return announcement
}

func NewFeatureAnnouncement(title, content, featureID string) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypeFeature).
		SetLevel(AnnouncementLevelInfo).
		SetTargetAudience(TargetAudienceAll).
		SetDisplaySettings(false, true, false).
		SetAction(ActionTypeDeeplink, "", "体验新功能")

	announcement.AnnouncementValues.RelatedFeatureID = &featureID
	announcement.AddTag("feature")
	announcement.AddTag("update")

	return announcement
}

func NewSecurityAnnouncement(title, content string) *Announcement {
	announcement := NewAnnouncementV2()
	announcement.SetTitle(title).
		SetContent(content).
		SetType(AnnouncementTypeSecurity).
		SetLevel(AnnouncementLevelWarning).
		SetTargetAudience(TargetAudienceAll).
		SetDisplaySettings(true, true, true).
		SetInteractionSettings(true, true).
		SetPushSettings(true, title, content, utils.TimeNowMilli()).
		SetEmailSettings(true, title, "security_alert")

	announcement.AddTag("security")
	announcement.AddTag("important")
	announcement.AddTag("notification")

	return announcement
}
