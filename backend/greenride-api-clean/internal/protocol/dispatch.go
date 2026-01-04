package protocol

// Dispatch 派单记录（输入参数）
type Dispatch struct {
	DispatchID       string  `json:"dispatch_id"`
	OrderID          string  `json:"order_id"`
	DriverID         string  `json:"driver_id"`
	Round            int     `json:"round"`
	RoundSeq         int     `json:"round_seq"`
	DispatchedAt     int64   `json:"dispatched_at"`
	ExpiredAt        int64   `json:"expired_at"`
	Status           string  `json:"status"`
	RespondedAt      int64   `json:"responded_at"`
	RejectReason     string  `json:"reject_reason"`
	RejectReasonType string  `json:"reject_reason_type"`
	DriverDistance   float64 `json:"driver_distance"`
	DriverLatitude   float64 `json:"driver_latitude"`
	DriverLongitude  float64 `json:"driver_longitude"`
	StrategyConfig   string  `json:"strategy_config"`
	ScheduledTime    int64   `json:"scheduled_time"`
	Level            string  `json:"level"`
	VehicleType      string  `json:"vehicle_type"`
	MaxDrivers       int     `json:"max_drivers"`      // 本轮派单司机数量
	SearchRadius     float64 `json:"search_radius"`    // 搜索半径(公里)
	PriceMultiplier  float64 `json:"price_multiplier"` // 价格倍数
	PassengerCount   int     `json:"passenger_count"`  // 乘客数量
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
}

// DispatchDriver 合格司机
type DispatchDriver struct {
	IsEligible        bool    `json:"is_eligible"`
	DriverID          string  `json:"driver_id"`
	RoundSeq          int     `json:"round_seq"`
	CanAcceptNewOrder bool    `json:"can_accept_new_order"`
	Distance          float64 `json:"distance"`          // 距离(公里)
	EstimatedArrival  int64   `json:"estimated_arrival"` // 预计到达时间
	DispatchOrder     int     `json:"dispatch_order"`    // 派单顺序
	WaitTimeMinutes   int     `json:"wait_time_minutes"`
	DistanceScore     float64 `json:"distance_score"`
	TimeScore         float64 `json:"time_score"`
	QueueScore        float64 `json:"queue_score"`
	RatingScore       float64 `json:"rating_score"`
	ExperienceScore   float64 `json:"experience_score"`
	FinalScore        float64 `json:"final_score"`
	RejectReason      string  `json:"reject_reason,omitempty"`
}

// DispatchResult 派单响应
type DispatchResult struct {
	Success     bool              `json:"success"`
	DispatchID  string            `json:"dispatch_id"`
	DriverCount int               `json:"driver_count"`
	Drivers     []*DispatchDriver `json:"selected_drivers"`
	Message     string            `json:"message"`
}
