package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"
)

// Feedback 反馈表 - 用户、司机的反馈和建议管理
type Feedback struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FeedbackID string `json:"feedback_id" gorm:"column:feedback_id;type:varchar(64);uniqueIndex"`
	Salt       string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*FeedbackValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type FeedbackValues struct {
	// 关联信息
	UserID  *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`   // 反馈用户ID
	OrderID *string `json:"order_id" gorm:"column:order_id;type:varchar(64);index"` // 相关订单ID（如果适用），关联主表order_id

	// 反馈类型和分类
	FeedbackType *string `json:"feedback_type" gorm:"column:feedback_type;type:varchar(50);index"` // complaint, suggestion, compliment, bug_report, feature_request, safety_issue
	Category     *string `json:"category" gorm:"column:category;type:varchar(50);index"`           // service, driver, vehicle, app, payment, safety, other
	Subcategory  *string `json:"subcategory" gorm:"column:subcategory;type:varchar(50)"`           // 子分类

	// 反馈内容
	Title       *string `json:"title" gorm:"column:title;type:varchar(255)"`     // 反馈标题
	Content     *string `json:"content" gorm:"column:content;type:text"`         // 反馈内容
	Description *string `json:"description" gorm:"column:description;type:text"` // 详细描述

	// 评分信息
	Rating        *int `json:"rating" gorm:"column:rating;type:int"`                 // 评分 1-5星
	ServiceRating *int `json:"service_rating" gorm:"column:service_rating;type:int"` // 服务评分
	DriverRating  *int `json:"driver_rating" gorm:"column:driver_rating;type:int"`   // 司机评分
	VehicleRating *int `json:"vehicle_rating" gorm:"column:vehicle_rating;type:int"` // 车辆评分
	AppRating     *int `json:"app_rating" gorm:"column:app_rating;type:int"`         // 应用评分

	// 严重程度和优先级
	Severity *string `json:"severity" gorm:"column:severity;type:varchar(20);index"` // low, medium, high, critical
	Priority *string `json:"priority" gorm:"column:priority;type:varchar(20);index"` // low, medium, high, urgent
	Impact   *string `json:"impact" gorm:"column:impact;type:varchar(20)"`           // minor, moderate, major, severe

	// 状态管理
	Status         *string `json:"status" gorm:"column:status;type:varchar(30);index;default:'pending'"` // pending, reviewing, in_progress, resolved, closed, rejected
	Resolution     *string `json:"resolution" gorm:"column:resolution;type:text"`                        // 解决方案
	ResolutionNote *string `json:"resolution_note" gorm:"column:resolution_note;type:text"`              // 解决说明

	// 处理信息
	AssignedTo *string `json:"assigned_to" gorm:"column:assigned_to;type:varchar(64)"` // 分配给的管理员ID
	AssignedAt *int64  `json:"assigned_at" gorm:"column:assigned_at"`                  // 分配时间
	HandledBy  *string `json:"handled_by" gorm:"column:handled_by;type:varchar(64)"`   // 处理人ID
	HandledAt  *int64  `json:"handled_at" gorm:"column:handled_at"`                    // 处理时间
	ResolvedAt *int64  `json:"resolved_at" gorm:"column:resolved_at"`                  // 解决时间
	ClosedAt   *int64  `json:"closed_at" gorm:"column:closed_at"`                      // 关闭时间

	// 联系信息
	ContactName  *string `json:"contact_name" gorm:"column:contact_name;type:varchar(100)"`   // 联系人姓名
	ContactPhone *string `json:"contact_phone" gorm:"column:contact_phone;type:varchar(30)"`  // 联系电话
	ContactEmail *string `json:"contact_email" gorm:"column:contact_email;type:varchar(255)"` // 联系邮箱

	// 位置信息
	Latitude  *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8)"`   // 事发纬度
	Longitude *float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8)"` // 事发经度
	Location  *string  `json:"location" gorm:"column:location;type:varchar(255)"`    // 位置描述

	// 时间信息
	IncidentTime *int64 `json:"incident_time" gorm:"column:incident_time"` // 事件发生时间
	ReportedAt   *int64 `json:"reported_at" gorm:"column:reported_at"`     // 报告时间

	// 设备和环境信息
	DeviceType  *string `json:"device_type" gorm:"column:device_type;type:varchar(50)"`    // iOS, Android, Web
	AppVersion  *string `json:"app_version" gorm:"column:app_version;type:varchar(50)"`    // 应用版本
	OSVersion   *string `json:"os_version" gorm:"column:os_version;type:varchar(50)"`      // 操作系统版本
	DeviceModel *string `json:"device_model" gorm:"column:device_model;type:varchar(100)"` // 设备型号

	// 附件信息
	Attachments *string `json:"attachments" gorm:"column:attachments;type:json"` // JSON数组：附件URL列表
	Screenshots *string `json:"screenshots" gorm:"column:screenshots;type:json"` // JSON数组：截图URL列表
	AudioFiles  *string `json:"audio_files" gorm:"column:audio_files;type:json"` // JSON数组：音频文件URL列表
	VideoFiles  *string `json:"video_files" gorm:"column:video_files;type:json"` // JSON数组：视频文件URL列表

	// 相关人员和证据
	Witnesses *string `json:"witnesses" gorm:"column:witnesses;type:json"` // JSON数组：证人信息
	Evidence  *string `json:"evidence" gorm:"column:evidence;type:text"`   // 证据描述

	// 满意度和后续
	SatisfactionRating *int    `json:"satisfaction_rating" gorm:"column:satisfaction_rating;type:int"`    // 解决满意度评分 1-5
	FollowUpRequired   *bool   `json:"follow_up_required" gorm:"column:follow_up_required;default:false"` // 是否需要后续跟进
	FollowUpNote       *string `json:"follow_up_note" gorm:"column:follow_up_note;type:text"`             // 后续跟进说明
	FollowUpDate       *int64  `json:"follow_up_date" gorm:"column:follow_up_date"`                       // 后续跟进日期

	// 匿名和隐私
	IsAnonymous  *bool `json:"is_anonymous" gorm:"column:is_anonymous;default:false"`  // 是否匿名反馈
	IsPublic     *bool `json:"is_public" gorm:"column:is_public;default:false"`        // 是否公开反馈
	AllowContact *bool `json:"allow_contact" gorm:"column:allow_contact;default:true"` // 是否允许联系

	// 处理时效
	ResponseTime   *int    `json:"response_time" gorm:"column:response_time;type:int"`     // 响应时间(小时)
	ResolutionTime *int    `json:"resolution_time" gorm:"column:resolution_time;type:int"` // 解决时间(小时)
	SLALevel       *string `json:"sla_level" gorm:"column:sla_level;type:varchar(20)"`     // standard, priority, emergency

	// 分析和统计
	ViewCount    *int `json:"view_count" gorm:"column:view_count;type:int;default:0"`       // 查看次数
	UpdateCount  *int `json:"update_count" gorm:"column:update_count;type:int;default:0"`   // 更新次数
	CommentCount *int `json:"comment_count" gorm:"column:comment_count;type:int;default:0"` // 评论次数

	// 情感分析
	SentimentScore *float64 `json:"sentiment_score" gorm:"column:sentiment_score;type:decimal(4,2)"` // 情感评分 -1.0 到 1.0
	EmotionType    *string  `json:"emotion_type" gorm:"column:emotion_type;type:varchar(30)"`        // angry, frustrated, satisfied, happy, neutral
	ToneAnalysis   *string  `json:"tone_analysis" gorm:"column:tone_analysis;type:varchar(50)"`      // formal, casual, aggressive, polite

	// 关键词和标签
	Keywords *string `json:"keywords" gorm:"column:keywords;type:varchar(500)"` // 关键词
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"`         // 标签
	AutoTags *string `json:"auto_tags" gorm:"column:auto_tags;type:json"`       // JSON数组：自动生成的标签

	// 相似度和关联
	SimilarFeedbacks *string `json:"similar_feedbacks" gorm:"column:similar_feedbacks;type:json"` // JSON数组：相似反馈ID列表
	RelatedTickets   *string `json:"related_tickets" gorm:"column:related_tickets;type:json"`     // JSON数组：相关工单ID列表
	DuplicateOf      *string `json:"duplicate_of" gorm:"column:duplicate_of;type:varchar(64)"`    // 重复的原始反馈ID

	// 处理成本和影响
	HandlingCost       *float64 `json:"handling_cost" gorm:"column:handling_cost;type:decimal(10,2)"`             // 处理成本
	CompensationAmount *float64 `json:"compensation_amount" gorm:"column:compensation_amount;type:decimal(10,2)"` // 补偿金额
	RefundAmount       *float64 `json:"refund_amount" gorm:"column:refund_amount;type:decimal(10,2)"`             // 退款金额

	// 系统集成
	TicketNumber    *string `json:"ticket_number" gorm:"column:ticket_number;type:varchar(50)"` // 工单号
	ExternalRef     *string `json:"external_ref" gorm:"column:external_ref;type:varchar(100)"`  // 外部系统引用
	IntegrationData *string `json:"integration_data" gorm:"column:integration_data;type:json"`  // 集成数据

	// 质量控制
	QualityScore   *float64 `json:"quality_score" gorm:"column:quality_score;type:decimal(3,2)"` // 反馈质量评分 0-5.0
	IsValidated    *bool    `json:"is_validated" gorm:"column:is_validated;default:false"`       // 是否已验证
	ValidationNote *string  `json:"validation_note" gorm:"column:validation_note;type:text"`     // 验证说明

	// 通知和提醒
	NotificationsSent *string `json:"notifications_sent" gorm:"column:notifications_sent;type:json"`  // JSON数组：已发送的通知
	RemindersSent     *int    `json:"reminders_sent" gorm:"column:reminders_sent;type:int;default:0"` // 已发送提醒次数
	LastReminderAt    *int64  `json:"last_reminder_at" gorm:"column:last_reminder_at"`                // 最后提醒时间

	// 备注和元数据
	InternalNotes  *string `json:"internal_notes" gorm:"column:internal_notes;type:text"`   // 内部备注
	PublicResponse *string `json:"public_response" gorm:"column:public_response;type:text"` // 公开回复
	Metadata       *string `json:"metadata" gorm:"column:metadata;type:json"`               // 附加元数据

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Feedback) TableName() string {
	return "t_feedbacks"
}

