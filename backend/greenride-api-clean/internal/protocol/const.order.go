package protocol

// 订单拒绝原因枚举
const (
	// 距离相关
	RejectReasonDistanceFar = "distance_too_far"      // 距离太远
	RejectReasonPickupIssue = "pickup_location_issue" // 接客地点有问题

	// 时间相关
	RejectReasonDriverBusy = "driver_busy" // 司机忙碌中
	RejectReasonTrafficJam = "traffic_jam" // 交通拥堵
	RejectReasonBreakTime  = "break_time"  // 休息时间

	// 订单相关
	RejectReasonDestinationUnsafe = "destination_unsafe" // 目的地不安全
	RejectReasonShortTrip         = "short_trip"         // 行程太短
	RejectReasonPaymentMethod     = "payment_method"     // 支付方式问题

	// 外部因素
	RejectReasonWeatherBad        = "weather_bad"        // 天气恶劣
	RejectReasonVehicleIssue      = "vehicle_issue"      // 车辆问题
	RejectReasonPersonalEmergency = "personal_emergency" // 个人紧急情况

	// 其他
	RejectReasonOther = "other" // 其他原因（需要自定义文本）
)

// GetValidRejectReasons 获取所有有效的拒绝原因枚举
func GetValidRejectReasons() []string {
	return []string{
		RejectReasonDistanceFar,
		RejectReasonPickupIssue,
		RejectReasonDriverBusy,
		RejectReasonTrafficJam,
		RejectReasonBreakTime,
		RejectReasonDestinationUnsafe,
		RejectReasonShortTrip,
		RejectReasonPaymentMethod,
		RejectReasonWeatherBad,
		RejectReasonVehicleIssue,
		RejectReasonPersonalEmergency,
		RejectReasonOther,
	}
}

// IsValidRejectReason 验证拒绝原因是否有效
func IsValidRejectReason(reason string) bool {
	validReasons := GetValidRejectReasons()
	for _, valid := range validReasons {
		if reason == valid {
			return true
		}
	}
	return false
}
