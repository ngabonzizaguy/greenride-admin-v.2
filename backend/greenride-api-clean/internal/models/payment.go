package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Payment 支付表 - 基于最新设计文档
type Payment struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PaymentID string `json:"payment_id" gorm:"column:payment_id;type:varchar(64);uniqueIndex:idx_payment_id"`
	Salt      string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PaymentValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PaymentValues struct {
	// 关联信息
	OrderID    *string `json:"order_id" gorm:"column:order_id;type:varchar(64);index:idx_order_id"`
	OriOrderID *string `json:"ori_order_id" gorm:"column:ori_order_id;type:varchar(64)"`
	OrderType  *string `json:"order_type" gorm:"column:order_type;type:varchar(32)"`
	OrderSku   *string `json:"order_sku" gorm:"column:order_sku;type:varchar(255)"`
	UserID     *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index:idx_user_id"`

	// 客户信息
	Phone       *string `json:"phone" gorm:"column:phone;type:varchar(32)"`
	Email       *string `json:"email" gorm:"column:email;type:varchar(255)"`
	AccountName *string `json:"account_name" gorm:"column:account_name;type:varchar(255)"`
	AccountNo   *string `json:"account_no" gorm:"column:account_no;type:varchar(64)"`

	// 支付信息
	PaymentMethod *string `json:"payment_method" gorm:"column:payment_method;type:varchar(32);index:idx_payment_method"` // visa, mastercard, paypal, stripe, alipay, cash, wallet, bank_transfer

	// 状态信息
	Status        *string `json:"status" gorm:"column:status;type:varchar(32);index:idx_status"`
	ChannelStatus *string `json:"channel_status" gorm:"column:channel_status;type:varchar(32)"`
	ResCode       *string `json:"res_code" gorm:"column:res_code;type:varchar(32)"`
	ResMsg        *string `json:"res_msg" gorm:"column:res_msg;type:varchar(255)"`
	Reason        *string `json:"reason" gorm:"column:reason;type:varchar(255)"` // Deprecated: use ResMsg instead

	// 金额信息
	Currency          *string          `json:"currency" gorm:"column:currency;type:varchar(10);default:'USD'"`
	Amount            *decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(20,6)"`
	UsdAmount         *decimal.Decimal `json:"usd_amount" gorm:"column:usd_amount;type:decimal(20,6)"`
	UsdRate           *decimal.Decimal `json:"usd_rate" gorm:"column:usd_rate;type:decimal(20,6)"`
	ReceivedCcy       *string          `json:"received_ccy" gorm:"column:received_ccy;type:varchar(10)"`
	ReceivedAmount    *decimal.Decimal `json:"received_amount" gorm:"column:received_amount;type:decimal(20,6)"`
	ReceivedUSDAmount *decimal.Decimal `json:"received_usd_amount" gorm:"column:received_usd_amount;type:decimal(20,6)"`

	// 支付渠道信息
	CardID            *string          `json:"card_id" gorm:"column:card_id;type:varchar(64)"`
	ChannelCode       *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32)"`
	ChannelPaymentID  *string          `json:"channel_payment_id" gorm:"column:channel_payment_id;type:varchar(128)"`
	ChannelAccountID  *string          `json:"channel_account_id" gorm:"column:channel_account_id;type:varchar(64)"`
	ChannelFeeCcy     *string          `json:"channel_fee_ccy" gorm:"column:channel_fee_ccy;type:varchar(10)"`
	ChannelFeeAmount  *decimal.Decimal `json:"channel_fee_amount" gorm:"column:channel_fee_amount;type:decimal(20,6)"`
	ChannelPaidCcy    *string          `json:"channel_paid_ccy" gorm:"column:channel_paid_ccy;type:varchar(10)"`
	ChannelPaidAmount *decimal.Decimal `json:"channel_paid_amount" gorm:"column:channel_paid_amount;type:decimal(20,6)"`

	RedirectURL *string `json:"redirect_url" gorm:"column:redirect_url;type:varchar(255)"` // 支付重定向URL
	ReturnURL   *string `json:"return_url" gorm:"column:return_url;type:varchar(512)"`     // 支付结果返回页地址
	CallbackUrl *string `json:"callback_url" gorm:"column:callback_url;type:varchar(255)"` // 支付回调URL

	// 重试信息
	RetryCount *int `json:"retry_count" gorm:"column:retry_count;default:0"`

	// 退款信息
	RefundAmount *decimal.Decimal `json:"refund_amount" gorm:"column:refund_amount;type:decimal(20,6)"`
	RefundReason *string          `json:"refund_reason" gorm:"column:refund_reason;type:varchar(255)"`
	RefundedAt   *int64           `json:"refunded_at" gorm:"column:refunded_at"`

	// 描述和备注
	Description *string `json:"description" gorm:"column:description;type:varchar(500)"`
	Remark      *string `json:"remark" gorm:"column:remark;type:text"`

	// 扩展信息
	Metadata protocol.MapData `json:"metadata" gorm:"column:metadata;type:json;serializer:json"` // JSON格式的额外信息

	// 时间戳和版本
	ExpiredAt   *int64 `json:"expired_at" gorm:"column:expired_at"`
	CompletedAt *int64 `json:"completed_at" gorm:"column:completed_at"`
	Version     *int64 `json:"version" gorm:"column:version;default:1"`
	UpdatedAt   int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Payment) TableName() string {
	return "t_payments"
}

