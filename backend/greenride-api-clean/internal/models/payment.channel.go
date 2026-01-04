package models

import "greenride/internal/protocol"

// PaymentChannels 支付渠道信息
type PaymentChannels struct {
	ID             int64            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AccountID      string           `json:"account_id" gorm:"column:account_id;type:varchar(64);uniqueIndex:idx_account_id"`
	Name           string           `json:"name" gorm:"column:name;type:varchar(100)"`
	ChannelCode    string           `json:"channel_code" gorm:"column:channel_code;type:varchar(32);index:idx_channel_code"`
	PaymentMethods []string         `json:"payment_methods" gorm:"column:payment_methods;type:json;serializer:json"` // 支持的支付方式列表
	Status         string           `json:"status" gorm:"column:status;type:varchar(20);default:'inactive'"`
	Config         protocol.MapData `json:"config" gorm:"column:config;type:json;serializer:json"` // JSON格式的配置信息
	Remark         string           `json:"remark" gorm:"column:remark;type:varchar(500)"`
	CreatedAt      int64            `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      int64            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (PaymentChannels) TableName() string {
	return "t_payment_channels"
}

// GetPaymentMethods 获取支持的支付方式列表
func (p *PaymentChannels) GetPaymentMethods() []string {
	if p.PaymentMethods == nil {
		return []string{}
	}
	return p.PaymentMethods
}

// SetPaymentMethods 设置支持的支付方式列表
func (p *PaymentChannels) SetPaymentMethods(methods []string) *PaymentChannels {
	p.PaymentMethods = methods
	return p
}

// SupportsPaymentMethod 检查是否支持指定的支付方式
func (p *PaymentChannels) SupportsPaymentMethod(method string) bool {
	if p.PaymentMethods == nil {
		return false
	}
	for _, supportedMethod := range p.PaymentMethods {
		if supportedMethod == method {
			return true
		}
	}
	return false
}

// GetActiveChannels 获取所有活跃的支付渠道
func GetActiveChannels() ([]*PaymentChannels, error) {
	var channels []*PaymentChannels

	// 查询活跃的支付渠道
	err := DB.Where("status = ?", protocol.StatusActive).Find(&channels).Error
	if err != nil {
		return nil, err
	}

	return channels, nil
}

// GetActiveChannelsByPaymentMethod 根据支付方式获取活跃的支付渠道
func GetActiveChannelsByPaymentMethod(paymentMethod string) ([]*PaymentChannels, error) {
	var channels []*PaymentChannels

	// 查询活跃且支持指定支付方式的渠道
	// 使用JSON_CONTAINS查询JSON数组中是否包含特定值
	err := DB.Where("status = ? AND JSON_CONTAINS(payment_methods, ?)",
		protocol.StatusActive,
		`"`+paymentMethod+`"`).Find(&channels).Error
	if err != nil {
		return nil, err
	}

	return channels, nil
}
