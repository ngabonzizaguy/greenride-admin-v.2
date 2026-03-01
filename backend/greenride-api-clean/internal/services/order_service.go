package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/i18n"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"gorm.io/gorm"
)

// OrderService 通用订单服务 - 整合了所有订单相关功能
type OrderService struct {
}

var (
	orderServiceInstance *OrderService
	orderServiceOnce     sync.Once
)

// GetOrderService 获取订单服务单例
func GetOrderService() *OrderService {
	if orderServiceInstance == nil {
		SetupOrderService()
	}
	return orderServiceInstance
}

// SetupOrderService 设置订单服务
func SetupOrderService() {
	orderServiceOnce.Do(func() {
		orderServiceInstance = &OrderService{}
	})
}

// ============================================================================
// 基础订单管理功能
// ============================================================================

// GetOrdersByUser 根据用户获取订单列表
func (s *OrderService) GetOrdersByUser(req *protocol.UserRidesRequest) ([]*protocol.Order, int64) {
	query := models.DB.Model(&models.Order{})

	// 如果没有指定用户类型，则自动识别
	userType := req.UserType
	if userType == "" {
		// 通过用户服务获取用户信息来确定用户类型
		userService := GetUserService()
		user := userService.GetUserByID(req.UserID)
		if user == nil {
			return nil, 0
		}
		userType = user.GetUserType()
	}

	// 根据用户类型设置查询条件
	orderBy := "created_at DESC"
	switch userType {
	case protocol.UserTypePassenger:
		query = query.Where("user_id = ?", req.UserID)
	case protocol.UserTypeDriver:
		query = query.Where("provider_id = ?", req.UserID)
		orderBy = "accepted_at DESC"
	default:
		return nil, 0
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 订单类型过滤
	if req.OrderType != "" {
		query = query.Where("order_type = ?", req.OrderType)
	}

	// 日期过滤 (使用时间戳毫秒)
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", *req.EndDate)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0
	}

	// 获取订单列表
	var orders []*models.Order
	var orderList []string
	offset := (req.Page - 1) * req.Limit
	if err := query.Select([]string{"order_id"}).Offset(offset).Limit(req.Limit).Order(orderBy).Find(&orderList).Error; err != nil {
		return nil, 0
	}
	// 根据订单ID列表获取完整的订单信息
	if len(orderList) > 0 {
		if err := models.DB.Where("order_id IN ?", orderList).Order(orderBy).Find(&orders).Error; err != nil {
			return nil, 0
		}
	}
	list := []*protocol.Order{}
	for _, order := range orders {
		info := s.GetOrderInfoSanitized(order, req.UserID, userType)
		list = append(list, info)
	}

	return list, total
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderID, cancelledBy, reason string) protocol.ErrorCode {
	order := models.GetOrderByID(orderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return protocol.OrderNotFound
	}
	// 检查订单是否可以取消
	if !s.CanCancelOrder(order.GetStatus()) {
		return protocol.CancellationNotAllowed
	}

	// 使用事务处理取消逻辑
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		values := &models.OrderValues{}
		values.CancelOrder(cancelledBy, reason).
			SetCompletedAt(utils.TimeNowMilli())

		if err := models.UpdateOrder(tx, order, values); err != nil {
			return err
		}

		// Cancel all pending dispatch records for this order
		if err := tx.Model(&models.DispatchRecord{}).
			Where("order_id = ? AND status = ?", orderID, protocol.StatusPending).
			UpdateColumn("status", protocol.StatusCancelled).Error; err != nil {
			log.Get().Warnf("Failed to cancel pending dispatches for order %s: %v", orderID, err)
		}

		// 恢复用户优惠券状态
		if err := models.ResetUserPromotionsByIDs(tx, order.GetUserPromotionIDs()); err != nil {
			log.Get().Errorf("取消订单时恢复用户优惠券失败: %v", err)
			//return err
		}

		return nil
	})

	if err != nil {
		return protocol.DatabaseError
	}

	// 发送FCM通知
	go s.NotifyOrderCancelled(orderID)

	return protocol.Success
}

// ReleaseOrder 释放订单
func (s *OrderService) ReleaseOrder(orderID, releasedBy, reason string) protocol.ErrorCode {
	order := models.GetOrderByID(orderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return protocol.OrderNotFound
	}
	// 检查订单是否可以取消
	if !s.CanCancelOrder(order.GetStatus()) {
		return protocol.CancellationNotAllowed
	}

	// 使用事务处理取消逻辑
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		values := &models.OrderValues{}
		values.CancelOrder(releasedBy, reason).
			SetCompletedAt(utils.TimeNowMilli())

		if err := models.UpdateOrder(tx, order, values); err != nil {
			return err
		}

		// 恢复用户优惠券状态
		if err := models.ResetUserPromotionsByIDs(tx, order.GetUserPromotionIDs()); err != nil {
			log.Get().Errorf("取消订单时恢复用户优惠券失败: %v", err)
			//return err
		}

		return nil
	})

	if err != nil {
		return protocol.DatabaseError
	}

	// 发送FCM通知
	go s.NotifyOrderCancelled(orderID)

	return protocol.Success
}

// CancelOrderRequest 使用请求结构体取消订单
func (s *OrderService) CancelOrderRequest(req *protocol.CancelOrderRequest) protocol.ErrorCode {
	// 获取当前用户信息，确定用户类型
	userService := GetUserService()
	user := userService.GetUserByID(req.UserID)
	userType := protocol.UserTypePassenger
	if user != nil && user.IsDriver() {
		userType = protocol.UserTypeDriver
	}

	// Resolve cancellation reason (backward-compatible)
	reason := s.resolveCancelReason(req, userType)
	if reason == "" {
		return protocol.InvalidParams // must provide either reason or reason_key
	}

	// 保存订单旧状态
	oldOrder := models.GetOrderByID(req.OrderID)
	if oldOrder == nil {
		return protocol.OrderNotFound
	}

	// 调用取消订单的内部方法
	errCode := s.CancelOrder(req.OrderID, req.UserID, reason)

	// 如果取消成功，记录历史
	if errCode == protocol.Success {
		go func() {
			newOrder := models.GetOrderByID(req.OrderID)
			GetOrderHistoryService().RecordOrderCancelled(oldOrder, newOrder, req.UserID, userType, reason)
		}()
	}

	return errCode
}

// resolveCancelReason resolves the cancellation reason from either the new
// reason_key/custom_reason fields or the legacy free-form reason field.
func (s *OrderService) resolveCancelReason(req *protocol.CancelOrderRequest, userType string) string {
	// New flow: reason_key provided
	if req.ReasonKey != "" {
		if req.ReasonKey == "other" {
			if req.CustomReason != "" {
				return req.CustomReason
			}
			return "Other"
		}
		// Validate key and resolve to label
		label := protocol.GetCancelReasonLabel(req.ReasonKey, userType)
		if label != "" {
			return label
		}
		// Unknown key — treat as free-form text
		return req.ReasonKey
	}

	// Legacy flow: free-form reason field
	return req.Reason
}

// CanCancelOrder 检查订单是否可以取消
func (s *OrderService) CanCancelOrder(status string) bool {
	cancellableStates := []string{
		protocol.StatusRequested,
		protocol.StatusAccepted,
		protocol.StatusDriverArrived,
		protocol.StatusDriverComing,
		protocol.StatusInProgress, // allow admin/system to cancel stuck in-progress orders
	}

	return slices.Contains(cancellableStates, status)
}

func (s *OrderService) GetOrderInfoByID(orderId string) *protocol.Order {
	order := models.GetOrderByID(orderId)
	if order == nil {
		return nil
	}
	return s.GetOrderInfo(order)
}

func (s *OrderService) GetOrderInfo(order *models.Order) *protocol.Order {
	info := order.Protocol()
	if detail := models.GetOrderDetail(order.OrderID, order.GetOrderType()); detail != nil {
		info.Details = detail.Protocol()
		if detail.VehicleID != nil {
			vehicle := models.GetVehicleByID(*detail.VehicleID)
			if vehicle != nil {
				info.Vehicle = vehicle.Protocol()
			}
		}
	}
	user := models.GetUserByID(order.GetUserID())
	if user != nil {
		info.Passenger = user.Protocol()
	}
	if order.GetProviderID() != "" {
		provider := models.GetUserByID(order.GetProviderID())
		if provider != nil {
			info.Driver = provider.Protocol()
		}
	}
	ratings := models.GetRatingsByOrderID(order.OrderID)
	for _, rating := range ratings {
		switch rating.RaterType {
		case protocol.UserTypePassenger, protocol.UserTypeUser:
			info.PassengerRatings = append(info.PassengerRatings, rating.Protocol())
		case protocol.UserTypeDriver, protocol.UserTypeProvider:
			info.DriverRatings = append(info.DriverRatings, rating.Protocol())
		}
	}

	return info
}

