package protocol

// =============================================================================
// 通用订单管理请求和响应结构体
// =============================================================================

// OrderIDRequest 订单ID请求结构体
type OrderIDRequest struct {
	UserID  string `json:"user_id"`                     // 用户ID
	OrderID string `json:"order_id" binding:"required"` // 订单ID
}

// RatingIDRequest 评价ID请求结构体
type RatingIDRequest struct {
	RatingID int64  `json:"rating_id" binding:"required"` // 评价ID
	UserID   string `json:"user_id" binding:"required"`   // 用户ID
}

// OrderAcceptRequest 接受订单请求结构体
type OrderAcceptRequest struct {
	UserID  string `json:"user_id" binding:"required"`  // 用户ID（司机ID）
	OrderID string `json:"order_id" binding:"required"` // 订单ID
}

type OrderActionRequest struct {
	DispatchId   string  `json:"dispatch_id"`                  // 派单ID
	UserID       string  `json:"user_id"`                      // 用户ID
	OrderID      string  `json:"order_id"`                     // 订单ID
	ActionType   string  `json:"action_type"`                  // 动作类型
	Latitude     float64 `json:"latitude" binding:"required"`  // 到达纬度
	Longitude    float64 `json:"longitude" binding:"required"` // 到达经度
	RejectReason string  `json:"reject_reason,omitempty"`      // 拒绝原因（枚举值或自定义文本）
}

type UserOnlineRequest struct {
	UserID    string  `json:"user_id"`
	VehicleID string  `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`  // 到达纬度
	Longitude float64 `json:"longitude" binding:"required"` // 到达经度
}

// OrderStatusUpdateRequest 订单状态更新请求结构体
type OrderStatusUpdateRequest struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
	Status  string `json:"status" binding:"required"`   // 状态
}

// OrderCancelAPIRequest 取消订单请求结构体
type OrderCancelAPIRequest struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
	Reason  string `json:"reason,omitempty"`            // 取消原因
}

// PaymentConfirmRequest 确认收款请求结构体
type PaymentConfirmRequest struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
}

// OrderRatingCreateRequest 创建评价请求结构体
type OrderRatingCreateRequest struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
	Rating  int    `json:"rating" binding:"required"`   // 评分
	Comment string `json:"comment,omitempty"`           // 评价内容
}

// RatingReplyRequest 回复评价请求结构体
type RatingReplyRequest struct {
	RatingID string `json:"rating_id" binding:"required"` // 评价ID
	Reply    string `json:"reply" binding:"required"`     // 回复内容
}

// =============================================================================
// 订单预估和创建请求响应结构体
// =============================================================================

// EstimateRequest 内部价格计算请求 - 扩展了更多内部字段
type EstimateRequest struct {
	// 基础信息
	UserID          string `json:"user_id"`
	SessionID       string `json:"session_id"`
	VehicleCategory string `json:"vehicle_category"`     // sedan, suv, mpv, van, hatchback
	VehicleLevel    string `json:"vehicle_level"`        // economy, comfort, premium, luxury
	OrderType       string `json:"order_type,omitempty"` // 订单类型，默认为网约车
	UserCategory    string `json:"user_category"`        // 用户类别

	// 行程信息
	PassengerCount    int     `json:"passenger_count"`
	PickupLatitude    float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude   float64 `json:"pickup_longitude" binding:"required"`
	PickupAddress     string  `json:"pickup_address,omitempty"`
	PickupLandmark    string  `json:"pickup_landmark,omitempty"`
	DropoffLatitude   float64 `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude  float64 `json:"dropoff_longitude" binding:"required"`
	DropoffAddress    string  `json:"dropoff_address,omitempty"`
	DropoffLandmark   string  `json:"dropoff_landmark,omitempty"`
	EstimatedDistance float64 `json:"estimated_distance"`
	EstimatedDuration int     `json:"estimated_duration"`

	// 价格相关
	Currency  string  `json:"currency"`   // 币种
	BasePrice float64 `json:"base_price"` // 基础价格（内部计算用）

	// 上下文信息
	ScheduledAt int64 `json:"scheduled_at,omitempty"` // 预定时间(时间戳毫秒)，0表示即时
	RequestedAt int64 `json:"requested_at,omitempty"` // 请求时间(时间戳毫秒)，0表示当前时间

	ServiceArea   string   `json:"service_area"`
	PromoCodes    []string `json:"promo_codes,omitempty"`
	PaymentMethod string   `json:"payment_method,omitempty"`

	// 客户端信息
	Platform   string `json:"platform,omitempty"`
	AppVersion string `json:"app_version,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	RequestIP  string `json:"request_ip,omitempty"`

	// 设置选项
	SnapshotDuration int64 `json:"snapshot_duration"` // 快照有效期(分钟)，默认30分钟
}

