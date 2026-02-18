package config

// DispatchConfig 派单配置
type DispatchConfig struct {
	Enabled         bool                   `mapstructure:"enabled"`          // 启用新派单系统
	DriverSelection *DriverSelectionConfig `mapstructure:"driver_selection"` // 司机筛选配置
	TimeWindow      *TimeWindowConfig      `mapstructure:"time_window"`      // 时间窗口配置
	Scoring         *ScoringConfig         `mapstructure:"scoring"`          // 评分配置
	Rounds          *RoundsConfig          `mapstructure:"rounds"`           // 派单轮次配置

	MaxRounds                 int     `mapstructure:"max_rounds" json:"max_rounds"`                             // 最大派单轮次
	TimeoutSeconds            int     `mapstructure:"timeout_seconds" json:"timeout_seconds"`                   // 司机响应超时时间
	MaxDistance               float64 `mapstructure:"max_distance" json:"max_distance"`                         // 最大距离(公里)
	MaxWaitTime               int     `mapstructure:"max_wait_time" json:"max_wait_time"`                       // 最大等待时间(分钟)
	MaxIdleTime               int     `mapstructure:"max_idle_time" json:"max_idle_time"`                       // 最大空闲时间(分钟)
	MaxIdleDistance           float64 `mapstructure:"max_idle_distance" json:"max_idle_distance"`               // 最大空闲距离(公里)
	MaxNextOrderDelayTime     int     `mapstructure:"max_next_order_delay_time" json:"max_next_order_delay_time"`     // 最大下单延迟时间(分钟)
	MaxNextOrderDelayDistance float64 `mapstructure:"max_next_order_delay_distance" json:"max_next_order_delay_distance"` // 最大下单延迟距离(公里)

	// 评分权重
	DistanceWeight   float64 `mapstructure:"distance_weight" json:"distance_weight"`
	TimeWeight       float64 `mapstructure:"time_weight" json:"time_weight"`
	RatingWeight     float64 `mapstructure:"rating_weight" json:"rating_weight"`
	QueueWeight      float64 `mapstructure:"queue_weight" json:"queue_weight"`
	ExperienceWeight float64 `mapstructure:"experience_weight" json:"experience_weight"`
}

// RoundStrategy 轮次策略
type RoundStrategy struct {
	MaxDrivers            int     `json:"max_drivers"`
	SearchRadius          float64 `json:"search_radius"`
	PriceMultiplier       float64 `json:"price_multiplier"`
	MinRatingScore        float64 `json:"min_rating_score"`
	MinAcceptanceRate     float64 `json:"min_acceptance_rate"`
	MaxConsecutiveRejects int     `json:"max_consecutive_rejects"`
}

// TimeWindowConfig 时间窗口配置
type TimeWindowConfig struct {
	MaxWaitTime       int  `json:"max_wait_time"`       // 最大等待时间(分钟)
	LocationTolerance int  `json:"location_tolerance"`  // 位置容差(米)
	RouteCheckEnabled bool `json:"route_check_enabled"` // 是否启用路线检查
}

func (d *TimeWindowConfig) Validate() {
	// 设置 time_window 默认值
	if d.MaxWaitTime == 0 {
		d.MaxWaitTime = 60 // 默认60分钟
	}
	if d.LocationTolerance == 0 {
		d.LocationTolerance = 2000 // 默认2公里
	}
}

// DispatchStatus 派单状态
type DispatchStatus struct {
	OrderID          string            `json:"order_id"`
	DispatchID       string            `json:"dispatch_id"`
	Round            int               `json:"round"`
	Status           string            `json:"status"` // pending, completed, failed, timeout
	DriverCount      int               `json:"driver_count"`
	Responses        map[string]string `json:"responses"` // driver_id -> response (accept/reject/timeout)
	SelectedDriverID string            `json:"selected_driver_id,omitempty"`
	CreatedAt        int64             `json:"created_at"`
	UpdatedAt        int64             `json:"updated_at"`
}

