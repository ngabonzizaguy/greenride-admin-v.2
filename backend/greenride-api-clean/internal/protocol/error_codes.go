package protocol

// ErrorCode 错误码类型
type ErrorCode string

// 系统级错误码 (1000-1999)
const (
	Success         ErrorCode = "0000" // 成功
	SystemError     ErrorCode = "1000" // 系统错误
	DatabaseError   ErrorCode = "1001" // 数据库错误
	CacheError      ErrorCode = "1002" // 缓存错误
	NetworkError    ErrorCode = "1003" // 网络错误
	ServiceUnavail  ErrorCode = "1004" // 服务不可用
	InternalError   ErrorCode = "1005" // 内部错误
	ConfigError     ErrorCode = "1006" // 配置错误
	FileError       ErrorCode = "1007" // 文件操作错误
	JSONParseError  ErrorCode = "1008" // JSON解析错误
	ThirdPartyError ErrorCode = "1009" // 第三方服务错误
	MaintenanceMode ErrorCode = "1100" // 系统维护中
)

// 请求相关错误码 (2000-2999)
const (
	InvalidRequest     ErrorCode = "2000" // 无效请求
	InvalidParams      ErrorCode = "2001" // 参数错误
	MissingParams      ErrorCode = "2002" // 缺少参数
	InvalidJSON        ErrorCode = "2003" // JSON格式错误
	InvalidMethod      ErrorCode = "2004" // 请求方法错误
	RequestTooLarge    ErrorCode = "2005" // 请求体过大
	RateLimitExceeded  ErrorCode = "2006" // 请求频率限制
	InvalidContentType ErrorCode = "2007" // 内容类型错误
	InvalidEncoding    ErrorCode = "2008" // 编码错误
	RequestTimeout     ErrorCode = "2009" // 请求超时
)

// 认证相关错误码 (3000-3999)
const (
	AuthenticationFailed    ErrorCode = "3000" // 认证失败
	InvalidToken            ErrorCode = "3001" // 无效令牌
	TokenExpired            ErrorCode = "3002" // 令牌过期
	InvalidCredentials      ErrorCode = "3003" // 凭据无效
	AccessDenied            ErrorCode = "3004" // 访问被拒绝
	PermissionDenied        ErrorCode = "3005" // 权限不足
	AccountDisabled         ErrorCode = "3006" // 账户被禁用
	AccountLocked           ErrorCode = "3007" // 账户被锁定
	LoginRequired           ErrorCode = "3008" // 需要登录
	RefreshTokenExpired     ErrorCode = "3009" // 刷新令牌过期
	InvalidSignature        ErrorCode = "3010" // 签名无效
	InvalidAPIKey           ErrorCode = "3011" // API密钥无效
	InsufficientPermissions ErrorCode = "3012" // 权限不足
	TwoFactorRequired       ErrorCode = "3013" // 需要双因子认证
	InvalidTwoFactorCode    ErrorCode = "3014" // 双因子认证码无效
	AccountSuspended        ErrorCode = "3015" // 账户被暂停
	IPNotAllowed            ErrorCode = "3016" // IP地址不被允许
	SessionLimitExceeded    ErrorCode = "3017" // 会话限制超出
)

// 用户相关错误码 (4000-4999)
const (
	UserNotFound             ErrorCode = "4000" // 用户不存在
	UserAlreadyExists        ErrorCode = "4001" // 用户已存在
	EmailAlreadyExists       ErrorCode = "4002" // 邮箱已存在
	PhoneAlreadyExists       ErrorCode = "4003" // 手机号已存在
	UsernameAlreadyExists    ErrorCode = "4004" // 用户名已存在
	InvalidPassword          ErrorCode = "4005" // 密码错误
	WeakPassword             ErrorCode = "4006" // 密码强度不够
	PasswordMismatch         ErrorCode = "4007" // 密码不匹配
	InvalidEmail             ErrorCode = "4008" // 邮箱格式错误
	InvalidPhone             ErrorCode = "4009" // 手机号格式错误
	UserNotActive            ErrorCode = "4010" // 用户未激活
	UserNotVerified          ErrorCode = "4011" // 用户未验证
	ProfileIncomplete        ErrorCode = "4012" // 个人资料不完整
	AgeRestriction           ErrorCode = "4013" // 年龄限制
	InvalidUserType          ErrorCode = "4014" // 用户类型无效
	UserTypeExists           ErrorCode = "4015" // 该邮箱/手机号已被此用户类型注册
	UsernameExists           ErrorCode = "4016" // 用户名已存在
	EmailExists              ErrorCode = "4017" // 邮箱已存在
	UserHasUnCompletedOrders ErrorCode = "4018" // 用户有未完成的订单
)

