package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// UserPaymentMethod 用户支付方式关联表 - 用户与支付方式的多对多关联关系
type UserPaymentMethod struct {
	ID                  int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserPaymentMethodID string `json:"user_payment_method_id" gorm:"column:user_payment_method_id;type:varchar(64);uniqueIndex"`
	Salt                string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*UserPaymentMethodValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type UserPaymentMethodValues struct {
	// 关联信息
	UserID          *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`                     // 用户ID
	PaymentMethodID *string `json:"payment_method_id" gorm:"column:payment_method_id;type:varchar(64);index"` // 支付方式ID

	// 关联状态
	Status      *string `json:"status" gorm:"column:status;type:varchar(30);index;default:'active'"` // active, inactive, deleted, suspended
	IsDefault   *bool   `json:"is_default" gorm:"column:is_default;default:false"`                   // 是否为默认支付方式
	IsPrimary   *bool   `json:"is_primary" gorm:"column:is_primary;default:false"`                   // 是否为主要支付方式
	IsPreferred *bool   `json:"is_preferred" gorm:"column:is_preferred;default:false"`               // 是否为偏好支付方式

	// 验证状态
	IsVerified         *bool   `json:"is_verified" gorm:"column:is_verified;default:false"`                    // 是否验证
	VerifiedAt         *int64  `json:"verified_at" gorm:"column:verified_at"`                                  // 验证时间
	VerificationMethod *string `json:"verification_method" gorm:"column:verification_method;type:varchar(50)"` // manual, automatic, otp, biometric

	// 绑定信息
	BoundAt     *int64  `json:"bound_at" gorm:"column:bound_at"`                           // 绑定时间
	BoundBy     *string `json:"bound_by" gorm:"column:bound_by;type:varchar(64)"`          // 绑定操作者ID（用户本人或管理员）
	BoundMethod *string `json:"bound_method" gorm:"column:bound_method;type:varchar(50)"`  // manual, automatic, import, api
	BoundSource *string `json:"bound_source" gorm:"column:bound_source;type:varchar(100)"` // 绑定来源 app, web, admin_panel, third_party

	// 支付限制和配置
	DailyLimit             *float64 `json:"daily_limit" gorm:"column:daily_limit;type:decimal(12,2)"`                           // 日支付限额
	MonthlyLimit           *float64 `json:"monthly_limit" gorm:"column:monthly_limit;type:decimal(12,2)"`                       // 月支付限额
	SingleTransactionLimit *float64 `json:"single_transaction_limit" gorm:"column:single_transaction_limit;type:decimal(12,2)"` // 单笔限额

	// 支付类型配置
	AllowedPaymentTypes *string `json:"allowed_payment_types" gorm:"column:allowed_payment_types;type:json"` // JSON数组：允许的支付类型
	RestrictedServices  *string `json:"restricted_services" gorm:"column:restricted_services;type:json"`     // JSON数组：限制的服务类型

	// 使用统计
	TotalTransactions      *int     `json:"total_transactions" gorm:"column:total_transactions;type:int;default:0"`           // 总交易次数
	TotalAmount            *float64 `json:"total_amount" gorm:"column:total_amount;type:decimal(12,2);default:0.00"`          // 总交易金额
	SuccessfulTransactions *int     `json:"successful_transactions" gorm:"column:successful_transactions;type:int;default:0"` // 成功交易次数
	FailedTransactions     *int     `json:"failed_transactions" gorm:"column:failed_transactions;type:int;default:0"`         // 失败交易次数
	LastUsedAt             *int64   `json:"last_used_at" gorm:"column:last_used_at"`                                          // 最后使用时间
	FirstUsedAt            *int64   `json:"first_used_at" gorm:"column:first_used_at"`                                        // 首次使用时间

	// 成功率和可靠性
	SuccessRate         *float64 `json:"success_rate" gorm:"column:success_rate;type:decimal(5,2);default:0.00"`           // 成功率百分比
	AverageResponseTime *int     `json:"average_response_time" gorm:"column:average_response_time;type:int;default:0"`     // 平均响应时间(毫秒)
	ReliabilityScore    *float64 `json:"reliability_score" gorm:"column:reliability_score;type:decimal(3,2);default:0.00"` // 可靠性评分(0-5.0)

	// 安全相关
	SecurityLevel     *string `json:"security_level" gorm:"column:security_level;type:varchar(20);default:'medium'"` // low, medium, high, maximum
	RequiresBiometric *bool   `json:"requires_biometric" gorm:"column:requires_biometric;default:false"`             // 是否需要生物识别
	RequiresPin       *bool   `json:"requires_pin" gorm:"column:requires_pin;default:false"`                         // 是否需要PIN
	RequiresOTP       *bool   `json:"requires_otp" gorm:"column:requires_otp;default:false"`                         // 是否需要OTP

	// 风险评估
	RiskLevel       *string  `json:"risk_level" gorm:"column:risk_level;type:varchar(20);default:'low'"`   // low, medium, high
	FraudScore      *float64 `json:"fraud_score" gorm:"column:fraud_score;type:decimal(5,2);default:0.00"` // 欺诈风险评分(0-100)
	IsBlacklisted   *bool    `json:"is_blacklisted" gorm:"column:is_blacklisted;default:false"`            // 是否在黑名单
	BlacklistReason *string  `json:"blacklist_reason" gorm:"column:blacklist_reason;type:varchar(255)"`    // 黑名单原因

	// 地理和设备限制
	AllowedCountries    *string `json:"allowed_countries" gorm:"column:allowed_countries;type:json"`           // JSON数组：允许的国家
	RestrictedCountries *string `json:"restricted_countries" gorm:"column:restricted_countries;type:json"`     // JSON数组：限制的国家
	AllowedDevices      *string `json:"allowed_devices" gorm:"column:allowed_devices;type:json"`               // JSON数组：允许的设备
	DeviceFingerprint   *string `json:"device_fingerprint" gorm:"column:device_fingerprint;type:varchar(255)"` // 设备指纹

	// 时间限制
	TimeRestrictions *string `json:"time_restrictions" gorm:"column:time_restrictions;type:json"` // JSON对象：时间限制配置
	ValidFrom        *int64  `json:"valid_from" gorm:"column:valid_from"`                         // 有效期开始时间
	ValidUntil       *int64  `json:"valid_until" gorm:"column:valid_until"`                       // 有效期结束时间

	// 自动充值配置
	AutoReloadEnabled   *bool    `json:"auto_reload_enabled" gorm:"column:auto_reload_enabled;default:false"`         // 是否启用自动充值
	AutoReloadThreshold *float64 `json:"auto_reload_threshold" gorm:"column:auto_reload_threshold;type:decimal(8,2)"` // 自动充值触发阈值
	AutoReloadAmount    *float64 `json:"auto_reload_amount" gorm:"column:auto_reload_amount;type:decimal(8,2)"`       // 自动充值金额

	// 通知配置
	NotifyOnSuccess      *bool    `json:"notify_on_success" gorm:"column:notify_on_success;default:true"`                // 成功时通知
	NotifyOnFailure      *bool    `json:"notify_on_failure" gorm:"column:notify_on_failure;default:true"`                // 失败时通知
	NotifyOnLargeAmount  *bool    `json:"notify_on_large_amount" gorm:"column:notify_on_large_amount;default:true"`      // 大额交易通知
	LargeAmountThreshold *float64 `json:"large_amount_threshold" gorm:"column:large_amount_threshold;type:decimal(8,2)"` // 大额交易阈值

	// 费用配置
	FeePercentage *float64 `json:"fee_percentage" gorm:"column:fee_percentage;type:decimal(5,4);default:0.0000"` // 费用百分比
	FixedFee      *float64 `json:"fixed_fee" gorm:"column:fixed_fee;type:decimal(8,2);default:0.00"`             // 固定费用
	MinimumFee    *float64 `json:"minimum_fee" gorm:"column:minimum_fee;type:decimal(8,2);default:0.00"`         // 最小费用
	MaximumFee    *float64 `json:"maximum_fee" gorm:"column:maximum_fee;type:decimal(8,2)"`                      // 最大费用

	// 优先级和排序
	Priority     *int `json:"priority" gorm:"column:priority;type:int;default:100"`           // 优先级（数字越小优先级越高）
	DisplayOrder *int `json:"display_order" gorm:"column:display_order;type:int;default:999"` // 显示顺序

	// 别名和个性化
	UserAlias *string `json:"user_alias" gorm:"column:user_alias;type:varchar(100)"` // 用户自定义别名
	IconType  *string `json:"icon_type" gorm:"column:icon_type;type:varchar(30)"`    // 图标类型
	Color     *string `json:"color" gorm:"column:color;type:varchar(20)"`            // 显示颜色

	// 同步和更新
	LastSyncAt   *int64  `json:"last_sync_at" gorm:"column:last_sync_at"`                                  // 最后同步时间
	SyncStatus   *string `json:"sync_status" gorm:"column:sync_status;type:varchar(30);default:'pending'"` // pending, syncing, synced, failed
	SyncAttempts *int    `json:"sync_attempts" gorm:"column:sync_attempts;type:int;default:0"`             // 同步尝试次数

	// 备注和元数据
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`       // 备注
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 标签
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 附加元数据

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (UserPaymentMethod) TableName() string {
	return "t_user_payment_methods"
}

// 状态常量
const (
	UserPaymentMethodStatusActive    = "active"
	UserPaymentMethodStatusInactive  = "inactive"
	UserPaymentMethodStatusDeleted   = "deleted"
	UserPaymentMethodStatusSuspended = "suspended"
)

// 验证方式常量
const (
	UserPaymentVerificationManual    = "manual"
	UserPaymentVerificationAutomatic = "automatic"
	UserPaymentVerificationOTP       = "otp"
	UserPaymentVerificationBiometric = "biometric"
)

// 绑定方式常量
const (
	UserPaymentBoundMethodManual    = "manual"
	UserPaymentBoundMethodAutomatic = "automatic"
	UserPaymentBoundMethodImport    = "import"
	UserPaymentBoundMethodAPI       = "api"
)

// 绑定来源常量
const (
	UserPaymentBoundSourceApp        = "app"
	UserPaymentBoundSourceWeb        = "web"
	UserPaymentBoundSourceAdminPanel = "admin_panel"
	UserPaymentBoundSourceThirdParty = "third_party"
)

// 安全级别常量
const (
	UserPaymentSecurityLow     = "low"
	UserPaymentSecurityMedium  = "medium"
	UserPaymentSecurityHigh    = "high"
	UserPaymentSecurityMaximum = "maximum"
)

// 同步状态常量
const (
	UserPaymentSyncStatusPending = "pending"
	UserPaymentSyncStatusSyncing = "syncing"
	UserPaymentSyncStatusSynced  = "synced"
	UserPaymentSyncStatusFailed  = "failed"
)

// 创建新的用户支付方式关联对象
func NewUserPaymentMethodV2() *UserPaymentMethod {
	return &UserPaymentMethod{
		UserPaymentMethodID: utils.GenerateUserPaymentMethodID(),
		Salt:                utils.GenerateSalt(),
		UserPaymentMethodValues: &UserPaymentMethodValues{
			Status:                 utils.StringPtr(UserPaymentMethodStatusActive),
			IsDefault:              utils.BoolPtr(false),
			IsPrimary:              utils.BoolPtr(false),
			IsPreferred:            utils.BoolPtr(false),
			IsVerified:             utils.BoolPtr(false),
			BoundMethod:            utils.StringPtr(UserPaymentBoundMethodManual),
			BoundSource:            utils.StringPtr(UserPaymentBoundSourceApp),
			TotalTransactions:      utils.IntPtr(0),
			TotalAmount:            utils.Float64Ptr(0.00),
			SuccessfulTransactions: utils.IntPtr(0),
			FailedTransactions:     utils.IntPtr(0),
			SuccessRate:            utils.Float64Ptr(0.00),
			AverageResponseTime:    utils.IntPtr(0),
			ReliabilityScore:       utils.Float64Ptr(0.00),
			SecurityLevel:          utils.StringPtr(UserPaymentSecurityMedium),
			RequiresBiometric:      utils.BoolPtr(false),
			RequiresPin:            utils.BoolPtr(false),
			RequiresOTP:            utils.BoolPtr(false),
			RiskLevel:              utils.StringPtr(protocol.LevelLow),
			FraudScore:             utils.Float64Ptr(0.00),
			IsBlacklisted:          utils.BoolPtr(false),
			AutoReloadEnabled:      utils.BoolPtr(false),
			NotifyOnSuccess:        utils.BoolPtr(true),
			NotifyOnFailure:        utils.BoolPtr(true),
			NotifyOnLargeAmount:    utils.BoolPtr(true),
			FeePercentage:          utils.Float64Ptr(0.0000),
			FixedFee:               utils.Float64Ptr(0.00),
			MinimumFee:             utils.Float64Ptr(0.00),
			Priority:               utils.IntPtr(100),
			DisplayOrder:           utils.IntPtr(999),
			SyncStatus:             utils.StringPtr(UserPaymentSyncStatusPending),
			SyncAttempts:           utils.IntPtr(0),
		},
	}
}

// SetValues 更新UserPaymentMethodV2Values中的非nil值
func (u *UserPaymentMethodValues) SetValues(values *UserPaymentMethodValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		u.UserID = values.UserID
	}
	if values.PaymentMethodID != nil {
		u.PaymentMethodID = values.PaymentMethodID
	}
	if values.Status != nil {
		u.Status = values.Status
	}
	if values.IsDefault != nil {
		u.IsDefault = values.IsDefault
	}
	if values.IsVerified != nil {
		u.IsVerified = values.IsVerified
	}
	if values.DailyLimit != nil {
		u.DailyLimit = values.DailyLimit
	}
	if values.SecurityLevel != nil {
		u.SecurityLevel = values.SecurityLevel
	}
	if values.Notes != nil {
		u.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		u.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (u *UserPaymentMethodValues) GetUserID() string {
	if u.UserID == nil {
		return ""
	}
	return *u.UserID
}

func (u *UserPaymentMethodValues) GetPaymentMethodID() string {
	if u.PaymentMethodID == nil {
		return ""
	}
	return *u.PaymentMethodID
}

func (u *UserPaymentMethodValues) GetStatus() string {
	if u.Status == nil {
		return UserPaymentMethodStatusActive
	}
	return *u.Status
}

func (u *UserPaymentMethodValues) GetIsDefault() bool {
	if u.IsDefault == nil {
		return false
	}
	return *u.IsDefault
}

func (u *UserPaymentMethodValues) GetIsPrimary() bool {
	if u.IsPrimary == nil {
		return false
	}
	return *u.IsPrimary
}

func (u *UserPaymentMethodValues) GetIsPreferred() bool {
	if u.IsPreferred == nil {
		return false
	}
	return *u.IsPreferred
}

func (u *UserPaymentMethodValues) GetIsVerified() bool {
	if u.IsVerified == nil {
		return false
	}
	return *u.IsVerified
}

func (u *UserPaymentMethodValues) GetSecurityLevel() string {
	if u.SecurityLevel == nil {
		return UserPaymentSecurityMedium
	}
	return *u.SecurityLevel
}

func (u *UserPaymentMethodValues) GetRiskLevel() string {
	if u.RiskLevel == nil {
		return protocol.LevelLow
	}
	return *u.RiskLevel
}

func (u *UserPaymentMethodValues) GetTotalTransactions() int {
	if u.TotalTransactions == nil {
		return 0
	}
	return *u.TotalTransactions
}

func (u *UserPaymentMethodValues) GetTotalAmount() float64 {
	if u.TotalAmount == nil {
		return 0.00
	}
	return *u.TotalAmount
}

func (u *UserPaymentMethodValues) GetSuccessRate() float64 {
	if u.SuccessRate == nil {
		return 0.00
	}
	return *u.SuccessRate
}

func (u *UserPaymentMethodValues) GetReliabilityScore() float64 {
	if u.ReliabilityScore == nil {
		return 0.00
	}
	return *u.ReliabilityScore
}

func (u *UserPaymentMethodValues) GetPriority() int {
	if u.Priority == nil {
		return 100
	}
	return *u.Priority
}

func (u *UserPaymentMethodValues) GetUserAlias() string {
	if u.UserAlias == nil {
		return ""
	}
	return *u.UserAlias
}

// Setter 方法
func (u *UserPaymentMethodValues) SetUserID(userID string) *UserPaymentMethodValues {
	u.UserID = &userID
	return u
}

func (u *UserPaymentMethodValues) SetPaymentMethodID(paymentMethodID string) *UserPaymentMethodValues {
	u.PaymentMethodID = &paymentMethodID
	return u
}

func (u *UserPaymentMethodValues) SetStatus(status string) *UserPaymentMethodValues {
	u.Status = &status
	return u
}

func (u *UserPaymentMethodValues) SetDefault(isDefault bool) *UserPaymentMethodValues {
	u.IsDefault = &isDefault
	return u
}

func (u *UserPaymentMethodValues) SetPrimary(isPrimary bool) *UserPaymentMethodValues {
	u.IsPrimary = &isPrimary
	return u
}

func (u *UserPaymentMethodValues) SetPreferred(isPreferred bool) *UserPaymentMethodValues {
	u.IsPreferred = &isPreferred
	return u
}

func (u *UserPaymentMethodValues) SetVerified(isVerified bool, method string) *UserPaymentMethodValues {
	u.IsVerified = &isVerified
	u.VerificationMethod = &method
	if isVerified {
		now := utils.TimeNowMilli()
		u.VerifiedAt = &now
	}
	return u
}

func (u *UserPaymentMethodValues) SetBoundInfo(boundBy, method, source string) *UserPaymentMethodValues {
	u.BoundBy = &boundBy
	u.BoundMethod = &method
	u.BoundSource = &source
	now := utils.TimeNowMilli()
	u.BoundAt = &now
	return u
}

func (u *UserPaymentMethodValues) SetLimits(daily, monthly, singleTransaction float64) *UserPaymentMethodValues {
	u.DailyLimit = &daily
	u.MonthlyLimit = &monthly
	u.SingleTransactionLimit = &singleTransaction
	return u
}

func (u *UserPaymentMethodValues) SetSecurity(level string, requiresBiometric, requiresPin, requiresOTP bool) *UserPaymentMethodValues {
	u.SecurityLevel = &level
	u.RequiresBiometric = &requiresBiometric
	u.RequiresPin = &requiresPin
	u.RequiresOTP = &requiresOTP
	return u
}

func (u *UserPaymentMethodValues) SetRisk(level string, fraudScore float64) *UserPaymentMethodValues {
	u.RiskLevel = &level
	u.FraudScore = &fraudScore
	return u
}

func (u *UserPaymentMethodValues) SetAutoReload(enabled bool, threshold, amount float64) *UserPaymentMethodValues {
	u.AutoReloadEnabled = &enabled
	u.AutoReloadThreshold = &threshold
	u.AutoReloadAmount = &amount
	return u
}

func (u *UserPaymentMethodValues) SetNotifications(onSuccess, onFailure, onLargeAmount bool, largeAmountThreshold float64) *UserPaymentMethodValues {
	u.NotifyOnSuccess = &onSuccess
	u.NotifyOnFailure = &onFailure
	u.NotifyOnLargeAmount = &onLargeAmount
	u.LargeAmountThreshold = &largeAmountThreshold
	return u
}

func (u *UserPaymentMethodValues) SetFees(percentage, fixed, minimum, maximum float64) *UserPaymentMethodValues {
	u.FeePercentage = &percentage
	u.FixedFee = &fixed
	u.MinimumFee = &minimum
	u.MaximumFee = &maximum
	return u
}

func (u *UserPaymentMethodValues) SetPersonalization(alias, iconType, color string, priority int) *UserPaymentMethodValues {
	u.UserAlias = &alias
	u.IconType = &iconType
	u.Color = &color
	u.Priority = &priority
	return u
}

func (u *UserPaymentMethodValues) SetValidity(validFrom, validUntil int64) *UserPaymentMethodValues {
	u.ValidFrom = &validFrom
	u.ValidUntil = &validUntil
	return u
}

// 业务方法
func (u *UserPaymentMethod) IsActive() bool {
	return u.GetStatus() == UserPaymentMethodStatusActive
}

func (u *UserPaymentMethod) IsInactive() bool {
	return u.GetStatus() == UserPaymentMethodStatusInactive
}

func (u *UserPaymentMethod) IsDeleted() bool {
	return u.GetStatus() == UserPaymentMethodStatusDeleted
}

func (u *UserPaymentMethod) IsSuspended() bool {
	return u.GetStatus() == UserPaymentMethodStatusSuspended
}

func (u *UserPaymentMethod) IsDefault() bool {
	return u.GetIsDefault()
}

func (u *UserPaymentMethod) IsPrimary() bool {
	return u.GetIsPrimary()
}

func (u *UserPaymentMethod) IsPreferred() bool {
	return u.GetIsPreferred()
}

func (u *UserPaymentMethod) IsVerified() bool {
	return u.GetIsVerified()
}

func (u *UserPaymentMethod) IsBlacklisted() bool {
	if u.UserPaymentMethodValues.IsBlacklisted == nil {
		return false
	}
	return *u.UserPaymentMethodValues.IsBlacklisted
}

func (u *UserPaymentMethod) IsHighRisk() bool {
	return u.GetRiskLevel() == protocol.LevelHigh
}

func (u *UserPaymentMethod) IsLowRisk() bool {
	return u.GetRiskLevel() == protocol.LevelLow
}

func (u *UserPaymentMethod) IsHighSecurity() bool {
	level := u.GetSecurityLevel()
	return level == UserPaymentSecurityHigh || level == UserPaymentSecurityMaximum
}

func (u *UserPaymentMethod) RequiresBiometric() bool {
	if u.UserPaymentMethodValues.RequiresBiometric == nil {
		return false
	}
	return *u.UserPaymentMethodValues.RequiresBiometric
}

func (u *UserPaymentMethod) RequiresPin() bool {
	if u.UserPaymentMethodValues.RequiresPin == nil {
		return false
	}
	return *u.UserPaymentMethodValues.RequiresPin
}

func (u *UserPaymentMethod) RequiresOTP() bool {
	if u.UserPaymentMethodValues.RequiresOTP == nil {
		return false
	}
	return *u.UserPaymentMethodValues.RequiresOTP
}

func (u *UserPaymentMethod) HasAutoReload() bool {
	if u.UserPaymentMethodValues.AutoReloadEnabled == nil {
		return false
	}
	return *u.UserPaymentMethodValues.AutoReloadEnabled
}

func (u *UserPaymentMethod) IsExpired() bool {
	if u.UserPaymentMethodValues.ValidUntil == nil {
		return false
	}
	return *u.UserPaymentMethodValues.ValidUntil < utils.TimeNowMilli()
}

func (u *UserPaymentMethod) IsValid() bool {
	now := utils.TimeNowMilli()

	// 检查有效期
	if u.UserPaymentMethodValues.ValidFrom != nil && *u.UserPaymentMethodValues.ValidFrom > now {
		return false
	}
	if u.UserPaymentMethodValues.ValidUntil != nil && *u.UserPaymentMethodValues.ValidUntil < now {
		return false
	}

	// 检查状态
	return u.IsActive() && u.IsVerified() && !u.IsBlacklisted()
}

func (u *UserPaymentMethod) CanProcess(amount float64) bool {
	if !u.IsValid() {
		return false
	}

	// 检查单笔限额
	if u.UserPaymentMethodValues.SingleTransactionLimit != nil {
		if amount > *u.UserPaymentMethodValues.SingleTransactionLimit {
			return false
		}
	}

	return true
}

func (u *UserPaymentMethod) IsReliable() bool {
	score := u.GetReliabilityScore()
	successRate := u.GetSuccessRate()
	return score >= 4.0 && successRate >= 95.0
}

func (u *UserPaymentMethod) IsFrequentlyUsed() bool {
	return u.GetTotalTransactions() >= 50
}

// 统计更新方法
func (u *UserPaymentMethodValues) UpdateTransactionStats(amount float64, success bool) *UserPaymentMethodValues {
	// 更新总交易次数和金额
	totalTx := u.GetTotalTransactions() + 1
	totalAmount := u.GetTotalAmount() + amount
	u.TotalTransactions = &totalTx
	u.TotalAmount = &totalAmount

	// 更新成功/失败次数
	if success {
		successTx := u.GetSuccessfulTransactions() + 1
		u.SuccessfulTransactions = &successTx
	} else {
		failedTx := u.GetFailedTransactions() + 1
		u.FailedTransactions = &failedTx
	}

	// 更新成功率
	successRate := float64(u.GetSuccessfulTransactions()) / float64(totalTx) * 100.0
	u.SuccessRate = &successRate

	// 更新使用时间
	now := utils.TimeNowMilli()
	u.LastUsedAt = &now
	if u.FirstUsedAt == nil {
		u.FirstUsedAt = &now
	}

	// 更新可靠性评分
	u.UpdateReliabilityScore()

	return u
}

func (u *UserPaymentMethodValues) GetSuccessfulTransactions() int {
	if u.SuccessfulTransactions == nil {
		return 0
	}
	return *u.SuccessfulTransactions
}

func (u *UserPaymentMethodValues) GetFailedTransactions() int {
	if u.FailedTransactions == nil {
		return 0
	}
	return *u.FailedTransactions
}

func (u *UserPaymentMethodValues) UpdateReliabilityScore() *UserPaymentMethodValues {
	score := 0.0

	// 成功率权重 (40%)
	successRate := u.GetSuccessRate()
	if successRate >= 95.0 {
		score += 2.0
	} else if successRate >= 90.0 {
		score += 1.5
	} else if successRate >= 80.0 {
		score += 1.0
	} else {
		score += 0.5
	}

	// 交易次数权重 (30%)
	totalTx := u.GetTotalTransactions()
	if totalTx >= 100 {
		score += 1.5
	} else if totalTx >= 50 {
		score += 1.0
	} else if totalTx >= 10 {
		score += 0.5
	}

	// 响应时间权重 (20%)
	avgResponseTime := u.GetAverageResponseTime()
	if avgResponseTime <= 1000 {
		score += 1.0
	} else if avgResponseTime <= 3000 {
		score += 0.7
	} else if avgResponseTime <= 5000 {
		score += 0.4
	}

	// 风险评分权重 (10%)
	if u.GetRiskLevel() == protocol.LevelLow {
		score += 0.5
	} else if u.GetRiskLevel() == protocol.LevelMedium {
		score += 0.3
	}

	// 确保评分在0-5范围内
	if score > 5.0 {
		score = 5.0
	}

	u.ReliabilityScore = &score
	return u
}

func (u *UserPaymentMethodValues) GetAverageResponseTime() int {
	if u.AverageResponseTime == nil {
		return 0
	}
	return *u.AverageResponseTime
}

func (u *UserPaymentMethodValues) UpdateResponseTime(responseTimeMs int) *UserPaymentMethodValues {
	currentAvg := u.GetAverageResponseTime()
	totalTx := u.GetTotalTransactions()

	// 计算新的平均响应时间
	if totalTx <= 1 {
		u.AverageResponseTime = &responseTimeMs
	} else {
		newAvg := (currentAvg*(totalTx-1) + responseTimeMs) / totalTx
		u.AverageResponseTime = &newAvg
	}

	return u
}

// 黑名单管理
func (u *UserPaymentMethodValues) AddToBlacklist(reason string) *UserPaymentMethodValues {
	u.IsBlacklisted = utils.BoolPtr(true)
	u.BlacklistReason = &reason
	u.SetStatus(UserPaymentMethodStatusSuspended)
	return u
}

func (u *UserPaymentMethodValues) RemoveFromBlacklist() *UserPaymentMethodValues {
	u.IsBlacklisted = utils.BoolPtr(false)
	u.BlacklistReason = nil
	u.SetStatus(UserPaymentMethodStatusActive)
	return u
}

// 同步状态管理
func (u *UserPaymentMethodValues) StartSync() *UserPaymentMethodValues {
	u.SyncStatus = utils.StringPtr(UserPaymentSyncStatusSyncing)
	now := utils.TimeNowMilli()
	u.LastSyncAt = &now
	return u
}

func (u *UserPaymentMethodValues) CompleteSync(success bool) *UserPaymentMethodValues {
	if success {
		u.SyncStatus = utils.StringPtr(UserPaymentSyncStatusSynced)
		u.SyncAttempts = utils.IntPtr(0)
	} else {
		u.SyncStatus = utils.StringPtr(UserPaymentSyncStatusFailed)
		attempts := u.GetSyncAttempts() + 1
		u.SyncAttempts = &attempts
	}

	now := utils.TimeNowMilli()
	u.LastSyncAt = &now
	return u
}

func (u *UserPaymentMethodValues) GetSyncAttempts() int {
	if u.SyncAttempts == nil {
		return 0
	}
	return *u.SyncAttempts
}

// 便捷创建方法
func NewUserCreditCardPayment(userID, paymentMethodID string) *UserPaymentMethod {
	upm := NewUserPaymentMethodV2()
	upm.SetUserID(userID).
		SetPaymentMethodID(paymentMethodID).
		SetSecurity(UserPaymentSecurityHigh, false, true, false).
		SetLimits(5000.0, 50000.0, 2000.0).
		SetNotifications(true, true, true, 1000.0).
		SetPersonalization("My Credit Card", "credit_card", "#1976D2", 1)

	return upm
}

func NewUserMobileWalletPayment(userID, paymentMethodID string) *UserPaymentMethod {
	upm := NewUserPaymentMethodV2()
	upm.SetUserID(userID).
		SetPaymentMethodID(paymentMethodID).
		SetSecurity(UserPaymentSecurityMedium, true, false, false).
		SetLimits(2000.0, 20000.0, 1000.0).
		SetNotifications(true, true, true, 500.0).
		SetPersonalization("Mobile Wallet", "mobile_wallet", "#4CAF50", 2)

	return upm
}

func NewUserBankAccountPayment(userID, paymentMethodID string) *UserPaymentMethod {
	upm := NewUserPaymentMethodV2()
	upm.SetUserID(userID).
		SetPaymentMethodID(paymentMethodID).
		SetSecurity(UserPaymentSecurityMaximum, false, true, true).
		SetLimits(10000.0, 100000.0, 5000.0).
		SetNotifications(true, true, true, 2000.0).
		SetPersonalization("Bank Account", "bank", "#FF9800", 3).
		SetAutoReload(true, 500.0, 1000.0)

	return upm
}

func NewUserCashPayment(userID, paymentMethodID string) *UserPaymentMethod {
	upm := NewUserPaymentMethodV2()
	upm.SetUserID(userID).
		SetPaymentMethodID(paymentMethodID).
		SetSecurity(UserPaymentSecurityLow, false, false, false).
		SetLimits(500.0, 5000.0, 200.0).
		SetNotifications(false, false, false, 0.0).
		SetPersonalization("Cash", "cash", "#795548", 10)

	return upm
}