// GetOrderInfoSanitized returns order info with phone numbers masked based on
// who is requesting. Only the assigned driver (after acceptance) can see the
// passenger phone, and only the passenger (after a driver accepts) can see the
// driver phone. Everyone else sees masked values.
func (s *OrderService) GetOrderInfoSanitized(order *models.Order, requesterID, requesterType string) *protocol.Order {
	info := s.GetOrderInfo(order)
	if info == nil {
		return nil
	}

	// Statuses where contact info is revealed between assigned parties
	revealStatuses := map[string]bool{
		protocol.StatusAccepted:      true,
		protocol.StatusDriverComing:  true,
		protocol.StatusDriverArrived: true,
		protocol.StatusInProgress:    true,
	}

	orderStatus := order.GetStatus()
	isAssignedDriver := requesterType == protocol.UserTypeDriver && order.GetProviderID() == requesterID && order.GetProviderID() != ""
	isPassenger := requesterType == protocol.UserTypePassenger && order.GetUserID() == requesterID

	// Mask passenger phone unless requester is the assigned driver with an active status
	if info.Passenger != nil {
		if !(isAssignedDriver && revealStatuses[orderStatus]) {
			info.Passenger.Phone = utils.MaskPhone(info.Passenger.Phone)
		}
	}

	// Mask driver phone unless requester is the passenger with an active status
	if info.Driver != nil {
		if !(isPassenger && revealStatuses[orderStatus]) {
			info.Driver.Phone = utils.MaskPhone(info.Driver.Phone)
		}
	}

	// Mask phone fields in OrderDetail as well
	if info.Details != nil {
		if !(isAssignedDriver && revealStatuses[orderStatus]) {
			info.Details.PassengerPhone = utils.MaskPhone(info.Details.PassengerPhone)
		}
		if !(isPassenger && revealStatuses[orderStatus]) {
			info.Details.DriverPhone = utils.MaskPhone(info.Details.DriverPhone)
		}
	}

	return info
}

// GetOrderContactInfo returns the counterpart's phone number for calling,
// only if the requester is the assigned driver or the passenger on an active order.
func (s *OrderService) GetOrderContactInfo(req *protocol.OrderContactRequest) (*protocol.OrderContactResponse, protocol.ErrorCode) {
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return nil, protocol.OrderNotFound
	}

	revealStatuses := map[string]bool{
		protocol.StatusAccepted:      true,
		protocol.StatusDriverComing:  true,
		protocol.StatusDriverArrived: true,
		protocol.StatusInProgress:    true,
	}

	if !revealStatuses[order.GetStatus()] {
		return &protocol.OrderContactResponse{Allowed: false}, protocol.Success
	}

	// Requester is the assigned driver -> return passenger phone
	if order.GetProviderID() == req.UserID && order.GetProviderID() != "" {
		passenger := models.GetUserByID(order.GetUserID())
		if passenger != nil {
			return &protocol.OrderContactResponse{
				Allowed: true,
				Phone:   passenger.GetPhone(),
				Name:    passenger.GetFullName(),
			}, protocol.Success
		}
	}

	// Requester is the passenger -> return driver phone
	if order.GetUserID() == req.UserID && order.GetProviderID() != "" {
		driver := models.GetUserByID(order.GetProviderID())
		if driver != nil {
			return &protocol.OrderContactResponse{
				Allowed: true,
				Phone:   driver.GetPhone(),
				Name:    driver.GetDisplayName(),
			}, protocol.Success
		}
	}

	return &protocol.OrderContactResponse{Allowed: false}, protocol.Success
}

// GetOrderETA returns live ETA from the assigned driver's location to the
// pickup (pre-trip) or dropoff (in-progress). Uses Haversine rough estimate;
// Google Directions integration will be added in Phase 3.
func (s *OrderService) GetOrderETA(req *protocol.OrderETARequest) (*protocol.OrderETAResponse, protocol.ErrorCode) {
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return nil, protocol.OrderNotFound
	}

	// Permission: only passenger or assigned driver
	if order.GetUserID() != req.UserID && order.GetProviderID() != req.UserID {
		return nil, protocol.PermissionDenied
	}

	detail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
	if detail == nil {
		return nil, protocol.OrderNotFound
	}

	resp := &protocol.OrderETAResponse{
		OrderID:         order.OrderID,
		PickupLatitude:  detail.GetPickupLatitude(),
		PickupLongitude: detail.GetPickupLongitude(),
		UpdatedAt:       utils.TimeNowMilli(),
		Mode:            "rough",
	}

	// If no driver assigned yet, return estimate from order detail
	if order.GetProviderID() == "" {
		resp.ETAMinutes = detail.GetEstimatedDuration() / 60
		resp.DistanceKm = detail.GetEstimatedDistance()
		return resp, protocol.Success
	}

	driver := models.GetUserByID(order.GetProviderID())
	if driver == nil {
		return resp, protocol.Success
	}

	resp.DriverLatitude = driver.GetLatitude()
	resp.DriverLongitude = driver.GetLongitude()

	// Determine destination based on order status
	var destLat, destLng float64
	switch order.GetStatus() {
	case protocol.StatusAccepted, protocol.StatusDriverComing, protocol.StatusDriverArrived:
		// Driver heading to pickup
		destLat = detail.GetPickupLatitude()
		destLng = detail.GetPickupLongitude()
	case protocol.StatusInProgress:
		// Driver heading to dropoff
		destLat = detail.GetDropoffLatitude()
		destLng = detail.GetDropoffLongitude()
	default:
		return resp, protocol.Success
	}

	// Try Google Directions API for accurate ETA
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	route, err := GetGoogleService().CalculateRidehailingRoute(ctx, driver.GetLatitude(), driver.GetLongitude(), destLat, destLng, false)
	cancel()

	if err == nil && route != nil && route.Duration != nil && route.Duration.Value > 0 {
		etaMin := (route.Duration.Value + 59) / 60 // ceil seconds to minutes
		if etaMin < 1 {
			etaMin = 1
		}
		resp.ETAMinutes = etaMin
		if route.Distance != nil && route.Distance.Value > 0 {
			resp.DistanceKm = float64(route.Distance.Value) / 1000.0
		} else {
			resp.DistanceKm = utils.CalculateDistanceHaversine(driver.GetLatitude(), driver.GetLongitude(), destLat, destLng)
		}
		resp.Mode = "accurate"
		return resp, protocol.Success
	}

	// Fallback to rough estimate
	distKm := utils.CalculateDistanceHaversine(driver.GetLatitude(), driver.GetLongitude(), destLat, destLng)
	etaMin := int(distKm * 2) // rough: 2 min per km
	if etaMin < 1 && distKm > 0 {
		etaMin = 1
	}
	resp.DistanceKm = distKm
	resp.ETAMinutes = etaMin

	return resp, protocol.Success
}

// ============================================================================
// RideOrderService 功能融入 - 网约车订单管理
// ============================================================================

// EstimateOrder 预估订单费用
func (s *OrderService) EstimateOrder(req *protocol.EstimateRequest) (*protocol.OrderPrice, protocol.ErrorCode) {
	// 0. 向后兼容性处理
	if req.VehicleCategory == "" {
		req.VehicleCategory = "sedan" // 默认小车
	}
	if req.VehicleLevel == "" {
		req.VehicleLevel = "economy" // 默认经济型
	}

	// 1. 使用 GoogleService 获取准确的路线信息
	if req.EstimatedDistance == 0 || req.EstimatedDuration == 0 {
		s.EnrichRouteByGoogleMap(req)
	}

	// 2. 使用 PriceRuleService 进行价格计算
	pricingService := GetPriceRuleService()
	req.SnapshotDuration = 30

	// 3. 调用价格引擎进行计算
	snapshot := pricingService.EstimatePrice(req)

	// 保存价格快照
	if err := pricingService.SavePriceSnapshot(snapshot); err != nil {
		return nil, protocol.DatabaseError
	}

	// 4. 从快照转换为 OrderPrice
	orderPrice := snapshot.Protocol()

	return orderPrice, protocol.Success
}

