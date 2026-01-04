package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

// 预定义的原因描述（英文）
const (
	ReasonOrderCreated    = "Order created"
	ReasonOrderDispatched = "Order dispatched"
	ReasonOrderAssigned   = "Order assigned to driver"
	ReasonOrderAccepted   = "Driver accepted the order"
	ReasonOrderRejected   = "Driver rejected the order"
	ReasonDriverArrived   = "Driver arrived at pickup location"
	ReasonOrderStarted    = "Trip started"
	ReasonOrderCompleted  = "Trip completed"
	ReasonOrderCancelled  = "Order cancelled"
	ReasonPaymentSuccess  = "Payment successful"
	ReasonSystemProcess   = "System automated process"
)

// OrderHistoryService 订单历史记录服务
type OrderHistoryService struct{}

// RecordHistory 记录订单历史 - 单一入口函数
func (s *OrderHistoryService) RecordHistory(
	actionType string,
	order *models.Order,
	newOrder *models.Order,
	operatorID string,
	operatorType string,
	reason string,
) error {
	historyLog := models.NewOrderHistoryLog(actionType)

	// 填充订单信息
	if order != nil {
		historyLog.FillOrderInfo(order)
	}

	// 设置操作者信息
	historyLog.OperatorID = operatorID
	historyLog.OperatorType = operatorType
	historyLog.Reason = reason

	// 设置变更后信息
	if newOrder != nil {
		historyLog.FillValues(newOrder)
	} else if order != nil {
		historyLog.FillValues(order)
	}

	return models.CreateOrderHistoryLog(historyLog)
}

// 以下是简化版的订单历史记录功能

// RecordOrderCreated 记录订单创建历史
func (s *OrderHistoryService) RecordOrderCreated(order *models.Order, userID string) error {
	return s.RecordHistory(
		protocol.ActionOrderCreated,
		nil,                        // 没有原订单
		order,                      // 新创建的订单
		userID,                     // 用户ID
		protocol.UserTypePassenger, // 乘客创建订单
		ReasonOrderCreated,         // 原因
	)
}

// RecordOrderAccepted 记录订单接受历史
func (s *OrderHistoryService) RecordOrderAccepted(oldOrder *models.Order, newOrder *models.Order, driverID string) error {
	return s.RecordHistory(
		protocol.ActionOrderAccepted,
		oldOrder,                // 原订单
		newOrder,                // 更新后的订单
		driverID,                // 司机ID
		protocol.UserTypeDriver, // 司机操作
		ReasonOrderAccepted,     // 原因
	)
}

// RecordOrderRejected 记录订单拒绝历史
func (s *OrderHistoryService) RecordOrderRejected(oldOrder *models.Order, newOrder *models.Order, driverID string, reason string) error {
	reasonText := ReasonOrderRejected
	if reason != "" {
		reasonText = reason
	}
	return s.RecordHistory(
		protocol.ActionOrderRejected,
		oldOrder,                // 原订单
		newOrder,                // 更新后的订单
		driverID,                // 司机ID
		protocol.UserTypeDriver, // 司机操作
		reasonText,              // 原因
	)
}

// RecordOrderStarted 记录订单开始历史
func (s *OrderHistoryService) RecordOrderStarted(oldOrder *models.Order, newOrder *models.Order, driverID string) error {
	return s.RecordHistory(
		protocol.ActionOrderStarted,
		oldOrder,                // 原订单
		newOrder,                // 更新后的订单
		driverID,                // 司机ID
		protocol.UserTypeDriver, // 司机操作
		ReasonOrderStarted,      // 原因
	)
}

// RecordOrderFinished 记录订单结束历史
func (s *OrderHistoryService) RecordOrderFinished(oldOrder *models.Order, newOrder *models.Order, driverID string) error {
	return s.RecordHistory(
		protocol.ActionOrderCompleted,
		oldOrder,                // 原订单
		newOrder,                // 更新后的订单
		driverID,                // 司机ID
		protocol.UserTypeDriver, // 司机操作
		ReasonOrderCompleted,    // 原因
	)
}

// RecordOrderCancelled 记录订单取消历史
func (s *OrderHistoryService) RecordOrderCancelled(oldOrder *models.Order, newOrder *models.Order, cancelledBy string, cancellerType string, reason string) error {
	// 根据操作者类型选择正确的取消动作类型
	actionType := protocol.ActionOrderCancelled
	switch cancellerType {
	case protocol.UserTypePassenger:
		actionType = protocol.ActionOrderCancelledByUser
	case protocol.UserTypeDriver:
		actionType = protocol.ActionOrderCancelledByProvider
	case "system":
		actionType = protocol.ActionOrderCancelledBySystem
	}

	return s.RecordHistory(
		actionType,
		oldOrder,      // 原订单
		newOrder,      // 更新后的订单
		cancelledBy,   // 取消者ID
		cancellerType, // 取消者类型
		reason,        // 取消原因
	)
}

// RecordOrderPaymentCompleted 记录支付完成历史
func (s *OrderHistoryService) RecordOrderPaymentCompleted(oldOrder *models.Order, newOrder *models.Order, userID string) error {
	return s.RecordHistory(
		protocol.ActionPaymentCompleted,
		oldOrder,                   // 原订单
		newOrder,                   // 更新后的订单
		userID,                     // 用户ID
		protocol.UserTypePassenger, // 乘客支付
		ReasonPaymentSuccess,       // 原因
	)
}

// GetOrderHistoryService 获取订单历史服务实例
func GetOrderHistoryService() *OrderHistoryService {
	return &OrderHistoryService{}
}