// 验证相关错误码 (5000-5999)
const (
	VerificationCodeRequired   ErrorCode = "5000" // 需要验证码
	InvalidVerificationCode    ErrorCode = "5001" // 验证码错误
	VerificationCodeExpired    ErrorCode = "5002" // 验证码过期
	VerificationLimitReached   ErrorCode = "5003" // 验证次数限制
	EmailNotVerified           ErrorCode = "5004" // 邮箱未验证
	PhoneNotVerified           ErrorCode = "5005" // 手机号未验证
	AccountNotVerified         ErrorCode = "5006" // 账户未验证
	VerificationFailed         ErrorCode = "5007" // 验证失败
	IDVerificationRequired     ErrorCode = "5008" // 需要身份验证
	InvalidIDDocument          ErrorCode = "5009" // 身份证件无效
	IDVerificationFailed       ErrorCode = "5010" // 身份验证失败
	KYCRequired                ErrorCode = "5011" // 需要KYC验证
	KYCFailed                  ErrorCode = "5012" // KYC验证失败
	BackgroundCheckRequired    ErrorCode = "5013" // 需要背景调查
	BackgroundCheckFailed      ErrorCode = "5014" // 背景调查失败
	VerificationRequired       ErrorCode = "5015" // 需要验证
	TooManyAttempts            ErrorCode = "5016" // 尝试次数过多
	VerificationCooldown       ErrorCode = "5017" // 验证冷却期
	SMSServiceError            ErrorCode = "5018" // 短信服务错误
	EmailServiceError          ErrorCode = "5019" // 邮件服务错误
	InvalidVerificationMethod  ErrorCode = "5020" // 无效验证方式
	VerificationCodeSendFailed ErrorCode = "5021" // 验证码发送失败
)

// Checkout 相关错误码 (5100-5199)
const (
	CheckoutNotFound ErrorCode = "5100" // Checkout不存在
	CheckoutExpired  ErrorCode = "5101" // Checkout已过期
	Unauthorized     ErrorCode = "5102" // 未授权访问
)

// 行程相关错误码 (6000-6999)
const (
	RideNotFound                   ErrorCode = "6000" // 行程不存在
	RideAlreadyExists              ErrorCode = "6001" // 行程已存在
	RideNotAvailable               ErrorCode = "6002" // 行程不可用
	RideAlreadyBooked              ErrorCode = "6003" // 行程已被预订
	RideAlreadyCancelled           ErrorCode = "6004" // 行程已取消
	RideAlreadyCompleted           ErrorCode = "6005" // 行程已完成
	RideNotStarted                 ErrorCode = "6006" // 行程未开始
	RideInProgress                 ErrorCode = "6007" // 行程进行中
	InvalidRideStatus              ErrorCode = "6008" // 行程状态无效
	InvalidPickupLocation          ErrorCode = "6009" // 上车地点无效
	InvalidDropoffLocation         ErrorCode = "6010" // 下车地点无效
	InvalidRideDate                ErrorCode = "6011" // 行程日期无效
	InsufficientSeats              ErrorCode = "6012" // 座位不足
	DriverCannotBook               ErrorCode = "6013" // 司机不能预订自己的行程
	BookingDeadlinePassed          ErrorCode = "6014" // 预订截止时间已过
	CancellationNotAllowed         ErrorCode = "6015" // 不允许取消
	RideNotBookedByUser            ErrorCode = "6016" // 用户未预订此行程
	DriverOffline                  ErrorCode = "6017" // 司机离线
	DriverHasActiveOrder           ErrorCode = "6018" // 司机有未完成的订单，无法接单
	DriverHasActiveOrderInProgress ErrorCode = "6019" // 司机有在途订单，不能开启新行程
)