// EnrichRouteByGoogleMap 使用 GoogleService 丰富路线数据；若 Google 不可用或失败则用 Haversine 距离 + 估算时长，避免始终返回最低价
func (s *OrderService) EnrichRouteByGoogleMap(req *protocol.EstimateRequest) {
	googleService := GetGoogleService()
	if googleService != nil {
		avoidTolls := req.VehicleLevel == "economy"
		ctx := context.Background()
		route, err := googleService.CalculateRidehailingRoute(ctx,
			req.PickupLatitude, req.PickupLongitude,
			req.DropoffLatitude, req.DropoffLongitude,
			avoidTolls)
		if err == nil && route != nil {
			if route.Distance != nil {
				req.EstimatedDistance = float64(route.Distance.Value) / 1000.0 // 米 -> 公里
			}
			if route.Duration != nil {
				req.EstimatedDuration = route.Duration.Value / 60 // 秒 -> 分钟
			}
			if req.EstimatedDistance > 0 || req.EstimatedDuration > 0 {
				return
			}
		}
	}

	// Fallback: Haversine 距离 + 粗略时长（约 2.5 分钟/公里，与 app 定价一致，避免总是最低价）
	if req.EstimatedDistance == 0 && (req.PickupLatitude != 0 || req.PickupLongitude != 0) && (req.DropoffLatitude != 0 || req.DropoffLongitude != 0) {
		req.EstimatedDistance = utils.CalculateDistanceHaversine(
			req.PickupLatitude, req.PickupLongitude,
			req.DropoffLatitude, req.DropoffLongitude)
		if req.EstimatedDistance > 0 && req.EstimatedDuration == 0 {
			// 城市路况约 2.5 min/km
			req.EstimatedDuration = int(req.EstimatedDistance*2.5 + 0.5)
			if req.EstimatedDuration < 1 {
				req.EstimatedDuration = 1
			}
			log.Get().Infof("EnrichRouteByGoogleMap: using haversine fallback distance=%.2f km, duration=%d min", req.EstimatedDistance, req.EstimatedDuration)
		}
	}
}

// ValidateRideOrderSnapshot 验证行程订单快照
func (s *OrderService) ValidateRideOrderSnapshot(orderID string, snapshot map[string]interface{}) error {
	order := models.GetRideOrderByOrderID(orderID)
	if order == nil {
		return errors.New("ride order not found")
	}

	// 验证关键字段是否匹配
	if pickup, ok := snapshot["pickup_address"].(string); ok {
		if order.PickupAddress != nil && *order.PickupAddress != pickup {
			return errors.New("pickup address mismatch")
		}
	}

	if dropoff, ok := snapshot["dropoff_address"].(string); ok {
		if order.DropoffAddress != nil && *order.DropoffAddress != dropoff {
			return errors.New("dropoff address mismatch")
		}
	}

	return nil
}

func (s *OrderService) ArrivedOrder(req *protocol.OrderActionRequest) protocol.ErrorCode {
	user := models.GetUserByID(req.UserID)
	if !user.IsDriver() {
		return protocol.AccessDenied
	}
	if req.OrderID == "" {
		return protocol.InvalidParams
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return protocol.OrderNotFound
	}
	if order.GetProviderID() != user.UserID {
		return protocol.AccessDenied
	}
	if order.GetStatus() == protocol.StatusDriverArrived {
		return protocol.Success
	}
	if order.GetStatus() != protocol.StatusDriverComing && order.GetStatus() != protocol.StatusAccepted {
		return protocol.InvalidRideStatus // 司机状态不正确，无法标记到达
	}
	//检查当前司机是否有其他进行中的订单
	if s.CountActiveRideOrdersByDriver(user.UserID, order.OrderID) > 0 {
		return protocol.DriverHasActiveOrderInProgress // 司机有在途订单，不能开启新行程
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		orderValues := models.OrderValues{}
		orderValues.SetStatus(protocol.StatusDriverArrived)
		if err := models.UpdateOrder(tx, order, &orderValues); err != nil {
			return err
		}
		detail := models.OrderDetailValues{}
		detail.SetDriverArrived()
		return tx.Table(models.GetTableNameByOrderType(order.GetOrderType())).Where("order_id = ?", order.OrderID).UpdateColumns(detail).Error
	})
	if err != nil {
		return protocol.DatabaseError
	}
	go GetUserService().RefreshDriverOrderQueue(req.UserID)
	// 发送FCM通知
	go s.NotifyDriverArrived(req.OrderID)

	return protocol.Success
}

