package protocol

// Order History Log Action Types
const (
	// 订单创建相关
	ActionOrderCreated   = "order_created"
	ActionOrderRequested = "order_requested"
	ActionOrderScheduled = "order_scheduled"

	// 派单相关
	ActionOrderDispatched = "order_dispatched"
	ActionOrderAssigned   = "order_assigned"
	ActionOrderAccepted   = "order_accepted"
	ActionOrderRejected   = "order_rejected"

	// 行程相关
	ActionDriverArrived   = "driver_arrived"
	ActionOrderStarted    = "order_started"
	ActionOrderInProgress = "order_in_progress"
	ActionOrderCompleted  = "order_completed"

	// 支付相关
	ActionPaymentInitiated = "payment_initiated"
	ActionPaymentCompleted = "payment_completed"
	ActionPaymentFailed    = "payment_failed"
	ActionOrderRefunded    = "order_refunded"

	// 取消相关
	ActionOrderCancelled           = "order_cancelled"
	ActionOrderCancelledByUser     = "order_cancelled_by_user"
	ActionOrderCancelledByProvider = "order_cancelled_by_provider"
	ActionOrderCancelledBySystem   = "order_cancelled_by_system"

	// 状态变更相关
	ActionStatusChanged   = "status_changed"
	ActionMetadataChanged = "metadata_changed"

	// 系统相关
	ActionSystemProcessed    = "system_processed"
	ActionManualIntervention = "manual_intervention"
)
