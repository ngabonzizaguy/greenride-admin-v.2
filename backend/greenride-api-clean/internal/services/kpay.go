package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"greenride/internal/config"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	KpaySandboxURL = "https://pay.esicia.com"
	KpayLiveURL    = "https://pay.esicia.rw"
)

// KPayStatusMapping KPay渠道状态映射到系统状态
// 根据KPay API文档的错误码表
var KPayStatusMapping = map[string]string{
	"0":   protocol.StatusPending, // No error. Transaction being processed
	"01":  protocol.StatusSuccess, // Successful payment
	"02":  protocol.StatusFailed,  // Payment failed
	"03":  protocol.StatusPending, // Pending transaction
	"401": protocol.StatusFailed,  // Missing authentication header
	"500": protocol.StatusFailed,  // Non HTTPS request
	"600": protocol.StatusFailed,  // Invalid username/password combination
	"601": protocol.StatusFailed,  // Invalid remote user
	"602": protocol.StatusFailed,  // Location/IP not whitelisted
	"603": protocol.StatusFailed,  // Empty parameter - missing required parameters
	"604": protocol.StatusFailed,  // Unknown retailer
	"605": protocol.StatusFailed,  // Retailer not enabled
	"606": protocol.StatusFailed,  // Error processing
	"607": protocol.StatusFailed,  // Failed mobile money transaction
	"608": protocol.StatusFailed,  // Used ref id - error uniqueness
	"609": protocol.StatusFailed,  // Unknown Payment method
	"610": protocol.StatusFailed,  // Unknown or not enabled Financial institution
	"611": protocol.StatusFailed,  // Transaction not found
}

// KPayService 为 KPay 支付渠道服务
// 实现 PaymentChannel 接口
type KPayService struct {
	config *KPayConfig
}

// KPayConfig KPay支付配置
type KPayConfig struct {
	Username     string `json:"username"`      // 基本认证用户名
	Password     string `json:"password"`      // 基本认证密码
	BaseURL      string `json:"base_url"`      // API基础URL
	ReturnURL    string `json:"return_url"`    // 异步通知回调URL
	RedirectURL  string `json:"redirect_url"`  // 支付完成后重定向URL
	LogoURL      string `json:"logo_url"`      // 结账页面Logo URL
	RetailerID   string `json:"retailer_id"`   // 零售商ID
	Timeout      int    `json:"timeout"`       // 请求超时时间(秒)
	DefaultPhone string `json:"default_phone"` // 默认手机号（当客户未提供时使用）
	DefaultEmail string `json:"default_email"` // 默认邮箱（当客户未提供时使用）
}

func (c *KPayConfig) Validate() error {
	// 根据KPay API文档，只需要Basic认证（用户名和密码）
	if c.Username == "" {
		return errors.New("KPay 用户名不能为空")
	}
	if c.Password == "" {
		return errors.New("KPay 密码不能为空")
	}
	if c.RetailerID == "" {
		return errors.New("KPay 零售商ID不能为空")
	}
	if c.Timeout <= 0 {
		c.Timeout = 30 // 默认30秒
	}
	// 只有当BaseURL为空时才设置默认值
	if c.BaseURL == "" {
		c.BaseURL = KpaySandboxURL // 默认沙盒环境，生产环境需要明确设置
	}
	// 从全局配置中获取缺失的URL配置
	if _cfg := config.Get(); _cfg != nil {
		if cfg := _cfg.KPay; cfg != nil {
			if c.LogoURL == "" && cfg.LogoURL != "" {
				c.LogoURL = cfg.LogoURL
			}
			if c.ReturnURL == "" && cfg.CallbackUrl != "" {
				c.ReturnURL = cfg.CallbackUrl
			}
			if c.RedirectURL == "" && cfg.ReturnURL != "" {
				c.RedirectURL = cfg.ReturnURL
			}
		}
	}
	return nil
}

