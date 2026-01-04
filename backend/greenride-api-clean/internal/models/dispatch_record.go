package models

import (
	"encoding/json"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// DispatchRecord 派单记录模型
type DispatchRecord struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	DispatchID   string `json:"dispatch_id" gorm:"column:dispatch_id;type:varchar(64);uniqueIndex"`
	OrderID      string `json:"order_id" gorm:"column:order_id;type:varchar(64);index"`
	DriverID     string `json:"driver_id" gorm:"column:driver_id;type:varchar(64);index"`
	Round        int    `json:"round" gorm:"column:round;type:int"`
	RoundSeq     int    `json:"round_seq" gorm:"column:round_seq;type:int"`
	DispatchedAt int64  `json:"dispatched_at" gorm:"column:dispatched_at;type:bigint"`
	ExpiredAt    int64  `json:"expired_at" gorm:"column:expired_at;type:bigint"`
	CreatedAt    int64  `json:"created_at" gorm:"column:created_at;type:bigint;autoCreateTime:milli"`
	*DispatchRecordValues
}

// DispatchRecordValues 用于数据操作的值对象
type DispatchRecordValues struct {
	RespondedAt      *int64   `json:"responded_at" gorm:"column:responded_at;type:bigint"`
	Status           *string  `json:"status" gorm:"column:status;type:varchar(20);default:'pending'"`
	RejectReason     *string  `json:"reject_reason" gorm:"column:reject_reason;type:varchar(255)"`          // 保持向后兼容
	RejectReasonType *string  `json:"reject_reason_type" gorm:"column:reject_reason_type;type:varchar(50)"` // 拒绝原因类型（枚举）
	DriverDistance   *float64 `json:"driver_distance" gorm:"column:driver_distance;type:double"`
	DriverLatitude   *float64 `json:"driver_latitude" gorm:"column:driver_latitude;type:double"`   // 司机纬度
	DriverLongitude  *float64 `json:"driver_longitude" gorm:"column:driver_longitude;type:double"` // 司机经度
	StrategyConfig   *string  `json:"strategy_config" gorm:"column:strategy_config;type:json"`
	UpdatedAt        int64    `json:"updated_at" gorm:"column:updated_at;type:bigint;;autoUpdateTime:milli"`
}

// TableName 指定表名
func (DispatchRecord) TableName() string {
	return "t_dispatch_records"
}

func (d *DispatchRecordValues) GetRespondedAt() int64 {
	if d.RespondedAt == nil {
		return 0
	}
	return *d.RespondedAt
}

func (d *DispatchRecordValues) GetStatus() string {
	if d.Status == nil {
		return "pending"
	}
	return *d.Status
}

func (d *DispatchRecordValues) GetRejectReason() string {
	if d.RejectReason == nil {
		return ""
	}
	return *d.RejectReason
}

func (d *DispatchRecordValues) GetRejectReasonType() string {
	if d.RejectReasonType == nil {
		return ""
	}
	return *d.RejectReasonType
}

func (d *DispatchRecordValues) GetDriverDistance() float64 {
	if d.DriverDistance == nil {
		return 0
	}
	return *d.DriverDistance
}

func (d *DispatchRecordValues) GetDriverLatitude() float64 {
	if d.DriverLatitude == nil {
		return 0
	}
	return *d.DriverLatitude
}

func (d *DispatchRecordValues) GetDriverLongitude() float64 {
	if d.DriverLongitude == nil {
		return 0
	}
	return *d.DriverLongitude
}

func (d *DispatchRecordValues) GetStrategyConfig() string {
	if d.StrategyConfig == nil {
		return ""
	}
	return *d.StrategyConfig
}

func (d *DispatchRecordValues) SetRespondedAt(respondedAt int64) *DispatchRecordValues {
	d.RespondedAt = &respondedAt
	return d
}

func (d *DispatchRecordValues) SetStatus(status string) *DispatchRecordValues {
	d.Status = &status
	return d
}

func (d *DispatchRecordValues) SetRejectReason(reason string) *DispatchRecordValues {
	d.RejectReason = &reason
	return d
}

func (d *DispatchRecordValues) SetRejectReasonType(reasonType string) *DispatchRecordValues {
	d.RejectReasonType = &reasonType
	return d
}

func (d *DispatchRecordValues) SetDriverDistance(distance float64) *DispatchRecordValues {
	d.DriverDistance = &distance
	return d
}

