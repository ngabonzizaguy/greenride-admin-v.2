package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// Notification 通知表 - 基于最新设计文档
type Notification struct {
	ID             int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	NotificationID string `json:"notification_id" gorm:"column:notification_id;type:varchar(64);uniqueIndex"`
	Salt           string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*NotificationValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type NotificationValues struct {
	UserID   *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver, admin

	// 通知基本信息
	Type     *string `json:"type" gorm:"column:type;type:varchar(64);index"`         // 通知类型
	Category *string `json:"category" gorm:"column:category;type:varchar(64);index"` // 通知分类
	Title    *string `json:"title" gorm:"column:title;type:varchar(255)"`            // 通知标题
	Content  *string `json:"content" gorm:"column:content;type:text"`                // 通知内容
	Summary  *string `json:"summary" gorm:"column:summary;type:varchar(500)"`        // 通知摘要

	// 多语言支持
	TitleEn   *string `json:"title_en" gorm:"column:title_en;type:varchar(255)"`
	ContentEn *string `json:"content_en" gorm:"column:content_en;type:text"`

	// 状态信息
	Status     *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'pending'"` // pending, sent, delivered, read, failed
	Priority   *string `json:"priority" gorm:"column:priority;type:varchar(32);default:'normal'"`    // low, normal, high, urgent
	IsRead     *bool   `json:"is_read" gorm:"column:is_read;default:false"`
	IsArchived *bool   `json:"is_archived" gorm:"column:is_archived;default:false"`

	// 发送渠道
	Channels *string `json:"channels" gorm:"column:channels;type:json"` // ["push", "email", "sms"] JSON数组

	// 关联信息
	RelatedType *string `json:"related_type" gorm:"column:related_type;type:varchar(64)"`   // order, payment, user, vehicle等
	RelatedID   *string `json:"related_id" gorm:"column:related_id;type:varchar(64);index"` // 关联对象ID

	// 推送信息
	PushToken *string `json:"push_token" gorm:"column:push_token;type:varchar(512)"`
	PushTitle *string `json:"push_title" gorm:"column:push_title;type:varchar(255)"`
	PushBody  *string `json:"push_body" gorm:"column:push_body;type:text"`
	PushData  *string `json:"push_data" gorm:"column:push_data;type:json"` // 推送附加数据
	PushBadge *int    `json:"push_badge" gorm:"column:push_badge"`         // iOS badge数量
	PushSound *string `json:"push_sound" gorm:"column:push_sound;type:varchar(100)"`
	PushIcon  *string `json:"push_icon" gorm:"column:push_icon;type:varchar(255)"`

	// 邮件信息
	EmailTo      *string `json:"email_to" gorm:"column:email_to;type:varchar(255)"`
	EmailSubject *string `json:"email_subject" gorm:"column:email_subject;type:varchar(255)"`
	EmailBody    *string `json:"email_body" gorm:"column:email_body;type:text"`
	EmailHTML    *string `json:"email_html" gorm:"column:email_html;type:text"`

	// 短信信息
	SMSTo      *string `json:"sms_to" gorm:"column:sms_to;type:varchar(32)"`
	SMSContent *string `json:"sms_content" gorm:"column:sms_content;type:varchar(500)"`

	// 时间信息
	ScheduledAt *int64 `json:"scheduled_at" gorm:"column:scheduled_at"` // 计划发送时间
	SentAt      *int64 `json:"sent_at" gorm:"column:sent_at"`           // 实际发送时间
	DeliveredAt *int64 `json:"delivered_at" gorm:"column:delivered_at"` // 送达时间
	ReadAt      *int64 `json:"read_at" gorm:"column:read_at"`           // 读取时间
	ExpiresAt   *int64 `json:"expires_at" gorm:"column:expires_at"`     // 过期时间

	// 重试信息
	RetryCount   *int    `json:"retry_count" gorm:"column:retry_count;default:0"`
	MaxRetries   *int    `json:"max_retries" gorm:"column:max_retries;default:3"`
	LastRetryAt  *int64  `json:"last_retry_at" gorm:"column:last_retry_at"`
	ErrorMessage *string `json:"error_message" gorm:"column:error_message;type:text"`

	// 扩展信息
	Tags     *string `json:"tags" gorm:"column:tags;type:json"`         // 标签数组
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 额外元数据

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Notification) TableName() string {
	return "t_notifications"
}

// 通知状态常量
const (
	NotificationStatusPending   = "pending"
	NotificationStatusSent      = "sent"
	NotificationStatusDelivered = "delivered"
	NotificationStatusRead      = "read"
	NotificationStatusFailed    = "failed"
)

// 通知优先级常量
const (
	NotificationPriorityLow    = "low"
	NotificationPriorityNormal = "normal"
	NotificationPriorityHigh   = "high"
	NotificationPriorityUrgent = "urgent"
)

// 通知类型常量
const (
	NotificationTypeOrderUpdate    = "order_update"
	NotificationTypePaymentSuccess = "payment_success"
	NotificationTypePaymentFailed  = "payment_failed"
	NotificationTypeDriverArrived  = "driver_arrived"
	NotificationTypeRideStarted    = "ride_started"
	NotificationTypeRideCompleted  = "ride_completed"
	NotificationTypePromotion      = "promotion"
	NotificationTypeSystem         = "system"
)

// 通知分类常量
const (
	NotificationCategoryTransaction = "transaction"
	NotificationCategoryOrder       = "order"
	NotificationCategoryMarketing   = "marketing"
	NotificationCategorySystem      = "system"
	NotificationCategoryAlert       = "alert"
)

// 创建新的通知对象
func NewNotificationV2() *Notification {
	return &Notification{
		NotificationID: utils.GenerateNotificationID(),
		Salt:           utils.GenerateSalt(),
		NotificationValues: &NotificationValues{
			UserType:   utils.StringPtr(protocol.UserTypePassenger),
			Status:     utils.StringPtr(NotificationStatusPending),
			Priority:   utils.StringPtr(NotificationPriorityNormal),
			IsRead:     utils.BoolPtr(false),
			IsArchived: utils.BoolPtr(false),
			RetryCount: utils.IntPtr(0),
			MaxRetries: utils.IntPtr(3),
		},
	}
}

// SetValues 更新NotificationV2Values中的非nil值
func (n *NotificationValues) SetValues(values *NotificationValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		n.UserID = values.UserID
	}
	if values.UserType != nil {
		n.UserType = values.UserType
	}
	if values.Type != nil {
		n.Type = values.Type
	}
	if values.Category != nil {
		n.Category = values.Category
	}
	if values.Title != nil {
		n.Title = values.Title
	}
	if values.Content != nil {
		n.Content = values.Content
	}
	if values.Status != nil {
		n.Status = values.Status
	}
	if values.Priority != nil {
		n.Priority = values.Priority
	}
	if values.IsRead != nil {
		n.IsRead = values.IsRead
	}
	if values.RelatedType != nil {
		n.RelatedType = values.RelatedType
	}
	if values.RelatedID != nil {
		n.RelatedID = values.RelatedID
	}
	if values.Channels != nil {
		n.Channels = values.Channels
	}
	if values.Metadata != nil {
		n.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		n.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (n *NotificationValues) GetUserID() string {
	if n.UserID == nil {
		return ""
	}
	return *n.UserID
}

func (n *NotificationValues) GetType() string {
	if n.Type == nil {
		return ""
	}
	return *n.Type
}

func (n *NotificationValues) GetCategory() string {
	if n.Category == nil {
		return ""
	}
	return *n.Category
}

func (n *NotificationValues) GetTitle() string {
	if n.Title == nil {
		return ""
	}
	return *n.Title
}

func (n *NotificationValues) GetContent() string {
	if n.Content == nil {
		return ""
	}
	return *n.Content
}

func (n *NotificationValues) GetStatus() string {
	if n.Status == nil {
		return NotificationStatusPending
	}
	return *n.Status
}

func (n *NotificationValues) GetPriority() string {
	if n.Priority == nil {
		return NotificationPriorityNormal
	}
	return *n.Priority
}

func (n *NotificationValues) GetIsRead() bool {
	if n.IsRead == nil {
		return false
	}
	return *n.IsRead
}

func (n *NotificationValues) GetIsArchived() bool {
	if n.IsArchived == nil {
		return false
	}
	return *n.IsArchived
}

func (n *NotificationValues) GetRetryCount() int {
	if n.RetryCount == nil {
		return 0
	}
	return *n.RetryCount
}

func (n *NotificationValues) GetMaxRetries() int {
	if n.MaxRetries == nil {
		return 3
	}
	return *n.MaxRetries
}

func (n *NotificationValues) GetRelatedType() string {
	if n.RelatedType == nil {
		return ""
	}
	return *n.RelatedType
}

func (n *NotificationValues) GetRelatedID() string {
	if n.RelatedID == nil {
		return ""
	}
	return *n.RelatedID
}

// Setter 方法
func (n *NotificationValues) SetUserID(userID string) *NotificationValues {
	n.UserID = &userID
	return n
}

func (n *NotificationValues) SetType(notType string) *NotificationValues {
	n.Type = &notType
	return n
}

func (n *NotificationValues) SetTitle(title string) *NotificationValues {
	n.Title = &title
	return n
}

func (n *NotificationValues) SetContent(content string) *NotificationValues {
	n.Content = &content
	return n
}

func (n *NotificationValues) SetStatus(status string) *NotificationValues {
	n.Status = &status
	return n
}

func (n *NotificationValues) SetPriority(priority string) *NotificationValues {
	n.Priority = &priority
	return n
}

func (n *NotificationValues) SetRelated(relatedType, relatedID string) *NotificationValues {
	n.RelatedType = &relatedType
	n.RelatedID = &relatedID
	return n
}

func (n *NotificationValues) MarkAsRead() *NotificationValues {
	n.IsRead = utils.BoolPtr(true)
	now := utils.TimeNowMilli()
	n.ReadAt = &now
	return n
}

func (n *NotificationValues) MarkAsArchived() *NotificationValues {
	n.IsArchived = utils.BoolPtr(true)
	return n
}

// 业务方法
func (n *Notification) IsPending() bool {
	return n.GetStatus() == NotificationStatusPending
}

func (n *Notification) IsSent() bool {
	return n.GetStatus() == NotificationStatusSent
}

func (n *Notification) IsDelivered() bool {
	return n.GetStatus() == NotificationStatusDelivered
}

func (n *Notification) IsRead() bool {
	return n.GetIsRead()
}

func (n *Notification) IsFailed() bool {
	return n.GetStatus() == NotificationStatusFailed
}

func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return utils.TimeNowMilli() > *n.ExpiresAt
}

func (n *Notification) CanRetry() bool {
	return n.GetRetryCount() < n.GetMaxRetries() && n.IsFailed()
}

func (n *Notification) IsHighPriority() bool {
	priority := n.GetPriority()
	return priority == NotificationPriorityHigh || priority == NotificationPriorityUrgent
}

// 状态更新方法
func (n *NotificationValues) MarkAsSent() error {
	n.SetStatus(NotificationStatusSent)
	now := utils.TimeNowMilli()
	n.SentAt = &now
	return nil
}

func (n *NotificationValues) MarkAsDelivered() error {
	n.SetStatus(NotificationStatusDelivered)
	now := utils.TimeNowMilli()
	n.DeliveredAt = &now
	return nil
}

func (n *NotificationValues) MarkAsFailed(errorMsg string) error {
	n.SetStatus(NotificationStatusFailed)
	n.ErrorMessage = &errorMsg

	// 增加重试次数
	retryCount := n.GetRetryCount() + 1
	n.RetryCount = &retryCount
	now := utils.TimeNowMilli()
	n.LastRetryAt = &now

	return nil
}

func (n *NotificationValues) SetScheduledTime(scheduledAt int64) *NotificationValues {
	n.ScheduledAt = &scheduledAt
	return n
}

func (n *NotificationValues) SetExpiryTime(expiresAt int64) *NotificationValues {
	n.ExpiresAt = &expiresAt
	return n
}

// 推送相关方法
func (n *NotificationValues) SetPushInfo(token, title, body string, data map[string]interface{}) error {
	n.PushToken = &token
	n.PushTitle = &title
	n.PushBody = &body

	if data != nil {
		dataJSON, err := utils.ToJSON(data)
		if err != nil {
			return fmt.Errorf("failed to marshal push data: %v", err)
		}
		n.PushData = &dataJSON
	}

	return nil
}

// 邮件相关方法
func (n *NotificationValues) SetEmailInfo(to, subject, body, html string) *NotificationValues {
	n.EmailTo = &to
	n.EmailSubject = &subject
	n.EmailBody = &body
	if html != "" {
		n.EmailHTML = &html
	}
	return n
}

// 短信相关方法
func (n *NotificationValues) SetSMSInfo(to, content string) *NotificationValues {
	n.SMSTo = &to
	n.SMSContent = &content
	return n
}

// 设置发送渠道
func (n *NotificationValues) SetChannels(channels []string) error {
	if len(channels) == 0 {
		return fmt.Errorf("channels cannot be empty")
	}

	channelsJSON, err := utils.ToJSON(channels)
	if err != nil {
		return fmt.Errorf("failed to marshal channels: %v", err)
	}

	n.Channels = &channelsJSON
	return nil
}

// 设置标签
func (n *NotificationValues) SetTags(tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	tagsJSON, err := utils.ToJSON(tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v", err)
	}

	n.Tags = &tagsJSON
	return nil
}

// 检查是否应该发送
func (n *Notification) ShouldSend() bool {
	// 已过期不发送
	if n.IsExpired() {
		return false
	}

	// 非待发送状态不发送
	if !n.IsPending() {
		return false
	}

	// 检查计划发送时间
	if n.ScheduledAt != nil && utils.TimeNowMilli() < *n.ScheduledAt {
		return false
	}

	return true
}

// 获取显示标题（支持多语言）
func (n *NotificationValues) GetDisplayTitle(lang string) string {
	if lang == "en" && n.TitleEn != nil && *n.TitleEn != "" {
		return *n.TitleEn
	}
	return n.GetTitle()
}

// 获取显示内容（支持多语言）
func (n *NotificationValues) GetDisplayContent(lang string) string {
	if lang == "en" && n.ContentEn != nil && *n.ContentEn != "" {
		return *n.ContentEn
	}
	return n.GetContent()
}
