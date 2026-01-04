package models

import (
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PaymentRouters 支付路由表
type PaymentRouters struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	RouterID  string `json:"router_id" gorm:"column:router_id;type:varchar(64);uniqueIndex:idx_router_id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	*PaymentRouterValues
}

// PaymentRouterValues 支付路由配置值
type PaymentRouterValues struct {
	Name             *string          `json:"name" gorm:"column:name;type:varchar(100)"`                                                      // 路由规则名称
	ChannelCode      *string          `json:"channel_code" gorm:"column:channel_code;type:varchar(32);index:idx_channel_code"`                // 渠道代码
	ChannelAccountID *string          `json:"channel_account_id" gorm:"column:channel_account_id;type:varchar(64);index:idx_channel_account"` // 渠道账户ID
	PaymentMethod    *string          `json:"payment_method" gorm:"column:payment_method;type:varchar(32);index:idx_payment_method"`          // 支付方式
	Currency         *string          `json:"currency" gorm:"column:currency;type:varchar(8);index:idx_currency"`                             // 货币类型
	MinAmount        *decimal.Decimal `json:"min_amount" gorm:"column:min_amount;type:decimal(15,2)"`                                         // 最小金额
	MaxAmount        *decimal.Decimal `json:"max_amount" gorm:"column:max_amount;type:decimal(15,2)"`                                         // 最大金额
	Priority         *int             `json:"priority" gorm:"column:priority;type:int;default:100;index:idx_priority"`                        // 优先级(数字越大优先级越高)
	Status           *string          `json:"status" gorm:"column:status;type:varchar(20);default:'inactive';index:idx_status"`               // 状态
	EffectiveTime    *int64           `json:"effective_time" gorm:"column:effective_time;type:bigint"`                                        // 生效时间
	ExpireTime       *int64           `json:"expire_time" gorm:"column:expire_time;type:bigint"`                                              // 失效时间
	Region           *string          `json:"region" gorm:"column:region;type:varchar(32)"`                                                   // 地区限制
	Remark           *string          `json:"remark" gorm:"column:remark;type:varchar(500)"`                                                  // 备注
	UpdatedAt        *int64           `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`                                       // 更新时间
}

func (PaymentRouters) TableName() string {
	return "t_payment_routers"
}

// NewPaymentRouter 创建新的支付路由
func NewPaymentRouter() *PaymentRouters {
	status := protocol.StatusInactive
	priority := 100
	return &PaymentRouters{
		RouterID: utils.GenerateID(),
		PaymentRouterValues: &PaymentRouterValues{
			Status:   &status,
			Priority: &priority,
		},
	}
}

// Getter methods for PaymentRouterValues
func (p *PaymentRouterValues) GetName() string {
	if p.Name == nil {
		return ""
	}
	return *p.Name
}

func (p *PaymentRouterValues) GetChannelCode() string {
	if p.ChannelCode == nil {
		return ""
	}
	return *p.ChannelCode
}

func (p *PaymentRouterValues) GetChannelAccountID() string {
	if p.ChannelAccountID == nil {
		return ""
	}
	return *p.ChannelAccountID
}

func (p *PaymentRouterValues) GetPaymentMethod() string {
	if p.PaymentMethod == nil {
		return ""
	}
	return *p.PaymentMethod
}

func (p *PaymentRouterValues) GetCurrency() string {
	if p.Currency == nil {
		return ""
	}
	return *p.Currency
}

func (p *PaymentRouterValues) GetMinAmount() decimal.Decimal {
	if p.MinAmount == nil {
		return decimal.Zero
	}
	return *p.MinAmount
}

func (p *PaymentRouterValues) GetMaxAmount() decimal.Decimal {
	if p.MaxAmount == nil {
		return decimal.Zero
	}
	return *p.MaxAmount
}

func (p *PaymentRouterValues) GetPriority() int {
	if p.Priority == nil {
		return 100
	}
	return *p.Priority
}

func (p *PaymentRouterValues) GetStatus() string {
	if p.Status == nil {
		return protocol.StatusInactive
	}
	return *p.Status
}

func (p *PaymentRouterValues) GetEffectiveTime() int64 {
	if p.EffectiveTime == nil {
		return 0
	}
	return *p.EffectiveTime
}

func (p *PaymentRouterValues) GetExpireTime() int64 {
	if p.ExpireTime == nil {
		return 0
	}
	return *p.ExpireTime
}

func (p *PaymentRouterValues) GetRegion() string {
	if p.Region == nil {
		return ""
	}
	return *p.Region
}

func (p *PaymentRouterValues) GetRemark() string {
	if p.Remark == nil {
		return ""
	}
	return *p.Remark
}

// Setter methods for PaymentRouterValues (支持链式调用)
func (p *PaymentRouterValues) SetName(name string) *PaymentRouterValues {
	p.Name = &name
	return p
}

func (p *PaymentRouterValues) SetChannelCode(channelCode string) *PaymentRouterValues {
	p.ChannelCode = &channelCode
	return p
}

func (p *PaymentRouterValues) SetChannelAccountID(accountID string) *PaymentRouterValues {
	p.ChannelAccountID = &accountID
	return p
}

func (p *PaymentRouterValues) SetPaymentMethod(method string) *PaymentRouterValues {
	p.PaymentMethod = &method
	return p
}

func (p *PaymentRouterValues) SetCurrency(currency string) *PaymentRouterValues {
	p.Currency = &currency
	return p
}

func (p *PaymentRouterValues) SetMinAmount(amount decimal.Decimal) *PaymentRouterValues {
	p.MinAmount = &amount
	return p
}

func (p *PaymentRouterValues) SetMaxAmount(amount decimal.Decimal) *PaymentRouterValues {
	p.MaxAmount = &amount
	return p
}

func (p *PaymentRouterValues) SetAmountRange(min, max decimal.Decimal) *PaymentRouterValues {
	p.MinAmount = &min
	p.MaxAmount = &max
	return p
}

func (p *PaymentRouterValues) SetPriority(priority int) *PaymentRouterValues {
	p.Priority = &priority
	return p
}

func (p *PaymentRouterValues) SetStatus(status string) *PaymentRouterValues {
	p.Status = &status
	return p
}

func (p *PaymentRouterValues) SetEffectiveTime(time int64) *PaymentRouterValues {
	p.EffectiveTime = &time
	return p
}

func (p *PaymentRouterValues) SetExpireTime(time int64) *PaymentRouterValues {
	p.ExpireTime = &time
	return p
}

func (p *PaymentRouterValues) SetTimeRange(effective, expire int64) *PaymentRouterValues {
	p.EffectiveTime = &effective
	p.ExpireTime = &expire
	return p
}

func (p *PaymentRouterValues) SetRegion(region string) *PaymentRouterValues {
	p.Region = &region
	return p
}

func (p *PaymentRouterValues) SetRemark(remark string) *PaymentRouterValues {
	p.Remark = &remark
	return p
}

// SetValues 更新PaymentRouterValues中的非空值
func (p *PaymentRouterValues) SetValues(values *PaymentRouterValues) {
	if values == nil {
		return
	}

	if values.Name != nil && *values.Name != "" {
		p.Name = values.Name
	}
	if values.ChannelCode != nil && *values.ChannelCode != "" {
		p.ChannelCode = values.ChannelCode
	}
	if values.ChannelAccountID != nil && *values.ChannelAccountID != "" {
		p.ChannelAccountID = values.ChannelAccountID
	}
	if values.PaymentMethod != nil && *values.PaymentMethod != "" {
		p.PaymentMethod = values.PaymentMethod
	}
	if values.Currency != nil && *values.Currency != "" {
		p.Currency = values.Currency
	}
	if values.MinAmount != nil {
		p.MinAmount = values.MinAmount
	}
	if values.MaxAmount != nil {
		p.MaxAmount = values.MaxAmount
	}
	if values.Priority != nil && *values.Priority != 0 {
		p.Priority = values.Priority
	}
	if values.Status != nil && *values.Status != "" {
		p.Status = values.Status
	}
	if values.EffectiveTime != nil {
		p.EffectiveTime = values.EffectiveTime
	}
	if values.ExpireTime != nil {
		p.ExpireTime = values.ExpireTime
	}
	if values.Region != nil && *values.Region != "" {
		p.Region = values.Region
	}
	if values.Remark != nil && *values.Remark != "" {
		p.Remark = values.Remark
	}
	if values.UpdatedAt != nil && *values.UpdatedAt > 0 {
		p.UpdatedAt = values.UpdatedAt
	}
}

// 路由匹配相关方法

// IsActive 检查路由是否激活
func (p *PaymentRouterValues) IsActive() bool {
	return p.Status != nil && *p.Status == protocol.StatusActive
}

// IsInAmountRange 检查金额是否在路由范围内
func (p *PaymentRouterValues) IsInAmountRange(currency string, amount decimal.Decimal) bool {
	// 检查货币类型：支持通配符 * 或精确匹配
	if p.GetCurrency() == "" || p.GetCurrency() == "*" || p.GetCurrency() != currency {
		return false
	}
	// 检查最小金额
	if p.MinAmount != nil && p.MinAmount.GreaterThan(decimal.Zero) && amount.LessThan(*p.MinAmount) {
		return false
	}

	// 检查最大金额
	if p.MaxAmount != nil && p.MaxAmount.GreaterThan(decimal.Zero) && amount.GreaterThan(*p.MaxAmount) {
		return false
	}

	return true
}

// IsInTimeRange 检查是否在有效时间范围内
func (p *PaymentRouterValues) IsInTimeRange(timestamp int64) bool {
	// 检查生效时间
	if p.EffectiveTime != nil && timestamp < *p.EffectiveTime {
		return false
	}

	// 检查失效时间
	if p.ExpireTime != nil && timestamp > *p.ExpireTime {
		return false
	}

	return true
}

// IsEffective 检查路由是否当前有效
func (p *PaymentRouterValues) IsEffective() bool {
	if !p.IsActive() {
		return false
	}

	now := utils.TimeNowMilli()
	return p.IsInTimeRange(now)
}

// 数据库查询方法

// GetPaymentRouter 根据条件获取匹配的支付路由
// 按优先级排序，返回最高优先级的路由
func GetPaymentRouter(paymentMethod, currency string, amount decimal.Decimal) (*PaymentRouters, error) {
	var router PaymentRouters

	query := DB.Where("status = ? AND payment_method = ? AND currency = ?",
		protocol.StatusActive, paymentMethod, currency)

	// 金额范围检查
	query = query.Where("(min_amount IS NULL OR min_amount <= ?) AND (max_amount IS NULL OR max_amount >= ?)",
		amount, amount)

	// 时间范围检查
	now := utils.TimeNowMilli()
	query = query.Where("(effective_time IS NULL OR effective_time <= ?) AND (expire_time IS NULL OR expire_time >= ?)",
		now, now)

	// 按优先级排序，获取第一个匹配的路由
	err := query.Order("priority DESC").First(&router).Error
	if err != nil {
		return nil, err
	}

	return &router, nil
}

// GetPaymentRouterWithRegion 根据条件获取激活状态的支付路由列表
// 只返回状态为激活的记录，不做任何业务逻辑判断
func GetPaymentRouterWithRegion() ([]*PaymentRouters, error) {
	var routers []*PaymentRouters

	query := DB.Where("status = ?", protocol.StatusActive)

	// 时间范围检查
	now := utils.TimeNowMilli()
	query = query.Where("(effective_time IS NULL OR effective_time <= ?) AND (expire_time IS NULL OR expire_time >= ?)",
		now, now)

	// 按优先级排序，获取所有匹配的路由
	err := query.Order("priority DESC").Find(&routers).Error
	if err != nil {
		return nil, err
	}

	return routers, nil
}

// GetActiveRouters 获取所有活跃路由
func GetActiveRouters() ([]*PaymentRouters, error) {
	var routers []*PaymentRouters

	err := DB.Where("status = ?", protocol.StatusActive).
		Order("priority DESC, payment_method ASC").
		Find(&routers).Error
	if err != nil {
		return nil, err
	}

	return routers, nil
}

// GetRoutersByChannelAccountID 根据渠道账户ID获取路由
func GetRoutersByChannelAccountID(channelAccountID string) ([]*PaymentRouters, error) {
	var routers []*PaymentRouters

	err := DB.Where("channel_account_id = ?", channelAccountID).
		Order("priority DESC").
		Find(&routers).Error
	if err != nil {
		return nil, err
	}

	return routers, nil
}

// GetRouterByID 根据路由ID获取路由
func GetRouterByID(routerID string) (*PaymentRouters, error) {
	var router PaymentRouters

	err := DB.Where("router_id = ?", routerID).First(&router).Error
	if err != nil {
		return nil, err
	}

	return &router, nil
}

// UpdateRouterValues 更新路由配置
func UpdateRouterValues(tx *gorm.DB, router *PaymentRouters, values *PaymentRouterValues) error {
	defer func() {
		router.SetValues(values)
	}()
	if err := tx.Model(router).UpdateColumns(values).Error; err != nil {
		return err
	}
	return nil
}

// GetAllRouters 获取所有路由（分页）
func GetAllRouters(offset, limit int) ([]*PaymentRouters, int64, error) {
	var routers []*PaymentRouters
	var total int64

	// 获取总数
	if err := DB.Model(&PaymentRouters{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := DB.Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&routers).Error
	if err != nil {
		return nil, 0, err
	}

	return routers, total, nil
}

// GetAvailablePaymentMethods 根据条件获取可用的支付方式列表
// 按优先级高到低排序，返回支付方式码列表
func GetAvailablePaymentMethods(currency string, amount decimal.Decimal) ([]string, error) {
	var routers []*PaymentRouters
	// 查询所有符合条件的路由记录
	query := DB.Where("status = ? ", protocol.StatusActive)
	// 时间范围检查
	now := utils.TimeNowMilli()
	query = query.Where("(effective_time IS NULL OR effective_time <= ?) AND (expire_time IS NULL OR expire_time >= ?)",
		now, now)

	// 按优先级排序（数字越大优先级越高）
	err := query.Order("priority DESC, created_at DESC").Find(&routers).Error
	if err != nil {
		return nil, err
	}

	// 用于存储已添加的支付方式，避免重复
	methodSet := make(map[string]bool)
	var methods []string

	// 轮询处理每个路由记录
	for _, router := range routers {
		// 检查金额范围
		if !router.IsInAmountRange(currency, amount) {
			continue // 金额不在范围内，跳过
		}

		paymentMethod := router.GetPaymentMethod()
		if paymentMethod == "*" || paymentMethod == "" {
			// 通配符情况，需要从PaymentChannels表获取支持的支付方式
			channelAccountID := router.GetChannelAccountID()
			if channelAccountID != "" {
				channelMethods, err := GetPaymentMethodsFromChannel(channelAccountID)
				if err != nil {
					continue // 出错时跳过这个路由
				}
				// 添加渠道支持的所有支付方式
				for _, method := range channelMethods {
					if method != "" && !methodSet[method] {
						methodSet[method] = true
						methods = append(methods, method)
					}
				}
			}
			continue
		}
		// 具体的支付方式
		if !methodSet[paymentMethod] {
			methodSet[paymentMethod] = true
			methods = append(methods, paymentMethod)
		}
	}

	return methods, nil
}

// GetPaymentMethodsFromChannel 根据渠道账户ID获取支持的支付方式
func GetPaymentMethodsFromChannel(channelAccountID string) ([]string, error) {
	var channel PaymentChannels

	err := DB.Where("account_id = ? AND status = ?", channelAccountID, protocol.StatusActive).
		First(&channel).Error
	if err != nil {
		return nil, err
	}

	return channel.GetPaymentMethods(), nil
}
