package models

import (
	"fmt"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"
	"time"
)

// PaymentMethod 支付方式表 - 用户绑定的银行卡、移动支付等支付方式
type PaymentMethod struct {
	ID              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PaymentMethodID string `json:"payment_method_id" gorm:"column:payment_method_id;type:varchar(64);uniqueIndex"`
	Salt            string `json:"salt" gorm:"column:salt;type:varchar(256)"`
	*PaymentMethodValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
}

type PaymentMethodValues struct {
	// 基本信息
	UserID   *string `json:"user_id" gorm:"column:user_id;type:varchar(64);index"`
	UserType *string `json:"user_type" gorm:"column:user_type;type:varchar(32);index;default:'user'"` // user, driver

	// 支付方式类型
	PaymentType  *string `json:"payment_type" gorm:"column:payment_type;type:varchar(50);index"`    // bank_card, credit_card, mobile_money, paypal, stripe, alipay, wechat
	ProviderName *string `json:"provider_name" gorm:"column:provider_name;type:varchar(100);index"` // Equity Bank, KCB, MTN, Airtel, PayPal, Stripe
	SubType      *string `json:"sub_type" gorm:"column:sub_type;type:varchar(50)"`                  // visa, mastercard, mtn_momo, airtel_money

	// 账户信息 (加密存储)
	AccountNumber *string `json:"account_number" gorm:"column:account_number;type:varchar(255)"` // 银行账号/手机号/邮箱等 (加密)
	AccountName   *string `json:"account_name" gorm:"column:account_name;type:varchar(255)"`     // 账户名称 (加密)
	MaskedNumber  *string `json:"masked_number" gorm:"column:masked_number;type:varchar(100)"`   // 脱敏显示号码

	// 银行卡信息
	CardNumber     *string `json:"card_number" gorm:"column:card_number;type:varchar(255)"`           // 银行卡号 (加密)
	CardHolderName *string `json:"card_holder_name" gorm:"column:card_holder_name;type:varchar(255)"` // 持卡人姓名 (加密)
	ExpiryMonth    *int    `json:"expiry_month" gorm:"column:expiry_month;type:int"`                  // 过期月份
	ExpiryYear     *int    `json:"expiry_year" gorm:"column:expiry_year;type:int"`                    // 过期年份
	CVV            *string `json:"cvv" gorm:"column:cvv;type:varchar(255)"`                           // CVV (加密)
	BankName       *string `json:"bank_name" gorm:"column:bank_name;type:varchar(100)"`               // 银行名称
	BankCode       *string `json:"bank_code" gorm:"column:bank_code;type:varchar(20)"`                // 银行代码
	BranchName     *string `json:"branch_name" gorm:"column:branch_name;type:varchar(100)"`           // 分行名称
	SwiftCode      *string `json:"swift_code" gorm:"column:swift_code;type:varchar(20)"`              // SWIFT代码

	// 移动支付信息
	PhoneNumber   *string `json:"phone_number" gorm:"column:phone_number;type:varchar(255)"`    // 手机号 (加密)
	MobileCarrier *string `json:"mobile_carrier" gorm:"column:mobile_carrier;type:varchar(50)"` // 运营商: MTN, Airtel, Tigo

	// 第三方支付信息
	ExternalID            *string `json:"external_id" gorm:"column:external_id;type:varchar(255)"`                           // 第三方平台ID
	PaypalEmail           *string `json:"paypal_email" gorm:"column:paypal_email;type:varchar(255)"`                         // PayPal邮箱 (加密)
	StripeCustomerID      *string `json:"stripe_customer_id" gorm:"column:stripe_customer_id;type:varchar(255)"`             // Stripe客户ID
	StripePaymentMethodID *string `json:"stripe_payment_method_id" gorm:"column:stripe_payment_method_id;type:varchar(255)"` // Stripe支付方式ID
	AlipayAccount         *string `json:"alipay_account" gorm:"column:alipay_account;type:varchar(255)"`                     // 支付宝账号 (加密)
	WechatAccount         *string `json:"wechat_account" gorm:"column:wechat_account;type:varchar(255)"`                     // 微信账号 (加密)

	// 状态管理
	Status     *string `json:"status" gorm:"column:status;type:varchar(32);index;default:'active'"` // active, inactive, suspended, expired, blocked
	IsDefault  *bool   `json:"is_default" gorm:"column:is_default;default:false"`                   // 是否默认支付方式
	IsVerified *bool   `json:"is_verified" gorm:"column:is_verified;default:false"`                 // 是否已验证
	VerifiedAt *int64  `json:"verified_at" gorm:"column:verified_at"`                               // 验证时间

	// 验证信息
	VerificationMethod   *string `json:"verification_method" gorm:"column:verification_method;type:varchar(50)"` // sms, email, micro_transfer, document
	VerificationCode     *string `json:"verification_code" gorm:"column:verification_code;type:varchar(10)"`     // 验证码
	VerificationExpiry   *int64  `json:"verification_expiry" gorm:"column:verification_expiry"`                  // 验证码过期时间
	VerificationAttempts *int    `json:"verification_attempts" gorm:"column:verification_attempts;default:0"`    // 验证尝试次数

	// 使用统计
	FirstUsedAt *int64   `json:"first_used_at" gorm:"column:first_used_at"`                               // 首次使用时间
	LastUsedAt  *int64   `json:"last_used_at" gorm:"column:last_used_at"`                                 // 最后使用时间
	UsageCount  *int     `json:"usage_count" gorm:"column:usage_count;default:0"`                         // 使用次数
	TotalAmount *float64 `json:"total_amount" gorm:"column:total_amount;type:decimal(15,2);default:0.00"` // 累计交易金额
	LastAmount  *float64 `json:"last_amount" gorm:"column:last_amount;type:decimal(10,2)"`                // 最后一次交易金额

	// 安全设置
	DailyLimit   *float64 `json:"daily_limit" gorm:"column:daily_limit;type:decimal(10,2)"`     // 日限额
	MonthlyLimit *float64 `json:"monthly_limit" gorm:"column:monthly_limit;type:decimal(15,2)"` // 月限额
	SingleLimit  *float64 `json:"single_limit" gorm:"column:single_limit;type:decimal(10,2)"`   // 单笔限额
	RequireAuth  *bool    `json:"require_auth" gorm:"column:require_auth;default:false"`        // 是否需要二次认证

	// 风控信息
	RiskLevel     *string  `json:"risk_level" gorm:"column:risk_level;type:varchar(20);default:'low'"` // low, medium, high
	RiskScore     *float64 `json:"risk_score" gorm:"column:risk_score;type:decimal(5,2);default:0.00"` // 风险评分
	RiskFlags     *string  `json:"risk_flags" gorm:"column:risk_flags;type:json"`                      // 风险标记
	BlockedReason *string  `json:"blocked_reason" gorm:"column:blocked_reason;type:varchar(500)"`      // 封禁原因
	BlockedAt     *int64   `json:"blocked_at" gorm:"column:blocked_at"`                                // 封禁时间
	BlockedBy     *string  `json:"blocked_by" gorm:"column:blocked_by;type:varchar(64)"`               // 封禁操作员

	// 附加信息
	DisplayName *string `json:"display_name" gorm:"column:display_name;type:varchar(100)"`            // 显示名称
	IconURL     *string `json:"icon_url" gorm:"column:icon_url;type:varchar(500)"`                    // 图标URL
	Currency    *string `json:"currency" gorm:"column:currency;type:varchar(10);default:'RWF'"`       // 支持币种
	CountryCode *string `json:"country_code" gorm:"column:country_code;type:varchar(5);default:'RW'"` // 国家代码

	// 元数据
	Metadata *string `json:"metadata" gorm:"column:metadata;type:json"` // 附加元数据
	Notes    *string `json:"notes" gorm:"column:notes;type:text"`       // 备注信息
	Tags     *string `json:"tags" gorm:"column:tags;type:varchar(500)"` // 标签 (逗号分隔)

	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (PaymentMethod) TableName() string {
	return "t_payment_methods"
}

// 支付方式类型常量
const (
	PaymentTypeBankCard    = "bank_card"
	PaymentTypeCreditCard  = "credit_card"
	PaymentTypeMobileMoney = "mobile_money"
	PaymentTypePaypal      = "paypal"
	PaymentTypeStripe      = "stripe"
	PaymentTypeAlipay      = "alipay"
	PaymentTypeWechat      = "wechat"
)

// 支付状态常量
const (
	PaymentMethodStatusActive    = "active"
	PaymentMethodStatusInactive  = "inactive"
	PaymentMethodStatusSuspended = "suspended"
	PaymentMethodStatusExpired   = "expired"
	PaymentMethodStatusBlocked   = "blocked"
)

// 验证方式常量
const (
	VerificationMethodSMS           = "sms"
	VerificationMethodEmail         = "email"
	VerificationMethodMicroTransfer = "micro_transfer"
	VerificationMethodDocument      = "document"
)

// 风险等级常量
const (
	RiskLevelLow    = "low"
	RiskLevelMedium = "medium"
	RiskLevelHigh   = "high"
)

// 创建新的支付方式对象
func NewPaymentMethod() *PaymentMethod {
	return &PaymentMethod{
		PaymentMethodID: utils.GeneratePaymentMethodID(),
		Salt:            utils.GenerateSalt(),
		PaymentMethodValues: &PaymentMethodValues{
			UserType:             utils.StringPtr(protocol.UserTypePassenger),
			Status:               utils.StringPtr(PaymentMethodStatusActive),
			IsDefault:            utils.BoolPtr(false),
			IsVerified:           utils.BoolPtr(false),
			UsageCount:           utils.IntPtr(0),
			TotalAmount:          utils.Float64Ptr(0.00),
			VerificationAttempts: utils.IntPtr(0),
			RequireAuth:          utils.BoolPtr(false),
			RiskLevel:            utils.StringPtr(RiskLevelLow),
			RiskScore:            utils.Float64Ptr(0.00),
			Currency:             utils.StringPtr("RWF"),
			CountryCode:          utils.StringPtr("RW"),
		},
	}
}

// SetValues 更新PaymentMethodValues中的非nil值
func (p *PaymentMethodValues) SetValues(values *PaymentMethodValues) {
	if values == nil {
		return
	}

	if values.UserID != nil {
		p.UserID = values.UserID
	}
	if values.PaymentType != nil {
		p.PaymentType = values.PaymentType
	}
	if values.ProviderName != nil {
		p.ProviderName = values.ProviderName
	}
	if values.AccountNumber != nil {
		p.AccountNumber = values.AccountNumber
	}
	if values.AccountName != nil {
		p.AccountName = values.AccountName
	}
	if values.Status != nil {
		p.Status = values.Status
	}
	if values.IsDefault != nil {
		p.IsDefault = values.IsDefault
	}
	if values.DisplayName != nil {
		p.DisplayName = values.DisplayName
	}
	if values.Notes != nil {
		p.Notes = values.Notes
	}
	if values.UpdatedAt > 0 {
		p.UpdatedAt = values.UpdatedAt
	}
}

// Getter 方法
func (p *PaymentMethodValues) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return *p.UserID
}

