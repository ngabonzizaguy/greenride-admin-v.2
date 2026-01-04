package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PriceSnapshot 价格快照表 - 简化版，专注价格信息
type PriceSnapshot struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SnapshotID string `json:"snapshot_id" gorm:"column:snapshot_id;type:varchar(64);uniqueIndex"`
	*PriceSnapshotValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PriceSnapshotValues struct {
	// 用户标识
	UserID *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`

	// 请求上下文信息
	Distance        *float64 `json:"distance" gorm:"column:distance;type:decimal(10,2)"`               // 距离（公里）
	Duration        *int     `json:"duration" gorm:"column:duration;type:int"`                         // 预计用时（分钟）
	OrderType       *string  `json:"order_type" gorm:"column:order_type;type:varchar(50)"`             // 订单类型
	VehicleCategory *string  `json:"vehicle_category" gorm:"column:vehicle_category;type:varchar(50)"` // 车辆类别
	VehicleLevel    *string  `json:"vehicle_level" gorm:"column:vehicle_level;type:varchar(50)"`       // 车辆级别

	// 价格详情 - 核心数据
	BaseFare       *decimal.Decimal `json:"base_fare" gorm:"column:base_fare;type:decimal(20,6)"`
	DistanceFare   *decimal.Decimal `json:"distance_fare" gorm:"column:distance_fare;type:decimal(20,6)"`
	TimeFare       *decimal.Decimal `json:"time_fare" gorm:"column:time_fare;type:decimal(20,6)"`
	SurgeFare      *decimal.Decimal `json:"surge_fare" gorm:"column:surge_fare;type:decimal(20,6)"`
	ServiceFee     *decimal.Decimal `json:"service_fee" gorm:"column:service_fee;type:decimal(20,6)"`
	OriginalFare   *decimal.Decimal `json:"original_fare" gorm:"column:total_fare;type:decimal(20,6)"`       // 优惠前原始费用
	DiscountedFare *decimal.Decimal `json:"discounted_fare" gorm:"column:estimated_fare;type:decimal(20,6)"` // 优惠后折扣费用
	Currency       *string          `json:"currency" gorm:"column:currency;type:varchar(3);default:'RWF'"`

	// 优惠相关
	DiscountAmount    *decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:decimal(20,6)"`
	PromoCodes        []string         `json:"promo_codes" gorm:"column:promo_codes;type:json;serializer:json"`
	PromoDiscount     *decimal.Decimal `json:"promo_discount" gorm:"column:promo_discount;type:decimal(20,6)"`
	UserPromoDiscount *decimal.Decimal `json:"user_promo_discount" gorm:"column:user_promo_discount;type:decimal(20,6)"`      // 用户优惠券折扣金额
	UserPromotionIDs  []string         `json:"user_promotion_ids" gorm:"column:user_promotion_ids;type:json;serializer:json"` // 使用的用户优惠券ID列表

	// 价格明细分解 - 使用结构体数组
	Breakdowns []*protocol.PriceRuleResult `json:"breakdowns" gorm:"column:breakdowns;type:json;serializer:json"`

	// 计算元信息
	CalculationTime *int64  `json:"calculation_time_ms" gorm:"column:calculation_time_ms;type:bigint"` // 计算耗时（毫秒）
	RulesEvaluated  *int    `json:"rules_evaluated" gorm:"column:rules_evaluated;type:int"`            // 评估的规则数量
	RulesApplied    *int    `json:"rules_applied" gorm:"column:rules_applied;type:int"`                // 应用的规则数量
	EngineVersion   *string `json:"engine_version" gorm:"column:engine_version;type:varchar(20)"`      // 引擎版本

	// 状态管理
	Status      *string `json:"status" gorm:"column:status;type:varchar(32);default:'active'"` // active, expired
	ScheduledAt *int64  `json:"scheduled_at" gorm:"column:scheduled_at;type:bigint;index"`     // 预定时间
	ExpiresAt   *int64  `json:"expires_at" gorm:"column:expires_at;type:bigint;index"`         // 过期时间

	// 业务关联
	OrderID *string `json:"order_id" gorm:"column:order_id;type:varchar(64);index"` // 关联订单ID

	// 扩展信息 - 使用 MapData
	Metadata protocol.MapData `json:"metadata" gorm:"column:metadata;type:json;serializer:json"`

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (PriceSnapshot) TableName() string {
	return "t_price_snapshots"
}

// 创建新的价格快照对象
func NewPriceSnapshot(userID string) *PriceSnapshot {
	expiresAt := utils.TimeNowMilli() + 30*60*1000 // 30分钟有效期
	zero := decimal.Zero

	return &PriceSnapshot{
		SnapshotID: utils.GenerateSnapshotID(),
		PriceSnapshotValues: &PriceSnapshotValues{
			UserID:            &userID,
			Status:            utils.StringPtr(protocol.StatusActive),
			ExpiresAt:         &expiresAt,
			Currency:          utils.StringPtr(protocol.CurrencyRWF),
			VehicleCategory:   utils.StringPtr(protocol.VehicleCategorySedan), // 默认小车
			VehicleLevel:      utils.StringPtr(protocol.VehicleLevelEconomy),  // 默认经济型
			PromoCodes:        []string{},                                     // 初始化为空数组，避免 JSON 错误
			UserPromotionIDs:  []string{},                                     // 初始化为空数组，避免 JSON 错误
			Breakdowns:        []*protocol.PriceRuleResult{},                  // 初始化为空数组
			Metadata:          make(protocol.MapData),
			DiscountAmount:    &decimal.Zero, // 初始化为0
			PromoDiscount:     &decimal.Zero, // 初始化为0
			UserPromoDiscount: &decimal.Zero, // 初始化为0
			BaseFare:          &zero,         // 初始化为0
			DistanceFare:      &zero,         // 初始化为0
			TimeFare:          &zero,         // 初始化为0
			SurgeFare:         &zero,         // 初始化为0
			ServiceFee:        &zero,         // 初始化为0
			OriginalFare:      &zero,         // 初始化为0
			DiscountedFare:    &zero,         // 初始化为0
		},
	}
}

func (v *PriceSnapshotValues) GetScheduledAt() int64 {
	if v.ScheduledAt == nil {
		return 0
	}
	return *v.ScheduledAt
}

func (v *PriceSnapshotValues) SetScheduledAt(scheduled_at int64) *PriceSnapshotValues {
	v.ScheduledAt = &scheduled_at
	return v
}

// Getter 方法
func (p *PriceSnapshotValues) GetUserID() string {
	return utils.SafeStringDeref(p.UserID)
}

func (p *PriceSnapshotValues) GetDistance() float64 {
	return utils.SafeFloat64Deref(p.Distance)
}

func (p *PriceSnapshotValues) GetDuration() int {
	if p.Duration == nil {
		return 0
	}
	return *p.Duration
}

func (p *PriceSnapshotValues) GetOrderType() string {
	return utils.SafeStringDeref(p.OrderType)
}
func (p *PriceSnapshotValues) GetBaseFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.BaseFare)
}