// Getter 方法
// 状态相关
func (p *PaymentValues) GetStatus() string {
	if p.Status == nil || *p.Status == "" {
		return ""
	}
	return *p.Status
}

func (p *PaymentValues) GetChannelStatus() string {
	if p.ChannelStatus == nil {
		return ""
	}
	return *p.ChannelStatus
}

func (p *PaymentValues) GetResCode() string {
	if p.ResCode == nil {
		return ""
	}
	return *p.ResCode
}

func (p *PaymentValues) GetResMsg() string {
	if p.ResMsg == nil {
		return ""
	}
	return *p.ResMsg
}

// 订单和业务相关
func (p *PaymentValues) GetOrderID() string {
	if p.OrderID == nil {
		return ""
	}
	return *p.OrderID
}

func (p *PaymentValues) GetOriOrderID() string {
	if p.OriOrderID == nil {
		return ""
	}
	return *p.OriOrderID
}

func (p *PaymentValues) GetOrderType() string {
	if p.OrderType == nil {
		return ""
	}
	return *p.OrderType
}

func (p *PaymentValues) GetOrderSku() string {
	if p.OrderSku == nil {
		return ""
	}
	return *p.OrderSku
}

func (p *PaymentValues) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return *p.UserID
}

func (p *PaymentValues) GetPaymentMethod() string {
	if p.PaymentMethod == nil {
		return ""
	}
	return *p.PaymentMethod
}

// 客户信息相关
func (p *PaymentValues) GetPhone() string {
	if p.Phone == nil {
		return ""
	}
	return *p.Phone
}

func (p *PaymentValues) GetEmail() string {
	if p.Email == nil {
		return ""
	}
	return *p.Email
}

func (p *PaymentValues) GetAccountName() string {
	if p.AccountName == nil {
		return ""
	}
	return *p.AccountName
}

func (p *PaymentValues) GetAccountNo() string {
	if p.AccountNo == nil {
		return ""
	}
	return *p.AccountNo
}

// 金额相关
func (p *PaymentValues) GetCurrency() string {
	if p.Currency == nil {
		return ""
	}
	return *p.Currency
}

func (p *PaymentValues) GetAmount() decimal.Decimal {
	if p.Amount == nil {
		return decimal.Zero
	}
	return *p.Amount
}