func NewKPayService(config *models.PaymentChannels) *KPayService {
	cfg := config.Config
	if len(cfg) == 0 {
		return nil
	}
	var kpayConfig KPayConfig
	cfg.ToObject(&kpayConfig)
	return NewKPayServiceWithConfig(&kpayConfig)
}

func NewKPayServiceWithConfig(kpayConfig *KPayConfig) *KPayService {
	if kpayConfig == nil {
		return nil
	}
	if err := kpayConfig.Validate(); err != nil {
		return nil
	}
	return &KPayService{
		config: kpayConfig,
	}
}

// Pay 处理支付请求
// 实现 PaymentChannel 接口的 Pay 方法
func (s *KPayService) Pay(payment *models.Payment) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		Status:        protocol.StatusFailed,
		ChannelStatus: protocol.StatusFailed,
		OrderType:     protocol.PaymentTypePayment,
		ChannelCode:   protocol.PaymentChannelKPay,
	}
	// 准备请求数据
	if payment.GetAmount().LessThanOrEqual(decimal.Zero) {
		result.ResCode = protocol.ResCodeInvalidAmount
		return result
	}
	phone := strings.ReplaceAll(strings.ReplaceAll(payment.GetPhone(), "-", ""), "+", "")
	// 构建支付请求 - 根据KPay API文档格式
	req := protocol.MapData{
		"action":      "pay", // KPay API必需参数
		"msisdn":      phone,
		"details":     payment.GetOrderSku(),
		"refid":       payment.PaymentID,
		"amount":      payment.GetAmount().String(),
		"currency":    payment.GetCurrency(), // KPay支持RWF
		"email":       payment.GetEmail(),
		"cname":       payment.GetAccountName(),
		"cnumber":     payment.UserID,
		"retailerid":  s.config.RetailerID,
		"returl":      fmt.Sprintf("%v/%v", s.config.ReturnURL, payment.PaymentID),
		"redirecturl": payment.GetReturnURL(),
		"logourl":     s.config.LogoURL,
	}
	channelPaymentMethod := ""
	bankId := ""
	needUrl := false
	// 从元数据中提取支付方式（如果有）
	// 映射支付方式 - 根据KPay文档支持的支付方式
	switch payment.GetPaymentMethod() {
	case protocol.PaymentMethodVisa, protocol.PaymentMethodMaster, protocol.PaymentMethodCard:
		channelPaymentMethod = "cc" // Credit Card
		bankId = "000"
		needUrl = true
	case protocol.PaymentMethodMomo:
		channelPaymentMethod = "momo" // MTN Mobile Money
		bankId = "63510"              // MTN MOMO bank code
	case protocol.PaymentMethodAirtel:
		channelPaymentMethod = "momo" // Airtel Money
		bankId = "63514"              // AIRTEL MONEY bank code
	case protocol.PaymentMethodSpenn:
		channelPaymentMethod = "spenn"
		bankId = "63502" // SPENN bank code
	default:
		result.ResCode = protocol.ResCodeUnsupportedPaymentMethod
		return result
	}
	req.Set("pmethod", channelPaymentMethod)
	req.Set("bankid", bankId)
	// 调用封装的InitiatePayment方法
	resp, err := s.InitiatePayment(req)
	if err != nil {
		return &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeRequestFailed,
			OrderType:     protocol.PaymentTypePayment,
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}
	resp.Set("need_url", needUrl)
	// 使用ResolveResponse统一处理响应
	result = s.ResolveResponse(resp)
	return result
}

// Refund 处理退款请求
// 实现 PaymentChannel 接口的 Refund 方法
func (s *KPayService) Refund(payment *models.Payment) *protocol.ChannelResult {
	// 返回成功结果，具体逻辑待实现
	return &protocol.ChannelResult{
		Status:           protocol.StatusRefunded,
		ChannelStatus:    protocol.StatusSuccess,
		ResCode:          protocol.Success.GetCode(),
		OrderType:        protocol.PaymentTypeRefund,
		ChannelCode:      protocol.PaymentChannelKPay,
		ChannelPaymentID: utils.GenerateID(),
	}
}

