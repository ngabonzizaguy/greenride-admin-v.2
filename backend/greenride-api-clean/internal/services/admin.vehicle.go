package services

import (
	"errors"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"log"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type AdminVehicleService struct {
}

var (
	adminVehicleInstance *AdminVehicleService
	adminVehicleOnce     sync.Once
)

func GetAdminVehicleService() *AdminVehicleService {
	adminVehicleOnce.Do(func() {
		SetupAdminVehicleService()
	})
	return adminVehicleInstance
}
func SetupAdminVehicleService() {
	adminVehicleInstance = &AdminVehicleService{}
}

// GetVehicleList 获取车辆列表
func (s *AdminVehicleService) GetVehicleList(req *protocol.VehicleListRequest) ([]*models.Vehicle, int64, protocol.ErrorCode) {
	var vehicles []*models.Vehicle
	var total int64

	query := models.GetDB().Model(&models.Vehicle{})

	// 车辆类型过滤
	if req.TypeID != "" {
		query = query.Where("type_id = ?", req.TypeID)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Level != "" {
		query = query.Where("level = ?", req.Level)
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 关键字搜索（只搜索车牌号和品牌）
	if req.Keyword != "" {
		searchTerm := "%" + req.Keyword + "%"
		query = query.Where("plate_number LIKE ? OR brand LIKE ?", searchTerm, searchTerm)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	// 分页查询，按创建时间倒序
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&vehicles).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	return vehicles, total, protocol.Success
}

// GetVehicleByID 根据ID获取车辆
func (s *AdminVehicleService) GetVehicleByID(vehicleID string) *models.Vehicle {
	return models.FindVehicleByID(vehicleID)
}

// UpdateVehicle 更新车辆信息
func (s *AdminVehicleService) UpdateVehicle(req *protocol.VehicleUpdateRequest) protocol.ErrorCode {
	vehicle := s.GetVehicleByID(req.VehicleID)
	if vehicle == nil {
		return protocol.VehicleNotFound
	}

	// 更新基础信息 - 使用链式调用
	if req.Brand != nil {
		vehicle.SetBrand(*req.Brand)
	}
	if req.Model != nil {
		vehicle.SetModel(*req.Model)
	}
	if req.Year != nil {
		vehicle.SetYear(*req.Year)
	}
	if req.Color != nil {
		vehicle.SetColor(*req.Color)
	}
	if req.PlateNumber != nil {
		vehicle.SetPlateNumber(*req.PlateNumber)
	}
	if req.VIN != nil {
		vehicle.SetVIN(*req.VIN)
	}

	// 更新车辆类型和配置 - 使用链式调用
	if req.TypeID != nil {
		vehicle.SetTypeID(*req.TypeID)
	}
	if req.Category != nil {
		vehicle.SetCategory(*req.Category)
	}
	if req.Level != nil {
		vehicle.SetLevel(*req.Level)
	}
	if req.SeatCapacity != nil {
		vehicle.SetSeatCapacity(*req.SeatCapacity)
	}
	if req.FuelType != nil {
		vehicle.SetFuelType(*req.FuelType)
	}
	if req.Transmission != nil {
		vehicle.SetTransmission(*req.Transmission)
	}
	if req.EngineSize != nil {
		vehicle.SetEngineSize(*req.EngineSize)
	}

	// 更新注册和保险信息 - 使用链式调用
	if req.RegistrationNumber != nil {
		vehicle.SetRegistrationNumber(*req.RegistrationNumber)
	}
	if req.RegistrationExpiry != nil {
		// 将字符串日期转换为Unix时间戳（这里需要根据实际日期格式进行转换）
		// 暂时跳过，需要具体的日期格式规范
		// vehicle.SetRegistrationExpiry(*req.RegistrationExpiry)
	}
	if req.InsuranceCompany != nil {
		vehicle.SetInsuranceCompany(*req.InsuranceCompany)
	}
	if req.InsurancePolicyNumber != nil {
		vehicle.SetInsurancePolicyNumber(*req.InsurancePolicyNumber)
	}
	if req.InsuranceExpiry != nil {
		// 将字符串日期转换为Unix时间戳（这里需要根据实际日期格式进行转换）
		// 暂时跳过，需要具体的日期格式规范
		// vehicle.SetInsuranceExpiry(*req.InsuranceExpiry)
	}

	// 更新维护信息 - 使用链式调用
	if req.LastServiceDate != nil {
		// 将字符串日期转换为Unix时间戳（这里需要根据实际日期格式进行转换）
		// 暂时跳过，需要具体的日期格式规范
		// vehicle.SetLastServiceDate(*req.LastServiceDate)
	}
	if req.NextServiceDue != nil {
		// 将字符串日期转换为Unix时间戳（这里需要根据实际日期格式进行转换）
		// 暂时跳过，需要具体的日期格式规范
		// vehicle.SetNextServiceDue(*req.NextServiceDue)
	}
	// TotalMileage field has been removed
	// if req.TotalMileage != nil {
	//	vehicle.SetTotalMileage(*req.TotalMileage)
	// }
	// ServiceMileage field has been removed
	// if req.ServiceMileage != nil {
	//	vehicle.SetServiceMileage(*req.ServiceMileage)
	// }

	// 更新车辆特性 - 使用链式调用
	// HasAirConditioner field has been removed
	// if req.HasAirConditioner != nil {
	//	vehicle.SetHasAirConditioner(*req.HasAirConditioner)
	// }
	// HasGPS field has been removed
	// if req.HasGPS != nil {
	//	vehicle.SetHasGPS(*req.HasGPS)
	// }
	// HasWiFi field has been removed
	// if req.HasWiFi != nil {
	//	vehicle.SetHasWiFi(*req.HasWiFi)
	// }
	// HasCharger field has been removed
	// if req.HasCharger != nil {
	//	vehicle.SetHasCharger(*req.HasCharger)
	// }
	// HasBluetooth field has been removed
	// if req.HasBluetooth != nil {
	//	vehicle.SetHasBluetooth(*req.HasBluetooth)
	// }

	// 更新文档和图片
	if len(req.Photos) > 0 {
		vehicle.SetPhotos(req.Photos)
	}
	if len(req.Documents) > 0 {
		vehicle.SetDocuments(req.Documents)
	}

	// 更新备注 - 使用链式调用
	if req.Notes != nil {
		vehicle.SetNotes(*req.Notes)
	}

	// 司机分配（可选）：为了保证一致性，在事务内处理：
	// - 解绑时写 NULL
	// - 绑定时校验司机存在且为 driver，并确保同一司机不会同时绑定多台车
	if err := models.GetDB().Transaction(func(tx *gorm.DB) error {
		if req.DriverID != nil {
			driverID := strings.TrimSpace(*req.DriverID)

			// Unassign
			if driverID == "" {
				vehicle.SetDriver("")
			} else {
				// Validate driver exists and is actually a driver
				driver := models.GetUserByID(driverID)
				if driver == nil {
					return errors.New(string(protocol.UserNotFound))
				}
				if driver.GetUserType() != protocol.UserTypeDriver {
					return errors.New(string(protocol.InvalidUserType))
				}

				// Ensure this driver is not assigned to another vehicle
				now := utils.TimeNowMilli()
				if err := tx.Model(&models.Vehicle{}).
					Where("driver_id = ? AND vehicle_id <> ?", driverID, vehicle.VehicleID).
					Updates(map[string]any{"driver_id": nil, "updated_at": now}).Error; err != nil {
					return err
				}

				vehicle.SetDriver(driverID)
			}
		}

		if err := tx.Save(vehicle).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		// Map known business errors surfaced via sentinel strings above
		if err.Error() == string(protocol.UserNotFound) {
			return protocol.UserNotFound
		}
		if err.Error() == string(protocol.InvalidUserType) {
			return protocol.InvalidUserType
		}
		return protocol.SystemError
	}

	return protocol.Success
}

// UpdateVehicleStatus 更新车辆状态
func (s *AdminVehicleService) UpdateVehicleStatus(vehicleID string, status *string, verifyStatus *string) protocol.ErrorCode {
	vehicle := s.GetVehicleByID(vehicleID)
	if vehicle == nil {
		return protocol.VehicleNotFound
	}

	// 更新状态 - 使用链式调用
	if status != nil {
		vehicle.SetStatus(*status)
	}
	if verifyStatus != nil {
		vehicle.SetVerifyStatus(*verifyStatus)
	}

	if err := models.GetDB().Save(vehicle).Error; err != nil {
		return protocol.SystemError
	}
	return protocol.Success
}

// SearchVehicles 搜索车辆（高级搜索）
func (s *AdminVehicleService) SearchVehicles(req *protocol.VehicleSearchRequest) ([]*protocol.Vehicle, int64, protocol.ErrorCode) {
	var vehicles []*models.Vehicle
	var total int64

	query := models.GetDB().Model(&models.Vehicle{})

	// 构建搜索条件（只搜索车牌号和品牌）
	if req.Keyword != "" {
		searchTerm := "%" + req.Keyword + "%"
		query = query.Where("plate_number LIKE ? OR brand LIKE ?",
			searchTerm, searchTerm)
	}

	if req.TypeID != "" {
		query = query.Where("type_id = ?", req.TypeID)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Level != "" {
		query = query.Where("level = ?", req.Level)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.DriverID != "" {
		query = query.Where("driver_id = ?", req.DriverID)
	}

	if req.IsVerified != nil {
		query = query.Where("is_verified = ?", *req.IsVerified)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if req.YearFrom != nil {
		query = query.Where("year >= ?", *req.YearFrom)
	}

	if req.YearTo != nil {
		query = query.Where("year <= ?", *req.YearTo)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	// 分页和排序
	page := 1
	limit := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	offset := (page - 1) * limit
	sortBy := "created_at DESC"

	if err := query.Order(sortBy).Offset(offset).Limit(limit).Find(&vehicles).Error; err != nil {
		return nil, 0, protocol.SystemError
	}
	var list []*protocol.Vehicle
	for _, vehicle := range vehicles {
		list = append(list, GetVehicleService().GetVehicleInfo(vehicle))
	}
	return list, total, protocol.Success
}

// GetDriverList 获取司机列表
func (s *AdminVehicleService) GetDriverList(keyword string, status string, verificationStatus string, page, limit int) ([]*models.User, int64, protocol.ErrorCode) {
	var drivers []*models.User
	var total int64

	query := models.GetDB().Model(&models.User{}).Where("user_type = ?", protocol.UserTypeDriver)

	// 关键字搜索（只搜索姓名和手机号）
	if keyword != "" {
		searchTerm := "%" + keyword + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR phone LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 认证状态过滤
	if verificationStatus == "verified" {
		query = query.Where("is_email_verified = ? AND is_phone_verified = ?", true, true)
	} else if verificationStatus == "unverified" {
		query = query.Where("is_email_verified = ? OR is_phone_verified = ?", false, false)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Failed to count drivers: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&drivers).Error; err != nil {
		log.Printf("Failed to get driver list: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	return drivers, total, protocol.Success
}

// GetDriverByID 获取司机详情
func (s *AdminVehicleService) GetDriverByID(driverID string) *models.User {
	var driver models.User
	if err := models.GetDB().Where("user_id = ? AND user_type = ?", driverID, protocol.UserTypeDriver).First(&driver).Error; err != nil {
		return nil
	}
	return &driver
}

// UpdateDriver 更新司机信息
func (s *AdminVehicleService) UpdateDriver(driverID string, updates map[string]interface{}) protocol.ErrorCode {
	// 验证司机是否存在
	var driver models.User
	if err := models.GetDB().Where("user_id = ? AND user_type = ?", driverID, protocol.UserTypeDriver).First(&driver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.UserNotFound
		}
		return protocol.SystemError
	}

	// 更新司机信息
	if err := models.GetDB().Model(&driver).Updates(updates).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// UpdateDriverStatus 更新司机状态
func (s *AdminVehicleService) UpdateDriverStatus(driverID string, status string, isActive *bool) protocol.ErrorCode {
	// 验证司机是否存在
	var driver models.User
	if err := models.GetDB().Where("user_id = ? AND user_type = ?", driverID, protocol.UserTypeDriver).First(&driver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.UserNotFound
		}
		return protocol.SystemError
	}

	updates := map[string]interface{}{
		"status": status,
	}

	if isActive != nil {
		updates["is_active"] = *isActive
	}

	if err := models.GetDB().Model(&driver).Updates(updates).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// VerifyDriver 审核司机认证
func (s *AdminVehicleService) VerifyDriver(driverID string, isEmailVerified, isPhoneVerified *bool, verifiedBy string) protocol.ErrorCode {
	// 验证司机是否存在
	var driver models.User
	if err := models.GetDB().Where("user_id = ? AND user_type = ?", driverID, protocol.UserTypeDriver).First(&driver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.UserNotFound
		}
		return protocol.SystemError
	}

	updates := map[string]interface{}{}

	if isEmailVerified != nil {
		updates["is_email_verified"] = *isEmailVerified
		if *isEmailVerified {
			updates["email_verified_at"] = utils.TimeNowMilli()
		}
	}

	if isPhoneVerified != nil {
		updates["is_phone_verified"] = *isPhoneVerified
		if *isPhoneVerified {
			updates["phone_verified_at"] = utils.TimeNowMilli()
		}
	}

	if err := models.GetDB().Model(&driver).Updates(updates).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// CreateVehicle 创建车辆
func (s *AdminVehicleService) CreateVehicle(req *protocol.VehicleCreateRequest) (*models.Vehicle, protocol.ErrorCode) {
	// 验证年份有效性
	if req.Year < 1900 || req.Year > 2030 {
		return nil, protocol.InvalidYear
	}

	// 验证车牌号是否已存在
	if existingVehicle := models.FindVehicleByPlateNumber(req.PlateNumber); existingVehicle != nil {
		return nil, protocol.PlateNumberExists
	}

	// 验证VIN是否已存在（如果提供）
	if req.VIN != "" {
		if existingVehicle := models.FindVehicleByVIN(req.VIN); existingVehicle != nil {
			return nil, protocol.VINExists
		}
	}
	if req.Color == "" {
		req.Color = "White" // 默认白色
	}

	// 创建新车辆
	vehicle := models.NewVehicle()

	// 设置基础信息
	vehicle.SetVehicleInfo(req.Brand, req.Model, req.PlateNumber, req.Year).
		SetColor(req.Color)

	if req.VIN != "" {
		vehicle.SetVIN(req.VIN)
	}

	// 设置车辆配置
	if req.Category != "" {
		vehicle.SetCategory(req.Category)
	}
	if req.Level != "" {
		vehicle.SetLevel(req.Level)
	}
	vehicle.SetSeatCapacity(req.SeatCapacity).
		SetFuelType(req.FuelType).
		SetTransmission(req.Transmission)

	if req.EngineSize != "" {
		vehicle.SetEngineSize(req.EngineSize)
	}

	// 设置注册和保险信息
	if req.RegistrationNumber != "" {
		if req.RegistrationExpiry != nil {
			// 这里应该解析日期字符串转换为时间戳，暂时跳过
			vehicle.RegistrationNumber = &req.RegistrationNumber
		} else {
			vehicle.RegistrationNumber = &req.RegistrationNumber
		}
	}

	if req.InsuranceCompany != "" {
		if req.InsuranceExpiry != nil {
			// 这里应该解析日期字符串转换为时间戳，暂时跳过
			vehicle.InsuranceCompany = &req.InsuranceCompany
			vehicle.InsurancePolicyNumber = &req.InsurancePolicyNumber
		} else {
			vehicle.InsuranceCompany = &req.InsuranceCompany
			vehicle.InsurancePolicyNumber = &req.InsurancePolicyNumber
		}
	}

	// 设置车辆特性
	// HasAirConditioner field has been removed
	// if req.HasAirConditioner != nil {
	//	vehicle.HasAirConditioner = req.HasAirConditioner
	// }
	// HasGPS field has been removed
	// if req.HasGPS != nil {
	//	vehicle.HasGPS = req.HasGPS
	// }
	// HasWiFi field has been removed
	// if req.HasWiFi != nil {
	//	vehicle.HasWiFi = req.HasWiFi
	// }
	// HasCharger field has been removed
	// if req.HasCharger != nil {
	//	vehicle.HasCharger = req.HasCharger
	// }
	// HasBluetooth field has been removed
	// if req.HasBluetooth != nil {
	//	vehicle.HasBluetooth = req.HasBluetooth
	// }

	// 设置文档和图片
	if len(req.Photos) > 0 {
		vehicle.SetPhotos(req.Photos)
	}
	if len(req.Documents) > 0 {
		vehicle.SetDocuments(req.Documents)
	}

	// 设置备注 - 使用链式调用
	if req.Notes != "" {
		vehicle.SetNotes(req.Notes)
	}

	// 设置创建时间
	vehicle.CreatedAt = utils.TimeNowMilli()

	// 保存到数据库
	if err := models.GetDB().Create(vehicle).Error; err != nil {
		log.Printf("Failed to create vehicle: %v", err)
		return nil, protocol.VehicleCreateFailed
	}

	return vehicle, protocol.Success
}

// DeleteVehicle 删除车辆（软删除）
func (s *AdminVehicleService) DeleteVehicle(vehicleID string) protocol.ErrorCode {
	vehicle := s.GetVehicleByID(vehicleID)
	if vehicle == nil {
		return protocol.VehicleNotFound
	}

	// 检查车辆是否正在使用中
	if vehicle.GetDriverID() != "" {
		return protocol.VehicleInUse
	}

	// 执行软删除
	if err := models.GetDB().Where("vehicle_id = ?", vehicleID).Delete(&models.Vehicle{}).Error; err != nil {
		log.Printf("Failed to delete vehicle %s: %v", vehicleID, err)
		return protocol.VehicleDeleteFailed
	}

	log.Printf("Vehicle %s deleted successfully", vehicleID)
	return protocol.Success
}

// SearchDrivers 搜索司机（高级搜索）
func (s *AdminVehicleService) SearchDrivers(params map[string]any) ([]*models.User, int64, error) {
	var drivers []*models.User
	var total int64

	query := models.GetDB().Model(&models.User{}).Where("user_type = ?", protocol.UserTypeDriver)

	// 构建搜索条件
	if keyword, ok := params["keyword"].(string); ok && keyword != "" {
		searchTerm := "%" + keyword + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR phone LIKE ? OR driver_license_number LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if status, ok := params["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if onlineStatus, ok := params["online_status"].(string); ok && onlineStatus != "" {
		query = query.Where("online_status = ?", onlineStatus)
	}

	if isEmailVerified, ok := params["is_email_verified"].(bool); ok {
		query = query.Where("is_email_verified = ?", isEmailVerified)
	}

	if isPhoneVerified, ok := params["is_phone_verified"].(bool); ok {
		query = query.Where("is_phone_verified = ?", isPhoneVerified)
	}

	if isActive, ok := params["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if minScore, ok := params["min_driver_score"].(float64); ok {
		query = query.Where("driver_score >= ?", minScore)
	}

	if maxScore, ok := params["max_driver_score"].(float64); ok {
		query = query.Where("driver_score <= ?", maxScore)
	}

	if minRides, ok := params["min_total_rides"].(int); ok {
		query = query.Where("total_rides_as_driver >= ?", minRides)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页和排序
	page := 1
	limit := 10
	if p, ok := params["page"].(int); ok && p > 0 {
		page = p
	}
	if l, ok := params["limit"].(int); ok && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit
	orderBy := "created_at DESC"
	if order, ok := params["order_by"].(string); ok && order != "" {
		orderBy = order
	}

	if err := query.Offset(offset).Limit(limit).Order(orderBy).Find(&drivers).Error; err != nil {
		return nil, 0, err
	}

	return drivers, total, nil
}
