package protocol

// =============================================================================
// 车辆管理协议结构体
// =============================================================================

// Vehicle 车辆信息响应结构体
type Vehicle struct {
	// 基础信息
	VehicleID   string `json:"vehicle_id"`
	DriverID    string `json:"driver_id,omitempty"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	Color       string `json:"color"`
	PlateNumber string `json:"plate_number"`
	VIN         string `json:"vin,omitempty"`

	// 车辆类型和配置
	TypeID       string `json:"type_id"`  // 关联VehicleType表的ID
	Category     string `json:"category"` // sedan, suv, mpv, van, hatchback
	Level        string `json:"level"`    // economy, comfort, premium, luxury
	SeatCapacity int    `json:"seat_capacity"`
	FuelType     string `json:"fuel_type"`    // petrol, diesel, electric, hybrid
	Transmission string `json:"transmission"` // manual, automatic
	EngineSize   string `json:"engine_size,omitempty"`

	// 状态信息
	Status       string `json:"status"`        // active, inactive, maintenance, retired
	VerifyStatus string `json:"verify_status"` // unverified, pending, active, inactive, maintenance, suspended, banned, retired

	// 注册和保险信息
	RegistrationNumber    string  `json:"registration_number,omitempty"`
	RegistrationExpiry    *string `json:"registration_expiry,omitempty"`
	InsuranceCompany      string  `json:"insurance_company,omitempty"`
	InsurancePolicyNumber string  `json:"insurance_policy_number,omitempty"`
	InsuranceExpiry       *string `json:"insurance_expiry,omitempty"`

	// 位置信息
	CurrentLatitude   *float64 `json:"current_latitude,omitempty"`
	CurrentLongitude  *float64 `json:"current_longitude,omitempty"`
	LocationUpdatedAt *int64   `json:"location_updated_at,omitempty"`

	// 维护信息
	LastServiceDate string `json:"last_service_date,omitempty"`
	NextServiceDue  string `json:"next_service_due,omitempty"`
	TotalMileage    int    `json:"total_mileage"`
	ServiceMileage  int    `json:"service_mileage,omitempty"`

	// 车辆特性
	HasAirConditioner bool `json:"has_air_conditioner"`
	HasGPS            bool `json:"has_gps"`
	HasWiFi           bool `json:"has_wifi"`
	HasCharger        bool `json:"has_charger"`
	HasBluetooth      bool `json:"has_bluetooth"`

	// 文档和图片
	Photos    []string `json:"photos,omitempty"`
	Documents []string `json:"documents,omitempty"`

	// 评价和统计
	Rating        float64 `json:"rating"`
	TotalRides    int     `json:"total_rides"`
	TotalDistance float64 `json:"total_distance"`

	// 验证状态
	DocumentsVerified bool   `json:"documents_verified"`
	InspectionPassed  bool   `json:"inspection_passed"`
	InspectionDate    string `json:"inspection_date,omitempty"`
	InspectionExpiry  string `json:"inspection_expiry,omitempty"`
	VerifiedAt        int64  `json:"verified_at,omitempty"`

	// 使用信息
	LastUsedAt int64  `json:"last_used_at,omitempty"`
	Notes      string `json:"notes,omitempty"`

	// 时间戳
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`

	Driver *User `json:"driver,omitempty"`
}

// =============================================================================
// 客户端车辆请求结构体
// =============================================================================

// UserVehicleRequest 用户车辆查询请求结构体（司机专用）
type UserVehicleRequest struct {
	// 空结构体，保持POST一致性，司机查看分派车辆信息
}

// =============================================================================
// 管理后台车辆请求结构体
// =============================================================================