// 订单管理相关错误码 (6500-6599)
const (
	OrderNotFound             ErrorCode = "6500" // 订单不存在
	OrderAlreadyCancelled     ErrorCode = "6501" // 订单已取消
	OrderAlreadyCompleted     ErrorCode = "6502" // 订单已完成
	OrderCannotCancel         ErrorCode = "6503" // 订单无法取消
	OrderCannotUpdate         ErrorCode = "6504" // 订单无法更新
	InvalidOrderStatus        ErrorCode = "6505" // 订单状态无效
	InvalidPaymentStatus      ErrorCode = "6506" // 支付状态无效
	OrderSearchFailed         ErrorCode = "6507" // 订单搜索失败
	OrderUpdateFailed         ErrorCode = "6508" // 订单更新失败
	OrderCancelFailed         ErrorCode = "6509" // 订单取消失败
	OrderPermissionDenied     ErrorCode = "6510" // 订单操作权限不足
	OrderStatusUpdateFailed   ErrorCode = "6511" // 订单状态更新失败
	PaymentStatusUpdateFailed ErrorCode = "6512" // 支付状态更新失败
	OnlinePaymentPending      ErrorCode = "6513" // 有在线支付正在进行，禁止现金支付
)

// 车辆管理相关错误码 (6700-6799)
const (
	VehicleNotFound             ErrorCode = "6700" // 车辆不存在
	VehicleAlreadyExists        ErrorCode = "6701" // 车辆已存在
	VehicleNotAvailable         ErrorCode = "6702" // 车辆不可用
	VehicleAlreadyBound         ErrorCode = "6703" // 车辆已被绑定
	VehicleNotBound             ErrorCode = "6704" // 车辆未绑定
	VehicleNotVerified          ErrorCode = "6705" // 车辆未验证
	VehicleVerificationFailed   ErrorCode = "6706" // 车辆验证失败
	InvalidVehicleStatus        ErrorCode = "6707" // 车辆状态无效
	InvalidVehicleType          ErrorCode = "6708" // 车辆类型无效
	PlateNumberExists           ErrorCode = "6709" // 车牌号已存在
	VINExists                   ErrorCode = "6710" // VIN码已存在
	VehicleInUse                ErrorCode = "6711" // 车辆使用中
	InvalidYear                 ErrorCode = "6712" // 无效年份
	VehicleServiceUnavailable   ErrorCode = "6730"
	VehicleLocationUpdateFailed ErrorCode = "6731"
	VehicleMaintenanceFailed    ErrorCode = "6732"
	VehicleNeedsMaintenance     ErrorCode = "6733"
	VehicleInvalidStatus        ErrorCode = "6734"
	NotImplemented              ErrorCode = "6735"
	VehicleDocumentsMissing     ErrorCode = "6713" // 车辆文档缺失
	VehicleInsuranceExpired     ErrorCode = "6714" // 车辆保险过期
	VehicleRegistrationExpired  ErrorCode = "6715" // 车辆注册过期
	VehicleInspectionExpired    ErrorCode = "6716" // 车辆检验过期
	VehicleCreateFailed         ErrorCode = "6717" // 车辆创建失败
	VehicleUpdateFailed         ErrorCode = "6718" // 车辆更新失败
	VehicleDeleteFailed         ErrorCode = "6719" // 车辆删除失败
	VehicleSearchFailed         ErrorCode = "6720" // 车辆搜索失败
	VehiclePermissionDenied     ErrorCode = "6721" // 车辆操作权限不足
	VehicleStatusUpdateFailed   ErrorCode = "6722" // 车辆状态更新失败
	VehicleBindFailed           ErrorCode = "6723" // 车辆绑定失败
	VehicleUnbindFailed         ErrorCode = "6724" // 车辆解绑失败
	VehicleNotAssigned          ErrorCode = "6725" // 车辆未分派
	VehicleAssignFailed         ErrorCode = "6726" // 车辆分派失败
)

// 支付相关错误码 (7000-7999)
const (
	PaymentRequired           ErrorCode = "7000" // 需要支付
	PaymentFailed             ErrorCode = "7001" // 支付失败
	PaymentMethodRequired     ErrorCode = "7002" // 需要支付方式
	InvalidPaymentMethod      ErrorCode = "7003" // 支付方式无效
	PaymentAlreadyMade        ErrorCode = "7004" // 已经支付
	RefundFailed              ErrorCode = "7005" // 退款失败
	InsufficientFunds         ErrorCode = "7006" // 余额不足
	TransactionNotFound       ErrorCode = "7007" // 交易不存在
	DuplicateTransaction      ErrorCode = "7008" // 重复交易
	PaymentTimeout            ErrorCode = "7009" // 支付超时
	NoAvailablePaymentService ErrorCode = "7010" // 无可用支付服务
)