func (p *PriceSnapshotValues) GetDistanceFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.DistanceFare)
}

func (p *PriceSnapshotValues) GetTimeFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.TimeFare)
}

func (p *PriceSnapshotValues) GetSurgeFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.SurgeFare)
}

func (p *PriceSnapshotValues) GetServiceFee() decimal.Decimal {
	return utils.SafeDecimalDeref(p.ServiceFee)
}

func (p *PriceSnapshotValues) GetDiscountedFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.DiscountedFare)
}

func (p *PriceSnapshotValues) GetOriginalFare() decimal.Decimal {
	return utils.SafeDecimalDeref(p.OriginalFare)
}

func (p *PriceSnapshotValues) GetCurrency() string {
	if p.Currency == nil {
		return protocol.CurrencyRWF
	}
	return *p.Currency
}

func (p *PriceSnapshotValues) GetDiscountAmount() decimal.Decimal {
	if p.DiscountAmount == nil {
		return decimal.Zero
	}
	return *p.DiscountAmount
}

func (p *PriceSnapshotValues) GetPromoCodes() []string {
	if p.PromoCodes == nil {
		return []string{}
	}
	return p.PromoCodes
}

func (p *PriceSnapshotValues) GetPromoDiscount() decimal.Decimal {
	if p.PromoDiscount == nil {
		return decimal.Zero
	}
	return *p.PromoDiscount
}

