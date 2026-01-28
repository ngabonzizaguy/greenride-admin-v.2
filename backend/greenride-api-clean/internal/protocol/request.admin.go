package protocol

// ============================================================================
// 管理员相关请求结构体
// ============================================================================

// AdminLoginRequest 管理员登录请求结构体
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// AdminChangePasswordRequest 管理员修改密码请求结构体
type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// AdminResetPasswordRequest 管理员重置密码请求结构体
type AdminResetPasswordRequest struct {
	TargetAdminID string `json:"target_admin_id" binding:"required"` // 目标管理员ID
	NewPassword   string `json:"new_password" binding:"required"`    // 新密码
}

// AdminCreateRequest 创建管理员请求结构体
type AdminCreateRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email      string `json:"email" binding:"required,email"`           // 邮箱
	Password   string `json:"password" binding:"required,min=6"`        // 密码
	Role       string `json:"role" binding:"required"`                  // 角色
	Department string `json:"department"`                               // 部门
	JobTitle   string `json:"job_title"`                                // 职位
	FirstName  string `json:"first_name"`                               // 名字
	LastName   string `json:"last_name"`                                // 姓氏
}

// SearchRequest 统一的搜索请求结构体
type SearchRequest struct {
	Keyword  string `json:"keyword,omitempty"`   // 搜索关键字
	Page     int    `json:"page,omitempty"`      // 页码，默认1
	Limit    int    `json:"limit,omitempty"`     // 每页数量，默认10
	UserType string `json:"user_type,omitempty"` // 用户类型 (仅用户搜索)

	// 司机相关搜索参数
	Status          string   `json:"status,omitempty"`            // 用户/司机状态
	OnlineStatus    string   `json:"online_status,omitempty"`     // 在线状态 (仅司机)
	IsEmailVerified *bool    `json:"is_email_verified,omitempty"` // 邮箱是否验证
	IsPhoneVerified *bool    `json:"is_phone_verified,omitempty"` // 手机是否验证
	IsActive        *bool    `json:"is_active,omitempty"`         // 是否活跃
	MinDriverScore  *float64 `json:"min_driver_score,omitempty"`  // 最低司机评分 (仅司机)
	MaxDriverScore  *float64 `json:"max_driver_score,omitempty"`  // 最高司机评分 (仅司机)
	MinTotalRides   *int     `json:"min_total_rides,omitempty"`   // 最少接单数 (仅司机)
}

// VehicleSearchRequest 车辆搜索请求结构体
type VehicleSearchRequest struct {
	Keyword    string `json:"keyword,omitempty"`     // 搜索关键字
	Page       int    `json:"page,omitempty"`        // 页码，默认1
	Limit      int    `json:"limit,omitempty"`       // 每页数量，默认10
	TypeID     string `json:"type_id,omitempty"`     // 车辆类型ID（关联VehicleType表）
	Category   string `json:"category,omitempty"`    // 车辆分类
	Level      string `json:"level,omitempty"`       // 服务级别
	Status     string `json:"status,omitempty"`      // 车辆状态
	DriverID   string `json:"driver_id,omitempty"`   // 司机ID
	IsVerified *bool  `json:"is_verified,omitempty"` // 是否已验证
	IsActive   *bool  `json:"is_active,omitempty"`   // 是否活跃
	YearFrom   *int   `json:"year_from,omitempty"`   // 年份范围开始
	YearTo     *int   `json:"year_to,omitempty"`     // 年份范围结束
}

// DriverSearchRequest 司机搜索请求结构体
type DriverSearchRequest struct {
	Keyword         string   `json:"keyword,omitempty"`           // 搜索关键字
	Page            int      `json:"page,omitempty"`              // 页码，默认1
	Limit           int      `json:"limit,omitempty"`             // 每页数量，默认10
	Status          string   `json:"status,omitempty"`            // 司机状态
	OnlineStatus    string   `json:"online_status,omitempty"`     // 在线状态
	IsEmailVerified *bool    `json:"is_email_verified,omitempty"` // 邮箱是否验证
	IsPhoneVerified *bool    `json:"is_phone_verified,omitempty"` // 手机是否验证
	IsActive        *bool    `json:"is_active,omitempty"`         // 是否活跃
	MinDriverScore  *float64 `json:"min_driver_score,omitempty"`  // 最低司机评分
	MaxDriverScore  *float64 `json:"max_driver_score,omitempty"`  // 最高司机评分
	MinTotalRides   *int     `json:"min_total_rides,omitempty"`   // 最少接单数
}