func (p *PaymentValues) GetUsdAmount() decimal.Decimal {
	if p.UsdAmount == nil {
		return decimal.Zero
	}
	return *p.UsdAmount
}

func (p *PaymentValues) GetUsdRate() decimal.Decimal {
	if p.UsdRate == nil {
		return decimal.NewFromInt(1)
	}
	return *p.UsdRate
}

func (p *PaymentValues) GetReceivedCcy() string {
	if p.ReceivedCcy == nil {
		return ""
	}
	return *p.ReceivedCcy
}

func (p *PaymentValues) GetReceivedAmount() decimal.Decimal {
	if p.ReceivedAmount == nil {
		return decimal.Zero
	}
	return *p.ReceivedAmount
}

func (p *PaymentValues) GetReceivedUSDAmount() decimal.Decimal {
	if p.ReceivedUSDAmount == nil {
		return decimal.Zero
	}
	return *p.ReceivedUSDAmount
}

// 支付渠道相关
func (p *PaymentValues) GetCardID() string {
	if p.CardID == nil {
		return ""
	}
	return *p.CardID
}

func (p *PaymentValues) GetChannelCode() string {
	if p.ChannelCode == nil {
		return ""
	}
	return *p.ChannelCode
}

func (p *PaymentValues) GetChannelPaymentID() string {
	if p.ChannelPaymentID == nil {
		return ""
	}
	return *p.ChannelPaymentID
}

func (p *PaymentValues) GetChannelAccountID() string {
	if p.ChannelAccountID == nil {
		return ""
	}
	return *p.ChannelAccountID
}

func (p *PaymentValues) GetChannelFeeCcy() string {
	if p.ChannelFeeCcy == nil {
		return ""
	}
	return *p.ChannelFeeCcy
}

func (p *PaymentValues) GetChannelFeeAmount() decimal.Decimal {
	if p.ChannelFeeAmount == nil {
		return decimal.Zero
	}
	return *p.ChannelFeeAmount
}

func (p *PaymentValues) GetChannelPaidCcy() string {
	if p.ChannelPaidCcy == nil {
		return ""
	}
	return *p.ChannelPaidCcy
}

func (p *PaymentValues) GetChannelPaidAmount() decimal.Decimal {
	if p.ChannelPaidAmount == nil {
		return decimal.Zero
	}
	return *p.ChannelPaidAmount
}

func (p *PaymentValues) GetRedirectURL() string {
	if p.RedirectURL == nil {
		return ""
	}
	return *p.RedirectURL
}

func (p *PaymentValues) GetReturnURL() string {
	if p.ReturnURL == nil {
		return ""
	}
	return *p.ReturnURL
}

func (p *PaymentValues) GetCallbackUrl() string {
	if p.CallbackUrl == nil {
		return ""
	}
	return *p.CallbackUrl
}

// 退款相关
func (p *PaymentValues) GetRefundAmount() decimal.Decimal {
	if p.RefundAmount == nil {
		return decimal.Zero
	}
	return *p.RefundAmount
}

func (p *PaymentValues) GetRefundReason() string {
	if p.RefundReason == nil {
		return ""
	}
	return *p.RefundReason
}

func (p *PaymentValues) GetRefundedAt() int64 {
	if p.RefundedAt == nil {
		return 0
	}
	return *p.RefundedAt
}

// 其他信息
func (p *PaymentValues) GetDescription() string {
	if p.Description == nil {
		return ""
	}
	return *p.Description
}

func (p *PaymentValues) GetRemark() string {
	if p.Remark == nil {
		return ""
	}
	return *p.Remark
}

func (p *PaymentValues) GetMetadata() protocol.MapData {
	if p.Metadata == nil {
		return protocol.MapData{}
	}
	return p.Metadata
}

func (p *PaymentValues) GetRetryCount() int {
	if p.RetryCount == nil {
		return 0
	}
	return *p.RetryCount
}

