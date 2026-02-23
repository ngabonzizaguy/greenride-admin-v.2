package protocol

// Order 统一订单实体
type Order struct {
	// 订单基础信息
	OrderID       string `json:"order_id"`
	OrderType     string `json:"order_type"`
	UserID        string `json:"user_id"`
	ProviderID    string `json:"provider_id,omitempty"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`

	// 金额信息
	OriginalAmount      float64 `json:"original_amount"`   // 预估原始金额 (优惠前)
	DiscountedAmount    float64 `json:"discounted_amount"` // 优惠后总金额
	PaymentAmount       float64 `json:"payment_amount"`    // 最终付款金额
	TotalDiscountAmount float64 `json:"discount_amount"`   // 总优惠金额
	PlatformFee         float64 `json:"platform_fee"`
	PromoDiscount       float64 `json:"promo_discount"`
	UserPromoDiscount   float64 `json:"user_promo_discount"` // 用户优惠券折扣金额
	CancellationFee     float64 `json:"cancellation_fee"`
	Currency            string  `json:"currency"`

	// 支付信息
	PaymentMethod      string   `json:"payment_method"`
	PaymentID          string   `json:"payment_id,omitempty"`
	ChannelPaymentID   string   `json:"channel_payment_id,omitempty"`
	PaymentResult      string   `json:"payment_result,omitempty"`       // 支付结果详情
	PaymentRedirectURL string   `json:"payment_redirect_url,omitempty"` // 支付重定向URL
	PromoCodes         []string `json:"promo_codes,omitempty"`
	UserPromotionIDs   []string `json:"user_promotion_ids,omitempty"` // 使用的用户优惠券ID列表

	// 时间信息
	ScheduleType string `json:"schedule_type"`
	ScheduledAt  int64  `json:"scheduled_at,omitempty"`
	CreatedAt    int64  `json:"created_at"`
	AcceptedAt   int64  `json:"accepted_at,omitempty"`
	StartedAt    int64  `json:"started_at,omitempty"`
	EndedAt      int64  `json:"ended_at,omitempty"`
	CompletedAt  int64  `json:"completed_at,omitempty"`
	CancelledAt  int64  `json:"cancelled_at,omitempty"`

	// 取消信息
	CancelledBy  string `json:"cancelled_by,omitempty"`
	CancelReason string `json:"cancel_reason,omitempty"`

	// 派单相关信息
	DispatchID          string `json:"dispatch_id,omitempty"`   // 当司机从「派给我的单」接单时必填，用于 first-accept-wins
	DispatchStatus      string `json:"dispatch_status"`
	CurrentRound        int    `json:"current_round"`
	MaxRounds           int    `json:"max_rounds"`
	DispatchStartedAt   int64  `json:"dispatch_started_at,omitempty"`
	LastDispatchedAt    int64  `json:"last_dispatched_at,omitempty"`
	AutoDispatchEnabled bool   `json:"auto_dispatch_enabled"`
	DispatchStrategy    string `json:"dispatch_strategy,omitempty"`
	NextStrategy        string `json:"next_strategy,omitempty"`

	// 其他信息
	Notes    string         `json:"notes,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Version  int64          `json:"version"`

	// 订单详情（兼容所有订单类型）
	Details *OrderDetail `json:"details,omitempty"`

	Passenger        *User     `json:"passenger,omitempty"`
	Driver           *User     `json:"driver,omitempty"`
	Vehicle          *Vehicle  `json:"vehicle,omitempty"`
	PassengerRatings []*Rating `json:"passenger_ratings,omitempty"`
	DriverRatings    []*Rating `json:"driver_ratings,omitempty"`
}

// OrderDetail 统一订单详情（兼容所有订单类型）
type OrderDetail struct {
	// 通用字段
	OrderID         string `json:"order_id,omitempty"`
	OrderType       string `json:"order_type,omitempty"`
	VehicleCategory string `json:"vehicle_category,omitempty"`
	VehicleLevel    string `json:"vehicle_level,omitempty"`

	// 网约车详情字段（当order_type=ride时使用）
	PassengerID       string  `json:"passenger_id,omitempty"`
	PassengerName     string  `json:"passenger_name,omitempty"`
	PassengerPhone    string  `json:"passenger_phone,omitempty"`
	PassengerCount    int     `json:"passenger_count,omitempty"`
	RideType          string  `json:"ride_type,omitempty"`
	PickupAddress     string  `json:"pickup_address,omitempty"`
	PickupLatitude    float64 `json:"pickup_latitude"`
	PickupLongitude   float64 `json:"pickup_longitude"`
	PickupLandmark    string  `json:"pickup_landmark,omitempty"`
	DropoffAddress    string  `json:"dropoff_address,omitempty"`
	DropoffLatitude   float64 `json:"dropoff_latitude"`
	DropoffLongitude  float64 `json:"dropoff_longitude"`
	DropoffLandmark   string  `json:"dropoff_landmark,omitempty"`
	EstimatedDistance float64 `json:"estimated_distance,omitempty"`
	EstimatedDuration int     `json:"estimated_duration,omitempty"`
	ActualDistance    float64 `json:"actual_distance,omitempty"`
	ActualDuration    int     `json:"actual_duration,omitempty"`

	// 司机信息
	DriverID     string  `json:"driver_id,omitempty"`
	DriverName   string  `json:"driver_name,omitempty"`
	DriverPhone  string  `json:"driver_phone,omitempty"`
	DriverRating float64 `json:"driver_rating,omitempty"`

	// 车辆信息
	VehicleID    string `json:"vehicle_id,omitempty"`
	LicensePlate string `json:"license_plate,omitempty"`
	VehicleModel string `json:"vehicle_model,omitempty"`
	VehicleColor string `json:"vehicle_color,omitempty"`
	VehiclePhoto string `json:"vehicle_photo,omitempty"`

	// 外卖详情字段（当order_type=delivery时使用，预留）
	RestaurantID      string  `json:"restaurant_id,omitempty"`
	RestaurantName    string  `json:"restaurant_name,omitempty"`
	RestaurantAddress string  `json:"restaurant_address,omitempty"`
	RestaurantPhone   string  `json:"restaurant_phone,omitempty"`
	DeliveryAddress   string  `json:"delivery_address,omitempty"`
	DeliveryLatitude  float64 `json:"delivery_latitude,omitempty"`
	DeliveryLongitude float64 `json:"delivery_longitude,omitempty"`
	CourierID         string  `json:"courier_id,omitempty"`
	CourierName       string  `json:"courier_name,omitempty"`
	CourierPhone      string  `json:"courier_phone,omitempty"`
	CourierRating     float64 `json:"courier_rating,omitempty"`

	// 购物详情字段（当order_type=shopping时使用，预留）
	StoreID         string `json:"store_id,omitempty"`
	StoreName       string `json:"store_name,omitempty"`
	StoreAddress    string `json:"store_address,omitempty"`
	ShippingAddress string `json:"shipping_address,omitempty"`
	TrackingNumber  string `json:"tracking_number,omitempty"`
}