func (d *DispatchRecordValues) SetDriverLatitude(latitude float64) *DispatchRecordValues {
	d.DriverLatitude = &latitude
	return d
}

func (d *DispatchRecordValues) SetDriverLongitude(longitude float64) *DispatchRecordValues {
	d.DriverLongitude = &longitude
	return d
}

func (d *DispatchRecordValues) SetDriverCoordinates(latitude, longitude float64) *DispatchRecordValues {
	d.DriverLatitude = &latitude
	d.DriverLongitude = &longitude
	return d
}

func (d *DispatchRecordValues) SetDriverLocation(location string) *DispatchRecordValues {
	// 尝试解析JSON格式的位置信息
	var locationData struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.Unmarshal([]byte(location), &locationData); err == nil {
		d.DriverLatitude = &locationData.Latitude
		d.DriverLongitude = &locationData.Longitude
	}
	return d
}

func (d *DispatchRecordValues) SetStrategyConfig(config string) *DispatchRecordValues {
	d.StrategyConfig = &config
	return d
}

// NewDispatchRecord 创建新的派单记录
func NewDispatchRecord() *DispatchRecord {
	return &DispatchRecord{
		DispatchID: utils.GenerateDispatchID(),
		DispatchRecordValues: &DispatchRecordValues{
			Status: stringPtr("pending"),
		},
	}
}

// stringPtr 辅助函数，返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// SetValues 更新DispatchRecordValues中的非nil值
func (d *DispatchRecordValues) SetValues(values *DispatchRecordValues) {
	if values == nil {
		return
	}
	if values.RespondedAt != nil {
		d.RespondedAt = values.RespondedAt
	}
	if values.Status != nil {
		d.Status = values.Status
	}
	if values.RejectReason != nil {
		d.RejectReason = values.RejectReason
	}
	if values.RejectReasonType != nil {
		d.RejectReasonType = values.RejectReasonType
	}
	if values.DriverDistance != nil {
		d.DriverDistance = values.DriverDistance
	}
	if values.DriverLatitude != nil {
		d.DriverLatitude = values.DriverLatitude
	}
	if values.DriverLongitude != nil {
		d.DriverLongitude = values.DriverLongitude
	}
	if values.StrategyConfig != nil {
		d.StrategyConfig = values.StrategyConfig
	}
}

// Business methods
func (d *DispatchRecord) IsPending() bool {
	return d.GetStatus() == "pending"
}

func (d *DispatchRecord) IsAccepted() bool {
	return d.GetStatus() == "accepted"
}

func (d *DispatchRecord) IsRejected() bool {
	return d.GetStatus() == "rejected"
}

func (d *DispatchRecord) IsTimeout() bool {
	return d.GetStatus() == "timeout"
}

func (d *DispatchRecord) HasExpired() bool {
	return d.ExpiredAt > 0 && d.ExpiredAt < utils.TimeNowMilli()
}

func GetDispatchByID(dispatchId string) *DispatchRecord {
	var record DispatchRecord
	err := GetDB().Where("dispatch_id = ?", dispatchId).First(&record).Error
	if err != nil {
		return nil
	}
	return &record
}

type DispatchRecords []*DispatchRecord

// ToProtocolList 转换为协议列表
func (list DispatchRecords) Protocol() []*protocol.Dispatch {
	if list == nil {
		return nil
	}
	result := make([]*protocol.Dispatch, 0, len(list))
	for _, item := range list {
		if item != nil {
			result = append(result, item.Protocol())
		}
	}
	return result
}

func (t *DispatchRecord) Protocol() *protocol.Dispatch {
	if t == nil {
		return nil
	}
	return &protocol.Dispatch{
		DispatchID:      t.DispatchID,
		OrderID:         t.OrderID,
		DriverID:        t.DriverID,
		Round:           t.Round,
		RoundSeq:        t.RoundSeq,
		DispatchedAt:    t.DispatchedAt,
		ExpiredAt:       t.ExpiredAt,
		RespondedAt:     t.GetRespondedAt(),
		Status:          t.GetStatus(),
		RejectReason:    t.GetRejectReason(),
		DriverDistance:  t.GetDriverDistance(),
		DriverLatitude:  t.GetDriverLatitude(),
		DriverLongitude: t.GetDriverLongitude(),
		StrategyConfig:  t.GetStrategyConfig(),
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