func (p *PriceSnapshotValues) GetUserPromoDiscount() decimal.Decimal {
	if p.UserPromoDiscount == nil {
		return decimal.Zero
	}
	return *p.UserPromoDiscount
}

func (p *PriceSnapshotValues) GetUserPromotionIDs() []string {
	if p.UserPromotionIDs == nil {
		return []string{}
	}
	return p.UserPromotionIDs
}

func (p *PriceSnapshotValues) GetBreakdowns() []*protocol.PriceRuleResult {
	if p.Breakdowns == nil {
		return []*protocol.PriceRuleResult{}
	}
	return p.Breakdowns
}

func (p *PriceSnapshotValues) GetStatus() string {
	if p.Status == nil {
		return protocol.StatusActive
	}
	return *p.Status
}

func (p *PriceSnapshotValues) GetExpiresAt() int64 {
	if p.ExpiresAt == nil {
		return 0
	}
	return *p.ExpiresAt
}

func (p *PriceSnapshotValues) GetOrderID() string {
	return utils.SafeStringDeref(p.OrderID)
}

func (p *PriceSnapshotValues) GetMetadata() protocol.MapData {
	if p.Metadata == nil {
		return protocol.MapData{}
	}
	return p.Metadata
}

func (p *PriceSnapshotValues) GetUpdatedAt() int64 {
	return p.UpdatedAt
}

func (p *PriceSnapshotValues) GetCalculationTime() int64 {
	if p.CalculationTime == nil {
		return 0
	}
	return *p.CalculationTime
}

func (p *PriceSnapshotValues) GetRulesEvaluated() int {
	if p.RulesEvaluated == nil {
		return 0
	}
	return *p.RulesEvaluated
}

func (p *PriceSnapshotValues) GetRulesApplied() int {
	if p.RulesApplied == nil {
		return 0
	}
	return *p.RulesApplied
}

func (p *PriceSnapshotValues) GetEngineVersion() string {
	if p.EngineVersion == nil {
		return "v1.0.0"
	}
	return *p.EngineVersion
}
func (p *PriceSnapshotValues) GetVehicleCategory() string {
	if p.VehicleCategory == nil {
		return ""
	}
	return *p.VehicleCategory
}
func (p *PriceSnapshotValues) GetVehicleLevel() string {
	if p.VehicleLevel == nil {
		return ""
	}
	return *p.VehicleLevel
}
func (p *PriceSnapshotValues) SetVehicleCategory(category string) *PriceSnapshotValues {
	p.VehicleCategory = &category
	return p
}
func (p *PriceSnapshotValues) SetVehicleLevel(level string) *PriceSnapshotValues {
	p.VehicleLevel = &level
	return p
}

// Setter 方法
func (p *PriceSnapshotValues) SetUserID(userID string) *PriceSnapshotValues {
	p.UserID = &userID
	return p
}

func (p *PriceSnapshotValues) SetDistance(distance float64) *PriceSnapshotValues {
	p.Distance = &distance
	return p
}

func (p *PriceSnapshotValues) SetDuration(duration int) *PriceSnapshotValues {
	p.Duration = &duration
	return p
}

func (p *PriceSnapshotValues) SetOrderType(orderType string) *PriceSnapshotValues {
	p.OrderType = &orderType
	return p
}

