package models

import (
	"fmt"
	"strings"

	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// VehicleFilter 车辆筛选条件 - 精确匹配车辆类别和服务级别组合
type VehicleFilter struct {
	Category string `json:"category"` // sedan, suv, mpv, van, hatchback
	Level    string `json:"level"`    // economy, comfort, premium, luxury
}

// PriceRule 价格规则表 - 动态定价、优惠券、折扣等规则配置
type PriceRule struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RuleID string `json:"rule_id" gorm:"column:rule_id;type:varchar(64);uniqueIndex"`
	Salt   string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PriceRuleValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PriceRuleValues struct {
	// 基本信息
	RuleName    *string `json:"rule_name" gorm:"column:rule_name;type:varchar(255);index"`
	DisplayName *string `json:"display_name" gorm:"column:display_name;type:varchar(255)"`
	Description *string `json:"description" gorm:"column:description;type:text"`
	Category    *string `json:"category" gorm:"column:category;type:varchar(100);index"` // base_pricing, surge_pricing, discount, promotion, special_offer

	// 规则类型
	RuleType     *string `json:"rule_type" gorm:"column:rule_type;type:varchar(50);index"`   // percentage, fixed_amount, multiplier, tiered, custom
	DiscountType *string `json:"discount_type" gorm:"column:discount_type;type:varchar(50)"` // percentage, fixed, buy_x_get_y, free_delivery
	PricingModel *string `json:"pricing_model" gorm:"column:pricing_model;type:varchar(50)"` // distance_based, time_based, fixed_rate, dynamic

	// 适用范围
	VehicleFilters  []*VehicleFilter `json:"vehicle_filters" gorm:"column:vehicle_filters;type:json;serializer:json"`   // 车辆筛选条件数组
	ServiceAreas    []string         `json:"service_areas" gorm:"column:service_areas;type:json;serializer:json"`       // 适用区域数组
	UserCategories  []string         `json:"user_categories" gorm:"column:user_categories;type:json;serializer:json"`   // 用户类型数组 (new, premium, regular)
	ApplicableRides []string         `json:"applicable_rides" gorm:"column:applicable_rides;type:json;serializer:json"` // 适用订单类型数组

	// 时间范围
	StartedAt *int64  `json:"started_at" gorm:"column:started_at"`                            // 生效开始时间
	EndedAt   *int64  `json:"ended_at" gorm:"column:ended_at"`                                // 生效结束时间
	TimeZone  *string `json:"timezone" gorm:"column:timezone;type:varchar(50);default:'UTC'"` // 时区

	// 时间限制
	DayOfWeek     *string              `json:"day_of_week" gorm:"column:day_of_week;type:varchar(20)"`                // 星期限制：1,2,3,4,5,6,7
	TimeSlots     []*protocol.TimeSlot `json:"time_slots" gorm:"column:time_slots;type:json;serializer:json"`         // 时间段限制数组
	ExcludedDates []string             `json:"excluded_dates" gorm:"column:excluded_dates;type:json;serializer:json"` // 排除日期数组 ["2024-01-01", "2024-12-25"]
	IncludedDates []string             `json:"included_dates" gorm:"column:included_dates;type:json;serializer:json"` // 特定日期数组

	// 价格计算
	BaseRate      *float64 `json:"base_rate" gorm:"column:base_rate;type:decimal(10,2)"`            // 基础价格
	PerKmRate     *float64 `json:"per_km_rate" gorm:"column:per_km_rate;type:decimal(8,2)"`         // 每公里价格
	PerMinuteRate *float64 `json:"per_minute_rate" gorm:"column:per_minute_rate;type:decimal(6,2)"` // 每分钟价格
	MinimumFare   *float64 `json:"minimum_fare" gorm:"column:minimum_fare;type:decimal(10,2)"`      // 最低收费
	MaximumFare   *float64 `json:"maximum_fare" gorm:"column:maximum_fare;type:decimal(10,2)"`      // 最高收费

	// 折扣/加价
	DiscountAmount  *float64 `json:"discount_amount" gorm:"column:discount_amount;type:decimal(10,2)"`   // 固定折扣金额
	DiscountPercent *float64 `json:"discount_percent" gorm:"column:discount_percent;type:decimal(5,2)"`  // 折扣百分比
	SurgeMultiplier *float64 `json:"surge_multiplier" gorm:"column:surge_multiplier;type:decimal(4,2)"`  // 涌潮加价倍数
	MaxDiscount     *float64 `json:"max_discount" gorm:"column:max_discount;type:decimal(10,2)"`         // 最大折扣金额
	MinOrderAmount  *float64 `json:"min_order_amount" gorm:"column:min_order_amount;type:decimal(10,2)"` // 最小订单金额

	// 阶梯定价
	TieredRules *protocol.TieredRuleConfig `json:"tiered_rules" gorm:"column:tiered_rules;type:json;serializer:json"` // 阶梯价格规则配置

	// 动态定价参数
	DemandFactor   *float64                 `json:"demand_factor" gorm:"column:demand_factor;type:decimal(4,2);default:1.0"`   // 需求系数
	SupplyFactor   *float64                 `json:"supply_factor" gorm:"column:supply_factor;type:decimal(4,2);default:1.0"`   // 供给系数
	WeatherFactor  *float64                 `json:"weather_factor" gorm:"column:weather_factor;type:decimal(4,2);default:1.0"` // 天气系数
	EventFactor    *float64                 `json:"event_factor" gorm:"column:event_factor;type:decimal(4,2);default:1.0"`     // 事件系数
	DynamicFactors *protocol.DynamicFactors `json:"dynamic_factors" gorm:"column:dynamic_factors;type:json;serializer:json"`   // 动态因素配置

	// 使用限制
	MaxUsagePerUser *int `json:"max_usage_per_user" gorm:"column:max_usage_per_user;type:int"` // 每用户最大使用次数
	MaxUsagePerDay  *int `json:"max_usage_per_day" gorm:"column:max_usage_per_day;type:int"`   // 每日最大使用次数
	MaxUsageTotal   *int `json:"max_usage_total" gorm:"column:max_usage_total;type:int"`       // 总最大使用次数
	UsageCount      *int `json:"usage_count" gorm:"column:usage_count;type:int;default:0"`     // 已使用次数

	// 条件限制
	MinDistance *float64 `json:"min_distance" gorm:"column:min_distance;type:decimal(8,2)"` // 最小距离
	MaxDistance *float64 `json:"max_distance" gorm:"column:max_distance;type:decimal(8,2)"` // 最大距离
	MinDuration *int     `json:"min_duration" gorm:"column:min_duration;type:int"`          // 最小时长(分钟)
	MaxDuration *int     `json:"max_duration" gorm:"column:max_duration;type:int"`          // 最大时长(分钟)

	// 组合规则
	StackableRules []string `json:"stackable_rules" gorm:"column:stackable_rules;type:json;serializer:json"` // 可叠加的规则ID数组
	ExclusiveRules []string `json:"exclusive_rules" gorm:"column:exclusive_rules;type:json;serializer:json"` // 互斥的规则ID数组
	Priority       *int     `json:"priority" gorm:"column:priority;type:int;default:100"`                    // 优先级(数字越小优先级越高)

	// 状态管理
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'draft'"` // draft, active, paused, expired, deleted
	IsGlobal     *int    `json:"is_global" gorm:"column:is_global;type:tinyint;default:0"`           // 是否全局规则 0:否 1:是
	AutoApply    *int    `json:"auto_apply" gorm:"column:auto_apply;type:tinyint;default:0"`         // 是否自动应用 0:否 1:是
	RequiresCode *int    `json:"requires_code" gorm:"column:requires_code;type:tinyint;default:0"`   // 是否需要优惠码 0:否 1:是

	// 促销码相关
	PromoCode     *string `json:"promo_code" gorm:"column:promo_code;type:varchar(50);index"`         // 促销码
	CaseSensitive *int    `json:"case_sensitive" gorm:"column:case_sensitive;type:tinyint;default:0"` // 促销码是否区分大小写 0:否 1:是

	// 审批信息
	CreatedBy     *string `json:"created_by" gorm:"column:created_by;type:varchar(64)"`   // 创建人
	ApprovedBy    *string `json:"approved_by" gorm:"column:approved_by;type:varchar(64)"` // 审批人
	ApprovedAt    *int64  `json:"approved_at" gorm:"column:approved_at"`                  // 审批时间
	ApprovalNotes *string `json:"approval_notes" gorm:"column:approval_notes;type:text"`  // 审批备注

	// 统计信息
	ViewCount     *int     `json:"view_count" gorm:"column:view_count;type:int;default:0"`                      // 查看次数
	ClickCount    *int     `json:"click_count" gorm:"column:click_count;type:int;default:0"`                    // 点击次数
	UsageToday    *int     `json:"usage_today" gorm:"column:usage_today;type:int;default:0"`                    // 今日使用次数
	RevenueImpact *float64 `json:"revenue_impact" gorm:"column:revenue_impact;type:decimal(15,2);default:0.00"` // 收入影响
	CostSaved     *float64 `json:"cost_saved" gorm:"column:cost_saved;type:decimal(15,2);default:0.00"`         // 为用户节省的费用

	// 元数据
	Metadata map[string]any `json:"metadata" gorm:"column:metadata;type:json;serializer:json"` // 附加元数据对象
	Tags     []string       `json:"tags" gorm:"column:tags;type:json;serializer:json"`         // 标签数组
	Notes    *string        `json:"notes" gorm:"column:notes;type:text"`                       // 备注

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type PriceRules []*PriceRule

func (PriceRule) TableName() string {
	return "t_price_rules"
}

// NewPriceRule 创建新的价格规则对象
func NewPriceRule() *PriceRule {
	return &PriceRule{
		RuleID: utils.GeneratePriceRuleID(),
		Salt:   utils.GenerateSalt(),
		PriceRuleValues: &PriceRuleValues{
			Category:      utils.StringPtr(protocol.PriceRuleCategoryDiscount),
			RuleType:      utils.StringPtr(protocol.PriceRuleTypePercentage),
			PricingModel:  utils.StringPtr(protocol.PricingModelDistanceBased),
			TimeZone:      utils.StringPtr("UTC"),
			Status:        utils.StringPtr(protocol.StatusDraft),
			IsGlobal:      utils.IntPtr(0),
			AutoApply:     utils.IntPtr(0),
			RequiresCode:  utils.IntPtr(0),
			CaseSensitive: utils.IntPtr(0),
			Priority:      utils.IntPtr(100),
			UsageCount:    utils.IntPtr(0),
			ViewCount:     utils.IntPtr(0),
			ClickCount:    utils.IntPtr(0),
			UsageToday:    utils.IntPtr(0),
			RevenueImpact: utils.Float64Ptr(0.00),
			CostSaved:     utils.Float64Ptr(0.00),
			DemandFactor:  utils.Float64Ptr(1.0),
			SupplyFactor:  utils.Float64Ptr(1.0),
			WeatherFactor: utils.Float64Ptr(1.0),
			EventFactor:   utils.Float64Ptr(1.0),
		},
	}
}

// SetValues 更新PriceRuleV2Values中的非nil值
func (p *PriceRuleValues) SetValues(values *PriceRuleValues) {
	if values == nil {
		return
	}

	if values.RuleName != nil {
		p.RuleName = values.RuleName
	}
	if values.DisplayName != nil {
		p.DisplayName = values.DisplayName
	}
	if values.Description != nil {
		p.Description = values.Description
	}
	if values.Category != nil {
		p.Category = values.Category
	}
	if values.RuleType != nil {
		p.RuleType = values.RuleType
	}
	if values.Status != nil {
		p.Status = values.Status
	}
	if values.BaseRate != nil {
		p.BaseRate = values.BaseRate
	}
	if values.DiscountPercent != nil {
		p.DiscountPercent = values.DiscountPercent
	}
	if values.PromoCode != nil {
		p.PromoCode = values.PromoCode
	}
	if values.VehicleFilters != nil {
		p.VehicleFilters = values.VehicleFilters
	}
	if values.Notes != nil {
		p.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		p.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (p *PriceRuleValues) GetRuleName() string {
	if p.RuleName == nil {
		return ""
	}
	return *p.RuleName
}

func (p *PriceRuleValues) GetDisplayName() string {
	if p.DisplayName == nil {
		return p.GetRuleName()
	}
	return *p.DisplayName
}

func (p *PriceRuleValues) GetCategory() string {
	if p.Category == nil {
		return protocol.PriceRuleCategoryDiscount
	}
	return *p.Category
}

func (p *PriceRuleValues) GetRuleType() string {
	if p.RuleType == nil {
		return protocol.PriceRuleTypePercentage
	}
	return *p.RuleType
}

func (p *PriceRuleValues) GetPricingModel() string {
	if p.PricingModel == nil {
		return protocol.PricingModelDistanceBased
	}
	return *p.PricingModel
}

func (p *PriceRuleValues) GetStatus() string {
	if p.Status == nil {
		return protocol.StatusDraft
	}
	return *p.Status
}

func (p *PriceRuleValues) GetBaseRate() float64 {
	if p.BaseRate == nil {
		return 0.0
	}
	return *p.BaseRate
}

func (p *PriceRuleValues) GetDiscountPercent() float64 {
	if p.DiscountPercent == nil {
		return 0.0
	}
	return *p.DiscountPercent
}

func (p *PriceRuleValues) GetDiscountAmount() float64 {
	if p.DiscountAmount == nil {
		return 0.0
	}
	return *p.DiscountAmount
}

func (p *PriceRuleValues) GetMinimumFare() float64 {
	if p.MinimumFare == nil {
		return 0.0
	}
	return *p.MinimumFare
}

func (p *PriceRuleValues) GetMaximumFare() float64 {
	if p.MaximumFare == nil {
		return 999999.0
	}
	return *p.MaximumFare
}

func (p *PriceRuleValues) GetSurgeMultiplier() float64 {
	if p.SurgeMultiplier == nil {
		return 1.0
	}
	return *p.SurgeMultiplier
}

func (p *PriceRuleValues) GetUsageCount() int {
	if p.UsageCount == nil {
		return 0
	}
	return *p.UsageCount
}

func (p *PriceRuleValues) GetMaxUsageTotal() int {
	if p.MaxUsageTotal == nil {
		return 999999
	}
	return *p.MaxUsageTotal
}

func (p *PriceRuleValues) GetPriority() int {
	if p.Priority == nil {
		return 100
	}
	return *p.Priority
}

func (p *PriceRuleValues) GetIsGlobal() bool {
	if p.IsGlobal == nil {
		return false
	}
	return *p.IsGlobal == 1
}

func (p *PriceRuleValues) GetAutoApply() bool {
	if p.AutoApply == nil {
		return false
	}
	return *p.AutoApply == 1
}

func (p *PriceRuleValues) GetRequiresCode() bool {
	if p.RequiresCode == nil {
		return false
	}
	return *p.RequiresCode == 1
}

func (p *PriceRuleValues) GetDescription() string {
	if p.Description == nil {
		return ""
	}
	return *p.Description
}

func (p *PriceRuleValues) GetDiscountType() string {
	if p.DiscountType == nil {
		return protocol.DiscountTypePercentage
	}
	return *p.DiscountType
}

func (p *PriceRuleValues) GetVehicleFilters() []*VehicleFilter {
	if p.VehicleFilters == nil {
		return []*VehicleFilter{}
	}
	return p.VehicleFilters
}

func (p *PriceRuleValues) GetServiceAreas() []string {
	if p.ServiceAreas == nil {
		return []string{}
	}
	return p.ServiceAreas
}

func (p *PriceRuleValues) GetUserCategories() []string {
	if p.UserCategories == nil {
		return []string{}
	}
	return p.UserCategories
}

func (p *PriceRuleValues) GetApplicableRides() []string {
	if p.ApplicableRides == nil {
		return []string{}
	}
	return p.ApplicableRides
}

func (p *PriceRuleValues) GetStartedAt() int64 {
	if p.StartedAt == nil {
		return 0
	}
	return *p.StartedAt
}

func (p *PriceRuleValues) GetEndedAt() int64 {
	if p.EndedAt == nil {
		return 0
	}
	return *p.EndedAt
}

func (p *PriceRuleValues) GetTimeZone() string {
	if p.TimeZone == nil {
		return "UTC"
	}
	return *p.TimeZone
}

func (p *PriceRuleValues) GetDayOfWeek() string {
	if p.DayOfWeek == nil {
		return ""
	}
	return *p.DayOfWeek
}

func (p *PriceRuleValues) GetTimeSlots() []*protocol.TimeSlot {
	if p.TimeSlots == nil {
		return []*protocol.TimeSlot{}
	}
	return p.TimeSlots
}

func (p *PriceRuleValues) GetExcludedDates() []string {
	if p.ExcludedDates == nil {
		return []string{}
	}
	return p.ExcludedDates
}

func (p *PriceRuleValues) GetIncludedDates() []string {
	if p.IncludedDates == nil {
		return []string{}
	}
	return p.IncludedDates
}

func (p *PriceRuleValues) GetPerKmRate() float64 {
	if p.PerKmRate == nil {
		return 0.0
	}
	return *p.PerKmRate
}

func (p *PriceRuleValues) GetPerMinuteRate() float64 {
	if p.PerMinuteRate == nil {
		return 0.0
	}
	return *p.PerMinuteRate
}

func (p *PriceRuleValues) GetMaxDiscount() float64 {
	if p.MaxDiscount == nil {
		return 999999.0
	}
	return *p.MaxDiscount
}

func (p *PriceRuleValues) GetMinOrderAmount() float64 {
	if p.MinOrderAmount == nil {
		return 0.0
	}
	return *p.MinOrderAmount
}

func (p *PriceRuleValues) GetTieredRules() *protocol.TieredRuleConfig {
	return p.TieredRules
}

func (p *PriceRuleValues) GetDemandFactor() float64 {
	if p.DemandFactor == nil {
		return 1.0
	}
	return *p.DemandFactor
}

func (p *PriceRuleValues) GetSupplyFactor() float64 {
	if p.SupplyFactor == nil {
		return 1.0
	}
	return *p.SupplyFactor
}

func (p *PriceRuleValues) GetWeatherFactor() float64 {
	if p.WeatherFactor == nil {
		return 1.0
	}
	return *p.WeatherFactor
}

func (p *PriceRuleValues) GetEventFactor() float64 {
	if p.EventFactor == nil {
		return 1.0
	}
	return *p.EventFactor
}

func (p *PriceRuleValues) GetDynamicFactors() *protocol.DynamicFactors {
	return p.DynamicFactors
}

func (p *PriceRuleValues) GetMaxUsagePerUser() int {
	if p.MaxUsagePerUser == nil {
		return 999999
	}
	return *p.MaxUsagePerUser
}

func (p *PriceRuleValues) GetMaxUsagePerDay() int {
	if p.MaxUsagePerDay == nil {
		return 999999
	}
	return *p.MaxUsagePerDay
}

func (p *PriceRuleValues) GetMinDistance() float64 {
	if p.MinDistance == nil {
		return 0.0
	}
	return *p.MinDistance
}

func (p *PriceRuleValues) GetMaxDistance() float64 {
	if p.MaxDistance == nil {
		return 999999.0
	}
	return *p.MaxDistance
}

func (p *PriceRuleValues) GetMinDuration() int {
	if p.MinDuration == nil {
		return 0
	}
	return *p.MinDuration
}

func (p *PriceRuleValues) GetMaxDuration() int {
	if p.MaxDuration == nil {
		return 999999
	}
	return *p.MaxDuration
}

func (p *PriceRuleValues) GetStackableRules() []string {
	if p.StackableRules == nil {
		return []string{}
	}
	return p.StackableRules
}

func (p *PriceRuleValues) GetExclusiveRules() []string {
	if p.ExclusiveRules == nil {
		return []string{}
	}
	return p.ExclusiveRules
}

func (p *PriceRuleValues) GetCaseSensitive() bool {
	if p.CaseSensitive == nil {
		return false
	}
	return *p.CaseSensitive == 1
}

func (p *PriceRuleValues) GetPromoCode() string {
	if p.PromoCode == nil {
		return ""
	}
	return *p.PromoCode
}

func (p *PriceRuleValues) GetCreatedBy() string {
	if p.CreatedBy == nil {
		return ""
	}
	return *p.CreatedBy
}

func (p *PriceRuleValues) GetApprovedBy() string {
	if p.ApprovedBy == nil {
		return ""
	}
	return *p.ApprovedBy
}

func (p *PriceRuleValues) GetApprovedAt() int64 {
	if p.ApprovedAt == nil {
		return 0
	}
	return *p.ApprovedAt
}

func (p *PriceRuleValues) GetApprovalNotes() string {
	if p.ApprovalNotes == nil {
		return ""
	}
	return *p.ApprovalNotes
}

func (p *PriceRuleValues) GetViewCount() int {
	if p.ViewCount == nil {
		return 0
	}
	return *p.ViewCount
}

func (p *PriceRuleValues) GetClickCount() int {
	if p.ClickCount == nil {
		return 0
	}
	return *p.ClickCount
}

func (p *PriceRuleValues) GetUsageToday() int {
	if p.UsageToday == nil {
		return 0
	}
	return *p.UsageToday
}

func (p *PriceRuleValues) GetRevenueImpact() float64 {
	if p.RevenueImpact == nil {
		return 0.0
	}
	return *p.RevenueImpact
}

func (p *PriceRuleValues) GetCostSaved() float64 {
	if p.CostSaved == nil {
		return 0.0
	}
	return *p.CostSaved
}

func (p *PriceRuleValues) GetMetadata() map[string]any {
	if p.Metadata == nil {
		return make(map[string]any)
	}
	return p.Metadata
}

func (p *PriceRuleValues) GetTags() []string {
	if p.Tags == nil {
		return []string{}
	}
	return p.Tags
}

func (p *PriceRuleValues) GetNotes() string {
	if p.Notes == nil {
		return ""
	}
	return *p.Notes
}

// Setter 方法
func (p *PriceRuleValues) SetRuleName(name string) *PriceRuleValues {
	p.RuleName = &name
	return p
}

func (p *PriceRuleValues) SetDisplayName(name string) *PriceRuleValues {
	p.DisplayName = &name
	return p
}

func (p *PriceRuleValues) SetDescription(desc string) *PriceRuleValues {
	p.Description = &desc
	return p
}

func (p *PriceRuleValues) SetCategory(category string) *PriceRuleValues {
	p.Category = &category
	return p
}

func (p *PriceRuleValues) SetRuleType(ruleType string) *PriceRuleValues {
	p.RuleType = &ruleType
	return p
}

func (p *PriceRuleValues) SetPricingModel(model string) *PriceRuleValues {
	p.PricingModel = &model
	return p
}

func (p *PriceRuleValues) SetStatus(status string) *PriceRuleValues {
	p.Status = &status
	return p
}

func (p *PriceRuleValues) SetTimeRange(startedAt, endedAt int64) *PriceRuleValues {
	p.StartedAt = &startedAt
	p.EndedAt = &endedAt
	return p
}

func (p *PriceRuleValues) SetBasePricing(baseRate, perKmRate, perMinuteRate float64) *PriceRuleValues {
	p.BaseRate = &baseRate
	p.PerKmRate = &perKmRate
	p.PerMinuteRate = &perMinuteRate
	return p
}

func (p *PriceRuleValues) SetFareLimits(minimum, maximum float64) *PriceRuleValues {
	p.MinimumFare = &minimum
	p.MaximumFare = &maximum
	return p
}

func (p *PriceRuleValues) SetDiscountPercent(percent float64) *PriceRuleValues {
	p.DiscountPercent = &percent
	p.RuleType = utils.StringPtr(protocol.PriceRuleTypePercentage)
	return p
}

func (p *PriceRuleValues) SetDiscountAmount(amount float64) *PriceRuleValues {
	p.DiscountAmount = &amount
	p.RuleType = utils.StringPtr(protocol.PriceRuleTypeFixedAmount)
	return p
}

func (p *PriceRuleValues) SetSurgeMultiplier(multiplier float64) *PriceRuleValues {
	p.SurgeMultiplier = &multiplier
	p.Category = utils.StringPtr(protocol.PriceRuleCategorySurgePricing)
	p.RuleType = utils.StringPtr(protocol.PriceRuleTypeMultiplier)
	return p
}

func (p *PriceRuleValues) SetUsageLimits(perUser, perDay, total int) *PriceRuleValues {
	if perUser > 0 {
		p.MaxUsagePerUser = &perUser
	}
	if perDay > 0 {
		p.MaxUsagePerDay = &perDay
	}
	if total > 0 {
		p.MaxUsageTotal = &total
	}
	return p
}

func (p *PriceRuleValues) SetPriority(priority int) *PriceRuleValues {
	p.Priority = &priority
	return p
}

func (p *PriceRuleValues) SetGlobal(global bool) *PriceRuleValues {
	value := 0
	if global {
		value = 1
	}
	p.IsGlobal = &value
	return p
}

func (p *PriceRuleValues) SetAutoApply(autoApply bool) *PriceRuleValues {
	value := 0
	if autoApply {
		value = 1
	}
	p.AutoApply = &value
	return p
}

func (p *PriceRuleValues) SetPromoCode(code string, caseSensitive bool) *PriceRuleValues {
	p.PromoCode = &code
	caseValue := 0
	if caseSensitive {
		caseValue = 1
	}
	p.CaseSensitive = &caseValue
	p.RequiresCode = utils.IntPtr(1)
	return p
}

func (p *PriceRuleValues) SetDiscountType(discountType string) *PriceRuleValues {
	p.DiscountType = &discountType
	return p
}

func (p *PriceRuleValues) SetVehicleFilters(filters []*VehicleFilter) *PriceRuleValues {
	p.VehicleFilters = filters
	return p
}

// AddVehicleFilter 添加单个车辆筛选条件
func (p *PriceRuleValues) AddVehicleFilter(category, level string) *PriceRuleValues {
	filter := &VehicleFilter{
		Category: category,
		Level:    level,
	}

	if p.VehicleFilters == nil {
		p.VehicleFilters = []*VehicleFilter{}
	}

	p.VehicleFilters = append(p.VehicleFilters, filter)
	return p
}

func (p *PriceRuleValues) SetServiceAreas(serviceAreas []string) *PriceRuleValues {
	p.ServiceAreas = serviceAreas
	return p
}

func (p *PriceRuleValues) SetUserCategories(userCategories []string) *PriceRuleValues {
	p.UserCategories = userCategories
	return p
}

func (p *PriceRuleValues) SetApplicableRides(applicableRides []string) *PriceRuleValues {
	p.ApplicableRides = applicableRides
	return p
}

func (p *PriceRuleValues) SetTimeZone(timeZone string) *PriceRuleValues {
	p.TimeZone = &timeZone
	return p
}

func (p *PriceRuleValues) SetDayOfWeek(dayOfWeek string) *PriceRuleValues {
	p.DayOfWeek = &dayOfWeek
	return p
}

func (p *PriceRuleValues) SetTimeSlots(timeSlots []*protocol.TimeSlot) *PriceRuleValues {
	p.TimeSlots = timeSlots
	return p
}

func (p *PriceRuleValues) SetExcludedDates(excludedDates []string) *PriceRuleValues {
	p.ExcludedDates = excludedDates
	return p
}

func (p *PriceRuleValues) SetIncludedDates(includedDates []string) *PriceRuleValues {
	p.IncludedDates = includedDates
	return p
}

func (p *PriceRuleValues) SetPerKmRate(rate float64) *PriceRuleValues {
	p.PerKmRate = &rate
	return p
}

func (p *PriceRuleValues) SetPerMinuteRate(rate float64) *PriceRuleValues {
	p.PerMinuteRate = &rate
	return p
}

func (p *PriceRuleValues) SetMaxDiscount(maxDiscount float64) *PriceRuleValues {
	p.MaxDiscount = &maxDiscount
	return p
}

func (p *PriceRuleValues) SetMinOrderAmount(minAmount float64) *PriceRuleValues {
	p.MinOrderAmount = &minAmount
	return p
}

func (p *PriceRuleValues) SetTieredRules(tieredRules *protocol.TieredRuleConfig) *PriceRuleValues {
	p.TieredRules = tieredRules
	return p
}

func (p *PriceRuleValues) SetDemandFactor(factor float64) *PriceRuleValues {
	p.DemandFactor = &factor
	return p
}

func (p *PriceRuleValues) SetSupplyFactor(factor float64) *PriceRuleValues {
	p.SupplyFactor = &factor
	return p
}

func (p *PriceRuleValues) SetWeatherFactor(factor float64) *PriceRuleValues {
	p.WeatherFactor = &factor
	return p
}

func (p *PriceRuleValues) SetEventFactor(factor float64) *PriceRuleValues {
	p.EventFactor = &factor
	return p
}

func (p *PriceRuleValues) SetDynamicFactors(factors *protocol.DynamicFactors) *PriceRuleValues {
	p.DynamicFactors = factors
	return p
}

func (p *PriceRuleValues) SetMaxUsagePerUser(maxUsage int) *PriceRuleValues {
	p.MaxUsagePerUser = &maxUsage
	return p
}

func (p *PriceRuleValues) SetMaxUsagePerDay(maxUsage int) *PriceRuleValues {
	p.MaxUsagePerDay = &maxUsage
	return p
}

func (p *PriceRuleValues) SetMaxUsageTotal(maxUsage int) *PriceRuleValues {
	p.MaxUsageTotal = &maxUsage
	return p
}

func (p *PriceRuleValues) SetUsageCount(count int) *PriceRuleValues {
	p.UsageCount = &count
	return p
}

func (p *PriceRuleValues) SetMinDistance(minDistance float64) *PriceRuleValues {
	p.MinDistance = &minDistance
	return p
}

func (p *PriceRuleValues) SetMaxDistance(maxDistance float64) *PriceRuleValues {
	p.MaxDistance = &maxDistance
	return p
}

func (p *PriceRuleValues) SetMinDuration(minDuration int) *PriceRuleValues {
	p.MinDuration = &minDuration
	return p
}

func (p *PriceRuleValues) SetMaxDuration(maxDuration int) *PriceRuleValues {
	p.MaxDuration = &maxDuration
	return p
}

func (p *PriceRuleValues) SetStackableRules(stackableRules []string) *PriceRuleValues {
	p.StackableRules = stackableRules
	return p
}

func (p *PriceRuleValues) SetExclusiveRules(exclusiveRules []string) *PriceRuleValues {
	p.ExclusiveRules = exclusiveRules
	return p
}

func (p *PriceRuleValues) SetCaseSensitive(caseSensitive bool) *PriceRuleValues {
	value := 0
	if caseSensitive {
		value = 1
	}
	p.CaseSensitive = &value
	return p
}

func (p *PriceRuleValues) SetApprovedBy(adminID string) *PriceRuleValues {
	p.ApprovedBy = &adminID
	return p
}

func (p *PriceRuleValues) SetApprovedAt(timestamp int64) *PriceRuleValues {
	p.ApprovedAt = &timestamp
	return p
}

func (p *PriceRuleValues) SetApprovalNotes(notes string) *PriceRuleValues {
	p.ApprovalNotes = &notes
	return p
}

func (p *PriceRuleValues) SetViewCount(count int) *PriceRuleValues {
	p.ViewCount = &count
	return p
}

func (p *PriceRuleValues) SetClickCount(count int) *PriceRuleValues {
	p.ClickCount = &count
	return p
}

func (p *PriceRuleValues) SetUsageToday(count int) *PriceRuleValues {
	p.UsageToday = &count
	return p
}

func (p *PriceRuleValues) SetRevenueImpact(impact float64) *PriceRuleValues {
	p.RevenueImpact = &impact
	return p
}

func (p *PriceRuleValues) SetCostSaved(saved float64) *PriceRuleValues {
	p.CostSaved = &saved
	return p
}

func (p *PriceRuleValues) SetMetadata(metadata map[string]any) *PriceRuleValues {
	p.Metadata = metadata
	return p
}

func (p *PriceRuleValues) SetTags(tags []string) *PriceRuleValues {
	p.Tags = tags
	return p
}

func (p *PriceRuleValues) SetNotes(notes string) *PriceRuleValues {
	p.Notes = &notes
	return p
}

func (p *PriceRuleValues) SetCreatedBy(adminID string) *PriceRuleValues {
	p.CreatedBy = &adminID
	return p
}

// 业务方法
func (p *PriceRule) IsActive() bool {
	return p.GetStatus() == protocol.StatusActive
}

func (p *PriceRule) IsDraft() bool {
	return p.GetStatus() == protocol.StatusDraft
}

func (p *PriceRule) IsPaused() bool {
	return p.GetStatus() == protocol.StatusPaused
}

func (p *PriceRule) IsExpired() bool {
	if p.GetStatus() == protocol.StatusExpired {
		return true
	}

	// 检查时间是否过期
	if p.EndedAt != nil {
		return utils.TimeNowMilli() > *p.EndedAt
	}

	return false
}

func (p *PriceRule) IsUsageLimitReached() bool {
	maxUsage := p.GetMaxUsageTotal()
	currentUsage := p.GetUsageCount()
	return currentUsage >= maxUsage
}

func (p *PriceRule) CanApply() bool {
	return p.IsActive() && !p.IsExpired() && !p.IsUsageLimitReached()
}

func (p *PriceRule) IsValidForTime(timestamp int64) bool {
	if p.StartedAt != nil && timestamp < *p.StartedAt {
		return false
	}
	if p.EndedAt != nil && timestamp > *p.EndedAt {
		return false
	}
	return true
}

func (p *PriceRule) IsDiscountRule() bool {
	category := p.GetCategory()
	return category == protocol.PriceRuleCategoryDiscount || category == protocol.PriceRuleCategoryPromotion
}

func (p *PriceRule) IsSurgeRule() bool {
	return p.GetCategory() == protocol.PriceRuleCategorySurgePricing
}

func (p *PriceRule) IsBasePricingRule() bool {
	return p.GetCategory() == protocol.PriceRuleCategoryBasePricing
}

// 应用范围检查
func (p *PriceRuleValues) IsApplicableToVehicle(category, level string) bool {
	filters := p.GetVehicleFilters()

	// 如果没有设置筛选条件，则适用所有车辆
	if len(filters) == 0 {
		return true
	}

	// 检查是否匹配任一筛选条件
	for _, filter := range filters {
		if (filter.Category == "*" || filter.Category == "" || filter.Category == category) &&
			(filter.Level == "*" || filter.Level == "" || filter.Level == level) {
			return true
		}
	}

	return false
}

func (p *PriceRuleValues) IsApplicableToServiceArea(areaID string) bool {
	if len(p.ServiceAreas) == 0 {
		return true // 如果没有限制，则适用所有区域
	}

	for _, area := range p.ServiceAreas {
		if area == areaID {
			return true
		}
	}

	return false
}

func (p *PriceRuleValues) IsApplicableToUserCategory(category string) bool {
	if len(p.UserCategories) == 0 {
		return true // 如果没有限制，则适用所有用户
	}

	for _, cat := range p.UserCategories {
		if cat == category {
			return true
		}
	}

	return false
}

// 价格计算方法
func (p *PriceRuleValues) CalculateDiscount(originalAmount float64) float64 {
	ruleType := p.GetRuleType()

	switch ruleType {
	case protocol.PriceRuleTypePercentage:
		discount := originalAmount * (p.GetDiscountPercent() / 100.0)
		if p.MaxDiscount != nil && discount > *p.MaxDiscount {
			discount = *p.MaxDiscount
		}
		return discount

	case protocol.PriceRuleTypeFixedAmount:
		discount := p.GetDiscountAmount()
		if discount > originalAmount {
			discount = originalAmount // 不能超过原价
		}
		return discount

	default:
		return 0.0
	}
}

func (p *PriceRuleValues) ApplySurgeMultiplier(baseAmount float64) float64 {
	if p.SurgeMultiplier == nil {
		return baseAmount
	}

	result := baseAmount * (*p.SurgeMultiplier)

	// 检查最大限额
	if p.MaximumFare != nil && result > *p.MaximumFare {
		result = *p.MaximumFare
	}

	return result
}

func (p *PriceRuleValues) CalculateFare(distance float64, duration int) float64 {
	baseRate := p.GetBaseRate()
	perKmRate := 0.0
	perMinuteRate := 0.0

	if p.PerKmRate != nil {
		perKmRate = *p.PerKmRate
	}
	if p.PerMinuteRate != nil {
		perMinuteRate = *p.PerMinuteRate
	}

	fare := baseRate + (distance * perKmRate) + (float64(duration) * perMinuteRate)

	// 应用最低和最高限额
	if p.MinimumFare != nil && fare < *p.MinimumFare {
		fare = *p.MinimumFare
	}
	if p.MaximumFare != nil && fare > *p.MaximumFare {
		fare = *p.MaximumFare
	}

	return fare
}

// 状态管理方法
func (p *PriceRuleValues) Activate() *PriceRuleValues {
	p.SetStatus(protocol.StatusActive)
	return p
}

func (p *PriceRuleValues) Pause() *PriceRuleValues {
	p.SetStatus(protocol.StatusPaused)
	return p
}

func (p *PriceRuleValues) Expire() *PriceRuleValues {
	p.SetStatus(protocol.StatusExpired)
	return p
}

func (p *PriceRuleValues) Approve(adminID string) *PriceRuleValues {
	p.ApprovedBy = &adminID
	now := utils.TimeNowMilli()
	p.ApprovedAt = &now
	p.Activate() // 审批后自动激活
	return p
}

// 使用统计更新
func (p *PriceRuleValues) IncrementUsage() *PriceRuleValues {
	count := p.GetUsageCount() + 1
	p.UsageCount = &count

	// 增加今日使用次数
	todayCount := 0
	if p.UsageToday != nil {
		todayCount = *p.UsageToday
	}
	todayCount++
	p.UsageToday = &todayCount

	return p
}

func (p *PriceRuleValues) IncrementView() *PriceRuleValues {
	count := 0
	if p.ViewCount != nil {
		count = *p.ViewCount
	}
	count++
	p.ViewCount = &count
	return p
}

func (p *PriceRuleValues) IncrementClick() *PriceRuleValues {
	count := 0
	if p.ClickCount != nil {
		count = *p.ClickCount
	}
	count++
	p.ClickCount = &count
	return p
}

func (p *PriceRuleValues) UpdateRevenueImpact(amount float64) *PriceRuleValues {
	current := 0.0
	if p.RevenueImpact != nil {
		current = *p.RevenueImpact
	}
	current += amount
	p.RevenueImpact = &current
	return p
}

func (p *PriceRuleValues) UpdateCostSaved(amount float64) *PriceRuleValues {
	current := 0.0
	if p.CostSaved != nil {
		current = *p.CostSaved
	}
	current += amount
	p.CostSaved = &current
	return p
}

// 检查促销码匹配
func (p *PriceRuleValues) MatchesPromoCode(inputCode string) bool {
	if !p.GetRequiresCode() || p.PromoCode == nil {
		return false
	}

	ruleCode := *p.PromoCode

	// 检查是否区分大小写
	if p.CaseSensitive != nil && *p.CaseSensitive == 1 {
		return ruleCode == inputCode
	}

	return strings.EqualFold(ruleCode, inputCode)
}

// 时间检查方法
func (p *PriceRuleValues) IsValidForDayOfWeek(dayOfWeek int) bool {
	if p.DayOfWeek == nil {
		return true
	}

	allowedDays := *p.DayOfWeek
	dayStr := fmt.Sprintf("%d", dayOfWeek)

	return strings.Contains(allowedDays, dayStr)
}

func (p *PriceRuleValues) IsValidForTimeSlot(hour, minute int) bool {
	if len(p.TimeSlots) == 0 {
		return true
	}

	currentTime := hour*60 + minute

	for _, slot := range p.TimeSlots {
		startTime := slot.StartHour*60 + slot.StartMinute
		endTime := slot.EndHour*60 + slot.EndMinute

		if currentTime >= startTime && currentTime <= endTime {
			return true
		}
	}

	return false
}

func GetPriceRuleByID(ruleID string) *PriceRule {
	var rule PriceRule
	err := GetDB().Where("rule_id = ?", ruleID).First(&rule).Error
	if err != nil {
		return nil
	}
	return &rule
}

func GetActivePriceRules() []*PriceRule {
	var rules []*PriceRule
	err := GetDB().Where("status = ?", protocol.StatusActive).Order("priority DESC").Find(&rules).Error
	if err != nil {
		return nil
	}
	return rules
}

// 便捷创建方法
func NewPercentageDiscountRule(name string, percent float64, promoCode string) *PriceRule {
	rule := NewPriceRule()
	rule.SetRuleName(name).
		SetDisplayName(name).
		SetCategory(protocol.PriceRuleCategoryDiscount).
		SetDiscountPercent(percent)

	if promoCode != "" {
		rule.SetPromoCode(promoCode, false)
	} else {
		rule.SetAutoApply(true)
	}

	return rule
}

func NewFixedDiscountRule(name string, amount float64, promoCode string) *PriceRule {
	rule := NewPriceRule()
	rule.SetRuleName(name).
		SetDisplayName(name).
		SetCategory(protocol.PriceRuleCategoryDiscount).
		SetDiscountAmount(amount)

	if promoCode != "" {
		rule.SetPromoCode(promoCode, false)
	} else {
		rule.SetAutoApply(true)
	}

	return rule
}

func NewSurgeRule(name string, multiplier float64) *PriceRule {
	rule := NewPriceRule()
	rule.SetRuleName(name).
		SetDisplayName(name).
		SetSurgeMultiplier(multiplier).
		SetGlobal(true).
		SetAutoApply(true)

	return rule
}

func NewBasePricingRule(name string, baseRate, perKmRate, perMinuteRate float64) *PriceRule {
	rule := NewPriceRule()
	rule.SetRuleName(name).
		SetDisplayName(name).
		SetCategory(protocol.PriceRuleCategoryBasePricing).
		SetBasePricing(baseRate, perKmRate, perMinuteRate).
		SetGlobal(true).
		SetAutoApply(true)

	return rule
}
