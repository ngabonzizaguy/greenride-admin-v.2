package protocol

import (
	"greenride/internal/utils"
)

// DriverRuntime 司机实时数据完整格式
type DriverRuntime struct {
	DriverID           string             `json:"driver_id"`
	OnlineStatus       string             `json:"online_status"` // online, offline, busy
	Latitude           float64            `json:"latitude"`
	Longitude          float64            `json:"longitude"`
	Heading            float64            `json:"heading"`             // 行驶方向
	Speed              float64            `json:"speed"`               // 速度(km/h)
	Accuracy           float64            `json:"accuracy"`            // 位置精度(米)
	LocationUpdatedAt  int64              `json:"location_updated_at"` // 位置更新时间
	CurrentOrder       *QueuedOrderData   `json:"current_order"`       // 使用轻量级订单结构
	QueuedOrders       []*QueuedOrderData `json:"queued_orders"`       // 排队订单列表，使用轻量级结构
	MaxQueueCapacity   int                `json:"max_queue_capacity"`  // 最大队列容量
	ConsecutiveRejects int                `json:"consecutive_rejects"` // 连续拒单次数
	LastDispatchAt     int64              `json:"last_dispatch_at"`    // 上次派单时间
	LastResponseAt     int64              `json:"last_response_at"`    // 上次响应时间
	AcceptanceRate     float64            `json:"acceptance_rate"`     // 实时接单率
	Rating             float64            `json:"rating"`              // 实时评分
	ExperienceLevel    int                `json:"experience_level"`    // 经验级别
	VehicleID          string             `json:"vehicle_id"`          // 绑定车辆ID
	LastHeartbeatAt    int64              `json:"last_heartbeat_at"`   // 最后心跳时间戳
	NextAvailableAt    int64              `json:"next_available_at"`   // 下次可用时间
	UpdatedAt          int64              `json:"updated_at"`          // 最后更新时间
	Version            int64              `json:"version"`             // 版本号(乐观锁)
}

// QueuedOrderData 司机排队订单轻量级数据
type QueuedOrderData struct {
	OrderID           string  `json:"order_id"`
	Status            string  `json:"status"`
	ScheduledAt       int64   `json:"scheduled_at,omitempty"`
	StartAt           int64   `json:"start_at,omitempty"`
	EndAt             int64   `json:"end_at,omitempty"`
	PickupLatitude    float64 `json:"pickup_latitude"`
	PickupLongitude   float64 `json:"pickup_longitude"`
	DropoffLatitude   float64 `json:"dropoff_latitude"`
	DropoffLongitude  float64 `json:"dropoff_longitude"`
	EstimatedDuration int     `json:"estimated_duration"` // 分钟
	PassengerCount    int     `json:"passenger_count"`
}

// IsOnline 检查司机是否在线
func (d *DriverRuntime) IsOnline() bool {
	return d.OnlineStatus == StatusOnline
}

// IsAvailable 检查司机是否可接单
func (d *DriverRuntime) IsAvailable() bool {
	if !d.IsOnline() {
		return false
	}

	// 检查心跳超时
	if utils.TimeNowMilli()-d.LastHeartbeatAt > 120 { // 2分钟
		//return false
	}

	// 检查数据时效性
	if utils.TimeNowMilli()-d.UpdatedAt > 300 { // 5分钟
		//return false
	}

	return true
}

// HasCurrentOrder 检查是否有当前订单
func (d *DriverRuntime) HasCurrentOrder() bool {
	return d.CurrentOrder != nil && d.CurrentOrder.OrderID != ""
}

// GetTotalQueueSize 获取总队列大小（包括当前订单）
func (d *DriverRuntime) GetTotalQueueSize() int {
	size := len(d.QueuedOrders)
	if d.HasCurrentOrder() {
		size++
	}
	return size
}

// CanAcceptMoreOrders 检查是否还能接更多订单
func (d *DriverRuntime) CanAcceptMoreOrders() bool {
	if d.MaxQueueCapacity <= 0 {
		return true // 无限制
	}
	return d.GetTotalQueueSize() < d.MaxQueueCapacity
}

// GetCacheKey 获取Redis存储键
func (d *DriverRuntime) GetCacheKey() string {
	return "driver:realtime:" + d.DriverID
}
