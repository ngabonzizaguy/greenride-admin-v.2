package services

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

var (
	dispatchServiceInstance *DispatchService
	dispatchServiceOnce     sync.Once
)

// DispatchService 简化的派单服务
type DispatchService struct {
	config *config.DispatchConfig // 内置派单配置
}

func GetDispatchService() *DispatchService {
	if dispatchServiceInstance == nil {
		SetupDispatchService()
	}
	return dispatchServiceInstance
}

func SetupDispatchService() {
	dispatchServiceOnce.Do(func() {
		cfg := config.Get().Dispatch
		if cfg == nil {
			return
		}
		dispatchServiceInstance = &DispatchService{
			config: cfg,
		}
	})
}

// ================ 核心派单函数 ================
// StartAutoDispatch 启动自动派单
func (s *DispatchService) StartAutoDispatch(order *protocol.Order) (result *protocol.DispatchResult) {
	log.Get().Infof("Starting auto dispatch for order %s, round %d", order.OrderID, order.CurrentRound)
	result = &protocol.DispatchResult{
		Success:     false,
		DriverCount: 0,
	}
	vehicleCategory := ""
	if order.Details != nil {
		vehicleCategory = order.Details.VehicleCategory
	}
	log.Get().Infof("[Dispatch] Order %s: max_distance=%.1f km, vehicle_category=%s",
		order.OrderID, s.config.MaxDistance, vehicleCategory)
	driver_list := s.FindEligibleDrivers(order)
	if len(driver_list) == 0 {
		result.Message = "No eligible drivers found"
		log.Get().Warnf("[Dispatch] Order %s: No drivers found with bound vehicles (category=%s)",
			order.OrderID, vehicleCategory)
		return
	}
	log.Get().Infof("[Dispatch] Order %s: Found %d drivers with matching vehicles", order.OrderID, len(driver_list))
	runtime_list := GetUserService().GetDriversRuntime(driver_list)
	if len(runtime_list) == 0 {
		result.Message = "No online drivers found"
		log.Get().Warnf("[Dispatch] Order %s: None of %d drivers have runtime data", order.OrderID, len(driver_list))
		return
	}
	log.Get().Infof("[Dispatch] Order %s: %d drivers have runtime data", order.OrderID, len(runtime_list))
	// 4. 评估每个司机
	var eligible_drivers []*protocol.DispatchDriver
	for _, item := range runtime_list {
		// 评估司机是否适合接单
		eligibleDriver := s.EvaluateDriverForOrder(item, order)
		if !eligibleDriver.IsEligible {
			log.Get().Infof("[Dispatch] Order %s: Driver %s rejected — %s (dist=%.1f km, status=%s)",
				order.OrderID, item.DriverID, eligibleDriver.RejectReason, eligibleDriver.Distance, item.OnlineStatus)
			continue
		}
		eligible_drivers = append(eligible_drivers, eligibleDriver)
	}
	log.Get().Infof("[Dispatch] Order %s: %d/%d drivers eligible after evaluation",
		order.OrderID, len(eligible_drivers), len(runtime_list))

	// 2. 司机按评分排序
	sort.Slice(eligible_drivers, func(i, j int) bool {
		return eligible_drivers[i].FinalScore > eligible_drivers[j].FinalScore
	})

	// 3. 执行派单
	records := s.ExecuteDispatch(eligible_drivers, order)
	if len(records) == 0 {
		result.Message = "Failed to execute dispatch"
		return
	}

	// 返回派单结果
	result.Success = true
	result.DriverCount = len(eligible_drivers)
	result.Drivers = eligible_drivers

	return
}

// FindEligibleDrivers 查找符合条件的司机（广播模式：只要求有车辆，不按车型筛选；在线+无当前订单在 EvaluateDriverForOrder 中过滤）
func (s *DispatchService) FindEligibleDrivers(order *protocol.Order) []string {
	// Broadcast: include all drivers who have a vehicle. Eligibility (online, no active ride) is enforced in EvaluateDriverForOrder.
	// Vehicle category/level are not required to match — order goes to all such drivers; first to accept wins.
	driverList := models.FindDriversByVehicle("", "")
	return driverList
}

