package protocol

// 支付渠道相关常量
const (
	// 支付渠道
	PaymentChannelStripe  = "stripe"  // Stripe支付
	PaymentChannelCard    = "card"    // 卡支付
	PaymentChannelKPay    = "kpay"    // KPay支付
	PaymentChannelCash    = "cash"    // 现金支付
	PaymentChannelSandbox = "sandbox" // 沙盒支付
	PaymentChannelMoMo    = "momo"    // MTN MoMo Direct API

	// 支付方式
	PayMethodCard         = "card"          // 信用卡/借记卡
	PayMethodApplePay     = "apple_pay"     // Apple Pay
	PayMethodGooglePlay   = "google_play"   // Google Play
	PayMethodAlipay       = "alipay"        // 支付宝
	PayMethodWechatPay    = "wechat_pay"    // 微信支付
	PayMethodBankTransfer = "bank_transfer" // 银行转账
	PayMethodCrypto       = "crypto"        // 加密货币

	// 交易类型
	PaymentTypePayment = "payment" // 支付
	PaymentTypeRefund  = "refund"  // 退款
)

// PaymentMethodsRequest 获取支付方式列表请求
type PaymentMethodsRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// ChannelResult 支付渠道处理结果
type ChannelResult struct {
	PaymentID        string  `json:"payment_id"`              // 支付订单ID
	Status           string  `json:"status"`                  // 支付状态
	ChannelStatus    string  `json:"channel_status"`          // 渠道状态
	ResCode          string  `json:"res_code"`                // 响应码
	ResMsg           string  `json:"res_msg"`                 // 响应消息
	OrderType        string  `json:"type"`                    // 交易类型
	ChannelCode      string  `json:"channel_code"`            // 渠道代码
	ChannelPaymentID string  `json:"channel_payment_id"`      // 渠道订单ID
	RedirectURL      string  `json:"redirect_url,omitempty"`  // 重定向URL
	CallbackData     string  `json:"callback_data,omitempty"` // 回调数据
	Metadata         MapData `json:"metadata,omitempty"`      // 扩展元数据 (e.g., Stripe client_secret)
}

type CancelPaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  string `json:"user_id,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// PaymentListRequest 支付列表请求
type PaymentListRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	OrderID   string `json:"order_id,omitempty"`
	OrderType string `json:"order_type,omitempty"`
	Status    string `json:"status,omitempty"`
	StartDate *int64 `json:"start_date,omitempty"`
	EndDate   *int64 `json:"end_date,omitempty"`
	Page      int    `json:"page" binding:"min=1"`
	Limit     int    `json:"limit" binding:"min=1,max=100"`
}

// PaymentListResponse 支付列表响应
type PaymentListResponse struct {
	Payments []*Payment `json:"payments"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	Limit    int        `json:"limit"`
}

// PaymentDetailRequest 支付详情请求
type PaymentDetailRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

// PaymentStatusRequest 支付状态查询请求
type PaymentStatusRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	OrderID   string `json:"order_id,omitempty"`
}

// RefundRequest 退款请求
type RefundRequest struct {
	PaymentID    string  `json:"payment_id" binding:"required"`
	RefundAmount float64 `json:"refund_amount,omitempty"` // 空则全额退款
	RefundReason string  `json:"refund_reason" binding:"required"`
	UserID       string  `json:"user_id" binding:"required"`
}

// RefundResponse 退款响应
type RefundResponse struct {
	RefundID     string  `json:"refund_id"`
	Status       string  `json:"status"`
	RefundAmount float64 `json:"refund_amount"`
	Message      string  `json:"message,omitempty"`
}

// Payment 支付协议层结构
type Payment struct {
	// 基础信息
	ID        int64  `json:"id"`
	PaymentID string `json:"payment_id"`

	// 关联信息
	OrderID    string `json:"order_id"`
	OriOrderID string `json:"ori_order_id,omitempty"`
	OrderType  string `json:"order_type"`
	OrderSku   string `json:"order_sku,omitempty"`
	UserID     string `json:"user_id"`

	// 支付信息
	PaymentMethod string `json:"payment_method"`

	// 客户信息
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	AccountNo   string `json:"account_no,omitempty"`

	// 状态信息
	Status        string `json:"status"`
	ChannelStatus string `json:"channel_status,omitempty"`
	ResCode       string `json:"res_code,omitempty"`
	ResMsg        string `json:"res_msg,omitempty"`

	// 金额信息
	Currency          string  `json:"currency"`
	Amount            float64 `json:"amount"`
	UsdAmount         float64 `json:"usd_amount"`
	UsdRate           float64 `json:"usd_rate"`
	ReceivedCcy       string  `json:"received_ccy,omitempty"`
	ReceivedAmount    float64 `json:"received_amount"`
	ReceivedUSDAmount float64 `json:"received_usd_amount"`

	// 支付渠道信息
	CardID            string  `json:"card_id,omitempty"`
	ChannelCode       string  `json:"channel_code,omitempty"`
	ChannelPaymentID  string  `json:"channel_payment_id,omitempty"`
	ChannelAccountID  string  `json:"channel_account_id,omitempty"`
	ChannelFeeCcy     string  `json:"channel_fee_ccy,omitempty"`
	ChannelFeeAmount  float64 `json:"channel_fee_amount"`
	ChannelPaidCcy    string  `json:"channel_paid_ccy,omitempty"`
	ChannelPaidAmount float64 `json:"channel_paid_amount"`

	RedirectURL string `json:"redirect_url,omitempty"` // 支付重定向URL
	ReturnURL   string `json:"return_url,omitempty"`   // 支付结果返回页地址
	CallbackUrl string `json:"callback_url,omitempty"` // 支付回调URL

	// 重试信息
	RetryCount int `json:"retry_count"`

	// 退款信息
	RefundAmount float64 `json:"refund_amount"`
	RefundReason string  `json:"refund_reason,omitempty"`
	RefundedAt   int64   `json:"refunded_at,omitempty"`

	// 描述和备注
	Description string `json:"description,omitempty"`
	Remark      string `json:"remark,omitempty"`

	// 扩展信息
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 扩展元数据

	// 时间戳信息
	ExpiredAt   int64 `json:"expired_at,omitempty"`
	CompletedAt int64 `json:"completed_at,omitempty"`
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
	Version     int64 `json:"version"`
}