func (p *PriceSnapshotValues) SetBaseFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.BaseFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetDistanceFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.DistanceFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetTimeFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.TimeFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetSurgeFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.SurgeFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetServiceFee(fee decimal.Decimal) *PriceSnapshotValues {
	p.ServiceFee = &fee
	return p
}

func (p *PriceSnapshotValues) SetDiscountedFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.DiscountedFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetOriginalFare(fare decimal.Decimal) *PriceSnapshotValues {
	p.OriginalFare = &fare
	return p
}

func (p *PriceSnapshotValues) SetCurrency(currency string) *PriceSnapshotValues {
	p.Currency = &currency
	return p
}

func (p *PriceSnapshotValues) SetDiscountAmount(amount decimal.Decimal) *PriceSnapshotValues {
	p.DiscountAmount = &amount
	return p
}

func (p *PriceSnapshotValues) SetPromoCodes(codes []string) *PriceSnapshotValues {
	if codes == nil {
		p.PromoCodes = []string{} // 设置为空数组而不是 nil
	} else {
		p.PromoCodes = codes
	}
	return p
}

func (p *PriceSnapshotValues) SetPromoDiscount(discount decimal.Decimal) *PriceSnapshotValues {
	p.PromoDiscount = &discount
	return p
}

func (p *PriceSnapshotValues) SetUserPromoDiscount(discount decimal.Decimal) *PriceSnapshotValues {
	p.UserPromoDiscount = &discount
	return p
}

func (p *PriceSnapshotValues) SetUserPromotionIDs(promotionIDs []string) *PriceSnapshotValues {
	if promotionIDs == nil {
		p.UserPromotionIDs = []string{} // 设置为空数组而不是 nil
	} else {
		p.UserPromotionIDs = promotionIDs
	}
	return p
}

func (p *PriceSnapshotValues) SetBreakdowns(rules []*protocol.PriceRuleResult) *PriceSnapshotValues {
	if rules == nil {
		p.Breakdowns = []*protocol.PriceRuleResult{} // 设置为空数组而不是 nil
	} else {
		p.Breakdowns = rules
	}
	return p
}

func (p *PriceSnapshotValues) SetStatus(status string) *PriceSnapshotValues {
	p.Status = &status
	return p
}

func (p *PriceSnapshotValues) SetExpiresAt(timestamp int64) *PriceSnapshotValues {
	p.ExpiresAt = &timestamp
	return p
}

func (p *PriceSnapshotValues) SetOrderID(orderID string) *PriceSnapshotValues {
	p.OrderID = &orderID
	return p
}

func (p *PriceSnapshotValues) SetMetadata(metadata protocol.MapData) *PriceSnapshotValues {
	p.Metadata = metadata
	return p
}

func (p *PriceSnapshotValues) SetCalculationTime(time int64) *PriceSnapshotValues {
	p.CalculationTime = &time
	return p
}

func (p *PriceSnapshotValues) SetRulesEvaluated(count int) *PriceSnapshotValues {
	p.RulesEvaluated = &count
	return p
}

func (p *PriceSnapshotValues) SetRulesApplied(count int) *PriceSnapshotValues {
	p.RulesApplied = &count
	return p
}

func (p *PriceSnapshotValues) SetEngineVersion(version string) *PriceSnapshotValues {
	p.EngineVersion = &version
	return p
}

// 业务方法
func (p *PriceSnapshotValues) IsValid() bool {
	return p.GetStatus() == protocol.StatusActive && !p.IsExpired()
}

func (p *PriceSnapshotValues) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return utils.TimeNowMilli() > *p.ExpiresAt
}

func (p *PriceSnapshotValues) MarkAsExpired() *PriceSnapshotValues {
	p.Status = utils.StringPtr(protocol.StatusExpired)
	return p
}

// GetPriceSnapshotByID 根据快照ID获取价格快照
func GetPriceSnapshotByID(snapshotID string) *PriceSnapshot {
	if snapshotID == "" {
		return nil
	}

	var snapshot PriceSnapshot
	if err := GetDB().Where("snapshot_id = ?", snapshotID).First(&snapshot).Error; err != nil {
		return nil
	}

	return &snapshot
}