// 价格相关错误码 (7100-7199)
const (
	PriceIDNotFound     ErrorCode = "7100" // 价格ID不存在
	PriceIDExpired      ErrorCode = "7101" // 价格ID已过期
	PriceIDAlreadyUsed  ErrorCode = "7102" // 价格ID已被使用
	PriceMismatch       ErrorCode = "7103" // 价格不匹配
	PriceIDLocked       ErrorCode = "7104" // 价格ID已锁定
	PriceCalculateError ErrorCode = "7105" // 价格计算错误
	PriceIDUserMismatch ErrorCode = "7106" // 价格ID用户不匹配
)

// 评价相关错误码 (8000-8999)
const (
	RatingRequired      ErrorCode = "8000" // 需要评价
	RatingAlreadyExists ErrorCode = "8001" // 评价已存在
	InvalidRating       ErrorCode = "8002" // 评价无效
	RatingNotFound      ErrorCode = "8003" // 评价不存在
	CannotRateSelf      ErrorCode = "8004" // 不能评价自己
	RatingNotAllowed    ErrorCode = "8005" // 不允许评价
)

// 消息推送相关错误码 (9000-9999)
const (
	PushNotificationFailed ErrorCode = "9000" // 推送失败
	InvalidFCMToken        ErrorCode = "9001" // FCM令牌无效
	MessageTooLong         ErrorCode = "9002" // 消息过长
	InvalidRecipient       ErrorCode = "9003" // 接收者无效
	TemplateNotFound       ErrorCode = "9004" // 模板不存在
)

// 工具服务相关错误码 (9100-9199)
const (
	ExportFailed   ErrorCode = "9100" // 导出失败
	ImportFailed   ErrorCode = "9101" // 导入失败
	ProcessFailed  ErrorCode = "9102" // 处理失败
	GenerateFailed ErrorCode = "9103" // 生成失败
)

// 地理位置相关错误码 (9500-9599)
const (
	LocationRequired     ErrorCode = "9500" // 需要位置信息
	InvalidLocation      ErrorCode = "9501" // 位置信息无效
	LocationNotFound     ErrorCode = "9502" // 位置未找到
	RouteNotFound        ErrorCode = "9503" // 路线未找到
	DistanceTooLong      ErrorCode = "9504" // 距离太远
	LocationServiceError ErrorCode = "9505" // 位置服务错误
)

// 优惠码相关错误码 (10000-10999)
const (
	PromotionNotFound           ErrorCode = "10000" // 优惠码不存在
	PromotionAlreadyExists      ErrorCode = "10001" // 优惠码已存在
	PromotionExpired            ErrorCode = "10002" // 优惠码已过期
	PromotionInactive           ErrorCode = "10003" // 优惠码未激活
	PromotionUsageLimitExceeded ErrorCode = "10004" // 使用次数超限
	PromotionNotApplicable      ErrorCode = "10005" // 优惠码不适用
	PromotionAlreadyUsed        ErrorCode = "10006" // 优惠码已使用
	PromotionCreationFailed     ErrorCode = "10007" // 优惠码创建失败
	PromotionUpdateFailed       ErrorCode = "10008" // 优惠码更新失败
	PromotionDeleteFailed       ErrorCode = "10009" // 优惠码删除失败
	PromotionApprovalFailed     ErrorCode = "10010" // 优惠码审批失败
	PromotionStatusInvalid      ErrorCode = "10011" // 优惠码状态无效
	PromotionSearchFailed       ErrorCode = "10012" // 优惠码搜索失败
	PromotionPermissionDenied   ErrorCode = "10013" // 优惠码操作权限不足
	PromotionBudgetExceeded     ErrorCode = "10014" // 优惠码预算超限
	PromotionTimeInvalid        ErrorCode = "10015" // 优惠码时间设置无效
	PromotionValueInvalid       ErrorCode = "10016" // 优惠码数值无效
	PromotionTargetInvalid      ErrorCode = "10017" // 优惠码目标用户无效
	PromotionCombinationInvalid ErrorCode = "10018" // 优惠码组合无效
	PromotionSecurityViolation  ErrorCode = "10019" // 优惠码安全违规
	PromotionInUse              ErrorCode = "10020" // 优惠码正在使用中，无法删除
)

