package protocol

// PaymentRouteRequest 支付路由请求参数
type PaymentRouteRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"` // 支付方式
	Currency      string `json:"currency" binding:"required"`       // 货币类型
	Region        string `json:"region"`                            // 地区/国家（可选）
	Amount        string `json:"amount" binding:"required"`         // 支付金额
}

// Validate 验证请求参数
func (r *PaymentRouteRequest) Validate() error {
	if r.PaymentMethod == "" {
		return &ValidationError{Field: "payment_method", Message: "支付方式不能为空"}
	}
	if r.Currency == "" {
		return &ValidationError{Field: "currency", Message: "货币类型不能为空"}
	}
	if r.Amount == "" {
		return &ValidationError{Field: "amount", Message: "金额不能为空"}
	}
	return nil
}
