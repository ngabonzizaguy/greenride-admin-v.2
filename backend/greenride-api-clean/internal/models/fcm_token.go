package models

import (
	"greenride/internal/log"
	"greenride/internal/protocol"
	"time"

	"gorm.io/gorm"
)

// FCMToken FCM推送令牌表
type FCMToken struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID     string `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;uniqueIndex:idx_user_token"`
	Token      string `json:"token" gorm:"column:token;type:varchar(500);not null;uniqueIndex:idx_user_token"`
	DeviceID   string `json:"device_id" gorm:"column:device_id;type:varchar(200);index"`
	Platform   string `json:"platform" gorm:"column:platform;type:varchar(20)"` // ios, android, web
	AppID      string `json:"app_id" gorm:"column:app_id;type:varchar(100)"`    // 应用标识
	Status     string `json:"status" gorm:"column:status;type:varchar(20)"`     // active, inactive
	CreatedAt  int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	LastUsedAt *int64 `json:"last_used_at" gorm:"column:last_used_at"` // 最后使用时间
}

// TableName 表名
func (f *FCMToken) TableName() string {
	return "t_fcm_tokens"
}

type FCMTokenValues struct {
	UserID     *string `json:"user_id"`
	Token      *string `json:"token"`
	DeviceID   *string `json:"device_id"`
	Platform   *string `json:"platform"`
	AppID      *string `json:"app_id"`
	Status     *string `json:"status"`
	UpdatedAt  int64   `json:"updated_at"`
	LastUsedAt *int64  `json:"last_used_at"`
}

// Getter methods for FCMTokenValues
func (f *FCMTokenValues) GetUserID() string {
	if f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *FCMTokenValues) GetToken() string {
	if f.Token != nil {
		return *f.Token
	}
	return ""
}

func (f *FCMTokenValues) GetDeviceID() string {
	if f.DeviceID != nil {
		return *f.DeviceID
	}
	return ""
}

func (f *FCMTokenValues) GetPlatform() string {
	if f.Platform != nil {
		return *f.Platform
	}
	return ""
}

func (f *FCMTokenValues) GetAppID() string {
	if f.AppID != nil {
		return *f.AppID
	}
	return ""
}

func (f *FCMTokenValues) GetStatus() string {
	if f.Status != nil {
		return *f.Status
	}
	return ""
}

func (f *FCMTokenValues) GetLastUsedAt() *int64 {
	return f.LastUsedAt
}

// Setter methods for FCMTokenValues
func (f *FCMTokenValues) SetUserID(userID string) *FCMTokenValues {
	f.UserID = &userID
	return f
}

func (f *FCMTokenValues) SetToken(token string) *FCMTokenValues {
	f.Token = &token
	return f
}

func (f *FCMTokenValues) SetDeviceID(deviceID string) *FCMTokenValues {
	f.DeviceID = &deviceID
	return f
}

func (f *FCMTokenValues) SetPlatform(platform string) *FCMTokenValues {
	f.Platform = &platform
	return f
}

func (f *FCMTokenValues) SetAppID(appID string) *FCMTokenValues {
	f.AppID = &appID
	return f
}

func (f *FCMTokenValues) SetStatus(status string) *FCMTokenValues {
	f.Status = &status
	return f
}

func (f *FCMTokenValues) SetLastUsedAt(timestamp int64) *FCMTokenValues {
	f.LastUsedAt = &timestamp
	return f
}

func (f *FCMTokenValues) UpdateLastUsed() *FCMTokenValues {
	now := time.Now().UnixMilli()
	f.LastUsedAt = &now
	f.UpdatedAt = now
	return f
}

func (f *FCMTokenValues) SetActiveStatus() *FCMTokenValues {
	f.SetStatus(protocol.StatusActive)
	return f
}

func (f *FCMTokenValues) SetInactiveStatus() *FCMTokenValues {
	f.SetStatus(protocol.StatusInactive)
	return f
}

func GetFcmTokenByToken(token string) *FCMToken {
	var fcmToken FCMToken
	if err := GetDB().Where("token = ?", token).First(&fcmToken).Error; err != nil {
		return nil
	}
	return &fcmToken
}

// GetFcmTokenByUserIDAndToken 根据用户ID和Token获取FCM Token记录
func GetFcmTokenByUserIDAndToken(userID, token string) *FCMToken {
	var fcmToken FCMToken
	if err := GetDB().Where("user_id = ? AND token = ?", userID, token).First(&fcmToken).Error; err != nil {
		return nil
	}
	return &fcmToken
}

// CreateFCMToken 创建FCM Token记录
func CreateFCMToken(token *FCMToken) error {
	if err := GetDB().Create(token).Error; err != nil {
		log.Get().Errorf("failed to create fcm token: %v", err)
		return err
	}
	return nil
}

// GetFCMTokensByUserID 根据用户ID获取所有激活的token
func GetFCMTokensByUserID(userID string) []*FCMToken {
	var tokens []*FCMToken
	if err := GetDB().Where("user_id = ? AND status = ?", userID, protocol.StatusActive).Find(&tokens).Error; err != nil {
		log.Get().Errorf("failed to get fcm tokens by user id: %v", err)
		return nil
	}
	return tokens
}

// UpsertFCMToken 插入或更新FCM Token（根据用户ID和Token去重）
func UpsertFCMToken(token *FCMToken) error {
	// 先检查是否存在相同用户ID和Token的记录
	var existing FCMToken
	err := GetDB().Where("user_id = ? AND token = ?", token.UserID, token.Token).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		token.Status = protocol.StatusActive
		return CreateFCMToken(token)
	} else if err != nil {
		log.Get().Errorf("failed to check existing fcm token: %v", err)
		return err
	}

	// 存在，更新记录
	existing.DeviceID = token.DeviceID
	existing.Platform = token.Platform
	existing.AppID = token.AppID
	existing.Status = protocol.StatusActive
	existing.UpdatedAt = time.Now().UnixMilli()

	return GetDB().Save(&existing).Error
}

// CleanupInactiveFCMTokens 清理非激活的token
func CleanupInactiveFCMTokens() error {
	if err := GetDB().Where("status = ?", protocol.StatusInactive).Delete(&FCMToken{}).Error; err != nil {
		log.Get().Errorf("failed to cleanup inactive fcm tokens: %v", err)
		return err
	}
	return nil
}
