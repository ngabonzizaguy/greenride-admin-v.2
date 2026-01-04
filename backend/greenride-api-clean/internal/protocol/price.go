package protocol

// VehicleFilter 车辆筛选条件 - 精确匹配车辆类别和服务级别组合
type VehicleFilter struct {
	Category string `json:"category"` // sedan, suv, mpv, van, hatchback
	Level    string `json:"level"`    // economy, comfort, premium, luxury
}

// OrderPrice 价格实体 - 用于价格预估响应和价格管理
type OrderPrice struct {
	// 基础价格信息
	Currency        string  `json:"currency"`         // 货币
	Distance        float64 `json:"distance"`         // 距离（公里）
	Duration        int     `json:"duration"`         // 预计用时（分钟）
	OrderType       string  `json:"order_type"`       // 订单类型
	VehicleCategory string  `json:"vehicle_category"` // 车辆类别
	VehicleLevel    string  `json:"vehicle_level"`    // 车辆级别

	// 价格分解（扁平化）
	*PriceBreakdown
	// 价格锁定信息（扁平化）
	PriceID   string `json:"price_id"`   // 价格ID
	ExpiresAt int64  `json:"expires_at"` // 过期时间戳
	IsLocked  bool   `json:"is_locked"`  // 是否已锁定

	// 计算元信息（扁平化）
	CalculationTime int64  `json:"calculation_time_ms"` // 计算耗时（毫秒）
	RulesEvaluated  int    `json:"rules_evaluated"`     // 评估的规则数量
	RulesApplied    int    `json:"rules_applied"`       // 应用的规则数量
	EngineVersion   string `json:"engine_version"`      // 引擎版本

	// 价格明细分解
	Breakdowns []*PriceRuleResult `json:"breakdowns,omitempty"` // 价格明细分解
}

// PriceBreakdown 价格明细
type PriceBreakdown struct {
	BaseFare          float64 `json:"base_fare"`           // 起步价
	DistanceFare      float64 `json:"distance_fare"`       // 里程费
	TimeFare          float64 `json:"time_fare"`           // 时长费
	ServiceFee        float64 `json:"service_fee"`         // 服务费
	SurgeFare         float64 `json:"surge_fare"`          // 高峰期费用
	DiscountAmount    float64 `json:"discount_amount"`     // 总折扣金额
	PromoDiscount     float64 `json:"promo_discount"`      // 优惠码折扣金额
	UserPromoDiscount float64 `json:"user_promo_discount"` // 用户优惠券折扣金额
	OriginalFare      float64 `json:"original_fare"`       // 优惠前原始费用 (所有费用项的总和)
	DiscountedFare    float64 `json:"discounted_fare"`     // 优惠后折扣费用 (应用所有优惠后的价格)
}

// UpdatePriceRuleRequest 更新价格规则请求
type UpdatePriceRuleRequest struct {
	UserID          string           `json:"user_id"`                    // 更新者用户ID
	RuleID          string           `json:"rule_id" binding:"required"` // 价格规则ID (对应 models.PriceRule.RuleID)
	RuleName        *string          `json:"rule_name"`
	Description     *string          `json:"description"`
	Category        *string          `json:"category"`
	RuleType        *string          `json:"rule_type"`
	DiscountAmount  *float64         `json:"discount_amount"`
	DiscountPercent *float64         `json:"discount_percent"`
	SurgeMultiplier *float64         `json:"surge_multiplier"`
	MinimumFare     *float64         `json:"minimum_fare"`
	MaximumFare     *float64         `json:"maximum_fare"`
	Priority        *int             `json:"priority"`
	Status          *string          `json:"status"`
	StartDate       *int64           `json:"start_date"`
	EndDate         *int64           `json:"end_date"`
	VehicleFilters  []*VehicleFilter `json:"vehicle_filters"` // 车辆筛选条件数组
	ServiceAreas    []string         `json:"service_areas"`
	UserCategories  []string         `json:"user_categories"`
	ApplicableRides []string         `json:"applicable_rides"`
	MaxUsagePerUser *int             `json:"max_usage_per_user"`
	MaxUsageTotal   *int             `json:"max_usage_total"`
	PromoCode       *string          `json:"promo_code"`
	RequiresCode    *int             `json:"requires_code"` // 0:否 1:是
	Metadata        map[string]any   `json:"metadata"`
}

// TimeSlot 时间段结构
type TimeSlot struct {
	StartHour   int `json:"start_hour"`   // 开始小时
	StartMinute int `json:"start_minute"` // 开始分钟
	EndHour     int `json:"end_hour"`     // 结束小时
	EndMinute   int `json:"end_minute"`   // 结束分钟
}

// TierConfig 阶梯配置结构 - 匹配数据库JSON格式
type TierConfig struct {
	MinDuration *float64 `json:"min_duration"` // 最小时长（秒），使用指针以支持null
	MaxDuration *float64 `json:"max_duration"` // 最大时长（秒），null表示无上限
	Rate        float64  `json:"rate"`         // 费率
	Description string   `json:"description"`  // 描述
}

// TieredRuleConfig 阶梯规则配置 - 匹配数据库JSON格式
type TieredRuleConfig struct {
	Unit              string        `json:"unit"`               // 单位: minute, kilometer 等
	Tiers             []*TierConfig `json:"tiers"`              // 阶梯配置数组
	CalculationMethod string        `json:"calculation_method"` // 计算方法: step, progressive 等
}

// DynamicFactors 动态因素配置
type DynamicFactors struct {
	WeatherConditions map[string]float64 `json:"weather_conditions"` // 天气系数 {"rain": 1.2, "snow": 1.5}
	TimeOfDay         map[string]float64 `json:"time_of_day"`        // 时段系数 {"peak": 1.3, "off_peak": 0.9}
	EventTypes        map[string]float64 `json:"event_types"`        // 事件系数 {"holiday": 1.4, "concert": 1.6}
}
