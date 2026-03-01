package services

import (
	"time"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
)

const systemConfigCacheKey = "greenride:system_config"
const systemConfigCacheTTL = 30 // seconds

// SystemConfigService 系统配置服务
type SystemConfigService struct{}

var systemConfigServiceInstance *SystemConfigService

// GetSystemConfigService 获取系统配置服务实例
func GetSystemConfigService() *SystemConfigService {
	if systemConfigServiceInstance == nil {
		systemConfigServiceInstance = &SystemConfigService{}
	}
	return systemConfigServiceInstance
}

// GetConfig 获取系统配置（先查Redis缓存，再查DB）
func (s *SystemConfigService) GetConfig() *protocol.SystemConfigResponse {
	// Try cache first
	cached, err := models.GetObjectFromCache[protocol.SystemConfigResponse](systemConfigCacheKey)
	if err == nil && cached != nil {
		return cached
	}

	// Cache miss - load from DB
	db := models.GetDB()
	if db == nil {
		return &protocol.SystemConfigResponse{MaintenanceMode: false}
	}

	var config models.SystemConfig
	if err := db.First(&config).Error; err != nil {
		// Not found - create default
		config = *models.NewSystemConfig()
		if createErr := db.Create(&config).Error; createErr != nil {
			log.Get().Errorf("Failed to create default system config: %v", createErr)
			return &protocol.SystemConfigResponse{MaintenanceMode: false}
		}
	}

	resp := toSystemConfigResponse(&config)

	// Cache the result
	_ = models.SetObjectCache(systemConfigCacheKey, resp, time.Duration(systemConfigCacheTTL)*time.Second)

	return resp
}

// IsMaintenanceMode 快速检查维护模式状态（供中间件使用）
func (s *SystemConfigService) IsMaintenanceMode() bool {
	config := s.GetConfig()
	return config != nil && config.MaintenanceMode
}

// UpdateConfig 更新系统配置
func (s *SystemConfigService) UpdateConfig(req *protocol.SystemConfigUpdateRequest, adminID string) error {
	db := models.GetDB()
	if db == nil {
		return nil
	}

	var config models.SystemConfig
	if err := db.First(&config).Error; err != nil {
		config = *models.NewSystemConfig()
	}

	updates := make(map[string]interface{})

	if req.MaintenanceMode != nil {
		updates["maintenance_mode"] = *req.MaintenanceMode
		// When enabling maintenance, record the start time
		if *req.MaintenanceMode {
			updates["maintenance_started_at"] = time.Now().UnixMilli()
		} else {
			updates["maintenance_started_at"] = 0
		}
	}
	if req.MaintenanceMessage != nil {
		updates["maintenance_message"] = *req.MaintenanceMessage
	}
	if req.MaintenancePhone != nil {
		updates["maintenance_phone"] = *req.MaintenancePhone
	}
	if req.UpdateNoticeEnabled != nil {
		updates["update_notice_enabled"] = *req.UpdateNoticeEnabled
	}
	if req.UpdateNoticeTitle != nil {
		updates["update_notice_title"] = *req.UpdateNoticeTitle
	}
	if req.UpdateNoticeMessage != nil {
		updates["update_notice_message"] = *req.UpdateNoticeMessage
	}
	if req.ForceUpdateEnabled != nil {
		updates["force_update_enabled"] = *req.ForceUpdateEnabled
	}
	if req.MinimumAppVersion != nil {
		updates["minimum_app_version"] = *req.MinimumAppVersion
	}
	if req.LatestAppVersion != nil {
		updates["latest_app_version"] = *req.LatestAppVersion
	}
	if req.AndroidStoreURL != nil {
		updates["android_store_url"] = *req.AndroidStoreURL
	}
	if req.IOSStoreURL != nil {
		updates["ios_store_url"] = *req.IOSStoreURL
	}

	updates["updated_by"] = adminID

	var dbErr error
	if config.ID == 0 {
		// Apply updates to new config
		if req.MaintenanceMode != nil {
			config.SetMaintenanceMode(*req.MaintenanceMode)
		}
		if req.MaintenanceMessage != nil {
			config.SetMaintenanceMessage(*req.MaintenanceMessage)
		}
		if req.MaintenancePhone != nil {
			config.SetMaintenancePhone(*req.MaintenancePhone)
		}
		if req.UpdateNoticeEnabled != nil {
			config.SetUpdateNoticeEnabled(*req.UpdateNoticeEnabled)
		}
		if req.UpdateNoticeTitle != nil {
			config.SetUpdateNoticeTitle(*req.UpdateNoticeTitle)
		}
		if req.UpdateNoticeMessage != nil {
			config.SetUpdateNoticeMessage(*req.UpdateNoticeMessage)
		}
		if req.ForceUpdateEnabled != nil {
			config.SetForceUpdateEnabled(*req.ForceUpdateEnabled)
		}
		if req.MinimumAppVersion != nil {
			config.SetMinimumAppVersion(*req.MinimumAppVersion)
		}
		if req.LatestAppVersion != nil {
			config.SetLatestAppVersion(*req.LatestAppVersion)
		}
		if req.AndroidStoreURL != nil {
			config.SetAndroidStoreURL(*req.AndroidStoreURL)
		}
		if req.IOSStoreURL != nil {
			config.SetIOSStoreURL(*req.IOSStoreURL)
		}
		config.SetUpdatedBy(adminID)
		dbErr = db.Create(&config).Error
	} else {
		dbErr = db.Model(&config).Updates(updates).Error
	}

	if dbErr != nil {
		return dbErr
	}

	// Invalidate cache so next read picks up the new value
	_ = models.Delete(systemConfigCacheKey)

	return nil
}

func toSystemConfigResponse(config *models.SystemConfig) *protocol.SystemConfigResponse {
	return &protocol.SystemConfigResponse{
		MaintenanceMode:     config.GetMaintenanceMode(),
		MaintenanceMessage:  config.GetMaintenanceMessage(),
		MaintenancePhone:    config.GetMaintenancePhone(),
		MaintenanceStartAt:  config.MaintenanceStartAt,
		UpdateNoticeEnabled: config.GetUpdateNoticeEnabled(),
		UpdateNoticeTitle:   config.GetUpdateNoticeTitle(),
		UpdateNoticeMessage: config.GetUpdateNoticeMessage(),
		ForceUpdateEnabled:  config.GetForceUpdateEnabled(),
		MinimumAppVersion:   config.GetMinimumAppVersion(),
		LatestAppVersion:    config.GetLatestAppVersion(),
		AndroidStoreURL:     config.GetAndroidStoreURL(),
		IOSStoreURL:         config.GetIOSStoreURL(),
	}
}