// GetMessage 获取错误码对应的英文消息
func (code ErrorCode) GetMessage() string {
	messages := map[ErrorCode]string{
		// 系统级错误码
		Success:         "Success",
		SystemError:     "System error",
		DatabaseError:   "Database error",
		CacheError:      "Cache error",
		NetworkError:    "Network error",
		ServiceUnavail:  "Service unavailable",
		InternalError:   "Internal error",
		ConfigError:     "Configuration error",
		FileError:       "File operation error",
		JSONParseError:  "JSON parsing error",
		ThirdPartyError: "Third-party service error",
		MaintenanceMode: "System is under maintenance",

		// 请求相关错误码
		InvalidRequest:     "Invalid request",
		InvalidParams:      "Invalid parameters",
		MissingParams:      "Missing required parameters",
		InvalidJSON:        "Invalid JSON format",
		InvalidMethod:      "Invalid request method",
		RequestTooLarge:    "Request body too large",
		RateLimitExceeded:  "Rate limit exceeded",
		InvalidContentType: "Invalid content type",
		InvalidEncoding:    "Invalid encoding",
		RequestTimeout:     "Request timeout",

		// 认证相关错误码
		AuthenticationFailed:    "Authentication failed",
		InvalidToken:            "Invalid token",
		TokenExpired:            "Token expired",
		InvalidCredentials:      "Invalid credentials",
		AccessDenied:            "Access denied",
		PermissionDenied:        "Permission denied",
		AccountDisabled:         "Account disabled",
		AccountLocked:           "Account locked",
		LoginRequired:           "Login required",
		RefreshTokenExpired:     "Refresh token expired",
		InvalidSignature:        "Invalid signature",
		InvalidAPIKey:           "Invalid API key",
		InsufficientPermissions: "Insufficient permissions",
		TwoFactorRequired:       "Two-factor authentication required",
		InvalidTwoFactorCode:    "Invalid two-factor authentication code",
		AccountSuspended:        "Account suspended",
		IPNotAllowed:            "IP address not allowed",
		SessionLimitExceeded:    "Session limit exceeded",

		// 用户相关错误码
		UserNotFound:             "User not found",
		UserAlreadyExists:        "User already exists",
		EmailAlreadyExists:       "Email already exists",
		PhoneAlreadyExists:       "Phone number already exists",
		UsernameAlreadyExists:    "Username already exists",
		InvalidPassword:          "Invalid password",
		WeakPassword:             "Password is too weak",
		PasswordMismatch:         "Password confirmation does not match",
		InvalidEmail:             "Invalid email format",
		InvalidPhone:             "Invalid phone number format",
		UserNotActive:            "User account is not active",
		UserNotVerified:          "User is not verified",
		ProfileIncomplete:        "Profile is incomplete",
		AgeRestriction:           "Age restriction applies",
		InvalidUserType:          "Invalid user type",
		UserTypeExists:           "This email/phone is already registered with this user type",
		UsernameExists:           "Username already exists",
		EmailExists:              "Email already exists",
		UserHasUnCompletedOrders: "User has uncompleted orders",

		// 验证相关错误码
		VerificationRequired:       "Verification required",
		InvalidVerificationCode:    "Invalid verification code",
		VerificationCodeExpired:    "Verification code expired",
		VerificationFailed:         "Verification failed",
		TooManyAttempts:            "Too many attempts",
		VerificationCooldown:       "Please wait before requesting another verification code",
		EmailNotVerified:           "Email not verified",
		PhoneNotVerified:           "Phone number not verified",
		SMSServiceError:            "SMS service error",
		EmailServiceError:          "Email service error",
		InvalidVerificationMethod:  "Invalid verification method",
		VerificationCodeSendFailed: "Failed to send verification code",

		// 行程相关错误码
		RideNotFound:           "Ride not found",
		RideAlreadyExists:      "Ride already exists",
		RideNotAvailable:       "Ride not available",
		RideAlreadyBooked:      "Ride already booked",
		RideAlreadyCancelled:   "Ride already cancelled",
		RideAlreadyCompleted:   "Ride already completed",
		RideNotStarted:         "Ride not started",
		RideInProgress:         "Ride in progress",
		InvalidRideStatus:      "Invalid ride status",
		InvalidPickupLocation:  "Invalid pickup location",
		InvalidDropoffLocation: "Invalid dropoff location",
		InvalidRideDate:        "Invalid ride date",
		InsufficientSeats:      "Insufficient seats",
		DriverCannotBook:       "Driver cannot book their own ride",
		BookingDeadlinePassed:  "Booking deadline has passed",
		CancellationNotAllowed: "Cancellation not allowed",
		RideNotBookedByUser:    "User has not booked this ride",
		DriverOffline:          "Driver is offline",
		DriverHasActiveOrder:   "Driver has active orders and cannot accept new orders",

		// 订单管理相关错误码
		OrderNotFound:             "Order not found",
		OrderAlreadyCancelled:     "Order already cancelled",
		OrderAlreadyCompleted:     "Order already completed",
		OrderCannotCancel:         "Order cannot be cancelled",
		OrderCannotUpdate:         "Order cannot be updated",
		InvalidOrderStatus:        "Invalid order status",
		InvalidPaymentStatus:      "Invalid payment status",
		OrderSearchFailed:         "Order search failed",
		OrderUpdateFailed:         "Order update failed",
		OrderCancelFailed:         "Order cancellation failed",
		OrderPermissionDenied:     "Insufficient permissions for order operation",
		OrderStatusUpdateFailed:   "Order status update failed",
		PaymentStatusUpdateFailed: "Payment status update failed",
		OnlinePaymentPending:      "Online payment is pending, cash payment is not allowed",

		// 车辆管理相关错误码
		VehicleNotFound:             "Vehicle not found",
		VehicleAlreadyExists:        "Vehicle already exists",
		VehicleNotAvailable:         "Vehicle not available",
		VehicleAlreadyBound:         "Vehicle already bound",
		VehicleNotBound:             "Vehicle not bound",
		VehicleNotVerified:          "Vehicle not verified",
		VehicleVerificationFailed:   "Vehicle verification failed",
		InvalidVehicleStatus:        "Invalid vehicle status",
		InvalidVehicleType:          "Invalid vehicle type",
		PlateNumberExists:           "Plate number already exists",
		VINExists:                   "VIN already exists",
		VehicleInUse:                "Vehicle is in use",
		VehicleNeedsMaintenance:     "Vehicle needs maintenance",
		VehicleDocumentsMissing:     "Vehicle documents missing",
		VehicleInsuranceExpired:     "Vehicle insurance expired",
		VehicleRegistrationExpired:  "Vehicle registration expired",
		VehicleInspectionExpired:    "Vehicle inspection expired",
		VehicleCreateFailed:         "Vehicle creation failed",
		VehicleUpdateFailed:         "Vehicle update failed",
		VehicleDeleteFailed:         "Vehicle deletion failed",
		VehicleSearchFailed:         "Vehicle search failed",
		VehiclePermissionDenied:     "Vehicle operation permission denied",
		VehicleStatusUpdateFailed:   "Vehicle status update failed",
		VehicleBindFailed:           "Vehicle binding failed",
		VehicleUnbindFailed:         "Vehicle unbinding failed",
		VehicleNotAssigned:          "Vehicle not assigned",
		VehicleAssignFailed:         "Vehicle assignment failed",
		VehicleServiceUnavailable:   "Vehicle service unavailable",
		VehicleLocationUpdateFailed: "Vehicle location update failed",
		VehicleMaintenanceFailed:    "Vehicle maintenance failed",
		VehicleInvalidStatus:        "Invalid vehicle status",
		NotImplemented:              "Feature not implemented",

		// 支付相关错误码
		PaymentRequired:       "Payment required",
		PaymentFailed:         "Payment failed",
		PaymentMethodRequired: "Payment method required",
		InvalidPaymentMethod:  "Invalid payment method",
		PaymentAlreadyMade:    "Payment already made",
		RefundFailed:          "Refund failed",
		InsufficientFunds:     "Insufficient funds",
		TransactionNotFound:   "Transaction not found",
		DuplicateTransaction:  "Duplicate transaction",
		PaymentTimeout:        "Payment timeout",

		// 评价相关错误码
		RatingRequired:      "Rating required",
		RatingAlreadyExists: "Rating already exists",
		InvalidRating:       "Invalid rating",
		RatingNotFound:      "Rating not found",
		CannotRateSelf:      "Cannot rate yourself",
		RatingNotAllowed:    "Rating not allowed",

		// 消息推送相关错误码
		PushNotificationFailed: "Push notification failed",
		InvalidFCMToken:        "Invalid FCM token",
		MessageTooLong:         "Message too long",
		InvalidRecipient:       "Invalid recipient",
		TemplateNotFound:       "Template not found",

		// 工具服务相关错误码
		ExportFailed:   "Export failed",
		ImportFailed:   "Import failed",
		ProcessFailed:  "Process failed",
		GenerateFailed: "Generate failed",

		// 地理位置相关错误码
		LocationRequired:     "Location required",
		InvalidLocation:      "Invalid location",
		LocationNotFound:     "Location not found",
		RouteNotFound:        "Route not found",
		DistanceTooLong:      "Distance too long",
		LocationServiceError: "Location service error",

		// 优惠码相关错误码
		PromotionNotFound:           "Promo code not found",
		PromotionAlreadyExists:      "Promo code already exists",
		PromotionExpired:            "Promo code has expired",
		PromotionInactive:           "Promo code is not active",
		PromotionUsageLimitExceeded: "Promo code usage limit exceeded",
		PromotionNotApplicable:      "Promo code is not applicable",
		PromotionAlreadyUsed:        "Promo code already used",
		PromotionCreationFailed:     "Promo code creation failed",
		PromotionUpdateFailed:       "Promo code update failed",
		PromotionDeleteFailed:       "Promo code deletion failed",
		PromotionApprovalFailed:     "Promo code approval failed",
		PromotionStatusInvalid:      "Invalid promo code status",
		PromotionSearchFailed:       "Promo code search failed",
		PromotionPermissionDenied:   "Insufficient permissions for promo code operation",
		PromotionBudgetExceeded:     "Promo code budget exceeded",
		PromotionTimeInvalid:        "Invalid promo code time settings",
		PromotionValueInvalid:       "Invalid promo code value",
		PromotionTargetInvalid:      "Invalid promo code target users",
		PromotionCombinationInvalid: "Invalid promo code combination",
		PromotionSecurityViolation:  "Promo code security violation",
		PromotionInUse:              "Promo code is in use and cannot be deleted",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}
	return "Unknown error"
}

