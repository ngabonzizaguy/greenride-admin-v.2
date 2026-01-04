package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"log"
	"time"
)

// Promotion 促销代码表 - 简化版优惠券管理
type Promotion struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PromotionID string `json:"promotion_id" gorm:"column:promotion_id;type:varchar(64);uniqueIndex"`
	*PromotionValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PromotionValues struct {
	// 基本信息
	Code        *string `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex"` // 优惠码
	Title       *string `json:"title" gorm:"column:title;type:varchar(255)"`          // 优惠券标题
	Description *string `json:"description" gorm:"column:description;type:text"`      // 优惠券描述

	// 优惠信息 - 核心字段
	DiscountType      *string  `json:"discount_type" gorm:"column:discount_type;type:varchar(20)"`               // percentage, fixed_amount
	DiscountValue     *float64 `json:"discount_value" gorm:"column:discount_value;type:decimal(10,2)"`           // 优惠金额或百分比
	MaxDiscountAmount *float64 `json:"max_discount_amount" gorm:"column:max_discount_amount;type:decimal(10,2)"` // 最大优惠金额
	MinOrderAmount    *float64 `json:"min_order_amount" gorm:"column:min_order_amount;type:decimal(10,2)"`       // 最小订单金额

	// 使用限制
	UsageLimit     *int `json:"usage_limit" gorm:"column:usage_limit;type:int"`                     // 总使用次数限制
	UsageCount     *int `json:"usage_count" gorm:"column:usage_count;type:int;default:0"`           // 已使用次数
	UserUsageLimit *int `json:"user_usage_limit" gorm:"column:user_usage_limit;type:int;default:1"` // 单用户使用限制

	// 时间限制
	StartDate *int64 `json:"start_date" gorm:"column:start_date"` // 开始时间
	EndDate   *int64 `json:"end_date" gorm:"column:end_date"`     // 结束时间

	// 地域和车型限制
	ValidCities       *string `json:"valid_cities" gorm:"column:valid_cities;type:text"`               // 有效城市
	ValidVehicleTypes *string `json:"valid_vehicle_types" gorm:"column:valid_vehicle_types;type:text"` // 有效车型

	// 状态管理
	Status       *string `json:"status" gorm:"column:status;type:varchar(30);index;default:'active'"` // active, inactive, expired, suspended, deleted
	StatusReason *string `json:"status_reason" gorm:"column:status_reason;type:varchar(255)"`         // 状态变更原因

	// 时间戳
	ActivatedAt *int64 `json:"activated_at" gorm:"column:activated_at"` // 激活时间
	SuspendedAt *int64 `json:"suspended_at" gorm:"column:suspended_at"` // 暂停时间
	DeletedAt   *int64 `json:"deleted_at" gorm:"column:deleted_at"`     // 删除时间

	// 审批相关
	ApprovalStatus *string `json:"approval_status" gorm:"column:approval_status;type:varchar(30);default:'pending'"` // pending, approved, rejected
	ApprovedBy     *string `json:"approved_by" gorm:"column:approved_by;type:varchar(100)"`                          // 审批者
	ApprovedAt     *int64  `json:"approved_at" gorm:"column:approved_at"`                                            // 审批时间
	ApprovalNotes  *string `json:"approval_notes" gorm:"column:approval_notes;type:text"`                            // 审批备注

	// 管理信息
	CreatedBy *string `json:"created_by" gorm:"column:created_by;type:varchar(100)"` // 创建者
	UpdatedBy *string `json:"updated_by" gorm:"column:updated_by;type:varchar(100)"` // 更新者
	Priority  *int    `json:"priority" gorm:"column:priority;type:int;default:0"`    // 优先级
	Tags      *string `json:"tags" gorm:"column:tags;type:varchar(500)"`             // 标签

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Promotion) TableName() string {
	return "t_promotions"
}

// 创建新的促销代码对象
func NewPromotion() *Promotion {
	return &Promotion{
		PromotionID: utils.GeneratePromotionID(),
		PromotionValues: &PromotionValues{
			DiscountType:   utils.StringPtr(protocol.PromoDiscountTypePercentage),
			DiscountValue:  utils.Float64Ptr(10.0),
			UsageLimit:     utils.IntPtr(1000),
			UsageCount:     utils.IntPtr(0),
			UserUsageLimit: utils.IntPtr(1),
			Status:         utils.StringPtr(protocol.StatusActive),
			ApprovalStatus: utils.StringPtr(protocol.StatusPending),
			Priority:       utils.IntPtr(0),
		},
	}
}