func (p *PaymentValues) GetExpiredAt() int64 {
	if p.ExpiredAt == nil {
		return 0
	}
	return *p.ExpiredAt
}

func (p *PaymentValues) GetCompletedAt() int64 {
	if p.CompletedAt == nil {
		return 0
	}
	return *p.CompletedAt
}

func (p *PaymentValues) GetVersion() int64 {
	if p.Version == nil {
		return 1
	}
	return *p.Version
}

// Setter 方法
// 状态相关
func (p *PaymentValues) SetStatus(status string) *PaymentValues {
	p.Status = &status
	return p
}

func (p *PaymentValues) SetChannelStatus(status string) *PaymentValues {
	p.ChannelStatus = &status
	return p
}

func (p *PaymentValues) SetResCode(code string) *PaymentValues {
	p.ResCode = &code
	return p
}

func (p *PaymentValues) SetResMsg(msg string) *PaymentValues {
	p.ResMsg = &msg
	return p
}

// 订单和业务相关
func (p *PaymentValues) SetOrderID(orderID string) *PaymentValues {
	p.OrderID = &orderID
	return p
}

func (p *PaymentValues) SetOriOrderID(oriOrderID string) *PaymentValues {
	p.OriOrderID = &oriOrderID
	return p
}

func (p *PaymentValues) SetOrderType(orderType string) *PaymentValues {
	p.OrderType = &orderType
	return p
}

func (p *PaymentValues) SetOrderSku(orderSku string) *PaymentValues {
	p.OrderSku = &orderSku
	return p
}

func (p *PaymentValues) SetUserID(userID string) *PaymentValues {
	p.UserID = &userID
	return p
}

func (p *PaymentValues) SetPaymentMethod(method string) *PaymentValues {
	p.PaymentMethod = &method
	return p
}

// 客户信息相关
func (p *PaymentValues) SetPhone(phone string) *PaymentValues {
	p.Phone = &phone
	return p
}

func (p *PaymentValues) SetEmail(email string) *PaymentValues {
	p.Email = &email
	return p
}

func (p *PaymentValues) SetAccountName(name string) *PaymentValues {
	p.AccountName = &name
	return p
}

func (p *PaymentValues) SetAccountNo(accountNo string) *PaymentValues {
	p.AccountNo = &accountNo
	return p
}

// 金额相关
func (p *PaymentValues) SetCurrency(currency string) *PaymentValues {
	p.Currency = &currency
	return p
}

func (p *PaymentValues) SetAmount(amount decimal.Decimal) *PaymentValues {
	p.Amount = &amount
	return p
}

func (p *PaymentValues) SetUsdAmount(amount decimal.Decimal) *PaymentValues {
	p.UsdAmount = &amount
	return p
}

func (p *PaymentValues) SetUsdRate(rate decimal.Decimal) *PaymentValues {
	p.UsdRate = &rate
	return p
}

func (p *PaymentValues) SetReceivedCcy(ccy string) *PaymentValues {
	p.ReceivedCcy = &ccy
	return p
}

func (p *PaymentValues) SetReceivedAmount(amount decimal.Decimal) *PaymentValues {
	p.ReceivedAmount = &amount
	return p
}

func (p *PaymentValues) SetReceivedUSDAmount(amount decimal.Decimal) *PaymentValues {
	p.ReceivedUSDAmount = &amount
	return p
}

// 支付渠道相关
func (p *PaymentValues) SetCardID(cardID string) *PaymentValues {
	p.CardID = &cardID
	return p
}

func (p *PaymentValues) SetChannelCode(code string) *PaymentValues {
	p.ChannelCode = &code
	return p
}

func (p *PaymentValues) SetChannelPaymentID(paymentID string) *PaymentValues {
	p.ChannelPaymentID = &paymentID
	return p
}

func (p *PaymentValues) SetChannelAccountID(accountID string) *PaymentValues {
	p.ChannelAccountID = &accountID
	return p
}