// CreateOrderRequest 创建订单请求（扁平化结构）
type CreateOrderRequest struct {
	UserID string `json:"user_id"` // 用户ID

	// 价格信息
	PriceID string `json:"price_id,omitempty" binding:"required"` // 价格ID，来自预估接口

	// 扩展信息
	Notes string `json:"notes,omitempty"`

	// 手动指定司机（可选）
	// - 为空：走自动派单
	// - 不为空：将订单预分配给该司机，并仅向该司机发送派单
	ProviderID string `json:"provider_id,omitempty"`
}

// =============================================================================
// 订单查询和管理请求结构体
// =============================================================================

// GetOrdersRequest 获取订单列表请求
type GetVehiclesRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted in_progress completed cancelled"`
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID      string `json:"order_id" binding:"required"`
	UserID       string `json:"user_id"`                  // 内部设置
	Reason       string `json:"reason,omitempty"`         // 取消原因 (旧版兼容, free-form text)
	ReasonKey    string `json:"reason_key,omitempty"`     // 预定义取消原因 key (新版)
	CustomReason string `json:"custom_reason,omitempty"`  // 自定义原因 (当 reason_key="other" 时使用)
}

// =============================================================================
// 订单评价相关请求结构体
// =============================================================================
// CreateOrderRatingRequest 创建订单评价请求（包含订单ID）
type CreateOrderRatingRequest struct {
	OrderID     string   `json:"order_id" binding:"required"`
	UserID      string   `json:"user_id"`                               // 评价者ID
	Rating      float64  `json:"rating" binding:"required,min=1,max=5"` // 评分 1-5
	Comment     string   `json:"comment"`                               // 评价内容
	Tags        []string `json:"tags,omitempty"`                        // 评价标签(JSON格式)
	IsAnonymous bool     `json:"is_anonymous,omitempty"`                // 是否匿名评价
}

// UpdateOrderRatingRequest 更新订单评价API请求（包含评价ID）
type UpdateOrderRatingRequest struct {
	UserID   string  `json:"user_id"`                               // 评价者ID
	RatingID int64   `json:"rating_id" binding:"required"`          // 评价ID
	Rating   float64 `json:"rating" binding:"required,min=1,max=5"` // 评分 1-5
	Comment  string  `json:"comment,omitempty"`                     // 评价内容
}

// ReplyToRatingRequest 回复评价请求
type ReplyToRatingRequest struct {
	UserID   string `json:"user_id"`                      // 评价者ID
	RatingID string `json:"rating_id" binding:"required"` // 评价ID
	Reply    string `json:"reply" binding:"required"`     // 回复内容
}

// GetNearbyOrdersRequest 获取附近订单请求
type GetNearbyOrdersRequest struct {
	Latitude    float64 `form:"latitude"`
	Longitude   float64 `form:"longitude" `
	Radius      float64 `form:"radius,default=10"`    // 半径（公里）
	OrderType   string  `form:"order_type,omitempty"` // 订单类型过滤
	Limit       int     `form:"limit,default=10"`     // 数量限制
	RequesterID string  `json:"-"`                    // set from auth context
}

// =============================================================================
// 价格快照相关请求和响应结构体
// =============================================================================

// OrderPaymentRequest 付款请求
type OrderPaymentRequest struct {
	Language      string `json:"-"`
	OrderID       string `json:"order_id" binding:"required"`
	UserID        string `json:"user_id"`                  // 内部设置
	PaymentMethod string `json:"payment_method,omitempty"` // 内部设置
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	AccountNo     string `json:"account_no,omitempty"`
	AccountName   string `json:"account_name,omitempty"`
}

type OrderPaymentResult struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	ExpiredAt   int64  `json:"expired_at,omitempty"`
}

// GetOrdersResponse 订单列表响应
type GetOrdersResponse struct {
	Orders []*Order `json:"orders"`
	Total  int64    `json:"total"`
}

// GetNearbyOrdersResponse 附近订单响应
type GetNearbyOrdersResponse struct {
	Orders []*Order `json:"orders"`
	Count  int      `json:"count"`
}