// VehicleListRequest 车辆列表请求结构体（管理后台）
type VehicleListRequest struct {
	Keyword      string `json:"keyword,omitempty"`       // 搜索关键字
	Page         int    `json:"page,omitempty"`          // 页码，默认1
	Limit        int    `json:"limit,omitempty"`         // 每页数量，默认10
	TypeID       string `json:"type_id,omitempty"`       // 车辆类型ID
	Category     string `json:"category,omitempty"`      // 车辆分类
	Level        string `json:"level,omitempty"`         // 服务级别
	Status       string `json:"status,omitempty"`        // 车辆状态
	DriverID     string `json:"driver_id,omitempty"`     // 司机ID
	VerifyStatus string `json:"verify_status,omitempty"` // 验证状态
	YearFrom     *int   `json:"year_from,omitempty"`     // 年份范围开始
	YearTo       *int   `json:"year_to,omitempty"`       // 年份范围结束
}

// VehicleDetailRequest 车辆详情请求结构体（管理后台）
type VehicleDetailRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"` // 车辆ID
}

// VehicleCreateRequest 车辆创建请求结构体（管理后台）
type VehicleCreateRequest struct {
	Brand        string `json:"brand" binding:"required"`              // 品牌
	Model        string `json:"model" binding:"required"`              // 型号
	Year         int    `json:"year" binding:"required"`               // 年份
	Color        string `json:"color,omitempty"`                       // 颜色
	PlateNumber  string `json:"plate_number" binding:"required"`       // 车牌号
	VIN          string `json:"vin,omitempty"`                         // 车辆识别号
	Category     string `json:"category,omitempty" binding:"required"` // 车辆分类
	Level        string `json:"level,omitempty" binding:"required"`    // 服务级别
	SeatCapacity int    `json:"seat_capacity,omitempty"`               // 座位数
	FuelType     string `json:"fuel_type,omitempty"`                   // 燃料类型
	Transmission string `json:"transmission,omitempty"`                // 变速箱类型
	EngineSize   string `json:"engine_size,omitempty"`                 // 引擎大小

	// 注册和保险信息
	RegistrationNumber    string  `json:"registration_number,omitempty"`     // 注册号
	RegistrationExpiry    *string `json:"registration_expiry,omitempty"`     // 注册到期日
	InsuranceCompany      string  `json:"insurance_company,omitempty"`       // 保险公司
	InsurancePolicyNumber string  `json:"insurance_policy_number,omitempty"` // 保险单号
	InsuranceExpiry       *string `json:"insurance_expiry,omitempty"`        // 保险到期日

	// 车辆特性
	HasAirConditioner *bool `json:"has_air_conditioner,omitempty"` // 是否有空调
	HasGPS            *bool `json:"has_gps,omitempty"`             // 是否有GPS
	HasWiFi           *bool `json:"has_wifi,omitempty"`            // 是否有WiFi
	HasCharger        *bool `json:"has_charger,omitempty"`         // 是否有充电器
	HasBluetooth      *bool `json:"has_bluetooth,omitempty"`       // 是否有蓝牙

	// 文档和图片
	Photos    []string `json:"photos,omitempty"`    // 车辆照片
	Documents []string `json:"documents,omitempty"` // 车辆文档

	// 其他信息
	Notes string `json:"notes,omitempty"` // 备注
}

