package protocol

// Promotion 统一优惠码实体（对应前端展示）
type Promotion struct {
	// 基础信息
	PromotionID string `json:"promotion_id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PromoType   string `json:"promo_type"` // discount, cashback, free_ride, upgrade, points

	// 优惠信息
	DiscountType      string  `json:"discount_type"` // percentage, fixed_amount, free
	DiscountValue     float64 `json:"discount_value"`
	MaxDiscountAmount float64 `json:"max_discount_amount"`
	MinOrderAmount    float64 `json:"min_order_amount"`

	// 使用限制
	UsageLimit     int `json:"usage_limit"`
	UsageCount     int `json:"usage_count"`
	UserUsageLimit int `json:"user_usage_limit"`

	// 时间信息
	StartDate int64 `json:"start_date"`
	EndDate   int64 `json:"end_date"`
	ValidDays int   `json:"valid_days"`

	// 状态信息
	Status         string `json:"status"` // active, inactive, expired, suspended, deleted
	IsActive       bool   `json:"is_active"`
	IsPublic       bool   `json:"is_public"`
	ApprovalStatus string `json:"approval_status"` // pending, approved, rejected

	// 限制条件
	TargetUserType    string `json:"target_user_type"`
	ValidCities       string `json:"valid_cities"`
	ValidVehicleTypes string `json:"valid_vehicle_types"`

	// 统计信息
	ViewCount      int     `json:"view_count"`
	ClaimCount     int     `json:"claim_count"`
	ConversionRate float64 `json:"conversion_rate"`

	// 管理信息
	CreatedBy  string `json:"created_by"`
	ApprovedBy string `json:"approved_by"`
	Priority   int    `json:"priority"`
	Tags       string `json:"tags"`

	// 时间戳
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// PromotionDetail 优惠码详情实体（包含扩展信息）
type PromotionDetail struct {
	*Promotion

	// 扩展时间信息
	ValidHours      string `json:"valid_hours,omitempty"`
	ValidWeekdays   string `json:"valid_weekdays,omitempty"`
	DailyUsageLimit int    `json:"daily_usage_limit,omitempty"`

	// 扩展地理限制
	ValidServiceAreas string `json:"valid_service_areas,omitempty"`
	ExcludedAreas     string `json:"excluded_areas,omitempty"`

	// 扩展用户限制
	TargetUserSegment string `json:"target_user_segment,omitempty"`
	MinUserLevel      int    `json:"min_user_level,omitempty"`
	ExcludedUserIDs   string `json:"excluded_user_ids,omitempty"`

	// 扩展订单限制
	ValidServiceTypes string  `json:"valid_service_types,omitempty"`
	MinDistance       float64 `json:"min_distance,omitempty"`
	MaxDistance       float64 `json:"max_distance,omitempty"`
	MinDuration       int     `json:"min_duration,omitempty"`

	// 组合限制
	CombinableWithOther bool   `json:"combinable_with_other,omitempty"`
	ExclusiveWith       string `json:"exclusive_with,omitempty"`
	RequiredPromotions  string `json:"required_promotions,omitempty"`

	// 发放设置
	IssueType      string `json:"issue_type,omitempty"`
	AutoIssueRule  string `json:"auto_issue_rule,omitempty"`
	IssueStartDate int64  `json:"issue_start_date,omitempty"`
	IssueEndDate   int64  `json:"issue_end_date,omitempty"`

	// 显示设置
	IsAutomatic     bool    `json:"is_automatic,omitempty"`
	RequiresCode    bool    `json:"requires_code,omitempty"`
	Weight          float64 `json:"weight,omitempty"`
	SortOrder       int     `json:"sort_order,omitempty"`
	ShowInList      bool    `json:"show_in_list,omitempty"`
	ShowDescription bool    `json:"show_description,omitempty"`
	IconURL         string  `json:"icon_url,omitempty"`
	BannerURL       string  `json:"banner_url,omitempty"`
	ColorScheme     string  `json:"color_scheme,omitempty"`

	// 营销信息
	MarketingText  string `json:"marketing_text,omitempty"`
	ShareText      string `json:"share_text,omitempty"`
	ShareURL       string `json:"share_url,omitempty"`
	LandingPageURL string `json:"landing_page_url,omitempty"`

	// 渠道限制
	ValidChannels    string `json:"valid_channels,omitempty"`
	ValidPlatforms   string `json:"valid_platforms,omitempty"`
	ValidAppVersions string `json:"valid_app_versions,omitempty"`

	// 分析追踪
	CampaignID      string `json:"campaign_id,omitempty"`
	SourceChannel   string `json:"source_channel,omitempty"`
	ReferrerID      string `json:"referrer_id,omitempty"`
	AttributionData string `json:"attribution_data,omitempty"`

	// 扩展统计信息
	ShareCount int `json:"share_count,omitempty"`

	// 成本和收益
	CostPerUse   float64 `json:"cost_per_use,omitempty"`
	TotalBudget  float64 `json:"total_budget,omitempty"`
	UsedBudget   float64 `json:"used_budget,omitempty"`
	EstimatedROI float64 `json:"estimated_roi,omitempty"`

	// A/B测试
	ExperimentID string `json:"experiment_id,omitempty"`
	VariantID    string `json:"variant_id,omitempty"`
	TestGroup    string `json:"test_group,omitempty"`
	ControlGroup bool   `json:"control_group,omitempty"`

	// 安全设置
	SecurityLevel  string  `json:"security_level,omitempty"`
	AntiAbuseRules string  `json:"anti_abuse_rules,omitempty"`
	RiskScore      float64 `json:"risk_score,omitempty"`

	// 通知设置
	NotifyOnClaim    bool `json:"notify_on_claim,omitempty"`
	NotifyOnUse      bool `json:"notify_on_use,omitempty"`
	NotifyOnExpiry   bool `json:"notify_on_expiry,omitempty"`
	ExpiryNotifyDays int  `json:"expiry_notify_days,omitempty"`

	// 创建者信息
	CreatorType string `json:"creator_type,omitempty"`
	CreatorID   string `json:"creator_id,omitempty"`

	// 审批信息
	ApprovedAt    int64  `json:"approved_at,omitempty"`
	ApprovalNotes string `json:"approval_notes,omitempty"`

	// 时间记录
	LastUsedAt    int64 `json:"last_used_at,omitempty"`
	LastClaimedAt int64 `json:"last_claimed_at,omitempty"`
	LastViewedAt  int64 `json:"last_viewed_at,omitempty"`
	ActivatedAt   int64 `json:"activated_at,omitempty"`
	SuspendedAt   int64 `json:"suspended_at,omitempty"`

	// 扩展信息
	CustomFields  string `json:"custom_fields,omitempty"`
	InternalNotes string `json:"internal_notes,omitempty"`
	ExternalNotes string `json:"external_notes,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
}

