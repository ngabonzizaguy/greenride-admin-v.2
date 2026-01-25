package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserLocationCache 用户位置缓存结构
type UserLocationCache struct {
	UserID       string  `json:"user_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	OnlineStatus string  `json:"online_status"`
	UpdatedAt    int64   `json:"updated_at"`
}

// UserService 用户服务
type UserService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

var (
	userServiceInstance *UserService
	userServiceOnce     sync.Once
)

// GetUserService 获取用户服务单例
func GetUserService() *UserService {
	userServiceOnce.Do(func() {
		SetupUserService()
	})
	return userServiceInstance
}

// SetupUserService 设置用户服务
func SetupUserService() {
	userServiceInstance = NewUserService()
}

// NewUserService 创建用户服务实例（用于测试）
func NewUserService() *UserService {
	return &UserService{
		db:          models.GetDB(),
		redisClient: models.Redis,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) protocol.ErrorCode {
	if user == nil {
		return protocol.InvalidParams
	}

	// 检查必填字段
	if user.UserID == "" {
		return protocol.InvalidParams
	}

	userType := user.GetUserType()

	// 检查邮箱和用户类型组合是否已存在
	if user.GetEmail() != "" && s.IsEmailExistsWithType(user.GetEmail(), userType) {
		return protocol.EmailAlreadyExists
	}

	// 检查手机号和用户类型组合是否已存在
	if user.GetPhone() != "" && s.IsPhoneExistsWithType(user.GetPhone(), userType) {
		return protocol.PhoneAlreadyExists
	}

	if err := s.db.Create(user).Error; err != nil {
		return protocol.SystemError
	}
	return protocol.Success
}

// GetUserByID 根据用户ID获取用户
func (s *UserService) GetUserByID(userID string) *models.User {
	var user models.User
	err := s.db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*models.User, protocol.ErrorCode) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, protocol.Success
		}
		return nil, protocol.SystemError
	}
	return &user, protocol.Success
}

// GetUserByPhoneAndType 根据手机号和用户类型获取用户
func (s *UserService) GetUserByPhoneAndType(phone, userType string) *models.User {
	return models.GetUserByPhoneAndType(phone, userType)
}

func (s *UserService) UserOnline(req protocol.UserOnlineRequest) protocol.ErrorCode {
	vehicle := models.GetVehicleByID(req.VehicleID) // 预加载车辆缓存
	if vehicle == nil {
		log.Printf("No vehicle found %s", req.VehicleID)
		return protocol.VehicleNotFound
	}
	values := &models.UserValues{}
	values.SetActiveStatus(protocol.StatusOnline)

	err := models.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("user_id = ?", req.UserID).Updates(values).Error; err != nil {
			return err
		}

		vehicleValues := &models.VehicleValues{}
		vehicleValues.SetDriver(req.UserID) // 绑定
		return tx.Model(&models.Vehicle{}).Where("vehicle_id = ?", req.VehicleID).Updates(vehicleValues).Error
	})
	if err != nil {
		log.Printf("Failed to set user online: %v", err)
		return protocol.SystemError
	}
	vehicle = models.GetVehicleByDriverID(req.UserID) // 预加载车辆缓存
	if vehicle == nil {
		log.Printf("No vehicle found for driver %s", req.UserID)
		return protocol.SystemError
	}
	go s.RefreshDriverRuntimeCache(req.UserID)
	return protocol.Success
}

func (s *UserService) UserOffline(userID string) protocol.ErrorCode {
	values := &models.UserValues{}
	values.SetActiveStatus(protocol.StatusOffline)

	err := models.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("user_id = ?", userID).Updates(values).Error; err != nil {
			return err
		}

		vehicleValues := &models.VehicleValues{}
		vehicleValues.SetDriver("") // 解绑司机
		return tx.Model(&models.Vehicle{}).Where("driver_id = ?", userID).Updates(vehicleValues).Error
	})
	if err != nil {
		log.Printf("Failed to set user offline: %v", err)
		return protocol.SystemError
	}

	vehicle := models.GetVehicleByDriverID(userID) // 预加载车辆缓存
	if vehicle != nil {
		log.Printf("vehicle still for driver %s", userID)
		return protocol.SystemError
	}
	go s.RefreshDriverRuntimeCache(userID)
	return protocol.Success
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.User, values *models.UserValues) protocol.ErrorCode {
	if user == nil || values == nil {
		return protocol.InvalidParams
	}

	defer func() {
		// 使用SetValues方法更新非空字段
		user.SetValues(values)
	}()

	if err := s.db.Model(user).UpdateColumns(values).Error; err != nil {
		return protocol.SystemError
	}
	return protocol.Success
}

// IsEmailExists 检查邮箱是否存在
func (s *UserService) IsEmailExists(email string) bool {
	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// IsPhoneExists 检查手机号是否存在
func (s *UserService) IsPhoneExists(phone string) bool {
	var count int64
	s.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// IsEmailExistsWithType 检查邮箱和用户类型组合是否存在
func (s *UserService) IsEmailExistsWithType(email, userType string) bool {
	var count int64
	s.db.Model(&models.User{}).Where("email = ? AND user_type = ?", email, userType).Count(&count)
	return count > 0
}

// IsPhoneExistsWithType 检查手机号和用户类型组合是否存在
func (s *UserService) IsPhoneExistsWithType(phone, userType string) bool {
	var count int64
	s.db.Model(&models.User{}).Where("phone = ? AND user_type = ?", phone, userType).Count(&count)
	return count > 0
}

// GetUserByEmailAndType 根据邮箱和用户类型获取用户
func (s *UserService) GetUserByEmailAndType(email, userType string) *models.User {
	return models.GetUserByEmailAndType(email, userType)
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *models.User, password string) bool {
	if user == nil || user.Password == nil {
		return false
	}
	return utils.VerifyPassword(password, user.Salt, *user.Password)
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(userID, newPassword string) protocol.ErrorCode {
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	// 生成新的哈希密码
	hashedPassword, err := utils.HashPassword(newPassword, user.Salt)
	if err != nil {
		return protocol.SystemError
	}

	// 更新密码
	values := &models.UserValues{}
	values.Password = &hashedPassword

	return s.UpdateUser(user, values)
}

// ActivateUser 激活用户
func (s *UserService) ActivateUser(userID string) protocol.ErrorCode {
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	values := &models.UserValues{}
	values.SetStatus(protocol.StatusActive)

	return s.UpdateUser(user, values)
}

// DeactivateUser 停用用户
func (s *UserService) DeactivateUser(userID string) protocol.ErrorCode {
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	values := &models.UserValues{}
	values.SetStatus(protocol.StatusInactive)

	return s.UpdateUser(user, values)
}

// SearchUsers 搜索用户（支持关键字搜索和用户类型筛选）
func (s *UserService) SearchUsers(req *protocol.SearchRequest) ([]*models.User, int64, protocol.ErrorCode) {
	var users []*models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 用户类型过滤
	if req.UserType != "" {
		query = query.Where("user_type = ?", req.UserType)
	}

	// 关键字搜索（只搜索用户名和手机号）
	if req.Keyword != "" {
		searchTerm := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR phone LIKE ?", searchTerm, searchTerm)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	return users, total, protocol.Success
}

// UpdateUserStatus 统一的用户状态更新方法
func (s *UserService) UpdateUserStatus(userID string, status string, isActive bool) protocol.ErrorCode {
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	// 验证状态值是否有效
	if status != protocol.StatusActive && status != protocol.StatusInactive &&
		status != protocol.StatusSuspended && status != protocol.StatusBanned {
		return protocol.InvalidParams
	}

	values := &models.UserValues{}
	values.SetStatus(status)

	return s.UpdateUser(user, values)
}

// UpdateUserID 更新用户ID（用于删除账户时的特殊处理）
func (s *UserService) UpdateUserID(pk int64, newUserID string) protocol.ErrorCode {
	if err := s.db.Model(&models.User{}).Where("id = ?", pk).UpdateColumn("user_id", newUserID).Error; err != nil {
		return protocol.SystemError
	}
	return protocol.Success
}

// VerifyUser 审核用户认证（替代VerifyDriver，统一处理用户认证）
func (s *UserService) VerifyUser(userID string, isEmailVerified *bool, isPhoneVerified *bool, verifiedBy string) protocol.ErrorCode {
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	values := &models.UserValues{}
	now := utils.TimeNowMilli()

	// 设置邮箱验证状态
	if isEmailVerified != nil {
		values.IsEmailVerified = isEmailVerified
		if *isEmailVerified {
			values.EmailVerifiedAt = &now
		}
	}

	// 设置手机验证状态
	if isPhoneVerified != nil {
		values.IsPhoneVerified = isPhoneVerified
		if *isPhoneVerified {
			values.PhoneVerifiedAt = &now
		}
	}

	return s.UpdateUser(user, values)
}

// =============================================================================
// 司机位置管理功能
// =============================================================================

// UpdateUserLocation 更新用户位置
func (s *UserService) UpdateUserLocation(req *protocol.UpdateLocationRequest) protocol.ErrorCode {
	user := s.GetUserByID(req.UserID)
	if user == nil {
		return protocol.UserNotFound
	}

	// 统一时间源 - 关键点：确保User表和历史表使用完全相同的时间戳
	timestamp := req.UpdatedAt
	if timestamp == 0 {
		timestamp = utils.TimeNowMilli()
	}

	// 1. 更新用户表中的最新位置
	values := &models.UserValues{}
	values.SetLatitude(req.Latitude).
		SetLongitude(req.Longitude).
		SetLocationUpdatedAt(timestamp).
		SetOnlineStatus(protocol.StatusOnline)

	if err := models.GetDB().Model(&models.User{}).Where("user_id = ?", req.UserID).UpdateColumns(values).Error; err != nil {
		log.Printf("Failed to update driver location: %v", err)
		return protocol.SystemError
	}
	// 2. 插入位置历史记录
	locationHistory := models.NewUserLocationHistory(req.UserID, req.Latitude, req.Longitude, protocol.StatusOnline, timestamp)
	if err := models.GetDB().Create(locationHistory).Error; err != nil {
		log.Printf("Failed to insert driver location history: %v", err)
	}
	if user.IsDriver() {
		go s.RefreshDriverLocationRuntimeCache(req.UserID)
	}
	return protocol.Success
}

func (s *UserService) NewDriverRuntime(driverID string) *protocol.DriverRuntime {
	info := &protocol.DriverRuntime{
		CurrentOrder: &protocol.QueuedOrderData{},
		QueuedOrders: []*protocol.QueuedOrderData{},
	}
	//
	return info
}

func (s *UserService) GetDriverRuntime(driverID string) (data *protocol.DriverRuntime) {
	data = &protocol.DriverRuntime{
		DriverID: driverID,
	}
	if err := models.GetObjectCache(data.GetCacheKey(), data); err == nil {
		log.Printf("Failed to get driver runtime from Redis: %v", err)
		return data
	}
	return s.RefreshDriverRuntimeCache(driverID)
}

func (s *UserService) RefreshDriverRuntimeCache(driverID string) (data *protocol.DriverRuntime) {
	//完成数据读取和填充
	user := s.GetUserByID(driverID)
	if user == nil {
		log.Printf("Driver not found in database: %s", driverID)
		return nil
	}

	// 检查用户类型
	if !user.IsDriver() {
		return nil
	}
	data = &protocol.DriverRuntime{
		DriverID:           user.UserID,
		OnlineStatus:       user.GetOnlineStatus(),
		Latitude:           user.GetLatitude(),
		Longitude:          user.GetLongitude(),
		Heading:            0.0,
		Speed:              0.0,
		Accuracy:           10.0,
		LocationUpdatedAt:  user.GetLocationUpdatedAt(),
		QueuedOrders:       []*protocol.QueuedOrderData{},
		MaxQueueCapacity:   user.GetMaxQueueCapacity(),
		ConsecutiveRejects: 0,
		LastDispatchAt:     0,
		LastResponseAt:     0,
		AcceptanceRate:     1.0,
		Rating:             user.GetRating(),
		ExperienceLevel:    1,
		LastHeartbeatAt:    user.GetLocationUpdatedAt(),
		NextAvailableAt:    0,
		UpdatedAt:          utils.TimeNowMilli(),
	}
	defer func() {
		if err := models.SetObjectCache(data.GetCacheKey(), data, 5*time.Minute); err != nil {
			log.Printf("Failed to set driver runtime to Redis: %v", err)
		}
	}()
	vehicle := models.GetVehicleByDriverID(user.UserID)
	if vehicle != nil {
		data.VehicleID = vehicle.VehicleID
	}

	orderIds := []string{}
	if user.GetCurrentOrderID() != "" {
		orderIds = append(orderIds, user.GetCurrentOrderID())
	}
	if len(user.QueuedOrderIds) > 0 {
		orderIds = append(orderIds, user.QueuedOrderIds...)
	}
	if len(orderIds) > 0 {
		order_list := models.GetOrderListByID(user.QueuedOrderIds)
		for _, order := range order_list {
			qorder := &protocol.QueuedOrderData{
				OrderID:     order.OrderID,
				Status:      order.GetStatus(),
				ScheduledAt: order.GetScheduledAt(),
				StartAt:     order.GetStartedAt(),
				EndAt:       order.GetEndedAt(),
			}

			detail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
			if detail != nil {
				qorder.PickupLatitude = detail.GetPickupLatitude()
				qorder.PickupLongitude = detail.GetPickupLongitude()
				qorder.DropoffLatitude = detail.GetDropoffLatitude()
				qorder.DropoffLongitude = detail.GetDropoffLongitude()
				qorder.EstimatedDuration = detail.GetEstimatedDuration()
				qorder.PassengerCount = detail.GetPassengerCount()
			}

			data.QueuedOrders = append(data.QueuedOrders, qorder)
			if data.CurrentOrder == nil && order.OrderID == user.GetCurrentOrderID() {
				data.CurrentOrder = qorder
			}
		}
	}
	return data
}

func (s *UserService) RefreshDriverLocationRuntimeCache(driverID string) (data *protocol.DriverRuntime) {
	data = &protocol.DriverRuntime{
		DriverID: driverID,
	}
	if err := models.GetObjectCache(data.GetCacheKey(), data); err == nil || data == nil {
		return s.RefreshDriverRuntimeCache(driverID)
	}
	defer func() {
		if err := models.SetObjectCache(data.GetCacheKey(), data, 5*time.Minute); err != nil {
			log.Printf("Failed to set driver runtime to Redis: %v", err)
		}
	}()
	user := models.GetUserByID(driverID)
	data.Latitude = user.GetLatitude()
	data.Longitude = user.GetLongitude()
	data.LocationUpdatedAt = user.GetLocationUpdatedAt()
	return data
}

func (s *UserService) GetDriversRuntime(drivers []string) (list []*protocol.DriverRuntime) {
	// 构建Redis键列表
	keys := make([]string, len(drivers))
	for i, driverID := range drivers {
		keys[i] = fmt.Sprintf("driver:realtime:%s", driverID)
	}
	driverLib := map[string]*protocol.DriverRuntime{}
	// 批量获取司机实时数据
	cacheList := models.BatchGetObjectCache(keys, &protocol.DriverRuntime{})
	if len(cacheList) > 0 {
		// 转换结果格式并设置DriverID
		for _, item := range cacheList {
			if realtimeData, ok := item.(*protocol.DriverRuntime); ok {
				list = append(list, realtimeData)
				driverLib[realtimeData.DriverID] = realtimeData
			}
		}
	}
	lessList := []string{}
	for _, driverID := range drivers {
		if _, ok := driverLib[driverID]; !ok {
			lessList = append(lessList, driverID)
		}
	}
	if len(lessList) > 0 {
		for _, driverID := range lessList {
			data := s.GetDriverRuntime(driverID)
			if data != nil {
				list = append(list, data)
			}
		}
	}

	return list
}

func (s *UserService) RefreshUserOrderQueue(userID string) {
	var order_list []*models.Order
	if err := models.GetDB().Select([]string{"order_id", "status"}).
		Where("user_id = ? ", userID).
		Where("status IN ?", []string{protocol.StatusRequested, protocol.StatusDriverArrived, protocol.StatusDriverComing, protocol.StatusAccepted, protocol.StatusInProgress}).
		Order("scheduled_at ASC").Find(&order_list).Error; err != nil {
		log.Printf("Failed to refresh user %s queue: %v", userID, err)
		return
	}
	values := &models.UserValues{}
	values.SetCurrentOrderId("").
		SetQueuedOrderIds([]string{}).
		SetQueueUpdatedAt(utils.TimeNowMilli())
	for _, order := range order_list {
		if slices.Contains(protocol.CURRENT_RIDE_STATUS, order.GetStatus()) {
			values.SetCurrentOrderId(order.OrderID)
			continue
		}
		values.QueuedOrderIds = append(values.QueuedOrderIds, order.OrderID)
	}

	if err := models.GetDB().Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(values).Error; err != nil {
		log.Printf("Failed to update user %s queue: %v", userID, err)
	}
}

func (s *UserService) RefreshDriverOrderQueue(driverID string) {
	var order_list []*models.Order
	if err := models.GetDB().Select([]string{"order_id", "status"}).
		Where("provider_id = ? ", driverID).
		Where("status IN ?", []string{protocol.StatusRequested, protocol.StatusDriverArrived, protocol.StatusDriverComing, protocol.StatusAccepted, protocol.StatusInProgress}).
		Order("scheduled_at DESC").Find(&order_list).Error; err != nil {
		log.Printf("Failed to refresh driver %s queue: %v", driverID, err)
		return
	}
	values := &models.UserValues{}
	values.SetCurrentOrderId("").
		SetQueuedOrderIds([]string{}).
		SetQueueUpdatedAt(utils.TimeNowMilli())
	for _, order := range order_list {
		if slices.Contains(protocol.CURRENT_RIDE_STATUS, order.GetStatus()) {
			values.SetCurrentOrderId(order.OrderID)
			continue
		}
		values.QueuedOrderIds = append(values.QueuedOrderIds, order.OrderID)
	}
	if err := models.GetDB().Model(&models.User{}).
		Where("user_id = ?", driverID).
		Updates(values).Error; err != nil {
		log.Printf("Failed to update driver %s queue: %v", driverID, err)
	}
	s.RefreshDriverRuntimeCache(driverID)
}

// UploadAvatarRequest 头像上传请求参数
type UploadAvatarRequest struct {
	UserID    string    `json:"user_id"`   // 用户ID
	Reader    io.Reader `json:"-"`         // 文件读取器 (io.Reader)
	Extension string    `json:"extension"` // 文件扩展名
	Filename  string    `json:"filename"`  // 原始文件名
	Size      int64     `json:"size"`      // 文件大小
}

// UploadAvatarResponse 头像上传响应
type UploadAvatarResponse struct {
	AvatarURL string             `json:"avatar_url"` // 新的头像URL
	ErrorCode protocol.ErrorCode `json:"error_code"` // 错误码
}

// UploadAvatar 上传用户头像（通用方法，可在API和Admin中复用）
func (s *UserService) UploadAvatar(ctx context.Context, req *UploadAvatarRequest) *UploadAvatarResponse {
	// 检查AWS S3服务是否可用
	awsService, available := GetAWSServiceSafe()
	if !available {
		return &UploadAvatarResponse{
			ErrorCode: protocol.ServiceUnavail,
		}
	}

	// 上传头像到S3
	avatarURL, err := awsService.UploadUserAvatar(ctx, req.UserID, req.Reader, req.Extension)
	if err != nil {
		return &UploadAvatarResponse{
			ErrorCode: protocol.ThirdPartyError,
		}
	}

	// 更新用户头像URL到数据库
	values := &models.UserValues{}
	values.Avatar = &avatarURL

	user := s.GetUserByID(req.UserID)
	if user == nil {
		return &UploadAvatarResponse{
			ErrorCode: protocol.UserNotFound,
		}
	}

	if errCode := s.UpdateUser(user, values); errCode != protocol.Success {
		return &UploadAvatarResponse{
			ErrorCode: errCode,
		}
	}

	return &UploadAvatarResponse{
		AvatarURL: avatarURL,
		ErrorCode: protocol.Success,
	}
}

// DeleteUserByID 删除用户（软删除）
// 参数：
//   - userID: 要删除的用户ID
//   - reason: 删除原因（可选）
//   - isPassengerOnly: 是否只允许删除乘客类型用户
//
// 返回：
//   - protocol.ErrorCode: 错误码
func (s *UserService) DeleteUserByID(userID string, reason string, isPassengerOnly bool) protocol.ErrorCode {
	// 获取用户信息
	user := s.GetUserByID(userID)
	if user == nil {
		return protocol.UserNotFound
	}

	// 检查用户类型权限
	if isPassengerOnly && !user.IsPassenger() {
		return protocol.PermissionDenied
	}

	// 检查是否有未完成的订单
	if count := models.CountUnCompletedOrdersByUserID(user.UserID); count > 0 {
		return protocol.UserHasUnCompletedOrders
	}

	// 生成删除时间戳（毫秒）
	timestamp := utils.TimeNowMilli()

	// 更新用户信息，添加删除标记
	values := &models.UserValues{}
	values.SetStatus(protocol.StatusDeleted).
		SetDeletedAt(timestamp)

	// 添加删除前缀和时间戳后缀
	if user.GetEmail() != "" && !strings.HasPrefix(user.GetEmail(), "deleted_") {
		deletedEmail := "deleted_" + user.GetEmail()
		values.SetEmail(deletedEmail)
	}
	if user.GetPhone() != "" && !strings.HasPrefix(user.GetPhone(), "deleted_") {
		deletedPhone := "deleted_" + user.GetPhone()
		values.SetPhone(deletedPhone)
	}

	// 先更新 UserValues 字段
	if errCode := s.UpdateUser(user, values); errCode != protocol.Success {
		return errCode
	}

	// 禁用用户促销信息
	if err := models.DisabledUserPromotionByUserID(s.db, user.UserID); err != nil {
		log.Printf("Failed to disable user promotions - user_id: %s, error: %s", user.UserID, err.Error())
		// 这里不返回错误，因为主要的删除操作已经完成
	}

	log.Printf("User deleted successfully - user_id: %s, user_type: %s, reason: %s", user.UserID, user.GetUserType(), reason)
	return protocol.Success
}

// CreateUserByAdmin 管理员创建用户（密码为空，直接发放优惠券）
func (s *UserService) CreateUserByAdmin(req *protocol.AdminCreateUserRequest, adminID string) (*models.User, protocol.ErrorCode) {
	// 验证基本参数
	if req.Phone == "" {
		return nil, protocol.MissingParams
	}

	// 设置默认用户类型
	if req.UserType == "" {
		req.UserType = protocol.UserTypePassenger
	}

	// 检查手机号是否已存在（同类型用户）
	if s.IsPhoneExistsWithType(req.Phone, req.UserType) {
		return nil, protocol.PhoneAlreadyExists
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if req.Email != "" && s.IsEmailExistsWithType(req.Email, req.UserType) {
		return nil, protocol.EmailAlreadyExists
	}

	// 创建用户
	user := models.NewUser()

	// 设置用户信息（密码为空）
	user.UserValues.
		SetPhone(req.Phone).
		SetUsername(req.Username).
		SetUserType(req.UserType).
		SetStatus(protocol.StatusActive). // 直接设置为激活状态
		SetIsPhoneVerified(true)          // 管理员创建默认手机已验证

	// 设置可选字段
	if req.Email != "" {
		user.UserValues.SetEmail(req.Email).SetIsEmailVerified(true)
	}
	if req.FirstName != "" {
		user.UserValues.SetFirstName(req.FirstName)
	}
	if req.LastName != "" {
		user.UserValues.SetLastName(req.LastName)
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		log.Printf("CreateUserByAdmin failed to create user - phone: %s, error: %s", req.Phone, err.Error())
		return nil, protocol.DatabaseError
	}

	log.Printf("CreateUserByAdmin created new user - user_id: %s, phone: %s, user_type: %s, admin_id: %s",
		user.UserID, req.Phone, req.UserType, adminID)

	// 直接同步发放优惠券（如果是乘客类型）
	if user.IsPassenger() {
		if err := s.issueWelcomeCouponSync(user.UserID, adminID); err != nil {
			// 记录错误但不影响用户创建
			log.Printf("Failed to issue welcome coupon for admin-created user: user_id=%s, admin_id=%s, error=%v",
				user.UserID, adminID, err)
		}
	}

	return user, protocol.Success
}

// issueWelcomeCouponSync 同步发放欢迎优惠券（管理员创建用户时使用）
func (s *UserService) issueWelcomeCouponSync(userID string, adminID string) error {
	// 检查用户是否已有欢迎优惠券
	if models.CheckUserHasWelcomeCoupon(userID) {
		log.Printf("User %s already has welcome coupon", userID)
		return nil
	}

	// 获取或创建默认欢迎优惠券模板
	promotion := models.GetOrCreateDefaultWelcomePromotion()
	if promotion == nil {
		return fmt.Errorf("failed to get or create welcome promotion template")
	}

	// 创建用户优惠券实例
	userPromotion := models.CreateWelcomePromoForUser(userID, promotion)

	// 设置来源为管理员创建
	userPromotion.SetSource("admin_created")
	userPromotion.SetBatchID(fmt.Sprintf("admin_%s_%d", adminID, utils.TimeNowMilli()))

	// 保存到数据库（同步处理）
	if err := models.CreateUserPromotionInDB(userPromotion); err != nil {
		return fmt.Errorf("failed to create user promotion: %v", err)
	}

	log.Printf("Welcome coupon issued successfully by admin %s for user %s: %s",
		adminID, userID, userPromotion.GetCode())
	return nil
}

// =============================================================================
// Nearby Drivers Service (for passengers to find drivers)
// =============================================================================

// GetNearbyDrivers 获取附近在线司机列表
func (s *UserService) GetNearbyDrivers(latitude, longitude, radiusKm float64, limit int) ([]*protocol.NearbyDriver, error) {
	// Set defaults
	if radiusKm <= 0 {
		radiusKm = 5.0 // Default 5km radius
	}
	if limit <= 0 || limit > 50 {
		limit = 20 // Default 20 drivers, max 50
	}

	// Query online drivers within radius using Haversine formula
	// Distance in km = 6371 * acos(cos(radians(lat1)) * cos(radians(lat2)) * cos(radians(lng2) - radians(lng1)) + sin(radians(lat1)) * sin(radians(lat2)))
	// Uses t_user_location_history for recent location updates (within last 5 minutes = 300000ms)
	query := `
		SELECT 
			u.user_id,
			u.first_name,
			u.last_name,
			u.display_name,
			u.avatar,
			u.phone,
			COALESCE(ulh.latitude, u.latitude) as latitude,
			COALESCE(ulh.longitude, u.longitude) as longitude,
			COALESCE(ulh.heading, 0) as heading,
			COALESCE(ulh.online_status, u.online_status) as online_status,
			u.score,
			u.total_rides,
			v.brand,
			v.model,
			v.plate_number,
			v.color,
			v.category,
			(6371 * acos(
				LEAST(1.0, GREATEST(-1.0,
					cos(radians(?)) * cos(radians(COALESCE(ulh.latitude, u.latitude))) * cos(radians(COALESCE(ulh.longitude, u.longitude)) - radians(?)) +
					sin(radians(?)) * sin(radians(COALESCE(ulh.latitude, u.latitude)))
				))
			)) AS distance_km
		FROM t_users u
		LEFT JOIN (
			SELECT user_id, latitude, longitude, heading, online_status, recorded_at,
				ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY recorded_at DESC) as rn
			FROM t_user_location_history
			WHERE recorded_at > (UNIX_TIMESTAMP(NOW()) * 1000 - 300000)
		) ulh ON u.user_id = ulh.user_id AND ulh.rn = 1
		LEFT JOIN t_vehicles v ON u.user_id = v.driver_id AND v.status = 'active'
		WHERE u.user_type = 'driver'
			AND (u.online_status = 'online' OR u.online_status = 'busy' OR ulh.online_status = 'online' OR ulh.online_status = 'busy')
			AND u.status = 'active'
			AND (u.latitude IS NOT NULL OR ulh.latitude IS NOT NULL)
			AND (u.longitude IS NOT NULL OR ulh.longitude IS NOT NULL)
			AND u.deleted_at IS NULL
		HAVING distance_km <= ?
		ORDER BY distance_km ASC
		LIMIT ?
	`

	type driverRow struct {
		UserID       string   `gorm:"column:user_id"`
		FirstName    *string  `gorm:"column:first_name"`
		LastName     *string  `gorm:"column:last_name"`
		DisplayName  *string  `gorm:"column:display_name"`
		Avatar       *string  `gorm:"column:avatar"`
		Phone        *string  `gorm:"column:phone"`
		Latitude     *float64 `gorm:"column:latitude"`
		Longitude    *float64 `gorm:"column:longitude"`
		Heading      float64  `gorm:"column:heading"`
		OnlineStatus *string  `gorm:"column:online_status"`
		Score        *float64 `gorm:"column:score"`
		TotalRides   *int     `gorm:"column:total_rides"`
		Brand        *string  `gorm:"column:brand"`
		Model        *string  `gorm:"column:model"`
		PlateNumber  *string  `gorm:"column:plate_number"`
		Color        *string  `gorm:"column:color"`
		Category     *string  `gorm:"column:category"`
		DistanceKm   float64  `gorm:"column:distance_km"`
	}

	var rows []driverRow
	if err := s.db.Raw(query, latitude, longitude, latitude, radiusKm, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	drivers := make([]*protocol.NearbyDriver, 0, len(rows))
	for _, row := range rows {
		// Build driver name
		name := ""
		if row.DisplayName != nil && *row.DisplayName != "" {
			name = *row.DisplayName
		} else if row.FirstName != nil || row.LastName != nil {
			firstName := ""
			lastName := ""
			if row.FirstName != nil {
				firstName = *row.FirstName
			}
			if row.LastName != nil {
				lastName = *row.LastName
			}
			name = strings.TrimSpace(firstName + " " + lastName)
		}
		if name == "" {
			name = "Driver"
		}

		// Calculate ETA (rough estimate: assume 30 km/h average speed in city)
		etaMinutes := int(row.DistanceKm * 2) // 2 minutes per km
		if etaMinutes < 1 {
			etaMinutes = 1
		}

		// Determine online/busy status
		isOnline := true
		isBusy := false
		if row.OnlineStatus != nil {
			isOnline = *row.OnlineStatus == "online" || *row.OnlineStatus == "busy"
			isBusy = *row.OnlineStatus == "busy"
		}

		driver := &protocol.NearbyDriver{
			DriverID:   row.UserID,
			Name:       name,
			DistanceKm: row.DistanceKm,
			ETAMinutes: etaMinutes,
			IsOnline:   isOnline,
			IsBusy:     isBusy,
			Heading:    row.Heading,
		}

		// Set optional fields
		if row.Avatar != nil {
			driver.PhotoURL = *row.Avatar
		} else {
			driver.PhotoURL = protocol.GenerateDefaultAvatar(row.UserID)
		}
		if row.Phone != nil {
			driver.Phone = *row.Phone
		}
		if row.Latitude != nil {
			driver.Latitude = *row.Latitude
		}
		if row.Longitude != nil {
			driver.Longitude = *row.Longitude
		}
		if row.Score != nil {
			driver.Rating = *row.Score
		} else {
			driver.Rating = 5.0 // Default rating
		}
		if row.TotalRides != nil {
			driver.TotalRides = *row.TotalRides
		}
		if row.Brand != nil {
			driver.VehicleBrand = *row.Brand
		}
		if row.Model != nil {
			driver.VehicleModel = *row.Model
		}
		if row.PlateNumber != nil {
			driver.PlateNumber = *row.PlateNumber
		}
		if row.Color != nil {
			driver.VehicleColor = *row.Color
		}
		if row.Category != nil {
			driver.VehicleType = *row.Category
			driver.VehicleCategory = *row.Category
		} else {
			driver.VehicleType = "sedan"
			driver.VehicleCategory = "sedan"
		}

		drivers = append(drivers, driver)
	}

	return drivers, nil
}
