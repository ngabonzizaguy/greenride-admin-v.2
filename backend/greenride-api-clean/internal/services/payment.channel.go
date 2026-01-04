package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
)

type ChannelPaymentRequest struct {
	Phone         string
	Email         string
	AccountNo     string
	AccountName   string
	PaymentMethod string
	AuthToken     string
	Order         *models.Order
	User          *models.User
}

// PaymentChannel 支付渠道接口
type PaymentChannel interface {
	// Pay 处理支付请求
	Pay(payment *models.Payment) *protocol.ChannelResult

	// Refund 处理退款请求
	Refund(payment *models.Payment) *protocol.ChannelResult

	// Status 查询支付状态
	Status(payment *models.Payment) *protocol.ChannelResult
}

// PaymentChannels 渠道服务映射
var PaymentChannels = make(map[string]PaymentChannel) // key: accountID