// EvaluateDriverForOrder 评估单个司机是否适合接单（仅强制：在线、无当前订单；其余为可选）
func (s *DispatchService) EvaluateDriverForOrder(rt *protocol.DriverRuntime, order *protocol.Order) (driver *protocol.DispatchDriver) {
	driver = &protocol.DispatchDriver{
		DriverID:   rt.DriverID,
		IsEligible: false,
	}
	// 1. Mandatory: must be online
	if !rt.IsAvailable() {
		driver.RejectReason = "Driver not available"
		return
	}
	// 2. Mandatory: must have no active ride (one ride at a time)
	if rt.HasCurrentOrder() {
		driver.RejectReason = "Driver has active ride"
		return
	}
	// 3. Optional: queue capacity (if configured)
	if !rt.CanAcceptMoreOrders() {
		driver.RejectReason = "Driver queue is full"
		return
	}

	driver.Distance = 0
	if order.Details != nil {
		driver.Distance = utils.CalculateDistanceHaversine(
			rt.Latitude,
			rt.Longitude,
			order.Details.PickupLatitude,
			order.Details.PickupLongitude,
		)
	}
	// 4. Optional: distance filter only when MaxDistance > 0
	if s.config.MaxDistance > 0 && driver.Distance > s.config.MaxDistance {
		driver.RejectReason = "Distance too far"
		return
	}

	timeWindow := s.analyzeDriverTimeWindow(rt, order)
	if !timeWindow.CanAcceptNewOrder {
		driver.RejectReason = "Cannot accept new order due to time window"
		driver.WaitTimeMinutes = timeWindow.WaitTimeMinutes
		return
	}

	driver.IsEligible = true
	driver.CanAcceptNewOrder = timeWindow.CanAcceptNewOrder
	driver.WaitTimeMinutes = timeWindow.WaitTimeMinutes

	return
}

// ExecuteDispatch 执行派单
func (s *DispatchService) ExecuteDispatch(drivers []*protocol.DispatchDriver, order *protocol.Order) (list []string) {
	dispatchedAt := utils.TimeNowMilli()
	for idx, driver := range drivers {
		// 创建派单记录
		record := &models.DispatchRecord{
			DriverID:             driver.DriverID,
			OrderID:              order.OrderID,
			DispatchID:           utils.GenerateDispatchID(),
			Round:                order.CurrentRound,
			DispatchedAt:         dispatchedAt,
			ExpiredAt:            order.ScheduledAt, // 过期时间
			RoundSeq:             idx + 1,           // 本轮派单顺序
			DispatchRecordValues: &models.DispatchRecordValues{},
			CreatedAt:            dispatchedAt,
		}
		if err := models.GetDB().Create(record).Error; err != nil {
			log.Get().Warnf("Warning: failed to create dispatch record for driver %s: %v", driver.DriverID, err)
			continue
		}
		// 异步发送推送通知
		go s.SendDispatchNotifications(record)
		// 记录派单ID
		list = append(list, record.DispatchID)
	}

	return list
}

