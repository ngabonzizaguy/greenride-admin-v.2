package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"slices"
	"strings"
)

// Admin 管理员表 - 后台管理员账户、权限和登录管理
type Admin struct {
	ID      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AdminID string `json:"admin_id" gorm:"column:admin_id;type:varchar(64);uniqueIndex"`
	Salt    string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*AdminValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type AdminValues struct {
	// 基本信息
	Username *string `json:"username" gorm:"column:username;type:varchar(50);uniqueIndex"`
	Email    *string `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex"`
	Phone    *string `json:"phone" gorm:"column:phone;type:varchar(20);index"`

	// 个人信息
	FirstName *string `json:"first_name" gorm:"column:first_name;type:varchar(100)"`
	LastName  *string `json:"last_name" gorm:"column:last_name;type:varchar(100)"`
	FullName  *string `json:"full_name" gorm:"column:full_name;type:varchar(200)"`
	Avatar    *string `json:"avatar" gorm:"column:avatar;type:varchar(500)"`

	// 认证信息
	PasswordHash     *string `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	PasswordSalt     *string `json:"-" gorm:"column:password_salt;type:varchar(255)"`
	EmailVerified    *bool   `json:"email_verified" gorm:"column:email_verified;default:false"`
	PhoneVerified    *bool   `json:"phone_verified" gorm:"column:phone_verified;default:false"`
	TwoFactorEnabled *bool   `json:"two_factor_enabled" gorm:"column:two_factor_enabled;default:false"`
	TwoFactorSecret  *string `json:"-" gorm:"column:two_factor_secret;type:varchar(255)"`

	// 角色和权限
	Role        *string `json:"role" gorm:"column:role;type:varchar(50);index"`        // super_admin, admin, moderator, support, analyst
	Permissions *string `json:"permissions" gorm:"column:permissions;type:json"`       // JSON数组存储权限列表
	Department  *string `json:"department" gorm:"column:department;type:varchar(100)"` // 部门
	JobTitle    *string `json:"job_title" gorm:"column:job_title;type:varchar(100)"`   // 职位

	// 状态管理
	Status       *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"`          // active, inactive, suspended, locked
	ActiveStatus *string `json:"active_status" gorm:"column:active_status;type:varchar(32);default:'offline'"` // online, offline, busy

	// 登录相关
	LastLoginAt    *int64  `json:"last_login_at" gorm:"column:last_login_at"`
	LastLoginIP    *string `json:"last_login_ip" gorm:"column:last_login_ip;type:varchar(45)"`
	LoginCount     *int    `json:"login_count" gorm:"column:login_count;default:0"`
	FailedAttempts *int    `json:"failed_attempts" gorm:"column:failed_attempts;default:0"`
	LastFailedAt   *int64  `json:"last_failed_at" gorm:"column:last_failed_at"`
	LockedUntil    *int64  `json:"locked_until" gorm:"column:locked_until"`

	// 会话管理
	CurrentSessionID      *string `json:"current_session_id" gorm:"column:current_session_id;type:varchar(255)"`
	SessionCount          *int    `json:"session_count" gorm:"column:session_count;default:0"`
	MaxConcurrentSessions *int    `json:"max_concurrent_sessions" gorm:"column:max_concurrent_sessions;default:3"`

	// 安全设置
	PasswordChangedAt  *int64  `json:"password_changed_at" gorm:"column:password_changed_at"`
	MustChangePassword *bool   `json:"must_change_password" gorm:"column:must_change_password;default:false"`
	AllowedIPs         *string `json:"allowed_ips" gorm:"column:allowed_ips;type:text"` // 允许的IP地址列表
	SecurityQuestion   *string `json:"security_question" gorm:"column:security_question;type:varchar(255)"`
	SecurityAnswerHash *string `json:"-" gorm:"column:security_answer_hash;type:varchar(255)"`

	// 审计信息
	CreatedBy        *string `json:"created_by" gorm:"column:created_by;type:varchar(64)"`
	LastUpdatedBy    *string `json:"last_updated_by" gorm:"column:last_updated_by;type:varchar(64)"`
	ApprovedBy       *string `json:"approved_by" gorm:"column:approved_by;type:varchar(64)"`
	ApprovedAt       *int64  `json:"approved_at" gorm:"column:approved_at"`
	SuspendedBy      *string `json:"suspended_by" gorm:"column:suspended_by;type:varchar(64)"`
	SuspendedAt      *int64  `json:"suspended_at" gorm:"column:suspended_at"`
	SuspensionReason *string `json:"suspension_reason" gorm:"column:suspension_reason;type:varchar(500)"`

	// 工作时间和区域
	WorkingHours    *string `json:"working_hours" gorm:"column:working_hours;type:json"` // JSON对象存储工作时间配置
	WorkingTimeZone *string `json:"working_timezone" gorm:"column:working_timezone;type:varchar(50)"`
	WorkingAreas    *string `json:"working_areas" gorm:"column:working_areas;type:json"` // JSON数组存储负责区域

	// 联系信息
	OfficePhone      *string `json:"office_phone" gorm:"column:office_phone;type:varchar(20)"`
	OfficeAddress    *string `json:"office_address" gorm:"column:office_address;type:varchar(500)"`
	EmergencyContact *string `json:"emergency_contact" gorm:"column:emergency_contact;type:varchar(255)"`

	// 偏好设置
	Language    *string `json:"language" gorm:"column:language;type:varchar(10);default:'en'"`
	Timezone    *string `json:"timezone" gorm:"column:timezone;type:varchar(50);default:'UTC'"`
	DateFormat  *string `json:"date_format" gorm:"column:date_format;type:varchar(20);default:'YYYY-MM-DD'"`
	TimeFormat  *string `json:"time_format" gorm:"column:time_format;type:varchar(20);default:'24h'"`
	Preferences *string `json:"preferences" gorm:"column:preferences;type:json"` // JSON对象存储个人偏好

	// 统计数据
	ActionsToday     *int `json:"actions_today" gorm:"column:actions_today;default:0"`
	ActionsThisWeek  *int `json:"actions_this_week" gorm:"column:actions_this_week;default:0"`
	ActionsThisMonth *int `json:"actions_this_month" gorm:"column:actions_this_month;default:0"`
	TotalActions     *int `json:"total_actions" gorm:"column:total_actions;default:0"`

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"`
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 逗号分隔的标签

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (Admin) TableName() string {
	return "t_admins"
}

// 管理员状态常量
const (
	AdminStatusActive    = "active"
	AdminStatusInactive  = "inactive"
	AdminStatusSuspended = "suspended"
	AdminStatusLocked    = "locked"
)

// 在线状态常量
const (
	AdminActiveStatusOnline  = "online"
	AdminActiveStatusOffline = "offline"
	AdminActiveStatusBusy    = "busy"
)

// 管理员角色常量
const (
	AdminRoleSuperAdmin = "super_admin"
	AdminRoleAdmin      = "admin"
	AdminRoleModerator  = "moderator"
	AdminRoleSupport    = "support"
	AdminRoleAnalyst    = "analyst"
)

// 权限常量
const (
	PermissionUserManagement      = "user_management"
	PermissionDriverManagement    = "driver_management"
	PermissionVehicleManagement   = "vehicle_management"
	PermissionOrderManagement     = "order_management"
	PermissionPaymentManagement   = "payment_management"
	PermissionFinancialManagement = "financial_management"
	PermissionSystemConfig        = "system_config"
	PermissionAnalytics           = "analytics"
	PermissionCustomerSupport     = "customer_support"
	PermissionAdminManagement     = "admin_management"
	PermissionAuditLogs           = "audit_logs"
	PermissionEmergencyActions    = "emergency_actions"
)

// 创建新的管理员对象
func NewAdminV2() *Admin {
	return &Admin{
		AdminID: utils.GenerateAdminID(),
		Salt:    utils.GenerateSalt(),
		AdminValues: &AdminValues{
			Status:                utils.StringPtr(AdminStatusActive),
			ActiveStatus:          utils.StringPtr(AdminActiveStatusOffline),
			EmailVerified:         utils.BoolPtr(false),
			PhoneVerified:         utils.BoolPtr(false),
			TwoFactorEnabled:      utils.BoolPtr(false),
			LoginCount:            utils.IntPtr(0),
			FailedAttempts:        utils.IntPtr(0),
			SessionCount:          utils.IntPtr(0),
			MaxConcurrentSessions: utils.IntPtr(3),
			MustChangePassword:    utils.BoolPtr(false),
			Language:              utils.StringPtr("en"),
			Timezone:              utils.StringPtr("UTC"),
			DateFormat:            utils.StringPtr("YYYY-MM-DD"),
			TimeFormat:            utils.StringPtr("24h"),
			ActionsToday:          utils.IntPtr(0),
			ActionsThisWeek:       utils.IntPtr(0),
			ActionsThisMonth:      utils.IntPtr(0),
			TotalActions:          utils.IntPtr(0),
		},
	}
}

// SetValues 更新AdminV2Values中的非nil值
func (a *AdminValues) SetValues(values *AdminValues) {
	if values == nil {
		return
	}

	if values.Username != nil {
		a.Username = values.Username
	}
	if values.Email != nil {
		a.Email = values.Email
	}
	if values.Phone != nil {
		a.Phone = values.Phone
	}
	if values.FirstName != nil {
		a.FirstName = values.FirstName
	}
	if values.LastName != nil {
		a.LastName = values.LastName
	}
	if values.FullName != nil {
		a.FullName = values.FullName
	}
	if values.Role != nil {
		a.Role = values.Role
	}
	if values.Status != nil {
		a.Status = values.Status
	}
	if values.Department != nil {
		a.Department = values.Department
	}
	if values.JobTitle != nil {
		a.JobTitle = values.JobTitle
	}
	if values.Permissions != nil {
		a.Permissions = values.Permissions
	}
	if values.Notes != nil {
		a.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		a.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (a *AdminValues) GetUsername() string {
	if a.Username == nil {
		return ""
	}
	return *a.Username
}

func (a *AdminValues) GetEmail() string {
	if a.Email == nil {
		return ""
	}
	return *a.Email
}

func (a *AdminValues) GetFullName() string {
	if a.FullName == nil {
		return ""
	}
	return *a.FullName
}

func (a *AdminValues) GetRole() string {
	if a.Role == nil {
		return AdminRoleSupport
	}
	return *a.Role
}

func (a *AdminValues) GetStatus() string {
	if a.Status == nil {
		return AdminStatusActive
	}
	return *a.Status
}

func (a *AdminValues) GetActiveStatus() string {
	if a.ActiveStatus == nil {
		return AdminActiveStatusOffline
	}
	return *a.ActiveStatus
}

func (a *AdminValues) GetLoginCount() int {
	if a.LoginCount == nil {
		return 0
	}
	return *a.LoginCount
}

func (a *AdminValues) GetFailedAttempts() int {
	if a.FailedAttempts == nil {
		return 0
	}
	return *a.FailedAttempts
}

func (a *AdminValues) GetSessionCount() int {
	if a.SessionCount == nil {
		return 0
	}
	return *a.SessionCount
}

func (a *AdminValues) GetMaxConcurrentSessions() int {
	if a.MaxConcurrentSessions == nil {
		return 3
	}
	return *a.MaxConcurrentSessions
}

func (a *AdminValues) GetEmailVerified() bool {
	if a.EmailVerified == nil {
		return false
	}
	return *a.EmailVerified
}

func (a *AdminValues) GetTwoFactorEnabled() bool {
	if a.TwoFactorEnabled == nil {
		return false
	}
	return *a.TwoFactorEnabled
}

func (a *AdminValues) GetMustChangePassword() bool {
	if a.MustChangePassword == nil {
		return false
	}
	return *a.MustChangePassword
}

func (a *AdminValues) GetActionsToday() int {
	if a.ActionsToday == nil {
		return 0
	}
	return *a.ActionsToday
}

func (a *AdminValues) GetActionsThisWeek() int {
	if a.ActionsThisWeek == nil {
		return 0
	}
	return *a.ActionsThisWeek
}

func (a *AdminValues) GetActionsThisMonth() int {
	if a.ActionsThisMonth == nil {
		return 0
	}
	return *a.ActionsThisMonth
}

func (a *AdminValues) GetTotalActions() int {
	if a.TotalActions == nil {
		return 0
	}
	return *a.TotalActions
}

// Setter 方法
func (a *AdminValues) SetUsername(username string) *AdminValues {
	a.Username = &username
	return a
}

func (a *AdminValues) SetEmail(email string) *AdminValues {
	a.Email = &email
	return a
}

func (a *AdminValues) SetPhone(phone string) *AdminValues {
	a.Phone = &phone
	return a
}

func (a *AdminValues) SetFullName(firstName, lastName string) *AdminValues {
	a.FirstName = &firstName
	a.LastName = &lastName
	fullName := strings.TrimSpace(firstName + " " + lastName)
	a.FullName = &fullName
	return a
}

func (a *AdminValues) SetRole(role string) *AdminValues {
	a.Role = &role
	return a
}

func (a *AdminValues) SetStatus(status string) *AdminValues {
	a.Status = &status
	return a
}

func (a *AdminValues) SetActiveStatus(status string) *AdminValues {
	a.ActiveStatus = &status
	return a
}

func (a *AdminValues) SetDepartment(department string) *AdminValues {
	a.Department = &department
	return a
}

func (a *AdminValues) SetJobTitle(title string) *AdminValues {
	a.JobTitle = &title
	return a
}

func (a *AdminValues) SetPasswordHash(hash, salt string) *AdminValues {
	a.PasswordHash = &hash
	a.PasswordSalt = &salt
	now := utils.TimeNowMilli()
	a.PasswordChangedAt = &now
	return a
}

func (a *AdminValues) SetEmailVerified(verified bool) *AdminValues {
	a.EmailVerified = &verified
	return a
}

func (a *AdminValues) SetTwoFactorEnabled(enabled bool) *AdminValues {
	a.TwoFactorEnabled = &enabled
	return a
}

func (a *AdminValues) SetTwoFactorSecret(secret string) *AdminValues {
	a.TwoFactorSecret = &secret
	return a
}

func (a *AdminValues) SetMustChangePassword(must bool) *AdminValues {
	a.MustChangePassword = &must
	return a
}

// 业务方法
func (a *Admin) IsActive() bool {
	return a.GetStatus() == AdminStatusActive
}

func (a *Admin) IsInactive() bool {
	return a.GetStatus() == AdminStatusInactive
}

func (a *Admin) IsSuspended() bool {
	return a.GetStatus() == AdminStatusSuspended
}

func (a *Admin) IsLocked() bool {
	return a.GetStatus() == AdminStatusLocked
}

func (a *Admin) IsOnline() bool {
	return a.GetActiveStatus() == AdminActiveStatusOnline
}

func (a *Admin) IsOffline() bool {
	return a.GetActiveStatus() == AdminActiveStatusOffline
}

func (a *Admin) IsBusy() bool {
	return a.GetActiveStatus() == AdminActiveStatusBusy
}

func (a *Admin) IsSuperAdmin() bool {
	return a.GetRole() == AdminRoleSuperAdmin
}

func (a *Admin) IsAdmin() bool {
	return a.GetRole() == AdminRoleAdmin
}

func (a *Admin) CanLogin() bool {
	return a.IsActive() && !a.IsAccountLocked()
}

func (a *Admin) IsAccountLocked() bool {
	if a.LockedUntil == nil {
		return false
	}
	return utils.TimeNowMilli() < *a.LockedUntil
}

func (a *Admin) HasExceededSessionLimit() bool {
	return a.GetSessionCount() >= a.GetMaxConcurrentSessions()
}

func (a *Admin) ShouldForcePasswordChange() bool {
	return a.GetMustChangePassword()
}

// 权限相关方法
func (a *AdminValues) HasPermission(permission string) bool {
	if a.Permissions == nil {
		return false
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return false
	}

	return slices.Contains(permissions, permission)
}

func (a *AdminValues) AddPermission(permission string) error {
	var permissions []string
	if a.Permissions != nil {
		if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
			return fmt.Errorf("failed to parse existing permissions: %v", err)
		}
	}

	// 避免重复添加
	for _, p := range permissions {
		if p == permission {
			return nil
		}
	}

	permissions = append(permissions, permission)
	permissionsJSON, err := utils.ToJSON(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *AdminValues) RemovePermission(permission string) error {
	if a.Permissions == nil {
		return nil
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return fmt.Errorf("failed to parse existing permissions: %v", err)
	}

	var newPermissions []string
	for _, p := range permissions {
		if p != permission {
			newPermissions = append(newPermissions, p)
		}
	}

	permissionsJSON, err := utils.ToJSON(newPermissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *AdminValues) SetPermissions(permissions []string) error {
	permissionsJSON, err := utils.ToJSON(permissions)
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %v", err)
	}

	a.Permissions = &permissionsJSON
	return nil
}

func (a *AdminValues) GetPermissions() []string {
	if a.Permissions == nil {
		return []string{}
	}

	var permissions []string
	if err := utils.FromJSON(*a.Permissions, &permissions); err != nil {
		return []string{}
	}

	return permissions
}

// 登录相关方法
func (a *AdminValues) RecordLogin(ip string, sessionID string) *AdminValues {
	now := utils.TimeNowMilli()
	a.LastLoginAt = &now
	a.LastLoginIP = &ip
	a.CurrentSessionID = &sessionID

	loginCount := a.GetLoginCount() + 1
	a.LoginCount = &loginCount

	sessionCount := a.GetSessionCount() + 1
	a.SessionCount = &sessionCount

	// 重置失败次数
	a.FailedAttempts = utils.IntPtr(0)

	return a
}

func (a *AdminValues) RecordFailedLogin() *AdminValues {
	now := utils.TimeNowMilli()
	a.LastFailedAt = &now

	attempts := a.GetFailedAttempts() + 1
	a.FailedAttempts = &attempts

	// 如果失败次数过多，锁定账户
	if attempts >= 5 {
		lockUntil := now + (30 * 60 * 1000) // 锁定30分钟
		a.LockedUntil = &lockUntil
	}

	return a
}

func (a *AdminValues) Logout() *AdminValues {
	a.CurrentSessionID = nil
	sessionCount := a.GetSessionCount() - 1
	if sessionCount < 0 {
		sessionCount = 0
	}
	a.SessionCount = &sessionCount
	a.SetActiveStatus(AdminActiveStatusOffline)

	return a
}

func (a *AdminValues) UnlockAccount() *AdminValues {
	a.LockedUntil = nil
	a.FailedAttempts = utils.IntPtr(0)
	return a
}

// 状态管理方法
func (a *AdminValues) Activate(adminID string) *AdminValues {
	a.SetStatus(AdminStatusActive)
	a.LastUpdatedBy = &adminID
	return a
}

func (a *AdminValues) Deactivate(adminID string) *AdminValues {
	a.SetStatus(AdminStatusInactive)
	a.LastUpdatedBy = &adminID
	return a
}

func (a *AdminValues) Suspend(adminID, reason string) *AdminValues {
	a.SetStatus(AdminStatusSuspended)
	a.SuspendedBy = &adminID
	a.SuspensionReason = &reason
	now := utils.TimeNowMilli()
	a.SuspendedAt = &now
	a.LastUpdatedBy = &adminID
	return a
}

func (a *AdminValues) Unsuspend(adminID string) *AdminValues {
	a.SetStatus(AdminStatusActive)
	a.SuspendedBy = nil
	a.SuspendedAt = nil
	a.SuspensionReason = nil
	a.LastUpdatedBy = &adminID
	return a
}

func (a *AdminValues) Lock(adminID string) *AdminValues {
	a.SetStatus(AdminStatusLocked)
	now := utils.TimeNowMilli()
	a.LockedUntil = &now
	a.LastUpdatedBy = &adminID
	return a
}

// 审批方法
func (a *AdminValues) Approve(adminID string) *AdminValues {
	a.ApprovedBy = &adminID
	now := utils.TimeNowMilli()
	a.ApprovedAt = &now
	a.SetStatus(AdminStatusActive)
	return a
}

// 统计更新方法
func (a *AdminValues) IncrementActions() *AdminValues {
	today := a.GetActionsToday() + 1
	week := a.GetActionsThisWeek() + 1
	month := a.GetActionsThisMonth() + 1
	total := a.GetTotalActions() + 1

	a.ActionsToday = &today
	a.ActionsThisWeek = &week
	a.ActionsThisMonth = &month
	a.TotalActions = &total

	return a
}

// 工作时间设置
func (a *AdminValues) SetWorkingHours(workingHours map[string]interface{}) error {
	workingHoursJSON, err := utils.ToJSON(workingHours)
	if err != nil {
		return fmt.Errorf("failed to marshal working hours: %v", err)
	}

	a.WorkingHours = &workingHoursJSON
	return nil
}

func (a *AdminValues) SetWorkingAreas(areas []string) error {
	areasJSON, err := utils.ToJSON(areas)
	if err != nil {
		return fmt.Errorf("failed to marshal working areas: %v", err)
	}

	a.WorkingAreas = &areasJSON
	return nil
}

// 偏好设置
func (a *AdminValues) SetPreferences(preferences map[string]interface{}) error {
	preferencesJSON, err := utils.ToJSON(preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %v", err)
	}

	a.Preferences = &preferencesJSON
	return nil
}

// IP白名单管理
func (a *AdminValues) SetAllowedIPs(ips []string) *AdminValues {
	allowedIPs := strings.Join(ips, ",")
	a.AllowedIPs = &allowedIPs
	return a
}

func (a *AdminValues) IsIPAllowed(ip string) bool {
	if a.AllowedIPs == nil || *a.AllowedIPs == "" {
		return true // 如果没有设置白名单，允许所有IP
	}

	allowedIPs := strings.Split(*a.AllowedIPs, ",")
	for _, allowedIP := range allowedIPs {
		if strings.TrimSpace(allowedIP) == ip {
			return true
		}
	}

	return false
}

// 角色级权限设置
func GetRolePermissions(role string) []string {
	switch role {
	case AdminRoleSuperAdmin:
		return []string{
			PermissionUserManagement,
			PermissionDriverManagement,
			PermissionVehicleManagement,
			PermissionOrderManagement,
			PermissionPaymentManagement,
			PermissionFinancialManagement,
			PermissionSystemConfig,
			PermissionAnalytics,
			PermissionCustomerSupport,
			PermissionAdminManagement,
			PermissionAuditLogs,
			PermissionEmergencyActions,
		}
	case AdminRoleAdmin:
		return []string{
			PermissionUserManagement,
			PermissionDriverManagement,
			PermissionVehicleManagement,
			PermissionOrderManagement,
			PermissionPaymentManagement,
			PermissionAnalytics,
			PermissionCustomerSupport,
			PermissionAuditLogs,
		}
	case AdminRoleModerator:
		return []string{
			PermissionUserManagement,
			PermissionDriverManagement,
			PermissionOrderManagement,
			PermissionCustomerSupport,
		}
	case AdminRoleSupport:
		return []string{
			PermissionCustomerSupport,
			PermissionOrderManagement,
		}
	case AdminRoleAnalyst:
		return []string{
			PermissionAnalytics,
		}
	default:
		return []string{}
	}
}

// 创建具有指定角色的管理员
func NewAdminV2WithRole(username, email, role string) *Admin {
	admin := NewAdminV2()
	admin.SetUsername(username).
		SetEmail(email).
		SetRole(role)

	// 设置角色对应的权限
	rolePermissions := GetRolePermissions(role)
	admin.SetPermissions(rolePermissions)

	return admin
}

func GetAdminByID(adminID string) *Admin {
	var admin Admin
	if err := GetDB().Where("id = ?", adminID).First(&admin).Error; err != nil {
		return nil
	}
	return &admin
}

// NewAdminInfoV2FromModel 从模型创建管理员信息
func (admin *Admin) Protocol() protocol.Admin {
	if admin == nil {
		return protocol.Admin{}
	}

	adminInfo := protocol.Admin{
		AdminID:   admin.AdminID,
		Username:  admin.GetUsername(),
		Email:     admin.GetEmail(),
		FullName:  admin.GetFullName(),
		CreatedAt: admin.CreatedAt,
	}

	// 处理可能为空的指针字段
	if admin.Role != nil {
		adminInfo.Role = *admin.Role
	}

	if admin.Department != nil {
		adminInfo.Department = *admin.Department
	}

	if admin.Status != nil {
		adminInfo.Status = *admin.Status
	}

	if admin.ActiveStatus != nil {
		adminInfo.ActiveStatus = *admin.ActiveStatus
	}

	if admin.LastLoginAt != nil {
		adminInfo.LastLoginAt = admin.LastLoginAt
	}

	return adminInfo
}

func (a *AdminValues) GetFirstName() string {
	if a.FirstName == nil {
		return ""
	}
	return *a.FirstName
}

func (a *AdminValues) GetLastName() string {
	if a.LastName == nil {
		return ""
	}
	return *a.LastName
}