// DriverTimeWindow 司机时间窗口
type DriverTimeWindow struct {
	CanAcceptNewOrder     bool    `json:"can_accept_new_order"`    // 是否可以接新单
	CurrentOrderID        string  `json:"current_order_id"`        // 当前订单ID
	EstimatedCompleteTime int64   `json:"estimated_complete_time"` // 预计完成时间(时间戳)
	WaitTimeMinutes       int     `json:"wait_time_minutes"`       // 等待时间(分钟)
	RouteMatchScore       float64 `json:"route_match_score"`       // 路线匹配度(0-1)
	DistanceToDestination float64 `json:"distance_to_destination"` // 到目标的距离(米)
}

// DriverSelectionConfig 司机筛选配置
type DriverSelectionConfig struct {
	UseGeolocation bool `mapstructure:"use_geolocation"` // 是否使用地理位置筛选
}

// ScoringConfig 评分配置
type ScoringConfig struct {
	NormalizeScores bool           `mapstructure:"normalize_scores"` // 是否标准化评分
	Factors         ScoringFactors `mapstructure:"factors"`          // 评分因子
}

func (d *ScoringConfig) Validate() {
	// 设置 scoring.factors 默认值
	if d.Factors.Rating == 0 {
		d.Factors.Rating = 0.4
	}
	if d.Factors.AcceptanceRate == 0 {
		d.Factors.AcceptanceRate = 0.3
	}
	if d.Factors.Distance == 0 {
		d.Factors.Distance = 0.2
	}
	if d.Factors.ExperienceLevel == 0 {
		d.Factors.ExperienceLevel = 0.1
	}
}

// ScoringFactors 评分因子权重
type ScoringFactors struct {
	Rating          float64 `mapstructure:"rating"`           // 司机评分权重
	AcceptanceRate  float64 `mapstructure:"acceptance_rate"`  // 接单率权重
	Distance        float64 `mapstructure:"distance"`         // 距离权重
	ExperienceLevel float64 `mapstructure:"experience_level"` // 经验级别权重
}

// RoundsConfig 派单轮次配置
type RoundsConfig struct {
	MaxRounds       int             `mapstructure:"max_rounds"`        // 最大轮次(1)
	DriversPerRound int             `mapstructure:"drivers_per_round"` // 每轮司机数量(0=全部)
	RoundStrategys  []RoundStrategy `mapstructure:"round_strategys"`   // 轮次策略
}

func (d *RoundsConfig) Validate() {
	// 设置 rounds 默认值
	if d.MaxRounds == 0 {
		d.MaxRounds = 1 // 默认最大1轮
	}
	if d.DriversPerRound == 0 {
		d.DriversPerRound = 0 // 默认全部司机(0)
	}
}

// Validate 验证并设置派单配置默认值
func (d *DispatchConfig) Validate() {
	if d == nil {
		return
	}
	if d.TimeoutSeconds == 0 {
		d.TimeoutSeconds = 5 * 60 // 默认5分钟
	}
	if d.DriverSelection == nil {
		d.DriverSelection = &DriverSelectionConfig{
			UseGeolocation: false,
		}
	}

	if d.TimeWindow == nil {
		d.TimeWindow = &TimeWindowConfig{
			MaxWaitTime:       60,   // 默认60分钟
			LocationTolerance: 2000, // 默认2公里
			RouteCheckEnabled: true,
		}
	}
	d.TimeWindow.Validate()

	if d.Scoring == nil {
		d.Scoring = &ScoringConfig{
			NormalizeScores: true,
			Factors: ScoringFactors{
				Rating:          0.4,
				AcceptanceRate:  0.3,
				Distance:        0.2,
				ExperienceLevel: 0.1,
			},
		}
	}
	d.Scoring.Validate()

	if d.Rounds == nil {
		d.Rounds = &RoundsConfig{
			MaxRounds:       1,
			DriversPerRound: 0,
			RoundStrategys: []RoundStrategy{
				{
					MaxDrivers:            0,
					SearchRadius:          0.0,
					PriceMultiplier:       1.0,
					MinRatingScore:        4.5,
					MinAcceptanceRate:     0.7,
					MaxConsecutiveRejects: 3,
				},
			},
		}
	}
	d.Rounds.Validate()
}