func (s *DispatchService) SendDispatchNotifications(record *models.DispatchRecord) {
	if record == nil {
		return
	}
	// 获取订单信息
	order := models.GetOrderByID(record.OrderID)
	if order == nil {
		return
	}
	if order.GetProviderID() != "" && order.GetProviderID() != record.DriverID {
		return
	}

	// 获取司机信息
	driver := models.GetUserByID(record.DriverID)
	if driver == nil {
		return
	}

	// 获取订单详情
	orderDetail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
	if orderDetail == nil {
		return
	}

	// 准备消息参数
	params := map[string]any{
		"to":                record.DriverID,
		"dispatch_id":       record.DispatchID,
		"OrderID":           order.OrderID,
		"order_id":          order.OrderID,
		"OrderStatus":       order.GetStatus(),
		"order_status":      order.GetStatus(),
		"PickupAddress":     orderDetail.GetPickupAddress(),
		"DropoffAddress":    orderDetail.GetDropoffAddress(),
		"pickup_latitude":   orderDetail.GetPickupLatitude(),
		"pickup_longitude":  orderDetail.GetPickupLongitude(),
		"dropoff_latitude":  orderDetail.GetDropoffLatitude(),
		"dropoff_longitude": orderDetail.GetDropoffLongitude(),
		"Amount":            order.GetPaymentAmount().StringFixed(2),
		"Currency":          order.GetCurrency(),
		"Distance":          orderDetail.GetEstimatedDistance(),
		"Duration":          orderDetail.GetEstimatedDuration(),
		"msg_type":          protocol.FCMMessageTypeOrder,
		"notification_type": protocol.NotificationTypeNewOrderAvailable,
		"timeout_seconds":   s.config.TimeoutSeconds,
	}

	// 添加乘客信息
	if order.GetUserID() != "" {
		passenger := models.GetUserByID(order.GetUserID())
		if passenger != nil {
			params["PassengerName"] = passenger.GetFullName()
		}
	}

	// Add driver-to-pickup ETA and distance (Google API with rough fallback)
	if driver.GetLatitude() != 0 && orderDetail.GetPickupLatitude() != 0 {
		distKm := utils.CalculateDistanceHaversine(
			driver.GetLatitude(), driver.GetLongitude(),
			orderDetail.GetPickupLatitude(), orderDetail.GetPickupLongitude(),
		)

		// Try Google Directions API for accurate ETA
		etaMin := 0
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		route, err := GetGoogleService().CalculateRidehailingRoute(
			ctx, driver.GetLatitude(), driver.GetLongitude(),
			orderDetail.GetPickupLatitude(), orderDetail.GetPickupLongitude(), false,
		)
		cancel()
		if err == nil && route != nil && route.Duration != nil && route.Duration.Value > 0 {
			etaMin = int((route.Duration.Value + 59) / 60) // ceil seconds→minutes
			if route.Distance != nil && route.Distance.Value > 0 {
				distKm = float64(route.Distance.Value) / 1000.0 // use routing distance
			}
		} else {
			// Fallback to rough estimate
			if err != nil {
				log.Get().Warnf("[Dispatch] Google ETA failed for driver %s → pickup: %v, using rough estimate", record.DriverID, err)
			}
			etaMin = int(distKm * 2) // rough: 2 min/km
		}
		if etaMin < 2 {
			etaMin = 2 // minimum 2 minutes (more realistic for urban driving)
		}
		params["DriverToPickupETA"] = etaMin
		params["DriverToPickupDistance"] = fmt.Sprintf("%.1f", distKm)
	}

	// 创建消息对象
	message := &Message{
		Type:     protocol.MsgTypeDriverNewOrder,
		Channels: []string{protocol.MsgChannelFcm},
		Params:   params,
		Language: getUserLanguage(driver),
	}

	// 使用消息服务发送
	err := GetMessageService().SendMessage(message)
	if err != nil {
		log.Get().Errorf("Failed to send dispatch notification to driver %s: %v", record.DriverID, err)
	} else {
		log.Get().Infof("Dispatch notification sent to driver %s for order %s", record.DriverID, order.OrderID)
	}
}

// ================ 司机响应处理 ================

// HandleDriverAccept 处理司机接单
func (s *DispatchService) HandleDriverAccept(dispatchID, driverID string, latitude, longitude float64) protocol.ErrorCode {
	// 获取派单记录
	record := models.GetDispatchByID(dispatchID)
	if record == nil || record.DriverID != driverID {
		return protocol.OrderNotFound
	}

	if record.GetStatus() != protocol.StatusPending {
		return protocol.InvalidOrderStatus
	}

	// 更新派单记录为已接受
	values := &models.DispatchRecordValues{}
	values.SetStatus(protocol.StatusAccepted).
		SetRespondedAt(utils.TimeNowMilli())

	if err := models.GetDB().Model(record).UpdateColumns(values).Error; err != nil {
		return protocol.SystemError
	}
	otherValues := &models.DispatchRecordValues{}
	otherValues.SetStatus(protocol.StatusCancelled)
	if err := models.GetDB().Model(&models.DispatchRecord{}).
		Where("order_id = ?", record.OrderID).
		Where("dispatch_id != ?", record.DispatchID).
		Where("status = ?", protocol.StatusPending).
		UpdateColumns(otherValues).Error; err != nil {
		log.Get().Warnf("Warning: failed to cancel other dispatches for order %s: %v", record.OrderID, err)
	}
	log.Get().Infof("Driver %s accepted dispatch %s for order %s", record.DriverID, dispatchID, record.OrderID)
	return protocol.Success
}

