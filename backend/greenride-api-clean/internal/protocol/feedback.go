package protocol

// ============================================================================
// Public API - Feedback Request/Response
// ============================================================================

// FeedbackRequest 简化的反馈请求结构
type FeedbackRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

// FeedbackResponse 反馈响应结构
type FeedbackResponse struct {
	FeedbackID string `json:"feedback_id"`
}

// ============================================================================
// Admin API - Feedback Search/Management
// ============================================================================

// FeedbackSearchRequest 反馈搜索请求
type FeedbackSearchRequest struct {
	Keyword      string `json:"keyword,omitempty"`       // 搜索关键字
	Page         int    `json:"page,omitempty"`          // 页码，默认1
	Limit        int    `json:"limit,omitempty"`         // 每页数量，默认10
	Status       string `json:"status,omitempty"`        // 状态过滤
	FeedbackType string `json:"feedback_type,omitempty"` // 反馈类型
	Category     string `json:"category,omitempty"`      // 分类
	Severity     string `json:"severity,omitempty"`      // 严重程度
	Priority     string `json:"priority,omitempty"`      // 优先级
	UserID       string `json:"user_id,omitempty"`       // 用户ID
	StartDate    *int64 `json:"start_date,omitempty"`    // 开始日期 (时间戳毫秒)
	EndDate      *int64 `json:"end_date,omitempty"`      // 结束日期 (时间戳毫秒)
}

// FeedbackIDRequest 反馈ID请求
type FeedbackIDRequest struct {
	FeedbackID string `json:"feedback_id" binding:"required"`
}

// FeedbackUpdateRequest 反馈更新请求
type FeedbackUpdateRequest struct {
	FeedbackID     string  `json:"feedback_id" binding:"required"` // 反馈ID
	Status         *string `json:"status,omitempty"`               // 状态
	Priority       *string `json:"priority,omitempty"`             // 优先级
	Severity       *string `json:"severity,omitempty"`             // 严重程度
	AssignedTo     *string `json:"assigned_to,omitempty"`          // 分配给
	Resolution     *string `json:"resolution,omitempty"`           // 解决方案
	ResolutionNote *string `json:"resolution_note,omitempty"`      // 解决说明
	InternalNotes  *string `json:"internal_notes,omitempty"`       // 内部备注
	PublicResponse *string `json:"public_response,omitempty"`      // 公开回复
}

// ============================================================================
// Admin API - Feedback Response DTOs
// ============================================================================