func (p *PaymentMethodValues) GetPaymentType() string {
	if p.PaymentType == nil {
		return ""
	}
	return *p.PaymentType
}

func (p *PaymentMethodValues) GetProviderName() string {
	if p.ProviderName == nil {
		return ""
	}
	return *p.ProviderName
}

func (p *PaymentMethodValues) GetAccountName() string {
	if p.AccountName == nil {
		return ""
	}
	return *p.AccountName
}

func (p *PaymentMethodValues) GetMaskedNumber() string {
	if p.MaskedNumber == nil {
		return ""
	}
	return *p.MaskedNumber
}

func (p *PaymentMethodValues) GetStatus() string {
	if p.Status == nil {
		return PaymentMethodStatusActive
	}
	return *p.Status
}

func (p *PaymentMethodValues) GetIsDefault() bool {
	if p.IsDefault == nil {
		return false
	}
	return *p.IsDefault
}

func (p *PaymentMethodValues) GetIsVerified() bool {
	if p.IsVerified == nil {
		return false
	}
	return *p.IsVerified
}

func (p *PaymentMethodValues) GetUsageCount() int {
	if p.UsageCount == nil {
		return 0
	}
	return *p.UsageCount
}

func (p *PaymentMethodValues) GetTotalAmount() float64 {
	if p.TotalAmount == nil {
		return 0.00
	}
	return *p.TotalAmount
}