// OrderSearchRequest 订单搜索请求结构体
type OrderSearchRequest struct {
	Keyword       string   `json:"keyword,omitempty"`        // 搜索关键字
	Page          int      `json:"page,omitempty"`           // 页码，默认1
	Limit         int      `json:"limit,omitempty"`          // 每页数量，默认10
	OrderID       string   `json:"order_id,omitempty"`       // 订单ID
	OrderType     string   `json:"order_type,omitempty"`     // 订单类型
	Status        string   `json:"status,omitempty"`         // 订单状态
	PaymentStatus string   `json:"payment_status,omitempty"` // 支付状态
	UserID        string   `json:"user_id,omitempty"`        // 用户ID
	ProviderID    string   `json:"provider_id,omitempty"`    // 服务提供者ID
	StartDate     *int64   `json:"start_date,omitempty"`     // 开始日期 (时间戳毫秒)
	EndDate       *int64   `json:"end_date,omitempty"`       // 结束日期 (时间戳毫秒)
	MinAmount     *float64 `json:"min_amount,omitempty"`     // 最小金额
	MaxAmount     *float64 `json:"max_amount,omitempty"`     // 最大金额
}

// PromotionSearchRequest 优惠券搜索请求结构体
type PromotionSearchRequest struct {
	Keyword   string `json:"keyword,omitempty"`    // 搜索关键字
	Page      int    `json:"page,omitempty"`       // 页码，默认1
	Limit     int    `json:"limit,omitempty"`      // 每页数量，默认10
	Status    string `json:"status,omitempty"`     // 优惠券状态
	Type      string `json:"type,omitempty"`       // 优惠券类型
	IsActive  *bool  `json:"is_active,omitempty"`  // 是否活跃
	ValidFrom string `json:"valid_from,omitempty"` // 有效期开始
	ValidTo   string `json:"valid_to,omitempty"`   // 有效期结束
}

// IDRequest 通用ID请求结构体
type IDRequest struct {
	ID string `json:"id" binding:"required"` // ID
}

// UserIDRequest 用户ID请求结构体
type UserIDRequest struct {
	UserID string `json:"user_id" binding:"required"` // 用户ID
}

// UserRidesRequest 用户行程查询请求结构体（兼容司机和乘客）
type UserRidesRequest struct {
	UserID    string `json:"user_id"`              // 用户ID
	UserType  string `json:"user_type,omitempty"`  // 用户类型：user(乘客), driver(司机)，空则自动识别
	Page      int    `json:"page,omitempty"`       // 页码，默认1
	Limit     int    `json:"limit,omitempty"`      // 每页数量，默认10
	Status    string `json:"status,omitempty"`     // 订单状态过滤
	OrderType string `json:"order_type,omitempty"` // 订单类型过滤
	StartDate *int64 `json:"start_date,omitempty"` // 开始日期过滤 (时间戳毫秒)
	EndDate   *int64 `json:"end_date,omitempty"`   // 结束日期过滤 (时间戳毫秒)
}

// UserDispatchsRequest 用户派单记录查询请求结构体（兼容司机和乘客）
type UserDispatchsRequest struct {
	UserID    string `json:"user_id" binding:"required"` // 用户ID
	UserType  string `json:"user_type,omitempty"`        // 用户类型：user(乘客), driver(司机)，空则自动识别
	Page      int    `json:"page,omitempty"`             // 页码，默认1
	Limit     int    `json:"limit,omitempty"`            // 每页数量，默认10
	Status    string `json:"status,omitempty"`           // 订单状态过滤
	OrderType string `json:"order_type,omitempty"`       // 订单类型过滤
	StartDate *int64 `json:"start_date,omitempty"`       // 开始日期过滤 (时间戳毫秒)
	EndDate   *int64 `json:"end_date,omitempty"`         // 结束日期过滤 (时间戳毫秒)
}