// FeedbackListItem 反馈列表项
type FeedbackListItem struct {
	FeedbackID   string `json:"feedback_id"`
	Title        string `json:"title"`
	FeedbackType string `json:"feedback_type"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	Severity     string `json:"severity"`
	Priority     string `json:"priority"`
	Rating       int    `json:"rating"`
	UserName     string `json:"user_name,omitempty"`
	UserEmail    string `json:"user_email,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// FeedbackUserInfo 反馈用户信息
type FeedbackUserInfo struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// FeedbackDetail 反馈详情
type FeedbackDetail struct {
	FeedbackID   string            `json:"feedback_id"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Description  string            `json:"description,omitempty"`
	FeedbackType string            `json:"feedback_type"`
	Category     string            `json:"category"`
	Subcategory  string            `json:"subcategory,omitempty"`
	Status       string            `json:"status"`
	Severity     string            `json:"severity"`
	Priority     string            `json:"priority"`
	Rating       int               `json:"rating"`
	Resolution   string            `json:"resolution,omitempty"`
	User         *FeedbackUserInfo `json:"user,omitempty"`

	// Contact info
	ContactName  string `json:"contact_name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`

	// Assignment
	AssignedTo string `json:"assigned_to,omitempty"`
	HandledBy  string `json:"handled_by,omitempty"`

	// Related order
	OrderID string `json:"order_id,omitempty"`

	// Location
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Location  string   `json:"location,omitempty"`

	// Device info
	DeviceType string `json:"device_type,omitempty"`
	AppVersion string `json:"app_version,omitempty"`

	// Attachments
	Attachments []string `json:"attachments,omitempty"`

	// Notes
	InternalNotes  string `json:"internal_notes,omitempty"`
	PublicResponse string `json:"public_response,omitempty"`

	// Timestamps
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	ResolvedAt *int64 `json:"resolved_at,omitempty"`
}

// FeedbackStats 反馈统计
type FeedbackStats struct {
	TotalFeedback     int64   `json:"total_feedback"`
	PendingCount      int64   `json:"pending_count"`
	InProgressCount   int64   `json:"in_progress_count"`
	ResolvedCount     int64   `json:"resolved_count"`
	ClosedCount       int64   `json:"closed_count"`
	ComplaintCount    int64   `json:"complaint_count"`
	SuggestionCount   int64   `json:"suggestion_count"`
	ComplimentCount   int64   `json:"compliment_count"`
	BugReportCount    int64   `json:"bug_report_count"`
	AvgResponseTime   float64 `json:"avg_response_time"`   // 平均响应时间（小时）
	AvgResolutionTime float64 `json:"avg_resolution_time"` // 平均解决时间（小时）
}

// ============================================================================
// Admin API - Support Configuration
// ============================================================================

// SupportConfigResponse 支持配置响应
type SupportConfigResponse struct {
	SupportEmail         string `json:"support_email"`
	SupportPhone         string `json:"support_phone"`
	SupportHours         string `json:"support_hours"`
	EmergencyPhone       string `json:"emergency_phone"`
	WhatsAppNumber       string `json:"whatsapp_number"`
	ResponseTimeTarget   int    `json:"response_time_target"` // 小时
	AutoReplyEnabled     bool   `json:"auto_reply_enabled"`
	AutoReplyMessage     string `json:"auto_reply_message,omitempty"`
	EscalationEnabled    bool   `json:"escalation_enabled"`
	EscalationTimeout    int    `json:"escalation_timeout"` // 小时
	WorkdayStart         string `json:"workday_start,omitempty"`
	WorkdayEnd           string `json:"workday_end,omitempty"`
	WorkDays             string `json:"work_days,omitempty"`
	NotifyOnNewFeedback  bool   `json:"notify_on_new_feedback"`
	NotifyOnHighPriority bool   `json:"notify_on_high_priority"`
	NotifyOnSafetyIssue  bool   `json:"notify_on_safety_issue"`
	NotifyOnEscalation   bool   `json:"notify_on_escalation"`
	FAQUrl               string `json:"faq_url,omitempty"`
	HelpCenterUrl        string `json:"help_center_url,omitempty"`
	UpdatedAt            int64  `json:"updated_at"`
}

// SupportConfigUpdateRequest 支持配置更新请求
type SupportConfigUpdateRequest struct {
	SupportEmail         *string `json:"support_email,omitempty"`
	SupportPhone         *string `json:"support_phone,omitempty"`
	SupportHours         *string `json:"support_hours,omitempty"`
	EmergencyPhone       *string `json:"emergency_phone,omitempty"`
	WhatsAppNumber       *string `json:"whatsapp_number,omitempty"`
	ResponseTimeTarget   *int    `json:"response_time_target,omitempty"`
	AutoReplyEnabled     *bool   `json:"auto_reply_enabled,omitempty"`
	AutoReplyMessage     *string `json:"auto_reply_message,omitempty"`
	EscalationEnabled    *bool   `json:"escalation_enabled,omitempty"`
	EscalationTimeout    *int    `json:"escalation_timeout,omitempty"`
	WorkdayStart         *string `json:"workday_start,omitempty"`
	WorkdayEnd           *string `json:"workday_end,omitempty"`
	WorkDays             *string `json:"work_days,omitempty"`
	NotifyOnNewFeedback  *bool   `json:"notify_on_new_feedback,omitempty"`
	NotifyOnHighPriority *bool   `json:"notify_on_high_priority,omitempty"`
	NotifyOnSafetyIssue  *bool   `json:"notify_on_safety_issue,omitempty"`
	NotifyOnEscalation   *bool   `json:"notify_on_escalation,omitempty"`
	FAQUrl               *string `json:"faq_url,omitempty"`
	HelpCenterUrl        *string `json:"help_center_url,omitempty"`
}