// PromotionUsageStats 优惠码使用统计
type PromotionUsageStats struct {
	PromotionID     string       `json:"promotion_id"`
	Code            string       `json:"code"`
	Title           string       `json:"title"`
	TotalUsage      int          `json:"total_usage"`
	UniqueUsers     int          `json:"unique_users"`
	TotalDiscount   float64      `json:"total_discount"`
	AverageDiscount float64      `json:"average_discount"`
	ConversionRate  float64      `json:"conversion_rate"`
	UsageByDate     []DailyUsage `json:"usage_by_date,omitempty"`
	TopUsers        []UserUsage  `json:"top_users,omitempty"`
}

// DailyUsage 每日使用统计
type DailyUsage struct {
	Date        string  `json:"date"` // YYYY-MM-DD
	Count       int     `json:"count"`
	Discount    float64 `json:"discount"`
	UniqueUsers int     `json:"unique_users"`
}

// UserUsage 用户使用统计
type UserUsage struct {
	UserID        string  `json:"user_id"`
	Username      string  `json:"username"`
	UsageCount    int     `json:"usage_count"`
	TotalDiscount float64 `json:"total_discount"`
	LastUsedAt    int64   `json:"last_used_at"`
}

// UserPromotion 用户优惠券实体（用户领取的优惠券）
type UserPromotion struct {
	// 基础信息
	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	PromotionID string `json:"promotion_id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`

	// 优惠信息
	DiscountType      string  `json:"discount_type"` // percentage, fixed_amount
	DiscountValue     float64 `json:"discount_value"`
	MaxDiscountAmount float64 `json:"max_discount_amount"`
	MinOrderAmount    float64 `json:"min_order_amount"`

	// 使用状态
	Status     string  `json:"status"`      // available, used, expired
	IsUsed     bool    `json:"is_used"`     // 是否已使用
	UsedAt     int64   `json:"used_at"`     // 使用时间
	UsedAmount float64 `json:"used_amount"` // 使用的优惠金额
	OrderID    string  `json:"order_id"`    // 使用的订单ID

	// 有效期
	ExpiredAt int64 `json:"expired_at"` // 过期时间

	// 来源信息
	Source     string `json:"source"`      // 来源：system, admin, event, referral
	SourceID   string `json:"source_id"`   // 来源ID
	SourceDesc string `json:"source_desc"` // 来源描述
	IssuedBy   string `json:"issued_by"`   // 发放者

	// 时间戳
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// UserPromotionList 用户优惠券列表响应
type UserPromotionList struct {
	Total      int              `json:"total"`
	Available  int              `json:"available"` // 可用数量
	Used       int              `json:"used"`      // 已使用数量
	Expired    int              `json:"expired"`   // 已过期数量
	Promotions []*UserPromotion `json:"promotions"`
}

// UserPromotionStats 用户优惠券统计
type UserPromotionStats struct {
	UserID          string  `json:"user_id"`
	TotalReceived   int     `json:"total_received"`   // 总共收到
	TotalUsed       int     `json:"total_used"`       // 总共使用
	TotalExpired    int     `json:"total_expired"`    // 总共过期
	TotalSaved      float64 `json:"total_saved"`      // 总共节省金额
	Available       int     `json:"available"`        // 当前可用
	MostUsedType    string  `json:"most_used_type"`   // 最常使用的类型
	AverageDiscount float64 `json:"average_discount"` // 平均优惠金额
	LastUsedAt      int64   `json:"last_used_at"`     // 最后使用时间
}

// PromotionUsageDetail 优惠券使用详情
type PromotionUsageDetail struct {
	UserPromotionID int64   `json:"user_promotion_id"`
	UserID          string  `json:"user_id"`
	OrderID         string  `json:"order_id"`
	Code            string  `json:"code"`
	DiscountAmount  float64 `json:"discount_amount"`
	OriginalAmount  float64 `json:"original_amount"`
	PaymentAmount   float64 `json:"payment_amount"`
	UsedAt          int64   `json:"used_at"`
}