func (p *PaymentValues) SetChannelFee(ccy string, amount decimal.Decimal) *PaymentValues {
	p.ChannelFeeCcy = &ccy
	p.ChannelFeeAmount = &amount
	return p
}

func (p *PaymentValues) SetChannelPaid(ccy string, amount decimal.Decimal) *PaymentValues {
	p.ChannelPaidCcy = &ccy
	p.ChannelPaidAmount = &amount
	return p
}

func (p *PaymentValues) SetRedirectURL(url string) *PaymentValues {
	p.RedirectURL = &url
	return p
}

func (p *PaymentValues) SetReturnURL(url string) *PaymentValues {
	p.ReturnURL = &url
	return p
}

func (p *PaymentValues) SetCallbackUrl(url string) *PaymentValues {
	p.CallbackUrl = &url
	return p
}

// 退款相关
func (p *PaymentValues) SetRefund(amount decimal.Decimal, reason string) *PaymentValues {
	p.RefundAmount = &amount
	p.RefundReason = &reason
	return p
}

func (p *PaymentValues) SetRefundedAt(time int64) *PaymentValues {
	p.RefundedAt = &time
	return p
}

// 其他信息
func (p *PaymentValues) SetDescription(desc string) *PaymentValues {
	p.Description = &desc
	return p
}

func (p *PaymentValues) SetRemark(remark string) *PaymentValues {
	p.Remark = &remark
	return p
}

func (p *PaymentValues) SetMetadata(metadata protocol.MapData) *PaymentValues {
	if metadata == nil {
		p.Metadata = protocol.MapData{}
	} else {
		p.Metadata = metadata
	}
	return p
}

func (p *PaymentValues) SetRetryCount(count int) *PaymentValues {
	p.RetryCount = &count
	return p
}

// 时间和版本
func (p *PaymentValues) SetExpiredAt(time int64) *PaymentValues {
	p.ExpiredAt = &time
	return p
}

func (p *PaymentValues) SetCompletedAt(time int64) *PaymentValues {
	p.CompletedAt = &time
	return p
}

func (p *PaymentValues) IncrementVersion() *PaymentValues {
	if p.Version == nil {
		one := int64(1)
		p.Version = &one
	} else {
		newVersion := *p.Version + 1
		p.Version = &newVersion
	}
	return p
}

