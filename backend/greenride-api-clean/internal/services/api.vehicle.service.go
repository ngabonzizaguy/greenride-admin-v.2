package services

import (
	"errors"
	"sync"

	"greenride/internal/models"
	"greenride/internal/protocol"

	"gorm.io/gorm"
)

// VehicleService 车辆服务
type VehicleService struct {
	db *gorm.DB
}

var (
	vehicleServiceInstance *VehicleService
	vehicleServiceOnce     sync.Once
)

// GetVehicleService 获取车辆服务单例
func GetVehicleService() *VehicleService {
	vehicleServiceOnce.Do(func() {
		SetupVehicleService()
	})
	return vehicleServiceInstance
}

// SetupVehicleService 设置车辆服务
func SetupVehicleService() {
	vehicleServiceInstance = &VehicleService{
		db: models.GetDB(),
	}
}

// NewVehicleService 创建车辆服务实例（用于测试）
func NewVehicleService() *VehicleService {
	return &VehicleService{
		db: models.GetDB(),
	}
}

// GetVehicleByID 根据ID获取车辆
func (s *VehicleService) GetVehicleByID(vehicleID string) (*models.Vehicle, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle_id is required")
	}

	var vehicle models.Vehicle
	err := s.db.Where("vehicle_id = ?", vehicleID).First(&vehicle).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &vehicle, nil
}

// GetVehicleByDriverID 根据司机ID获取车辆
func (s *VehicleService) GetVehicleByDriverID(driverID string) *models.Vehicle {
	if driverID == "" {
		return nil
	}

	var vehicle models.Vehicle
	err := s.db.Where("driver_id = ?", driverID).First(&vehicle).Error
	if err != nil {
		return nil
	}
	return &vehicle
}

// GetVehiclesByOwnerID 根据车主ID获取车辆列表
func (s *VehicleService) GetVehiclesByOwnerID(ownerID string) ([]*models.Vehicle, error) {
	if ownerID == "" {
		return nil, errors.New("owner_id is required")
	}

	var vehicles []*models.Vehicle
	err := s.db.Where("owner_id = ?", ownerID).Find(&vehicles).Error
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

// CreateVehicle 创建车辆
func (s *VehicleService) CreateVehicle(vehicle *models.Vehicle) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}

	if vehicle.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}

	return s.db.Create(vehicle).Error
}

// UpdateVehicle 更新车辆
func (s *VehicleService) UpdateVehicle(vehicle *models.Vehicle, values *models.VehicleValues) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}

	// 直接更新VehicleV2，GORM会自动处理嵌入的VehicleV2Values
	result := s.db.Model(vehicle).Where("id = ?", vehicle.ID).Updates(values)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteVehicle 删除车辆（软删除）
func (s *VehicleService) DeleteVehicle(vehicleID string) error {
	if vehicleID == "" {
		return errors.New("vehicle_id is required")
	}

	return s.db.Where("vehicle_id = ?", vehicleID).Delete(&models.Vehicle{}).Error
}

// GetActiveVehicles 获取激活状态的车辆列表
func (s *VehicleService) GetActiveVehicles() ([]*models.Vehicle, error) {
	var vehicles []*models.Vehicle
	err := s.db.Where("is_active = ?", true).Find(&vehicles).Error
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

// GetAvailableVehicles 获取可用车辆列表
func (s *VehicleService) GetAvailableVehicles() ([]*models.Vehicle, error) {
	var vehicles []*models.Vehicle
	err := s.db.Where("is_available = ? AND is_active = ?", true, true).Find(&vehicles).Error
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (s *VehicleService) GetVehicles(req *protocol.GetVehiclesRequest) (list []*protocol.Vehicle, total int64) {
	var vehicles models.Vehicles
	query := models.GetDB().Model(&models.Vehicle{}).Where("status = ?", protocol.StatusActive)
	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0
	}
	if req.Limit > 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Offset(offset).Limit(req.Limit)
	}
	// 获取车辆列表
	var vehicleList []string
	if err := query.Select([]string{"vehicle_id"}).Order("created_at DESC").Find(&vehicleList).Error; err != nil {
		return nil, 0
	}
	// 根据车辆ID列表获取完整的车辆信息
	if len(vehicleList) > 0 {
		if err := s.db.Where("vehicle_id IN ?", vehicleList).Find(&vehicles).Error; err != nil {
			return nil, 0
		}
	}
	for _, item := range vehicles {
		list = append(list, s.GetVehicleInfo(item))
	}
	return
}

func (s *VehicleService) GetVehicleInfoByID(vehicleID string) *protocol.Vehicle {
	vehicle := models.GetVehicleByID(vehicleID)
	return s.GetVehicleInfo(vehicle)
}

func (s VehicleService) GetVehicleInfo(vehicle *models.Vehicle) *protocol.Vehicle {
	if vehicle == nil {
		return nil
	}
	info := vehicle.Protocol()
	if info.DriverID != "" {
		driver := models.GetUserByID(vehicle.GetDriverID())
		if driver != nil {
			info.Driver = driver.Protocol()
		}
	}
	return info
}