// Status 查询支付状态
// 实现 PaymentChannel 接口的 Status 方法
func (s *KPayService) Status(payment *models.Payment) *protocol.ChannelResult {
	// 调用封装的CheckPaymentStatus方法
	result, err := s.CheckPaymentStatus(payment.PaymentID)
	if err != nil {
		return &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeRequestFailed,
			OrderType:     protocol.PaymentTypePayment,
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}
	return result
}

// ----- 以下是新增的KPay接口封装 -----

// MapStatus 映射KPay状态到系统状态
func MapStatus(kpayStatusID string) (status string) {
	// 获取系统状态
	if systemStatus, ok := KPayStatusMapping[kpayStatusID]; ok {
		return systemStatus
	}
	return protocol.StatusPending
}

// DoPost 统一处理KPay接口请求
// 封装了完整的错误处理逻辑，包括网络错误、解析错误等
func (s *KPayService) DoPost(reqData protocol.MapData, requestType string) (protocol.MapData, *protocol.ChannelResult) {
	// 记录请求开始
	fmt.Printf("[KPay] %s 请求开始: %s\n", requestType, reqData.ToJson())

	// 验证请求数据
	if len(reqData) == 0 {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeMissingFields,
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}

	// 发送HTTP请求
	body, resp, err := utils.PostJsonDataWithHeader(
		s.config.BaseURL,
		[]byte(reqData.ToJson()),
		s.GetHeaders(),
	)

	// 处理网络层错误
	if err != nil {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeRequestFailed,
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}
	// 检查HTTP状态码
	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 300 {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.GetResCodeByStatusCode(statusCode),
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}

	// 检查响应体是否为空
	if body == "" {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeInvalidResponse,
			ChannelCode:   protocol.PaymentChannelKPay,
		}
	}

	// 解析JSON响应
	respData := protocol.MapData{}
	if err := json.Unmarshal([]byte(body), &respData); err != nil {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeResponseParseFailed,
			ChannelCode:   protocol.PaymentChannelKPay,
			CallbackData:  body, // 保存原始响应用于调试
		}
	}

	// 验证响应格式的基本完整性
	if len(respData) == 0 {
		return nil, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeInvalidResponse,
			ChannelCode:   protocol.PaymentChannelKPay,
			CallbackData:  body,
		}
	}

	// 检查是否包含明显的错误响应
	if respData.Has("error") {
		return respData, &protocol.ChannelResult{
			Status:        protocol.StatusFailed,
			ChannelStatus: protocol.StatusFailed,
			ResCode:       protocol.ResCodeBusinessError,
			ChannelCode:   protocol.PaymentChannelKPay,
			CallbackData:  body,
		}
	}

	return respData, nil
}

// GetHeaders 获取KPay请求头
func (s *KPayService) GetHeaders() map[string]string {
	headers := map[string]string{
		"Authorization": s.GetBase64Auth(),
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"User-Agent":    "GreenRide-KPay-Client/1.0",
	}

	return headers
}

func (s *KPayService) GetBase64Auth() string {
	return "Basic " + utils.GetBase64(fmt.Sprintf("%s:%s", s.config.Username, s.config.Password))
}

