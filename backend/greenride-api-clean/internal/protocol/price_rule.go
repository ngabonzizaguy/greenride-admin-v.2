package protocol

// PriceRule 价格规则对外输出结构
type PriceRule struct {
	ID     int64  `json:"id,omitempty"`
	RuleID string `json:"rule_id,omitempty"`

	// 基本信息
	RuleName    string `json:"rule_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"` // base_pricing, surge_pricing, discount, promotion, special_offer

	// 规则类型
	RuleType     string `json:"rule_type,omitempty"`     // percentage, fixed_amount, multiplier, tiered, custom
	DiscountType string `json:"discount_type,omitempty"` // percentage, fixed, buy_x_get_y, free_delivery
	PricingModel string `json:"pricing_model,omitempty"` // distance_based, time_based, fixed_rate, dynamic

	// 适用范围
	VehicleTypes    []string `json:"vehicle_types,omitempty"`    // 适用车型数组
	ServiceAreas    []string `json:"service_areas,omitempty"`    // 适用区域数组
	UserCategories  []string `json:"user_categories,omitempty"`  // 用户类型数组 (new, premium, regular)
	ApplicableRides []string `json:"applicable_rides,omitempty"` // 适用订单类型数组

	// 时间范围
	StartedAt int64  `json:"started_at,omitempty"` // 生效开始时间
	EndedAt   int64  `json:"ended_at,omitempty"`   // 生效结束时间
	TimeZone  string `json:"timezone,omitempty"`   // 时区

	// 时间限制
	DayOfWeek     string      `json:"day_of_week,omitempty"`    // 星期限制：1,2,3,4,5,6,7
	TimeSlots     []*TimeSlot `json:"time_slots,omitempty"`     // 时间段限制数组
	ExcludedDates []string    `json:"excluded_dates,omitempty"` // 排除日期数组 ["2024-01-01", "2024-12-25"]
	IncludedDates []string    `json:"included_dates,omitempty"` // 特定日期数组

	// 价格计算
	BaseRate      float64 `json:"base_rate,omitempty"`       // 基础价格
	PerKmRate     float64 `json:"per_km_rate,omitempty"`     // 每公里价格
	PerMinuteRate float64 `json:"per_minute_rate,omitempty"` // 每分钟价格
	MinimumFare   float64 `json:"minimum_fare,omitempty"`    // 最低收费
	MaximumFare   float64 `json:"maximum_fare,omitempty"`    // 最高收费

	// 折扣/加价
	DiscountAmount  float64 `json:"discount_amount,omitempty"`  // 固定折扣金额
	DiscountPercent float64 `json:"discount_percent,omitempty"` // 折扣百分比
	SurgeMultiplier float64 `json:"surge_multiplier,omitempty"` // 涌潮加价倍数
	MaxDiscount     float64 `json:"max_discount,omitempty"`     // 最大折扣金额
	MinOrderAmount  float64 `json:"min_order_amount,omitempty"` // 最小订单金额

	// 阶梯定价
	TieredRules *TieredRuleConfig `json:"tiered_rules,omitempty"` // 阶梯价格规则配置

	// 动态定价参数
	DemandFactor   float64         `json:"demand_factor,omitempty"`   // 需求系数
	SupplyFactor   float64         `json:"supply_factor,omitempty"`   // 供给系数
	WeatherFactor  float64         `json:"weather_factor,omitempty"`  // 天气系数
	EventFactor    float64         `json:"event_factor,omitempty"`    // 事件系数
	DynamicFactors *DynamicFactors `json:"dynamic_factors,omitempty"` // 动态因素配置

	// 使用限制
	MaxUsagePerUser int `json:"max_usage_per_user,omitempty"` // 每用户最大使用次数
	MaxUsagePerDay  int `json:"max_usage_per_day,omitempty"`  // 每日最大使用次数
	MaxUsageTotal   int `json:"max_usage_total,omitempty"`    // 总最大使用次数
	UsageCount      int `json:"usage_count,omitempty"`        // 已使用次数

	// 条件限制
	MinDistance float64 `json:"min_distance,omitempty"` // 最小距离
	MaxDistance float64 `json:"max_distance,omitempty"` // 最大距离
	MinDuration int     `json:"min_duration,omitempty"` // 最小时长(分钟)
	MaxDuration int     `json:"max_duration,omitempty"` // 最大时长(分钟)

	// 组合规则
	StackableRules []string `json:"stackable_rules,omitempty"` // 可叠加的规则ID数组
	ExclusiveRules []string `json:"exclusive_rules,omitempty"` // 互斥的规则ID数组
	Priority       int      `json:"priority,omitempty"`        // 优先级(数字越小优先级越高)

	// 状态管理
	Status       string `json:"status,omitempty"`        // draft, active, paused, expired, deleted
	IsGlobal     bool   `json:"is_global,omitempty"`     // 是否全局规则
	AutoApply    bool   `json:"auto_apply,omitempty"`    // 是否自动应用
	RequiresCode bool   `json:"requires_code,omitempty"` // 是否需要优惠码

	// 促销码相关
	PromoCode     string `json:"promo_code,omitempty"`     // 促销码
	CaseSensitive bool   `json:"case_sensitive,omitempty"` // 促销码是否区分大小写

	// 审批信息
	CreatedBy     string `json:"created_by,omitempty"`     // 创建人
	ApprovedBy    string `json:"approved_by,omitempty"`    // 审批人
	ApprovedAt    int64  `json:"approved_at,omitempty"`    // 审批时间
	ApprovalNotes string `json:"approval_notes,omitempty"` // 审批备注

	// 统计信息
	ViewCount     int     `json:"view_count,omitempty"`     // 查看次数
	ClickCount    int     `json:"click_count,omitempty"`    // 点击次数
	UsageToday    int     `json:"usage_today,omitempty"`    // 今日使用次数
	RevenueImpact float64 `json:"revenue_impact,omitempty"` // 收入影响
	CostSaved     float64 `json:"cost_saved,omitempty"`     // 为用户节省的费用

	// 元数据
	Metadata map[string]any `json:"metadata,omitempty"` // 附加元数据对象
	Tags     []string       `json:"tags,omitempty"`     // 标签数组
	Notes    string         `json:"notes,omitempty"`    // 备注

	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
}
