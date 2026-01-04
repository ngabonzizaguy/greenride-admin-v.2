package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// OrderDetail 统一订单详情实体，兼容所有业务子订单字段
type OrderDetail struct {
	// 基础字段
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	OrderID string `json:"order_id" gorm:"column:order_id;type:varchar(64)"`
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*OrderDetailValues
	// 时间戳
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type OrderDetailValues struct {
	// 网约车详情字段（当order_type=ride时使用）
	VehicleID         *string  `json:"vehicle_id,omitempty" gorm:"column:vehicle_id;type:varchar(64)"`
	VehicleCategory   *string  `json:"vehicle_category,omitempty" gorm:"column:vehicle_category;type:varchar(64)"`
	VehicleLevel      *string  `json:"vehicle_level,omitempty" gorm:"column:vehicle_level;type:varchar(64)"`
	PassengerCount    *int     `json:"passenger_count,omitempty" gorm:"column:passenger_count"`
	PickupAddress     *string  `json:"pickup_address,omitempty" gorm:"column:pickup_address;type:text"`
	PickupLatitude    *float64 `json:"pickup_latitude,omitempty" gorm:"column:pickup_latitude;type:decimal(10,8)"`
	PickupLongitude   *float64 `json:"pickup_longitude,omitempty" gorm:"column:pickup_longitude;type:decimal(11,8)"`
	PickupLandmark    *string  `json:"pickup_landmark,omitempty" gorm:"column:pickup_landmark;type:varchar(255)"`
	DropoffAddress    *string  `json:"dropoff_address,omitempty" gorm:"column:dropoff_address;type:text"`
	DropoffLatitude   *float64 `json:"dropoff_latitude,omitempty" gorm:"column:dropoff_latitude;type:decimal(10,8)"`
	DropoffLongitude  *float64 `json:"dropoff_longitude,omitempty" gorm:"column:dropoff_longitude;type:decimal(11,8)"`
	DropoffLandmark   *string  `json:"dropoff_landmark,omitempty" gorm:"column:dropoff_landmark;type:varchar(255)"`
	EstimatedDistance *float64 `json:"estimated_distance,omitempty" gorm:"column:estimated_distance;type:decimal(8,2)"`
	EstimatedDuration *int     `json:"estimated_duration,omitempty" gorm:"column:estimated_duration"`
	ActualDistance    *float64 `json:"actual_distance,omitempty" gorm:"column:actual_distance;type:decimal(8,2)"`
	ActualDuration    *int     `json:"actual_duration,omitempty" gorm:"column:actual_duration"`
	BaseFare          *float64 `json:"base_fare,omitempty" gorm:"column:base_fare;type:decimal(10,2)"`
	DistanceFare      *float64 `json:"distance_fare,omitempty" gorm:"column:distance_fare;type:decimal(10,2)"`
	TimeFare          *float64 `json:"time_fare,omitempty" gorm:"column:time_fare;type:decimal(10,2)"`
	SurgeFare         *float64 `json:"surge_fare,omitempty" gorm:"column:surge_fare;type:decimal(10,2)"`
	TotalFare         *float64 `json:"total_fare,omitempty" gorm:"column:total_fare;type:decimal(10,2)"`
	DriverEnRouteAt   *int64   `json:"driver_en_route_at,omitempty" gorm:"column:driver_en_route_at"`
	ArrivedAt         *int64   `json:"arrived_at,omitempty" gorm:"column:arrived_at"`
	RouteData         *string  `json:"route_data,omitempty" gorm:"column:route_data;type:json"`

	// 外卖详情字段（当order_type=delivery时使用，预留）
	RestaurantID      *string  `json:"restaurant_id,omitempty" gorm:"column:restaurant_id;type:varchar(64)"`
	RestaurantName    *string  `json:"restaurant_name,omitempty" gorm:"column:restaurant_name;type:varchar(255)"`
	RestaurantAddress *string  `json:"restaurant_address,omitempty" gorm:"column:restaurant_address;type:text"`
	RestaurantPhone   *string  `json:"restaurant_phone,omitempty" gorm:"column:restaurant_phone;type:varchar(50)"`
	DeliveryAddress   *string  `json:"delivery_address,omitempty" gorm:"column:delivery_address;type:text"`
	DeliveryLatitude  *float64 `json:"delivery_latitude,omitempty" gorm:"column:delivery_latitude;type:decimal(10,8)"`
	DeliveryLongitude *float64 `json:"delivery_longitude,omitempty" gorm:"column:delivery_longitude;type:decimal(11,8)"`
	CourierID         *string  `json:"courier_id,omitempty" gorm:"column:courier_id;type:varchar(64)"`
	CourierName       *string  `json:"courier_name,omitempty" gorm:"column:courier_name;type:varchar(255)"`
	CourierPhone      *string  `json:"courier_phone,omitempty" gorm:"column:courier_phone;type:varchar(50)"`
	CourierRating     *float64 `json:"courier_rating,omitempty" gorm:"column:courier_rating;type:decimal(3,2)"`

	// 购物详情字段（当order_type=shopping时使用，预留）
	StoreID         *string `json:"store_id,omitempty" gorm:"column:store_id;type:varchar(64)"`
	StoreName       *string `json:"store_name,omitempty" gorm:"column:store_name;type:varchar(255)"`
	StoreAddress    *string `json:"store_address,omitempty" gorm:"column:store_address;type:text"`
	ShippingAddress *string `json:"shipping_address,omitempty" gorm:"column:shipping_address;type:text"`
	TrackingNumber  *string `json:"tracking_number,omitempty" gorm:"column:tracking_number;type:varchar(100)"`
}

// TableName 根据订单类型返回对应的表名
func (o *OrderDetail) TableName() string {
	// 这个方法不会被直接调用，因为我们会在查询时动态指定表名
	return ""
}

// GetTableNameByOrderType 根据订单类型返回对应的表名
func GetTableNameByOrderType(orderType string) string {
	switch orderType {
	case "ride":
		return "t_ride_orders"
	default:
		return "" // 默认表名
	}
}

func GetOrderDetail(orderID, orderType string) *OrderDetail {
	var orderDetail *OrderDetail
	err := DB.Table(GetTableNameByOrderType(orderType)).Where("order_id = ?", orderID).First(&orderDetail).Error
	if err != nil {
		return nil
	}
	return orderDetail
}

// Protocol 将OrderDetail转换为protocol.OrderDetail
func (o *OrderDetail) Protocol() *protocol.OrderDetail {
	detail := &protocol.OrderDetail{}

	if o.VehicleCategory != nil {
		detail.VehicleCategory = *o.VehicleCategory
	}
	if o.VehicleLevel != nil {
		detail.VehicleLevel = *o.VehicleLevel
	}
	// 网约车详情字段
	if o.PassengerCount != nil {
		detail.PassengerCount = *o.PassengerCount
	}
	if o.PickupAddress != nil {
		detail.PickupAddress = *o.PickupAddress
	}
	if o.PickupLatitude != nil {
		detail.PickupLatitude = *o.PickupLatitude
	}
	if o.PickupLongitude != nil {
		detail.PickupLongitude = *o.PickupLongitude
	}
	if o.PickupLandmark != nil {
		detail.PickupLandmark = *o.PickupLandmark
	}
	if o.DropoffAddress != nil {
		detail.DropoffAddress = *o.DropoffAddress
	}
	if o.DropoffLatitude != nil {
		detail.DropoffLatitude = *o.DropoffLatitude
	}
	if o.DropoffLongitude != nil {
		detail.DropoffLongitude = *o.DropoffLongitude
	}
	if o.DropoffLandmark != nil {
		detail.DropoffLandmark = *o.DropoffLandmark
	}
	if o.EstimatedDistance != nil {
		detail.EstimatedDistance = *o.EstimatedDistance
	}
	if o.EstimatedDuration != nil {
		detail.EstimatedDuration = *o.EstimatedDuration
	}
	if o.ActualDistance != nil {
		detail.ActualDistance = *o.ActualDistance
	}
	if o.ActualDuration != nil {
		detail.ActualDuration = *o.ActualDuration
	}
	if o.VehicleID != nil {
		detail.VehicleID = *o.VehicleID
	}

	// 外卖详情字段
	if o.RestaurantID != nil {
		detail.RestaurantID = *o.RestaurantID
	}
	if o.RestaurantName != nil {
		detail.RestaurantName = *o.RestaurantName
	}
	if o.RestaurantAddress != nil {
		detail.RestaurantAddress = *o.RestaurantAddress
	}
	if o.RestaurantPhone != nil {
		detail.RestaurantPhone = *o.RestaurantPhone
	}
	if o.DeliveryAddress != nil {
		detail.DeliveryAddress = *o.DeliveryAddress
	}
	if o.DeliveryLatitude != nil {
		detail.DeliveryLatitude = *o.DeliveryLatitude
	}
	if o.DeliveryLongitude != nil {
		detail.DeliveryLongitude = *o.DeliveryLongitude
	}
	if o.CourierID != nil {
		detail.CourierID = *o.CourierID
	}
	if o.CourierName != nil {
		detail.CourierName = *o.CourierName
	}
	if o.CourierPhone != nil {
		detail.CourierPhone = *o.CourierPhone
	}
	if o.CourierRating != nil {
		detail.CourierRating = *o.CourierRating
	}

	// 购物详情字段
	if o.StoreID != nil {
		detail.StoreID = *o.StoreID
	}
	if o.StoreName != nil {
		detail.StoreName = *o.StoreName
	}
	if o.StoreAddress != nil {
		detail.StoreAddress = *o.StoreAddress
	}
	if o.ShippingAddress != nil {
		detail.ShippingAddress = *o.ShippingAddress
	}
	if o.TrackingNumber != nil {
		detail.TrackingNumber = *o.TrackingNumber
	}

	return detail
}

func (o *OrderDetail) GetVehicleCategory() string {
	if o.VehicleCategory == nil {
		return ""
	}
	return *o.VehicleCategory
}
func (o *OrderDetail) GetVehicleLevel() string {
	if o.VehicleLevel == nil {
		return ""
	}
	return *o.VehicleLevel
}
func (o *OrderDetail) SetVehicleCategory(category string) *OrderDetail {
	o.VehicleCategory = &category
	return o
}
func (o *OrderDetail) SetVehicleLevel(level string) *OrderDetail {
	o.VehicleLevel = &level
	return o
}

// GetPassengerCount 获取乘客数量（网约车）
func (o *OrderDetail) GetPassengerCount() int {
	if o.PassengerCount == nil {
		return 1
	}
	return *o.PassengerCount
}

// GetTotalFare 获取总费用（网约车）
func (o *OrderDetail) GetTotalFare() float64 {
	if o.TotalFare == nil {
		return 0
	}
	return *o.TotalFare
}

// SetPickupLocation 设置上车地点（网约车）
func (o *OrderDetail) SetPickupLocation(address string, lat, lng float64, landmark string) {
	o.PickupAddress = &address
	o.PickupLatitude = &lat
	o.PickupLongitude = &lng
	if landmark != "" {
		o.PickupLandmark = &landmark
	}
}

// SetDropoffLocation 设置下车地点（网约车）
func (o *OrderDetail) SetDropoffLocation(address string, lat, lng float64, landmark string) {
	o.DropoffAddress = &address
	o.DropoffLatitude = &lat
	o.DropoffLongitude = &lng
	if landmark != "" {
		o.DropoffLandmark = &landmark
	}
}

// SetDriverEnRoute 设置司机出发时间（网约车）
func (o *OrderDetail) SetDriverEnRoute() {
	now := utils.TimeNowMilli()
	o.DriverEnRouteAt = &now
}

// SetDriverArrived 设置司机到达时间（网约车）
func (o *OrderDetail) SetDriverArrived() *OrderDetail {
	now := utils.TimeNowMilli()
	o.ArrivedAt = &now
	return o
}

// CalculateTotalFare 计算总费用（网约车）
func (o *OrderDetail) CalculateTotalFare() {
	total := float64(0)

	if o.BaseFare != nil {
		total += *o.BaseFare
	}
	if o.DistanceFare != nil {
		total += *o.DistanceFare
	}
	if o.TimeFare != nil {
		total += *o.TimeFare
	}
	if o.SurgeFare != nil {
		total += *o.SurgeFare
	}

	o.TotalFare = &total
}

// ================== 网约车字段 Getter/Setter ==================

// GetVehicleID 获取车辆ID
func (o *OrderDetailValues) GetVehicleID() string {
	if o.VehicleID == nil {
		return ""
	}
	return *o.VehicleID
}

// SetVehicleID 设置车辆ID
func (o *OrderDetailValues) SetVehicleID(vehicleID string) *OrderDetailValues {
	o.VehicleID = &vehicleID
	return o
}

// GetPickupAddress 获取上车地址
func (o *OrderDetailValues) GetPickupAddress() string {
	if o.PickupAddress == nil {
		return ""
	}
	return *o.PickupAddress
}

// SetPickupAddress 设置上车地址
func (o *OrderDetailValues) SetPickupAddress(address string) *OrderDetailValues {
	o.PickupAddress = &address
	return o
}

// GetPickupLatitude 获取上车纬度
func (o *OrderDetailValues) GetPickupLatitude() float64 {
	if o.PickupLatitude == nil {
		return 0
	}
	return *o.PickupLatitude
}

// SetPickupLatitude 设置上车纬度
func (o *OrderDetailValues) SetPickupLatitude(lat float64) *OrderDetailValues {
	o.PickupLatitude = &lat
	return o
}

// GetPickupLongitude 获取上车经度
func (o *OrderDetailValues) GetPickupLongitude() float64 {
	if o.PickupLongitude == nil {
		return 0
	}
	return *o.PickupLongitude
}

// SetPickupLongitude 设置上车经度
func (o *OrderDetailValues) SetPickupLongitude(lng float64) *OrderDetailValues {
	o.PickupLongitude = &lng
	return o
}

// GetPickupLandmark 获取上车地标
func (o *OrderDetailValues) GetPickupLandmark() string {
	if o.PickupLandmark == nil {
		return ""
	}
	return *o.PickupLandmark
}

// SetPickupLandmark 设置上车地标
func (o *OrderDetailValues) SetPickupLandmark(landmark string) *OrderDetailValues {
	o.PickupLandmark = &landmark
	return o
}

// GetDropoffAddress 获取下车地址
func (o *OrderDetailValues) GetDropoffAddress() string {
	if o.DropoffAddress == nil {
		return ""
	}
	return *o.DropoffAddress
}

// SetDropoffAddress 设置下车地址
func (o *OrderDetailValues) SetDropoffAddress(address string) *OrderDetailValues {
	o.DropoffAddress = &address
	return o
}

// GetDropoffLatitude 获取下车纬度
func (o *OrderDetailValues) GetDropoffLatitude() float64 {
	if o.DropoffLatitude == nil {
		return 0
	}
	return *o.DropoffLatitude
}

// SetDropoffLatitude 设置下车纬度
func (o *OrderDetailValues) SetDropoffLatitude(lat float64) *OrderDetailValues {
	o.DropoffLatitude = &lat
	return o
}

// GetDropoffLongitude 获取下车经度
func (o *OrderDetailValues) GetDropoffLongitude() float64 {
	if o.DropoffLongitude == nil {
		return 0
	}
	return *o.DropoffLongitude
}

// SetDropoffLongitude 设置下车经度
func (o *OrderDetailValues) SetDropoffLongitude(lng float64) *OrderDetailValues {
	o.DropoffLongitude = &lng
	return o
}

// GetDropoffLandmark 获取下车地标
func (o *OrderDetailValues) GetDropoffLandmark() string {
	if o.DropoffLandmark == nil {
		return ""
	}
	return *o.DropoffLandmark
}

// SetDropoffLandmark 设置下车地标
func (o *OrderDetailValues) SetDropoffLandmark(landmark string) *OrderDetailValues {
	o.DropoffLandmark = &landmark
	return o
}

// GetEstimatedDistance 获取预估距离
func (o *OrderDetailValues) GetEstimatedDistance() float64 {
	if o.EstimatedDistance == nil {
		return 0
	}
	return *o.EstimatedDistance
}

// SetEstimatedDistance 设置预估距离
func (o *OrderDetailValues) SetEstimatedDistance(distance float64) *OrderDetailValues {
	o.EstimatedDistance = &distance
	return o
}

// GetEstimatedDuration 获取预估时长
func (o *OrderDetailValues) GetEstimatedDuration() int {
	if o.EstimatedDuration == nil {
		return 0
	}
	return *o.EstimatedDuration
}

// SetEstimatedDuration 设置预估时长
func (o *OrderDetailValues) SetEstimatedDuration(duration int) *OrderDetailValues {
	o.EstimatedDuration = &duration
	return o
}

// GetActualDistance 获取实际距离
func (o *OrderDetailValues) GetActualDistance() float64 {
	if o.ActualDistance == nil {
		return 0
	}
	return *o.ActualDistance
}

// SetActualDistance 设置实际距离
func (o *OrderDetailValues) SetActualDistance(distance float64) *OrderDetailValues {
	o.ActualDistance = &distance
	return o
}

// GetActualDuration 获取实际时长
func (o *OrderDetailValues) GetActualDuration() int {
	if o.ActualDuration == nil {
		return 0
	}
	return *o.ActualDuration
}

// SetActualDuration 设置实际时长
func (o *OrderDetailValues) SetActualDuration(duration int) *OrderDetailValues {
	o.ActualDuration = &duration
	return o
}

// GetBaseFare 获取基础费用
func (o *OrderDetailValues) GetBaseFare() float64 {
	if o.BaseFare == nil {
		return 0
	}
	return *o.BaseFare
}

// SetBaseFare 设置基础费用
func (o *OrderDetailValues) SetBaseFare(fare float64) *OrderDetailValues {
	o.BaseFare = &fare
	return o
}

// GetDistanceFare 获取距离费用
func (o *OrderDetailValues) GetDistanceFare() float64 {
	if o.DistanceFare == nil {
		return 0
	}
	return *o.DistanceFare
}

// SetDistanceFare 设置距离费用
func (o *OrderDetailValues) SetDistanceFare(fare float64) *OrderDetailValues {
	o.DistanceFare = &fare
	return o
}

// GetTimeFare 获取时间费用
func (o *OrderDetailValues) GetTimeFare() float64 {
	if o.TimeFare == nil {
		return 0
	}
	return *o.TimeFare
}

// SetTimeFare 设置时间费用
func (o *OrderDetailValues) SetTimeFare(fare float64) *OrderDetailValues {
	o.TimeFare = &fare
	return o
}

// GetSurgeFare 获取高峰费用
func (o *OrderDetailValues) GetSurgeFare() float64 {
	if o.SurgeFare == nil {
		return 0
	}
	return *o.SurgeFare
}

// SetSurgeFare 设置高峰费用
func (o *OrderDetailValues) SetSurgeFare(fare float64) *OrderDetailValues {
	o.SurgeFare = &fare
	return o
}

// GetDriverEnRouteAt 获取司机出发时间
func (o *OrderDetailValues) GetDriverEnRouteAt() int64 {
	if o.DriverEnRouteAt == nil {
		return 0
	}
	return *o.DriverEnRouteAt
}

// SetDriverEnRouteAt 设置司机出发时间
func (o *OrderDetailValues) SetDriverEnRouteAt(timestamp int64) *OrderDetailValues {
	o.DriverEnRouteAt = &timestamp
	return o
}

// GetArrivedAt 获取到达时间
func (o *OrderDetailValues) GetArrivedAt() int64 {
	if o.ArrivedAt == nil {
		return 0
	}
	return *o.ArrivedAt
}

// SetArrivedAt 设置到达时间
func (o *OrderDetailValues) SetArrivedAt(timestamp int64) *OrderDetailValues {
	o.ArrivedAt = &timestamp
	return o
}

// SetDriverArrived 设置司机到达时间（便利方法）
func (o *OrderDetailValues) SetDriverArrived() *OrderDetailValues {
	now := utils.TimeNowMilli()
	o.ArrivedAt = &now
	return o
}

// GetRouteData 获取路线数据
func (o *OrderDetailValues) GetRouteData() string {
	if o.RouteData == nil {
		return ""
	}
	return *o.RouteData
}

// SetRouteData 设置路线数据
func (o *OrderDetailValues) SetRouteData(data string) *OrderDetailValues {
	o.RouteData = &data
	return o
}

// ================== 外卖字段 Getter/Setter ==================

// GetRestaurantID 获取餐厅ID
func (o *OrderDetailValues) GetRestaurantID() string {
	if o.RestaurantID == nil {
		return ""
	}
	return *o.RestaurantID
}

// SetRestaurantID 设置餐厅ID
func (o *OrderDetailValues) SetRestaurantID(restaurantID string) *OrderDetailValues {
	o.RestaurantID = &restaurantID
	return o
}

// GetRestaurantName 获取餐厅名称
func (o *OrderDetailValues) GetRestaurantName() string {
	if o.RestaurantName == nil {
		return ""
	}
	return *o.RestaurantName
}

// SetRestaurantName 设置餐厅名称
func (o *OrderDetailValues) SetRestaurantName(name string) *OrderDetailValues {
	o.RestaurantName = &name
	return o
}

// GetRestaurantAddress 获取餐厅地址
func (o *OrderDetailValues) GetRestaurantAddress() string {
	if o.RestaurantAddress == nil {
		return ""
	}
	return *o.RestaurantAddress
}

// SetRestaurantAddress 设置餐厅地址
func (o *OrderDetailValues) SetRestaurantAddress(address string) *OrderDetailValues {
	o.RestaurantAddress = &address
	return o
}

// GetRestaurantPhone 获取餐厅电话
func (o *OrderDetailValues) GetRestaurantPhone() string {
	if o.RestaurantPhone == nil {
		return ""
	}
	return *o.RestaurantPhone
}

// SetRestaurantPhone 设置餐厅电话
func (o *OrderDetailValues) SetRestaurantPhone(phone string) *OrderDetailValues {
	o.RestaurantPhone = &phone
	return o
}

// GetDeliveryAddress 获取配送地址
func (o *OrderDetailValues) GetDeliveryAddress() string {
	if o.DeliveryAddress == nil {
		return ""
	}
	return *o.DeliveryAddress
}

// SetDeliveryAddress 设置配送地址
func (o *OrderDetailValues) SetDeliveryAddress(address string) *OrderDetailValues {
	o.DeliveryAddress = &address
	return o
}

// GetDeliveryLatitude 获取配送纬度
func (o *OrderDetailValues) GetDeliveryLatitude() float64 {
	if o.DeliveryLatitude == nil {
		return 0
	}
	return *o.DeliveryLatitude
}

// SetDeliveryLatitude 设置配送纬度
func (o *OrderDetailValues) SetDeliveryLatitude(lat float64) *OrderDetailValues {
	o.DeliveryLatitude = &lat
	return o
}

// GetDeliveryLongitude 获取配送经度
func (o *OrderDetailValues) GetDeliveryLongitude() float64 {
	if o.DeliveryLongitude == nil {
		return 0
	}
	return *o.DeliveryLongitude
}

// SetDeliveryLongitude 设置配送经度
func (o *OrderDetailValues) SetDeliveryLongitude(lng float64) *OrderDetailValues {
	o.DeliveryLongitude = &lng
	return o
}

// GetCourierID 获取快递员ID
func (o *OrderDetailValues) GetCourierID() string {
	if o.CourierID == nil {
		return ""
	}
	return *o.CourierID
}

// SetCourierID 设置快递员ID
func (o *OrderDetailValues) SetCourierID(courierID string) *OrderDetailValues {
	o.CourierID = &courierID
	return o
}

// GetCourierName 获取快递员姓名
func (o *OrderDetailValues) GetCourierName() string {
	if o.CourierName == nil {
		return ""
	}
	return *o.CourierName
}

// SetCourierName 设置快递员姓名
func (o *OrderDetailValues) SetCourierName(name string) *OrderDetailValues {
	o.CourierName = &name
	return o
}

// GetCourierPhone 获取快递员电话
func (o *OrderDetailValues) GetCourierPhone() string {
	if o.CourierPhone == nil {
		return ""
	}
	return *o.CourierPhone
}

// SetCourierPhone 设置快递员电话
func (o *OrderDetailValues) SetCourierPhone(phone string) *OrderDetailValues {
	o.CourierPhone = &phone
	return o
}

// GetCourierRating 获取快递员评分
func (o *OrderDetailValues) GetCourierRating() float64 {
	if o.CourierRating == nil {
		return 0
	}
	return *o.CourierRating
}

// SetCourierRating 设置快递员评分
func (o *OrderDetailValues) SetCourierRating(rating float64) *OrderDetailValues {
	o.CourierRating = &rating
	return o
}

// ================== 购物字段 Getter/Setter ==================

// GetStoreID 获取商店ID
func (o *OrderDetailValues) GetStoreID() string {
	if o.StoreID == nil {
		return ""
	}
	return *o.StoreID
}

// SetStoreID 设置商店ID
func (o *OrderDetailValues) SetStoreID(storeID string) *OrderDetailValues {
	o.StoreID = &storeID
	return o
}

// GetStoreName 获取商店名称
func (o *OrderDetailValues) GetStoreName() string {
	if o.StoreName == nil {
		return ""
	}
	return *o.StoreName
}

// SetStoreName 设置商店名称
func (o *OrderDetailValues) SetStoreName(name string) *OrderDetailValues {
	o.StoreName = &name
	return o
}

// GetStoreAddress 获取商店地址
func (o *OrderDetailValues) GetStoreAddress() string {
	if o.StoreAddress == nil {
		return ""
	}
	return *o.StoreAddress
}

// SetStoreAddress 设置商店地址
func (o *OrderDetailValues) SetStoreAddress(address string) *OrderDetailValues {
	o.StoreAddress = &address
	return o
}

// GetShippingAddress 获取收货地址
func (o *OrderDetailValues) GetShippingAddress() string {
	if o.ShippingAddress == nil {
		return ""
	}
	return *o.ShippingAddress
}

// SetShippingAddress 设置收货地址
func (o *OrderDetailValues) SetShippingAddress(address string) *OrderDetailValues {
	o.ShippingAddress = &address
	return o
}

// GetTrackingNumber 获取快递单号
func (o *OrderDetailValues) GetTrackingNumber() string {
	if o.TrackingNumber == nil {
		return ""
	}
	return *o.TrackingNumber
}

// SetTrackingNumber 设置快递单号
func (o *OrderDetailValues) SetTrackingNumber(number string) *OrderDetailValues {
	o.TrackingNumber = &number
	return o
}
