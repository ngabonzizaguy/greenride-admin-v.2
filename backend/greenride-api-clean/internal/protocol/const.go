package protocol

const (
	System  = "system"
	Default = "default"
)

const (
	EnvSandbox    = "sandbox"
	EnvProduction = "production"
)

// 通用状态常量
const (
	StatusOn  = "on"
	StatusOff = "off"
)

// 用户状态常量
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusSuspended = "suspended"
	StatusBanned    = "banned"
	StatusDeleted   = "deleted"

	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusFailed     = "failed"
	StatusRefunded   = "refunded"
	StatusApproved   = "approved" // 已审批
	StatusRejected   = "rejected" // 已拒绝
	StatusResolved   = "resolved" // 已解决

	StatusRequested     = "requested"      // 用户下单
	StatusAccepted      = "accepted"       // 司机接单
	StatusDriverComing  = "driver_coming"  // 司机前往
	StatusDriverArrived = "driver_arrived" // 司机到达
	StatusTripEnded     = "trip_ended"     // 行程结束
	StatusInProgress    = "in_progress"    // 行程进行中
	StatusPaid          = "paid"           // 已支付
	StatusCompleted     = "completed"      // 已完成
	StatusCancelled     = "cancelled"      // 已取消
	StatusOnline        = "online"
	StatusOffline       = "offline"
	StatusBusy          = "busy"
	StatusSuccess       = "success"

	// 反馈相关状态
	StatusReviewing = "reviewing" // 审核中
	StatusClosed    = "closed"    // 已关闭

	// 车辆特有状态
	StatusMaintenance = "maintenance" // 维护中
	StatusRetired     = "retired"     // 已退役
	StatusUnverified  = "unverified"  // 未验证
	StatusVerified    = "verified"    // 已验证

	// 价格规则状态常量
	StatusDraft      = "draft"
	StatusTesting    = "testing"
	StatusPaused     = "paused"
	StatusExpired    = "expired"
	StatusDeprecated = "deprecated"
	StatusArchived   = "archived"

	// 优惠券状态常量
	StatusAvailable = "available" // 可用
	StatusUsed      = "used"      // 已使用
)

var (
	CURRENT_RIDE_STATUS = []string{
		StatusInProgress,
		StatusDriverComing,
		StatusArchived,
	}
)

// 用户类型常量
const (
	UserTypePassenger = "passenger"
	UserTypeDriver    = "driver"
	UserTypeUser      = "user"
	UserTypeProvider  = "provider"
)

// 订单类型常量
const (
	RideOrder     = "ride"
	DeliveryOrder = "delivery"
	ShoppingOrder = "shopping"
)

// 服务类型常量
const (
	ServiceTypeStandard = "standard"
	ServiceTypePremium  = "premium"
	ServiceTypeLuxury   = "luxury"
)

// 计划类型常量
const (
	ScheduleTypeInstant   = "instant"
	ScheduleTypeScheduled = "scheduled"
)

// MessageType 消息类型常量
const (
	MsgTypePasswordReset       = "password_reset"
	MsgTypeAccountVerification = "account_verification"
	MsgTypePasswordUpdate      = "password_update"
)

const (
	VerifyCodeTypeRegister      = "register"       // 注册验证码
	VerifyCodeTypeLogin         = "login"          // 登录验证码
	VerifyCodeTypeResetPassword = "reset_password" // 重置密码验证码
)

// TaskStatus 任务状态
const (
	TaskStatusEnabled  = "enabled"  // 启用
	TaskStatusDisabled = "disabled" // 禁用
)

// Signal Handler Key 信号处理器标识
const (
	SignalProcessorHandler    = "signal_processor"          // 信号处理器
	OrderSubmitHandler        = "order_submit_handler"      // 订单提交处理器
	OrderCompleteHandler      = "order_complete_handler"    // 订单完成处理器
	OrderCancelHandler        = "order_cancel_handler"      // 订单取消处理器
	DriverLocationHandler     = "driver_location_handler"   // 司机位置更新处理器
	PaymentProcessedHandler   = "payment_processed_handler" // 支付处理器
	UserRegisteredHandler     = "user_registered_handler"   // 用户注册处理器
	UserVerifiedHandler       = "user_verified_handler"     // 用户验证处理器
	DispatchOptimizeHandler   = "dispatch_optimize_handler" // 调度优化处理器
	PriceUpdateHandler        = "price_update_handler"      // 价格更新处理器
	TodayStatsHandler         = "today_stats_handler"       // 当日统计处理器
	YesterdayStatsHandler     = "yesterday_stats_handler"   // 昨日统计处理器
	PaymentChannelSyncHandler = "payment_channel_sync"      // 支付渠道同步处理器
)

// Signal Type 信号类型
const (
	SignalOrderSubmit      = "order_submit"      // 订单提交信号
	SignalOrderComplete    = "order_complete"    // 订单完成信号
	SignalOrderCancel      = "order_cancel"      // 订单取消信号
	SignalDriverLocation   = "driver_location"   // 司机位置更新信号
	SignalPaymentProcessed = "payment_processed" // 支付处理信号
	SignalUserRegistered   = "user_registered"   // 用户注册信号
	SignalUserVerified     = "user_verified"     // 用户验证信号
	SignalDispatchOptimize = "dispatch_optimize" // 调度优化信号
	SignalPriceUpdate      = "price_update"      // 价格更新信号
)