// Getter 方法
func (p *PromotionValues) GetCode() string {
	if p.Code == nil {
		return ""
	}
	return *p.Code
}

func (p *PromotionValues) GetTitle() string {
	if p.Title == nil {
		return ""
	}
	return *p.Title
}

func (p *PromotionValues) GetDescription() string {
	if p.Description == nil {
		return ""
	}
	return *p.Description
}

func (p *PromotionValues) GetDiscountType() string {
	if p.DiscountType == nil {
		return protocol.PromoDiscountTypePercentage
	}
	return *p.DiscountType
}

func (p *PromotionValues) GetDiscountValue() float64 {
	if p.DiscountValue == nil {
		return 0.0
	}
	return *p.DiscountValue
}

func (p *PromotionValues) GetMaxDiscountAmount() float64 {
	if p.MaxDiscountAmount == nil {
		return 0.0
	}
	return *p.MaxDiscountAmount
}

func (p *PromotionValues) GetMinOrderAmount() float64 {
	if p.MinOrderAmount == nil {
		return 0.0
	}
	return *p.MinOrderAmount
}

func (p *PromotionValues) GetUsageLimit() int {
	if p.UsageLimit == nil {
		return 0
	}
	return *p.UsageLimit
}

func (p *PromotionValues) GetUsageCount() int {
	if p.UsageCount == nil {
		return 0
	}
	return *p.UsageCount
}

func (p *PromotionValues) GetUserUsageLimit() int {
	if p.UserUsageLimit == nil {
		return 1
	}
	return *p.UserUsageLimit
}

func (p *PromotionValues) GetStatus() string {
	if p.Status == nil {
		return protocol.StatusInactive
	}
	return *p.Status
}

func (p *PromotionValues) GetCreatedBy() string {
	if p.CreatedBy == nil {
		return ""
	}
	return *p.CreatedBy
}

func (p *PromotionValues) GetStartDate() int64 {
	if p.StartDate == nil {
		return 0
	}
	return *p.StartDate
}

func (p *PromotionValues) GetEndDate() int64 {
	if p.EndDate == nil {
		return 0
	}
	return *p.EndDate
}

func (p *PromotionValues) GetValidCities() string {
	if p.ValidCities == nil {
		return ""
	}
	return *p.ValidCities
}

func (p *PromotionValues) GetValidVehicleTypes() string {
	if p.ValidVehicleTypes == nil {
		return ""
	}
	return *p.ValidVehicleTypes
}

func (p *PromotionValues) GetStatusReason() string {
	if p.StatusReason == nil {
		return ""
	}
	return *p.StatusReason
}

func (p *PromotionValues) GetActivatedAt() int64 {
	if p.ActivatedAt == nil {
		return 0
	}
	return *p.ActivatedAt
}

func (p *PromotionValues) GetSuspendedAt() int64 {
	if p.SuspendedAt == nil {
		return 0
	}
	return *p.SuspendedAt
}

func (p *PromotionValues) GetDeletedAt() int64 {
	if p.DeletedAt == nil {
		return 0
	}
	return *p.DeletedAt
}

func (p *PromotionValues) GetApprovalStatus() string {
	if p.ApprovalStatus == nil {
		return protocol.StatusPending
	}
	return *p.ApprovalStatus
}

func (p *PromotionValues) GetApprovedBy() string {
	if p.ApprovedBy == nil {
		return ""
	}
	return *p.ApprovedBy
}

func (p *PromotionValues) GetApprovedAt() int64 {
	if p.ApprovedAt == nil {
		return 0
	}
	return *p.ApprovedAt
}

func (p *PromotionValues) GetApprovalNotes() string {
	if p.ApprovalNotes == nil {
		return ""
	}
	return *p.ApprovalNotes
}

func (p *PromotionValues) GetUpdatedBy() string {
	if p.UpdatedBy == nil {
		return ""
	}
	return *p.UpdatedBy
}

func (p *PromotionValues) GetPriority() int {
	if p.Priority == nil {
		return 0
	}
	return *p.Priority
}

func (p *PromotionValues) GetTags() string {
	if p.Tags == nil {
		return ""
	}
	return *p.Tags
}