func UpdatePriceOrderID(tx *gorm.DB, snapshotID, orderID string) error {
	if snapshotID == "" || orderID == "" {
		return nil
	}
	values := &PriceSnapshotValues{}
	values.SetOrderID(orderID)
	return tx.Model(&PriceSnapshot{}).Where("snapshot_id = ?", snapshotID).UpdateColumns(values).Error
}

func (t *PriceSnapshot) Protocol() *protocol.OrderPrice {
	if t == nil {
		return nil
	}

	toFloat64 := func(d decimal.Decimal) float64 {
		val, _ := d.Float64()
		return val
	}

	return &protocol.OrderPrice{
		// 请求上下文信息
		Currency:        t.GetCurrency(),
		Distance:        t.GetDistance(),
		Duration:        t.GetDuration(),
		OrderType:       t.GetOrderType(),
		VehicleCategory: t.GetVehicleCategory(),
		VehicleLevel:    t.GetVehicleLevel(),
		// 价格分解
		PriceBreakdown: &protocol.PriceBreakdown{
			BaseFare:          toFloat64(t.GetBaseFare()),
			DistanceFare:      toFloat64(t.GetDistanceFare()),
			TimeFare:          toFloat64(t.GetTimeFare()),
			SurgeFare:         toFloat64(t.GetSurgeFare()),
			ServiceFee:        toFloat64(t.GetServiceFee()),
			DiscountedFare:    toFloat64(t.GetDiscountedFare()),
			DiscountAmount:    func() float64 { val, _ := t.GetDiscountAmount().Float64(); return val }(),
			PromoDiscount:     func() float64 { val, _ := t.GetPromoDiscount().Float64(); return val }(),
			UserPromoDiscount: func() float64 { val, _ := t.GetUserPromoDiscount().Float64(); return val }(),
			OriginalFare:      toFloat64(t.GetOriginalFare()),
		},
		// 价格锁定信息
		PriceID:   t.SnapshotID,
		ExpiresAt: t.GetExpiresAt(),
		IsLocked:  false, // 默认未锁定

		// 计算元信息
		CalculationTime: t.GetCalculationTime(),
		RulesEvaluated:  t.GetRulesEvaluated(),
		RulesApplied:    t.GetRulesApplied(),
		EngineVersion:   t.GetEngineVersion(),

		// 价格明细分解
		Breakdowns: t.GetBreakdowns(),
	}
}

// =============================================================================
// 向后兼容方法 - 逐步迁移时保持兼容性
// =============================================================================

// GetTotalFare 获取总费用 - 向后兼容方法，实际调用GetOriginalFare
// @deprecated 请使用GetOriginalFare()替代
func (p *PriceSnapshotValues) GetTotalFare() float64 {
	val, _ := p.GetOriginalFare().Float64()
	return val
}

// SetTotalFare 设置总费用 - 向后兼容方法，实际调用SetOriginalFare
// @deprecated 请使用SetOriginalFare()替代
func (p *PriceSnapshotValues) SetTotalFare(fare float64) *PriceSnapshotValues {
	return p.SetOriginalFare(decimal.NewFromFloat(fare))
}

// GetEstimatedFare 获取预估费用 - 向后兼容方法，实际调用GetDiscountedFare
// @deprecated 请使用GetDiscountedFare()替代
func (p *PriceSnapshotValues) GetEstimatedFare() float64 {
	val, _ := p.GetDiscountedFare().Float64()
	return val
}

// SetEstimatedFare 设置预估费用 - 向后兼容方法，实际调用SetDiscountedFare
// @deprecated 请使用SetDiscountedFare()替代
func (p *PriceSnapshotValues) SetEstimatedFare(fare float64) *PriceSnapshotValues {
	return p.SetDiscountedFare(decimal.NewFromFloat(fare))
}
