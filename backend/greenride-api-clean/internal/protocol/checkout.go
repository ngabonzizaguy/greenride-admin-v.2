package protocol

type Checkout struct {
	UserID        string `json:"user_id"`
	PaymentID     string `json:"payment_id"`
	OrderID       string `json:"order_id"`
	OrderStatus   string `json:"order_status,omitempty"`
	PaymentStatus string `json:"payment_status,omitempty"`
	ExpiredAt     int64  `json:"expired_at"`
}

// CheckoutStatusRequest checkout状态查询请求
type CheckoutStatusRequest struct {
	CheckoutID string `json:"checkout_id" binding:"required"` // checkout ID
	UserID     string `json:"user_id" binding:"required"`     // 用户ID，用于验证权限
}
