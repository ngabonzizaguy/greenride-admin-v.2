package protocol

// 支付方式常量
const (
	PaymentMethodWallet    = "wallet"
	PaymentMethodCash      = "cash"
	PaymentMethodCard      = "card"
	PaymentMethodVisa      = "visa"
	PaymentMethodMaster    = "master"
	PaymentMethodMomo      = "momo"
	PaymentMethodAirtel    = "airtel"
	PaymentMethodSpenn     = "spenn"
	PaymentMethodAmex      = "amex"
	PaymentMethodPaypal    = "paypal"
	PaymentMethodStripe    = "stripe"
	PaymentMethodAlipay    = "alipay"
	PaymentMethodWechatPay = "wechatpay"
)

// 支付错误结果码(ResCode)常量 - 用于标记请求渠道接口的失败情况
const (
	// 请求相关错误
	ResCodeRequestFailed    = "REQUEST_FAILED"    // 请求渠道接口失败
	ResCodeRequestTimeout   = "REQUEST_TIMEOUT"   // 请求超时
	ResCodeConnectionFailed = "CONNECTION_FAILED" // 连接失败
	ResCodeNetworkError     = "NETWORK_ERROR"     // 网络错误

	// 响应解析相关错误
	ResCodeResponseParseFailed = "RESPONSE_PARSE_FAILED" // 响应结果解析失败
	ResCodeInvalidResponse     = "INVALID_RESPONSE"      // 无效响应格式
	ResCodeMissingFields       = "MISSING_FIELDS"        // 响应缺少必要字段
	ResCodeUnexpectedFormat    = "UNEXPECTED_FORMAT"     // 响应格式不符合预期

	// 渠道相关错误
	ResCodeChannelError       = "CHANNEL_ERROR"        // 渠道错误
	ResCodeChannelUnavailable = "CHANNEL_UNAVAILABLE"  // 渠道不可用
	ResCodeChannelMaintenance = "CHANNEL_MAINTENANCE"  // 渠道维护中
	ResCodeChannelRateLimited = "CHANNEL_RATE_LIMITED" // 渠道限流

	// 配置相关错误
	ResCodeConfigError   = "CONFIG_ERROR"   // 配置错误
	ResCodeMissingConfig = "MISSING_CONFIG" // 缺少配置
	ResCodeInvalidConfig = "INVALID_CONFIG" // 无效配置

	// 认证相关错误
	ResCodeAuthFailed         = "AUTH_FAILED"         // 认证失败
	ResCodeInvalidCredentials = "INVALID_CREDENTIALS" // 无效凭证
	ResCodeTokenExpired       = "TOKEN_EXPIRED"       // 令牌过期

	// 业务逻辑错误
	ResCodeBusinessError            = "BUSINESS_ERROR"             // 业务逻辑错误
	ResCodeInvalidAmount            = "INVALID_AMOUNT"             // 无效金额
	ResCodeInvalidCurrency          = "INVALID_CURRENCY"           // 无效货币
	ResCodeInsufficientFunds        = "INSUFFICIENT_FUNDS"         // 余额不足
	ResCodeUnsupportedPaymentMethod = "UNSUPPORTED_PAYMENT_METHOD" // 不支持的支付方式

	// 支付成功结果码
	ResCodeSandboxSuccess = "SANDBOX_SUCCESS" // 沙盒支付成功
	ResCodeCashSuccess    = "CASH_SUCCESS"    // 现金支付成功

	// 支付失败结果码
	ResCodePaymentFailed = "PAYMENT_FAILED" // 支付失败

	// 系统错误
	ResCodeSystemError   = "SYSTEM_ERROR"   // 系统错误
	ResCodeInternalError = "INTERNAL_ERROR" // 内部错误
	ResCodeUnknownError  = "UNKNOWN_ERROR"  // 未知错误
)

// GetResCodeByStatusCode 根据HTTP状态码获取对应的错误结果码
func GetResCodeByStatusCode(statusCode int) string {
	switch statusCode {
	case 401:
		return ResCodeAuthFailed
	case 403:
		return ResCodeInvalidCredentials
	case 404:
		return ResCodeChannelUnavailable
	case 429:
		return ResCodeChannelRateLimited
	case 500, 502, 503, 504:
		return ResCodeChannelError
	default:
		return ResCodeChannelError
	}
}