// 创建新的反馈对象
func NewFeedback() *Feedback {
	return &Feedback{
		FeedbackID: utils.GenerateFeedbackID(),
		Salt:       utils.GenerateSalt(),
		FeedbackValues: &FeedbackValues{
			FeedbackType:     utils.StringPtr(protocol.FeedbackTypeSuggestion),
			Category:         utils.StringPtr(protocol.FeedbackCategoryService),
			Severity:         utils.StringPtr(protocol.FeedbackSeverityMedium),
			Priority:         utils.StringPtr(protocol.FeedbackPriorityMedium),
			Status:           utils.StringPtr(protocol.StatusPending),
			IsAnonymous:      utils.BoolPtr(false),
			IsPublic:         utils.BoolPtr(false),
			AllowContact:     utils.BoolPtr(true),
			SLALevel:         utils.StringPtr(protocol.FeedbackSLAStandard),
			ViewCount:        utils.IntPtr(0),
			UpdateCount:      utils.IntPtr(0),
			CommentCount:     utils.IntPtr(0),
			EmotionType:      utils.StringPtr(protocol.EmotionTypeNeutral),
			IsValidated:      utils.BoolPtr(false),
			RemindersSent:    utils.IntPtr(0),
			FollowUpRequired: utils.BoolPtr(false),
			ReportedAt:       utils.Int64Ptr(utils.TimeNowMilli()),
		},
	}
}

