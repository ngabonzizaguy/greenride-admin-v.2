package models

import (
	"greenride/internal/utils"
)

// FCMMessageLog FCM消息日志表
type FCMMessageLog struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	MessageID string `json:"message_id" gorm:"column:message_id;type:varchar(126);uniqueIndex"`
	Salt      string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*FCMMessageLogValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type FCMMessageLogValues struct {
	UserID       *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType     *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver
	TokenID      *string `json:"token_id" gorm:"column:token_id;type:varchar(256);index"`
	FCMToken     *string `json:"fcm_token" gorm:"column:fcm_token;type:text"`
	FCMMessageID *string `json:"fcm_message_id" gorm:"column:fcm_message_id;type:varchar(256);index"` // Firebase返回的messageID
	Title        *string `json:"title" gorm:"column:title;type:varchar(255)"`
	Body         *string `json:"body" gorm:"column:body;type:text"`
	Data         *string `json:"data" gorm:"column:data;type:json"`
	ImageURL     *string `json:"image_url" gorm:"column:image_url;type:varchar(500)"`
	ClickAction  *string `json:"click_action" gorm:"column:click_action;type:varchar(255)"`
	Priority     *string `json:"priority" gorm:"column:priority;type:varchar(20);default:'normal'"` // high, normal
	CollapseKey  *string `json:"collapse_key" gorm:"column:collapse_key;type:varchar(64)"`

	// 发送状态
	Status          *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'pending'"` // pending, sent, failed, delivered, opened
	ResponseCode    *int    `json:"response_code" gorm:"column:response_code;type:int"`
	ResponseMessage *string `json:"response_message" gorm:"column:response_message;type:text"`
	ResponseData    *string `json:"response_data" gorm:"column:response_data;type:json"`
	SentAt          *int64  `json:"sent_at" gorm:"column:sent_at"`
	DeliveredAt     *int64  `json:"delivered_at" gorm:"column:delivered_at"`
	OpenedAt        *int64  `json:"opened_at" gorm:"column:opened_at"`
	FailedAt        *int64  `json:"failed_at" gorm:"column:failed_at"`

	// 重试信息
	RetryCount  *int   `json:"retry_count" gorm:"column:retry_count;type:int;default:0"`
	MaxRetries  *int   `json:"max_retries" gorm:"column:max_retries;type:int;default:3"`
	NextRetryAt *int64 `json:"next_retry_at" gorm:"column:next_retry_at"`

	// 分类和标签
	MessageType *string `json:"message_type" gorm:"column:message_type;type:varchar(64);index"` // booking, notification, marketing, system
	Category    *string `json:"category" gorm:"column:category;type:varchar(64);index"`
	Tags        *string `json:"tags" gorm:"column:tags;type:json"`

	// 环境信息
	IsSandbox *int `json:"is_sandbox" gorm:"column:is_sandbox;type:int;default:0"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (FCMMessageLog) TableName() string {
	return "t_fcm_message_logs"
}

// SetValues 更新FCMMessageLogV2Values中的非nil值
func (f *FCMMessageLogValues) SetValues(values *FCMMessageLogValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		f.UserID = values.UserID
	}
	if values.UserType != nil {
		f.UserType = values.UserType
	}
	if values.TokenID != nil {
		f.TokenID = values.TokenID
	}
	if values.FCMToken != nil {
		f.FCMToken = values.FCMToken
	}
	if values.FCMMessageID != nil {
		f.FCMMessageID = values.FCMMessageID
	}
	if values.Title != nil {
		f.Title = values.Title
	}
	if values.Body != nil {
		f.Body = values.Body
	}
	if values.Data != nil {
		f.Data = values.Data
	}
	if values.ImageURL != nil {
		f.ImageURL = values.ImageURL
	}
	if values.ClickAction != nil {
		f.ClickAction = values.ClickAction
	}
	if values.Priority != nil {
		f.Priority = values.Priority
	}
	if values.CollapseKey != nil {
		f.CollapseKey = values.CollapseKey
	}
	if values.Status != nil {
		f.Status = values.Status
	}
	if values.ResponseCode != nil {
		f.ResponseCode = values.ResponseCode
	}
	if values.ResponseMessage != nil {
		f.ResponseMessage = values.ResponseMessage
	}
	if values.ResponseData != nil {
		f.ResponseData = values.ResponseData
	}
	if values.SentAt != nil {
		f.SentAt = values.SentAt
	}
	if values.DeliveredAt != nil {
		f.DeliveredAt = values.DeliveredAt
	}
	if values.OpenedAt != nil {
		f.OpenedAt = values.OpenedAt
	}
	if values.FailedAt != nil {
		f.FailedAt = values.FailedAt
	}
	if values.RetryCount != nil {
		f.RetryCount = values.RetryCount
	}
	if values.MaxRetries != nil {
		f.MaxRetries = values.MaxRetries
	}
	if values.NextRetryAt != nil {
		f.NextRetryAt = values.NextRetryAt
	}
	if values.MessageType != nil {
		f.MessageType = values.MessageType
	}
	if values.Category != nil {
		f.Category = values.Category
	}
	if values.Tags != nil {
		f.Tags = values.Tags
	}
	if values.IsSandbox != nil {
		f.IsSandbox = values.IsSandbox
	}

	if values.UpdatedAt > 0 {
		f.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (f *FCMMessageLogValues) GetUserID() string {
	if f.UserID == nil {
		return ""
	}
	return *f.UserID
}

func (f *FCMMessageLogValues) GetUserType() string {
	if f.UserType == nil {
		return "user"
	}
	return *f.UserType
}

func (f *FCMMessageLogValues) GetTokenID() string {
	if f.TokenID == nil {
		return ""
	}
	return *f.TokenID
}

func (f *FCMMessageLogValues) GetFCMMessageID() string {
	if f.FCMMessageID == nil {
		return ""
	}
	return *f.FCMMessageID
}

func (f *FCMMessageLogValues) GetTitle() string {
	if f.Title == nil {
		return ""
	}
	return *f.Title
}

func (f *FCMMessageLogValues) GetBody() string {
	if f.Body == nil {
		return ""
	}
	return *f.Body
}

func (f *FCMMessageLogValues) GetStatus() string {
	if f.Status == nil {
		return "pending"
	}
	return *f.Status
}

func (f *FCMMessageLogValues) GetPriority() string {
	if f.Priority == nil {
		return "normal"
	}
	return *f.Priority
}

func (f *FCMMessageLogValues) GetRetryCount() int {
	if f.RetryCount == nil {
		return 0
	}
	return *f.RetryCount
}

func (f *FCMMessageLogValues) GetMaxRetries() int {
	if f.MaxRetries == nil {
		return 3
	}
	return *f.MaxRetries
}

func (f *FCMMessageLogValues) GetMessageType() string {
	if f.MessageType == nil {
		return ""
	}
	return *f.MessageType
}

// Setter 方法
func (f *FCMMessageLogValues) SetUserID(userID string) *FCMMessageLogValues {
	f.UserID = &userID
	return f
}

func (f *FCMMessageLogValues) SetUserType(userType string) *FCMMessageLogValues {
	f.UserType = &userType
	return f
}

func (f *FCMMessageLogValues) SetTokenID(tokenID string) *FCMMessageLogValues {
	f.TokenID = &tokenID
	return f
}

func (f *FCMMessageLogValues) SetFCMMessageID(fcmMessageID string) *FCMMessageLogValues {
	f.FCMMessageID = &fcmMessageID
	return f
}

func (f *FCMMessageLogValues) SetTitle(title string) *FCMMessageLogValues {
	f.Title = &title
	return f
}

func (f *FCMMessageLogValues) SetBody(body string) *FCMMessageLogValues {
	f.Body = &body
	return f
}

func (f *FCMMessageLogValues) SetData(data string) *FCMMessageLogValues {
	f.Data = &data
	return f
}

func (f *FCMMessageLogValues) SetStatus(status string) *FCMMessageLogValues {
	f.Status = &status
	return f
}

func (f *FCMMessageLogValues) SetPriority(priority string) *FCMMessageLogValues {
	f.Priority = &priority
	return f
}

func (f *FCMMessageLogValues) SetRetryCount(retryCount int) *FCMMessageLogValues {
	f.RetryCount = &retryCount
	return f
}

func (f *FCMMessageLogValues) SetMaxRetries(maxRetries int) *FCMMessageLogValues {
	f.MaxRetries = &maxRetries
	return f
}

func (f *FCMMessageLogValues) SetMessageType(messageType string) *FCMMessageLogValues {
	f.MessageType = &messageType
	return f
}

func (f *FCMMessageLogValues) SetCategory(category string) *FCMMessageLogValues {
	f.Category = &category
	return f
}

func (f *FCMMessageLogValues) SetResponseCode(responseCode int) *FCMMessageLogValues {
	f.ResponseCode = &responseCode
	return f
}

func (f *FCMMessageLogValues) SetResponseMessage(responseMessage string) *FCMMessageLogValues {
	f.ResponseMessage = &responseMessage
	return f
}

func (f *FCMMessageLogValues) SetSentAt(sentAt int64) *FCMMessageLogValues {
	f.SentAt = &sentAt
	return f
}

func (f *FCMMessageLogValues) SetDeliveredAt(deliveredAt int64) *FCMMessageLogValues {
	f.DeliveredAt = &deliveredAt
	return f
}

func (f *FCMMessageLogValues) SetOpenedAt(openedAt int64) *FCMMessageLogValues {
	f.OpenedAt = &openedAt
	return f
}

func (f *FCMMessageLogValues) SetFailedAt(failedAt int64) *FCMMessageLogValues {
	f.FailedAt = &failedAt
	return f
}

// 业务方法
func (f *FCMMessageLog) IsSandbox() bool {
	if f.FCMMessageLogValues == nil || f.FCMMessageLogValues.IsSandbox == nil {
		return false
	}
	return *f.FCMMessageLogValues.IsSandbox == 1
}

func (f *FCMMessageLog) IsPending() bool {
	return f.GetStatus() == "pending"
}

func (f *FCMMessageLog) IsSent() bool {
	return f.GetStatus() == "sent"
}

func (f *FCMMessageLog) IsFailed() bool {
	return f.GetStatus() == "failed"
}

func (f *FCMMessageLog) IsDelivered() bool {
	return f.GetStatus() == "delivered"
}

func (f *FCMMessageLog) IsOpened() bool {
	return f.GetStatus() == "opened"
}

func (f *FCMMessageLog) CanRetry() bool {
	return f.GetRetryCount() < f.GetMaxRetries() && f.IsFailed()
}

// 标记为已发送
func (f *FCMMessageLogValues) MarkAsSent() {
	status := "sent"
	f.Status = &status
	now := utils.TimeNowMilli()
	f.SentAt = &now
}

// 标记为失败
func (f *FCMMessageLogValues) MarkAsFailed(responseCode int, responseMessage string) {
	status := "failed"
	f.Status = &status
	f.ResponseCode = &responseCode
	f.ResponseMessage = &responseMessage
	now := utils.TimeNowMilli()
	f.FailedAt = &now

	// 增加重试次数
	retryCount := f.GetRetryCount() + 1
	f.RetryCount = &retryCount
}

// 标记为已送达
func (f *FCMMessageLogValues) MarkAsDelivered() {
	status := "delivered"
	f.Status = &status
	now := utils.TimeNowMilli()
	f.DeliveredAt = &now
}

// 标记为已打开
func (f *FCMMessageLogValues) MarkAsOpened() {
	status := "opened"
	f.Status = &status
	now := utils.TimeNowMilli()
	f.OpenedAt = &now
}

// 生成MessageID
func GenerateFCMMessageID() string {
	return "msg_" + utils.GenerateShortID()
}

// 创建新的FCM消息日志对象
func NewFCMMessageLogV2() *FCMMessageLog {
	return &FCMMessageLog{
		MessageID: utils.GenerateFCMMessageID(),
		Salt:      utils.GenerateSalt(),
		FCMMessageLogValues: &FCMMessageLogValues{
			UserType:   utils.StringPtr("user"),
			Status:     utils.StringPtr("pending"),
			Priority:   utils.StringPtr("normal"),
			RetryCount: utils.IntPtr(0),
			MaxRetries: utils.IntPtr(3),
		},
	}
}
