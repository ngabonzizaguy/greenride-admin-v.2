package models

import (
	"errors"
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Order 通用订单表 - 支持多种订单类型的抽象订单
type Order struct {
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	OrderID string `json:"order_id" gorm:"column:order_id;type:varchar(64);uniqueIndex"`
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*OrderValues
	Details   *OrderDetail `json:"details,omitempty" gorm:"-"` // 订单详情，不保存到数据库
	CreatedAt int64        `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type OrderValues struct {
	// 订单基础信息
	OrderType *string `json:"order_type" gorm:"column:order_type;type:varchar(32);index"` // ride, delivery, shopping

	// 用户关联信息
	UserID     *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`         // 下单用户ID
	ProviderID *string `json:"provider_id" gorm:"column:provider_id;type:varchar(64);index"` // 服务提供者ID

	// 订单状态
	Status        *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'requested'"`
	PaymentStatus *string `json:"payment_status" gorm:"column:payment_status;type:varchar(32);index;default:''"`
	ScheduleType  *string `json:"schedule_type" gorm:"column:schedule_type;type:varchar(32);default:'instant'"` // instant, scheduled

	// 金额信息
	Currency            *string          `json:"currency" gorm:"column:currency;type:varchar(3);default:'USD'"`
	OriginalAmount      *decimal.Decimal `json:"original_amount" gorm:"column:original_amount;type:decimal(20,6)"`             // 预估原始金额
	DiscountedAmount    *decimal.Decimal `json:"discounted_amount" gorm:"column:discounted_amount;type:decimal(20,6)"`         // 优惠后总金额
	PaymentAmount       *decimal.Decimal `json:"payment_amount" gorm:"column:payment_amount;type:decimal(20,6)"`               // 最终付款金额
	TotalDiscountAmount *decimal.Decimal `json:"total_discount_amount" gorm:"column:total_discount_amount;type:decimal(20,6)"` // 总优惠金额
	PlatformFee         *decimal.Decimal `json:"platform_fee" gorm:"column:platform_fee;type:decimal(20,6)"`                   // 支付信息
	PaymentMethod       *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32)"`                 // card, cash, wallet, visa, mastercard, paypal
	PaymentID           *string          `json:"payment_id" gorm:"column:payment_id;type:varchar(64)"`
	ChannelPaymentID    *string          `json:"channel_payment_id" gorm:"column:channel_payment_id;type:varchar(128)"`     // 渠道交易ID
	PaymentResult       *string          `json:"payment_result" gorm:"column:payment_result;type:text"`                     // 支付结果详情
	PaymentRedirectURL  *string          `json:"payment_redirect_url" gorm:"column:payment_redirect_url;type:varchar(512)"` // 支付重定向URL

	// 优惠信息
	PromoCodes        []string         `json:"promo_codes" gorm:"column:promo_codes;type:json;serializer:json"`
	PromoDiscount     *decimal.Decimal `json:"promo_discount" gorm:"column:promo_discount;type:decimal(20,6)"`
	UserPromoDiscount *decimal.Decimal `json:"user_promo_discount" gorm:"column:user_promo_discount;type:decimal(20,6)"`      // 用户优惠券折扣金额
	UserPromotionIDs  []string         `json:"user_promotion_ids" gorm:"column:user_promotion_ids;type:json;serializer:json"` // 使用的用户优惠券ID列表

	Sandbox *int `json:"sandbox" gorm:"column:sandbox;type:tinyint(1);default:0"` // 是否为沙箱订单

	// 时间信息
	ScheduledAt *int64 `json:"scheduled_at" gorm:"column:scheduled_at"` // 预约时间
	AcceptedAt  *int64 `json:"accepted_at" gorm:"column:accepted_at"`   // 接单时间
	StartedAt   *int64 `json:"started_at" gorm:"column:started_at"`     // 开始服务时间
	EndedAt     *int64 `json:"ended_at" gorm:"column:ended_at"`         // 结束时间
	CompletedAt *int64 `json:"completed_at" gorm:"column:completed_at"` // 完成时间
	CancelledAt *int64 `json:"cancelled_at" gorm:"column:cancelled_at"` // 取消时间
	ExpiredAt   *int64 `json:"expired_at" gorm:"column:expired_at"`     // 过期时间

	// 取消信息
	CancelledBy     *string          `json:"cancelled_by" gorm:"column:cancelled_by;type:varchar(64)"` // 取消者用户ID
	CancelReason    *string          `json:"cancel_reason" gorm:"column:cancel_reason;type:varchar(255)"`
	CancellationFee *decimal.Decimal `json:"cancellation_fee" gorm:"column:cancellation_fee;type:decimal(20,6)"`

	// 评价信息 - 已改用单独的Rating表，不再保存评价信息

	*OrderDispatchValues

	// 扩展信息
	Metadata map[string]any `json:"metadata" gorm:"column:metadata;type:json;serializer:json"` // JSON格式的扩展信息
	Notes    *string        `json:"notes" gorm:"column:notes;type:text"`                       // 备注信息
	Version  *int64         `json:"version" gorm:"column:version;default:1"`                   // 版本号（乐观锁）

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type OrderDispatchValues struct { // 派单相关字段
	DispatchStatus      *string `json:"dispatch_status" gorm:"column:dispatch_status;type:varchar(32);default:'not_started'"` // not_started, dispatching, completed, failed
	CurrentRound        *int    `json:"current_round" gorm:"column:current_round;default:0"`                                  // 当前派单轮次
	MaxRounds           *int    `json:"max_rounds" gorm:"column:max_rounds;default:4"`                                        // 最大派单轮次
	DispatchStartedAt   *int64  `json:"dispatch_started_at" gorm:"column:dispatch_started_at"`                                // 派单开始时间戳
	LastDispatchedAt    *int64  `json:"last_dispatched_at" gorm:"column:last_dispatched_at"`                                  // 最后一次派单时间戳
	AutoDispatchEnabled *bool   `json:"auto_dispatch_enabled" gorm:"column:auto_dispatch_enabled;default:true"`               // 是否启用自动派单

	// 策略配置（JSON存储）
	DispatchStrategy *string `json:"dispatch_strategy" gorm:"column:dispatch_strategy;type:json"` // 当前轮次策略配置
	NextStrategy     *string `json:"next_strategy" gorm:"column:next_strategy;type:json"`         // 下轮策略配置
}

func (Order) TableName() string {
	return "t_orders"
}

// 创建新的订单对象
func NewOrder() *Order {
	return &Order{
		OrderID: utils.GenerateOrderID(),
		Salt:    utils.GenerateSalt(),
		OrderValues: &OrderValues{
			Status:        utils.StringPtr(protocol.StatusRequested),
			PaymentStatus: utils.StringPtr(protocol.StatusPending),
			ScheduleType:  utils.StringPtr(protocol.ScheduleTypeInstant),
			Currency:      utils.StringPtr("USD"),
			Version:       utils.Int64Ptr(1),
			OrderDispatchValues: &OrderDispatchValues{
				AutoDispatchEnabled: utils.BoolPtr(true),
			},
		},
	}
}

// SetValues 更新OrderValues中的非nil值
func (o *OrderValues) SetValues(values *OrderValues) {
	if values == nil {
		return
	}
	// 基础信息
	if values.OrderType != nil {
		o.OrderType = values.OrderType
	}

	// 用户关联信息
	if values.UserID != nil {
		o.UserID = values.UserID
	}
	if values.ProviderID != nil {
		o.ProviderID = values.ProviderID
	}

	// 订单状态
	if values.Status != nil {
		o.Status = values.Status
	}
	if values.PaymentStatus != nil {
		o.PaymentStatus = values.PaymentStatus
	}
	if values.ScheduleType != nil {
		o.ScheduleType = values.ScheduleType
	}

	// 金额信息
	if values.Currency != nil {
		o.Currency = values.Currency
	}
	if values.OriginalAmount != nil {
		o.OriginalAmount = values.OriginalAmount
	}
	if values.DiscountedAmount != nil {
		o.DiscountedAmount = values.DiscountedAmount
	}
	if values.PaymentAmount != nil {
		o.PaymentAmount = values.PaymentAmount
	}
	if values.TotalDiscountAmount != nil {
		o.TotalDiscountAmount = values.TotalDiscountAmount
	}
	if values.PlatformFee != nil {
		o.PlatformFee = values.PlatformFee
	}

	// 支付信息
	if values.PaymentMethod != nil {
		o.PaymentMethod = values.PaymentMethod
	}
	if values.PaymentID != nil {
		o.PaymentID = values.PaymentID
	}
	if values.ChannelPaymentID != nil {
		o.ChannelPaymentID = values.ChannelPaymentID
	}
	if values.PaymentResult != nil {
		o.PaymentResult = values.PaymentResult
	}
	if values.PaymentRedirectURL != nil {
		o.PaymentRedirectURL = values.PaymentRedirectURL
	}

	// 优惠信息
	if values.PromoCodes != nil {
		o.PromoCodes = values.PromoCodes
	}
	if values.PromoDiscount != nil {
		o.PromoDiscount = values.PromoDiscount
	}
	if values.UserPromoDiscount != nil {
		o.UserPromoDiscount = values.UserPromoDiscount
	}
	if values.UserPromotionIDs != nil {
		o.UserPromotionIDs = values.UserPromotionIDs
	}
	if values.Sandbox != nil {
		o.Sandbox = values.Sandbox
	}

	// 时间信息
	if values.ScheduledAt != nil {
		o.ScheduledAt = values.ScheduledAt
	}
	if values.AcceptedAt != nil {
		o.AcceptedAt = values.AcceptedAt
	}
	if values.StartedAt != nil {
		o.StartedAt = values.StartedAt
	}
	if values.EndedAt != nil {
		o.EndedAt = values.EndedAt
	}
	if values.CompletedAt != nil {
		o.CompletedAt = values.CompletedAt
	}
	if values.ExpiredAt != nil {
		o.ExpiredAt = values.ExpiredAt
	}
	if values.CancelledAt != nil {
		o.CancelledAt = values.CancelledAt
	}

	// 取消信息
	if values.CancelledBy != nil {
		o.CancelledBy = values.CancelledBy
	}
	if values.CancelReason != nil {
		o.CancelReason = values.CancelReason
	}
	if values.CancellationFee != nil {
		o.CancellationFee = values.CancellationFee
	}

	// 派单信息 (OrderDispatchValues)
	if values.OrderDispatchValues != nil {
		if o.OrderDispatchValues == nil {
			o.OrderDispatchValues = &OrderDispatchValues{}
		}
		if values.DispatchStatus != nil {
			o.DispatchStatus = values.DispatchStatus
		}
		if values.CurrentRound != nil {
			o.CurrentRound = values.CurrentRound
		}
		if values.MaxRounds != nil {
			o.MaxRounds = values.MaxRounds
		}
		if values.DispatchStartedAt != nil {
			o.DispatchStartedAt = values.DispatchStartedAt
		}
		if values.LastDispatchedAt != nil {
			o.LastDispatchedAt = values.LastDispatchedAt
		}
		if values.AutoDispatchEnabled != nil {
			o.AutoDispatchEnabled = values.AutoDispatchEnabled
		}
		if values.DispatchStrategy != nil {
			o.DispatchStrategy = values.DispatchStrategy
		}
		if values.NextStrategy != nil {
			o.NextStrategy = values.NextStrategy
		}
	}

	// 扩展信息
	if values.Metadata != nil {
		o.Metadata = values.Metadata
	}
	if values.Notes != nil {
		o.Notes = values.Notes
	}
	if values.Version != nil {
		o.Version = values.Version
	}
	if values.UpdatedAt > 0 {
		o.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (o *OrderValues) GetExpiredAt() int64 {
	if o.ExpiredAt == nil {
		return 0
	}
	return *o.ExpiredAt
}
func (o *OrderValues) GetOrderType() string {
	if o.OrderType == nil {
		return ""
	}
	return *o.OrderType
}

func (o *OrderValues) GetUserID() string {
	if o.UserID == nil {
		return ""
	}
	return *o.UserID
}

func (o *OrderValues) GetProviderID() string {
	if o.ProviderID == nil {
		return ""
	}
	return *o.ProviderID
}

func (o *OrderValues) GetStatus() string {
	if o.Status == nil {
		return protocol.StatusRequested
	}
	return *o.Status
}

func (o *OrderValues) GetPaymentStatus() string {
	if o.PaymentStatus == nil {
		return protocol.StatusPending
	}
	return *o.PaymentStatus
}

func (o *OrderValues) GetScheduleType() string {
	if o.ScheduleType == nil {
		return protocol.ScheduleTypeInstant
	}
	return *o.ScheduleType
}

func (o *OrderValues) GetOriginalAmount() decimal.Decimal {
	if o.OriginalAmount == nil {
		return decimal.Zero
	}
	return *o.OriginalAmount
}

func (o *OrderValues) GetDiscountedAmount() decimal.Decimal {
	if o.DiscountedAmount == nil {
		return decimal.Zero
	}
	return *o.DiscountedAmount
}

func (o *OrderValues) GetPaymentAmount() decimal.Decimal {
	if o.PaymentAmount == nil {
		return decimal.Zero
	}
	return *o.PaymentAmount
}

func (o *OrderValues) GetCurrency() string {
	if o.Currency == nil {
		return "USD"
	}
	return *o.Currency
}

func (o *OrderValues) GetTotalDiscountAmount() decimal.Decimal {
	if o.TotalDiscountAmount == nil {
		return decimal.Zero
	}
	return *o.TotalDiscountAmount
}

func (o *OrderValues) GetPlatformFee() decimal.Decimal {
	if o.PlatformFee == nil {
		return decimal.Zero
	}
	return *o.PlatformFee
}

func (o *OrderValues) GetPaymentMethod() string {
	if o.PaymentMethod == nil {
		return ""
	}
	return *o.PaymentMethod
}

func (o *OrderValues) GetPaymentID() string {
	if o.PaymentID == nil {
		return ""
	}
	return *o.PaymentID
}

func (o *OrderValues) GetChannelPaymentID() string {
	if o.ChannelPaymentID == nil {
		return ""
	}
	return *o.ChannelPaymentID
}

func (o *OrderValues) GetPaymentResult() string {
	if o.PaymentResult == nil {
		return ""
	}
	return *o.PaymentResult
}

func (o *OrderValues) GetPaymentRedirectURL() string {
	if o.PaymentRedirectURL == nil {
		return ""
	}
	return *o.PaymentRedirectURL
}

func (o *OrderValues) GetPromoCodes() []string {
	if o.PromoCodes == nil {
		return []string{}
	}
	return o.PromoCodes
}

func (o *OrderValues) GetPromoDiscount() decimal.Decimal {
	if o.PromoDiscount == nil {
		return decimal.Zero
	}
	return *o.PromoDiscount
}

func (o *OrderValues) GetUserPromotionIDs() []string {
	if o.UserPromotionIDs == nil {
		return []string{}
	}
	return o.UserPromotionIDs
}

func (o *OrderValues) GetSandbox() int {
	if o.Sandbox == nil {
		return 0
	}
	return *o.Sandbox
}

func (o *OrderValues) GetScheduledAt() int64 {
	if o.ScheduledAt == nil {
		return 0
	}
	return *o.ScheduledAt
}

func (o *OrderValues) GetAcceptedAt() int64 {
	if o.AcceptedAt == nil {
		return 0
	}
	return *o.AcceptedAt
}

func (o *OrderValues) GetStartedAt() int64 {
	if o.StartedAt == nil {
		return 0
	}
	return *o.StartedAt
}

func (o *OrderValues) GetEndedAt() int64 {
	if o.EndedAt == nil {
		return 0
	}
	return *o.EndedAt
}

func (o *OrderValues) GetCompletedAt() int64 {
	if o.CompletedAt == nil {
		return 0
	}
	return *o.CompletedAt
}

func (o *OrderValues) GetCancelledAt() int64 {
	if o.CancelledAt == nil {
		return 0
	}
	return *o.CancelledAt
}

func (o *OrderValues) GetCancelledBy() string {
	if o.CancelledBy == nil {
		return ""
	}
	return *o.CancelledBy
}

func (o *OrderValues) GetCancelReason() string {
	if o.CancelReason == nil {
		return ""
	}
	return *o.CancelReason
}

func (o *OrderValues) GetCancellationFee() decimal.Decimal {
	if o.CancellationFee == nil {
		return decimal.Zero
	}
	return *o.CancellationFee
}

func (o *OrderValues) GetDispatchStatus() string {
	if o.DispatchStatus == nil {
		return "not_started"
	}
	return *o.DispatchStatus
}

func (o *OrderValues) GetCurrentRound() int {
	if o.CurrentRound == nil {
		return 0
	}
	return *o.CurrentRound
}

func (o *OrderValues) GetMaxRounds() int {
	if o.MaxRounds == nil {
		return 4
	}
	return *o.MaxRounds
}

func (o *OrderValues) GetDispatchStartedAt() int64 {
	if o.DispatchStartedAt == nil {
		return 0
	}
	return *o.DispatchStartedAt
}

func (o *OrderValues) GetLastDispatchedAt() int64 {
	if o.LastDispatchedAt == nil {
		return 0
	}
	return *o.LastDispatchedAt
}

func (o *OrderValues) GetAutoDispatchEnabled() bool {
	if o.AutoDispatchEnabled == nil {
		return true
	}
	return *o.AutoDispatchEnabled
}

func (o *OrderValues) GetDispatchStrategy() string {
	if o.DispatchStrategy == nil {
		return ""
	}
	return *o.DispatchStrategy
}

func (o *OrderValues) GetNextStrategy() string {
	if o.NextStrategy == nil {
		return ""
	}
	return *o.NextStrategy
}

func (o *OrderValues) GetMetadata() map[string]interface{} {
	if o.Metadata == nil {
		return make(map[string]interface{})
	}
	return o.Metadata
}

func (o *OrderValues) GetNotes() string {
	if o.Notes == nil {
		return ""
	}
	return *o.Notes
}

func (o *OrderValues) GetUpdatedAt() int64 {
	return o.UpdatedAt
}

// Setter 方法
func (o *OrderValues) SetExpiredAt(expiredAt int64) *OrderValues {
	o.ExpiredAt = &expiredAt
	return o
}
func (o *OrderValues) SetOrderType(orderType string) *OrderValues {
	o.OrderType = &orderType
	return o
}

func (o *OrderValues) SetUserID(userID string) *OrderValues {
	o.UserID = &userID
	return o
}

func (o *OrderValues) SetProviderID(providerID string) *OrderValues {
	o.ProviderID = &providerID
	return o
}

func (o *OrderValues) SetStatus(status string) *OrderValues {
	o.Status = &status
	return o
}

func (o *OrderValues) SetPaymentStatus(status string) *OrderValues {
	o.PaymentStatus = &status
	return o
}

func (o *OrderValues) SetScheduleType(scheduleType string) *OrderValues {
	o.ScheduleType = &scheduleType
	return o
}

func (o *OrderValues) SetAmounts(original, discounted, payment decimal.Decimal) *OrderValues {
	o.OriginalAmount = &original
	o.DiscountedAmount = &discounted
	o.PaymentAmount = &payment
	return o
}

func (o *OrderValues) SetPaymentMethod(paymentMethod string) *OrderValues {
	o.PaymentMethod = &paymentMethod
	return o
}

func (o *OrderValues) SetPaymentID(paymentID string) *OrderValues {
	o.PaymentID = &paymentID
	return o
}

func (o *OrderValues) SetChannelPaymentID(channelPaymentID string) *OrderValues {
	o.ChannelPaymentID = &channelPaymentID
	return o
}

func (o *OrderValues) SetPaymentResult(paymentResult string) *OrderValues {
	o.PaymentResult = &paymentResult
	return o
}

func (o *OrderValues) SetPaymentRedirectURL(redirectURL string) *OrderValues {
	o.PaymentRedirectURL = &redirectURL
	return o
}

func (o *OrderValues) SetEndedAt(endedAt int64) *OrderValues {
	o.EndedAt = &endedAt
	return o
}

func (o *OrderValues) SetCompletedAt(completedAt int64) *OrderValues {
	o.CompletedAt = &completedAt
	return o
}

func (o *OrderValues) SetCurrency(currency string) *OrderValues {
	o.Currency = &currency
	return o
}

func (o *OrderValues) SetOriginalAmount(amount decimal.Decimal) *OrderValues {
	o.OriginalAmount = &amount
	return o
}

func (o *OrderValues) SetDiscountedAmount(amount decimal.Decimal) *OrderValues {
	o.DiscountedAmount = &amount
	return o
}

func (o *OrderValues) SetPaymentAmount(amount decimal.Decimal) *OrderValues {
	o.PaymentAmount = &amount
	return o
}

func (o *OrderValues) SetTotalDiscountAmount(amount decimal.Decimal) *OrderValues {
	o.TotalDiscountAmount = &amount
	return o
}

func (o *OrderValues) SetPlatformFee(fee decimal.Decimal) *OrderValues {
	o.PlatformFee = &fee
	return o
}

func (o *OrderValues) SetPromoCodes(promoCodes []string) *OrderValues {
	o.PromoCodes = promoCodes
	return o
}

func (o *OrderValues) SetPromoDiscount(discount decimal.Decimal) *OrderValues {
	o.PromoDiscount = &discount
	return o
}

func (o *OrderValues) SetUserPromotionIDs(promotionIDs []string) *OrderValues {
	o.UserPromotionIDs = promotionIDs
	return o
}

func (o *OrderValues) SetSandbox(sandbox int) *OrderValues {
	o.Sandbox = &sandbox
	return o
}

func (o *OrderValues) SetScheduledAt(scheduledAt int64) *OrderValues {
	o.ScheduledAt = &scheduledAt
	return o
}

func (o *OrderValues) SetAcceptedAt(acceptedAt int64) *OrderValues {
	o.AcceptedAt = &acceptedAt
	return o
}

func (o *OrderValues) SetStartedAt(startedAt int64) *OrderValues {
	o.StartedAt = &startedAt
	return o
}

func (o *OrderValues) SetCancelledAt(cancelledAt int64) *OrderValues {
	o.CancelledAt = &cancelledAt
	return o
}

func (o *OrderValues) SetCancelledBy(cancelledBy string) *OrderValues {
	o.CancelledBy = &cancelledBy
	return o
}

func (o *OrderValues) SetCancelReason(reason string) *OrderValues {
	o.CancelReason = &reason
	return o
}

func (o *OrderValues) SetCancellationFee(fee decimal.Decimal) *OrderValues {
	o.CancellationFee = &fee
	return o
}

func (o *OrderValues) SetDispatchStatus(status string) *OrderValues {
	o.DispatchStatus = &status
	return o
}

func (o *OrderValues) SetCurrentRound(round int) *OrderValues {
	o.CurrentRound = &round
	return o
}

func (o *OrderValues) SetMaxRounds(rounds int) *OrderValues {
	o.MaxRounds = &rounds
	return o
}

func (o *OrderValues) SetDispatchStartedAt(startedAt int64) *OrderValues {
	o.DispatchStartedAt = &startedAt
	return o
}

func (o *OrderValues) SetLastDispatchedAt(dispatchedAt int64) *OrderValues {
	o.LastDispatchedAt = &dispatchedAt
	return o
}

func (o *OrderValues) SetAutoDispatchEnabled(enabled bool) *OrderValues {
	o.AutoDispatchEnabled = &enabled
	return o
}

func (o *OrderValues) SetDispatchStrategy(strategy string) *OrderValues {
	o.DispatchStrategy = &strategy
	return o
}

func (o *OrderValues) SetNextStrategy(strategy string) *OrderValues {
	o.NextStrategy = &strategy
	return o
}

func (o *OrderValues) SetMetadata(metadata map[string]interface{}) *OrderValues {
	o.Metadata = metadata
	return o
}

func (o *OrderValues) SetNotes(notes string) *OrderValues {
	o.Notes = &notes
	return o
}

func (o *OrderValues) SetUpdatedAt(updatedAt int64) *OrderValues {
	o.UpdatedAt = updatedAt
	return o
}

// 业务方法
func (o *Order) IsActive() bool {
	status := o.GetStatus()
	return status != protocol.StatusCompleted && status != protocol.StatusCancelled
}

func (o *Order) CanBeCancelled() bool {
	status := o.GetStatus()
	return status == protocol.StatusRequested || status == protocol.StatusAccepted
}

func (o *Order) IsCompleted() bool {
	return o.GetStatus() == protocol.StatusCompleted
}

func (o *Order) IsCancelled() bool {
	return o.GetStatus() == protocol.StatusCancelled
}
func (o *Order) IsFinished() bool {
	return o.GetStatus() == protocol.StatusTripEnded
}

func (o *Order) IsScheduled() bool {
	return o.GetScheduleType() == protocol.ScheduleTypeScheduled
}

func (o *Order) IsSandbox() bool {
	return o.GetSandbox() > 0
}

// 状态变更方法
func (o *OrderValues) AcceptOrder(providerID string) *OrderValues {
	o.SetStatus(protocol.StatusAccepted)
	o.SetProviderID(providerID)
	now := utils.TimeNowMilli()
	o.AcceptedAt = &now
	return o
}

func (o *OrderValues) StartOrder() *OrderValues {
	o.SetStatus(protocol.StatusInProgress)
	now := utils.TimeNowMilli()
	o.StartedAt = &now
	return o
}

func (o *OrderValues) FinishOrder() *OrderValues {
	o.SetStatus(protocol.StatusTripEnded).
		SetPaymentStatus(protocol.StatusPending).
		SetEndedAt(utils.TimeNowMilli())
	return o
}

func (o *OrderValues) CompleteOrder() *OrderValues {
	o.SetStatus(protocol.StatusCompleted).
		SetCompletedAt(utils.TimeNowMilli())
	return o
}

func (o *OrderValues) CancelOrder(cancelledBy, reason string) *OrderValues {
	o.SetStatus(protocol.StatusCancelled)
	o.CancelledBy = &cancelledBy
	o.CancelReason = &reason
	now := utils.TimeNowMilli()
	o.CancelledAt = &now
	return o
}

// 计算最终支付金额
func (o *OrderValues) CalculatePaymentAmount() {
	// 最终支付金额 = 优惠后总金额 + 平台费用
	payment := o.GetDiscountedAmount()

	// 添加平台费用
	if o.PlatformFee != nil {
		payment = payment.Add(*o.PlatformFee)
	}

	// 确保金额不为负
	if payment.LessThan(decimal.Zero) {
		payment = decimal.Zero
	}

	o.PaymentAmount = &payment
}

// OrderV2 实例方法
func (o *Order) GetID() string {
	return o.OrderID
}

func (o *Order) GetVersion() int64 {
	if o.Version == nil {
		return 1
	}
	return *o.Version
}

func (o *Order) GetOrderID() string {
	return o.OrderID
}

// GetDetails 获取订单详情
func (o *Order) GetDetails() *OrderDetail {
	return o.Details
}

func (o *Order) GetAcceptedAt() int64 {
	if o.AcceptedAt == nil {
		return 0
	}
	return *o.AcceptedAt
}

func (o *Order) GetDistance() float64 {
	// 这里应该从RideOrder表或缓存中获取距离
	// 暂时返回0，实际实现需要查询相关表
	return 0
}

func (o *Order) IncrementVersion() {
	version := o.GetVersion() + 1
	o.Version = &version
}

func (o *Order) SetVersion(version int64) {
	o.Version = &version
}

func (o *OrderValues) GetVersion() int64 {
	if o.Version == nil {
		return 1
	}
	return *o.Version
}

func (o *OrderValues) SetVersion(version int64) *OrderValues {
	o.Version = &version
	return o
}

func GetOrderByID(orderID string) *Order {
	if orderID == "" {
		return nil
	}

	// 1. 尝试从缓存获取
	cacheKey := fmt.Sprintf("order:%s", orderID)
	var cachedOrder Order
	err := GetObjectCache(cacheKey, &cachedOrder)
	if err == nil {
		return &cachedOrder
	}

	// 2. 缓存未命中，从数据库查询
	var order Order
	if err := DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return nil
	}

	// 3. 将查询结果存入缓存
	_ = SetObjectCache(cacheKey, &order, 30*time.Minute)

	return &order
}

func GetOrderListByID(orderIDs []string) []*Order {
	var orders []*Order
	if err := DB.Where("order_id IN ?", orderIDs).Find(&orders).Error; err != nil {
		return nil
	}
	return orders
}

func CountProcessingOrdersByUserID(userID string) int64 {
	var count int64
	query := DB.Model(&Order{}).Where("user_id=?", userID)
	query = query.Where("status NOT IN ?", []string{protocol.StatusTripEnded, protocol.StatusCompleted, protocol.StatusCancelled})
	if err := query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func CountUnCompletedOrdersByUserID(userID string) int64 {
	var count int64
	query := DB.Model(&Order{}).Where("user_id=?", userID)
	query = query.Where("(status NOT IN ?) or (status=? AND payment_status != ?)", []string{protocol.StatusCompleted, protocol.StatusCancelled}, protocol.StatusCompleted, protocol.StatusSuccess)
	if err := query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func UpdateOrder(tx *gorm.DB, order *Order, values *OrderValues) error {
	// 如果没有提供事务，使用默认DB
	if tx == nil {
		tx = DB
	}
	values.SetVersion(order.GetVersion() + 1)
	rs := tx.Model(order).UpdateColumns(values)
	if rs.Error != nil {
		return rs.Error
	}
	// 检查是否有行被更新（乐观锁检查）
	if rs.RowsAffected == 0 {
		return errors.New("order version conflict or order not found")
	}
	order.SetValues(values)

	// 清除缓存，确保下次查询获取最新数据
	ClearOrderCache(order.OrderID)

	return nil
}

// ToOrderDetail 将Order转换为基础的OrderDetail（不包含业务特定字段）
func (o *Order) ToOrderDetail() *protocol.OrderDetail {
	// 如果Details字段存在且不为空，使用其Protocol()方法转换
	if o.Details != nil {
		return o.Details.Protocol()
	}

	// 否则返回基础的OrderDetail
	detail := &protocol.OrderDetail{
		OrderType: o.GetOrderType(),
	}

	return detail
}

// Protocol 将Order转换为protocol.Order
func (o *Order) Protocol() *protocol.Order {
	// 将Decimal类型转换为float64
	originalAmount, _ := o.GetOriginalAmount().Float64()
	discountedAmount, _ := o.GetDiscountedAmount().Float64()
	paymentAmount, _ := o.GetPaymentAmount().Float64()
	totalDiscountAmount, _ := o.GetTotalDiscountAmount().Float64()
	platformFee, _ := o.GetPlatformFee().Float64()
	promoDiscount, _ := o.GetPromoDiscount().Float64()
	cancellationFee, _ := o.GetCancellationFee().Float64()

	info := &protocol.Order{
		OrderID:             o.OrderID,
		OrderType:           o.GetOrderType(),
		UserID:              o.GetUserID(),
		ProviderID:          o.GetProviderID(),
		Status:              o.GetStatus(),
		PaymentStatus:       o.GetPaymentStatus(),
		OriginalAmount:      originalAmount,      // 预估原始金额
		DiscountedAmount:    discountedAmount,    // 优惠后总金额
		PaymentAmount:       paymentAmount,       // 最终付款金额
		TotalDiscountAmount: totalDiscountAmount, // 总优惠金额
		PlatformFee:         platformFee,
		PromoDiscount:       promoDiscount,
		CancellationFee:     cancellationFee,
		Currency:            o.GetCurrency(),
		PaymentMethod:       o.GetPaymentMethod(),
		PaymentID:           o.GetPaymentID(),
		ChannelPaymentID:    o.GetChannelPaymentID(),
		PaymentResult:       o.GetPaymentResult(),
		PaymentRedirectURL:  o.GetPaymentRedirectURL(),
		PromoCodes:          o.GetPromoCodes(),
		UserPromotionIDs:    o.GetUserPromotionIDs(),
		ScheduleType:        o.GetScheduleType(),
		ScheduledAt:         o.GetScheduledAt(),
		CreatedAt:           o.CreatedAt,
		AcceptedAt:          o.GetAcceptedAt(),
		StartedAt:           o.GetStartedAt(),
		EndedAt:             o.GetEndedAt(),
		CompletedAt:         o.GetCompletedAt(),
		CancelledAt:         o.GetCancelledAt(),
		CancelledBy:         o.GetCancelledBy(),
		CancelReason:        o.GetCancelReason(),
		DispatchStatus:      o.GetDispatchStatus(),
		CurrentRound:        o.GetCurrentRound(),
		MaxRounds:           o.GetMaxRounds(),
		DispatchStartedAt:   o.GetDispatchStartedAt(),
		LastDispatchedAt:    o.GetLastDispatchedAt(),
		AutoDispatchEnabled: o.GetAutoDispatchEnabled(),
		DispatchStrategy:    o.GetDispatchStrategy(),
		NextStrategy:        o.GetNextStrategy(),
		Notes:               o.GetNotes(),
		Metadata:            o.Metadata,
		Version:             o.GetVersion(),
	}
	if o.Details != nil {
		info.Details = o.Details.Protocol()
	}

	return info
}

// ============================================================================
// 订单缓存管理功能
// ============================================================================

const (
	OrderCacheKeyPrefix  = "order:"
	OrderCacheExpiration = 30 * time.Minute
)

// ClearOrderCache 清除订单缓存 - 统一的缓存删除函数
func ClearOrderCache(orderID string) {
	if orderID == "" {
		return
	}
	cacheKey := fmt.Sprintf("%s%s", OrderCacheKeyPrefix, orderID)
	_ = Delete(cacheKey)
}

// ClearOrderCacheBatch 批量清除订单缓存
func ClearOrderCacheBatch(orderIDs []string) {
	if len(orderIDs) == 0 {
		return
	}

	keys := make([]string, len(orderIDs))
	for i, orderID := range orderIDs {
		keys[i] = fmt.Sprintf("%s%s", OrderCacheKeyPrefix, orderID)
	}
	_ = DelCache(keys...)
}

// RefreshOrderCache 刷新订单缓存
func RefreshOrderCache(orderID string) {
	if orderID == "" {
		return
	}

	// 清除旧缓存
	ClearOrderCache(orderID)

	// 重新查询并缓存
	_ = GetOrderByID(orderID)
}

// CreateOrderWithCache 创建订单并设置缓存
func CreateOrderWithCache(order *Order) error {
	if order == nil || order.OrderID == "" {
		return errors.New("invalid order")
	}

	// 1. 创建订单
	err := DB.Create(order).Error
	if err != nil {
		return err
	}

	// 2. 设置缓存
	cacheKey := fmt.Sprintf("%s%s", OrderCacheKeyPrefix, order.OrderID)
	_ = SetObjectCache(cacheKey, order, OrderCacheExpiration)

	return nil
}