// SetValues 更新FeedbackV2Values中的非nil值
func (f *FeedbackValues) SetValues(values *FeedbackValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		f.UserID = values.UserID
	}
	if values.FeedbackType != nil {
		f.FeedbackType = values.FeedbackType
	}
	if values.Category != nil {
		f.Category = values.Category
	}
	if values.Title != nil {
		f.Title = values.Title
	}
	if values.Content != nil {
		f.Content = values.Content
	}
	if values.Rating != nil {
		f.Rating = values.Rating
	}
	if values.Severity != nil {
		f.Severity = values.Severity
	}
	if values.Status != nil {
		f.Status = values.Status
	}
	if values.Resolution != nil {
		f.Resolution = values.Resolution
	}
	if values.UpdatedAt > 0 {
		f.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (f *FeedbackValues) GetUserID() string {
	if f.UserID == nil {
		return ""
	}
	return *f.UserID
}

func (f *FeedbackValues) GetOrderID() string {
	if f.OrderID == nil {
		return ""
	}
	return *f.OrderID
}

func (f *FeedbackValues) GetFeedbackType() string {
	if f.FeedbackType == nil {
		return protocol.FeedbackTypeSuggestion
	}
	return *f.FeedbackType
}

func (f *FeedbackValues) GetCategory() string {
	if f.Category == nil {
		return protocol.FeedbackCategoryService
	}
	return *f.Category
}

func (f *FeedbackValues) GetTitle() string {
	if f.Title == nil {
		return ""
	}
	return *f.Title
}

func (f *FeedbackValues) GetContent() string {
	if f.Content == nil {
		return ""
	}
	return *f.Content
}

func (f *FeedbackValues) GetRating() int {
	if f.Rating == nil {
		return 0
	}
	return *f.Rating
}

func (f *FeedbackValues) GetSeverity() string {
	if f.Severity == nil {
		return protocol.FeedbackSeverityMedium
	}
	return *f.Severity
}

func (f *FeedbackValues) GetPriority() string {
	if f.Priority == nil {
		return protocol.FeedbackPriorityMedium
	}
	return *f.Priority
}

func (f *FeedbackValues) GetStatus() string {
	if f.Status == nil {
		return protocol.StatusPending
	}
	return *f.Status
}

func (f *FeedbackValues) GetResolution() string {
	if f.Resolution == nil {
		return ""
	}
	return *f.Resolution
}

func (f *FeedbackValues) GetAssignedTo() string {
	if f.AssignedTo == nil {
		return ""
	}
	return *f.AssignedTo
}

func (f *FeedbackValues) GetIsAnonymous() bool {
	if f.IsAnonymous == nil {
		return false
	}
	return *f.IsAnonymous
}

func (f *FeedbackValues) GetIsPublic() bool {
	if f.IsPublic == nil {
		return false
	}
	return *f.IsPublic
}

func (f *FeedbackValues) GetAllowContact() bool {
	if f.AllowContact == nil {
		return true
	}
	return *f.AllowContact
}

func (f *FeedbackValues) GetSatisfactionRating() int {
	if f.SatisfactionRating == nil {
		return 0
	}
	return *f.SatisfactionRating
}

func (f *FeedbackValues) GetFollowUpRequired() bool {
	if f.FollowUpRequired == nil {
		return false
	}
	return *f.FollowUpRequired
}

func (f *FeedbackValues) GetQualityScore() float64 {
	if f.QualityScore == nil {
		return 0.0
	}
	return *f.QualityScore
}

// Setter 方法
func (f *FeedbackValues) SetUserID(userID string) *FeedbackValues {
	f.UserID = &userID
	return f
}

func (f *FeedbackValues) SetOrderID(orderID string) *FeedbackValues {
	f.OrderID = &orderID
	return f
}

func (f *FeedbackValues) SetFeedbackType(feedbackType string) *FeedbackValues {
	f.FeedbackType = &feedbackType
	return f
}

func (f *FeedbackValues) SetCategory(category string) *FeedbackValues {
	f.Category = &category
	return f
}

func (f *FeedbackValues) SetSubcategory(subcategory string) *FeedbackValues {
	f.Subcategory = &subcategory
	return f
}

func (f *FeedbackValues) SetContent(title, content, description string) *FeedbackValues {
	f.Title = &title
	f.Content = &content
	f.Description = &description
	return f
}

func (f *FeedbackValues) SetRating(rating int) *FeedbackValues {
	f.Rating = &rating
	return f
}

func (f *FeedbackValues) SetRatings(service, driver, vehicle, app int) *FeedbackValues {
	f.ServiceRating = &service
	f.DriverRating = &driver
	f.VehicleRating = &vehicle
	f.AppRating = &app
	return f
}

func (f *FeedbackValues) SetSeverity(severity string) *FeedbackValues {
	f.Severity = &severity
	return f
}

func (f *FeedbackValues) SetPriority(priority string) *FeedbackValues {
	f.Priority = &priority
	return f
}

func (f *FeedbackValues) SetStatus(status string) *FeedbackValues {
	f.Status = &status
	return f
}

func (f *FeedbackValues) SetResolution(resolution, note string) *FeedbackValues {
	f.Resolution = &resolution
	f.ResolutionNote = &note
	return f
}

func (f *FeedbackValues) SetAssignment(assignedTo string) *FeedbackValues {
	f.AssignedTo = &assignedTo
	now := utils.TimeNowMilli()
	f.AssignedAt = &now
	return f
}

func (f *FeedbackValues) SetHandler(handledBy string) *FeedbackValues {
	f.HandledBy = &handledBy
	now := utils.TimeNowMilli()
	f.HandledAt = &now
	return f
}

func (f *FeedbackValues) SetContact(name, phone, email string) *FeedbackValues {
	f.ContactName = &name
	f.ContactPhone = &phone
	f.ContactEmail = &email
	return f
}

func (f *FeedbackValues) SetLocation(lat, lng float64, location string) *FeedbackValues {
	f.Latitude = &lat
	f.Longitude = &lng
	f.Location = &location
	return f
}

func (f *FeedbackValues) SetDeviceInfo(deviceType, appVersion, osVersion, deviceModel string) *FeedbackValues {
	f.DeviceType = &deviceType
	f.AppVersion = &appVersion
	f.OSVersion = &osVersion
	f.DeviceModel = &deviceModel
	return f
}

func (f *FeedbackValues) SetPrivacy(isAnonymous, isPublic, allowContact bool) *FeedbackValues {
	f.IsAnonymous = &isAnonymous
	f.IsPublic = &isPublic
	f.AllowContact = &allowContact
	return f
}

func (f *FeedbackValues) SetSentiment(score float64, emotionType, toneAnalysis string) *FeedbackValues {
	f.SentimentScore = &score
	f.EmotionType = &emotionType
	f.ToneAnalysis = &toneAnalysis
	return f
}

func (f *FeedbackValues) SetCompensation(handlingCost, compensationAmount, refundAmount float64) *FeedbackValues {
	f.HandlingCost = &handlingCost
	f.CompensationAmount = &compensationAmount
	f.RefundAmount = &refundAmount
	return f
}

func (f *FeedbackValues) SetFollowUp(required bool, note string, date int64) *FeedbackValues {
	f.FollowUpRequired = &required
	f.FollowUpNote = &note
	f.FollowUpDate = &date
	return f
}

// 业务方法
func (f *Feedback) IsPending() bool {
	return f.GetStatus() == protocol.StatusPending
}

func (f *Feedback) IsReviewing() bool {
	return f.GetStatus() == protocol.StatusReviewing
}

func (f *Feedback) IsInProgress() bool {
	return f.GetStatus() == protocol.StatusInProgress
}

func (f *Feedback) IsResolved() bool {
	return f.GetStatus() == protocol.StatusResolved
}

func (f *Feedback) IsClosed() bool {
	return f.GetStatus() == protocol.StatusClosed
}

func (f *Feedback) IsRejected() bool {
	return f.GetStatus() == protocol.StatusRejected
}

func (f *Feedback) IsComplaint() bool {
	return f.GetFeedbackType() == protocol.FeedbackTypeComplaint
}

func (f *Feedback) IsSuggestion() bool {
	return f.GetFeedbackType() == protocol.FeedbackTypeSuggestion
}

func (f *Feedback) IsCompliment() bool {
	return f.GetFeedbackType() == protocol.FeedbackTypeCompliment
}

func (f *Feedback) IsBugReport() bool {
	return f.GetFeedbackType() == protocol.FeedbackTypeBugReport
}

func (f *Feedback) IsSafetyIssue() bool {
	return f.GetFeedbackType() == protocol.FeedbackTypeSafetyIssue
}

func (f *Feedback) IsHighPriority() bool {
	priority := f.GetPriority()
	return priority == protocol.FeedbackPriorityHigh || priority == protocol.FeedbackPriorityUrgent
}

func (f *Feedback) IsHighSeverity() bool {
	severity := f.GetSeverity()
	return severity == protocol.FeedbackSeverityHigh || severity == protocol.FeedbackSeverityCritical
}

func (f *Feedback) IsAnonymous() bool {
	return f.GetIsAnonymous()
}

func (f *Feedback) IsPublic() bool {
	return f.GetIsPublic()
}

func (f *Feedback) AllowsContact() bool {
	return f.GetAllowContact()
}

func (f *Feedback) HasLocation() bool {
	return f.FeedbackValues.Latitude != nil && f.FeedbackValues.Longitude != nil
}

func (f *Feedback) HasRating() bool {
	return f.GetRating() > 0
}

func (f *Feedback) IsPositiveRating() bool {
	return f.GetRating() >= 4
}

func (f *Feedback) IsNegativeRating() bool {
	return f.GetRating() <= 2
}

func (f *Feedback) RequiresFollowUp() bool {
	return f.GetFollowUpRequired()
}

func (f *Feedback) IsHighQuality() bool {
	return f.GetQualityScore() >= 4.0
}

func (f *Feedback) IsValidated() bool {
	if f.FeedbackValues.IsValidated == nil {
		return false
	}
	return *f.FeedbackValues.IsValidated
}

// 状态流转方法
func (f *FeedbackValues) StartReview() *FeedbackValues {
	f.SetStatus(protocol.StatusReviewing)
	return f
}

func (f *FeedbackValues) StartProgress() *FeedbackValues {
	f.SetStatus(protocol.StatusInProgress)
	return f
}

func (f *FeedbackValues) Resolve(resolution, note string) *FeedbackValues {
	f.SetStatus(protocol.StatusResolved)
	f.SetResolution(resolution, note)
	now := utils.TimeNowMilli()
	f.ResolvedAt = &now
	return f
}

func (f *FeedbackValues) Close() *FeedbackValues {
	f.SetStatus(protocol.StatusClosed)
	now := utils.TimeNowMilli()
	f.ClosedAt = &now
	return f
}

func (f *FeedbackValues) Reject(reason string) *FeedbackValues {
	f.SetStatus(protocol.StatusRejected)
	f.SetResolution("", reason)
	return f
}

// 时效计算
func (f *FeedbackValues) CalculateResponseTime() int {
	if f.ReportedAt == nil || f.HandledAt == nil {
		return 0
	}

	diff := *f.HandledAt - *f.ReportedAt
	return int(diff / 1000 / 3600) // 转换为小时
}

func (f *FeedbackValues) CalculateResolutionTime() int {
	if f.ReportedAt == nil || f.ResolvedAt == nil {
		return 0
	}

	diff := *f.ResolvedAt - *f.ReportedAt
	return int(diff / 1000 / 3600) // 转换为小时
}

func (f *FeedbackValues) UpdateTimes() *FeedbackValues {
	responseTime := f.CalculateResponseTime()
	resolutionTime := f.CalculateResolutionTime()

	if responseTime > 0 {
		f.ResponseTime = &responseTime
	}
	if resolutionTime > 0 {
		f.ResolutionTime = &resolutionTime
	}

	return f
}

// 附件管理
func (f *FeedbackValues) SetAttachments(attachments []string) error {
	attachmentsJSON, err := utils.ToJSON(attachments)
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %v", err)
	}

	f.Attachments = &attachmentsJSON
	return nil
}