// VehicleUpdateRequest 车辆更新请求结构体（管理后台）
type VehicleUpdateRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"` // 车辆ID

	// 司机分配（可选）
	// - nil: 不修改当前分配
	// - "":  取消分配（写入 NULL）
	// - 其他: 分配到该司机（会自动解除该司机在其他车辆上的绑定）
	DriverID *string `json:"driver_id,omitempty"` // 司机ID

	// 基础信息（可选）
	Brand       *string `json:"brand,omitempty"`        // 品牌
	Model       *string `json:"model,omitempty"`        // 型号
	Year        *int    `json:"year,omitempty"`         // 年份
	Color       *string `json:"color,omitempty"`        // 颜色
	PlateNumber *string `json:"plate_number,omitempty"` // 车牌号
	VIN         *string `json:"vin,omitempty"`          // 车辆识别号

	// 车辆类型和配置（可选）
	TypeID       *string `json:"type_id,omitempty"`       // 关联VehicleType表的ID
	Category     *string `json:"category,omitempty"`      // 车辆分类
	Level        *string `json:"level,omitempty"`         // 服务级别
	SeatCapacity *int    `json:"seat_capacity,omitempty"` // 座位数
	FuelType     *string `json:"fuel_type,omitempty"`     // 燃料类型
	Transmission *string `json:"transmission,omitempty"`  // 变速箱类型
	EngineSize   *string `json:"engine_size,omitempty"`   // 引擎大小

	// 注册和保险信息（可选）
	RegistrationNumber    *string `json:"registration_number,omitempty"`     // 注册号
	RegistrationExpiry    *string `json:"registration_expiry,omitempty"`     // 注册到期日
	InsuranceCompany      *string `json:"insurance_company,omitempty"`       // 保险公司
	InsurancePolicyNumber *string `json:"insurance_policy_number,omitempty"` // 保险单号
	InsuranceExpiry       *string `json:"insurance_expiry,omitempty"`        // 保险到期日

	// 维护信息（可选）
	LastServiceDate *string `json:"last_service_date,omitempty"` // 最后保养日期
	NextServiceDue  *string `json:"next_service_due,omitempty"`  // 下次保养日期
	TotalMileage    *int    `json:"total_mileage,omitempty"`     // 总里程
	ServiceMileage  *int    `json:"service_mileage,omitempty"`   // 保养里程

	// 车辆特性（可选）
	HasAirConditioner *bool `json:"has_air_conditioner,omitempty"` // 是否有空调
	HasGPS            *bool `json:"has_gps,omitempty"`             // 是否有GPS
	HasWiFi           *bool `json:"has_wifi,omitempty"`            // 是否有WiFi
	HasCharger        *bool `json:"has_charger,omitempty"`         // 是否有充电器
	HasBluetooth      *bool `json:"has_bluetooth,omitempty"`       // 是否有蓝牙

	// 文档和图片（可选）
	Photos    []string `json:"photos,omitempty"`    // 车辆照片
	Documents []string `json:"documents,omitempty"` // 车辆文档

	// 其他信息（可选）
	Notes *string `json:"notes,omitempty"` // 备注
}

// VehicleStatusUpdateRequest 车辆状态更新请求结构体（管理后台）
type VehicleStatusUpdateRequest struct {
	VehicleID    string  `json:"vehicle_id" binding:"required"` // 车辆ID
	Status       *string `json:"status,omitempty"`              // 车辆状态
	VerifyStatus *string `json:"verify_status,omitempty"`       // 验证状态
}

// VehicleDeleteRequest 车辆删除请求结构体（管理后台）
type VehicleDeleteRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"` // 车辆ID
	Reason    string `json:"reason,omitempty"`              // 删除原因
}

// VehicleVerifyRequest 车辆验证请求结构体（管理后台）
type VehicleVerifyRequest struct {
	VehicleID         string `json:"vehicle_id" binding:"required"` // 车辆ID
	DocumentsVerified *bool  `json:"documents_verified,omitempty"`  // 文档是否验证
	InspectionPassed  *bool  `json:"inspection_passed,omitempty"`   // 检验是否通过
	VerifiedBy        string `json:"verified_by,omitempty"`         // 验证者ID
	Notes             string `json:"notes,omitempty"`               // 验证备注
}

// VehicleAssignRequest 车辆分派请求结构体（管理后台）
type VehicleAssignRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"` // 车辆ID
	DriverID  string `json:"driver_id" binding:"required"`  // 司机ID
	Notes     string `json:"notes,omitempty"`               // 分派备注
}

// VehicleUnassignRequest 车辆取消分派请求结构体（管理后台）
type VehicleUnassignRequest struct {
	VehicleID string `json:"vehicle_id" binding:"required"` // 车辆ID
	Reason    string `json:"reason,omitempty"`              // 取消分派原因
}