func (s *OrderService) StartOrder(req *protocol.OrderActionRequest) protocol.ErrorCode {
	user := models.GetUserByID(req.UserID)
	if !user.IsDriver() {
		return protocol.AccessDenied
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder || order.GetProviderID() != user.UserID {
		return protocol.OrderNotFound
	}
	if order.GetStatus() == protocol.StatusInProgress {
		return protocol.Success
	}
	// 司机还未到达，无法开始行程
	if !slices.Contains([]string{protocol.StatusDriverArrived, protocol.StatusAccepted}, order.GetStatus()) {
		return protocol.InvalidRideStatus
	}
	//检查当前司机是否有其他进行中的订单
	if s.CountActiveRideOrdersByDriver(user.UserID, order.OrderID) > 0 {
		return protocol.DriverHasActiveOrderInProgress // 司机有在途订单，不能开启新行程
	}
	orderValues := models.OrderValues{}
	orderValues.StartOrder()
	if err := models.UpdateOrder(models.DB, order, &orderValues); err != nil {
		return protocol.DatabaseError
	}
	go GetUserService().RefreshDriverOrderQueue(req.UserID)
	// 发送FCM通知
	go s.NotifyTripStarted(req.OrderID)

	// 记录订单开始历史
	go func() {
		// 获取更新后的订单
		updatedOrder := models.GetOrderByID(req.OrderID)
		if updatedOrder != nil {
			GetOrderHistoryService().RecordOrderStarted(order, updatedOrder, req.UserID)
		}
	}()

	return protocol.Success
}

func (s *OrderService) FinishOrder(req *protocol.OrderActionRequest) protocol.ErrorCode {
	user := models.GetUserByID(req.UserID)
	if !user.IsDriver() {
		return protocol.AccessDenied
	}
	if req.OrderID == "" {
		return protocol.InvalidParams
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetOrderType() != protocol.RideOrder || order.GetProviderID() != user.UserID {
		return protocol.OrderNotFound
	}
	if order.GetStatus() == protocol.StatusDriverArrived {
		return protocol.Success
	}
	if order.GetStatus() != protocol.StatusInProgress {
		return protocol.RideNotStarted // 行程还没开始
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		orderValues := models.OrderValues{}
		orderValues.FinishOrder()
		return models.UpdateOrder(tx, order, &orderValues)
	})
	if err != nil {
		return protocol.DatabaseError
	}

	go GetUserService().RefreshDriverOrderQueue(req.UserID)
	// 发送FCM通知
	go s.NotifyTripEnded(req.OrderID)

	// Increment ride counts for both driver and passenger
	go s.incrementRideCountsForOrder(order)

	// 记录订单完成历史
	go func() {
		// 获取更新后的订单
		updatedOrder := models.GetOrderByID(req.OrderID)
		if updatedOrder != nil {
			GetOrderHistoryService().RecordOrderFinished(order, updatedOrder, req.UserID)
		}
	}()

	return protocol.Success
}

// GetActiveRideOrderByUser gets active ride order for a user
func (s *OrderService) GetActiveRideOrderByUser(userID string) (*models.Order, error) {
	var order models.Order
	err := models.DB.Where("user_id = ? AND status IN (?)", userID, []string{protocol.StatusRequested, protocol.StatusPending, protocol.StatusAccepted, protocol.StatusInProgress}).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CountActiveRideOrdersByUser counts active ride orders for a user
func (s *OrderService) CountActiveRideOrdersByUser(userID string) int64 {
	var count int64
	err := models.DB.Model(&models.Order{}).Where("user_id = ? AND status IN (?)", userID, []string{protocol.StatusRequested, protocol.StatusPending, protocol.StatusAccepted, protocol.StatusInProgress}).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}
func (s *OrderService) CountActiveRideOrdersByDriver(userID, orderID string) int64 {
	var count int64
	err := models.DB.Model(&models.Order{}).Where("provider_id = ? AND status IN (?) and order_id != ?", userID, []string{protocol.StatusDriverArrived, protocol.StatusInProgress}, orderID).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func (s *OrderService) GenerateOrderIDByType(orderType string) string {
	switch orderType {
	case protocol.RideOrder:
		return utils.GenerateRideOrderID()
	default:
		return utils.GenerateOrderID()
	}
}

// CreateOrder creates a new order with details from request
func (s *OrderService) CreateOrder(req *protocol.CreateOrderRequest) (*protocol.Order, protocol.ErrorCode) {
	log.Get().Infof("OrderService.CreateOrder: 开始创建订单，UserID=%s, PriceID=%s", req.UserID, req.PriceID)

	// 调用价格验证和锁定逻辑
	pricingService := GetPriceRuleService()
	price, errCode := pricingService.ValidateAndLockPriceID(req.PriceID)
	if errCode != protocol.Success {
		log.Get().Errorf("OrderService.CreateOrder: 价格验证失败，ErrorCode=%s", errCode)
		return nil, errCode
	}

	log.Get().Infof("OrderService.CreateOrder: 价格验证成功，OrderType=%s, SnapshotUserID=%s", price.GetOrderType(), price.GetUserID())

	// 根据订单类型创建相应的详情
	switch price.GetOrderType() {
	case protocol.RideOrder:
		log.Get().Infof("OrderService.CreateOrder: 订单类型为RideOrder，检查用户活跃订单")
		// 检查用户是否有活跃的网约车订单
		activeOrderCount := s.CountActiveRideOrdersByUser(req.UserID)
		if activeOrderCount > 0 {
			log.Get().Warnf("OrderService.CreateOrder: 用户有活跃订单，数量=%d", activeOrderCount)
			return nil, protocol.RideInProgress
		}
		log.Get().Infof("OrderService.CreateOrder: 用户没有活跃订单，继续创建")
	default:
		log.Get().Errorf("OrderService.CreateOrder: 不支持的订单类型=%s", price.GetOrderType())
		// 其他订单类型暂时不支持
		return nil, protocol.InvalidParams
	}
	// 从快照metadata获取业务参数
	snapshotMeta := price.GetMetadata()
	nowtime := utils.TimeNowMilli()
	// 创建订单metadata，只记录必要信息
	metadata := map[string]any{}
	// 只记录快照ID，不重复存储快照中的坐标等信息
	metadata["price_id"] = price.SnapshotID

	// 检查用户是否为sandbox用户
	user := models.GetUserByID(req.UserID)
	sandbox := 0
	if user != nil && user.IsSandbox() {
		sandbox = 1 // 如果用户是sandbox用户，设置订单的sandbox=1
	}

	// 创建订单对象
	order := &models.Order{
		OrderID: s.GenerateOrderIDByType(price.GetOrderType()),
		Salt:    utils.GenerateSalt(),
		OrderValues: &models.OrderValues{
			PromoDiscount:       price.PromoDiscount,
			UserPromoDiscount:   price.UserPromoDiscount,
			UserPromotionIDs:    price.UserPromotionIDs,
			OrderDispatchValues: &models.OrderDispatchValues{},
		},
	}

	// 确定要使用的用户ID：优先使用请求中的UserID，如果为空则使用快照中的UserID
	userID := req.UserID
	if userID == "" {
		userID = price.GetUserID()
	}
	orderConfig := config.Get().Order
	expiredAt := time.Now().Add(time.Duration(orderConfig.RideOrder.ExpireMinutes) * time.Minute).UnixMilli()

	order.SetScheduleType(protocol.ScheduleTypeInstant).
		SetStatus(protocol.StatusRequested).
		SetPaymentStatus(protocol.StatusPending).
		SetOrderType(price.GetOrderType()).
		SetUserID(userID).
		SetOriginalAmount(price.GetOriginalFare()).
		SetDiscountedAmount(price.GetDiscountedFare()).
		SetPaymentAmount(price.GetDiscountedFare()).
		SetTotalDiscountAmount(price.GetDiscountAmount()).
		SetCurrency(price.GetCurrency()).
		SetScheduledAt(price.GetScheduledAt()).
		SetNotes(req.Notes).
		SetMetadata(metadata).
		SetSandbox(sandbox).
		SetExpiredAt(expiredAt)

	//15分钟以上的预定视为预约订单
	if order.GetScheduledAt()-nowtime > int64(15*time.Minute) {
		order.SetScheduleType(protocol.ScheduleTypeScheduled)
	}
	dispatchCfg := config.Get().Dispatch
	order.SetAutoDispatchEnabled(dispatchCfg.Enabled)
	if order.GetAutoDispatchEnabled() {
		order.SetCurrentRound(1).
			SetDispatchStatus(protocol.StatusPending).
			SetMaxRounds(dispatchCfg.MaxRounds)
	}

	// Manual driver selection (optional):
	// If provider_id is specified, pre-assign the order to that driver and send dispatch only to them.
	selectedProviderID := strings.TrimSpace(req.ProviderID)
	var manualDispatchRecord *models.DispatchRecord
	if selectedProviderID != "" {
		driver := models.GetUserByID(selectedProviderID)
		if driver == nil {
			return nil, protocol.UserNotFound
		}
		if !driver.IsDriver() {
			return nil, protocol.InvalidUserType
		}
		if driver.GetStatus() != protocol.StatusActive {
			return nil, protocol.UserNotActive
		}
		if driver.GetOnlineStatus() != protocol.StatusOnline {
			return nil, protocol.DriverOffline
		}
		if s.CountActiveRideOrdersByDriver(selectedProviderID, "") > 0 {
			return nil, protocol.DriverHasActiveOrder
		}

		// Pre-assign provider and disable auto-dispatch to avoid notifying others.
		order.SetProviderID(selectedProviderID).
			SetAutoDispatchEnabled(false).
			SetCurrentRound(1).
			SetDispatchStatus(protocol.StatusPending).
			SetMaxRounds(1)
	}

	err := models.GetDB().Transaction(func(tx *gorm.DB) error {
		// 创建主订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 根据订单类型创建详情表
		switch price.GetOrderType() {
		case protocol.RideOrder:
			// 创建网约车详情，优先从快照获取信息
			rideOrder := models.NewRideOrder(order.OrderID)
			rideOrder.SetVehicleCategory(price.GetVehicleCategory()).
				SetVehicleLevel(price.GetVehicleLevel()).
				SetEstimatedDistance(price.GetDistance()).
				SetEstimatedDuration(price.GetDuration()).
				SetPassengerCount(snapshotMeta.GetInt("passenger_count")).
				SetPickupAddress(snapshotMeta.Get("pickup_address")).
				SetPickupLatitude(snapshotMeta.GetFloat64("pickup_latitude")).
				SetPickupLongitude(snapshotMeta.GetFloat64("pickup_longitude")).
				SetPickupLandmark(snapshotMeta.Get("pickup_landmark")).
				SetDropoffAddress(snapshotMeta.Get("dropoff_address")).
				SetDropoffLatitude(snapshotMeta.GetFloat64("dropoff_latitude")).
				SetDropoffLongitude(snapshotMeta.GetFloat64("dropoff_longitude")).
				SetDropoffLandmark(snapshotMeta.Get("dropoff_landmark"))

			if err := tx.Create(rideOrder).Error; err != nil {
				return err
			}
		}
		// If manually selecting a driver, create exactly one dispatch record for that driver.
		if selectedProviderID != "" {
			dispatchedAt := utils.TimeNowMilli()
			manualDispatchRecord = &models.DispatchRecord{
				DriverID:             selectedProviderID,
				OrderID:              order.OrderID,
				DispatchID:           utils.GenerateDispatchID(),
				Round:                1,
				DispatchedAt:         dispatchedAt,
				ExpiredAt:            order.GetScheduledAt(),
				RoundSeq:             1,
				DispatchRecordValues: &models.DispatchRecordValues{},
				CreatedAt:            dispatchedAt,
			}
			if err := tx.Create(manualDispatchRecord).Error; err != nil {
				return err
			}
		}

		// 更新价格快照的订单关联
		if err := models.UpdatePriceOrderID(tx, price.SnapshotID, order.OrderID); err != nil {
			return err
		}

		// 将使用的用户优惠券标记为已使用
		for _, item := range price.GetBreakdowns() {
			if item.Category != protocol.PriceRuleCategoryUserPromotion {
				continue
			}
			if err := models.UseUserPromotionByID(tx, item.RuleID, order.OrderID, item.Amount); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, protocol.DatabaseError
	}
	orderInfo := s.GetOrderInfo(order)

	go GetUserService().RefreshUserOrderQueue(req.UserID)
	// Dispatch:
	// - manual: notify selected provider only
	// - auto: start auto dispatch
	if manualDispatchRecord != nil {
		go GetDispatchService().SendDispatchNotifications(manualDispatchRecord)
	} else {
		//开始自动派单（异步）
		go s.DispatchOrder(orderInfo)
	}
	// 记录订单创建历史
	go GetOrderHistoryService().RecordOrderCreated(order, req.UserID)

	return orderInfo, protocol.Success
}

func (s *OrderService) DispatchOrderByID(orderID string) (result *protocol.DispatchResult) {
	order := s.GetOrderInfoByID(orderID)
	return s.DispatchOrder(order)
}

func (s *OrderService) DispatchOrder(order *protocol.Order) (result *protocol.DispatchResult) {
	result = &protocol.DispatchResult{
		Success: false,
	}
	if !order.AutoDispatchEnabled {
		result.Message = "Auto dispatch not enabled for this order"
		return
	}
	if order.Details == nil {
		order = s.GetOrderInfoByID(order.OrderID)
	}
	// 这里是派单的具体逻辑
	//time.Sleep(10 * time.Second)
	GetDispatchService().StartAutoDispatch(order)
	return
}

// AcceptOrder accepts an order
func (s *OrderService) AcceptOrder(req *protocol.OrderActionRequest) protocol.ErrorCode {
	user := models.GetUserByID(req.UserID)
	// 检查用户类型
	if !user.IsDriver() {
		return protocol.PermissionDenied
	}
	if user.GetStatus() != protocol.StatusActive || user.IsDeleted() {
		return protocol.AccountDisabled
	}
	if user.GetOnlineStatus() != protocol.StatusOnline {
		return protocol.DriverOffline // 司机未上线，无法接单
	}
	if models.CountProcessingOrdersByUserID(req.UserID) > 3 {
		return protocol.DriverHasActiveOrder // 司机有未完成的订单，无法接单
	}
	var order *models.Order
	if req.DispatchId != "" {
		dispatch := models.GetDispatchByID(req.DispatchId)
		if dispatch != nil {
			order = models.GetOrderByID(dispatch.OrderID)
		}
	} else if req.OrderID != "" {
		order = models.GetOrderByID(req.OrderID)
	}
	// 先获取要接受的订单信息，确定订单类型
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return protocol.OrderNotFound
	}

	// 检查订单状态是否可以接单
	if order.GetStatus() != protocol.StatusRequested {
		// 根据当前状态返回具体错误
		switch order.GetStatus() {
		case protocol.StatusAccepted:
			return protocol.RideAlreadyBooked // 订单已被接单
		case protocol.StatusCancelled:
			return protocol.RideAlreadyCancelled // 订单已取消
		case protocol.StatusCompleted:
			return protocol.RideAlreadyCompleted // 订单已完成
		default:
			return protocol.InvalidRideStatus // 订单状态不允许接单
		}
	}

	// 检查订单是否已被其他司机接单
	if order.GetProviderID() != "" && order.GetProviderID() != req.UserID {
		return protocol.RideAlreadyBooked // 订单已被其他司机接单
	}
	vehicle := models.GetVehicleByDriverID(req.UserID)
	if vehicle == nil || !vehicle.IsAvailable() {
		return protocol.VehicleNotAssigned // 司机没有分配车辆，无法接单
	}
	orderValues := models.OrderValues{}
	orderValues.SetAcceptedAt(utils.TimeNowMilli()).
		SetStatus(protocol.StatusAccepted).
		SetProviderID(req.UserID)
	// 使用事务确保订单和车辆信息的一致性
	hasAccepted := true
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 更新主订单表 (also allow pre-assigned driver to accept)
		rs := tx.Model(&models.Order{}).
			Where("order_id = ?", order.OrderID).
			Where("status = ?", protocol.StatusRequested).
			Where("provider_id IS NULL OR provider_id = '' OR provider_id = ?", req.UserID).
			UpdateColumns(orderValues)
		if rs.Error != nil {
			return rs.Error
		}
		if rs.RowsAffected == 0 {
			hasAccepted = false
			return nil
		}

		// 更新网约车订单表的车辆ID
		return tx.Model(&models.RideOrder{}).
			Where("order_id = ?", req.OrderID).
			Update("vehicle_id", vehicle.VehicleID).Error
	})
	if err != nil {
		return protocol.DatabaseError
	}
	if !hasAccepted {
		return protocol.RideAlreadyBooked // 订单已被其他司机接单
	}
	models.RefreshOrderCache(order.OrderID)
	go GetUserService().RefreshDriverOrderQueue(req.UserID)
	if req.DispatchId != "" {
		// 如果是调度系统分配的订单，更新调度记录
		GetDispatchService().HandleDriverAccept(req.DispatchId, req.UserID, req.Latitude, req.Longitude)
	}
	// 发送FCM通知给乘客（司机已接单）
	go s.NotifyOrderAccepted(order.OrderID)

	// 记录订单接受历史
	go func() {
		GetOrderHistoryService().RecordOrderAccepted(order, order, req.UserID)
	}()

	return protocol.Success
}

func (s *OrderService) RejectOrder(req *protocol.OrderActionRequest) protocol.ErrorCode {
	user := models.GetUserByID(req.UserID)
	// 检查用户类型
	userType := user.GetUserType()
	if userType != protocol.UserTypeDriver {
		return protocol.PermissionDenied
	}

	// 先获取要接受的订单信息，确定订单类型
	var order *models.Order
	if req.DispatchId != "" {
		dispatch := models.GetDispatchByID(req.DispatchId)
		if dispatch != nil {
			order = models.GetOrderByID(dispatch.OrderID)
		}
	} else if req.OrderID != "" {
		order = models.GetOrderByID(req.OrderID)
	}
	// 先获取要接受的订单信息，确定订单类型
	if order == nil || order.GetOrderType() != protocol.RideOrder {
		return protocol.OrderNotFound
	}

	// 检查订单状态是否可以接单
	if order.GetStatus() != protocol.StatusRequested {
		// 根据当前状态返回具体错误
		switch order.GetStatus() {
		case protocol.StatusAccepted:
			return protocol.RideAlreadyBooked // 订单已被接单
		case protocol.StatusCancelled:
			return protocol.RideAlreadyCancelled // 订单已取消
		case protocol.StatusCompleted:
			return protocol.RideAlreadyCompleted // 订单已完成
		default:
			return protocol.InvalidRideStatus // 订单状态不允许接单
		}
	}

	if req.DispatchId != "" {
		// 如果是调度系统分配的订单，更新调度记录
		GetDispatchService().HandleDriverAccept(req.DispatchId, req.UserID, req.Latitude, req.Longitude)
	}

	// 记录订单拒绝历史
	go func() {
		GetOrderHistoryService().RecordOrderRejected(order, order, req.UserID, req.RejectReason)
	}()

	return protocol.Success
}

// GetNearbyOrders gets nearby orders for the requesting driver. It returns:
// 1) Orders with status=requested, no provider_id, that are NOT in any pending dispatch (broadcast-style).
// 2) Orders that have a pending dispatch TO this driver (so drivers actually receive requests sent to them).
// For (2), each order includes dispatch_id so the app can send it on accept for first-accept-wins.
func (s *OrderService) GetNearbyOrders(req *protocol.GetNearbyOrdersRequest) (*protocol.GetNearbyOrdersResponse, protocol.ErrorCode) {
	seen := make(map[string]bool)
	var orderList []*protocol.Order

	// 1) Broadcast-style: requested, no provider, no pending dispatch for anyone
	subNoPending := models.GetDB().Model(&models.DispatchRecord{}).Select("order_id").Where("status = ?", protocol.StatusPending)
	query := models.GetDB().Model(&models.Order{}).
		Where("order_type = ?", req.OrderType).
		Where("status = ?", protocol.StatusRequested).
		Where("t_orders.provider_id IS NULL OR t_orders.provider_id = ''").
		Where("t_orders.order_id NOT IN (?)", subNoPending)
	if req.Radius != 0 && req.Latitude != 0 && req.Longitude != 0 {
		minLat, maxLat, minLng, maxLng := utils.CalculateCoordinateRange(req.Latitude, req.Longitude, req.Radius)
		query = query.Joins("JOIN t_ride_orders ON t_orders.order_id = t_ride_orders.order_id").
			Where("t_ride_orders.pickup_latitude BETWEEN ? AND ?", minLat, maxLat).
			Where("t_ride_orders.pickup_longitude BETWEEN ? AND ?", minLng, maxLng)
	}
	var broadcastIDs []string
	query.Limit(req.Limit).Order("t_orders.created_at DESC").Pluck("t_orders.order_id", &broadcastIDs)
	for _, orderID := range broadcastIDs {
		if seen[orderID] {
			continue
		}
		seen[orderID] = true
		order := models.GetOrderByID(orderID)
		if order == nil {
			continue
		}
		info := s.GetOrderInfoSanitized(order, req.RequesterID, protocol.UserTypeDriver)
		if info != nil {
			orderList = append(orderList, info)
		}
	}

	// 2) Dispatched to this driver: pending dispatch records for req.RequesterID
	var dispatchedToMe []struct {
		OrderID    string
		DispatchID string
	}
	models.GetDB().Model(&models.DispatchRecord{}).
		Select("order_id, dispatch_id").
		Where("driver_id = ?", req.RequesterID).
		Where("status = ?", protocol.StatusPending).
		Find(&dispatchedToMe)
	for _, row := range dispatchedToMe {
		if seen[row.OrderID] {
			continue
		}
		order := models.GetOrderByID(row.OrderID)
		if order == nil || order.GetStatus() != protocol.StatusRequested {
			continue
		}
		seen[row.OrderID] = true
		info := s.GetOrderInfoSanitized(order, req.RequesterID, protocol.UserTypeDriver)
		if info != nil {
			info.DispatchID = row.DispatchID
			orderList = append(orderList, info)
		}
	}

	// Keep a reasonable total (broadcast first, then dispatched-to-me)
	if len(orderList) > req.Limit {
		orderList = orderList[:req.Limit]
	}

	response := &protocol.GetNearbyOrdersResponse{
		Orders: orderList,
		Count:  len(orderList),
	}
	return response, protocol.Success
}

// PrepareCashPayment stores passenger-generated cash verification code on the order.
func (s *OrderService) PrepareCashPayment(req *protocol.OrderCashRequest) (*protocol.OrderCashResponse, protocol.ErrorCode) {
	user := models.GetUserByID(req.UserID)
	if user == nil {
		return nil, protocol.UserNotFound
	}
	if !user.IsPassenger() {
		return nil, protocol.PermissionDenied
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil || order.GetUserID() != user.UserID {
		return nil, protocol.OrderNotFound
	}
	if order.GetStatus() != protocol.StatusTripEnded {
		return nil, protocol.InvalidRideStatus
	}
	if order.GetPaymentStatus() == protocol.StatusSuccess {
		return &protocol.OrderCashResponse{
			OrderID:       order.OrderID,
			Status:        protocol.StatusSuccess,
			PaymentMethod: protocol.PaymentMethodCash,
			CashCode:      "",
		}, protocol.Success
	}

	code := strings.TrimSpace(req.CashCode)
	if len(code) < 4 || len(code) > 8 {
		return nil, protocol.InvalidParams
	}
	for _, ch := range code {
		if ch < '0' || ch > '9' {
			return nil, protocol.InvalidParams
		}
	}

	metadata := order.GetMetadata()
	if metadata == nil {
		metadata = map[string]any{}
	}
	metadata["cash_verification_code"] = code
	metadata["cash_requested_at"] = utils.TimeNowMilli()

	values := &models.OrderValues{}
	values.
		SetPaymentMethod(protocol.PaymentMethodCash).
		SetPaymentStatus(protocol.StatusPending).
		SetPaymentID("").
		SetChannelPaymentID("").
		SetPaymentRedirectURL("").
		SetPaymentResult("").
		SetMetadata(metadata)
	if err := models.UpdateOrder(models.DB, order, values); err != nil {
		return nil, protocol.DatabaseError
	}

	return &protocol.OrderCashResponse{
		OrderID:       order.OrderID,
		Status:        protocol.StatusPending,
		PaymentMethod: protocol.PaymentMethodCash,
		CashCode:      code,
	}, protocol.Success
}

func (s *OrderService) OrderPayment(req *protocol.OrderPaymentRequest) (result *protocol.OrderPaymentResult, errCode protocol.ErrorCode) {
	errCode = protocol.Success
	user := models.GetUserByID(req.UserID)
	if user == nil {
		errCode = protocol.UserNotFound
		return
	}
	// 设置用户ID
	req.UserID = user.UserID
	if req.PaymentMethod == protocol.PaymentMethodCash && !user.IsDriver() {
		errCode = protocol.PermissionDenied
		return
	}
	if req.PaymentMethod != protocol.PaymentMethodCash && !user.IsPassenger() {
		errCode = protocol.PermissionDenied
		return
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil || (order.GetProviderID() != user.UserID && order.GetUserID() != user.UserID) {
		errCode = protocol.OrderNotFound
		return
	}
	if req.PaymentMethod == protocol.PaymentMethodCash && user.IsDriver() {
		expectedCode := ""
		if rawCode, ok := order.GetMetadata()["cash_verification_code"]; ok && rawCode != nil {
			expectedCode = strings.TrimSpace(fmt.Sprintf("%v", rawCode))
		}
		if expectedCode != "" {
			givenCode := strings.TrimSpace(req.CashCode)
			if givenCode == "" || givenCode != expectedCode {
				errCode = protocol.InvalidParams
				return
			}
		}
	}
	if order.GetPaymentStatus() == protocol.StatusSuccess {
		result = &protocol.OrderPaymentResult{
			OrderID: order.OrderID,
			Status:  order.GetPaymentStatus(),
		}
		errCode = protocol.Success
		return
	} else if order.GetPaymentMethod() == req.PaymentMethod && order.GetPaymentStatus() == protocol.StatusPending && order.GetPaymentRedirectURL() != "" {
		result = &protocol.OrderPaymentResult{
			RedirectURL: order.GetPaymentRedirectURL(),
			OrderID:     order.OrderID,
			Status:      order.GetPaymentStatus(),
		}
		errCode = protocol.Success
		return
	}
	values := &models.OrderValues{}
	values.SetPaymentMethod(req.PaymentMethod)
	cresult, errCode := GetPaymentService().OrderPayment(&ChannelPaymentRequest{
		Phone:         req.Phone,
		Email:         req.Email,
		AccountNo:     req.AccountNo,
		AccountName:   req.AccountName,
		PaymentMethod: req.PaymentMethod,
		Order:         order,
		User:          user,
	})
	if errCode != protocol.Success {
		return
	}

	values.SetPaymentMethod(req.PaymentMethod).
		SetPaymentID(cresult.PaymentID).
		SetChannelPaymentID(cresult.ChannelPaymentID).
		SetPaymentRedirectURL(cresult.RedirectURL).
		SetPaymentResult("")
	switch cresult.Status {
	case protocol.StatusSuccess:
		values.CompleteOrder().
			SetPaymentStatus(protocol.StatusSuccess)
		if req.PaymentMethod == protocol.PaymentMethodCash {
			metadata := order.GetMetadata()
			delete(metadata, "cash_verification_code")
			delete(metadata, "cash_requested_at")
			values.SetMetadata(metadata)
		}
	case protocol.StatusFailed:
		// 根据req.Language，判断ResCode是否预定义的结果码，如果是，PaymentResult就用[ResCode]对应翻译
		paymentResult := s.getTranslatedPaymentResult(cresult.ResCode, cresult.ResMsg, req.Language)
		values.SetPaymentResult(paymentResult)
	}
	if err := models.UpdateOrder(models.DB, order, values); err != nil {
		log.Get().Errorf(
			"OrderPayment 更新订单支付信息失败: order_id=%s payment_method=%s payment_id=%s channel_payment_id=%s status=%s error=%v",
			order.OrderID,
			req.PaymentMethod,
			cresult.PaymentID,
			cresult.ChannelPaymentID,
			cresult.Status,
			err,
		)
		errCode = protocol.DatabaseError
		return
	}
	// 支付成功时，发送支付确认消息
	if order.GetPaymentStatus() == protocol.StatusSuccess {
		go s.NotifyPaymentConfirmed(req.OrderID)
		go s.incrementRideCountsForOrder(order)
	}

	result = &protocol.OrderPaymentResult{
		RedirectURL: order.GetPaymentRedirectURL(),
		OrderID:     order.OrderID,
		Status:      order.GetPaymentStatus(),
		Reason:      order.GetPaymentResult(),
	}
	return
}

func (s *OrderService) CheckOrderPayment(order_id, payment_id string) {
	order := models.GetOrderByID(order_id)
	if order == nil {
		return
	}
	if order.GetPaymentStatus() == protocol.StatusSuccess || order.GetPaymentStatus() == protocol.StatusFailed {
		return
	}
	if order.GetPaymentMethod() == protocol.PaymentMethodCash {
		return
	}
	payment := models.GetPaymentByID(payment_id)
	if payment == nil {
		payment = models.GetLastPaymentByOrderID(order_id)
	}
	if payment == nil {
		return
	}
	if order.GetPaymentStatus() == payment.GetStatus() {
		return
	}
	values := &models.OrderValues{}
	values.SetPaymentRedirectURL(payment.GetRedirectURL())
	switch payment.GetStatus() {
	case protocol.StatusSuccess:
		values.SetPaymentStatus(protocol.StatusSuccess).
			SetStatus(protocol.StatusCompleted).
			SetCompletedAt(payment.GetCompletedAt()).
			SetPaymentID(payment.PaymentID)
		// 处理成功状态
	case protocol.StatusFailed:
		// 处理失败状态
		values.SetPaymentStatus(protocol.StatusFailed).
			SetPaymentID(payment.PaymentID).
			SetPaymentResult(fmt.Sprintf("[%v]%v", payment.GetResCode(), payment.GetResMsg()))
	}
	if err := models.UpdateOrder(models.DB, order, values); err != nil {
		log.Get().Errorf("OrderService.CheckOrderPayment: 更新订单支付状态失败, OrderID=%s, Error=%v", order.OrderID, err)
		return
	}
	// 支付成功时，发送支付确认消息
	if order.GetPaymentStatus() == protocol.StatusSuccess {
		go s.NotifyPaymentConfirmed(order.OrderID)
		go s.incrementRideCountsForOrder(order)
	}
	log.Get().Info("OrderService.CheckOrderPayment: 订单支付状态已更新", "OrderID", order.OrderID, "PaymentStatus", order.GetPaymentStatus())
}

// incrementRideCountsForOrder increments total_rides for both driver and passenger when an order completes.
func (s *OrderService) incrementRideCountsForOrder(order *models.Order) {
	if order == nil {
		return
	}
	// Increment driver ride count
	if providerID := order.GetProviderID(); providerID != "" {
		driver := models.GetUserByID(providerID)
		if driver != nil {
			driverValues := &models.UserValues{}
			driverValues.SetTotalRides(driver.GetTotalRides() + 1)
			if err := models.GetDB().Model(&models.User{}).Where("user_id = ?", providerID).UpdateColumns(driverValues).Error; err != nil {
				log.Get().Warnf("Failed to increment ride count for driver %s: %v", providerID, err)
			}
		}
		// Increment vehicle ride count
		vehicle := models.GetVehicleByDriverID(providerID)
		if vehicle != nil {
			vehicle.IncrementRideCount()
			if err := models.GetDB().Model(&models.Vehicle{}).Where("vehicle_id = ?", vehicle.VehicleID).UpdateColumns(vehicle.VehicleValues).Error; err != nil {
				log.Get().Warnf("Failed to increment ride count for vehicle %s: %v", vehicle.VehicleID, err)
			}
		}
	}
	// Increment passenger ride count
	if userID := order.GetUserID(); userID != "" {
		passenger := models.GetUserByID(userID)
		if passenger != nil {
			passengerValues := &models.UserValues{}
			passengerValues.SetTotalRides(passenger.GetTotalRides() + 1)
			if err := models.GetDB().Model(&models.User{}).Where("user_id = ?", userID).UpdateColumns(passengerValues).Error; err != nil {
				log.Get().Warnf("Failed to increment ride count for passenger %s: %v", userID, err)
			}
		}
	}
}

// ============================================================================
// FCM 推送通知功能
// ============================================================================

// NotifyPassenger 通知乘客 - 接收订单对象作为参数
func (s *OrderService) NotifyPassenger(order *models.Order, notificationType string) error {
	if order == nil {
		return errors.New("order is nil")
	}

	if order.GetUserID() == "" {
		return errors.New("passenger not found in order")
	}

	// 获取乘客信息
	passenger := models.GetUserByID(order.GetUserID())
	if passenger == nil {
		return errors.New("passenger not found")
	}

	// 获取订单详情
	orderDetail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
	if orderDetail == nil {
		return errors.New("order detail not found")
	}

	// 根据通知类型映射到消息类型
	var msgType string
	switch notificationType {
	case protocol.NotificationTypeOrderAccepted:
		msgType = protocol.MsgTypePassengerOrderAccepted
	case protocol.NotificationTypeDriverArrived:
		msgType = protocol.MsgTypePassengerDriverArrived
	case protocol.NotificationTypeTripStarted:
		msgType = protocol.MsgTypePassengerTripStarted
	case protocol.NotificationTypeTripEnded:
		msgType = protocol.MsgTypePassengerTripEnded
	case protocol.NotificationTypePaymentConfirmed:
		msgType = protocol.MsgTypePassengerPaymentConfirmed
	case protocol.NotificationTypeOrderCancelled:
		msgType = protocol.MsgTypePassengerOrderCancelled
	default:
		return fmt.Errorf("unsupported notification type for passenger: %s", notificationType)
	}

	// 准备消息参数
	params := map[string]any{
		"to":                order.GetUserID(),
		"OrderID":           order.OrderID,
		"OrderStatus":       order.GetStatus(),
		"PickupAddress":     orderDetail.GetPickupAddress(),
		"DropoffAddress":    orderDetail.GetDropoffAddress(),
		"Amount":            order.GetPaymentAmount().StringFixed(2),
		"Currency":          order.GetCurrency(),
		"msg_type":          protocol.FCMMessageTypeOrder,
		"notification_type": notificationType,
		"order_status":      order.GetStatus(),
	}

	// 添加司机信息
	if order.GetProviderID() != "" {
		driver := models.GetUserByID(order.GetProviderID())
		if driver != nil {
			params["DriverName"] = driver.GetDisplayName()

			// 添加车辆信息
			if orderDetail.GetVehicleID() != "" {
				vehicle := models.GetVehicleByID(orderDetail.GetVehicleID())
				if vehicle != nil {
					params["PlateNumber"] = vehicle.GetPlateNumber()
				}
			}
		}
	}

	// 添加取消信息
	if notificationType == protocol.NotificationTypeOrderCancelled {
		if order.CancelledBy != nil {
			cancellerName := "System"
			if *order.CancelledBy == order.GetUserID() {
				cancellerName = "You"
			} else if *order.CancelledBy == order.GetProviderID() {
				cancellerName = "Driver"
			}
			params["CancelledBy"] = cancellerName
		}
		if order.CancelReason != nil {
			params["CancelReason"] = *order.CancelReason
		} else {
			params["CancelReason"] = "No reason provided"
		}
	}

	// Include ETA metadata in notification if present (set by NotifyOrderAccepted)
	if meta := order.GetMetadata(); meta != nil {
		if eta, ok := meta["driver_to_pickup_eta"]; ok {
			params["DriverToPickupETA"] = eta
		}
		if dist, ok := meta["driver_to_pickup_distance"]; ok {
			params["DriverToPickupDistance"] = dist
		}
	}

	// 创建消息对象
	message := &Message{
		Type:     msgType,
		Channels: []string{protocol.MsgChannelFcm},
		Params:   params,
		Language: getUserLanguage(passenger), // 使用用户偏好语言
	}

	// 使用消息服务发送
	return GetMessageService().SendMessage(message)
}

// NotifyDriver 通知司机 - 接收订单对象作为参数
func (s *OrderService) NotifyDriver(order *models.Order, notificationType string) error {
	if order == nil {
		return errors.New("order is nil")
	}

	if order.GetProviderID() == "" {
		return errors.New("driver not found in order")
	}

	// 获取司机信息
	driver := models.GetUserByID(order.GetProviderID())
	if driver == nil {
		return errors.New("driver not found")
	}

	// 获取订单详情
	orderDetail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
	if orderDetail == nil {
		return errors.New("order detail not found")
	}

	// 根据通知类型映射到消息类型
	var msgType string
	switch notificationType {
	case protocol.NotificationTypeTripEnded:
		msgType = protocol.MsgTypeDriverTripEnded
	case protocol.NotificationTypePaymentConfirmed:
		msgType = protocol.MsgTypeDriverPaymentConfirmed
	case protocol.NotificationTypeOrderCancelled:
		msgType = protocol.MsgTypeDriverOrderCancelled
	case protocol.NotificationTypeNewOrderAvailable:
		msgType = protocol.MsgTypeDriverNewOrder
	default:
		return fmt.Errorf("unsupported notification type for driver: %s", notificationType)
	}

	// 准备消息参数
	params := map[string]any{
		"to":                order.GetProviderID(),
		"OrderID":           order.OrderID,
		"OrderStatus":       order.GetStatus(),
		"PassengerName":     "Passenger", // 默认值
		"PickupAddress":     orderDetail.GetPickupAddress(),
		"DropoffAddress":    orderDetail.GetDropoffAddress(),
		"Amount":            order.GetPaymentAmount().StringFixed(2),
		"Currency":          order.GetCurrency(),
		"msg_type":          protocol.FCMMessageTypeOrder,
		"notification_type": notificationType,
		"order_status":      order.GetStatus(),
	}

	// 添加乘客信息
	if order.GetUserID() != "" {
		passenger := models.GetUserByID(order.GetUserID())
		if passenger != nil {
			params["PassengerName"] = passenger.GetFullName()
		}
	}

	// 添加取消原因信息
	if notificationType == protocol.NotificationTypeOrderCancelled && order.CancelReason != nil {
		params["CancelReason"] = *order.CancelReason
	}

	// 创建消息对象
	message := &Message{
		Type:     msgType,
		Channels: []string{protocol.MsgChannelFcm},
		Params:   params,
		Language: getUserLanguage(driver), // 使用用户偏好语言
	}

	// 使用消息服务发送
	return GetMessageService().SendMessage(message)
}

// NotifyOrderAccepted 通知乘客订单已被接单（司机触发）
// Includes driver-to-pickup ETA in the notification.
func (s *OrderService) NotifyOrderAccepted(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}

	// Calculate ETA from driver to pickup for the notification
	if order.GetProviderID() != "" {
		driver := models.GetUserByID(order.GetProviderID())
		detail := models.GetOrderDetail(order.OrderID, order.GetOrderType())
		if driver != nil && detail != nil && detail.GetPickupLatitude() != 0 {
			distKm := utils.CalculateDistanceHaversine(
				driver.GetLatitude(), driver.GetLongitude(),
				detail.GetPickupLatitude(), detail.GetPickupLongitude(),
			)
			etaMin := int(distKm * 2) // rough: 2 min/km
			if etaMin < 1 && distKm > 0 {
				etaMin = 1
			}
			// Store ETA in order metadata so NotifyPassenger can include it
			if order.GetMetadata() == nil {
				order.SetMetadata(map[string]any{})
			}
			meta := order.GetMetadata()
			meta["driver_to_pickup_eta"] = etaMin
			meta["driver_to_pickup_distance"] = fmt.Sprintf("%.1f", distKm)
			order.SetMetadata(meta)
		}
	}

	return s.NotifyPassenger(order, protocol.NotificationTypeOrderAccepted)
}

// NotifyDriverArrived 通知乘客司机已到达（司机触发）
func (s *OrderService) NotifyDriverArrived(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}
	return s.NotifyPassenger(order, protocol.NotificationTypeDriverArrived)
}

// NotifyTripStarted 通知乘客行程已开始（司机触发）
func (s *OrderService) NotifyTripStarted(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}
	return s.NotifyPassenger(order, protocol.NotificationTypeTripStarted)
}

// NotifyTripEnded 通知相关用户行程已结束（司机触发）
func (s *OrderService) NotifyTripEnded(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}

	// 通知乘客
	err1 := s.NotifyPassenger(order, protocol.NotificationTypeTripEnded)
	// 通知司机
	err2 := s.NotifyDriver(order, protocol.NotificationTypeTripEnded)

	// 返回第一个遇到的错误
	if err1 != nil {
		return err1
	}
	return err2
}

