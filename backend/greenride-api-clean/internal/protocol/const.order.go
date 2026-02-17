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

// =============================================================================
// 订单取消原因枚举 (Cancellation Reasons)
// =============================================================================

// CancelReason represents a predefined cancellation reason
type CancelReason struct {
	ID     string `json:"id"`     // unique identifier — Flutter app uses this for RadioListTile value
	Key    string `json:"key"`
	Label  string `json:"label"`
	Reason string `json:"reason"` // alias for Label — Flutter app reads this field
}

// 司机取消原因
const (
	CancelReasonDriverPassengerNoShow   = "passenger_not_at_pickup"
	CancelReasonDriverTooFar            = "too_far_away"
	CancelReasonDriverVehicleIssue      = "vehicle_issue"
	CancelReasonDriverEmergency         = "emergency"
	CancelReasonDriverOther             = "other"
)

// 乘客取消原因
const (
	CancelReasonPassengerDriverTooLong  = "driver_taking_too_long"
	CancelReasonPassengerChangedMind    = "changed_my_mind"
	CancelReasonPassengerFoundAnother   = "found_another_ride"
	CancelReasonPassengerEmergency      = "emergency"
	CancelReasonPassengerOther          = "other"
)

// GetDriverCancelReasons returns predefined cancellation reasons for drivers
func GetDriverCancelReasons() []CancelReason {
	return []CancelReason{
		{ID: CancelReasonDriverPassengerNoShow, Key: CancelReasonDriverPassengerNoShow, Label: "Passenger not at pickup", Reason: "Passenger not at pickup"},
		{ID: CancelReasonDriverTooFar, Key: CancelReasonDriverTooFar, Label: "Too far away", Reason: "Too far away"},
		{ID: CancelReasonDriverVehicleIssue, Key: CancelReasonDriverVehicleIssue, Label: "Vehicle issue", Reason: "Vehicle issue"},
		{ID: CancelReasonDriverEmergency, Key: CancelReasonDriverEmergency, Label: "Emergency", Reason: "Emergency"},
		{ID: CancelReasonDriverOther, Key: CancelReasonDriverOther, Label: "Other", Reason: "Other"},
	}
}

// GetPassengerCancelReasons returns predefined cancellation reasons for passengers
func GetPassengerCancelReasons() []CancelReason {
	return []CancelReason{
		{ID: CancelReasonPassengerDriverTooLong, Key: CancelReasonPassengerDriverTooLong, Label: "Driver taking too long", Reason: "Driver taking too long"},
		{ID: CancelReasonPassengerChangedMind, Key: CancelReasonPassengerChangedMind, Label: "Changed my mind", Reason: "Changed my mind"},
		{ID: CancelReasonPassengerFoundAnother, Key: CancelReasonPassengerFoundAnother, Label: "Found another ride", Reason: "Found another ride"},
		{ID: CancelReasonPassengerEmergency, Key: CancelReasonPassengerEmergency, Label: "Emergency", Reason: "Emergency"},
		{ID: CancelReasonPassengerOther, Key: CancelReasonPassengerOther, Label: "Other", Reason: "Other"},
	}
}

// GetCancelReasonsByUserType returns cancellation reasons for the given user type
func GetCancelReasonsByUserType(userType string) []CancelReason {
	if userType == UserTypeDriver || userType == UserTypeProvider {
		return GetDriverCancelReasons()
	}
	return GetPassengerCancelReasons()
}

// IsValidCancelReasonKey checks if a reason key is valid for the given user type
func IsValidCancelReasonKey(key, userType string) bool {
	reasons := GetCancelReasonsByUserType(userType)
	for _, r := range reasons {
		if r.Key == key {
			return true
		}
	}
	return false
}

// GetCancelReasonLabel returns the label for a reason key, or empty string if not found
func GetCancelReasonLabel(key, userType string) string {
	reasons := GetCancelReasonsByUserType(userType)
	for _, r := range reasons {
		if r.Key == key {
			return r.Label
		}
	}
	return ""
}