// UserStatusUpdateRequest 用户状态更新请求结构体
type UserStatusUpdateRequest struct {
	UserID   string `json:"user_id" binding:"required"`   // 用户ID
	Status   string `json:"status" binding:"required"`    // active, inactive, suspended, banned
	IsActive *bool  `json:"is_active" binding:"required"` // true, false
}

// UserUpdateRequest 用户信息更新请求结构体
type UserUpdateRequest struct {
	UserID      string   `json:"user_id" binding:"required"` // 用户ID
	UserType    *string  `json:"user_type,omitempty"`        // user, driver
	Email       *string  `json:"email,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	CountryCode *string  `json:"country_code,omitempty"`
	Username    *string  `json:"username,omitempty"`
	DisplayName *string  `json:"display_name,omitempty"`
	FirstName   *string  `json:"first_name,omitempty"`
	LastName    *string  `json:"last_name,omitempty"`
	Avatar      *string  `json:"avatar,omitempty"`
	Gender      *string  `json:"gender,omitempty"` // male, female, other
	Birthday    *int64   `json:"birthday,omitempty"`
	Language    *string  `json:"language,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
	Address     *string  `json:"address,omitempty"`
	City        *string  `json:"city,omitempty"`
	State       *string  `json:"state,omitempty"`
	Country     *string  `json:"country,omitempty"`
	PostalCode  *string  `json:"postal_code,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	Status      *string  `json:"status,omitempty"` // active, inactive, suspended, banned
	IsActive    *bool    `json:"is_active,omitempty"`
}

// UserVerifyRequest 用户认证请求结构体（替代DriverVerifyRequest，统一处理用户认证）
type UserVerifyRequest struct {
	UserID          string `json:"user_id" binding:"required"`  // 用户ID
	IsEmailVerified *bool  `json:"is_email_verified,omitempty"` // 邮箱是否验证
	IsPhoneVerified *bool  `json:"is_phone_verified,omitempty"` // 手机是否验证
	VerifiedBy      string `json:"verified_by,omitempty"`       // 验证者ID
}

// StatusUpdateRequest 状态更新请求结构体
type StatusUpdateRequest struct {
	ID     string `json:"id" binding:"required"`     // ID
	Status string `json:"status" binding:"required"` // 状态
}

// OrderCancelRequest 取消订单请求结构体
type OrderCancelRequest struct {
	ID     string `json:"id" binding:"required"` // 订单ID
	Reason string `json:"reason,omitempty"`      // 取消原因
}

// AdminOrderCancelRequest 管理员取消订单请求结构体
type AdminOrderCancelRequest struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
	Reason  string `json:"reason" binding:"required"`   // 取消原因
}

// AdminOrderStatusUpdateRequest 管理员订单状态更新请求结构体
type AdminOrderStatusUpdateRequest struct {
	Status        string `json:"status" binding:"required"` // 订单状态
	PaymentStatus string `json:"payment_status,omitempty"`  // 支付状态
	Notes         string `json:"notes,omitempty"`           // 备注信息
}

// ApprovePromotionRequest 优惠券审批请求结构体
type ApprovePromotionRequest struct {
	UserID      string `json:"user_id" `                        // 用户ID
	PromotionID string `json:"promotion_id" binding:"required"` // 优惠券ID
	Notes       string `json:"notes,omitempty"`                 // 审批备注
}

// DeletePromotionRequest 优惠券删除请求结构体
type DeletePromotionRequest struct {
	UserID      string `json:"user_id,omitempty"`               // 用户ID（后端自动填充）
	PromotionID string `json:"promotion_id" binding:"required"` // 优惠券ID
}

// CreatePromotionRequest 优惠券创建请求结构体
type CreatePromotionRequest struct {
	UserID            string   `json:"user_id"` // 创建者ID
	Code              string   `json:"code" binding:"required,min=3,max=50"`
	Title             string   `json:"title" binding:"required,max=255"`
	Description       string   `json:"description,omitempty"`
	Status            string   `json:"status,omitempty" binding:"omitempty,oneof=active inactive expired suspended deleted"`
	PromotionType     string   `json:"promotion_type" binding:"required,oneof=discount cashback free_ride upgrade points"`
	DiscountType      string   `json:"discount_type" binding:"required,oneof=percentage fixed_amount free"`
	DiscountValue     float64  `json:"discount_value" binding:"required,min=0"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty" binding:"omitempty,min=0"`
	MinOrderAmount    *float64 `json:"min_order_amount,omitempty" binding:"omitempty,min=0"`
	UsageLimit        *int     `json:"usage_limit,omitempty" binding:"omitempty,min=1"`
	UserUsageLimit    *int     `json:"user_usage_limit,omitempty" binding:"omitempty,min=1"`
	StartDate         *int64   `json:"start_date,omitempty"`
	EndDate           *int64   `json:"end_date,omitempty"`
	ValidDays         *int     `json:"valid_days,omitempty" binding:"omitempty,min=1"`
	TargetUserType    string   `json:"target_user_type,omitempty" binding:"omitempty,oneof=all new_user old_user vip driver"`
	ValidCities       string   `json:"valid_cities,omitempty"`
	ValidVehicleTypes string   `json:"valid_vehicle_types,omitempty"`
	IsPublic          *bool    `json:"is_public,omitempty"`
	Priority          *int     `json:"priority,omitempty" binding:"omitempty,min=0"`
	Tags              string   `json:"tags,omitempty"`
}

