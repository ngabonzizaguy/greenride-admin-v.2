package models

import (
	"log"
	"maps"

	"greenride/internal/protocol"
	"greenride/internal/utils"

	"gorm.io/gorm"
)

// UserPromotion user promotion association table
type UserPromotion struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID      string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`           // User ID
	PromotionID string `json:"promotion_id" gorm:"column:promotion_id;type:varchar(64);index"` // Promotion ID
	*UserPromotionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserPromotionValues struct {
	// Basic information (redundant storage for easy querying)
	Code        *string `json:"code" gorm:"column:code;type:varchar(50);index"`  // Promo code
	Title       *string `json:"title" gorm:"column:title;type:varchar(255)"`     // Promotion title
	Description *string `json:"description" gorm:"column:description;type:text"` // Promotion description

	// Discount information (snapshot to avoid affecting issued coupons when original promotion is modified)
	DiscountType      *string  `json:"discount_type" gorm:"column:discount_type;type:varchar(20)"`               // percentage, fixed_amount
	DiscountValue     *float64 `json:"discount_value" gorm:"column:discount_value;type:decimal(10,2)"`           // Discount amount or percentage
	MaxDiscountAmount *float64 `json:"max_discount_amount" gorm:"column:max_discount_amount;type:decimal(10,2)"` // Maximum discount amount
	MinOrderAmount    *float64 `json:"min_order_amount" gorm:"column:min_order_amount;type:decimal(10,2)"`       // Minimum order amount

	// Usage status
	Status     *string  `json:"status" gorm:"column:status;type:varchar(30);index;default:'available'"` //
	IsUsed     *int     `json:"is_used" gorm:"column:is_used;default:0"`                                // Whether used (0:unused, 1:used)
	UsedAt     *int64   `json:"used_at" gorm:"column:used_at"`                                          // Usage time
	UsedAmount *float64 `json:"used_amount" gorm:"column:used_amount;type:decimal(10,2);default:0.00"`  // Used discount amount
	OrderID    *string  `json:"order_id" gorm:"column:order_id;type:varchar(64)"`                       // Order ID where used

	// Validity period
	ExpiredAt *int64 `json:"expired_at" gorm:"column:expired_at"` // Expiration time (millisecond timestamp)

	// 来源信息
	Source     *string `json:"source" gorm:"column:source;type:varchar(50)"`            // 来源：system, admin, event, referral
	SourceID   *string `json:"source_id" gorm:"column:source_id;type:varchar(100)"`     // 来源ID
	SourceDesc *string `json:"source_desc" gorm:"column:source_desc;type:varchar(255)"` // 来源描述

	// 批次管理
	BatchID *string `json:"batch_id" gorm:"column:batch_id;type:varchar(64);index"` // 发放批次ID

	// 渠道分析
	Channel  *string `json:"channel" gorm:"column:channel;type:varchar(50)"`     // 注册渠道（app, web, h5等）
	CityCode *string `json:"city_code" gorm:"column:city_code;type:varchar(20)"` // 用户注册城市代码

	// A/B测试
	ExperimentID *string `json:"experiment_id" gorm:"column:experiment_id;type:varchar(64)"` // 实验ID
	VariantID    *string `json:"variant_id" gorm:"column:variant_id;type:varchar(64)"`       // 变体ID

	// 用户行为追踪
	ActivatedAt *int64 `json:"activated_at" gorm:"column:activated_at"`       // 激活时间
	ViewCount   *int   `json:"view_count" gorm:"column:view_count;default:0"` // 查看次数
	SharedAt    *int64 `json:"shared_at" gorm:"column:shared_at"`             // 分享时间

	// 发放信息
	IssuedBy *string `json:"issued_by" gorm:"column:issued_by;type:varchar(100)"` // 发放者

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (UserPromotion) TableName() string {
	return "t_user_promotions"
}

// 创建用户优惠券
func NewUserPromotion(userID string, promotion *Promotion) *UserPromotion {
	code := promotion.GetCode()
	title := promotion.GetTitle()
	description := promotion.GetDescription()
	discountType := promotion.GetDiscountType()
	discountValue := promotion.GetDiscountValue()
	maxDiscountAmount := promotion.GetMaxDiscountAmount()
	minOrderAmount := promotion.GetMinOrderAmount()
	status := protocol.StatusAvailable
	isUsed := 0
	source := protocol.UserPromotionSourceSystem
	issuedBy := "system"

	userPromo := &UserPromotion{
		UserID:      userID,
		PromotionID: promotion.PromotionID,
		UserPromotionValues: &UserPromotionValues{
			Code:              &code,
			Title:             &title,
			Description:       &description,
			DiscountType:      &discountType,
			DiscountValue:     &discountValue,
			MaxDiscountAmount: &maxDiscountAmount,
			MinOrderAmount:    &minOrderAmount,
			Status:            &status,
			IsUsed:            &isUsed,
			Source:            &source,
			IssuedBy:          &issuedBy,
		},
	}

	// 简化版本：用户优惠券永不过期（或可以设置默认过期时间）
	// 如果需要过期时间，可以设置默认30天
	// defaultValidDays := int64(30)
	// expiredAt := utils.TimeNowMilli() + defaultValidDays*24*3600*1000
	// userPromo.ExpiredAt = &expiredAt

	return userPromo
}

// 创建用户优惠券（带来源信息）
func NewUserPromotionWithSource(userID string, promotion *Promotion, source, sourceID, sourceDesc, issuedBy string) *UserPromotion {
	userPromo := NewUserPromotion(userID, promotion)
	userPromo.Source = &source
	userPromo.SourceID = &sourceID
	userPromo.SourceDesc = &sourceDesc
	userPromo.IssuedBy = &issuedBy
	return userPromo
}

// 业务方法
func (u *UserPromotion) IsAvailable() bool {
	return u.GetStatus() == protocol.StatusAvailable && !u.IsUsed() && !u.IsExpired()
}

func (u *UserPromotion) IsExpired() bool {
	if u.ExpiredAt == nil {
		return false
	}
	return *u.ExpiredAt < utils.TimeNowMilli()
}

func (u *UserPromotion) IsUsed() bool {
	return u.GetIsUsed() == 1 || u.GetStatus() == protocol.StatusUsed
}

func (u *UserPromotion) CanUse() bool {
	return u.IsAvailable()
}

// 使用优惠券
func (u *UserPromotionValues) Use(orderID string, usedAmount float64) *UserPromotionValues {
	status := protocol.StatusUsed
	u.Status = &status
	isUsed := 1
	u.IsUsed = &isUsed
	u.OrderID = &orderID
	u.UsedAmount = &usedAmount
	now := utils.TimeNowMilli()
	u.UsedAt = &now
	return u
}

// 标记过期
func (u *UserPromotionValues) MarkExpired() *UserPromotionValues {
	status := protocol.StatusExpired
	u.Status = &status
	return u
}

// 设置过期时间
func (u *UserPromotionValues) SetExpiry(expiredAt int64) *UserPromotionValues {
	u.ExpiredAt = &expiredAt
	return u
}

// 延长有效期
func (u *UserPromotionValues) ExtendExpiry(days int) *UserPromotionValues {
	if u.ExpiredAt == nil {
		// 如果没有过期时间，从现在开始计算
		expiredAt := utils.TimeNowMilli() + int64(days)*24*3600*1000
		u.ExpiredAt = &expiredAt
	} else {
		// 在原有过期时间基础上延长
		*u.ExpiredAt += int64(days) * 24 * 3600 * 1000
	}
	return u
}

// 重置为可用状态（仅限管理操作）
func (u *UserPromotionValues) Reset() *UserPromotionValues {
	status := protocol.StatusAvailable
	u.Status = &status
	isUsed := 0
	u.IsUsed = &isUsed
	u.OrderID = nil
	usedAmount := 0.0
	u.UsedAmount = &usedAmount
	u.UsedAt = nil
	return u
}

// Getter 方法
func (u *UserPromotionValues) GetStatus() string {
	if u.Status == nil {
		return ""
	}
	return *u.Status
}

func (u *UserPromotionValues) GetCode() string {
	if u.Code == nil {
		return ""
	}
	return *u.Code
}

func (u *UserPromotionValues) GetTitle() string {
	if u.Title == nil {
		return ""
	}
	return *u.Title
}

func (u *UserPromotionValues) GetDescription() string {
	if u.Description == nil {
		return ""
	}
	return *u.Description
}

func (u *UserPromotionValues) GetDiscountType() string {
	if u.DiscountType == nil {
		return ""
	}
	return *u.DiscountType
}

func (u *UserPromotionValues) GetDiscountValue() float64 {
	if u.DiscountValue == nil {
		return 0.0
	}
	return *u.DiscountValue
}

func (u *UserPromotionValues) GetMaxDiscountAmount() float64 {
	if u.MaxDiscountAmount == nil {
		return 0.0
	}
	return *u.MaxDiscountAmount
}

func (u *UserPromotionValues) GetMinOrderAmount() float64 {
	if u.MinOrderAmount == nil {
		return 0.0
	}
	return *u.MinOrderAmount
}

func (u *UserPromotionValues) GetSource() string {
	if u.Source == nil {
		return ""
	}
	return *u.Source
}

func (u *UserPromotionValues) GetSourceID() string {
	if u.SourceID == nil {
		return ""
	}
	return *u.SourceID
}

func (u *UserPromotionValues) GetSourceDesc() string {
	if u.SourceDesc == nil {
		return ""
	}
	return *u.SourceDesc
}

func (u *UserPromotionValues) GetOrderID() string {
	if u.OrderID == nil {
		return ""
	}
	return *u.OrderID
}

func (u *UserPromotionValues) GetUsedAt() int64 {
	if u.UsedAt == nil {
		return 0
	}
	return *u.UsedAt
}

func (u *UserPromotionValues) GetExpiredAt() int64 {
	if u.ExpiredAt == nil {
		return 0
	}
	return *u.ExpiredAt
}

func (u *UserPromotionValues) GetIsUsed() int {
	if u.IsUsed == nil {
		return 0
	}
	return *u.IsUsed
}

func (u *UserPromotionValues) GetUsedAmount() float64 {
	if u.UsedAmount == nil {
		return 0.0
	}
	return *u.UsedAmount
}

func (u *UserPromotionValues) GetIssuedBy() string {
	if u.IssuedBy == nil {
		return ""
	}
	return *u.IssuedBy
}

func (u *UserPromotionValues) GetUpdatedAt() int64 {
	return u.UpdatedAt
}

func (u *UserPromotionValues) GetBatchID() string {
	if u.BatchID == nil {
		return ""
	}
	return *u.BatchID
}

func (u *UserPromotionValues) GetChannel() string {
	if u.Channel == nil {
		return ""
	}
	return *u.Channel
}

func (u *UserPromotionValues) GetCityCode() string {
	if u.CityCode == nil {
		return ""
	}
	return *u.CityCode
}

func (u *UserPromotionValues) GetExperimentID() string {
	if u.ExperimentID == nil {
		return ""
	}
	return *u.ExperimentID
}

func (u *UserPromotionValues) GetVariantID() string {
	if u.VariantID == nil {
		return ""
	}
	return *u.VariantID
}

func (u *UserPromotionValues) GetActivatedAt() int64 {
	if u.ActivatedAt == nil {
		return 0
	}
	return *u.ActivatedAt
}

func (u *UserPromotionValues) GetViewCount() int {
	if u.ViewCount == nil {
		return 0
	}
	return *u.ViewCount
}

func (u *UserPromotionValues) GetSharedAt() int64 {
	if u.SharedAt == nil {
		return 0
	}
	return *u.SharedAt
}

// Setter 方法
func (u *UserPromotionValues) SetCode(code string) *UserPromotionValues {
	u.Code = &code
	return u
}

func (u *UserPromotionValues) SetTitle(title string) *UserPromotionValues {
	u.Title = &title
	return u
}

func (u *UserPromotionValues) SetDescription(description string) *UserPromotionValues {
	u.Description = &description
	return u
}

func (u *UserPromotionValues) SetDiscountType(discountType string) *UserPromotionValues {
	u.DiscountType = &discountType
	return u
}

func (u *UserPromotionValues) SetDiscountValue(discountValue float64) *UserPromotionValues {
	u.DiscountValue = &discountValue
	return u
}

func (u *UserPromotionValues) SetMaxDiscountAmount(maxDiscountAmount float64) *UserPromotionValues {
	u.MaxDiscountAmount = &maxDiscountAmount
	return u
}

func (u *UserPromotionValues) SetMinOrderAmount(minOrderAmount float64) *UserPromotionValues {
	u.MinOrderAmount = &minOrderAmount
	return u
}

func (u *UserPromotionValues) SetStatus(status string) *UserPromotionValues {
	u.Status = &status
	return u
}

func (u *UserPromotionValues) SetIsUsed(isUsed int) *UserPromotionValues {
	u.IsUsed = &isUsed
	return u
}

func (u *UserPromotionValues) SetUsedAt(usedAt int64) *UserPromotionValues {
	u.UsedAt = &usedAt
	return u
}

func (u *UserPromotionValues) SetUsedAmount(usedAmount float64) *UserPromotionValues {
	u.UsedAmount = &usedAmount
	return u
}

func (u *UserPromotionValues) SetOrderID(orderID string) *UserPromotionValues {
	u.OrderID = &orderID
	return u
}

func (u *UserPromotionValues) SetExpiredAt(expiredAt int64) *UserPromotionValues {
	u.ExpiredAt = &expiredAt
	return u
}

func (u *UserPromotionValues) SetSource(source string) *UserPromotionValues {
	u.Source = &source
	return u
}

func (u *UserPromotionValues) SetSourceID(sourceID string) *UserPromotionValues {
	u.SourceID = &sourceID
	return u
}

func (u *UserPromotionValues) SetSourceDesc(sourceDesc string) *UserPromotionValues {
	u.SourceDesc = &sourceDesc
	return u
}

func (u *UserPromotionValues) SetIssuedBy(issuedBy string) *UserPromotionValues {
	u.IssuedBy = &issuedBy
	return u
}

func (u *UserPromotionValues) SetBatchID(batchID string) *UserPromotionValues {
	u.BatchID = &batchID
	return u
}

func (u *UserPromotionValues) SetChannel(channel string) *UserPromotionValues {
	u.Channel = &channel
	return u
}

func (u *UserPromotionValues) SetCityCode(cityCode string) *UserPromotionValues {
	u.CityCode = &cityCode
	return u
}

func (u *UserPromotionValues) SetExperimentID(experimentID string) *UserPromotionValues {
	u.ExperimentID = &experimentID
	return u
}

func (u *UserPromotionValues) SetVariantID(variantID string) *UserPromotionValues {
	u.VariantID = &variantID
	return u
}

func (u *UserPromotionValues) SetActivatedAt(activatedAt int64) *UserPromotionValues {
	u.ActivatedAt = &activatedAt
	return u
}

func (u *UserPromotionValues) SetViewCount(viewCount int) *UserPromotionValues {
	u.ViewCount = &viewCount
	return u
}

func (u *UserPromotionValues) SetSharedAt(sharedAt int64) *UserPromotionValues {
	u.SharedAt = &sharedAt
	return u
}

// Protocol 转换
func (u *UserPromotion) Protocol() *protocol.UserPromotion {
	return &protocol.UserPromotion{
		ID:          u.ID,
		UserID:      u.UserID,
		PromotionID: u.PromotionID,
		Code:        u.GetCode(),
		Title:       u.GetTitle(),
		Description: u.GetDescription(),

		DiscountType:      u.GetDiscountType(),
		DiscountValue:     u.GetDiscountValue(),
		MaxDiscountAmount: u.GetMaxDiscountAmount(),
		MinOrderAmount:    u.GetMinOrderAmount(),

		Status:     u.GetStatus(),
		IsUsed:     u.IsUsed(),
		UsedAt:     u.GetUsedAt(),
		UsedAmount: u.GetUsedAmount(),
		OrderID:    u.GetOrderID(),

		ExpiredAt:  u.GetExpiredAt(),
		Source:     u.GetSource(),
		SourceID:   u.GetSourceID(),
		SourceDesc: u.GetSourceDesc(),
		IssuedBy:   u.GetIssuedBy(),

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// 便捷创建方法

// CreateWelcomePromoForUser 为新用户创建欢迎优惠券
func CreateWelcomePromoForUser(userID string, promotion *Promotion) *UserPromotion {
	return NewUserPromotionWithSource(
		userID,
		promotion,
		protocol.UserPromotionSourceWelcome,
		"new_user_welcome",
		"New User Welcome Promotion",
		"system",
	)
}

// CreateReferralPromoForUser 为推荐用户创建推荐优惠券
func CreateReferralPromoForUser(userID, referrerID string, promotion *Promotion) *UserPromotion {
	return NewUserPromotionWithSource(
		userID,
		promotion,
		protocol.UserPromotionSourceReferral,
		referrerID,
		"Referral Reward",
		"system",
	)
}

// CreateEventPromoForUser 为用户创建活动优惠券
func CreateEventPromoForUser(userID, eventID, eventName string, promotion *Promotion) *UserPromotion {
	return NewUserPromotionWithSource(
		userID,
		promotion,
		protocol.UserPromotionSourceEvent,
		eventID,
		eventName,
		"system",
	)
}

// CreateAdminPromoForUser 管理员为用户发放优惠券
func CreateAdminPromoForUser(userID, adminID, reason string, promotion *Promotion) *UserPromotion {
	return NewUserPromotionWithSource(
		userID,
		promotion,
		protocol.UserPromotionSourceAdmin,
		adminID,
		reason,
		adminID,
	)
}

// BatchCreatePromoForUsers 批量为用户创建优惠券
func BatchCreatePromoForUsers(userIDs []string, promotion *Promotion, source, sourceID, sourceDesc, issuedBy string) []*UserPromotion {
	userPromos := make([]*UserPromotion, 0, len(userIDs))
	for _, userID := range userIDs {
		userPromo := NewUserPromotionWithSource(userID, promotion, source, sourceID, sourceDesc, issuedBy)
		userPromos = append(userPromos, userPromo)
	}
	return userPromos
}

// ToRuleCompatibleFields 将用户优惠券字段映射为价格规则兼容格式
// 返回值: (ruleType, discountPercent, discountAmount)
func (u *UserPromotion) ToRuleCompatibleFields() (string, *float64, *float64) {
	if u == nil {
		return "", nil, nil
	}

	discountValue := u.GetDiscountValue()
	switch u.GetDiscountType() {
	case protocol.DiscountTypePercentage:
		return protocol.PriceRuleTypePercentage, &discountValue, nil
	case protocol.DiscountTypeFixed:
		return protocol.PriceRuleTypeFixedAmount, nil, &discountValue
	default:
		return "", nil, nil
	}
}

// GetRuleType 获取等价的价格规则类型
func (u *UserPromotion) GetRuleType() string {
	ruleType, _, _ := u.ToRuleCompatibleFields()
	return ruleType
}

// GetDiscountPercent 获取折扣百分比（用于价格规则兼容）
func (u *UserPromotion) GetDiscountPercent() *float64 {
	_, discountPercent, _ := u.ToRuleCompatibleFields()
	return discountPercent
}

// GetDiscountAmount 获取固定折扣金额（用于价格规则兼容）
func (u *UserPromotion) GetDiscountAmount() *float64 {
	_, _, discountAmount := u.ToRuleCompatibleFields()
	return discountAmount
}

// GetMaxDiscount 获取最大折扣限制（兼容价格规则）
func (u *UserPromotion) GetMaxDiscount() *float64 {
	return u.MaxDiscountAmount
}

// GetMinOrderAmountPtr 获取最小订单金额指针（兼容价格规则）
func (u *UserPromotion) GetMinOrderAmountPtr() *float64 {
	return u.MinOrderAmount
}

// ToPriceRule 将用户优惠券完全适配为 PriceRule 实体
// 这样可以复用现有的价格计算逻辑
func (u *UserPromotion) ToPriceRule() *PriceRule {
	if u == nil {
		return nil
	}
	// 获取字段映射
	ruleType, discountPercent, discountAmount := u.ToRuleCompatibleFields()
	status := protocol.StatusActive
	if !u.IsAvailable() {
		status = protocol.StatusInactive
	}
	// 创建 PriceRule 实体
	priceRule := &PriceRule{
		RuleID: u.PromotionID,
		PriceRuleValues: &PriceRuleValues{
			// 基本信息
			RuleName:    utils.StringPtr(u.GetTitle()),
			DisplayName: utils.StringPtr(u.GetTitle()),
			Description: utils.StringPtr(u.GetDescription()),
			Category:    utils.StringPtr(protocol.PriceRuleCategoryUserPromotion), // 用户优惠券归类为促销

			// 规则类型和折扣信息
			RuleType:        utils.StringPtr(ruleType),
			DiscountPercent: discountPercent,
			DiscountAmount:  discountAmount,
			MaxDiscount:     u.MaxDiscountAmount,
			MinOrderAmount:  u.MinOrderAmount,

			// 促销码信息
			PromoCode:     u.Code,
			RequiresCode:  utils.IntPtr(1), // 用户优惠券总是需要代码
			CaseSensitive: utils.IntPtr(0), // 默认不区分大小写

			// 状态信息 - 映射用户券状态到规则状态
			Status: &status,

			// 使用限制 - 用户券通常只能用一次
			MaxUsagePerUser: utils.IntPtr(1),
			MaxUsageTotal:   utils.IntPtr(1),
			UsageCount:      utils.IntPtr(0),

			// 时间限制
			StartedAt: nil, // 用户券发放时即生效
			EndedAt:   u.ExpiredAt,

			// 适用范围 - 用户券通常没有车型/地域限制，使用空数组表示适用所有
			VehicleFilters: []*VehicleFilter{},
			ServiceAreas:   []string{},

			// 自动应用设置 - 用户券需要手动输入代码，不自动应用
			AutoApply: utils.IntPtr(0),
			IsGlobal:  utils.IntPtr(0),

			// 审批信息 - 用户券已经通过发放流程，直接设为已审批
			ApprovedBy: u.IssuedBy,
			CreatedBy:  u.IssuedBy,
		},
	}

	return priceRule
}

func DisabledUserPromotionByUserID(tx *gorm.DB, userID string) error {
	if userID == "" {
		return nil
	}
	values := &UserPromotionValues{
		Status: utils.StringPtr(protocol.StatusInactive),
	}
	err := tx.Model(&UserPromotion{}).
		Where("user_id = ?", userID).
		UpdateColumns(values).Error
	if err != nil {
		log.Printf("Failed to disabled user promotions for user %s: %v", userID, err)
		return err
	}

	return nil
}

// GetAvailablePromotionsByUser 根据用户ID和优惠码查询可用的用户优惠券
func GetAvailablePromotionsByUser(userID string, codes []string) []*UserPromotion {
	if len(codes) == 0 {
		return nil
	}

	var promotions []*UserPromotion
	err := DB.Where("user_id = ? AND code IN ? AND status = ?",
		userID, codes, protocol.StatusAvailable).Find(&promotions).Error
	if err != nil {
		log.Printf("Failed to query user promotions: %v", err)
		return nil
	}

	// 进一步筛选可用的优惠券（检查过期时间等）
	var availablePromotions []*UserPromotion
	for _, promo := range promotions {
		if promo.CanUse() {
			availablePromotions = append(availablePromotions, promo)
		}
	}

	return availablePromotions
}

// GetUserPromotionByID 根据ID获取用户优惠券
func GetUserPromotionByID(id string) *UserPromotion {
	if id == "" {
		return nil
	}

	var promotion UserPromotion
	err := DB.Where("promotion_id = ?", id).First(&promotion).Error
	if err != nil {
		log.Printf("Failed to query user promotion: %v", err)
		return nil
	}

	return &promotion
}

// ResetUserPromotionsByIDs 批量重置用户优惠券状态为可用
// 传入用户优惠券ID列表，将已使用的优惠券恢复为可用状态
// 使用事务和直接批量更新，无需先查询，提高性能
func ResetUserPromotionsByIDs(tx *gorm.DB, userPromotionIDs []string) error {
	if len(userPromotionIDs) == 0 {
		return nil
	}

	// 直接批量更新：只更新 is_used = 1 的记录
	updateData := map[string]any{
		"status":      protocol.StatusAvailable,
		"is_used":     0,
		"order_id":    nil,
		"used_amount": 0.0,
		"used_at":     nil,
		"updated_at":  utils.TimeNowMilli(),
	}

	result := tx.Model(&UserPromotion{}).
		Where("promotion_id IN ? AND is_used = 1", userPromotionIDs).
		Updates(updateData)

	if result.Error != nil {
		log.Printf("Failed to batch reset user promotion status: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully reset %d user promotions to available status", result.RowsAffected)
	return nil
}

// UseUserPromotionsByIDs 批量将用户优惠券标记为已使用
// 传入用户优惠券ID列表和使用信息，将可用的优惠券标记为已使用状态
// 使用事务和直接批量更新，提高性能
func UseUserPromotionsByIDs(tx *gorm.DB, orderID string, promotionUsageMap map[string]float64) error {
	if len(promotionUsageMap) == 0 {
		return nil
	}

	// 基础更新数据
	now := utils.TimeNowMilli()
	baseUpdateData := map[string]any{
		"status":     protocol.StatusUsed,
		"is_used":    1,
		"order_id":   orderID,
		"used_at":    now,
		"updated_at": now,
	}

	// 为每个优惠券执行更新（因为used_amount字段每个不同，需要分别处理）
	successCount := 0
	for promotionID, usedAmount := range promotionUsageMap {
		updateData := make(map[string]any)
		maps.Copy(updateData, baseUpdateData)
		updateData["used_amount"] = usedAmount

		result := tx.Model(&UserPromotion{}).
			Where("promotion_id = ? AND is_used = 0", promotionID).
			Updates(updateData)

		if result.Error != nil {
			log.Printf("Failed to batch mark user promotions as used: ID=%s, error=%v", promotionID, result.Error)
			return result.Error
		}

		if result.RowsAffected > 0 {
			successCount++
		}
	}

	log.Printf("Successfully marked %d user promotions as used", successCount)
	return nil
}
func UseUserPromotionByID(tx *gorm.DB, promotionID, orderID string, usedAmount float64) error {
	if promotionID == "" || orderID == "" {
		return nil
	}

	now := utils.TimeNowMilli()
	updateData := map[string]any{
		"status":      protocol.StatusUsed,
		"is_used":     1,
		"order_id":    orderID,
		"used_amount": usedAmount,
		"used_at":     now,
		"updated_at":  now,
	}

	result := tx.Model(&UserPromotion{}).
		Where("promotion_id = ? AND is_used = 0", promotionID).
		Updates(updateData)

	if result.Error != nil {
		log.Printf("Failed to mark user promotion as used: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("No available user promotion found for marking: ID=%s", promotionID)
	} else {
		log.Printf("Successfully marked user promotion as used: ID=%s", promotionID)
	}

	return nil
}

// CheckUserHasWelcomeCoupon 检查用户是否已有欢迎优惠券
func CheckUserHasWelcomeCoupon(userID string) bool {
	var count int64
	GetDB().Model(&UserPromotion{}).
		Where("user_id = ? AND source = ?", userID, protocol.UserPromotionSourceWelcome).
		Count(&count)
	return count > 0
}

// CreateUserPromotionInDB 创建用户优惠券并保存到数据库
func CreateUserPromotionInDB(userPromotion *UserPromotion) error {
	return GetDB().Create(userPromotion).Error
}

// GetWelcomePromotionTemplate 获取欢迎优惠券模板
func GetWelcomePromotionTemplate(promotionCode string) *Promotion {
	var promotion Promotion
	err := GetDB().Where("code = ? AND status = ?", promotionCode, protocol.StatusActive).First(&promotion).Error
	if err != nil {
		return nil
	}
	return &promotion
}

// GetOrCreateDefaultWelcomePromotion 获取或创建默认欢迎优惠券模板
func GetOrCreateDefaultWelcomePromotion() *Promotion {
	const defaultWelcomeCode = "WELCOME_NEW_USER"

	// 先尝试获取现有的
	promotion := GetWelcomePromotionTemplate(defaultWelcomeCode)
	if promotion != nil {
		return promotion
	}

	// 如果不存在，创建默认的欢迎优惠券模板
	promotion = &Promotion{
		PromotionID: utils.GeneratePromotionID(),
		PromotionValues: &PromotionValues{
			Code:              utils.StringPtr(defaultWelcomeCode),
			Title:             utils.StringPtr("New User Welcome Coupon"),
			Description:       utils.StringPtr("Exclusive coupon for new registered users, enjoy discount on first use"),
			DiscountType:      utils.StringPtr(protocol.PromoDiscountTypePercentage),
			DiscountValue:     utils.Float64Ptr(10.0),  // 10% 折扣
			MaxDiscountAmount: utils.Float64Ptr(50.0),  // 最大优惠50元
			MinOrderAmount:    utils.Float64Ptr(100.0), // 最小订单100元
			Status:            utils.StringPtr(protocol.StatusActive),
			UsageLimit:        utils.IntPtr(10000), // 总使用限制: 10,000张
			UserUsageLimit:    utils.IntPtr(1),     // 每用户限制: 1张
			StartDate:         utils.Int64Ptr(utils.TimeNowMilli()),
			EndDate:           utils.Int64Ptr(utils.TimeNowMilli() + 365*24*3600*1000), // 1年有效期
			ValidCities:       utils.StringPtr("all"),
			ValidVehicleTypes: utils.StringPtr("all"),
			ApprovalStatus:    utils.StringPtr(protocol.StatusApproved),
			ApprovedBy:        utils.StringPtr("system"),
			ApprovedAt:        utils.Int64Ptr(utils.TimeNowMilli()),
			CreatedBy:         utils.StringPtr("system"),
			Priority:          utils.IntPtr(0), // 优先级: 0
			Tags:              utils.StringPtr(""),
		},
	}

	// 保存到数据库
	if err := GetDB().Create(promotion).Error; err != nil {
		log.Printf("Failed to create default welcome promotion: %v", err)
		return nil
	}

	log.Printf("Created default welcome promotion: %s", defaultWelcomeCode)
	return promotion
}