// GetCode 获取错误码字符串值
func (code ErrorCode) GetCode() string {
	return string(code)
}

// GetCodeInt 获取错误码数值（为了向后兼容）
func (code ErrorCode) GetCodeInt() int {
	switch code {
	case Success:
		return 0
	case SystemError:
		return 1000
	case DatabaseError:
		return 1001
	case CacheError:
		return 1002
	case NetworkError:
		return 1003
	case ServiceUnavail:
		return 1004
	case InternalError:
		return 1005
	case MaintenanceMode:
		return 1100
	case InvalidRequest:
		return 2000
	case InvalidParams:
		return 2001
	case MissingParams:
		return 2002
	case InvalidJSON:
		return 2003
	case AuthenticationFailed:
		return 3000
	case InvalidToken:
		return 3001
	case TokenExpired:
		return 3002
	case InvalidCredentials:
		return 3003
	case AccessDenied:
		return 3004
	case AccountDisabled:
		return 3006
	case UserNotFound:
		return 4000
	case UserAlreadyExists:
		return 4001
	case EmailAlreadyExists:
		return 4002
	case PhoneAlreadyExists:
		return 4003
	case InvalidPassword:
		return 4005
	case WeakPassword:
		return 4006
	case PasswordMismatch:
		return 4007
	case InvalidEmail:
		return 4008
	case InvalidPhone:
		return 4009
	case UserNotActive:
		return 4010
	case InvalidUserType:
		return 4014
	case UserTypeExists:
		return 4015
	case UserHasUnCompletedOrders:
		return 4018
	case VerificationRequired:
		return 5000
	case InvalidVerificationCode:
		return 5001
	case VerificationCodeExpired:
		return 5002
	case SMSServiceError:
		return 5008
	case EmailServiceError:
		return 5009
	case VerificationCodeSendFailed:
		return 5011
	case RideNotFound:
		return 6000
	case RideAlreadyBooked:
		return 6003
	case DriverOffline:
		return 6017
	case DriverHasActiveOrder:
		return 6018
	case OrderNotFound:
		return 6500
	case OrderAlreadyCancelled:
		return 6501
	case OrderAlreadyCompleted:
		return 6502
	case OrderCannotCancel:
		return 6503
	case OrderCannotUpdate:
		return 6504
	case InvalidOrderStatus:
		return 6505
	case InvalidPaymentStatus:
		return 6506
	case OrderSearchFailed:
		return 6507
	case OrderUpdateFailed:
		return 6508
	case OrderCancelFailed:
		return 6509
	case OrderPermissionDenied:
		return 6510
	case OrderStatusUpdateFailed:
		return 6511
	case PaymentStatusUpdateFailed:
		return 6512
	case VehicleNotFound:
		return 6700
	case VehicleAlreadyExists:
		return 6701
	case VehicleNotAvailable:
		return 6702
	case VehicleAlreadyBound:
		return 6703
	case VehicleNotBound:
		return 6704
	case VehicleNotVerified:
		return 6705
	case VehicleVerificationFailed:
		return 6706
	case InvalidVehicleStatus:
		return 6707
	case InvalidVehicleType:
		return 6708
	case PlateNumberExists:
		return 6709
	case VINExists:
		return 6710
	case VehicleInUse:
		return 6711
	case InvalidYear:
		return 6712
	case VehicleNeedsMaintenance:
		return 6733
	case VehicleDocumentsMissing:
		return 6713
	case VehicleInsuranceExpired:
		return 6714
	case VehicleRegistrationExpired:
		return 6715
	case VehicleInspectionExpired:
		return 6716
	case VehicleCreateFailed:
		return 6717
	case VehicleUpdateFailed:
		return 6718
	case VehicleDeleteFailed:
		return 6719
	case VehicleSearchFailed:
		return 6720
	case VehiclePermissionDenied:
		return 6721
	case VehicleServiceUnavailable:
		return 6730
	case VehicleLocationUpdateFailed:
		return 6731
	case VehicleMaintenanceFailed:
		return 6732
	case VehicleInvalidStatus:
		return 6734
	case NotImplemented:
		return 6735
	case VehicleStatusUpdateFailed:
		return 6722
	case VehicleBindFailed:
		return 6723
	case VehicleUnbindFailed:
		return 6724
	case VehicleNotAssigned:
		return 6725
	case VehicleAssignFailed:
		return 6726
	case PromotionNotFound:
		return 10000
	case PromotionAlreadyExists:
		return 10001
	case PromotionExpired:
		return 10002
	case PromotionInactive:
		return 10003
	case PromotionUsageLimitExceeded:
		return 10004
	case PromotionNotApplicable:
		return 10005
	case PromotionAlreadyUsed:
		return 10006
	case PromotionCreationFailed:
		return 10007
	case PromotionUpdateFailed:
		return 10008
	case PromotionDeleteFailed:
		return 10009
	case PromotionApprovalFailed:
		return 10010
	case PromotionStatusInvalid:
		return 10011
	case PromotionSearchFailed:
		return 10012
	case PromotionPermissionDenied:
		return 10013
	case PromotionBudgetExceeded:
		return 10014
	case PromotionTimeInvalid:
		return 10015
	case PromotionValueInvalid:
		return 10016
	case PromotionTargetInvalid:
		return 10017
	case PromotionCombinationInvalid:
		return 10018
	case PromotionSecurityViolation:
		return 10019
	case PromotionInUse:
		return 10020
	default:
		return 9999 // 未知错误
	}
}