// UpdatePromotionRequest 优惠券更新请求结构体
type UpdatePromotionRequest struct {
	UserID            string   `json:"user_id" ` // 更新者ID
	PromotionID       string   `json:"promotion_id" binding:"required"`
	Title             *string  `json:"title,omitempty" binding:"omitempty,max=255"`
	Description       *string  `json:"description,omitempty"`
	DiscountValue     *float64 `json:"discount_value,omitempty" binding:"omitempty,min=0"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty" binding:"omitempty,min=0"`
	MinOrderAmount    *float64 `json:"min_order_amount,omitempty" binding:"omitempty,min=0"`
	UsageLimit        *int     `json:"usage_limit,omitempty" binding:"omitempty,min=1"`
	UserUsageLimit    *int     `json:"user_usage_limit,omitempty" binding:"omitempty,min=1"`
	StartDate         *int64   `json:"start_date,omitempty"`
	EndDate           *int64   `json:"end_date,omitempty"`
	ValidCities       *string  `json:"valid_cities,omitempty"`
	ValidVehicleTypes *string  `json:"valid_vehicle_types,omitempty"`
	Priority          *int     `json:"priority,omitempty" binding:"omitempty,min=0"`
	Tags              *string  `json:"tags,omitempty"`
}

// UpdatePromotionStatusRequest 优惠券状态更新请求结构体
type UpdatePromotionStatusRequest struct {
	UserID      string `json:"user_id"`
	PromotionID string `json:"promotion_id" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=active inactive expired suspended deleted"`
	IsActive    *bool  `json:"is_active,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

// PromotionDetailRequest 优惠券详情请求结构体
type PromotionDetailRequest struct {
	PromotionID string `json:"promotion_id" binding:"required"` // 优惠券ID
}

// PromotionUsageRequest 优惠券使用统计请求结构体
type PromotionUsageRequest struct {
	PromotionID string `json:"promotion_id" binding:"required"` // 优惠券ID
	StartDate   *int64 `json:"start_date,omitempty"`            // 开始时间
	EndDate     *int64 `json:"end_date,omitempty"`              // 结束时间
}

// AdminUpdateRequest 管理员更新请求结构体
type AdminUpdateRequest struct {
	ID         string  `json:"id" binding:"required"` // 管理员ID
	Username   *string `json:"username,omitempty"`    // 用户名
	Email      *string `json:"email,omitempty"`       // 邮箱
	Role       *string `json:"role,omitempty"`        // 角色
	Department *string `json:"department,omitempty"`  // 部门
	JobTitle   *string `json:"job_title,omitempty"`   // 职位
	FirstName  *string `json:"first_name,omitempty"`  // 名字
	LastName   *string `json:"last_name,omitempty"`   // 姓氏
	Status     *string `json:"status,omitempty"`      // 状态
}

// SearchPriceRuleRequest 价格规则列表请求结构体
type SearchPriceRuleRequest struct {
	Page     int    `json:"page,omitempty"`      // 页码，默认1
	PageSize int    `json:"page_size,omitempty"` // 每页数量，默认20
	Category string `json:"category,omitempty"`  // 规则分类
	Status   string `json:"status,omitempty"`    // 规则状态
	Keyword  string `json:"keyword,omitempty"`   // 搜索关键字
}

// PriceRuleRequest 价格规则基本请求结构体（仅包含RuleID）
type PriceRuleRequest struct {
	UserID string `json:"user_id"`                    // 操作者用户ID
	RuleID string `json:"rule_id" binding:"required"` // 价格规则ID (对应 models.PriceRule.RuleID)
}

// UpdatePriceRuleStatusRequest 价格规则状态更新请求结构体
type UpdatePriceRuleStatusRequest struct {
	UserID string `json:"user_id"`                    // 操作者用户ID
	RuleID string `json:"rule_id" binding:"required"` // 价格规则ID (对应 models.PriceRule.RuleID)
	Status string `json:"status" binding:"required"`  // 新状态
}

// AdminOrderEstimateRequest 管理员订单预估请求结构体
type AdminOrderEstimateRequest struct {
	*EstimateRequest        // 直接嵌入EstimateRequest，继承所有字段
	AdminReason      string `json:"admin_reason,omitempty"` // 管理员操作原因（管理员特有字段）
}

// AdminCreateOrderRequest 管理员创建订单请求结构体
type AdminCreateOrderRequest struct {
	AdminID     string `json:"admin_id,omitempty"`          // 管理员ID（从上下文获取）
	UserID      string `json:"user_id" binding:"required"`  // 用户ID
	PriceID     string `json:"price_id" binding:"required"` // 价格快照ID
	Notes       string `json:"notes,omitempty"`             // 订单备注
	AdminReason string `json:"admin_reason,omitempty"`      // 管理员创建原因
	ProviderID  string `json:"provider_id,omitempty"`       // 可选：手动指定司机（provider_id）
}

// AdminCreateUserRequest 管理员创建用户请求结构体
type AdminCreateUserRequest struct {
	Phone     string `json:"phone" binding:"required"` // 手机号
	Username  string `json:"username,omitempty"`       // 用户名
	Email     string `json:"email,omitempty"`          // 邮箱
	FirstName string `json:"first_name,omitempty"`     // 名字
	LastName  string `json:"last_name,omitempty"`      // 姓氏
	UserType  string `json:"user_type,omitempty"`      // 用户类型，默认passenger
}

// AdminSendNotificationRequest 管理员发送通知请求结构体
type AdminSendNotificationRequest struct {
	Audience    string `json:"audience" binding:"required,oneof=all drivers users"` // 受众：all, drivers, users
	Type        string `json:"type" binding:"required"`                             // 通知类型
	Category    string `json:"category" binding:"required"`                         // 通知分类
	Title       string `json:"title" binding:"required"`                            // 通知标题
	Content     string `json:"content" binding:"required"`                          // 通知内容
	Summary     string `json:"summary,omitempty"`                                   // 通知摘要
	ScheduledAt *int64 `json:"scheduled_at,omitempty"`                              // 计划发送时间（时间戳毫秒）
}

// AdminNotificationSearchRequest 管理员通知搜索请求结构体
type AdminNotificationSearchRequest struct {
	Keyword  string `json:"keyword,omitempty"`   // 搜索关键字
	Page     int    `json:"page,omitempty"`      // 页码，默认1
	Limit    int    `json:"limit,omitempty"`     // 每页数量，默认10
	UserType string `json:"user_type,omitempty"` // 用户类型过滤
	Type     string `json:"type,omitempty"`      // 通知类型过滤
	Status   string `json:"status,omitempty"`    // 状态过滤
}

// NotificationIDRequest 通知ID请求结构体
type NotificationIDRequest struct {
	NotificationID string `json:"notification_id" binding:"required"` // 通知ID
}