// InitiatePayment 初始化支付
// 封装KPay支付接口
func (s *KPayService) InitiatePayment(req protocol.MapData) (protocol.MapData, error) {
	// 发送请求
	body, resp, err := utils.PostJsonDataWithHeader(s.config.BaseURL, []byte(req.ToJson()), s.GetHeaders())
	if err != nil {
		return nil, fmt.Errorf("发送支付请求失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	// 解析响应
	payResp := protocol.MapData{}
	if err := json.Unmarshal([]byte(body), &payResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return payResp, nil
}

// CheckPaymentStatus 查询支付状态
// 封装KPay状态查询接口
func (s *KPayService) CheckPaymentStatus(paymentID string) (*protocol.ChannelResult, error) {

	// 构建状态查询请求 - 根据KPay API文档
	statusReq := protocol.MapData{
		"action": "checkstatus", // KPay API必需参数
		"refid":  paymentID,
	}

	// 发送请求
	body, _, err := utils.PostJsonDataWithHeader(s.config.BaseURL, []byte(statusReq.ToJson()), s.GetHeaders())
	if err != nil {
		return nil, fmt.Errorf("发送状态查询请求失败: %v", err)
	}

	// 解析响应
	statusResp := protocol.MapData{}
	if err := json.Unmarshal([]byte(body), &statusResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 使用ResolveResponse统一处理响应
	result := s.ResolveResponse(statusResp)

	return result, nil
}

// ResolveResponse 解析KPay响应并构建ChannelResult
// 兼容支付请求、状态查询、回调等各种响应类型
func (s *KPayService) ResolveResponse(resp protocol.MapData) *protocol.ChannelResult {
	result := &protocol.ChannelResult{
		ChannelCode: protocol.PaymentChannelKPay,
		OrderType:   protocol.PaymentTypePayment,
	}

	// 处理支付请求响应（包含success字段）
	if resp.Has("success") {
		success := resp.GetInt("success")
		if success == 1 {
			// 支付请求成功
			result.Status = protocol.StatusSuccess
			if resp.GetBool("need_url") {
				result.Status = protocol.StatusPending
			}
			result.ChannelStatus = resp.Get("cod")
			result.ResCode = resp.Get("cod")
			result.ResMsg = resp.Get("reply")
			result.ChannelPaymentID = resp.Get("tid")
			result.RedirectURL = resp.Get("url")
		} else {
			// 支付请求失败
			result.Status = protocol.StatusFailed
			result.ChannelStatus = protocol.StatusFailed
			result.ResCode = fmt.Sprintf("%d", resp.GetInt("retcode"))
			result.ResMsg = resp.Get("reply")
		}
		return result
	}

	// 处理状态查询响应（包含statusid字段）
	if resp.Has("statusid") {
		result.ChannelStatus = resp.Get("statusid")
		result.ResCode = result.ChannelStatus
		result.Status = MapStatus(result.ChannelStatus)
		result.ResMsg = resp.Get("statusdesc")
		result.ChannelPaymentID = resp.Get("tid")
		result.CallbackData = resp.ToJson()
		return result
	}

	// 处理回调响应（通常包含tid和refid）
	if resp.Has("tid") {
		// 如果有明确的状态ID，使用它
		if resp.Has("statusid") {
			statusID := resp.Get("statusid")
			result.ChannelStatus = statusID
			result.ResCode = statusID
			result.Status = MapStatus(statusID)
		} else {
			// 否则假设是成功的
			result.Status = protocol.StatusSuccess
			result.ChannelStatus = "01"
			result.ResCode = "01"
		}
		result.ResMsg = resp.Get("statusdesc")
		result.ChannelPaymentID = resp.Get("tid")
		result.CallbackData = resp.ToJson()
		return result
	}

	// 处理错误响应
	if resp.Has("error") {
		result.Status = protocol.StatusFailed
		result.ChannelStatus = protocol.StatusFailed
		result.ResCode = ""
		result.ResMsg = resp.Get("error")
		return result
	}

	// 处理其他响应格式
	if resp.Has("retcode") {
		retCode := resp.GetInt("retcode")
		if retCode == 0 || retCode == 200 {
			result.Status = protocol.StatusSuccess
			result.ChannelStatus = fmt.Sprintf("%v", retCode)
			result.ResCode = result.ChannelStatus
		} else {
			result.Status = protocol.StatusFailed
			result.ChannelStatus = protocol.StatusFailed
			result.ResCode = fmt.Sprintf("%d", retCode)
		}
		result.ResMsg = resp.Get("reply")
		result.ChannelPaymentID = resp.Get("tid")
		return result
	}

	// 默认处理：如果无法识别响应格式，返回失败状态
	result.Status = protocol.StatusFailed
	result.ChannelStatus = protocol.StatusFailed
	result.ResCode = ""
	result.CallbackData = resp.ToJson()

	return result
}
