package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"math"
	"strings"
)

// ServiceArea 服务区域表 - 地理位置服务覆盖区域管理
type ServiceArea struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ServiceAreaID string `json:"service_area_id" gorm:"column:service_area_id;type:varchar(64);uniqueIndex"`
	Salt          string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*ServiceAreaValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type ServiceAreaValues struct {
	// 基本信息
	AreaName    *string `json:"area_name" gorm:"column:area_name;type:varchar(255);index"` // 区域名称
	DisplayName *string `json:"display_name" gorm:"column:display_name;type:varchar(255)"` // 显示名称
	Description *string `json:"description" gorm:"column:description;type:text"`           // 区域描述
	AreaType    *string `json:"area_type" gorm:"column:area_type;type:varchar(50);index"`  // city, suburb, district, zone, airport, special

	// 地理位置
	CenterLat *float64 `json:"center_lat" gorm:"column:center_lat;type:decimal(10,8)"` // 中心纬度
	CenterLng *float64 `json:"center_lng" gorm:"column:center_lng;type:decimal(11,8)"` // 中心经度
	Radius    *float64 `json:"radius" gorm:"column:radius;type:decimal(8,2)"`          // 服务半径(公里)

	// 边界定义
	Boundaries *string `json:"boundaries" gorm:"column:boundaries;type:json"` // JSON数组：边界坐标点
	Polygon    *string `json:"polygon" gorm:"column:polygon;type:json"`       // JSON数组：多边形坐标

	// 行政区划
	Country  *string `json:"country" gorm:"column:country;type:varchar(100);index"`          // 国家
	Province *string `json:"province" gorm:"column:province;type:varchar(100);index"`        // 省/州
	City     *string `json:"city" gorm:"column:city;type:varchar(100);index"`                // 城市
	District *string `json:"district" gorm:"column:district;type:varchar(100);index"`        // 区/县
	ZipCode  *string `json:"zip_code" gorm:"column:zip_code;type:varchar(20);index"`         // 邮政编码
	TimeZone *string `json:"timezone" gorm:"column:timezone;type:varchar(50);default:'UTC'"` // 时区

	// 层级关系
	ParentAreaID *string `json:"parent_area_id" gorm:"column:parent_area_id;type:varchar(64);index"` // 父级区域ID
	Level        *int    `json:"level" gorm:"column:level;type:int;default:1"`                       // 层级级别
	Priority     *int    `json:"priority" gorm:"column:priority;type:int;default:100"`               // 优先级

	// 服务配置
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`           // active, inactive, maintenance, planned
	IsActive     *bool   `json:"is_active" gorm:"column:is_active;default:true"`                                // 是否激活
	ServiceLevel *string `json:"service_level" gorm:"column:service_level;type:varchar(50);default:'standard'"` // basic, standard, premium, full

	// 运营时间
	OperatingHours *string `json:"operating_hours" gorm:"column:operating_hours;type:json"` // JSON对象：运营时间配置
	Is24Hours      *bool   `json:"is_24hours" gorm:"column:is_24hours;default:false"`       // 是否24小时服务

	// 支持的服务类型
	SupportedVehicleTypes *string `json:"supported_vehicle_types" gorm:"column:supported_vehicle_types;type:json"` // JSON数组：支持的车型
	SupportedServices     *string `json:"supported_services" gorm:"column:supported_services;type:json"`           // JSON数组：支持的服务类型

	// 定价配置
	PricingTier    *string  `json:"pricing_tier" gorm:"column:pricing_tier;type:varchar(30);default:'standard'"` // economy, standard, premium
	SurgeEnabled   *bool    `json:"surge_enabled" gorm:"column:surge_enabled;default:true"`                      // 是否启用涌潮定价
	BaseFareRate   *float64 `json:"base_fare_rate" gorm:"column:base_fare_rate;type:decimal(5,2);default:1.00"`  // 基础价格倍率
	SurgeThreshold *float64 `json:"surge_threshold" gorm:"column:surge_threshold;type:decimal(5,2);default:2.0"` // 涌潮触发阈值

	// 需求供给数据
	DemandLevel  *string  `json:"demand_level" gorm:"column:demand_level;type:varchar(20);default:'medium'"` // low, medium, high
	SupplyLevel  *string  `json:"supply_level" gorm:"column:supply_level;type:varchar(20);default:'medium'"` // low, medium, high
	CurrentSurge *float64 `json:"current_surge" gorm:"column:current_surge;type:decimal(4,2);default:1.0"`   // 当前涌潮倍数

	// 统计数据
	TotalRides      *int     `json:"total_rides" gorm:"column:total_rides;type:int;default:0"`                   // 总订单数
	TotalDrivers    *int     `json:"total_drivers" gorm:"column:total_drivers;type:int;default:0"`               // 总司机数
	ActiveDrivers   *int     `json:"active_drivers" gorm:"column:active_drivers;type:int;default:0"`             // 活跃司机数
	TotalUsers      *int     `json:"total_users" gorm:"column:total_users;type:int;default:0"`                   // 总用户数
	AverageWaitTime *int     `json:"average_wait_time" gorm:"column:average_wait_time;type:int;default:0"`       // 平均等待时间(秒)
	AverageRating   *float64 `json:"average_rating" gorm:"column:average_rating;type:decimal(3,2);default:0.00"` // 平均评分

	// 交通和基础设施
	TrafficLevel        *string `json:"traffic_level" gorm:"column:traffic_level;type:varchar(20);default:'medium'"`               // low, medium, high, heavy
	ParkingAvailability *string `json:"parking_availability" gorm:"column:parking_availability;type:varchar(20);default:'medium'"` // limited, medium, abundant
	PublicTransport     *bool   `json:"public_transport" gorm:"column:public_transport;default:false"`                             // 是否有公共交通
	AirportNearby       *bool   `json:"airport_nearby" gorm:"column:airport_nearby;default:false"`                                 // 是否靠近机场

	// 特殊配置
	IsAirportZone      *bool    `json:"is_airport_zone" gorm:"column:is_airport_zone;default:false"`           // 是否机场区域
	IsBusinessDistrict *bool    `json:"is_business_district" gorm:"column:is_business_district;default:false"` // 是否商务区
	IsResidential      *bool    `json:"is_residential" gorm:"column:is_residential;default:false"`             // 是否住宅区
	IsTouristArea      *bool    `json:"is_tourist_area" gorm:"column:is_tourist_area;default:false"`           // 是否旅游区
	IsHighCrimeArea    *bool    `json:"is_high_crime_area" gorm:"column:is_high_crime_area;default:false"`     // 是否高犯罪率区域
	AirportSurcharge   *float64 `json:"airport_surcharge" gorm:"column:airport_surcharge;type:decimal(8,2)"`   // 机场附加费

	// 天气和环境
	WeatherSensitive *bool   `json:"weather_sensitive" gorm:"column:weather_sensitive;default:false"` // 是否受天气影响
	AltitudeLevel    *int    `json:"altitude_level" gorm:"column:altitude_level;type:int"`            // 海拔高度(米)
	ClimateZone      *string `json:"climate_zone" gorm:"column:climate_zone;type:varchar(50)"`        // 气候带

	// 人口和经济数据
	PopulationDensity *string `json:"population_density" gorm:"column:population_density;type:varchar(20)"` // low, medium, high, very_high
	IncomeLevel       *string `json:"income_level" gorm:"column:income_level;type:varchar(20)"`             // low, medium, high
	BusinessActivity  *string `json:"business_activity" gorm:"column:business_activity;type:varchar(20)"`   // low, medium, high

	// 限制和规则
	VehicleRestrictions *string `json:"vehicle_restrictions" gorm:"column:vehicle_restrictions;type:json"` // JSON数组：车辆限制
	DriverRequirements  *string `json:"driver_requirements" gorm:"column:driver_requirements;type:json"`   // JSON数组：司机要求
	SpecialRules        *string `json:"special_rules" gorm:"column:special_rules;type:json"`               // JSON数组：特殊规则

	// 显示配置
	IconURL      *string `json:"icon_url" gorm:"column:icon_url;type:varchar(500)"`              // 图标URL
	Color        *string `json:"color" gorm:"column:color;type:varchar(20)"`                     // 显示颜色
	DisplayOrder *int    `json:"display_order" gorm:"column:display_order;type:int;default:999"` // 显示顺序

	// 联系信息
	EmergencyContact *string `json:"emergency_contact" gorm:"column:emergency_contact;type:varchar(255)"` // 紧急联系方式
	SupportContact   *string `json:"support_contact" gorm:"column:support_contact;type:varchar(255)"`     // 客服联系方式
	LocalPartner     *string `json:"local_partner" gorm:"column:local_partner;type:varchar(255)"`         // 本地合作伙伴

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 附加元数据
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 标签
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`       // 备注

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (ServiceArea) TableName() string {
	return "t_service_areas"
}

// 区域类型常量
const (
	AreaTypeCity     = "city"
	AreaTypeSuburb   = "suburb"
	AreaTypeDistrict = "district"
	AreaTypeZone     = "zone"
	AreaTypeAirport  = "airport"
	AreaTypeSpecial  = "special"
)

// 状态常量
const (
	ServiceAreaStatusActive      = "active"
	ServiceAreaStatusInactive    = "inactive"
	ServiceAreaStatusMaintenance = "maintenance"
	ServiceAreaStatusPlanned     = "planned"
)

// 服务级别常量
const (
	ServiceLevelBasic    = "basic"
	ServiceLevelStandard = "standard"
	ServiceLevelPremium  = "premium"
	ServiceLevelFull     = "full"
)

// 定价层级常量
const (
	PricingTierEconomy  = "economy"
	PricingTierStandard = "standard"
	PricingTierPremium  = "premium"
)

// 创建新的服务区域对象
func NewServiceAreaV2() *ServiceArea {
	return &ServiceArea{
		ServiceAreaID: utils.GenerateServiceAreaID(),
		Salt:          utils.GenerateSalt(),
		ServiceAreaValues: &ServiceAreaValues{
			AreaType:            utils.StringPtr(AreaTypeCity),
			TimeZone:            utils.StringPtr("UTC"),
			Level:               utils.IntPtr(1),
			Priority:            utils.IntPtr(100),
			Status:              utils.StringPtr(ServiceAreaStatusActive),
			IsActive:            utils.BoolPtr(true),
			ServiceLevel:        utils.StringPtr(ServiceLevelStandard),
			Is24Hours:           utils.BoolPtr(false),
			PricingTier:         utils.StringPtr(PricingTierStandard),
			SurgeEnabled:        utils.BoolPtr(true),
			BaseFareRate:        utils.Float64Ptr(1.00),
			SurgeThreshold:      utils.Float64Ptr(2.0),
			DemandLevel:         utils.StringPtr(protocol.LevelMedium),
			SupplyLevel:         utils.StringPtr(protocol.LevelMedium),
			CurrentSurge:        utils.Float64Ptr(1.0),
			TotalRides:          utils.IntPtr(0),
			TotalDrivers:        utils.IntPtr(0),
			ActiveDrivers:       utils.IntPtr(0),
			TotalUsers:          utils.IntPtr(0),
			AverageWaitTime:     utils.IntPtr(0),
			AverageRating:       utils.Float64Ptr(0.00),
			TrafficLevel:        utils.StringPtr(protocol.LevelMedium),
			ParkingAvailability: utils.StringPtr(protocol.LevelMedium),
			PublicTransport:     utils.BoolPtr(false),
			AirportNearby:       utils.BoolPtr(false),
			IsAirportZone:       utils.BoolPtr(false),
			IsBusinessDistrict:  utils.BoolPtr(false),
			IsResidential:       utils.BoolPtr(false),
			IsTouristArea:       utils.BoolPtr(false),
			IsHighCrimeArea:     utils.BoolPtr(false),
			WeatherSensitive:    utils.BoolPtr(false),
			DisplayOrder:        utils.IntPtr(999),
		},
	}
}

// SetValues 更新ServiceAreaV2Values中的非nil值
func (s *ServiceAreaValues) SetValues(values *ServiceAreaValues) {
	if values == nil {
		return
	}

	if values.AreaName != nil {
		s.AreaName = values.AreaName
	}
	if values.DisplayName != nil {
		s.DisplayName = values.DisplayName
	}
	if values.Description != nil {
		s.Description = values.Description
	}
	if values.AreaType != nil {
		s.AreaType = values.AreaType
	}
	if values.CenterLat != nil {
		s.CenterLat = values.CenterLat
	}
	if values.CenterLng != nil {
		s.CenterLng = values.CenterLng
	}
	if values.Country != nil {
		s.Country = values.Country
	}
	if values.City != nil {
		s.City = values.City
	}
	if values.Status != nil {
		s.Status = values.Status
	}
	if values.Notes != nil {
		s.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		s.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (s *ServiceAreaValues) GetAreaName() string {
	if s.AreaName == nil {
		return ""
	}
	return *s.AreaName
}

func (s *ServiceAreaValues) GetDisplayName() string {
	if s.DisplayName == nil {
		return s.GetAreaName()
	}
	return *s.DisplayName
}

func (s *ServiceAreaValues) GetAreaType() string {
	if s.AreaType == nil {
		return AreaTypeCity
	}
	return *s.AreaType
}

func (s *ServiceAreaValues) GetCenterLat() float64 {
	if s.CenterLat == nil {
		return 0.0
	}
	return *s.CenterLat
}

func (s *ServiceAreaValues) GetCenterLng() float64 {
	if s.CenterLng == nil {
		return 0.0
	}
	return *s.CenterLng
}

func (s *ServiceAreaValues) GetRadius() float64 {
	if s.Radius == nil {
		return 5.0 // 默认5公里
	}
	return *s.Radius
}

func (s *ServiceAreaValues) GetStatus() string {
	if s.Status == nil {
		return ServiceAreaStatusActive
	}
	return *s.Status
}

func (s *ServiceAreaValues) GetIsActive() bool {
	if s.IsActive == nil {
		return true
	}
	return *s.IsActive
}

func (s *ServiceAreaValues) GetServiceLevel() string {
	if s.ServiceLevel == nil {
		return ServiceLevelStandard
	}
	return *s.ServiceLevel
}

func (s *ServiceAreaValues) GetPricingTier() string {
	if s.PricingTier == nil {
		return PricingTierStandard
	}
	return *s.PricingTier
}

func (s *ServiceAreaValues) GetCurrentSurge() float64 {
	if s.CurrentSurge == nil {
		return 1.0
	}
	return *s.CurrentSurge
}

func (s *ServiceAreaValues) GetLevel() int {
	if s.Level == nil {
		return 1
	}
	return *s.Level
}

func (s *ServiceAreaValues) GetPriority() int {
	if s.Priority == nil {
		return 100
	}
	return *s.Priority
}

func (s *ServiceAreaValues) GetIs24Hours() bool {
	if s.Is24Hours == nil {
		return false
	}
	return *s.Is24Hours
}

func (s *ServiceAreaValues) GetSurgeEnabled() bool {
	if s.SurgeEnabled == nil {
		return true
	}
	return *s.SurgeEnabled
}

// Setter 方法
func (s *ServiceAreaValues) SetAreaName(name string) *ServiceAreaValues {
	s.AreaName = &name
	return s
}

func (s *ServiceAreaValues) SetDisplayName(name string) *ServiceAreaValues {
	s.DisplayName = &name
	return s
}

func (s *ServiceAreaValues) SetDescription(desc string) *ServiceAreaValues {
	s.Description = &desc
	return s
}

func (s *ServiceAreaValues) SetAreaType(areaType string) *ServiceAreaValues {
	s.AreaType = &areaType
	return s
}

func (s *ServiceAreaValues) SetLocation(lat, lng, radius float64) *ServiceAreaValues {
	s.CenterLat = &lat
	s.CenterLng = &lng
	s.Radius = &radius
	return s
}

func (s *ServiceAreaValues) SetAdministrative(country, province, city, district string) *ServiceAreaValues {
	s.Country = &country
	s.Province = &province
	s.City = &city
	s.District = &district
	return s
}

func (s *ServiceAreaValues) SetStatus(status string) *ServiceAreaValues {
	s.Status = &status
	return s
}

func (s *ServiceAreaValues) SetActive(active bool) *ServiceAreaValues {
	s.IsActive = &active
	return s
}

func (s *ServiceAreaValues) SetServiceLevel(level string) *ServiceAreaValues {
	s.ServiceLevel = &level
	return s
}

func (s *ServiceAreaValues) SetPricingTier(tier string) *ServiceAreaValues {
	s.PricingTier = &tier
	return s
}

func (s *ServiceAreaValues) SetHierarchy(parentID string, level int) *ServiceAreaValues {
	s.ParentAreaID = &parentID
	s.Level = &level
	return s
}

func (s *ServiceAreaValues) SetSurge(currentSurge float64) *ServiceAreaValues {
	s.CurrentSurge = &currentSurge
	return s
}

func (s *ServiceAreaValues) Set24Hours(is24Hours bool) *ServiceAreaValues {
	s.Is24Hours = &is24Hours
	return s
}

func (s *ServiceAreaValues) SetSpecialZones(isAirport, isBusiness, isResidential, isTourist bool) *ServiceAreaValues {
	s.IsAirportZone = &isAirport
	s.IsBusinessDistrict = &isBusiness
	s.IsResidential = &isResidential
	s.IsTouristArea = &isTourist
	return s
}

// 业务方法
func (s *ServiceArea) IsActive() bool {
	return s.GetStatus() == ServiceAreaStatusActive && s.GetIsActive()
}

func (s *ServiceArea) IsInactive() bool {
	return s.GetStatus() == ServiceAreaStatusInactive || !s.GetIsActive()
}

func (s *ServiceArea) IsUnderMaintenance() bool {
	return s.GetStatus() == ServiceAreaStatusMaintenance
}

func (s *ServiceArea) IsPlanned() bool {
	return s.GetStatus() == ServiceAreaStatusPlanned
}

func (s *ServiceArea) IsAirportZone() bool {
	if s.ServiceAreaValues.IsAirportZone == nil {
		return false
	}
	return *s.ServiceAreaValues.IsAirportZone
}

func (s *ServiceArea) IsBusinessDistrict() bool {
	if s.ServiceAreaValues.IsBusinessDistrict == nil {
		return false
	}
	return *s.ServiceAreaValues.IsBusinessDistrict
}

func (s *ServiceArea) IsTouristArea() bool {
	if s.ServiceAreaValues.IsTouristArea == nil {
		return false
	}
	return *s.ServiceAreaValues.IsTouristArea
}

func (s *ServiceArea) IsHighCrimeArea() bool {
	if s.ServiceAreaValues.IsHighCrimeArea == nil {
		return false
	}
	return *s.ServiceAreaValues.IsHighCrimeArea
}

func (s *ServiceArea) HasSurgeActive() bool {
	return s.GetSurgeEnabled() && s.GetCurrentSurge() > 1.0
}

func (s *ServiceArea) IsHighDemandArea() bool {
	if s.DemandLevel == nil {
		return false
	}
	return *s.DemandLevel == protocol.LevelHigh
}

func (s *ServiceArea) IsLowSupplyArea() bool {
	if s.SupplyLevel == nil {
		return false
	}
	return *s.SupplyLevel == protocol.LevelLow
}

func (s *ServiceArea) ShouldApplySurge() bool {
	return s.GetSurgeEnabled() && s.IsHighDemandArea() && s.IsLowSupplyArea()
}

func (s *ServiceArea) CanAcceptRides() bool {
	return s.IsActive() && !s.IsUnderMaintenance()
}

// 地理位置相关方法
func (s *ServiceAreaValues) IsLocationWithinRadius(lat, lng float64) bool {
	if s.CenterLat == nil || s.CenterLng == nil {
		return false
	}

	distance := s.CalculateDistance(lat, lng)
	return distance <= s.GetRadius()
}

func (s *ServiceAreaValues) CalculateDistance(lat, lng float64) float64 {
	centerLat := s.GetCenterLat()
	centerLng := s.GetCenterLng()

	return haversineDistance(centerLat, centerLng, lat, lng)
}

// Haversine公式计算两点间距离
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // 地球半径(公里)

	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *ServiceAreaValues) IsLocationInPolygon(lat, lng float64) bool {
	if s.Polygon == nil {
		return s.IsLocationWithinRadius(lat, lng)
	}

	var polygon [][]float64
	if err := utils.FromJSON(*s.Polygon, &polygon); err != nil {
		return s.IsLocationWithinRadius(lat, lng)
	}

	return pointInPolygon(lat, lng, polygon)
}

// 射线法判断点是否在多边形内
func pointInPolygon(lat, lng float64, polygon [][]float64) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		if len(polygon[i]) < 2 || len(polygon[j]) < 2 {
			j = i
			continue
		}

		yi := polygon[i][0]
		xi := polygon[i][1]
		yj := polygon[j][0]
		xj := polygon[j][1]

		if ((yi > lat) != (yj > lat)) && (lng < (xj-xi)*(lat-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// 支持的服务管理
func (s *ServiceAreaValues) SetSupportedVehicleTypes(vehicleTypes []string) error {
	typesJSON, err := utils.ToJSON(vehicleTypes)
	if err != nil {
		return fmt.Errorf("failed to marshal vehicle types: %v", err)
	}

	s.SupportedVehicleTypes = &typesJSON
	return nil
}

func (s *ServiceAreaValues) GetSupportedVehicleTypes() []string {
	if s.SupportedVehicleTypes == nil {
		return []string{}
	}

	var types []string
	if err := utils.FromJSON(*s.SupportedVehicleTypes, &types); err != nil {
		return []string{}
	}

	return types
}

func (s *ServiceAreaValues) SupportsVehicleType(vehicleType string) bool {
	types := s.GetSupportedVehicleTypes()
	if len(types) == 0 {
		return true // 如果没有限制，支持所有车型
	}

	for _, t := range types {
		if t == vehicleType {
			return true
		}
	}

	return false
}

func (s *ServiceAreaValues) SetSupportedServices(services []string) error {
	servicesJSON, err := utils.ToJSON(services)
	if err != nil {
		return fmt.Errorf("failed to marshal services: %v", err)
	}

	s.SupportedServices = &servicesJSON
	return nil
}

func (s *ServiceAreaValues) GetSupportedServices() []string {
	if s.SupportedServices == nil {
		return []string{}
	}

	var services []string
	if err := utils.FromJSON(*s.SupportedServices, &services); err != nil {
		return []string{}
	}

	return services
}

// 运营时间管理
func (s *ServiceAreaValues) SetOperatingHours(hours map[string]interface{}) error {
	hoursJSON, err := utils.ToJSON(hours)
	if err != nil {
		return fmt.Errorf("failed to marshal operating hours: %v", err)
	}

	s.OperatingHours = &hoursJSON
	return nil
}

func (s *ServiceAreaValues) IsOperatingAtTime(dayOfWeek int, hour, minute int) bool {
	if s.GetIs24Hours() {
		return true
	}

	if s.OperatingHours == nil {
		return true // 如果没有设置，默认全天运营
	}

	var hours map[string]interface{}
	if err := utils.FromJSON(*s.OperatingHours, &hours); err != nil {
		return true
	}

	dayKey := fmt.Sprintf("day_%d", dayOfWeek)
	daySchedule, exists := hours[dayKey]
	if !exists {
		return false
	}

	schedule, ok := daySchedule.(map[string]interface{})
	if !ok {
		return false
	}

	startHour := int(schedule["start_hour"].(float64))
	startMinute := int(schedule["start_minute"].(float64))
	endHour := int(schedule["end_hour"].(float64))
	endMinute := int(schedule["end_minute"].(float64))

	currentTime := hour*60 + minute
	startTime := startHour*60 + startMinute
	endTime := endHour*60 + endMinute

	return currentTime >= startTime && currentTime <= endTime
}

// 统计更新方法
func (s *ServiceAreaValues) UpdateRideStats(totalRides int) *ServiceAreaValues {
	s.TotalRides = &totalRides
	return s
}

func (s *ServiceAreaValues) UpdateDriverStats(totalDrivers, activeDrivers int) *ServiceAreaValues {
	s.TotalDrivers = &totalDrivers
	s.ActiveDrivers = &activeDrivers
	return s
}

func (s *ServiceAreaValues) UpdateUserStats(totalUsers int) *ServiceAreaValues {
	s.TotalUsers = &totalUsers
	return s
}

func (s *ServiceAreaValues) UpdateWaitTime(waitTimeSeconds int) *ServiceAreaValues {
	s.AverageWaitTime = &waitTimeSeconds
	return s
}

func (s *ServiceAreaValues) UpdateRating(newRating float64) *ServiceAreaValues {
	s.AverageRating = &newRating
	return s
}

func (s *ServiceAreaValues) UpdateDemandSupply(demandLevel, supplyLevel string) *ServiceAreaValues {
	s.DemandLevel = &demandLevel
	s.SupplyLevel = &supplyLevel
	return s
}

// 标签管理
func (s *ServiceAreaValues) AddTag(tag string) *ServiceAreaValues {
	var tags []string
	if s.Tags != nil && *s.Tags != "" {
		tags = strings.Split(*s.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return s
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	s.Tags = &tagsStr
	return s
}

func (s *ServiceAreaValues) HasTag(tag string) bool {
	if s.Tags == nil || *s.Tags == "" {
		return false
	}

	tags := strings.Split(*s.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 便捷创建方法
func NewCityServiceArea(name, country, city string, lat, lng, radius float64) *ServiceArea {
	area := NewServiceAreaV2()
	area.SetAreaName(name).
		SetDisplayName(name).
		SetAreaType(AreaTypeCity).
		SetLocation(lat, lng, radius).
		SetAdministrative(country, "", city, "").
		SetServiceLevel(ServiceLevelStandard).
		SetPricingTier(PricingTierStandard)

	// 设置支持的车型
	vehicleTypes := []string{"economy", "comfort", "premium"}
	area.SetSupportedVehicleTypes(vehicleTypes)

	// 设置支持的服务
	services := []string{"ride_hailing", "delivery", "intercity"}
	area.SetSupportedServices(services)

	return area
}

func NewAirportServiceArea(airportName, city string, lat, lng float64) *ServiceArea {
	area := NewServiceAreaV2()
	area.SetAreaName(airportName).
		SetDisplayName(airportName+" Airport").
		SetAreaType(AreaTypeAirport).
		SetLocation(lat, lng, 2.0). // 机场区域通常较小
		SetAdministrative("Rwanda", "", city, "").
		SetServiceLevel(ServiceLevelPremium).
		SetPricingTier(PricingTierPremium).
		SetSpecialZones(true, false, false, false)

	// 设置机场附加费
	area.AirportSurcharge = utils.Float64Ptr(500.0)
	area.AirportNearby = utils.BoolPtr(true)

	// 设置24小时服务
	area.Set24Hours(true)

	// 添加标签
	area.AddTag("airport")
	area.AddTag("premium")
	area.AddTag("24hours")

	return area
}

func NewBusinessDistrictArea(name, city string, lat, lng, radius float64) *ServiceArea {
	area := NewServiceAreaV2()
	area.SetAreaName(name).
		SetDisplayName(name+" Business District").
		SetAreaType(AreaTypeDistrict).
		SetLocation(lat, lng, radius).
		SetAdministrative("Rwanda", "", city, "").
		SetServiceLevel(ServiceLevelStandard).
		SetPricingTier(PricingTierStandard).
		SetSpecialZones(false, true, false, false)

	// 商务区通常交通拥堵，停车困难
	area.TrafficLevel = utils.StringPtr(protocol.LevelHigh)
	area.ParkingAvailability = utils.StringPtr("limited")
	area.BusinessActivity = utils.StringPtr(protocol.LevelHigh)

	// 添加标签
	area.AddTag("business")
	area.AddTag("high-demand")

	return area
}
