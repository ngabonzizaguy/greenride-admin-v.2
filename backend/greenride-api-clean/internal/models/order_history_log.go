package models

import (
	"encoding/json"

	"greenride/internal/log"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// OrderHistoryLog 订单历史记录模型
type OrderHistoryLog struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	HistoryID  string `json:"history_id" gorm:"column:history_id;type:varchar(64);uniqueIndex"`
	ActionType string `json:"action_type" gorm:"column:action_type;type:varchar(32);index"` // 使用protocol包中定义的常量
	OrderID    string `json:"order_id" gorm:"column:order_id;type:varchar(64);index"`
	OrderType  string `json:"order_type" gorm:"column:order_type;type:varchar(32);index"` // 与订单类型对齐

	// 用户信息 - 乘客
	UserID   string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType string `json:"user_type" gorm:"column:user_type;type:varchar(32)"`

	// 服务提供者信息 - 司机
	ProviderID   string `json:"provider_id" gorm:"column:provider_id;type:varchar(64);index"`
	ProviderType string `json:"provider_type" gorm:"column:provider_type;type:varchar(32)"`

	// 操作者信息
	OperatorID   string `json:"operator_id" gorm:"column:operator_id;type:varchar(64);index"`
	OperatorType string `json:"operator_type" gorm:"column:operator_type;type:varchar(32)"`

	// 状态信息
	FromStatus string `json:"from_status" gorm:"column:from_status;type:varchar(32)"`
	ToStatus   string `json:"to_status" gorm:"column:to_status;type:varchar(32);index"`

	// 详细变更记录
	Before string `json:"before" gorm:"column:before;type:json"`
	After  string `json:"after" gorm:"column:after;type:json"`

	// 其他信息
	Reason    string `json:"reason" gorm:"column:reason;type:varchar(255)"`
	IPAddress string `json:"ip_address" gorm:"column:ip_address;type:varchar(45)"`
	UserAgent string `json:"user_agent" gorm:"column:user_agent;type:varchar(255)"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

func (OrderHistoryLog) TableName() string {
	return "t_order_history_logs"
}

type OrderHistoryLogs []*OrderHistoryLog

// ToInfo 将历史记录转换为API响应格式
func (t *OrderHistoryLog) ToInfo() map[string]interface{} {
	return map[string]interface{}{
		"history_id":    t.HistoryID,
		"action_type":   t.ActionType,
		"order_id":      t.OrderID,
		"order_type":    t.OrderType,
		"user_id":       t.UserID,
		"user_type":     t.UserType,
		"provider_id":   t.ProviderID,
		"provider_type": t.ProviderType,
		"operator_id":   t.OperatorID,
		"operator_type": t.OperatorType,
		"from_status":   t.FromStatus,
		"to_status":     t.ToStatus,
		"reason":        t.Reason,
		"created_at":    t.CreatedAt,
	}
}

// ToInfos 批量转换历史记录
func (t OrderHistoryLogs) ToInfos() []map[string]interface{} {
	infos := make([]map[string]interface{}, len(t))
	for i, history := range t {
		infos[i] = history.ToInfo()
	}
	return infos
}

// NewOrderHistoryLog 创建新的订单历史记录
func NewOrderHistoryLog(actionType string) *OrderHistoryLog {
	return &OrderHistoryLog{
		HistoryID:  utils.GenerateUUID(),
		ActionType: actionType,
		OrderType:  "",
		Before:     "{}",
		After:      "{}",
	}
}

// FillOrderInfo 填充订单信息
func (t *OrderHistoryLog) FillOrderInfo(order *Order) {
	if order == nil {
		return
	}
	t.OrderID = order.OrderID
	t.OrderType = order.GetOrderType()
	t.UserID = order.GetUserID()
	t.UserType = protocol.UserTypePassenger
	t.ProviderID = order.GetProviderID()
	t.ProviderType = protocol.UserTypeDriver
	t.FromStatus = order.GetStatus()

	// 序列化整个订单对象作为变更前的快照
	beforeJSON, err := json.Marshal(order)
	if err == nil {
		t.Before = string(beforeJSON)
	} else {
		t.Before = "{}"
		log.Errorf("FillOrderInfo marshal error: %v", err)
	}
}

// FillValues 填充变更后的值
func (t *OrderHistoryLog) FillValues(order *Order) {
	if order == nil {
		return
	}
	t.ToStatus = order.GetStatus()

	// 序列化变更后的订单对象作为快照
	afterJSON, err := json.Marshal(order)
	if err == nil {
		t.After = string(afterJSON)
	} else {
		t.After = "{}"
		log.Errorf("FillValues marshal error: %v", err)
	}
}

// CreateOrderHistoryLog 创建订单历史记录
func CreateOrderHistoryLog(history *OrderHistoryLog) error {
	err := GetDB().Create(history).Error
	if err != nil {
		log.Errorf("CreateOrderHistoryLog err: %v", err)
	}
	return err
}

// GetOrderHistoryTrail 获取订单的完整历史记录
func GetOrderHistoryTrail(orderID string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("order_id = ?", orderID).Order("created_at ASC").Find(&logs).Error
	return logs, err
}

// GetUserOrderHistory 获取用户相关的订单历史记录
func GetUserOrderHistory(userID string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetProviderOrderHistory 获取服务提供者相关的订单历史记录
func GetProviderOrderHistory(providerID string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("provider_id = ?", providerID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetOperatorOrderHistory 获取操作者的订单历史记录
func GetOperatorOrderHistory(operatorID string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("operator_id = ?", operatorID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetOrderHistoryByActionType 根据操作类型获取历史记录
func GetOrderHistoryByActionType(actionType string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("action_type = ?", actionType).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetOrderHistoryByOrderType 根据订单类型获取历史记录
func GetOrderHistoryByOrderType(orderType string) (OrderHistoryLogs, error) {
	var logs OrderHistoryLogs
	err := GetDB().Where("order_type = ?", orderType).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// CreateOrderHistory 创建订单历史记录的辅助函数
func CreateOrderHistory(actionType string, order *Order, newOrder *Order, operatorID string, operatorType string, reason string, ipAddress, userAgent string) error {
	// 使用protocol包中定义的ActionType常量
	historyLog := NewOrderHistoryLog(actionType)

	// 设置基本信息
	if order != nil {
		historyLog.FillOrderInfo(order)
	}

	// 设置操作者信息
	historyLog.OperatorID = operatorID
	historyLog.OperatorType = operatorType

	// 设置其他信息
	historyLog.Reason = reason
	historyLog.IPAddress = ipAddress
	historyLog.UserAgent = userAgent

	// 设置新状态
	if newOrder != nil {
		historyLog.FillValues(newOrder)
	} else if order != nil {
		// 如果没有提供新订单对象，则使用原订单对象填充After字段
		historyLog.FillValues(order)
	}

	return CreateOrderHistoryLog(historyLog)
}

// 一些便捷方法，用于订单流程的各个节点记录

// RecordOrderCreated 记录订单创建
func RecordOrderCreated(order *Order, operatorID string, operatorType string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderCreated, nil, order, operatorID, operatorType, "订单创建", ipAddress, userAgent)
}

// RecordOrderAssigned 记录订单分配
func RecordOrderAssigned(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, reason string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderAssigned, oldOrder, newOrder, operatorID, operatorType, reason, ipAddress, userAgent)
}

// RecordOrderAccepted 记录订单接受
func RecordOrderAccepted(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderAccepted, oldOrder, newOrder, operatorID, operatorType, "订单接受", ipAddress, userAgent)
}

// RecordOrderStarted 记录订单开始
func RecordOrderStarted(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderStarted, oldOrder, newOrder, operatorID, operatorType, "订单开始", ipAddress, userAgent)
}

// RecordOrderCompleted 记录订单完成
func RecordOrderCompleted(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderCompleted, oldOrder, newOrder, operatorID, operatorType, "订单完成", ipAddress, userAgent)
}

// RecordOrderCancelled 记录订单取消
func RecordOrderCancelled(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, reason string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderCancelled, oldOrder, newOrder, operatorID, operatorType, reason, ipAddress, userAgent)
}

// RecordOrderRejected 记录订单拒绝
func RecordOrderRejected(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, reason string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderRejected, oldOrder, newOrder, operatorID, operatorType, reason, ipAddress, userAgent)
}

// RecordOrderRefunded 记录订单退款
func RecordOrderRefunded(oldOrder *Order, newOrder *Order, operatorID string, operatorType string, reason string, ipAddress, userAgent string) error {
	return CreateOrderHistory(protocol.ActionOrderRefunded, oldOrder, newOrder, operatorID, operatorType, reason, ipAddress, userAgent)
}
