package services

import (
	"log"
	"sync"

	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

type AdminAdminService struct {
}

var (
	adminAdminInstance *AdminAdminService
	adminAdminOnce     sync.Once
)

func GetAdminAdminService() *AdminAdminService {
	adminAdminOnce.Do(func() {
		SetupAdminAdminService()
	})
	return adminAdminInstance
}
func SetupAdminAdminService() {
	adminAdminInstance = &AdminAdminService{}
}

// GetAdminByUsername 根据用户名获取管理员
func (s *AdminAdminService) GetAdminByUsername(username string) *models.Admin {
	var admin models.Admin
	if err := models.GetDB().Where("username = ?", username).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

// GetAdminByEmail 根据邮箱获取管理员
func (s *AdminAdminService) GetAdminByEmail(email string) *models.Admin {
	var admin models.Admin
	if err := models.GetDB().Where("email = ?", email).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

// GetAdminByID 根据管理员ID获取管理员
func (s *AdminAdminService) GetAdminByID(adminID string) *models.Admin {
	var admin models.Admin
	if err := models.GetDB().Where("admin_id = ?", adminID).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

// VerifyPassword 验证管理员密码
func (s *AdminAdminService) VerifyPassword(admin *models.Admin, password string) bool {
	if admin == nil || admin.PasswordHash == nil {
		return false
	}
	return utils.VerifyPassword(password, admin.Salt, *admin.PasswordHash)
}

// UpdatePassword 更新管理员密码
func (s *AdminAdminService) UpdatePassword(adminID, newPassword string) protocol.ErrorCode {
	admin := s.GetAdminByID(adminID)
	if admin == nil {
		return protocol.UserNotFound
	}

	// 生成新的哈希密码
	hashedPassword, err := utils.HashPassword(newPassword, admin.Salt)
	if err != nil {
		return protocol.SystemError
	}

	// 更新密码
	values := &models.AdminValues{}
	values.SetPasswordHash(hashedPassword, admin.Salt)
	values.SetMustChangePassword(false)

	return s.UpdateAdmin(admin, values)
}

// UpdateAdmin 更新管理员信息
func (s *AdminAdminService) UpdateAdmin(admin *models.Admin, values *models.AdminValues) protocol.ErrorCode {
	if admin == nil || values == nil {
		return protocol.InvalidParams
	}

	defer func() {
		// 使用SetValues方法更新非空字段
		admin.SetValues(values)
	}()

	if err := models.GetDB().Model(admin).UpdateColumns(values).Error; err != nil {
		return protocol.SystemError
	}
	return protocol.Success
}

// RecordLogin 记录登录信息
func (s *AdminAdminService) RecordLogin(admin *models.Admin, ip, sessionID string) protocol.ErrorCode {
	values := &models.AdminValues{}
	values.RecordLogin(ip, sessionID)
	values.SetActiveStatus(models.AdminActiveStatusOnline)

	return s.UpdateAdmin(admin, values)
}

// RecordFailedLogin 记录失败登录
func (s *AdminAdminService) RecordFailedLogin(admin *models.Admin) protocol.ErrorCode {
	values := &models.AdminValues{}
	values.RecordFailedLogin()

	return s.UpdateAdmin(admin, values)
}

// Logout 登出
func (s *AdminAdminService) Logout(admin *models.Admin) protocol.ErrorCode {
	values := &models.AdminValues{}
	values.Logout()

	return s.UpdateAdmin(admin, values)
}

// CheckLoginPermission 检查登录权限
func (s *AdminAdminService) CheckLoginPermission(admin *models.Admin, ip string) protocol.ErrorCode {
	if admin == nil {
		return protocol.UserNotFound
	}

	// 检查账户状态
	if !admin.CanLogin() {
		if admin.IsLocked() {
			return protocol.AccountLocked
		}
		if admin.IsSuspended() {
			return protocol.AccountSuspended
		}
		if admin.IsInactive() {
			return protocol.AccountDisabled
		}
		return protocol.AccessDenied
	}

	// 检查IP白名单
	if !admin.IsIPAllowed(ip) {
		return protocol.IPNotAllowed
	}

	// 检查是否需要强制修改密码
	if admin.ShouldForcePasswordChange() {
		return protocol.WeakPassword // 使用现有的密码相关错误码
	}

	return protocol.Success
}

// ResetPassword 重置管理员密码（管理员操作）
func (s *AdminAdminService) ResetPassword(adminID, newPassword, operatorID string) protocol.ErrorCode {
	admin := s.GetAdminByID(adminID)
	if admin == nil {
		return protocol.UserNotFound
	}

	// 生成新的哈希密码
	hashedPassword, err := utils.HashPassword(newPassword, admin.Salt)
	if err != nil {
		return protocol.SystemError
	}

	values := &models.AdminValues{}
	values.SetPasswordHash(hashedPassword, admin.Salt)
	values.SetMustChangePassword(true) // 重置后强制修改密码
	values.LastUpdatedBy = &operatorID

	return s.UpdateAdmin(admin, values)
}

// GetAdminList 获取管理员列表
func (s *AdminAdminService) GetAdminList(keyword string, role string, status string, page, limit int) ([]*models.Admin, int64, protocol.ErrorCode) {
	var admins []*models.Admin
	var total int64

	query := models.GetDB().Model(&models.Admin{})

	// 角色过滤
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键字搜索（用户名、邮箱、全名）
	if keyword != "" {
		searchTerm := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Failed to count admins: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		log.Printf("Failed to get admin list: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	return admins, total, protocol.Success
}

// CreateAdmin 创建新管理员
func (s *AdminAdminService) CreateAdmin(username, email, password, role, department, jobTitle, creatorID string) (*models.Admin, protocol.ErrorCode) {
	// 检查用户名是否已存在
	existingAdmin := s.GetAdminByUsername(username)
	if existingAdmin != nil {
		return nil, protocol.UserAlreadyExists
	}

	// 检查邮箱是否已存在
	existingAdmin = s.GetAdminByEmail(email)
	if existingAdmin != nil {
		return nil, protocol.EmailAlreadyExists
	}

	// 创建新管理员
	admin := models.NewAdminV2()
	admin.SetUsername(username).
		SetEmail(email).
		SetRole(role)

	// 设置可选字段
	if department != "" {
		admin.SetDepartment(department)
	}
	if jobTitle != "" {
		admin.SetJobTitle(jobTitle)
	}

	// 设置创建者
	if creatorID != "" {
		admin.CreatedBy = &creatorID
		admin.LastUpdatedBy = &creatorID
	}

	// 生成密码哈希
	hashedPassword, err := utils.HashPassword(password, admin.Salt)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, protocol.InternalError
	}
	admin.SetPasswordHash(hashedPassword, admin.Salt)

	// 设置角色对应的权限
	rolePermissions := models.GetRolePermissions(role)
	admin.SetPermissions(rolePermissions)

	// 保存到数据库
	if err := models.GetDB().Create(admin).Error; err != nil {
		log.Printf("Failed to create admin: %v", err)
		return nil, protocol.DatabaseError
	}

	return admin, protocol.Success
}

// IsUsernameExists 检查用户名是否存在
func (s *AdminAdminService) IsUsernameExists(username string) bool {
	admin := s.GetAdminByUsername(username)
	return admin != nil
}

// IsEmailExists 检查邮箱是否存在
func (s *AdminAdminService) IsEmailExists(email string) bool {
	admin := s.GetAdminByEmail(email)
	return admin != nil
}

// DeleteAdmin 删除管理员（硬删除）
func (s *AdminAdminService) DeleteAdmin(adminID, operatorID string) protocol.ErrorCode {
	admin := s.GetAdminByID(adminID)
	if admin == nil {
		return protocol.UserNotFound
	}

	// 不能删除自己
	if admin.AdminID == operatorID {
		return protocol.InvalidParams
	}

	// 硬删除管理员记录
	if err := models.GetDB().Where("admin_id = ?", adminID).Delete(&models.Admin{}).Error; err != nil {
		log.Printf("failed to hard delete admin %s: %v", adminID, err)
		return protocol.DatabaseError
	}
	return protocol.Success
}

// EnsureDefaultAdmin 确保至少有一个管理员，如果没有则创建一个默认的
func (s *AdminAdminService) EnsureDefaultAdmin() {
	var count int64
	models.GetDB().Model(&models.Admin{}).Count(&count)
	if count == 0 {
		log.Printf("No admins found in database. Creating default admin account...")
		_, err := s.CreateAdmin("admin", "admin@greenrideafrica.com", "admin123", models.AdminRoleSuperAdmin, "IT", "System Admin", "SYSTEM")
		if err != protocol.Success {
			log.Printf("CRITICAL: Failed to create default admin: %v", err)
		} else {
			log.Printf("Successfully created default admin: admin / admin123")
		}
	}
}
