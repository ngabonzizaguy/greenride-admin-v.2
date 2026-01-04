package protocol

// FCM 通知类型常量
const (
	// 订单相关通知类型
	NotificationTypeOrderAccepted     = "order_accepted"      // 订单已被接单
	NotificationTypeDriverArrived     = "driver_arrived"      // 司机已到达
	NotificationTypeTripStarted       = "trip_started"        // 行程开始
	NotificationTypeTripEnded         = "trip_ended"          // 行程结束
	NotificationTypePaymentConfirmed  = "payment_confirmed"   // 支付确认
	NotificationTypeOrderCancelled    = "order_cancelled"     // 订单已取消
	NotificationTypeNewOrderAvailable = "new_order_available" // 新订单可用
)