// NotifyPaymentConfirmed 通知司机和乘客支付已确认
func (s *OrderService) NotifyPaymentConfirmed(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}

	// 通知乘客支付确认
	err1 := s.NotifyPassenger(order, protocol.NotificationTypePaymentConfirmed)

	// 通知司机支付确认
	var err2 error
	if order.GetProviderID() != "" {
		err2 = s.NotifyDriver(order, protocol.NotificationTypePaymentConfirmed)
	}

	// 返回第一个遇到的错误
	if err1 != nil {
		return err1
	}
	return err2
}

// NotifyOrderCancelled 通知相关用户订单已取消（任一方触发）
func (s *OrderService) NotifyOrderCancelled(orderID string) error {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return errors.New("order not found")
	}

	// 通知乘客
	err1 := s.NotifyPassenger(order, protocol.NotificationTypeOrderCancelled)

	// 如果有司机，也通知司机
	var err2 error
	if order.GetProviderID() != "" {
		err2 = s.NotifyDriver(order, protocol.NotificationTypeOrderCancelled)
	}

	// 返回第一个遇到的错误
	if err1 != nil {
		return err1
	}
	return err2
}

// getTranslatedPaymentResult 根据语言获取翻译后的支付结果
// 检查ResCode是否是预定义的结果码，如果是，则使用翻译
func (s *OrderService) getTranslatedPaymentResult(resCode, resMsg, language string) string {
	// 获取有效的语言代码，如果不支持则使用默认语言
	language = i18n.GetValidLanguage(language)

	// 检查ResCode是否是预定义的结果码
	if s.isPredefinedResCode(resCode) {
		// 使用ResCode作为翻译键进行翻译
		translatedMsg := i18n.Translate(resCode, language)
		// 如果翻译成功（不等于原键值），则使用翻译结果
		if translatedMsg != "" && translatedMsg != resCode {
			return fmt.Sprintf("[%v]%v", resCode, translatedMsg)
		}
	}

	// 如果不是预定义的结果码或翻译失败，则使用原始消息
	return fmt.Sprintf("[%v]%v", resCode, resMsg)
}

