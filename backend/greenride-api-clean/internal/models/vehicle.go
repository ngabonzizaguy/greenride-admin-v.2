package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"time"
)

// Vehicle 车辆表
type Vehicle struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	VehicleID string `json:"vehicle_id" gorm:"column:vehicle_id;type:varchar(64);uniqueIndex"`
	Salt      string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*VehicleValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type VehicleValues struct {
	DriverID *string `json:"driver_id" gorm:"column:driver_id;type:varchar(64);index"`

	// 基本信息
	Brand       *string `json:"brand" gorm:"column:brand;type:varchar(100)"`
	Model       *string `json:"model" gorm:"column:model;type:varchar(100)"`
	Year        *int    `json:"year" gorm:"column:year;type:int"`
	Color       *string `json:"color" gorm:"column:color;type:varchar(50)"`
	PlateNumber *string `json:"plate_number" gorm:"column:plate_number;type:varchar(20);index"`
	VIN         *string `json:"vin" gorm:"column:vin;type:varchar(50);index"` // Vehicle Identification Number

	// 车辆类型和规格
	TypeID       *string `json:"type_id" gorm:"column:type_id;type:varchar(64);index"`   // 关联VehicleType表的ID
	Category     *string `json:"category" gorm:"column:category;type:varchar(50);index"` // sedan, suv, mpv, van, hatchback
	Level        *string `json:"level" gorm:"column:level;type:varchar(50);index"`       // economy, comfort, premium, luxury
	SeatCapacity *int    `json:"seat_capacity" gorm:"column:seat_capacity;default:4"`
	FuelType     *string `json:"fuel_type" gorm:"column:fuel_type;type:varchar(20)"`       // gasoline, diesel, electric, hybrid
	Transmission *string `json:"transmission" gorm:"column:transmission;type:varchar(20)"` // manual, automatic
	EngineSize   *string `json:"engine_size" gorm:"column:engine_size;type:varchar(20)"`

	// 车辆状态
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`                 // active, inactive, maintenance, retired
	VerifyStatus *string `json:"verify_status" gorm:"column:verify_status;type:varchar(32);index;default:'verified'"` // unverified, pending, active, inactive, maintenance, suspended, banned, retired

	// 注册和保险信息
	RegistrationNumber    *string `json:"registration_number" gorm:"column:registration_number;type:varchar(50)"`
	RegistrationExpiry    *int64  `json:"registration_expiry" gorm:"column:registration_expiry"`
	InsuranceCompany      *string `json:"insurance_company" gorm:"column:insurance_company;type:varchar(100)"`
	InsurancePolicyNumber *string `json:"insurance_policy_number" gorm:"column:insurance_policy_number;type:varchar(100)"`
	InsuranceExpiry       *int64  `json:"insurance_expiry" gorm:"column:insurance_expiry"`

	// 位置信息
	CurrentLatitude   *float64 `json:"current_latitude" gorm:"column:current_latitude;type:decimal(10,8)"`
	CurrentLongitude  *float64 `json:"current_longitude" gorm:"column:current_longitude;type:decimal(11,8)"`
	LocationUpdatedAt *int64   `json:"location_updated_at" gorm:"column:location_updated_at"`

	// 维护信息

	// 车辆特性

	// 文档信息
	Photos    []string `json:"photos" gorm:"column:photos;type:json;serializer:json"`       // JSON array of photo URLs
	Documents []string `json:"documents" gorm:"column:documents;type:json;serializer:json"` // JSON array of document URLs

	// 评价和统计
	Rating *float64 `json:"rating" gorm:"column:rating;type:decimal(3,2);default:5.0"`

	// 验证状态

	// 时间戳
	VerifiedAt *int64 `json:"verified_at" gorm:"column:verified_at"`
	LastUsedAt *int64 `json:"last_used_at" gorm:"column:last_used_at"`

	// 扩展信息
	Notes    *string        `json:"notes" gorm:"column:notes;type:text"`
	Metadata map[string]any `json:"metadata" gorm:"column:metadata;type:json;serializer:json"` // JSON格式的额外信息

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Vehicle) TableName() string {
	return "t_vehicles"
}

// 创建新的车辆对象
func NewVehicle() *Vehicle {
	return &Vehicle{
		VehicleID: utils.GenerateVehicleID(),
		Salt:      utils.GenerateSalt(),
		VehicleValues: &VehicleValues{
			Status:       utils.StringPtr(protocol.StatusActive),
			VerifyStatus: utils.StringPtr(protocol.StatusVerified),
			SeatCapacity: utils.IntPtr(4),
			Rating:       utils.Float64Ptr(5.0),
			Photos:       []string{},
			Documents:    []string{},
			Metadata:     make(map[string]any),
		},
	}
}

// SetValues 更新VehicleValues中的非nil值
func (v *VehicleValues) SetValues(values *VehicleValues) {
	if values == nil {
		return
	}

	if values.DriverID != nil {
		v.DriverID = values.DriverID
	}
	if values.Brand != nil {
		v.Brand = values.Brand
	}
	if values.Model != nil {
		v.Model = values.Model
	}
	if values.Year != nil {
		v.Year = values.Year
	}
	if values.Color != nil {
		v.Color = values.Color
	}
	if values.PlateNumber != nil {
		v.PlateNumber = values.PlateNumber
	}
	if values.VIN != nil {
		v.VIN = values.VIN
	}
	if values.TypeID != nil {
		v.TypeID = values.TypeID
	}
	if values.Category != nil {
		v.Category = values.Category
	}
	if values.Level != nil {
		v.Level = values.Level
	}
	if values.SeatCapacity != nil {
		v.SeatCapacity = values.SeatCapacity
	}
	if values.FuelType != nil {
		v.FuelType = values.FuelType
	}
	if values.Status != nil {
		v.Status = values.Status
	}
	if values.VerifyStatus != nil {
		v.VerifyStatus = values.VerifyStatus
	}
	if values.Metadata != nil {
		v.Metadata = values.Metadata
	}
	if values.UpdatedAt > 0 {
		v.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (v *VehicleValues) GetStatus() string {
	if v.Status == nil {
		return protocol.StatusActive
	}
	return *v.Status
}

func (v *VehicleValues) GetVerifyStatus() string {
	if v.VerifyStatus == nil {
		return protocol.StatusUnverified
	}
	return *v.VerifyStatus
}

func (v *VehicleValues) IsAvailable() bool {
	return v.IsActive() && v.IsVerified()
}

func (v *VehicleValues) IsVerified() bool {
	return v.GetVerifyStatus() == protocol.StatusVerified
}

func (v *VehicleValues) IsActive() bool {
	return v.GetStatus() == protocol.StatusActive
}

func (v *VehicleValues) GetDriverID() string {
	if v.DriverID == nil {
		return ""
	}
	return *v.DriverID
}

func (v *VehicleValues) GetBrand() string {
	if v.Brand == nil {
		return ""
	}
	return *v.Brand
}

func (v *VehicleValues) GetModel() string {
	if v.Model == nil {
		return ""
	}
	return *v.Model
}

func (v *VehicleValues) GetPlateNumber() string {
	if v.PlateNumber == nil {
		return ""
	}
	return *v.PlateNumber
}

func (v *VehicleValues) GetTypeID() string {
	if v.TypeID == nil {
		return ""
	}
	return *v.TypeID
}

func (v *VehicleValues) GetCategory() string {
	if v.Category == nil {
		return "sedan"
	}
	return *v.Category
}

func (v *VehicleValues) GetLevel() string {
	if v.Level == nil {
		return "economy"
	}
	return *v.Level
}

func (v *VehicleValues) GetSeatCapacity() int {
	if v.SeatCapacity == nil {
		return 4
	}
	return *v.SeatCapacity
}

func (v *VehicleValues) GetRating() float64 {
	if v.Rating == nil {
		return 5.0
	}
	return *v.Rating
}

// 添加缺少的 Getter 方法
func (v *VehicleValues) GetYear() int {
	if v.Year == nil {
		return 0
	}
	return *v.Year
}

func (v *VehicleValues) GetColor() string {
	if v.Color == nil {
		return ""
	}
	return *v.Color
}

func (v *VehicleValues) GetVIN() string {
	if v.VIN == nil {
		return ""
	}
	return *v.VIN
}

func (v *VehicleValues) GetFuelType() string {
	if v.FuelType == nil {
		return ""
	}
	return *v.FuelType
}

func (v *VehicleValues) GetTransmission() string {
	if v.Transmission == nil {
		return ""
	}
	return *v.Transmission
}

func (v *VehicleValues) GetEngineSize() string {
	if v.EngineSize == nil {
		return ""
	}
	return *v.EngineSize
}

func (v *VehicleValues) GetRegistrationNumber() string {
	if v.RegistrationNumber == nil {
		return ""
	}
	return *v.RegistrationNumber
}

func (v *VehicleValues) GetRegistrationExpiry() int64 {
	if v.RegistrationExpiry == nil {
		return 0
	}
	return *v.RegistrationExpiry
}

func (v *VehicleValues) GetInsuranceCompany() string {
	if v.InsuranceCompany == nil {
		return ""
	}
	return *v.InsuranceCompany
}

func (v *VehicleValues) GetInsurancePolicyNumber() string {
	if v.InsurancePolicyNumber == nil {
		return ""
	}
	return *v.InsurancePolicyNumber
}

func (v *VehicleValues) GetInsuranceExpiry() int64 {
	if v.InsuranceExpiry == nil {
		return 0
	}
	return *v.InsuranceExpiry
}

func (v *VehicleValues) GetCurrentLatitude() float64 {
	if v.CurrentLatitude == nil {
		return 0.0
	}
	return *v.CurrentLatitude
}

func (v *VehicleValues) GetCurrentLongitude() float64 {
	if v.CurrentLongitude == nil {
		return 0.0
	}
	return *v.CurrentLongitude
}

func (v *VehicleValues) GetLocationUpdatedAt() int64 {
	if v.LocationUpdatedAt == nil {
		return 0
	}
	return *v.LocationUpdatedAt
}

func (v *VehicleValues) GetPhotos() []string {
	if v.Photos == nil {
		return []string{}
	}
	return v.Photos
}

func (v *VehicleValues) GetDocuments() []string {
	if v.Documents == nil {
		return []string{}
	}
	return v.Documents
}

func (v *VehicleValues) GetVerifiedAt() int64 {
	if v.VerifiedAt == nil {
		return 0
	}
	return *v.VerifiedAt
}

func (v *VehicleValues) GetLastUsedAt() int64 {
	if v.LastUsedAt == nil {
		return 0
	}
	return *v.LastUsedAt
}

func (v *VehicleValues) GetNotes() string {
	if v.Notes == nil {
		return ""
	}
	return *v.Notes
}

func (v *VehicleValues) GetMetadata() map[string]interface{} {
	if v.Metadata == nil {
		return make(map[string]interface{})
	}
	return v.Metadata
}

// Setter 方法
func (v *VehicleValues) SetStatus(status string) *VehicleValues {
	v.Status = &status
	return v
}

func (v *VehicleValues) SetVerifyStatus(status string) *VehicleValues {
	v.VerifyStatus = &status
	return v
}

func (v *VehicleValues) SetDriver(driverID string) *VehicleValues {
	if driverID == "" {
		v.DriverID = nil
		return v
	}
	v.DriverID = &driverID
	return v
}

func (v *VehicleValues) SetAvailable(available bool) *VehicleValues {
	if available {
		v.VerifyStatus = utils.StringPtr(protocol.StatusActive)
	} else {
		v.VerifyStatus = utils.StringPtr(protocol.StatusInactive)
	}
	return v
}

func (v *VehicleValues) SetVehicleInfo(brand, model, plateNumber string, year int) *VehicleValues {
	v.Brand = &brand
	v.Model = &model
	v.PlateNumber = &plateNumber
	v.Year = &year
	return v
}

func (v *VehicleValues) SetLocation(lat, lng float64) *VehicleValues {
	v.CurrentLatitude = &lat
	v.CurrentLongitude = &lng
	now := utils.TimeNowMilli()
	v.LocationUpdatedAt = &now
	return v
}

// 添加缺少的 Setter 方法
func (v *VehicleValues) SetYear(year int) *VehicleValues {
	v.Year = &year
	return v
}

func (v *VehicleValues) SetColor(color string) *VehicleValues {
	v.Color = &color
	return v
}

func (v *VehicleValues) SetVIN(vin string) *VehicleValues {
	v.VIN = &vin
	return v
}

func (v *VehicleValues) SetTypeID(typeID string) *VehicleValues {
	if typeID == "" {
		v.TypeID = nil
	} else {
		v.TypeID = &typeID
	}
	return v
}

// HasTypeID 检查车辆是否关联了VehicleType
func (v *VehicleValues) HasTypeID() bool {
	return v.TypeID != nil && *v.TypeID != ""
}

// ClearTypeID 清空车辆类型关联
func (v *VehicleValues) ClearTypeID() *VehicleValues {
	v.TypeID = nil
	return v
}

func (v *VehicleValues) SetCategory(category string) *VehicleValues {
	v.Category = &category
	return v
}

func (v *VehicleValues) SetLevel(level string) *VehicleValues {
	v.Level = &level
	return v
}

func (v *VehicleValues) SetSeatCapacity(capacity int) *VehicleValues {
	v.SeatCapacity = &capacity
	return v
}

func (v *VehicleValues) SetFuelType(fuelType string) *VehicleValues {
	v.FuelType = &fuelType
	return v
}

func (v *VehicleValues) SetTransmission(transmission string) *VehicleValues {
	v.Transmission = &transmission
	return v
}

func (v *VehicleValues) SetEngineSize(engineSize string) *VehicleValues {
	v.EngineSize = &engineSize
	return v
}

func (v *VehicleValues) SetVerified(verified bool) *VehicleValues {
	if verified {
		v.VerifyStatus = utils.StringPtr(protocol.StatusActive)
		now := utils.TimeNowMilli()
		v.VerifiedAt = &now
	} else {
		v.VerifyStatus = utils.StringPtr(protocol.StatusUnverified)
	}
	return v
}

func (v *VehicleValues) SetActive(active bool) *VehicleValues {
	if active {
		v.VerifyStatus = utils.StringPtr(protocol.StatusActive)
	} else {
		v.VerifyStatus = utils.StringPtr(protocol.StatusInactive)
	}
	return v
}

func (v *VehicleValues) SetRegistrationInfo(number string, expiry int64) *VehicleValues {
	v.RegistrationNumber = &number
	v.RegistrationExpiry = &expiry
	return v
}

func (v *VehicleValues) SetInsuranceInfo(company, policyNumber string, expiry int64) *VehicleValues {
	v.InsuranceCompany = &company
	v.InsurancePolicyNumber = &policyNumber
	v.InsuranceExpiry = &expiry
	return v
}

// 单独的特性setter方法
// 其他缺失的setter方法
func (v *VehicleValues) SetBrand(brand string) *VehicleValues {
	v.Brand = &brand
	return v
}

func (v *VehicleValues) SetModel(model string) *VehicleValues {
	v.Model = &model
	return v
}

func (v *VehicleValues) SetPlateNumber(plateNumber string) *VehicleValues {
	v.PlateNumber = &plateNumber
	return v
}

func (v *VehicleValues) SetRegistrationNumber(regNumber string) *VehicleValues {
	v.RegistrationNumber = &regNumber
	return v
}

func (v *VehicleValues) SetInsuranceCompany(company string) *VehicleValues {
	v.InsuranceCompany = &company
	return v
}

func (v *VehicleValues) SetInsurancePolicyNumber(policyNumber string) *VehicleValues {
	v.InsurancePolicyNumber = &policyNumber
	return v
}

func (v *VehicleValues) SetNotes(notes string) *VehicleValues {
	v.Notes = &notes
	return v
}

func (v *VehicleValues) SetPhotos(photos []string) *VehicleValues {
	v.Photos = photos
	return v
}

func (v *VehicleValues) SetDocuments(documents []string) *VehicleValues {
	v.Documents = documents
	return v
}

func (v *VehicleValues) SetRating(rating float64) *VehicleValues {
	v.Rating = &rating
	return v
}

func (v *VehicleValues) SetMetadata(metadata map[string]interface{}) *VehicleValues {
	v.Metadata = metadata
	return v
}

func (v *VehicleValues) SetLastUsed() *VehicleValues {
	now := utils.TimeNowMilli()
	v.LastUsedAt = &now
	return v
}

// 业务方法
func (v *Vehicle) CanBeUsed() bool {
	return v.IsAvailable() &&
		!v.IsRegistrationExpired() && !v.IsInsuranceExpired()
}

func (v *Vehicle) IsRegistrationExpired() bool {
	if v.RegistrationExpiry == nil {
		return false
	}
	return *v.RegistrationExpiry < utils.TimeNowMilli()
}

func (v *Vehicle) IsInsuranceExpired() bool {
	if v.InsuranceExpiry == nil {
		return false
	}
	return *v.InsuranceExpiry < utils.TimeNowMilli()
}

func (v *VehicleValues) MarkAsVerified() {
	v.VerifyStatus = utils.StringPtr(protocol.StatusVerified)
	now := utils.TimeNowMilli()
	v.VerifiedAt = &now
}

func (v *VehicleValues) UpdateLocation(lat, lng float64) {
	v.CurrentLatitude = &lat
	v.CurrentLongitude = &lng
	now := utils.TimeNowMilli()
	v.LocationUpdatedAt = &now
}

func (v *VehicleValues) IncrementRideCount() {
	now := utils.TimeNowMilli()
	v.LastUsedAt = &now
}

func (v *VehicleValues) UpdateRating(newRating float64) {
	v.Rating = &newRating
}

func (v *VehicleValues) SetMaintenanceMode() {
	v.SetStatus(protocol.StatusMaintenance)
	v.SetVerifyStatus(protocol.StatusMaintenance)
}

func (v *VehicleValues) SetActiveMode() {
	v.SetStatus(protocol.StatusActive)
	v.SetVerifyStatus(protocol.StatusActive)
}

func (v *VehicleValues) GetDisplayName() string {
	brand := v.GetBrand()
	model := v.GetModel()
	if brand == "" && model == "" {
		return "Unknown Vehicle"
	}
	if brand == "" {
		return model
	}
	if model == "" {
		return brand
	}
	return brand + " " + model
}

func GetVehicleByID(vehicleID string) *Vehicle {
	var vehicle Vehicle
	err := DB.Model(&vehicle).Where("vehicle_id = ?", vehicleID).First(&vehicle).Error
	if err != nil {
		return nil
	}
	return &vehicle
}

func GetVehicleByDriverID(driverID string) *Vehicle {
	var vehicle Vehicle
	err := DB.Model(&vehicle).Where("driver_id = ?", driverID).First(&vehicle).Error
	if err != nil {
		return nil
	}
	return &vehicle
}

type Vehicles []*Vehicle

// Protocol 转换为protocol.Vehicle列表
func (vs Vehicles) Protocol() []*protocol.Vehicle {
	list := make([]*protocol.Vehicle, 0, len(vs))
	for _, v := range vs {
		list = append(list, v.Protocol())
	}
	return list
}

// Protocol 将Vehicle转换为protocol.Vehicle
func (v *Vehicle) Protocol() *protocol.Vehicle {
	vehicle := &protocol.Vehicle{
		VehicleID:    v.VehicleID,
		DriverID:     v.GetDriverID(),
		Brand:        v.GetBrand(),
		Model:        v.GetModel(),
		Color:        v.GetColor(),
		PlateNumber:  v.GetPlateNumber(),
		TypeID:       v.GetTypeID(),
		Category:     v.GetCategory(),
		Level:        v.GetLevel(),
		SeatCapacity: v.GetSeatCapacity(),
		Status:       v.GetStatus(),
		VerifyStatus: v.GetVerifyStatus(),
		Rating:       v.GetRating(),
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}

	// 处理可选字段 - 使用 getter 方法避免 nil 检查
	if year := v.GetYear(); year > 0 {
		vehicle.Year = year
	}
	if vin := v.GetVIN(); vin != "" {
		vehicle.VIN = vin
	}
	if fuelType := v.GetFuelType(); fuelType != "" {
		vehicle.FuelType = fuelType
	}
	if transmission := v.GetTransmission(); transmission != "" {
		vehicle.Transmission = transmission
	}
	if engineSize := v.GetEngineSize(); engineSize != "" {
		vehicle.EngineSize = engineSize
	}

	// 注册和保险信息
	if regNumber := v.GetRegistrationNumber(); regNumber != "" {
		vehicle.RegistrationNumber = regNumber
	}
	if regExpiry := v.GetRegistrationExpiry(); regExpiry > 0 {
		expiry := utils.MilliToTime(regExpiry).Format(time.DateOnly)
		vehicle.RegistrationExpiry = &expiry
	}
	if insCompany := v.GetInsuranceCompany(); insCompany != "" {
		vehicle.InsuranceCompany = insCompany
	}
	if insPolicyNumber := v.GetInsurancePolicyNumber(); insPolicyNumber != "" {
		vehicle.InsurancePolicyNumber = insPolicyNumber
	}
	if insExpiry := v.GetInsuranceExpiry(); insExpiry > 0 {
		expiry := utils.MilliToTime(insExpiry).Format(time.DateOnly)
		vehicle.InsuranceExpiry = &expiry
	}

	// 位置信息
	if lat := v.GetCurrentLatitude(); lat != 0 {
		vehicle.CurrentLatitude = &lat
	}
	if lng := v.GetCurrentLongitude(); lng != 0 {
		vehicle.CurrentLongitude = &lng
	}
	if locUpdated := v.GetLocationUpdatedAt(); locUpdated > 0 {
		vehicle.LocationUpdatedAt = &locUpdated
	}

	// 维护信息 - 已移除

	// 车辆特性 - 已移除

	// 文档和图片 - 直接使用 getter 方法
	vehicle.Photos = v.GetPhotos()
	vehicle.Documents = v.GetDocuments()

	// 验证状态 - 已移除

	if verifiedAt := v.GetVerifiedAt(); verifiedAt > 0 {
		vehicle.VerifiedAt = verifiedAt
	}

	// 使用信息
	if lastUsed := v.GetLastUsedAt(); lastUsed > 0 {
		vehicle.LastUsedAt = lastUsed
	}
	if notes := v.GetNotes(); notes != "" {
		vehicle.Notes = notes
	}

	return vehicle
}

func FindDriversByVehicle(category, level string) []string {
	var driver_list []string
	query := GetDB().Model(&Vehicle{}).Select([]string{"driver_id"}).Where("driver_id is not null and driver_id !=''")
	if category != "" {
		query = query.Where("category=?", category)
	}
	if level != "" {
		query = query.Where("level=?", level)
	}
	if err := query.Find(&driver_list).Error; err != nil {
		return []string{}
	}
	return driver_list
}

// FindVehicleByPlateNumber 根据车牌号查找车辆
func FindVehicleByPlateNumber(plateNumber string) *Vehicle {
	var vehicle Vehicle
	if err := GetDB().Where("plate_number = ?", plateNumber).First(&vehicle).Error; err != nil {
		return nil
	}
	return &vehicle
}

// FindVehicleByVIN 根据VIN码查找车辆
func FindVehicleByVIN(vin string) *Vehicle {
	var vehicle Vehicle
	if err := GetDB().Where("vin = ?", vin).First(&vehicle).Error; err != nil {
		return nil
	}
	return &vehicle
}

// FindVehicleByID 根据车辆ID查找车辆
func FindVehicleByID(vehicleID string) *Vehicle {
	var vehicle Vehicle
	if err := GetDB().Where("vehicle_id = ?", vehicleID).First(&vehicle).Error; err != nil {
		return nil
	}
	return &vehicle
}