// Setter 方法
func (p *PromotionValues) SetCode(code string) *PromotionValues {
	p.Code = &code
	return p
}

func (p *PromotionValues) SetTitle(title string) *PromotionValues {
	p.Title = &title
	return p
}

func (p *PromotionValues) SetDescription(description string) *PromotionValues {
	p.Description = &description
	return p
}

func (p *PromotionValues) SetDiscountInfo(discountType string, value, maxAmount float64) *PromotionValues {
	p.DiscountType = &discountType
	p.DiscountValue = &value
	if maxAmount > 0 {
		p.MaxDiscountAmount = &maxAmount
	}
	return p
}

func (p *PromotionValues) SetMinOrderAmount(minAmount float64) *PromotionValues {
	p.MinOrderAmount = &minAmount
	return p
}

func (p *PromotionValues) SetUsageLimits(totalLimit, userLimit int) *PromotionValues {
	p.UsageLimit = &totalLimit
	p.UserUsageLimit = &userLimit
	return p
}

func (p *PromotionValues) SetStatus(status string) *PromotionValues {
	p.Status = &status
	return p
}
func (p *PromotionValues) SetCreatedBy(createdBy string) *PromotionValues {
	p.CreatedBy = &createdBy
	return p
}

func (p *PromotionValues) SetDiscountType(discountType string) *PromotionValues {
	p.DiscountType = &discountType
	return p
}

func (p *PromotionValues) SetDiscountValue(value float64) *PromotionValues {
	p.DiscountValue = &value
	return p
}

func (p *PromotionValues) SetMaxDiscountAmount(maxAmount float64) *PromotionValues {
	p.MaxDiscountAmount = &maxAmount
	return p
}

func (p *PromotionValues) SetUsageLimit(limit int) *PromotionValues {
	p.UsageLimit = &limit
	return p
}

func (p *PromotionValues) SetUsageCount(count int) *PromotionValues {
	p.UsageCount = &count
	return p
}

func (p *PromotionValues) SetUserUsageLimit(limit int) *PromotionValues {
	p.UserUsageLimit = &limit
	return p
}

func (p *PromotionValues) SetStartDate(startDate int64) *PromotionValues {
	p.StartDate = &startDate
	return p
}

func (p *PromotionValues) SetEndDate(endDate int64) *PromotionValues {
	p.EndDate = &endDate
	return p
}

func (p *PromotionValues) SetValidCities(cities string) *PromotionValues {
	p.ValidCities = &cities
	return p
}

func (p *PromotionValues) SetValidVehicleTypes(types string) *PromotionValues {
	p.ValidVehicleTypes = &types
	return p
}

func (p *PromotionValues) SetStatusReason(reason string) *PromotionValues {
	p.StatusReason = &reason
	return p
}

func (p *PromotionValues) SetActivatedAt(activatedAt int64) *PromotionValues {
	p.ActivatedAt = &activatedAt
	return p
}

func (p *PromotionValues) SetSuspendedAt(suspendedAt int64) *PromotionValues {
	p.SuspendedAt = &suspendedAt
	return p
}

func (p *PromotionValues) SetDeletedAt(deletedAt int64) *PromotionValues {
	p.DeletedAt = &deletedAt
	return p
}

func (p *PromotionValues) SetApprovalStatus(status string) *PromotionValues {
	p.ApprovalStatus = &status
	return p
}

func (p *PromotionValues) SetApprovedBy(approvedBy string) *PromotionValues {
	p.ApprovedBy = &approvedBy
	return p
}

func (p *PromotionValues) SetApprovedAt(approvedAt int64) *PromotionValues {
	p.ApprovedAt = &approvedAt
	return p
}

func (p *PromotionValues) SetApprovalNotes(notes string) *PromotionValues {
	p.ApprovalNotes = &notes
	return p
}

func (p *PromotionValues) SetUpdatedBy(updatedBy string) *PromotionValues {
	p.UpdatedBy = &updatedBy
	return p
}

func (p *PromotionValues) SetPriority(priority int) *PromotionValues {
	p.Priority = &priority
	return p
}

func (p *PromotionValues) SetTags(tags string) *PromotionValues {
	p.Tags = &tags
	return p
}

// 业务方法
func (p *Promotion) IsActive() bool {
	return p.GetStatus() == protocol.StatusActive
}

