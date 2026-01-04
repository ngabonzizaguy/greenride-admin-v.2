package services

import "sync"

type AdminUserService struct {
}

var (
	adminUserInstance *AdminUserService
	adminUserOnce     sync.Once
)

func GetAdminUserService() *AdminUserService {
	adminUserOnce.Do(func() {
		SetupAdminUserService()
	})
	return adminUserInstance
}
func SetupAdminUserService() {
	adminUserInstance = &AdminUserService{}
}