// isPredefinedResCode 检查ResCode是否是预定义的结果码
func (s *OrderService) isPredefinedResCode(resCode string) bool {
	// 定义所有预定义的ResCode常量
	predefinedResCodes := []string{
		// 请求相关错误
		protocol.ResCodeRequestFailed,
		protocol.ResCodeRequestTimeout,
		protocol.ResCodeConnectionFailed,
		protocol.ResCodeNetworkError,

		// 响应解析相关错误
		protocol.ResCodeResponseParseFailed,
		protocol.ResCodeInvalidResponse,
		protocol.ResCodeMissingFields,
		protocol.ResCodeUnexpectedFormat,

		// 渠道相关错误
		protocol.ResCodeChannelError,
		protocol.ResCodeChannelUnavailable,
		protocol.ResCodeChannelMaintenance,
		protocol.ResCodeChannelRateLimited,

		// 配置相关错误
		protocol.ResCodeConfigError,
		protocol.ResCodeMissingConfig,
		protocol.ResCodeInvalidConfig,

		// 认证相关错误
		protocol.ResCodeAuthFailed,
		protocol.ResCodeInvalidCredentials,
		protocol.ResCodeTokenExpired,

		// 业务逻辑错误
		protocol.ResCodeBusinessError,
		protocol.ResCodeInvalidAmount,
		protocol.ResCodeInvalidCurrency,
		protocol.ResCodeInsufficientFunds,
		protocol.ResCodeUnsupportedPaymentMethod,

		// 支付成功结果码
		protocol.ResCodeSandboxSuccess,
		protocol.ResCodeCashSuccess,

		// 支付失败结果码
		protocol.ResCodePaymentFailed,

		// 系统错误
		protocol.ResCodeSystemError,
		protocol.ResCodeInternalError,
		protocol.ResCodeUnknownError,
	}

	// 检查resCode是否在预定义列表中
	return slices.Contains(predefinedResCodes, resCode)
}