// IsSuccess 判断是否成功
func (code ErrorCode) IsSuccess() bool {
	return code == Success
}

// IsSystemError 判断是否系统错误 (1000-1999)
func (code ErrorCode) IsSystemError() bool {
	return code >= "1000" && code < "2000"
}

// IsRequestError 判断是否请求错误 (2000-2999)
func (code ErrorCode) IsRequestError() bool {
	return code >= "2000" && code < "3000"
}

// IsAuthError 判断是否认证错误 (3000-3999)
func (code ErrorCode) IsAuthError() bool {
	return code >= "3000" && code < "4000"
}

// IsUserError 判断是否用户错误 (4000-4999)
func (code ErrorCode) IsUserError() bool {
	return code >= "4000" && code < "5000"
}

// IsVerificationError 判断是否验证错误 (5000-5999)
func (code ErrorCode) IsVerificationError() bool {
	return code >= "5000" && code < "6000"
}

// IsBusinessError 判断是否业务错误 (6000+)
func (code ErrorCode) IsBusinessError() bool {
	return code >= "6000"
}

// IsPromotionError 判断是否优惠码错误 (10000-10999)
func (code ErrorCode) IsPromotionError() bool {
	return code >= "10000" && code < "11000"
}

// ServiceError 服务层错误类型
type ServiceError struct {
	Code    ErrorCode
	Message string
}

// Error 实现error接口
func (e *ServiceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

// NewServiceError 创建服务错误
func NewServiceError(code ErrorCode, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

// IsServiceError 判断是否为ServiceError类型
func IsServiceError(err error) (*ServiceError, bool) {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr, true
	}
	return nil, false
}