// HandleDriverReject 处理司机拒单
func (s *DispatchService) HandleDriverReject(dispatchID, driverID, reason string, latitude, longitude float64) protocol.ErrorCode {
	// 获取派单记录
	record := models.GetDispatchByID(dispatchID)
	if record == nil || record.DriverID != driverID {
		return protocol.OrderNotFound
	}
	if record.GetStatus() != protocol.StatusPending {
		return protocol.InvalidOrderStatus
	}

	// 更新派单记录为已接受
	values := &models.DispatchRecordValues{}
	values.SetStatus("rejected").
		SetRespondedAt(utils.TimeNowMilli()).
		SetDriverLatitude(latitude).
		SetDriverLongitude(longitude).
		SetRejectReason(reason)

	// 判断是否为枚举值
	if protocol.IsValidRejectReason(reason) {
		values.SetRejectReasonType(reason)
	}
	if err := models.GetDB().Model(record).UpdateColumns(values).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// HandleDriverTimeout 处理司机超时
func (s *DispatchService) HandleDriverTimeout(dispatchID, driverID string) protocol.ErrorCode {
	// 获取派单记录
	record := models.GetDispatchByID(dispatchID)
	if record == nil {
		return protocol.OrderNotFound
	}

	if record.DriverID != driverID {
		return protocol.MissingParams
	}

	if record.GetStatus() != protocol.StatusPending {
		return protocol.InvalidOrderStatus
	}

	// 更新派单记录为已接受
	values := &models.DispatchRecordValues{}
	values.SetStatus("timeout").
		SetRespondedAt(utils.TimeNowMilli())

	if err := models.GetDB().Model(record).UpdateColumns(values).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// ================ 司机数据获取函数 ================

// getNearbyDriverIDs 获取附近司机ID列表
func (s *DispatchService) getNearbyDriverIDs(latitude, longitude, radius float64) ([]string, error) {
	// 使用Redis GEO命令查找附近司机
	geoKey := "drivers:geo"

	result, err := models.GetRedis().GeoRadius(context.Background(), geoKey, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        "km",
		Sort:        "ASC", // 按距离升序
		WithCoord:   false,
		WithDist:    false,
		WithGeoHash: false,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get nearby drivers from geo: %w", err)
	}

	var driverIDs []string
	for _, location := range result {
		driverIDs = append(driverIDs, location.Name)
	}

	return driverIDs, nil
}

// ================ 司机筛选和评估函数 ================

// analyzeDriverTimeWindow 分析司机时间窗口
func (s *DispatchService) analyzeDriverTimeWindow(rt *protocol.DriverRuntime, order *protocol.Order) *config.DriverTimeWindow {
	timeWindow := &config.DriverTimeWindow{
		CanAcceptNewOrder: true,
		WaitTimeMinutes:   10,
		RouteMatchScore:   1.0,
	}

	return timeWindow
}

func (s *DispatchService) GetDispatchRecordsByUser(req *protocol.UserDispatchsRequest) (records []*models.DispatchRecord, total int64) {
	// 获取司机的派单记录
	if err := models.GetDB().Model(&models.DispatchRecord{}).
		Where("driver_id = ?", req.UserID).
		Count(&total).Error; err != nil {
		return
	}

	if err := models.GetDB().Model(&models.DispatchRecord{}).
		Where("driver_id = ?", req.UserID).
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&records).Error; err != nil {
		return
	}

	return
}
