package protocol

// 消息渠道常量
const (
	MsgChannelEmail = "email"
	MsgChannelFcm   = "fcm"
	MsgChannelSms   = "sms"
)

// 消息类型常量
const (
	// 通用消息类型
	MsgTypeGeneric         = "generic"
	MsgTypeVerifyCode      = "verify_code"
	MsgTypeRegisterSuccess = "register_success"

	// 乘客通知类型
	MsgTypePassengerOrderAccepted    = "passenger_order_accepted"
	MsgTypePassengerDriverArrived    = "passenger_driver_arrived"
	MsgTypePassengerTripStarted      = "passenger_trip_started"
	MsgTypePassengerTripEnded        = "passenger_trip_ended"
	MsgTypePassengerPaymentConfirmed = "passenger_payment_confirmed"
	MsgTypePassengerOrderCancelled   = "passenger_order_cancelled"

	// 司机通知类型
	MsgTypeDriverNewOrder         = "driver_new_order"
	MsgTypeDriverTripEnded        = "driver_trip_ended"
	MsgTypeDriverPaymentConfirmed = "driver_payment_confirmed"
	MsgTypeDriverOrderCancelled   = "driver_order_cancelled"
)

// 语言常量
const (
	LangEnglish     = "en"
	LangFrench      = "fr"
	LangChinese     = "zh"
	LangKinyarwanda = "rw"
)

// 地区常量
const (
	RegionRwanda   = "RW"
	RegionKenya    = "KE"
	RegionTanzania = "TZ"
	RegionUganda   = "UG"
)