func (p *PaymentMethodValues) GetRiskLevel() string {
	if p.RiskLevel == nil {
		return RiskLevelLow
	}
	return *p.RiskLevel
}

func (p *PaymentMethodValues) GetRiskScore() float64 {
	if p.RiskScore == nil {
		return 0.00
	}
	return *p.RiskScore
}

func (p *PaymentMethodValues) GetRequireAuth() bool {
	if p.RequireAuth == nil {
		return false
	}
	return *p.RequireAuth
}

// Setter 方法
func (p *PaymentMethodValues) SetUserID(userID string) *PaymentMethodValues {
	p.UserID = &userID
	return p
}

func (p *PaymentMethodValues) SetPaymentType(paymentType string) *PaymentMethodValues {
	p.PaymentType = &paymentType
	return p
}

func (p *PaymentMethodValues) SetProviderName(provider string) *PaymentMethodValues {
	p.ProviderName = &provider
	return p
}

func (p *PaymentMethodValues) SetAccountInfo(accountNumber, accountName string) *PaymentMethodValues {
	p.AccountNumber = &accountNumber
	p.AccountName = &accountName

	// 生成脱敏号码
	masked := p.MaskAccountNumber(accountNumber)
	p.MaskedNumber = &masked

	return p
}

func (p *PaymentMethodValues) SetCardInfo(cardNumber, holderName string, expiryMonth, expiryYear int) *PaymentMethodValues {
	p.CardNumber = &cardNumber
	p.CardHolderName = &holderName
	p.ExpiryMonth = &expiryMonth
	p.ExpiryYear = &expiryYear

	// 生成脱敏卡号
	masked := p.MaskCardNumber(cardNumber)
	p.MaskedNumber = &masked

	return p
}

