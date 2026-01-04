package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"
)

// UserAddress 用户地址表 - 用户常用地址和收货地址管理
type UserAddress struct {
	ID            int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserAddressID string `json:"user_address_id" gorm:"column:user_address_id;type:varchar(64);uniqueIndex"`
	Salt          string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*UserAddressValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserAddressValues struct {
	// 关联信息
	UserID *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"` // 用户ID

	// 地址类型和标识
	AddressType *string `json:"address_type" gorm:"column:address_type;type:varchar(30);index"` // home, work, other, pickup, dropoff, delivery
	Label       *string `json:"label" gorm:"column:label;type:varchar(100)"`                    // 地址标签名称，如"家"、"公司"等
	IsDefault   *bool   `json:"is_default" gorm:"column:is_default;default:false"`              // 是否为默认地址
	IsFavorite  *bool   `json:"is_favorite" gorm:"column:is_favorite;default:false"`            // 是否为收藏地址

	// 地理位置信息
	Latitude  *float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8);index"`   // 纬度
	Longitude *float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8);index"` // 经度
	Accuracy  *float64 `json:"accuracy" gorm:"column:accuracy;type:decimal(8,2)"`          // 定位精度（米）
	Altitude  *float64 `json:"altitude" gorm:"column:altitude;type:decimal(8,2)"`          // 海拔高度（米）

	// 详细地址信息
	FormattedAddress *string `json:"formatted_address" gorm:"column:formatted_address;type:text"`  // 格式化完整地址
	StreetNumber     *string `json:"street_number" gorm:"column:street_number;type:varchar(50)"`   // 门牌号
	StreetName       *string `json:"street_name" gorm:"column:street_name;type:varchar(200)"`      // 街道名称
	Neighborhood     *string `json:"neighborhood" gorm:"column:neighborhood;type:varchar(100)"`    // 社区/小区
	District         *string `json:"district" gorm:"column:district;type:varchar(100);index"`      // 区/县
	City             *string `json:"city" gorm:"column:city;type:varchar(100);index"`              // 城市
	Province         *string `json:"province" gorm:"column:province;type:varchar(100);index"`      // 省/州
	Country          *string `json:"country" gorm:"column:country;type:varchar(100);index"`        // 国家
	PostalCode       *string `json:"postal_code" gorm:"column:postal_code;type:varchar(20);index"` // 邮政编码

	// 建筑物信息
	BuildingName   *string `json:"building_name" gorm:"column:building_name;type:varchar(200)"`    // 建筑物名称
	BuildingNumber *string `json:"building_number" gorm:"column:building_number;type:varchar(50)"` // 建筑物编号
	Floor          *string `json:"floor" gorm:"column:floor;type:varchar(20)"`                     // 楼层
	Unit           *string `json:"unit" gorm:"column:unit;type:varchar(50)"`                       // 单元/房间号
	Entrance       *string `json:"entrance" gorm:"column:entrance;type:varchar(50)"`               // 入口

	// 联系信息
	ContactName  *string `json:"contact_name" gorm:"column:contact_name;type:varchar(100)"`   // 联系人姓名
	ContactPhone *string `json:"contact_phone" gorm:"column:contact_phone;type:varchar(30)"`  // 联系电话
	ContactEmail *string `json:"contact_email" gorm:"column:contact_email;type:varchar(255)"` // 联系邮箱

	// 地址状态
	Status             *string `json:"status" gorm:"column:status;type:varchar(20);index;default:'active'"`    // active, inactive, deleted
	IsVerified         *bool   `json:"is_verified" gorm:"column:is_verified;default:false"`                    // 是否验证
	VerifiedAt         *int64  `json:"verified_at" gorm:"column:verified_at"`                                  // 验证时间
	VerificationMethod *string `json:"verification_method" gorm:"column:verification_method;type:varchar(50)"` // manual, sms, email, gps

	// 访问频率统计
	UsageCount  *int   `json:"usage_count" gorm:"column:usage_count;type:int;default:0"` // 使用次数
	LastUsedAt  *int64 `json:"last_used_at" gorm:"column:last_used_at"`                  // 最后使用时间
	FirstUsedAt *int64 `json:"first_used_at" gorm:"column:first_used_at"`                // 首次使用时间

	// 地址质量评分
	QualityScore   *float64 `json:"quality_score" gorm:"column:quality_score;type:decimal(3,2)"`     // 地址质量评分(0-5.0)
	CompletionRate *float64 `json:"completion_rate" gorm:"column:completion_rate;type:decimal(5,2)"` // 地址完整度(0-100%)
	AccuracyLevel  *string  `json:"accuracy_level" gorm:"column:accuracy_level;type:varchar(20)"`    // high, medium, low

	// 服务区域信息
	ServiceAreaID   *string `json:"service_area_id" gorm:"column:service_area_id;type:varchar(64);index"` // 所属服务区域
	IsInServiceArea *bool   `json:"is_in_service_area" gorm:"column:is_in_service_area;default:true"`     // 是否在服务区域内
	ZoneType        *string `json:"zone_type" gorm:"column:zone_type;type:varchar(30)"`                   // residential, commercial, industrial, mixed

	// 交通信息
	NearbyLandmarks    *string `json:"nearby_landmarks" gorm:"column:nearby_landmarks;type:json"`       // JSON数组：附近地标
	PublicTransport    *string `json:"public_transport" gorm:"column:public_transport;type:json"`       // JSON数组：附近公共交通
	AccessInstructions *string `json:"access_instructions" gorm:"column:access_instructions;type:text"` // 到达指引
	ParkingInfo        *string `json:"parking_info" gorm:"column:parking_info;type:text"`               // 停车信息

	// 特殊标记
	IsAirport        *bool `json:"is_airport" gorm:"column:is_airport;default:false"`                 // 是否机场
	IsHospital       *bool `json:"is_hospital" gorm:"column:is_hospital;default:false"`               // 是否医院
	IsSchool         *bool `json:"is_school" gorm:"column:is_school;default:false"`                   // 是否学校
	IsShoppingMall   *bool `json:"is_shopping_mall" gorm:"column:is_shopping_mall;default:false"`     // 是否购物中心
	IsOfficeBuilding *bool `json:"is_office_building" gorm:"column:is_office_building;default:false"` // 是否写字楼
	IsResidential    *bool `json:"is_residential" gorm:"column:is_residential;default:true"`          // 是否住宅

	// 安全等级
	SafetyLevel   *string `json:"safety_level" gorm:"column:safety_level;type:varchar(20);default:'medium'"` // high, medium, low
	CrimeRate     *string `json:"crime_rate" gorm:"column:crime_rate;type:varchar(20);default:'medium'"`     // low, medium, high
	LightingLevel *string `json:"lighting_level" gorm:"column:lighting_level;type:varchar(20)"`              // good, average, poor

	// 访问限制
	AccessRestrictions  *string `json:"access_restrictions" gorm:"column:access_restrictions;type:json"`   // JSON数组：访问限制
	TimeRestrictions    *string `json:"time_restrictions" gorm:"column:time_restrictions;type:json"`       // JSON对象：时间限制
	VehicleRestrictions *string `json:"vehicle_restrictions" gorm:"column:vehicle_restrictions;type:json"` // JSON数组：车辆限制

	// 成本信息
	ParkingCost *float64 `json:"parking_cost" gorm:"column:parking_cost;type:decimal(8,2)"` // 停车费用
	AccessFee   *float64 `json:"access_fee" gorm:"column:access_fee;type:decimal(8,2)"`     // 通行费用
	TollInfo    *string  `json:"toll_info" gorm:"column:toll_info;type:text"`               // 过路费信息

	// 天气敏感性
	WeatherSensitive   *bool   `json:"weather_sensitive" gorm:"column:weather_sensitive;default:false"`    // 是否受天气影响
	FloodRisk          *string `json:"flood_risk" gorm:"column:flood_risk;type:varchar(20);default:'low'"` // low, medium, high
	AccessibilityNotes *string `json:"accessibility_notes" gorm:"column:accessibility_notes;type:text"`    // 可达性说明

	// 元数据
	Source           *string  `json:"source" gorm:"column:source;type:varchar(50)"`                       // manual, google_maps, baidu_maps, auto_detected
	SourceReference  *string  `json:"source_reference" gorm:"column:source_reference;type:varchar(255)"`  // 地址来源引用ID
	Confidence       *float64 `json:"confidence" gorm:"column:confidence;type:decimal(5,2)"`              // 置信度(0-100%)
	ValidationStatus *string  `json:"validation_status" gorm:"column:validation_status;type:varchar(30)"` // pending, validated, failed, manual

	// 个性化设置
	Nickname *string `json:"nickname" gorm:"column:nickname;type:varchar(100)"`    // 个性化昵称
	IconType *string `json:"icon_type" gorm:"column:icon_type;type:varchar(30)"`   // home, work, heart, star, custom
	Color    *string `json:"color" gorm:"column:color;type:varchar(20)"`           // 显示颜色
	Priority *int    `json:"priority" gorm:"column:priority;type:int;default:100"` // 显示优先级

	// 共享设置
	IsShared         *bool   `json:"is_shared" gorm:"column:is_shared;default:false"`             // 是否共享地址
	SharedWith       *string `json:"shared_with" gorm:"column:shared_with;type:json"`             // JSON数组：共享给谁
	SharePermissions *string `json:"share_permissions" gorm:"column:share_permissions;type:json"` // JSON对象：共享权限

	// 备注和标签
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`       // 备注
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 标签
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 附加元数据

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (UserAddress) TableName() string {
	return "t_user_addresses"
}

// 地址类型常量
const (
	AddressTypeHome     = "home"
	AddressTypeWork     = "work"
	AddressTypeOther    = "other"
	AddressTypePickup   = "pickup"
	AddressTypeDropoff  = "dropoff"
	AddressTypeDelivery = "delivery"
)

// 地址状态常量
const (
	AddressStatusActive   = "active"
	AddressStatusInactive = "inactive"
	AddressStatusDeleted  = "deleted"
)

// 验证方式常量
const (
	VerificationMethodManual = "manual"
	VerificationMethodGPS    = "gps"
)

// 区域类型常量
const (
	ZoneTypeResidential = "residential"
	ZoneTypeCommercial  = "commercial"
	ZoneTypeIndustrial  = "industrial"
	ZoneTypeMixed       = "mixed"
)

// 图标类型常量
const (
	IconTypeHome   = "home"
	IconTypeWork   = "work"
	IconTypeHeart  = "heart"
	IconTypeStar   = "star"
	IconTypeCustom = "custom"
	IconTypeOther  = "other"
)

// 创建新的用户地址对象
func NewUserAddressV2() *UserAddress {
	return &UserAddress{
		UserAddressID: utils.GenerateUserAddressID(),
		Salt:          utils.GenerateSalt(),
		UserAddressValues: &UserAddressValues{
			AddressType:      utils.StringPtr(AddressTypeOther),
			IsDefault:        utils.BoolPtr(false),
			IsFavorite:       utils.BoolPtr(false),
			Status:           utils.StringPtr(AddressStatusActive),
			IsVerified:       utils.BoolPtr(false),
			UsageCount:       utils.IntPtr(0),
			QualityScore:     utils.Float64Ptr(0.0),
			CompletionRate:   utils.Float64Ptr(0.0),
			AccuracyLevel:    utils.StringPtr(protocol.LevelMedium),
			IsInServiceArea:  utils.BoolPtr(true),
			ZoneType:         utils.StringPtr(ZoneTypeResidential),
			IsAirport:        utils.BoolPtr(false),
			IsHospital:       utils.BoolPtr(false),
			IsSchool:         utils.BoolPtr(false),
			IsShoppingMall:   utils.BoolPtr(false),
			IsOfficeBuilding: utils.BoolPtr(false),
			IsResidential:    utils.BoolPtr(true),
			SafetyLevel:      utils.StringPtr(protocol.LevelMedium),
			CrimeRate:        utils.StringPtr(protocol.LevelMedium),
			WeatherSensitive: utils.BoolPtr(false),
			FloodRisk:        utils.StringPtr(protocol.LevelLow),
			Source:           utils.StringPtr("manual"),
			Confidence:       utils.Float64Ptr(100.0),
			ValidationStatus: utils.StringPtr("pending"),
			IconType:         utils.StringPtr(IconTypeOther),
			Priority:         utils.IntPtr(100),
			IsShared:         utils.BoolPtr(false),
		},
	}
}

// SetValues 更新UserAddressV2Values中的非nil值
func (u *UserAddressValues) SetValues(values *UserAddressValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		u.UserID = values.UserID
	}
	if values.AddressType != nil {
		u.AddressType = values.AddressType
	}
	if values.Label != nil {
		u.Label = values.Label
	}
	if values.Latitude != nil {
		u.Latitude = values.Latitude
	}
	if values.Longitude != nil {
		u.Longitude = values.Longitude
	}
	if values.FormattedAddress != nil {
		u.FormattedAddress = values.FormattedAddress
	}
	if values.City != nil {
		u.City = values.City
	}
	if values.Country != nil {
		u.Country = values.Country
	}
	if values.Status != nil {
		u.Status = values.Status
	}
	if values.Notes != nil {
		u.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		u.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (u *UserAddressValues) GetUserID() string {
	if u.UserID == nil {
		return ""
	}
	return *u.UserID
}

func (u *UserAddressValues) GetAddressType() string {
	if u.AddressType == nil {
		return AddressTypeOther
	}
	return *u.AddressType
}

func (u *UserAddressValues) GetLabel() string {
	if u.Label == nil {
		return ""
	}
	return *u.Label
}

func (u *UserAddressValues) GetLatitude() float64 {
	if u.Latitude == nil {
		return 0.0
	}
	return *u.Latitude
}

func (u *UserAddressValues) GetLongitude() float64 {
	if u.Longitude == nil {
		return 0.0
	}
	return *u.Longitude
}

func (u *UserAddressValues) GetFormattedAddress() string {
	if u.FormattedAddress == nil {
		return ""
	}
	return *u.FormattedAddress
}

func (u *UserAddressValues) GetCity() string {
	if u.City == nil {
		return ""
	}
	return *u.City
}

func (u *UserAddressValues) GetCountry() string {
	if u.Country == nil {
		return ""
	}
	return *u.Country
}

func (u *UserAddressValues) GetStatus() string {
	if u.Status == nil {
		return AddressStatusActive
	}
	return *u.Status
}

func (u *UserAddressValues) GetIsDefault() bool {
	if u.IsDefault == nil {
		return false
	}
	return *u.IsDefault
}

func (u *UserAddressValues) GetIsFavorite() bool {
	if u.IsFavorite == nil {
		return false
	}
	return *u.IsFavorite
}

func (u *UserAddressValues) GetIsVerified() bool {
	if u.IsVerified == nil {
		return false
	}
	return *u.IsVerified
}

func (u *UserAddressValues) GetUsageCount() int {
	if u.UsageCount == nil {
		return 0
	}
	return *u.UsageCount
}

func (u *UserAddressValues) GetQualityScore() float64 {
	if u.QualityScore == nil {
		return 0.0
	}
	return *u.QualityScore
}

func (u *UserAddressValues) GetSafetyLevel() string {
	if u.SafetyLevel == nil {
		return protocol.LevelMedium
	}
	return *u.SafetyLevel
}

func (u *UserAddressValues) GetIconType() string {
	if u.IconType == nil {
		return IconTypeOther
	}
	return *u.IconType
}

// Setter 方法
func (u *UserAddressValues) SetUserID(userID string) *UserAddressValues {
	u.UserID = &userID
	return u
}

func (u *UserAddressValues) SetAddressType(addressType string) *UserAddressValues {
	u.AddressType = &addressType
	return u
}

func (u *UserAddressValues) SetLabel(label string) *UserAddressValues {
	u.Label = &label
	return u
}

func (u *UserAddressValues) SetLocation(lat, lng float64) *UserAddressValues {
	u.Latitude = &lat
	u.Longitude = &lng
	return u
}

func (u *UserAddressValues) SetFormattedAddress(address string) *UserAddressValues {
	u.FormattedAddress = &address
	return u
}

func (u *UserAddressValues) SetAddress(streetNumber, streetName, district, city, province, country string) *UserAddressValues {
	u.StreetNumber = &streetNumber
	u.StreetName = &streetName
	u.District = &district
	u.City = &city
	u.Province = &province
	u.Country = &country
	return u
}

func (u *UserAddressValues) SetBuildingInfo(buildingName, floor, unit string) *UserAddressValues {
	u.BuildingName = &buildingName
	u.Floor = &floor
	u.Unit = &unit
	return u
}

func (u *UserAddressValues) SetContact(name, phone, email string) *UserAddressValues {
	u.ContactName = &name
	u.ContactPhone = &phone
	u.ContactEmail = &email
	return u
}

func (u *UserAddressValues) SetStatus(status string) *UserAddressValues {
	u.Status = &status
	return u
}

func (u *UserAddressValues) SetDefault(isDefault bool) *UserAddressValues {
	u.IsDefault = &isDefault
	return u
}

func (u *UserAddressValues) SetFavorite(isFavorite bool) *UserAddressValues {
	u.IsFavorite = &isFavorite
	return u
}

func (u *UserAddressValues) SetVerified(isVerified bool, method string) *UserAddressValues {
	u.IsVerified = &isVerified
	u.VerificationMethod = &method
	if isVerified {
		now := utils.TimeNowMilli()
		u.VerifiedAt = &now
	}
	return u
}

func (u *UserAddressValues) SetQuality(score float64, completionRate float64, accuracyLevel string) *UserAddressValues {
	u.QualityScore = &score
	u.CompletionRate = &completionRate
	u.AccuracyLevel = &accuracyLevel
	return u
}

func (u *UserAddressValues) SetSafety(safetyLevel, crimeRate, lightingLevel string) *UserAddressValues {
	u.SafetyLevel = &safetyLevel
	u.CrimeRate = &crimeRate
	u.LightingLevel = &lightingLevel
	return u
}

func (u *UserAddressValues) SetSpecialTypes(isAirport, isHospital, isSchool, isShoppingMall, isOfficeBuilding bool) *UserAddressValues {
	u.IsAirport = &isAirport
	u.IsHospital = &isHospital
	u.IsSchool = &isSchool
	u.IsShoppingMall = &isShoppingMall
	u.IsOfficeBuilding = &isOfficeBuilding
	return u
}

func (u *UserAddressValues) SetPersonalization(nickname, iconType, color string, priority int) *UserAddressValues {
	u.Nickname = &nickname
	u.IconType = &iconType
	u.Color = &color
	u.Priority = &priority
	return u
}

// 业务方法
func (u *UserAddress) IsActive() bool {
	return u.GetStatus() == AddressStatusActive
}

func (u *UserAddress) IsDeleted() bool {
	return u.GetStatus() == AddressStatusDeleted
}

func (u *UserAddress) IsVerified() bool {
	return u.GetIsVerified()
}

func (u *UserAddress) IsDefault() bool {
	return u.GetIsDefault()
}

func (u *UserAddress) IsFavorite() bool {
	return u.GetIsFavorite()
}

func (u *UserAddress) IsHomeAddress() bool {
	return u.GetAddressType() == AddressTypeHome
}

func (u *UserAddress) IsWorkAddress() bool {
	return u.GetAddressType() == AddressTypeWork
}

func (u *UserAddress) IsSpecialLocation() bool {
	return u.IsAirport() || u.IsHospital() || u.IsSchool() || u.IsShoppingMall()
}

func (u *UserAddress) IsAirport() bool {
	if u.UserAddressValues.IsAirport == nil {
		return false
	}
	return *u.UserAddressValues.IsAirport
}

func (u *UserAddress) IsHospital() bool {
	if u.UserAddressValues.IsHospital == nil {
		return false
	}
	return *u.UserAddressValues.IsHospital
}

func (u *UserAddress) IsSchool() bool {
	if u.UserAddressValues.IsSchool == nil {
		return false
	}
	return *u.UserAddressValues.IsSchool
}

func (u *UserAddress) IsShoppingMall() bool {
	if u.UserAddressValues.IsShoppingMall == nil {
		return false
	}
	return *u.UserAddressValues.IsShoppingMall
}

func (u *UserAddress) IsHighQuality() bool {
	return u.GetQualityScore() >= 4.0
}

func (u *UserAddress) IsLowQuality() bool {
	return u.GetQualityScore() < 2.0
}

func (u *UserAddress) IsSafeLocation() bool {
	return u.GetSafetyLevel() == protocol.LevelHigh
}

func (u *UserAddress) IsUnsafeLocation() bool {
	return u.GetSafetyLevel() == protocol.LevelLow
}

func (u *UserAddress) IsFrequentlyUsed() bool {
	return u.GetUsageCount() >= 10
}

func (u *UserAddress) HasValidLocation() bool {
	lat := u.GetLatitude()
	lng := u.GetLongitude()
	return lat != 0.0 && lng != 0.0 && lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

func (u *UserAddress) IsInServiceArea() bool {
	if u.UserAddressValues.IsInServiceArea == nil {
		return true
	}
	return *u.UserAddressValues.IsInServiceArea
}

// 地址质量评估
func (u *UserAddressValues) CalculateCompletionRate() float64 {
	totalFields := 15
	completedFields := 0

	if u.FormattedAddress != nil && *u.FormattedAddress != "" {
		completedFields++
	}
	if u.StreetNumber != nil && *u.StreetNumber != "" {
		completedFields++
	}
	if u.StreetName != nil && *u.StreetName != "" {
		completedFields++
	}
	if u.District != nil && *u.District != "" {
		completedFields++
	}
	if u.City != nil && *u.City != "" {
		completedFields++
	}
	if u.Province != nil && *u.Province != "" {
		completedFields++
	}
	if u.Country != nil && *u.Country != "" {
		completedFields++
	}
	if u.PostalCode != nil && *u.PostalCode != "" {
		completedFields++
	}
	if u.Latitude != nil && *u.Latitude != 0 {
		completedFields++
	}
	if u.Longitude != nil && *u.Longitude != 0 {
		completedFields++
	}
	if u.BuildingName != nil && *u.BuildingName != "" {
		completedFields++
	}
	if u.Floor != nil && *u.Floor != "" {
		completedFields++
	}
	if u.Unit != nil && *u.Unit != "" {
		completedFields++
	}
	if u.ContactName != nil && *u.ContactName != "" {
		completedFields++
	}
	if u.ContactPhone != nil && *u.ContactPhone != "" {
		completedFields++
	}

	return float64(completedFields) / float64(totalFields) * 100.0
}

func (u *UserAddressValues) UpdateCompletionRate() *UserAddressValues {
	rate := u.CalculateCompletionRate()
	u.CompletionRate = &rate
	return u
}

func (u *UserAddressValues) CalculateQualityScore() float64 {
	score := 0.0

	// 基础信息完整性 (40%)
	completionRate := u.CalculateCompletionRate()
	score += (completionRate / 100.0) * 2.0

	// 位置精确度 (30%)
	if u.Accuracy != nil {
		if *u.Accuracy <= 10.0 {
			score += 1.5
		} else if *u.Accuracy <= 50.0 {
			score += 1.0
		} else {
			score += 0.5
		}
	}

	// 验证状态 (20%)
	if u.GetIsVerified() {
		score += 1.0
	}

	// 使用频率 (10%)
	usageCount := u.GetUsageCount()
	if usageCount >= 10 {
		score += 0.5
	} else if usageCount >= 5 {
		score += 0.3
	} else if usageCount >= 1 {
		score += 0.1
	}

	// 确保评分在0-5范围内
	if score > 5.0 {
		score = 5.0
	}

	return score
}

func (u *UserAddressValues) UpdateQualityScore() *UserAddressValues {
	score := u.CalculateQualityScore()
	u.QualityScore = &score
	return u
}

// 使用统计更新
func (u *UserAddressValues) IncrementUsage() *UserAddressValues {
	count := u.GetUsageCount() + 1
	u.UsageCount = &count

	now := utils.TimeNowMilli()
	u.LastUsedAt = &now

	if u.FirstUsedAt == nil {
		u.FirstUsedAt = &now
	}

	// 更新质量评分
	u.UpdateQualityScore()

	return u
}

// 地址验证
func (u *UserAddressValues) ValidateAddress() []string {
	var errors []string

	if u.FormattedAddress == nil || *u.FormattedAddress == "" {
		errors = append(errors, "formatted_address is required")
	}

	if u.Latitude == nil || u.Longitude == nil {
		errors = append(errors, "latitude and longitude are required")
	} else {
		lat := *u.Latitude
		lng := *u.Longitude
		if lat < -90 || lat > 90 {
			errors = append(errors, "latitude must be between -90 and 90")
		}
		if lng < -180 || lng > 180 {
			errors = append(errors, "longitude must be between -180 and 180")
		}
	}

	if u.City == nil || *u.City == "" {
		errors = append(errors, "city is required")
	}

	if u.Country == nil || *u.Country == "" {
		errors = append(errors, "country is required")
	}

	return errors
}

func (u *UserAddressValues) IsValid() bool {
	errors := u.ValidateAddress()
	return len(errors) == 0
}

// 地址距离计算
func (u *UserAddressValues) CalculateDistanceTo(lat, lng float64) float64 {
	if u.Latitude == nil || u.Longitude == nil {
		return -1
	}

	return haversineDistance(*u.Latitude, *u.Longitude, lat, lng)
}

func (u *UserAddressValues) IsNearby(lat, lng, radiusKm float64) bool {
	distance := u.CalculateDistanceTo(lat, lng)
	return distance >= 0 && distance <= radiusKm
}

// 地标管理
func (u *UserAddressValues) SetNearbyLandmarks(landmarks []string) error {
	landmarksJSON, err := utils.ToJSON(landmarks)
	if err != nil {
		return fmt.Errorf("failed to marshal landmarks: %v", err)
	}

	u.NearbyLandmarks = &landmarksJSON
	return nil
}

func (u *UserAddressValues) GetNearbyLandmarks() []string {
	if u.NearbyLandmarks == nil {
		return []string{}
	}

	var landmarks []string
	if err := utils.FromJSON(*u.NearbyLandmarks, &landmarks); err != nil {
		return []string{}
	}

	return landmarks
}

func (u *UserAddressValues) AddLandmark(landmark string) error {
	landmarks := u.GetNearbyLandmarks()

	// 避免重复
	for _, existing := range landmarks {
		if existing == landmark {
			return nil
		}
	}

	landmarks = append(landmarks, landmark)
	return u.SetNearbyLandmarks(landmarks)
}

// 标签管理
func (u *UserAddressValues) AddTag(tag string) *UserAddressValues {
	var tags []string
	if u.Tags != nil && *u.Tags != "" {
		tags = strings.Split(*u.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return u
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	u.Tags = &tagsStr
	return u
}

func (u *UserAddressValues) HasTag(tag string) bool {
	if u.Tags == nil || *u.Tags == "" {
		return false
	}

	tags := strings.Split(*u.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 便捷创建方法
func NewHomeAddress(userID, label string, lat, lng float64, address string) *UserAddress {
	addr := NewUserAddressV2()
	addr.SetUserID(userID).
		SetAddressType(AddressTypeHome).
		SetLabel(label).
		SetLocation(lat, lng).
		SetFormattedAddress(address).
		SetDefault(false).
		SetPersonalization(label, IconTypeHome, "#4CAF50", 1)

	addr.AddTag("home")
	addr.UpdateCompletionRate()
	addr.UpdateQualityScore()

	return addr
}

func NewWorkAddress(userID, label string, lat, lng float64, address string) *UserAddress {
	addr := NewUserAddressV2()
	addr.SetUserID(userID).
		SetAddressType(AddressTypeWork).
		SetLabel(label).
		SetLocation(lat, lng).
		SetFormattedAddress(address).
		SetDefault(false).
		SetPersonalization(label, IconTypeWork, "#FF9800", 2)

	addr.AddTag("work")
	addr.UpdateCompletionRate()
	addr.UpdateQualityScore()

	return addr
}

func NewPickupAddress(userID string, lat, lng float64, address string) *UserAddress {
	addr := NewUserAddressV2()
	addr.SetUserID(userID).
		SetAddressType(AddressTypePickup).
		SetLocation(lat, lng).
		SetFormattedAddress(address).
		SetPersonalization("Pickup", IconTypeStar, "#2196F3", 10)

	addr.AddTag("pickup")
	addr.UpdateCompletionRate()
	addr.UpdateQualityScore()

	return addr
}

func NewDropoffAddress(userID string, lat, lng float64, address string) *UserAddress {
	addr := NewUserAddressV2()
	addr.SetUserID(userID).
		SetAddressType(AddressTypeDropoff).
		SetLocation(lat, lng).
		SetFormattedAddress(address).
		SetPersonalization("Dropoff", IconTypeStar, "#F44336", 10)

	addr.AddTag("dropoff")
	addr.UpdateCompletionRate()
	addr.UpdateQualityScore()

	return addr
}