func (p *Promotion) IsUsageExceeded() bool {
	usageLimit := p.GetUsageLimit()
	if usageLimit <= 0 {
		return false
	}
	return p.GetUsageCount() >= usageLimit
}

func (p *Promotion) IsValid() bool {
	return p.IsActive() && !p.IsUsageExceeded()
}

func (p *Promotion) IsPercentageDiscount() bool {
	return p.GetDiscountType() == protocol.PromoDiscountTypePercentage
}

func (p *Promotion) IsFixedAmountDiscount() bool {
	return p.GetDiscountType() == protocol.PromoDiscountTypeFixedAmount
}

// 计算折扣金额
func (p *Promotion) CalculateDiscount(orderAmount float64) float64 {
	if !p.IsValid() {
		return 0.0
	}

	// 检查最小订单金额
	if p.PromotionValues.MinOrderAmount != nil && orderAmount < *p.PromotionValues.MinOrderAmount {
		return 0.0
	}

	discountValue := p.GetDiscountValue()

	switch p.GetDiscountType() {
	case protocol.PromoDiscountTypePercentage:
		discount := orderAmount * discountValue / 100.0
		// 应用最大折扣限制
		if p.PromotionValues.MaxDiscountAmount != nil && discount > *p.PromotionValues.MaxDiscountAmount {
			return *p.PromotionValues.MaxDiscountAmount
		}
		return discount

	case protocol.PromoDiscountTypeFixedAmount:
		// 固定金额折扣不能超过订单金额
		if discountValue > orderAmount {
			return orderAmount
		}
		return discountValue

	default:
		return 0.0
	}
}

// 使用优惠券
func (p *PromotionValues) IncrementUsage() *PromotionValues {
	count := p.GetUsageCount()
	count++
	p.UsageCount = &count
	return p
}

// 状态管理
func (p *PromotionValues) Activate() *PromotionValues {
	p.SetStatus(protocol.StatusActive)
	return p
}

func (p *PromotionValues) Deactivate() *PromotionValues {
	p.SetStatus(protocol.StatusInactive)
	return p
}

func (p *PromotionValues) MarkExpired() *PromotionValues {
	p.SetStatus(protocol.StatusExpired)
	return p
}

func (p *PromotionValues) Suspend() *PromotionValues {
	p.SetStatus(protocol.StatusSuspended)
	return p
}

func (p *PromotionValues) Delete() *PromotionValues {
	p.SetStatus(protocol.StatusDeleted)
	return p
}

// 审批状态管理
func (p *PromotionValues) MarkPendingApproval() *PromotionValues {
	p.SetApprovalStatus(protocol.StatusPending)
	return p
}

func (p *PromotionValues) Approve(approvedBy string, notes string) *PromotionValues {
	currentTime := time.Now().Unix()
	p.SetApprovalStatus(protocol.StatusApproved).
		SetApprovedBy(approvedBy).
		SetApprovedAt(currentTime).
		SetApprovalNotes(notes)
	return p
}

func (p *PromotionValues) Reject(rejectedBy string, notes string) *PromotionValues {
	currentTime := time.Now().Unix()
	p.SetApprovalStatus(protocol.StatusRejected).
		SetApprovedBy(rejectedBy).
		SetApprovedAt(currentTime).
		SetApprovalNotes(notes)
	return p
}

func GetPromotionsByUserID(userID string) []*Promotion {
	var promos []*Promotion
	err := GetDB().Where("created_by = ?", userID).Find(&promos).Error
	if err != nil {
		return nil
	}
	return promos
}

func GetPromotionByID(id string) *Promotion {
	var promo Promotion
	err := GetDB().Where("promotion_id = ?", id).First(&promo).Error
	if err != nil {
		return nil
	}
	return &promo
}

// 便捷创建方法
func NewPercentageDiscountCode(code, title string, percentage float64, maxAmount float64) *Promotion {
	promo := NewPromotion()
	promo.SetCode(code).
		SetTitle(title).
		SetDiscountInfo(protocol.PromoDiscountTypePercentage, percentage, maxAmount).
		SetUsageLimits(1000, 1).
		Activate()

	return promo
}

func NewFixedAmountDiscountCode(code, title string, amount float64, minOrder float64) *Promotion {
	promo := NewPromotion()
	promo.SetCode(code).
		SetTitle(title).
		SetDiscountInfo(protocol.PromoDiscountTypeFixedAmount, amount, 0).
		SetMinOrderAmount(minOrder).
		SetUsageLimits(500, 1).
		Activate()

	return promo
}