func (p *PaymentMethodValues) SetMobileInfo(phoneNumber, carrier string) *PaymentMethodValues {
	p.PhoneNumber = &phoneNumber
	p.MobileCarrier = &carrier

	// 生成脱敏手机号
	masked := p.MaskPhoneNumber(phoneNumber)
	p.MaskedNumber = &masked

	return p
}

func (p *PaymentMethodValues) SetStatus(status string) *PaymentMethodValues {
	p.Status = &status
	return p
}

func (p *PaymentMethodValues) SetAsDefault(isDefault bool) *PaymentMethodValues {
	p.IsDefault = &isDefault
	return p
}

func (p *PaymentMethodValues) SetDisplayName(name string) *PaymentMethodValues {
	p.DisplayName = &name
	return p
}

func (p *PaymentMethodValues) SetVerified(verified bool) *PaymentMethodValues {
	p.IsVerified = &verified
	if verified {
		now := utils.TimeNowMilli()
		p.VerifiedAt = &now
	}
	return p
}

func (p *PaymentMethodValues) SetLimits(daily, monthly, single float64) *PaymentMethodValues {
	if daily > 0 {
		p.DailyLimit = &daily
	}
	if monthly > 0 {
		p.MonthlyLimit = &monthly
	}
	if single > 0 {
		p.SingleLimit = &single
	}
	return p
}

func (p *PaymentMethodValues) SetRiskLevel(level string) *PaymentMethodValues {
	p.RiskLevel = &level
	return p
}

func (p *PaymentMethodValues) SetRiskScore(score float64) *PaymentMethodValues {
	p.RiskScore = &score
	return p
}

// 业务方法
func (p *PaymentMethod) IsActive() bool {
	return p.GetStatus() == PaymentMethodStatusActive
}

func (p *PaymentMethod) IsInactive() bool {
	return p.GetStatus() == PaymentMethodStatusInactive
}

func (p *PaymentMethod) IsBlocked() bool {
	return p.GetStatus() == PaymentMethodStatusBlocked
}

func (p *PaymentMethod) IsExpired() bool {
	if p.GetStatus() == PaymentMethodStatusExpired {
		return true
	}

	// 检查银行卡是否过期
	if p.ExpiryMonth != nil && p.ExpiryYear != nil {
		now := time.Now()
		expiryDate := time.Date(*p.ExpiryYear, time.Month(*p.ExpiryMonth), 1, 0, 0, 0, 0, time.UTC)
		return now.After(expiryDate)
	}

	return false
}

