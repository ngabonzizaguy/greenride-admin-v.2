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

	// 版本更新提示（由移动端在启动/首页读取并决定提示策略）
	UpdateNoticeEnabled *bool   `json:"update_notice_enabled" gorm:"column:update_notice_enabled;default:false"`
	UpdateNoticeTitle   *string `json:"update_notice_title" gorm:"column:update_notice_title;type:varchar(120)"`
	UpdateNoticeMessage *string `json:"update_notice_message" gorm:"column:update_notice_message;type:text"`
	ForceUpdateEnabled  *bool   `json:"force_update_enabled" gorm:"column:force_update_enabled;default:false"`
	MinimumAppVersion   *string `json:"minimum_app_version" gorm:"column:minimum_app_version;type:varchar(40)"`
	LatestAppVersion    *string `json:"latest_app_version" gorm:"column:latest_app_version;type:varchar(40)"`
	AndroidStoreURL     *string `json:"android_store_url" gorm:"column:android_store_url;type:varchar(500)"`
	IOSStoreURL         *string `json:"ios_store_url" gorm:"column:ios_store_url;type:varchar(500)"`

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
			MaintenanceMode:     utils.BoolPtr(false),
			MaintenanceMessage:  utils.StringPtr("We're currently improving your experience. Our services will resume shortly."),
			MaintenancePhone:    utils.StringPtr("6996"),
			UpdateNoticeEnabled: utils.BoolPtr(false),
			UpdateNoticeTitle:   utils.StringPtr("Update available"),
			UpdateNoticeMessage: utils.StringPtr("A new version of Green Ride is available. Please update for the best experience."),
			ForceUpdateEnabled:  utils.BoolPtr(false),
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

func (c *SystemConfigValues) GetUpdateNoticeEnabled() bool {
	if c == nil || c.UpdateNoticeEnabled == nil {
		return false
	}
	return *c.UpdateNoticeEnabled
}

func (c *SystemConfigValues) GetUpdateNoticeTitle() string {
	if c == nil || c.UpdateNoticeTitle == nil {
		return ""
	}
	return *c.UpdateNoticeTitle
}

func (c *SystemConfigValues) GetUpdateNoticeMessage() string {
	if c == nil || c.UpdateNoticeMessage == nil {
		return ""
	}
	return *c.UpdateNoticeMessage
}

func (c *SystemConfigValues) GetForceUpdateEnabled() bool {
	if c == nil || c.ForceUpdateEnabled == nil {
		return false
	}
	return *c.ForceUpdateEnabled
}

func (c *SystemConfigValues) GetMinimumAppVersion() string {
	if c == nil || c.MinimumAppVersion == nil {
		return ""
	}
	return *c.MinimumAppVersion
}

func (c *SystemConfigValues) GetLatestAppVersion() string {
	if c == nil || c.LatestAppVersion == nil {
		return ""
	}
	return *c.LatestAppVersion
}

func (c *SystemConfigValues) GetAndroidStoreURL() string {
	if c == nil || c.AndroidStoreURL == nil {
		return ""
	}
	return *c.AndroidStoreURL
}

func (c *SystemConfigValues) GetIOSStoreURL() string {
	if c == nil || c.IOSStoreURL == nil {
		return ""
	}
	return *c.IOSStoreURL
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

func (c *SystemConfigValues) SetUpdateNoticeEnabled(enabled bool) *SystemConfigValues {
	c.UpdateNoticeEnabled = &enabled
	return c
}

func (c *SystemConfigValues) SetUpdateNoticeTitle(title string) *SystemConfigValues {
	c.UpdateNoticeTitle = &title
	return c
}

func (c *SystemConfigValues) SetUpdateNoticeMessage(message string) *SystemConfigValues {
	c.UpdateNoticeMessage = &message
	return c
}

func (c *SystemConfigValues) SetForceUpdateEnabled(enabled bool) *SystemConfigValues {
	c.ForceUpdateEnabled = &enabled
	return c
}

func (c *SystemConfigValues) SetMinimumAppVersion(version string) *SystemConfigValues {
	c.MinimumAppVersion = &version
	return c
}

func (c *SystemConfigValues) SetLatestAppVersion(version string) *SystemConfigValues {
	c.LatestAppVersion = &version
	return c
}

func (c *SystemConfigValues) SetAndroidStoreURL(url string) *SystemConfigValues {
	c.AndroidStoreURL = &url
	return c
}

func (c *SystemConfigValues) SetIOSStoreURL(url string) *SystemConfigValues {
	c.IOSStoreURL = &url
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
