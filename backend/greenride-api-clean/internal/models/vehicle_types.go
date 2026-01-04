package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"
)

// VehicleType 车辆类型表 - 不同车型的配置、价格和特性
type VehicleType struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	VehicleTypeID string `json:"vehicle_type_id" gorm:"column:vehicle_type_id;type:varchar(64);uniqueIndex"`
	Salt          string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*VehicleTypeValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type VehicleTypeValues struct {
	// 基本信息
	TypeName    *string `json:"type_name" gorm:"column:type_name;type:varchar(100);index"` // 车型名称
	DisplayName *string `json:"display_name" gorm:"column:display_name;type:varchar(100)"` // 显示名称
	Description *string `json:"description" gorm:"column:description;type:text"`           // 车型描述
	Category    *string `json:"category" gorm:"column:category;type:varchar(50);index"`    // economy, comfort, premium, luxury, suv, van
	Level       *string `json:"level" gorm:"column:level;type:varchar(50);index"`          // economy, comfort, premium, luxury

	// 车辆规格
	Capacity      *int     `json:"capacity" gorm:"column:capacity;type:int"`                        // 载客数量
	LuggageSpace  *int     `json:"luggage_space" gorm:"column:luggage_space;type:int"`              // 行李箱数量
	MinEngineSize *float64 `json:"min_engine_size" gorm:"column:min_engine_size;type:decimal(4,2)"` // 最小发动机排量
	MaxEngineSize *float64 `json:"max_engine_size" gorm:"column:max_engine_size;type:decimal(4,2)"` // 最大发动机排量
	FuelType      *string  `json:"fuel_type" gorm:"column:fuel_type;type:varchar(30)"`              // petrol, diesel, electric, hybrid
	Transmission  *string  `json:"transmission" gorm:"column:transmission;type:varchar(30)"`        // manual, automatic

	// 价格配置 - 已移除

	// 附加费用 - 已移除

	// 涌潮定价 - 已移除

	// 服务特性
	Features       *string `json:"features" gorm:"column:features;type:json"`               // JSON数组：特性列表
	Amenities      *string `json:"amenities" gorm:"column:amenities;type:json"`             // JSON数组：设施列表
	SafetyFeatures *string `json:"safety_features" gorm:"column:safety_features;type:json"` // JSON数组：安全特性

	// 车型要求 - 已移除

	// 司机要求 - 已移除

	// 可用性设置
	Status         *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"` // active, inactive, deprecated
	IsActive       *bool   `json:"is_active" gorm:"column:is_active;default:true"`                      // 是否激活
	ServiceAreas   *string `json:"service_areas" gorm:"column:service_areas;type:json"`                 // JSON数组：服务区域
	OperatingHours *string `json:"operating_hours" gorm:"column:operating_hours;type:json"`             // JSON对象：运营时间

	// 调度优先级
	Priority        *int     `json:"priority" gorm:"column:priority;type:int;default:100"`                          // 调度优先级
	PopularityScore *float64 `json:"popularity_score" gorm:"column:popularity_score;type:decimal(5,2);default:0.0"` // 受欢迎程度

	// 显示设置
	IconURL      *string `json:"icon_url" gorm:"column:icon_url;type:varchar(500)"`              // 图标URL
	ImageURL     *string `json:"image_url" gorm:"column:image_url;type:varchar(500)"`            // 车型图片URL
	Color        *string `json:"color" gorm:"column:color;type:varchar(20)"`                     // 主题颜色
	DisplayOrder *int    `json:"display_order" gorm:"column:display_order;type:int;default:999"` // 显示顺序

	// 统计数据 - 已移除

	// 市场数据 - 已移除

	// 环保信息 - 已移除

	// 维护要求 - 已移除

	// 佣金设置 - 已移除

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 附加元数据
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 标签
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`       // 备注

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (VehicleType) TableName() string {
	return "t_vehicle_types"
}

// 以下常量已移至 protocol/const.vehicle.go
// VehicleCategoryEconomy, VehicleCategoryComfort, VehicleCategoryPremium, VehicleCategoryLuxury, VehicleCategorySUV, VehicleCategoryVan
// VehicleTypeFuelPetrol, VehicleTypeFuelDiesel, VehicleTypeFuelElectric, VehicleTypeFuelHybrid
// TransmissionManual, TransmissionAutomatic
// StatusActive, StatusInactive, StatusDeprecated
// LevelLow, LevelMedium, LevelHigh

// 创建新的车辆类型对象
func NewVehicleType() *VehicleType {
	return &VehicleType{
		VehicleTypeID: utils.GenerateVehicleTypeID(),
		Salt:          utils.GenerateSalt(),
		VehicleTypeValues: &VehicleTypeValues{
			Category:        utils.StringPtr(protocol.VehicleCategorySedan),
			Level:           utils.StringPtr(protocol.VehicleLevelEconomy),
			FuelType:        utils.StringPtr(protocol.FuelTypeGasoline),
			Transmission:    utils.StringPtr(protocol.TransmissionManual),
			Status:          utils.StringPtr(protocol.StatusActive),
			IsActive:        utils.BoolPtr(true),
			Priority:        utils.IntPtr(100),
			PopularityScore: utils.Float64Ptr(0.0),
			DisplayOrder:    utils.IntPtr(999),
		},
	}
}

// SetValues 更新VehicleTypeV2Values中的非nil值
func (v *VehicleTypeValues) SetValues(values *VehicleTypeValues) {
	if values == nil {
		return
	}

	if values.TypeName != nil {
		v.TypeName = values.TypeName
	}
	if values.DisplayName != nil {
		v.DisplayName = values.DisplayName
	}
	if values.Description != nil {
		v.Description = values.Description
	}
	if values.Category != nil {
		v.Category = values.Category
	}
	if values.Level != nil {
		v.Level = values.Level
	}
	if values.Capacity != nil {
		v.Capacity = values.Capacity
	}

	if values.Status != nil {
		v.Status = values.Status
	}
	if values.Notes != nil {
		v.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		v.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (v *VehicleTypeValues) GetTypeName() string {
	if v.TypeName == nil {
		return ""
	}
	return *v.TypeName
}

func (v *VehicleTypeValues) GetDisplayName() string {
	if v.DisplayName == nil {
		return v.GetTypeName()
	}
	return *v.DisplayName
}

func (v *VehicleTypeValues) GetCategory() string {
	if v.Category == nil {
		return ""
	}
	return *v.Category
}

func (v *VehicleTypeValues) GetLevel() string {
	if v.Level == nil {
		return ""
	}
	return *v.Level
}

func (v *VehicleTypeValues) GetCapacity() int {
	if v.Capacity == nil {
		return 0
	}
	return *v.Capacity
}

func (v *VehicleTypeValues) GetStatus() string {
	if v.Status == nil {
		return protocol.StatusActive
	}
	return *v.Status
}

func (v *VehicleTypeValues) GetIsActive() bool {
	if v.IsActive == nil {
		return true
	}
	return *v.IsActive
}

func (v *VehicleTypeValues) GetPriority() int {
	if v.Priority == nil {
		return 0
	}
	return *v.Priority
}

// Setter 方法
func (v *VehicleTypeValues) SetTypeName(name string) *VehicleTypeValues {
	v.TypeName = &name
	return v
}

func (v *VehicleTypeValues) SetDisplayName(name string) *VehicleTypeValues {
	v.DisplayName = &name
	return v
}

func (v *VehicleTypeValues) SetDescription(desc string) *VehicleTypeValues {
	v.Description = &desc
	return v
}

func (v *VehicleTypeValues) SetCategory(category string) *VehicleTypeValues {
	v.Category = &category
	return v
}

func (v *VehicleTypeValues) SetLevel(level string) *VehicleTypeValues {
	v.Level = &level
	return v
}

func (v *VehicleTypeValues) SetCapacity(capacity int) *VehicleTypeValues {
	v.Capacity = &capacity
	return v
}

func (v *VehicleTypeValues) SetSpecs(capacity, luggage int, engineMin, engineMax float64) *VehicleTypeValues {
	v.Capacity = &capacity
	v.LuggageSpace = &luggage
	v.MinEngineSize = &engineMin
	v.MaxEngineSize = &engineMax
	return v
}

func (v *VehicleTypeValues) SetFuelAndTransmission(fuel, transmission string) *VehicleTypeValues {
	v.FuelType = &fuel
	v.Transmission = &transmission
	return v
}

func (v *VehicleTypeValues) SetStatus(status string) *VehicleTypeValues {
	v.Status = &status
	return v
}

func (v *VehicleTypeValues) SetActive(active bool) *VehicleTypeValues {
	v.IsActive = &active
	return v
}

func (v *VehicleTypeValues) SetPriority(priority int) *VehicleTypeValues {
	v.Priority = &priority
	return v
}

func (v *VehicleTypeValues) SetDisplayOrder(order int) *VehicleTypeValues {
	v.DisplayOrder = &order
	return v
}

func (v *VehicleTypeValues) SetImages(iconURL, imageURL string) *VehicleTypeValues {
	v.IconURL = &iconURL
	v.ImageURL = &imageURL
	return v
}

// 业务方法
func (v *VehicleType) IsActive() bool {
	return v.GetStatus() == protocol.StatusActive && v.GetIsActive()
}

func (v *VehicleType) IsInactive() bool {
	return v.GetStatus() == protocol.StatusInactive || !v.GetIsActive()
}

func (v *VehicleType) IsDeprecated() bool {
	return v.GetStatus() == protocol.StatusDeprecated
}

func (v *VehicleType) IsElectric() bool {
	if v.FuelType == nil {
		return false
	}
	return *v.FuelType == protocol.FuelTypeElectric
}

func (v *VehicleType) IsLuxury() bool {
	return v.GetCategory() == protocol.VehicleLevelLuxury || v.GetCategory() == protocol.VehicleLevelPremium
}

func (v *VehicleType) IsHighCapacity() bool {
	return v.GetCapacity() > 4
}

func (v *VehicleType) CanAcceptRides() bool {
	return v.IsActive()
}

func (v *VehicleType) ShouldApplySurge() bool {
	return false // 涌潮定价相关字段已移除
}

// 特性管理方法
func (v *VehicleTypeValues) AddFeature(feature string) error {
	var features []string
	if v.Features != nil {
		if err := utils.FromJSON(*v.Features, &features); err != nil {
			return fmt.Errorf("failed to parse existing features: %v", err)
		}
	}

	// 避免重复添加
	for _, f := range features {
		if f == feature {
			return nil
		}
	}

	features = append(features, feature)
	featuresJSON, err := utils.ToJSON(features)
	if err != nil {
		return fmt.Errorf("failed to marshal features: %v", err)
	}

	v.Features = &featuresJSON
	return nil
}

func (v *VehicleTypeValues) SetFeatures(features []string) error {
	featuresJSON, err := utils.ToJSON(features)
	if err != nil {
		return fmt.Errorf("failed to marshal features: %v", err)
	}

	v.Features = &featuresJSON
	return nil
}

func (v *VehicleTypeValues) GetFeatures() []string {
	if v.Features == nil {
		return []string{}
	}

	var features []string
	if err := utils.FromJSON(*v.Features, &features); err != nil {
		return []string{}
	}

	return features
}

func (v *VehicleTypeValues) SetAmenities(amenities []string) error {
	amenitiesJSON, err := utils.ToJSON(amenities)
	if err != nil {
		return fmt.Errorf("failed to marshal amenities: %v", err)
	}

	v.Amenities = &amenitiesJSON
	return nil
}

func (v *VehicleTypeValues) GetAmenities() []string {
	if v.Amenities == nil {
		return []string{}
	}

	var amenities []string
	if err := utils.FromJSON(*v.Amenities, &amenities); err != nil {
		return []string{}
	}

	return amenities
}

func (v *VehicleTypeValues) SetServiceAreas(areas []string) error {
	areasJSON, err := utils.ToJSON(areas)
	if err != nil {
		return fmt.Errorf("failed to marshal service areas: %v", err)
	}

	v.ServiceAreas = &areasJSON
	return nil
}

func (v *VehicleTypeValues) IsAvailableInArea(areaID string) bool {
	if v.ServiceAreas == nil {
		return true // 如果没有限制，则在所有区域可用
	}

	var areas []string
	if err := utils.FromJSON(*v.ServiceAreas, &areas); err != nil {
		return false
	}

	for _, area := range areas {
		if area == areaID {
			return true
		}
	}

	return false
}

// 价格计算方法
// 统计更新方法
// 标签管理
func (v *VehicleTypeValues) AddTag(tag string) *VehicleTypeValues {
	var tags []string
	if v.Tags != nil && *v.Tags != "" {
		tags = strings.Split(*v.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return v
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	v.Tags = &tagsStr
	return v
}

func (v *VehicleTypeValues) HasTag(tag string) bool {
	if v.Tags == nil || *v.Tags == "" {
		return false
	}

	tags := strings.Split(*v.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 便捷创建方法
func NewEconomyVehicleType(name string, capacity int, baseRate, perKmRate float64) *VehicleType {
	vehicleType := NewVehicleType()
	vehicleType.SetTypeName(name).
		SetDisplayName(name).
		SetCategory("economy").
		SetLevel("economy").
		SetCapacity(capacity).
		SetPriority(100)

	// 设置基本特性
	features := []string{"Air Conditioning", "Radio", "Safe Driving"}
	vehicleType.SetFeatures(features)

	return vehicleType
}

func NewLuxuryVehicleType(name string, baseRate, perKmRate float64) *VehicleType {
	vehicleType := NewVehicleType()
	vehicleType.SetTypeName(name).
		SetDisplayName(name).
		SetCategory("luxury").
		SetLevel("luxury").
		SetCapacity(4).
		SetPriority(10)

	// 设置豪华特性
	features := []string{"Leather Seats", "Premium Audio", "Climate Control", "WiFi", "Water Bottles"}
	amenities := []string{"Phone Charger", "Tissues", "Umbrella"}
	vehicleType.SetFeatures(features)
	vehicleType.SetAmenities(amenities)

	// Driver requirements methods have been removed

	return vehicleType
}

func NewSUVVehicleType(name string, capacity int, baseRate, perKmRate float64) *VehicleType {
	vehicleType := NewVehicleType()
	vehicleType.SetTypeName(name).
		SetDisplayName(name).
		SetCategory(protocol.VehicleCategorySUV).
		SetLevel("comfort").
		SetCapacity(capacity).
		SetSpecs(capacity, 3, 2.0, 4.0).
		SetPriority(50)

	// 设置SUV特性
	features := []string{"High Clearance", "4WD", "Large Luggage Space", "Air Conditioning"}
	vehicleType.SetFeatures(features)

	return vehicleType
}

func NewElectricVehicleType(name string, capacity int, baseRate, perKmRate float64) *VehicleType {
	vehicleType := NewVehicleType()
	vehicleType.SetTypeName(name).
		SetDisplayName(name).
		SetCategory("comfort").
		SetLevel("comfort").
		SetCapacity(capacity).
		SetFuelAndTransmission(protocol.FuelTypeElectric, protocol.TransmissionAutomatic).
		SetPriority(20)

	// Environmental fields have been removed from the model
	// 设置电动车特性
	features := []string{"Zero Emissions", "Silent Operation", "Fast Charging", "Digital Dashboard"}
	vehicleType.SetFeatures(features)
	vehicleType.AddTag("eco-friendly")
	vehicleType.AddTag("electric")

	return vehicleType
}
