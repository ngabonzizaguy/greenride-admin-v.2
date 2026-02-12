package models

import (
	"greenride/internal/utils"
)

// SystemConfig 系统全局配置表
type SystemConfig struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ConfigID string `json:"config_id" gorm:"column:config_id;type:varchar(64);uniqueIndex"`
	*SystemConfigValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type SystemConfigValues struct {
	// 维护模式
	MaintenanceMode    *bool   `json:"maintenance_mode" gorm:"column:maintenance_mode;default:false"`
	MaintenanceMessage *string `json:"maintenance_message" gorm:"column:maintenance_message;type:text"`
	MaintenancePhone   *string `json:"maintenance_phone" gorm:"column:maintenance_phone;type:varchar(50)"`
	MaintenanceStartAt int64   `json:"maintenance_started_at" gorm:"column:maintenance_started_at"`

	// 元数据
	UpdatedBy *string `json:"updated_by" gorm:"column:updated_by;type:varchar(64)"`
	UpdatedAt int64   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (SystemConfig) TableName() string {
	return "t_system_config"
}

// NewSystemConfig 创建默认系统配置
func NewSystemConfig() *SystemConfig {
	return &SystemConfig{
		ConfigID: utils.GenerateID(),
		SystemConfigValues: &SystemConfigValues{
			MaintenanceMode:    utils.BoolPtr(false),
			MaintenanceMessage: utils.StringPtr("We're currently improving your experience. Our services will resume shortly."),
			MaintenancePhone:   utils.StringPtr("6996"),
		},
	}
}

// Getter methods
func (c *SystemConfigValues) GetMaintenanceMode() bool {
	if c == nil || c.MaintenanceMode == nil {
		return false
	}
	return *c.MaintenanceMode
}

func (c *SystemConfigValues) GetMaintenanceMessage() string {
	if c == nil || c.MaintenanceMessage == nil {
		return ""
	}
	return *c.MaintenanceMessage
}

func (c *SystemConfigValues) GetMaintenancePhone() string {
	if c == nil || c.MaintenancePhone == nil {
		return "6996"
	}
	return *c.MaintenancePhone
}

// Setter methods
func (c *SystemConfigValues) SetMaintenanceMode(enabled bool) *SystemConfigValues {
	c.MaintenanceMode = &enabled
	return c
}

func (c *SystemConfigValues) SetMaintenanceMessage(msg string) *SystemConfigValues {
	c.MaintenanceMessage = &msg
	return c
}

func (c *SystemConfigValues) SetMaintenancePhone(phone string) *SystemConfigValues {
	c.MaintenancePhone = &phone
	return c
}

func (c *SystemConfigValues) SetUpdatedBy(adminID string) *SystemConfigValues {
	c.UpdatedBy = &adminID
	return c
}

// GetSystemConfig 从数据库获取系统配置（单例）
func GetSystemConfig() *SystemConfig {
	db := GetDB()
	if db == nil {
		return NewSystemConfig()
	}

	var config SystemConfig
	if err := db.First(&config).Error; err != nil {
		return nil
	}
	return &config
}