type Promotions []*Promotion

func (p *Promotions) Protocol() []*protocol.Promotion {
	var list []*protocol.Promotion
	for _, promo := range *p {
		list = append(list, promo.Protocol())
	}
	return list
}

// Protocol 转换为协议层结构体
func (p *Promotion) Protocol() *protocol.Promotion {
	return &protocol.Promotion{
		PromotionID: p.PromotionID,
		Code:        p.GetCode(),
		Title:       p.GetTitle(),
		Description: p.GetDescription(),

		DiscountType:      p.GetDiscountType(),
		DiscountValue:     p.GetDiscountValue(),
		MaxDiscountAmount: p.GetMaxDiscountAmount(),
		MinOrderAmount:    p.GetMinOrderAmount(),

		UsageLimit:     p.GetUsageLimit(),
		UsageCount:     p.GetUsageCount(),
		UserUsageLimit: p.GetUserUsageLimit(),

		StartDate:         p.GetStartDate(),
		EndDate:           p.GetEndDate(),
		ValidCities:       p.GetValidCities(),
		ValidVehicleTypes: p.GetValidVehicleTypes(),

		Status:         p.GetStatus(),
		ApprovalStatus: p.GetApprovalStatus(),
		ApprovedBy:     p.GetApprovedBy(),

		CreatedBy: p.GetCreatedBy(),
		Priority:  p.GetPriority(),
		Tags:      p.GetTags(),

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.PromotionValues.UpdatedAt,
	}
}

// GetReferralPromotionTemplate 获取推荐优惠券模板
func GetReferralPromotionTemplate(promotionCode string) *Promotion {
	var promotion Promotion
	err := GetDB().Where("code = ? AND status = ?", promotionCode, protocol.StatusActive).First(&promotion).Error
	if err != nil {
		return nil
	}
	return &promotion
}

// GetOrCreateDefaultReferralPromotion 获取或创建默认推荐优惠券模板
func GetOrCreateDefaultReferralPromotion() *Promotion {
	const defaultReferralCode = "INVITE_NEW_USER"

	// 先尝试获取现有的
	promotion := GetReferralPromotionTemplate(defaultReferralCode)
	if promotion != nil {
		return promotion
	}

	// 如果不存在，创建默认的推荐优惠券模板
	promotion = &Promotion{
		PromotionID: utils.GeneratePromotionID(),
		PromotionValues: &PromotionValues{
			Code:              utils.StringPtr(defaultReferralCode),
			Title:             utils.StringPtr("邀请新用户奖励券"),
			Description:       utils.StringPtr("成功邀请新用户注册并完成首单后，邀请者可获得此优惠券"),
			DiscountType:      utils.StringPtr(protocol.PromoDiscountTypeFixedAmount),
			DiscountValue:     utils.Float64Ptr(20.0), // 固定金额20元
			MaxDiscountAmount: utils.Float64Ptr(20.0), // 最大优惠20元
			MinOrderAmount:    utils.Float64Ptr(50.0), // 最小订单50元
			Status:            utils.StringPtr(protocol.StatusActive),
			UsageLimit:        utils.IntPtr(50000), // 总使用限制: 50,000张
			UserUsageLimit:    utils.IntPtr(5),     // 每用户限制: 5张
			StartDate:         utils.Int64Ptr(utils.TimeNowMilli()),
			EndDate:           utils.Int64Ptr(utils.TimeNowMilli() + 365*24*3600*1000), // 1年有效期
			ValidCities:       utils.StringPtr("all"),
			ValidVehicleTypes: utils.StringPtr("all"),
			ApprovalStatus:    utils.StringPtr(protocol.StatusApproved),
			ApprovedBy:        utils.StringPtr("system"),
			ApprovedAt:        utils.Int64Ptr(utils.TimeNowMilli()),
			CreatedBy:         utils.StringPtr("system"),
			Priority:          utils.IntPtr(1), // 优先级: 1
			Tags:              utils.StringPtr("invite,referral,new_user"),
		},
	}

	// 保存到数据库
	if err := GetDB().Create(promotion).Error; err != nil {
		log.Printf("Failed to create default referral promotion: %v", err)
		return nil
	}

	log.Printf("Created default referral promotion: %s", defaultReferralCode)
	return promotion
}