func (p *PaymentMethod) IsBankCard() bool {
	paymentType := p.GetPaymentType()
	return paymentType == PaymentTypeBankCard || paymentType == PaymentTypeCreditCard
}

func (p *PaymentMethod) IsMobileMoney() bool {
	return p.GetPaymentType() == PaymentTypeMobileMoney
}

func (p *PaymentMethod) IsThirdParty() bool {
	paymentType := p.GetPaymentType()
	return paymentType == PaymentTypePaypal || paymentType == PaymentTypeStripe ||
		paymentType == PaymentTypeAlipay || paymentType == PaymentTypeWechat
}

func (p *PaymentMethod) CanUse() bool {
	return p.IsActive() && p.GetIsVerified() && !p.IsExpired()
}

func (p *PaymentMethod) IsHighRisk() bool {
	return p.GetRiskLevel() == RiskLevelHigh || p.GetRiskScore() >= 80.0
}

func (p *PaymentMethod) NeedsVerification() bool {
	return !p.GetIsVerified()
}

// 脱敏方法
func (p *PaymentMethodValues) MaskCardNumber(cardNumber string) string {
	if len(cardNumber) < 8 {
		return "****"
	}
	return cardNumber[:4] + "****" + cardNumber[len(cardNumber)-4:]
}

func (p *PaymentMethodValues) MaskAccountNumber(accountNumber string) string {
	if len(accountNumber) < 6 {
		return "****"
	}
	return "****" + accountNumber[len(accountNumber)-4:]
}

func (p *PaymentMethodValues) MaskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 8 {
		return "****"
	}
	return phoneNumber[:3] + "****" + phoneNumber[len(phoneNumber)-3:]
}

// 使用统计更新
func (p *PaymentMethodValues) RecordUsage(amount float64) *PaymentMethodValues {
	now := utils.TimeNowMilli()

	// 更新使用统计
	count := p.GetUsageCount() + 1
	p.UsageCount = &count

	total := p.GetTotalAmount() + amount
	p.TotalAmount = &total

	p.LastAmount = &amount
	p.LastUsedAt = &now

	// 如果是首次使用
	if p.FirstUsedAt == nil {
		p.FirstUsedAt = &now
	}

	return p
}

// 验证相关方法
func (p *PaymentMethodValues) StartVerification(method string) *PaymentMethodValues {
	p.VerificationMethod = &method

	// 生成验证码
	code := utils.GenerateVerifyCode()
	p.VerificationCode = &code

	// 设置过期时间 (15分钟)
	expiry := utils.TimeNowMilli() + (15 * 60 * 1000)
	p.VerificationExpiry = &expiry

	// 重置尝试次数
	p.VerificationAttempts = utils.IntPtr(0)

	return p
}

func (p *PaymentMethodValues) VerifyCode(inputCode string) bool {
	// 检查验证码是否过期
	if p.VerificationExpiry != nil && utils.TimeNowMilli() > *p.VerificationExpiry {
		return false
	}

	// 增加尝试次数
	attempts := 0
	if p.VerificationAttempts != nil {
		attempts = *p.VerificationAttempts
	}
	attempts++
	p.VerificationAttempts = &attempts

	// 检查验证码
	if p.VerificationCode != nil && *p.VerificationCode == inputCode {
		p.SetVerified(true)
		// 清除验证信息
		p.VerificationCode = nil
		p.VerificationExpiry = nil
		return true
	}

	return false
}

// 风控方法
func (p *PaymentMethodValues) AddRiskFlag(flag string) error {
	var flags []string
	if p.RiskFlags != nil {
		if err := utils.FromJSON(*p.RiskFlags, &flags); err != nil {
			return fmt.Errorf("failed to parse existing risk flags: %v", err)
		}
	}

	// 避免重复添加
	for _, existingFlag := range flags {
		if existingFlag == flag {
			return nil
		}
	}

	flags = append(flags, flag)
	flagsJSON, err := utils.ToJSON(flags)
	if err != nil {
		return fmt.Errorf("failed to marshal risk flags: %v", err)
	}

	p.RiskFlags = &flagsJSON
	return nil
}

