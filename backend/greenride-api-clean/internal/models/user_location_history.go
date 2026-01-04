package models

// UserLocationHistory 用户位置历史记录表
type UserLocationHistory struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(64);index;not null"`
	*UserLocationHistoryValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserLocationHistoryValues struct {
	Latitude       float64  `json:"latitude" gorm:"column:latitude;type:decimal(10,8);not null"`
	Longitude      float64  `json:"longitude" gorm:"column:longitude;type:decimal(11,8);not null"`
	Altitude       *float64 `json:"altitude" gorm:"column:altitude;type:decimal(8,2)"`
	Accuracy       *float64 `json:"accuracy" gorm:"column:accuracy;type:decimal(8,2)"`
	Heading        *float64 `json:"heading" gorm:"column:heading;type:decimal(6,2)"`
	Speed          *float64 `json:"speed" gorm:"column:speed;type:decimal(6,2)"`
	OnlineStatus   string   `json:"online_status" gorm:"column:online_status;type:varchar(20);not null"` // online, offline, busy
	CurrentOrderID *string  `json:"current_order_id" gorm:"column:current_order_id;type:varchar(64)"`
	LocationSource string   `json:"location_source" gorm:"column:location_source;type:varchar(20);default:'gps'"` // gps, network, manual
	RecordedAt     int64    `json:"recorded_at" gorm:"column:recorded_at;index;not null"`                         // 位置记录时间戳，与User表的location_updated_at保持一致
}

func (UserLocationHistory) TableName() string {
	return "t_user_location_history"
}

// NewUserLocationHistory 创建新的位置历史记录
func NewUserLocationHistory(userID string, lat, lng float64, onlineStatus string, recordedAt int64) *UserLocationHistory {
	return &UserLocationHistory{
		UserID: userID,
		UserLocationHistoryValues: &UserLocationHistoryValues{
			Latitude:       lat,
			Longitude:      lng,
			OnlineStatus:   onlineStatus,
			LocationSource: "gps",
			RecordedAt:     recordedAt, // 使用传入的统一时间戳
		},
	}
}

// SetValues 更新UserLocationHistoryValues中的非nil值
func (h *UserLocationHistoryValues) SetValues(values *UserLocationHistoryValues) {
	if values == nil {
		return
	}

	if values.Latitude != 0 {
		h.Latitude = values.Latitude
	}
	if values.Longitude != 0 {
		h.Longitude = values.Longitude
	}
	if values.Altitude != nil {
		h.Altitude = values.Altitude
	}
	if values.Accuracy != nil {
		h.Accuracy = values.Accuracy
	}
	if values.Heading != nil {
		h.Heading = values.Heading
	}
	if values.Speed != nil {
		h.Speed = values.Speed
	}
	if values.OnlineStatus != "" {
		h.OnlineStatus = values.OnlineStatus
	}
	if values.CurrentOrderID != nil {
		h.CurrentOrderID = values.CurrentOrderID
	}
	if values.LocationSource != "" {
		h.LocationSource = values.LocationSource
	}
	if values.RecordedAt > 0 {
		h.RecordedAt = values.RecordedAt
	}
}

// Setter 方法
func (h *UserLocationHistoryValues) SetLatitude(latitude float64) *UserLocationHistoryValues {
	h.Latitude = latitude
	return h
}

func (h *UserLocationHistoryValues) SetLongitude(longitude float64) *UserLocationHistoryValues {
	h.Longitude = longitude
	return h
}

func (h *UserLocationHistoryValues) SetAltitude(altitude float64) *UserLocationHistoryValues {
	h.Altitude = &altitude
	return h
}

func (h *UserLocationHistoryValues) SetAccuracy(accuracy float64) *UserLocationHistoryValues {
	h.Accuracy = &accuracy
	return h
}

func (h *UserLocationHistoryValues) SetHeading(heading float64) *UserLocationHistoryValues {
	h.Heading = &heading
	return h
}

func (h *UserLocationHistoryValues) SetSpeed(speed float64) *UserLocationHistoryValues {
	h.Speed = &speed
	return h
}

func (h *UserLocationHistoryValues) SetOnlineStatus(status string) *UserLocationHistoryValues {
	h.OnlineStatus = status
	return h
}

func (h *UserLocationHistoryValues) SetCurrentOrderID(orderID string) *UserLocationHistoryValues {
	h.CurrentOrderID = &orderID
	return h
}

func (h *UserLocationHistoryValues) SetLocationSource(source string) *UserLocationHistoryValues {
	h.LocationSource = source
	return h
}

func (h *UserLocationHistoryValues) SetRecordedAt(timestamp int64) *UserLocationHistoryValues {
	h.RecordedAt = timestamp
	return h
}

// Getter 方法
func (h *UserLocationHistoryValues) GetLatitude() float64 {
	return h.Latitude
}

func (h *UserLocationHistoryValues) GetLongitude() float64 {
	return h.Longitude
}

func (h *UserLocationHistoryValues) GetAltitude() float64 {
	if h.Altitude == nil {
		return 0.0
	}
	return *h.Altitude
}

func (h *UserLocationHistoryValues) GetAccuracy() float64 {
	if h.Accuracy == nil {
		return 0.0
	}
	return *h.Accuracy
}

func (h *UserLocationHistoryValues) GetHeading() float64 {
	if h.Heading == nil {
		return 0.0
	}
	return *h.Heading
}

func (h *UserLocationHistoryValues) GetSpeed() float64 {
	if h.Speed == nil {
		return 0.0
	}
	return *h.Speed
}

func (h *UserLocationHistoryValues) GetOnlineStatus() string {
	return h.OnlineStatus
}

func (h *UserLocationHistoryValues) GetCurrentOrderID() string {
	if h.CurrentOrderID == nil {
		return ""
	}
	return *h.CurrentOrderID
}

func (h *UserLocationHistoryValues) GetLocationSource() string {
	return h.LocationSource
}

func (h *UserLocationHistoryValues) GetRecordedAt() int64 {
	return h.RecordedAt
}

// 业务方法
func (h *UserLocationHistory) IsOnline() bool {
	return h.OnlineStatus == "online"
}

func (h *UserLocationHistory) IsBusy() bool {
	return h.OnlineStatus == "busy"
}

func (h *UserLocationHistory) HasCurrentOrder() bool {
	return h.CurrentOrderID != nil && *h.CurrentOrderID != ""
}

// 批量插入历史记录的辅助方法
func CreateLocationHistoryBatch(histories []*UserLocationHistory) error {
	if len(histories) == 0 {
		return nil
	}

	return DB.CreateInBatches(histories, 100).Error
}

// 获取用户位置历史记录
func GetUserLocationHistory(userID string, limit int, offset int) ([]*UserLocationHistory, error) {
	var histories []*UserLocationHistory
	err := DB.Where("user_id = ?", userID).
		Order("recorded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&histories).Error
	return histories, err
}

// 获取用户指定时间范围内的位置历史
func GetUserLocationHistoryByTimeRange(userID string, startTime, endTime int64) ([]*UserLocationHistory, error) {
	var histories []*UserLocationHistory
	err := DB.Where("user_id = ? AND recorded_at >= ? AND recorded_at <= ?", userID, startTime, endTime).
		Order("recorded_at ASC").
		Find(&histories).Error
	return histories, err
}

// 清理过期的位置历史记录
func CleanupOldLocationHistory(beforeTime int64) error {
	return DB.Where("recorded_at < ?", beforeTime).Delete(&UserLocationHistory{}).Error
}
