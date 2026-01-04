package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// RideOrder 行程订单表 - 基于最新设计文档
type RideOrder struct {
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	OrderID string `json:"order_id" gorm:"column:order_id;type:varchar(64);uniqueIndex"` // 关联OrderV2的order_id，使用唯一索引，不设置外键约束
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*RideOrderValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type RideOrderValues struct {
	// 网约车特有信息
	VehicleID *string `json:"vehicle_id" gorm:"column:vehicle_id;type:varchar(64);index"` // 车辆ID，关联车辆信息表

	VehicleCategory *string `json:"vehicle_category" gorm:"column:vehicle_category;type:varchar(64);index"` // 车辆类别，如经济型、舒适型等
	VehicleLevel    *string `json:"vehicle_level" gorm:"column:vehicle_level;type:varchar(64);index"`       // 车辆等级，如普通、豪华等

	// 乘客信息
	PassengerCount *int `json:"passenger_count" gorm:"column:passenger_count;default:1"` // 乘客数量，默认1人

	// 上车地点信息
	PickupAddress   *string  `json:"pickup_address" gorm:"column:pickup_address;type:text"`              // 上车地址详情
	PickupLatitude  *float64 `json:"pickup_latitude" gorm:"column:pickup_latitude;type:decimal(10,8)"`   // 上车地点纬度
	PickupLongitude *float64 `json:"pickup_longitude" gorm:"column:pickup_longitude;type:decimal(11,8)"` // 上车地点经度
	PickupLandmark  *string  `json:"pickup_landmark" gorm:"column:pickup_landmark;type:varchar(255)"`    // 上车地点地标信息

	// 下车地点信息
	DropoffAddress   *string  `json:"dropoff_address" gorm:"column:dropoff_address;type:text"`              // 下车地址详情
	DropoffLatitude  *float64 `json:"dropoff_latitude" gorm:"column:dropoff_latitude;type:decimal(10,8)"`   // 下车地点纬度
	DropoffLongitude *float64 `json:"dropoff_longitude" gorm:"column:dropoff_longitude;type:decimal(11,8)"` // 下车地点经度
	DropoffLandmark  *string  `json:"dropoff_landmark" gorm:"column:dropoff_landmark;type:varchar(255)"`    // 下车地点地标信息

	// 行程距离和时间统计
	EstimatedDistance *float64 `json:"estimated_distance" gorm:"column:estimated_distance;type:decimal(8,2)"` // 预估行程距离（公里）
	EstimatedDuration *int     `json:"estimated_duration" gorm:"column:estimated_duration"`                   // 预估行程时长（分钟）
	ActualDistance    *float64 `json:"actual_distance" gorm:"column:actual_distance;type:decimal(8,2)"`       // 实际行程距离（公里）
	ActualDuration    *int     `json:"actual_duration" gorm:"column:actual_duration"`                         // 实际行程时长（分钟）

	// 网约车费用明细
	BaseFare     *float64 `json:"base_fare" gorm:"column:base_fare;type:decimal(10,2)"`         // 起步价
	DistanceFare *float64 `json:"distance_fare" gorm:"column:distance_fare;type:decimal(10,2)"` // 距离费用
	TimeFare     *float64 `json:"time_fare" gorm:"column:time_fare;type:decimal(10,2)"`         // 时长费用
	SurgeFare    *float64 `json:"surge_fare" gorm:"column:surge_fare;type:decimal(10,2)"`       // 动态调价费用
	TotalFare    *float64 `json:"total_fare" gorm:"column:total_fare;type:decimal(10,2)"`       // 行程总费用（对应OrderV2.TotalAmount）

	// 网约车特有时间节点
	DriverEnRouteAt *int64 `json:"driver_en_route_at" gorm:"column:driver_en_route_at"` // 司机出发前往上车点时间戳
	ArrivedAt       *int64 `json:"arrived_at" gorm:"column:arrived_at"`                 // 司机到达上车点时间戳

	// 行程路径数据
	RouteData *string `json:"route_data" gorm:"column:route_data;type:json"` // 行程路径轨迹数据（JSON格式）

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 记录更新时间戳
}

func (RideOrder) TableName() string {
	return "t_ride_orders"
}

// 创建新的乘车订单对象
func NewRideOrder(orderID string) *RideOrder {
	return &RideOrder{
		OrderID:         orderID,
		Salt:            utils.GenerateSalt(),
		RideOrderValues: &RideOrderValues{},
	}
}

// SetValues 更新RideOrderV2Values中的非nil值
func (r *RideOrderValues) SetValues(values *RideOrderValues) {
	if values == nil {
		return
	}
	if values.VehicleID != nil {
		r.VehicleID = values.VehicleID
	}
	if values.VehicleCategory != nil {
		r.VehicleCategory = values.VehicleCategory
	}
	if values.VehicleLevel != nil {
		r.VehicleLevel = values.VehicleLevel
	}
	if values.PassengerCount != nil {
		r.PassengerCount = values.PassengerCount
	}
	if values.PickupAddress != nil {
		r.PickupAddress = values.PickupAddress
	}
	if values.PickupLatitude != nil {
		r.PickupLatitude = values.PickupLatitude
	}
	if values.PickupLongitude != nil {
		r.PickupLongitude = values.PickupLongitude
	}
	if values.DropoffAddress != nil {
		r.DropoffAddress = values.DropoffAddress
	}
	if values.DropoffLatitude != nil {
		r.DropoffLatitude = values.DropoffLatitude
	}
	if values.DropoffLongitude != nil {
		r.DropoffLongitude = values.DropoffLongitude
	}
	if values.TotalFare != nil {
		r.TotalFare = values.TotalFare
	}
	if values.RouteData != nil {
		r.RouteData = values.RouteData
	}
	if values.UpdatedAt > 0 {
		r.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法

func (p *RideOrderValues) GetVehicleCategory() string {
	if p.VehicleCategory == nil {
		return ""
	}
	return *p.VehicleCategory
}
func (p *RideOrderValues) GetVehicleLevel() string {
	if p.VehicleLevel == nil {
		return ""
	}
	return *p.VehicleLevel
}
func (p *RideOrderValues) SetVehicleCategory(category string) *RideOrderValues {
	p.VehicleCategory = &category
	return p
}
func (p *RideOrderValues) SetVehicleLevel(level string) *RideOrderValues {
	p.VehicleLevel = &level
	return p
}

func (r *RideOrderValues) GetPassengerCount() int {
	if r.PassengerCount == nil {
		return 1
	}
	return *r.PassengerCount
}

func (r *RideOrderValues) GetTotalFare() float64 {
	if r.TotalFare == nil {
		return 0
	}
	return *r.TotalFare
}

func (r *RideOrderValues) GetPickupLatitude() float64 {
	if r.PickupLatitude == nil {
		return 0
	}
	return *r.PickupLatitude
}

func (r *RideOrderValues) GetPickupLongitude() float64 {
	if r.PickupLongitude == nil {
		return 0
	}
	return *r.PickupLongitude
}

func (r *RideOrderValues) GetDropoffLatitude() float64 {
	if r.DropoffLatitude == nil {
		return 0
	}
	return *r.DropoffLatitude
}

func (r *RideOrderValues) GetDropoffLongitude() float64 {
	if r.DropoffLongitude == nil {
		return 0
	}
	return *r.DropoffLongitude
}

func (r *RideOrderValues) GetEstimatedDuration() int {
	if r.EstimatedDuration == nil {
		return 0
	}
	return *r.EstimatedDuration
}

func (r *RideOrderValues) GetEstimatedDistance() float64 {
	if r.EstimatedDistance == nil {
		return 0
	}
	return *r.EstimatedDistance
}

func (r *RideOrderValues) GetPickupAddress() string {
	if r.PickupAddress == nil {
		return ""
	}
	return *r.PickupAddress
}

func (r *RideOrderValues) GetDropoffAddress() string {
	if r.DropoffAddress == nil {
		return ""
	}
	return *r.DropoffAddress
}

// Setter 方法
func (r *RideOrderValues) SetPickupLocation(address string, lat, lng float64) *RideOrderValues {
	r.PickupAddress = &address
	r.PickupLatitude = &lat
	r.PickupLongitude = &lng
	return r
}

func (r *RideOrderValues) SetDropoffLocation(address string, lat, lng float64) *RideOrderValues {
	r.DropoffAddress = &address
	r.DropoffLatitude = &lat
	r.DropoffLongitude = &lng
	return r
}

func (r *RideOrderValues) SetVehicle(vehicleID string) *RideOrderValues {
	r.VehicleID = &vehicleID
	return r
}

// 添加所有字段的单独setter方法，支持链式调用
func (r *RideOrderValues) SetPassengerCount(count int) *RideOrderValues {
	r.PassengerCount = &count
	return r
}

func (r *RideOrderValues) SetPickupAddress(address string) *RideOrderValues {
	r.PickupAddress = &address
	return r
}

func (r *RideOrderValues) SetPickupLatitude(lat float64) *RideOrderValues {
	r.PickupLatitude = &lat
	return r
}

func (r *RideOrderValues) SetPickupLongitude(lng float64) *RideOrderValues {
	r.PickupLongitude = &lng
	return r
}

func (r *RideOrderValues) SetPickupLandmark(landmark string) *RideOrderValues {
	r.PickupLandmark = &landmark
	return r
}

func (r *RideOrderValues) SetDropoffAddress(address string) *RideOrderValues {
	r.DropoffAddress = &address
	return r
}

func (r *RideOrderValues) SetDropoffLatitude(lat float64) *RideOrderValues {
	r.DropoffLatitude = &lat
	return r
}

func (r *RideOrderValues) SetDropoffLongitude(lng float64) *RideOrderValues {
	r.DropoffLongitude = &lng
	return r
}

func (r *RideOrderValues) SetDropoffLandmark(landmark string) *RideOrderValues {
	r.DropoffLandmark = &landmark
	return r
}

func (r *RideOrderValues) SetEstimatedDistance(distance float64) *RideOrderValues {
	r.EstimatedDistance = &distance
	return r
}

func (r *RideOrderValues) SetEstimatedDuration(duration int) *RideOrderValues {
	r.EstimatedDuration = &duration
	return r
}

// 网约车特有业务方法
func (r *RideOrderValues) SetDriverEnRoute() {
	now := utils.TimeNowMilli()
	r.DriverEnRouteAt = &now
}

func (r *RideOrderValues) SetDriverArrived() *RideOrderValues {
	now := utils.TimeNowMilli()
	r.ArrivedAt = &now
	return r
}

// ToOrderDetail 将RideOrderV2转换为统一的OrderDetail
func (r *RideOrder) ToOrderDetail() *protocol.OrderDetail {
	detail := &protocol.OrderDetail{
		OrderID:   r.OrderID,
		OrderType: protocol.RideOrder,
	}

	// 填充网约车特有字段
	if r.PassengerCount != nil {
		detail.PassengerCount = *r.PassengerCount
	}
	if r.PickupAddress != nil {
		detail.PickupAddress = *r.PickupAddress
	}
	if r.PickupLatitude != nil {
		detail.PickupLatitude = *r.PickupLatitude
	}
	if r.PickupLongitude != nil {
		detail.PickupLongitude = *r.PickupLongitude
	}
	if r.PickupLandmark != nil {
		detail.PickupLandmark = *r.PickupLandmark
	}
	if r.DropoffAddress != nil {
		detail.DropoffAddress = *r.DropoffAddress
	}
	if r.DropoffLatitude != nil {
		detail.DropoffLatitude = *r.DropoffLatitude
	}
	if r.DropoffLongitude != nil {
		detail.DropoffLongitude = *r.DropoffLongitude
	}
	if r.DropoffLandmark != nil {
		detail.DropoffLandmark = *r.DropoffLandmark
	}
	if r.EstimatedDistance != nil {
		detail.EstimatedDistance = *r.EstimatedDistance
	}
	if r.EstimatedDuration != nil {
		detail.EstimatedDuration = *r.EstimatedDuration
	}
	if r.ActualDistance != nil {
		detail.ActualDistance = *r.ActualDistance
	}
	if r.ActualDuration != nil {
		detail.ActualDuration = *r.ActualDuration
	}
	if r.VehicleID != nil {
		detail.VehicleID = *r.VehicleID
	}

	return detail
}

func GetRideOrderByOrderID(orderID string) *RideOrder {
	var rideOrder RideOrder
	err := DB.Where("order_id = ?", orderID).First(&rideOrder).Error
	if err != nil {
		return nil
	}
	return &rideOrder
}