// SetValues 更新PaymentValues中的非nil/非空值
func (p *PaymentValues) SetValues(values *PaymentValues) {
	if values == nil {
		return
	}

	// 订单和业务相关
	if values.OrderID != nil {
		p.OrderID = values.OrderID
	}
	if values.OriOrderID != nil {
		p.OriOrderID = values.OriOrderID
	}
	if values.OrderType != nil {
		p.OrderType = values.OrderType
	}
	if values.OrderSku != nil {
		p.OrderSku = values.OrderSku
	}
	if values.UserID != nil {
		p.UserID = values.UserID
	}
	if values.PaymentMethod != nil {
		p.PaymentMethod = values.PaymentMethod
	}

	// 客户信息相关
	if values.Phone != nil {
		p.Phone = values.Phone
	}
	if values.Email != nil {
		p.Email = values.Email
	}
	if values.AccountName != nil {
		p.AccountName = values.AccountName
	}
	if values.AccountNo != nil {
		p.AccountNo = values.AccountNo
	}

	// 状态相关
	if values.Status != nil {
		p.Status = values.Status
	}
	if values.ChannelStatus != nil {
		p.ChannelStatus = values.ChannelStatus
	}
	if values.ResCode != nil {
		p.ResCode = values.ResCode
	}
	if values.ResMsg != nil {
		p.ResMsg = values.ResMsg
	}
	if values.Reason != nil {
		p.Reason = values.Reason
	}

	// 金额相关
	if values.Currency != nil {
		p.Currency = values.Currency
	}
	if values.Amount != nil {
		p.Amount = values.Amount
	}
	if values.UsdAmount != nil {
		p.UsdAmount = values.UsdAmount
	}
	if values.UsdRate != nil {
		p.UsdRate = values.UsdRate
	}
	if values.ReceivedCcy != nil {
		p.ReceivedCcy = values.ReceivedCcy
	}
	if values.ReceivedAmount != nil {
		p.ReceivedAmount = values.ReceivedAmount
	}
	if values.ReceivedUSDAmount != nil {
		p.ReceivedUSDAmount = values.ReceivedUSDAmount
	}

	// 支付渠道相关
	if values.CardID != nil {
		p.CardID = values.CardID
	}
	if values.ChannelCode != nil {
		p.ChannelCode = values.ChannelCode
	}
	if values.ChannelPaymentID != nil {
		p.ChannelPaymentID = values.ChannelPaymentID
	}
	if values.ChannelAccountID != nil {
		p.ChannelAccountID = values.ChannelAccountID
	}
	if values.ChannelFeeCcy != nil {
		p.ChannelFeeCcy = values.ChannelFeeCcy
	}
	if values.ChannelFeeAmount != nil {
		p.ChannelFeeAmount = values.ChannelFeeAmount
	}
	if values.ChannelPaidCcy != nil {
		p.ChannelPaidCcy = values.ChannelPaidCcy
	}
	if values.ChannelPaidAmount != nil {
		p.ChannelPaidAmount = values.ChannelPaidAmount
	}

	if values.RedirectURL != nil {
		p.RedirectURL = values.RedirectURL
	}
	if values.ReturnURL != nil {
		p.ReturnURL = values.ReturnURL
	}
	if values.CallbackUrl != nil {
		p.CallbackUrl = values.CallbackUrl
	}

	// 退款相关
	if values.RefundAmount != nil {
		p.RefundAmount = values.RefundAmount
	}
	if values.RefundReason != nil {
		p.RefundReason = values.RefundReason
	}
	if values.RefundedAt != nil {
		p.RefundedAt = values.RefundedAt
	}

	// 其他信息
	if values.Description != nil {
		p.Description = values.Description
	}
	if values.Remark != nil {
		p.Remark = values.Remark
	}
	if len(values.Metadata) > 0 {
		p.Metadata = values.Metadata
	}
	if values.RetryCount != nil {
		p.RetryCount = values.RetryCount
	}

	// 时间和版本
	if values.ExpiredAt != nil && *values.ExpiredAt > 0 {
		p.ExpiredAt = values.ExpiredAt
	}
	if values.CompletedAt != nil && *values.CompletedAt > 0 {
		p.CompletedAt = values.CompletedAt
	}
	if values.Version != nil && *values.Version > 0 {
		p.Version = values.Version
	}
	if values.UpdatedAt > 0 {
		p.UpdatedAt = values.UpdatedAt
	}
}

// NewPayment 创建新的支付对象
func NewPayment() *Payment {
	zero := decimal.NewFromInt(0)
	defaultCurrency := protocol.DefaultCurrency
	defaultStatus := protocol.StatusPending
	usdRate := decimal.NewFromInt(1)
	version := int64(1)

	return &Payment{
		PaymentID: utils.GeneratePaymentID(),
		Salt:      utils.GenerateSalt(),
		PaymentValues: &PaymentValues{
			Status:     &defaultStatus,
			Currency:   &defaultCurrency,
			Amount:     &zero,
			UsdAmount:  &zero,
			UsdRate:    &usdRate,
			RetryCount: utils.IntPtr(0),
			Version:    &version,
			Metadata:   protocol.MapData{},
		},
	}
}

func GetNotFailedPaymentByOrderID(orderID string) *Payment {
	var existingPayment Payment
	if err := DB.Where("order_id = ? AND status NOT IN (?)", orderID, []string{protocol.StatusFailed, protocol.StatusCancelled}).First(&existingPayment).Error; err != nil {
		return nil
	}
	return &existingPayment
}