func (f *FeedbackValues) GetAttachments() []string {
	if f.Attachments == nil {
		return []string{}
	}

	var attachments []string
	if err := utils.FromJSON(*f.Attachments, &attachments); err != nil {
		return []string{}
	}

	return attachments
}

func (f *FeedbackValues) AddAttachment(url string) error {
	attachments := f.GetAttachments()
	attachments = append(attachments, url)
	return f.SetAttachments(attachments)
}

// 标签管理
func (f *FeedbackValues) AddTag(tag string) *FeedbackValues {
	var tags []string
	if f.Tags != nil && *f.Tags != "" {
		tags = strings.Split(*f.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return f
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	f.Tags = &tagsStr
	return f
}

func (f *FeedbackValues) HasTag(tag string) bool {
	if f.Tags == nil || *f.Tags == "" {
		return false
	}

	tags := strings.Split(*f.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 质量评分计算
func (f *FeedbackValues) CalculateQualityScore() float64 {
	score := 0.0

	// 内容完整性 (40%)
	if f.Title != nil && len(*f.Title) > 0 {
		score += 1.0
	}
	if f.Content != nil && len(*f.Content) > 20 {
		score += 1.0
	}

	// 详细信息 (30%)
	if f.Description != nil && len(*f.Description) > 0 {
		score += 0.5
	}
	if f.Category != nil {
		score += 0.5
	}
	if f.Latitude != nil && f.Longitude != nil {
		score += 0.5
	}

	// 联系信息 (20%)
	if f.ContactName != nil && len(*f.ContactName) > 0 {
		score += 0.5
	}
	if f.ContactPhone != nil || f.ContactEmail != nil {
		score += 0.5
	}

	// 附件证据 (10%)
	attachments := f.GetAttachments()
	if len(attachments) > 0 {
		score += 0.5
	}

	// 确保评分在0-5范围内
	if score > 5.0 {
		score = 5.0
	}

	return score
}

func (f *FeedbackValues) UpdateQualityScore() *FeedbackValues {
	score := f.CalculateQualityScore()
	f.QualityScore = &score
	return f
}

// 统计更新
func (f *FeedbackValues) IncrementViewCount() *FeedbackValues {
	count := 0
	if f.ViewCount != nil {
		count = *f.ViewCount
	}
	count++
	f.ViewCount = &count
	return f
}

func (f *FeedbackValues) IncrementUpdateCount() *FeedbackValues {
	count := 0
	if f.UpdateCount != nil {
		count = *f.UpdateCount
	}
	count++
	f.UpdateCount = &count
	return f
}

// 便捷创建方法
func NewComplaintFeedback(userID, orderID string, title, content string, rating int) *Feedback {
	feedback := NewFeedback()
	feedback.SetUserID(userID).
		SetOrderID(orderID).
		SetFeedbackType(protocol.FeedbackTypeComplaint).
		SetCategory(protocol.FeedbackCategoryService).
		SetContent(title, content, "").
		SetRating(rating).
		SetSeverity(protocol.FeedbackSeverityHigh).
		SetPriority(protocol.FeedbackPriorityHigh)

	feedback.UpdateQualityScore()

	return feedback
}

func NewSuggestionFeedback(userID, title, content string) *Feedback {
	feedback := NewFeedback()
	feedback.SetUserID(userID).
		SetFeedbackType(protocol.FeedbackTypeSuggestion).
		SetCategory(protocol.FeedbackCategoryApp).
		SetContent(title, content, "").
		SetSeverity(protocol.FeedbackSeverityLow).
		SetPriority(protocol.FeedbackPriorityMedium)

	feedback.UpdateQualityScore()

	return feedback
}

func NewComplimentFeedback(userID, orderID string, title, content string, rating int) *Feedback {
	feedback := NewFeedback()
	feedback.SetUserID(userID).
		SetOrderID(orderID).
		SetFeedbackType(protocol.FeedbackTypeCompliment).
		SetCategory(protocol.FeedbackCategoryDriver).
		SetContent(title, content, "").
		SetRating(rating).
		SetSeverity(protocol.FeedbackSeverityLow).
		SetPriority(protocol.FeedbackPriorityLow).
		SetSentiment(0.8, protocol.EmotionTypeHappy, protocol.ToneAnalysisPolite)

	feedback.UpdateQualityScore()

	return feedback
}

func NewSafetyIssueFeedback(userID string, title, content string, lat, lng float64) *Feedback {
	feedback := NewFeedback()
	feedback.SetUserID(userID).
		SetFeedbackType(protocol.FeedbackTypeSafetyIssue).
		SetCategory(protocol.FeedbackCategorySafety).
		SetContent(title, content, "").
		SetLocation(lat, lng, "").
		SetSeverity(protocol.FeedbackSeverityCritical).
		SetPriority(protocol.FeedbackPriorityUrgent).
		SetFollowUp(true, "Safety issue requires immediate follow-up", utils.TimeNowMilli()+24*3600*1000)

	feedback.UpdateQualityScore()

	return feedback
}
