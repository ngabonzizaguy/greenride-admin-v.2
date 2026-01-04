package protocol

// FCMMessageData FCM消息数据
type FCMMessageData struct {
	MsgType string `json:"msg_type"` // 消息类型
}

// FCMMessage FCM消息结构
type FCMMessage struct {
	MessageID   string            `json:"message_id,omitempty"`   // 消息ID (单条消息时使用)
	BatchID     string            `json:"batch_id,omitempty"`     // 批次ID (用于标识一组消息)
	UserID      string            `json:"user_id,omitempty"`      // 目标用户ID
	Token       string            `json:"token,omitempty"`        // 单个设备token
	Tokens      []string          `json:"tokens,omitempty"`       // 多个设备token
	Topic       string            `json:"topic,omitempty"`        // 主题
	Title       string            `json:"title"`                  // 标题
	Body        string            `json:"body"`                   // 内容
	Data        map[string]string `json:"data,omitempty"`         // 自定义数据
	ImageURL    string            `json:"image_url,omitempty"`    // 图片URL
	ClickAction string            `json:"click_action,omitempty"` // 点击动作
	ErrorMsg    string            `json:"error_msg,omitempty"`    // 错误信息
}

// FCMTokenRegisterRequest FCM Token注册请求
type FCMTokenRegisterRequest struct {
	Token    string `json:"token" binding:"required"`                           // FCM Token (必填)
	Platform string `json:"platform" binding:"omitempty,oneof=ios android web"` // 平台 (可选: ios, android, web)
	DeviceID string `json:"device_id" binding:"omitempty"`                      // 设备ID
	AppID    string `json:"app_id" binding:"omitempty"`                         // 应用标识 (可选)
}

// FCMTokenListResponse FCM Token列表响应
type FCMTokenListResponse struct {
	Tokens []string `json:"tokens"` // Token列表
	Count  int      `json:"count"`  // Token数量
}

// FCMSendMessageRequest FCM发送消息请求（仅内部使用，不暴露API）
type FCMSendMessageRequest struct {
	UserID      string            `json:"user_id,omitempty"`      // 目标用户ID
	Token       string            `json:"token,omitempty"`        // 目标FCM Token
	Tokens      []string          `json:"tokens,omitempty"`       // 多个FCM Token
	Topic       string            `json:"topic,omitempty"`        // 主题
	Title       string            `json:"title"`                  // 消息标题 (必填)
	Body        string            `json:"body"`                   // 消息内容 (必填)
	Data        map[string]string `json:"data,omitempty"`         // 自定义数据 (可选)
	ImageURL    string            `json:"image_url,omitempty"`    // 图片URL (可选)
	ClickAction string            `json:"click_action,omitempty"` // 点击动作 (可选)
}

// FCMSendMessageResponse FCM发送消息响应
type FCMSendMessageResponse struct {
	Success      bool   `json:"success"`                 // 是否成功
	MessageID    string `json:"message_id,omitempty"`    // 消息ID (单条消息)
	SuccessCount int    `json:"success_count,omitempty"` // 成功数量 (多条消息)
	FailureCount int    `json:"failure_count,omitempty"` // 失败数量 (多条消息)
	Message      string `json:"message"`                 // 响应消息
}

// FCM消息类型常量
const (
	FCMMessageTypeOrder     = "order"     // 订单相关
	FCMMessageTypeDriver    = "driver"    // 司机相关
	FCMMessageTypePayment   = "payment"   // 支付相关
	FCMMessageTypePromotion = "promotion" // 促销相关
	FCMMessageTypeSystem    = "system"    // 系统通知
	FCMMessageTypeEmergency = "emergency" // 紧急通知
)

// FCM点击动作常量
const (
	FCMClickActionOpenApp     = "OPEN_APP"     // 打开应用
	FCMClickActionOpenOrder   = "OPEN_ORDER"   // 打开订单
	FCMClickActionOpenProfile = "OPEN_PROFILE" // 打开个人资料
	FCMClickActionOpenWallet  = "OPEN_WALLET"  // 打开钱包
)