func GetLastPaymentByOrderID(orderID string) *Payment {
	var existingPayment Payment
	if err := DB.Where("order_id = ?", orderID).Order("created_at DESC").First(&existingPayment).Error; err != nil {
		return nil
	}
	return &existingPayment
}

// GetPaymentByID 根据支付ID获取支付记录
func GetPaymentByID(paymentID string) *Payment {
	var payment Payment
	if err := DB.Where("payment_id = ?", paymentID).First(&payment).Error; err != nil {
		return nil
	}
	return &payment
}

func UpdatePaymentValues(tx *gorm.DB, payment *Payment, values *PaymentValues) error {
	defer func() {
		payment.SetValues(values)
	}()
	if err := tx.Model(payment).UpdateColumns(values).Error; err != nil {
		return err
	}
	return nil
}

// Protocol 将模型转换为协议层结构
func (p *Payment) Protocol() *protocol.Payment {
	if p == nil {
		return nil
	}

	// 获取 metadata，直接使用 MapData
	metadata := p.GetMetadata()
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &protocol.Payment{
		// 基础信息
		ID:        p.ID,
		PaymentID: p.PaymentID,

		// 关联信息
		OrderID:    p.GetOrderID(),
		OriOrderID: p.GetOriOrderID(),
		OrderType:  p.GetOrderType(),
		OrderSku:   p.GetOrderSku(),
		UserID:     p.GetUserID(),

		// 支付信息
		PaymentMethod: p.GetPaymentMethod(),

		// 客户信息
		Phone:       p.GetPhone(),
		Email:       p.GetEmail(),
		AccountName: p.GetAccountName(),
		AccountNo:   p.GetAccountNo(),

		// 状态信息
		Status:        p.GetStatus(),
		ChannelStatus: p.GetChannelStatus(),
		ResCode:       p.GetResCode(),
		ResMsg:        p.GetResMsg(),

		// 金额信息
		Currency:          p.GetCurrency(),
		Amount:            p.GetAmount().InexactFloat64(),
		UsdAmount:         p.GetUsdAmount().InexactFloat64(),
		UsdRate:           p.GetUsdRate().InexactFloat64(),
		ReceivedCcy:       p.GetReceivedCcy(),
		ReceivedAmount:    p.GetReceivedAmount().InexactFloat64(),
		ReceivedUSDAmount: p.GetReceivedUSDAmount().InexactFloat64(),

		// 支付渠道信息
		CardID:            p.GetCardID(),
		ChannelCode:       p.GetChannelCode(),
		ChannelPaymentID:  p.GetChannelPaymentID(),
		ChannelAccountID:  p.GetChannelAccountID(),
		ChannelFeeCcy:     p.GetChannelFeeCcy(),
		ChannelFeeAmount:  p.GetChannelFeeAmount().InexactFloat64(),
		ChannelPaidCcy:    p.GetChannelPaidCcy(),
		ChannelPaidAmount: p.GetChannelPaidAmount().InexactFloat64(),

		RedirectURL: p.GetRedirectURL(),
		ReturnURL:   p.GetReturnURL(),
		CallbackUrl: p.GetCallbackUrl(),

		// 重试信息
		RetryCount: p.GetRetryCount(),

		// 退款信息
		RefundAmount: p.GetRefundAmount().InexactFloat64(),
		RefundReason: p.GetRefundReason(),
		RefundedAt:   p.GetRefundedAt(),

		// 描述和备注
		Description: p.GetDescription(),
		Remark:      p.GetRemark(),

		// 扩展信息
		Metadata: metadata,

		// 时间戳信息
		ExpiredAt:   p.GetExpiredAt(),
		CompletedAt: p.GetCompletedAt(),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Version:     p.GetVersion(),
	}
}