func (p *PaymentMethodValues) Block(adminID, reason string) *PaymentMethodValues {
	p.SetStatus(PaymentMethodStatusBlocked)
	p.BlockedReason = &reason
	p.BlockedBy = &adminID
	now := utils.TimeNowMilli()
	p.BlockedAt = &now
	return p
}

func (p *PaymentMethodValues) Unblock() *PaymentMethodValues {
	p.SetStatus(PaymentMethodStatusActive)
	p.BlockedReason = nil
	p.BlockedBy = nil
	p.BlockedAt = nil
	return p
}

// 标签管理
func (p *PaymentMethodValues) AddTag(tag string) *PaymentMethodValues {
	var tags []string
	if p.Tags != nil && *p.Tags != "" {
		tags = strings.Split(*p.Tags, ",")
	}

	// 避免重复
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return p
		}
	}

	tags = append(tags, tag)
	tagsStr := strings.Join(tags, ",")
	p.Tags = &tagsStr
	return p
}

func (p *PaymentMethodValues) RemoveTag(tag string) *PaymentMethodValues {
	if p.Tags == nil || *p.Tags == "" {
		return p
	}

	tags := strings.Split(*p.Tags, ",")
	var newTags []string

	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) != tag {
			newTags = append(newTags, strings.TrimSpace(existingTag))
		}
	}

	tagsStr := strings.Join(newTags, ",")
	p.Tags = &tagsStr
	return p
}

func (p *PaymentMethodValues) HasTag(tag string) bool {
	if p.Tags == nil || *p.Tags == "" {
		return false
	}

	tags := strings.Split(*p.Tags, ",")
	for _, existingTag := range tags {
		if strings.TrimSpace(existingTag) == tag {
			return true
		}
	}

	return false
}

// 便捷创建方法
func NewBankCardPaymentMethod(userID, bankName, accountNumber, accountName string) *PaymentMethod {
	payment := NewPaymentMethod()
	payment.SetUserID(userID).
		SetPaymentType(PaymentTypeBankCard).
		SetProviderName(bankName).
		SetAccountInfo(accountNumber, accountName)

	payment.BankName = &bankName
	payment.SetDisplayName(bankName + " " + payment.GetMaskedNumber())

	return payment
}

func NewMobileMoneyPaymentMethod(userID, carrier, phoneNumber string) *PaymentMethod {
	payment := NewPaymentMethod()
	payment.SetUserID(userID).
		SetPaymentType(PaymentTypeMobileMoney).
		SetProviderName(carrier).
		SetMobileInfo(phoneNumber, carrier)

	payment.SetDisplayName(carrier + " " + payment.GetMaskedNumber())

	return payment
}

func NewCreditCardPaymentMethod(userID, cardNumber, holderName string, expiryMonth, expiryYear int) *PaymentMethod {
	payment := NewPaymentMethod()
	payment.SetUserID(userID).
		SetPaymentType(PaymentTypeCreditCard).
		SetCardInfo(cardNumber, holderName, expiryMonth, expiryYear)

	// 根据卡号判断卡类型
	cardType := "Card"
	if strings.HasPrefix(cardNumber, "4") {
		cardType = "Visa"
		payment.SubType = utils.StringPtr("visa")
	} else if strings.HasPrefix(cardNumber, "5") {
		cardType = "Mastercard"
		payment.SubType = utils.StringPtr("mastercard")
	}

	payment.SetProviderName(cardType)
	payment.SetDisplayName(cardType + " " + payment.GetMaskedNumber())

	return payment
}

func NewPaypalPaymentMethod(userID, email string) *PaymentMethod {
	payment := NewPaymentMethod()
	payment.SetUserID(userID).
		SetPaymentType(PaymentTypePaypal).
		SetProviderName("PayPal")

	payment.PaypalEmail = &email
	maskedEmail := strings.Repeat("*", 3) + email[strings.Index(email, "@"):]
	payment.MaskedNumber = &maskedEmail
	payment.SetDisplayName("PayPal " + maskedEmail)

	return payment
}