// UpdateLocationRequest 位置更新请求
type UpdateLocationRequest struct {
	UserID       string  `json:"user_id"`                                      // 内部设置
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`   // 纬度
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"` // 经度
	Heading      float64 `json:"heading,omitempty"`                            // 行驶方向 (0-360度)
	Speed        float64 `json:"speed,omitempty"`                              // 速度 km/h
	Accuracy     float64 `json:"accuracy,omitempty"`                           // GPS精度 (米)
	Altitude     float64 `json:"altitude,omitempty"`                           // 海拔 (米)
	OnlineStatus string  `json:"online_status,omitempty"`                      // online, busy, offline
	UpdatedAt    int64   `json:"updated_at,omitempty"`                         // 时间戳，毫秒
}

type UserPromotionsRequest struct {
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	UserID    string `json:"user_id,omitempty" `
	UserType  string `json:"user_type,omitempty"`  // 用户类型，自动识别
	Status    string `json:"status,omitempty"`     // 优惠券状态过滤
	StartDate int64  `json:"start_date,omitempty"` // 开始日期过滤 (时间戳毫秒)
	EndDate   int64  `json:"end_date,omitempty"`   // 结束日期过滤 (时间戳毫秒)
}

// =============================================================================
// Nearby Drivers Request/Response (for passengers to find drivers)
// =============================================================================

// GetNearbyDriversRequest 获取附近司机请求
type GetNearbyDriversRequest struct {
	Latitude  float64 `form:"latitude" binding:"required"`  // 乘客当前纬度
	Longitude float64 `form:"longitude" binding:"required"` // 乘客当前经度
	RadiusKm  float64 `form:"radius_km"`                    // 搜索半径（公里），默认5km
	Limit     int     `form:"limit"`                        // 返回数量限制，默认20
	EtaMode   string  `form:"eta_mode"`                     // ETA模式：rough|accurate|none
}

// NearbyDriver 附近司机信息
type NearbyDriver struct {
	DriverID        string  `json:"driver_id"`
	Name            string  `json:"name"`
	PhotoURL        string  `json:"photo_url,omitempty"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	DistanceKm      float64 `json:"distance_km"`
	ETAMinutes      int     `json:"eta_minutes"`
	VehicleType     string  `json:"vehicle_type"`
	VehicleCategory string  `json:"vehicle_category"` // sedan, suv, moto
	VehicleBrand    string  `json:"vehicle_brand"`
	VehicleModel    string  `json:"vehicle_model"`
	PlateNumber     string  `json:"plate_number"`
	VehicleColor    string  `json:"vehicle_color"`
	Rating          float64 `json:"rating"`
	TotalRides      int     `json:"total_rides"`
	IsOnline        bool    `json:"is_online"`
	IsBusy          bool    `json:"is_busy"`   // 是否正在接单
	Heading         float64 `json:"heading"`   // 行驶方向角度 (0-360)
	Phone           string  `json:"phone,omitempty"` // 联系电话
}

// GetNearbyDriversResponse 附近司机响应
type GetNearbyDriversResponse struct {
	Drivers []*NearbyDriver `json:"drivers"`
	Count   int             `json:"count"`
}

// =============================================================================
// Order Contact (Call Permission) Request/Response
// =============================================================================

// OrderContactRequest request to get contact info for calling
type OrderContactRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  string `json:"-"` // set from auth context
}

// OrderContactResponse returns whether calling is allowed and the phone number
type OrderContactResponse struct {
	Allowed bool   `json:"allowed"`
	Phone   string `json:"phone,omitempty"`
	Name    string `json:"name,omitempty"`
}

// =============================================================================
// Order ETA Request/Response
// =============================================================================

// OrderETARequest request for live ETA of an active order
type OrderETARequest struct {
	OrderID string `json:"order_id" binding:"required"`
	UserID  string `json:"-"` // set from auth context
}

// OrderETAResponse returns live ETA and driver position
type OrderETAResponse struct {
	OrderID         string  `json:"order_id"`
	ETAMinutes      int     `json:"eta_minutes"`
	DistanceKm      float64 `json:"distance_km"`
	DriverLatitude  float64 `json:"driver_latitude,omitempty"`
	DriverLongitude float64 `json:"driver_longitude,omitempty"`
	PickupLatitude  float64 `json:"pickup_latitude"`
	PickupLongitude float64 `json:"pickup_longitude"`
	Mode            string  `json:"mode"` // "rough" or "accurate"
	UpdatedAt       int64   `json:"updated_at"`
}
