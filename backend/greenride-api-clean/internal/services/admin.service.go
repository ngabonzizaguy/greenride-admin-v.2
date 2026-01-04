package services

import (
	"greenride/internal/models"
	"sync"

	"gorm.io/gorm"
)

// AdminService 管理员服务
type AdminService struct {
	db *gorm.DB
}

var (
	adminServiceInstance *AdminService
	adminServiceOnce     sync.Once
)

// GetAdminService 获取管理员服务单例
func GetAdminService() *AdminService {
	adminServiceOnce.Do(func() {
		SetupAdminService()
	})
	return adminServiceInstance
}

// SetupAdminService 设置管理员服务
func SetupAdminService() {
	adminServiceInstance = &AdminService{
		db: models.GetDB(),
	}
}
