package services

import (
	"fmt"
	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"slices"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// RouterInfo 支付路由信息
type RouterInfo struct {
	ChannelService   PaymentChannel `json:"-"`                  // 支付渠道服务实例（不序列化）
	ChannelCode      string         `json:"channel_code"`       // 渠道代码
	ChannelAccountID string         `json:"channel_account_id"` // 渠道账户ID
	RouterID         string         `json:"router_id"`          // 路由ID
	Priority         int            `json:"priority"`           // 优先级
	PaymentMethod    string         `json:"payment_method"`     // 支付方式
	Currency         string         `json:"currency"`           // 货币类型
	Region           string         `json:"region"`             // 地区限制
}

// GetChannel 获取支付渠道服务实例
func (r *RouterInfo) GetChannel() PaymentChannel {
	return r.ChannelService
}

// IsValid 检查路由信息是否有效
func (r *RouterInfo) IsValid() bool {
	return r.ChannelService != nil &&
		r.ChannelCode != "" &&
		r.ChannelAccountID != "" &&
		r.RouterID != ""
}

// PaymentService 支付服务
type PaymentService struct {
	config *config.PaymentConfig
}

var (
	paymentServiceInstance *PaymentService
	paymentServiceOnce     sync.Once
)

// GetPaymentService 获取现金支付服务单例
func GetPaymentService() *PaymentService {
	if paymentServiceInstance == nil {
		SetupPaymentService()
	}
	return paymentServiceInstance
}

// SetupPaymentService 设置现金支付服务
func SetupPaymentService() {
	paymentServiceOnce.Do(func() {
		paymentServiceInstance = &PaymentService{
			config: config.Get().Payment,
		}
	})
}

// NewPaymentService 创建现金支付服务实例（用于测试）
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// OrderPayment 处理订单支付
// 入参是 models.Order，返回值是支付渠道结果和错误码
// 处理过程：
// 1. 验证订单和支付信息
// 2. 检查是否已有支付记录
// 3. 创建支付记录
// 4. 调用支付渠道进行支付
// 5. 返回支付结果
func (s *PaymentService) OrderPayment(req *ChannelPaymentRequest) (*protocol.ChannelResult, protocol.ErrorCode) {
	order := req.Order
	if order == nil {
		return nil, protocol.InvalidParams
	}
	user := req.User
	if user == nil {
		user = models.GetUserByID(order.GetUserID())
	}
	if user == nil {
		return nil, protocol.InvalidParams
	}
	// 1. 验证订单信息
	// 检查订单状态
	if order.GetStatus() != protocol.StatusTripEnded && order.GetStatus() != protocol.StatusPending {
		return nil, protocol.InvalidOrderStatus
	}

	// 检查订单金额
	paymentAmount := order.GetPaymentAmount()
	if paymentAmount.IsZero() || paymentAmount.IsNegative() {
		return nil, protocol.InvalidParams
	}

	// 检查支付方式
	if req.PaymentMethod == "" {
		req.PaymentMethod = order.GetPaymentMethod()
	}
	if req.PaymentMethod == "" {
		return nil, protocol.PaymentMethodRequired
	}
	//获取支付配置
	cfg := s.config
	// 2. 检查是否已有非失败的支付记录
	payment := models.GetNotFailedPaymentByOrderID(order.OrderID)
	if payment != nil {
		// For MoMo: allow re-initiation when an existing pending request has expired,
		// otherwise users keep getting the same stale pending record with no fresh prompt.
		if req.PaymentMethod == protocol.PaymentMethodMomo &&
			payment.GetPaymentMethod() == protocol.PaymentMethodMomo &&
			payment.GetStatus() == protocol.StatusPending {
			now := utils.TimeNowMilli()
			expiredAt := payment.GetExpiredAt()
			isExpired := expiredAt > 0 && now >= expiredAt
			if isExpired {
				expiredMsg := "Previous pending MoMo request expired; creating a new request"
				recycleValues := &models.PaymentValues{}
				recycleValues.
					SetStatus(protocol.StatusFailed).
					SetResCode(protocol.ResCodeRequestTimeout).
					SetResMsg(expiredMsg)
				if err := models.UpdatePaymentValues(models.DB, payment, recycleValues); err != nil {
					log.Get().Errorf("过期MoMo支付记录标记失败: order_id=%s payment_id=%s error=%v", order.OrderID, payment.PaymentID, err)
					return nil, protocol.DatabaseError
				}
				log.Get().Infof("过期MoMo支付记录已回收: order_id=%s old_payment_id=%s", order.OrderID, payment.PaymentID)
				payment = nil
			}
		}
	}
	if payment != nil {
		if req.PaymentMethod == protocol.PaymentMethodCash {
			if payment.GetStatus() == protocol.StatusPending {
				return nil, protocol.OnlinePaymentPending
			}
		}
		// 如果已有支付记录且不是失败状态，直接返回相应状态
		return &protocol.ChannelResult{
			Status:           payment.GetStatus(),
			ChannelStatus:    payment.GetChannelStatus(),
			ResCode:          payment.GetResCode(),
			ResMsg:           payment.GetResMsg(),
			OrderType:        payment.GetOrderType(),
			PaymentID:        payment.PaymentID,
			ChannelCode:      payment.GetChannelCode(),
			ChannelPaymentID: payment.GetChannelPaymentID(),
			RedirectURL:      payment.GetRedirectURL(),
		}, protocol.Success
	}
	if req.Phone == "" {
		req.Phone = user.GetPhone()
	}
	if req.Email == "" {
		req.Email = user.GetEmail()
	}
	if req.AccountName == "" {
		req.AccountName = user.GetUsername()
	}

	// 3. 创建新的支付记录
	payment = models.NewPayment()
	payment.SetOrderID(order.OrderID).
		SetOrderType(order.GetOrderType()).
		SetOrderSku(fmt.Sprintf("%v Order [%v]%v", order.GetOrderType(), order.GetCurrency(), order.GetPaymentAmount())).
		SetUserID(order.GetUserID()).
		SetPaymentMethod(req.PaymentMethod).
		SetStatus(protocol.StatusPending).
		SetCurrency(order.GetCurrency()).
		SetPhone(req.Phone).
		SetEmail(req.Email).
		SetAccountNo(req.AccountNo).
		SetAccountName(req.AccountName).
		SetAmount(paymentAmount).
		SetReturnURL(fmt.Sprintf("%v?user_id=%v&checkout_id=%v", cfg.ReturnURL, order.GetUserID(), s.GetCheckoutID(payment)))
	// 5. 创建支付记录
	if err := models.DB.Create(payment).Error; err != nil {
		log.Get().Errorf(
			"创建支付记录失败: order_id=%s payment_method=%s user_id=%s amount=%s currency=%s error=%v",
			order.OrderID,
			req.PaymentMethod,
			order.GetUserID(),
			paymentAmount.String(),
			order.GetCurrency(),
			err,
		)
		return nil, protocol.DatabaseError
	}

	values := &models.PaymentValues{}
	// 6. 调用支付渠道进行支付
	var result *protocol.ChannelResult
	defer func() {
		result = &protocol.ChannelResult{
			Status:           payment.GetStatus(),
			ChannelStatus:    payment.GetChannelStatus(),
			ResCode:          payment.GetResCode(),
			ResMsg:           payment.GetResMsg(),
			OrderType:        payment.GetOrderType(),
			ChannelCode:      payment.GetChannelCode(),
			PaymentID:        payment.PaymentID,
			ChannelPaymentID: payment.GetChannelPaymentID(),
		}
	}()
	// 检查是否为沙盒订单，沙盒订单不请求真实支付渠道，直接返回成功
	if cfg.IsSandbox() && order.IsSandbox() {
		// 沙盒模式，直接返回成功结果
		result = &protocol.ChannelResult{
			Status:           protocol.StatusSuccess,
			ChannelStatus:    protocol.StatusSuccess,
			ResCode:          protocol.ResCodeSandboxSuccess,
			OrderType:        order.GetOrderType(),
			ChannelCode:      protocol.PaymentChannelSandbox,
			PaymentID:        payment.PaymentID,
			ChannelPaymentID: utils.GenerateSandboxChannelPaymentID(),
		}
	} else if req.PaymentMethod == protocol.PaymentMethodCash {
		// 现金支付，直接返回成功结果
		result = &protocol.ChannelResult{
			Status:           protocol.StatusSuccess,
			ChannelStatus:    protocol.StatusSuccess,
			ResCode:          protocol.ResCodeCashSuccess,
			OrderType:        order.GetOrderType(),
			PaymentID:        payment.PaymentID,
			ChannelCode:      protocol.PaymentChannelCash,
			ChannelPaymentID: utils.GenerateCashChannelPaymentID(),
		}
	} else {
		// 4. 使用支付路由获取支付渠道
		routeRequest := &protocol.PaymentRouteRequest{
			PaymentMethod: req.PaymentMethod,
			Currency:      order.GetCurrency(),
			Region:        user.GetCountryCode(), // 使用用户的国家代码作为地区信息
			Amount:        paymentAmount.String(),
		}

		routerInfo, errorCode := s.GetPaymentRouter(routeRequest)
		if errorCode != protocol.Success {
			log.Get().Errorf("获取支付路由失败: order_id=%s, payment_method=%s, currency=%s, region=%s, amount=%s, error_code=%s",
				order.OrderID, req.PaymentMethod, order.GetCurrency(), user.GetCountryCode(), paymentAmount.String(), errorCode)
			return nil, errorCode
		}

		// 设置支付记录的渠道信息
		values.SetChannelCode(routerInfo.ChannelCode).
			SetChannelAccountID(routerInfo.ChannelAccountID)

		// 获取支付渠道服务实例
		channel := routerInfo.ChannelService

		// 真实环境，调用支付渠道
		result = channel.Pay(payment)
	}

	// 7. 处理支付结果
	if result == nil {
		result = &protocol.ChannelResult{
			Status:  protocol.StatusPending,
			ResCode: protocol.ResCodePaymentFailed,
		}
	}
	result.PaymentID = payment.PaymentID
	nowtime := utils.TimeNowMilli()
	// 8. 更新支付记录状态
	values.SetStatus(result.Status).
		SetChannelStatus(result.ChannelStatus).
		SetResCode(result.ResCode).
		SetResMsg(result.ResMsg).
		SetChannelPaymentID(result.ChannelPaymentID).
		SetRedirectURL(result.RedirectURL).
		SetExpiredAt(nowtime + int64(cfg.PaymentTimeout*1000)) //支付超时
	// 如果支付完成，设置完成时间
	if result.Status == protocol.StatusSuccess {
		payment.SetCompletedAt(nowtime)
	}
	if err := models.UpdatePaymentValues(models.DB, payment, values); err != nil {
		log.Get().Errorf("update %v payment error:%v", payment.GetOrderID(), err.Error())
		return nil, protocol.DatabaseError
	}

	// 9. 处理后续操作
	go func() {
		// 记录支付日志
		logger := log.GetServiceLogger("payment")
		logger.Infof("支付完成 - 订单ID: %s, 支付ID: %s, 状态: %s, 渠道: %s, 金额: %v %s",
			order.OrderID, payment.PaymentID, payment.GetStatus(), payment.GetChannelCode(),
			payment.GetAmount(), payment.GetCurrency())

		// 如果支付成功，记录日志
		if payment.GetStatus() == protocol.StatusSuccess {
			// 支付成功信息
			paymentInfo := map[string]any{
				"order_id":    order.OrderID,
				"payment_id":  payment.PaymentID,
				"status":      payment.GetStatus(),
				"amount":      payment.GetAmount().String(),
				"currency":    payment.GetCurrency(),
				"paid_at":     payment.CompletedAt,
				"description": payment.GetDescription(),
			}

			logger.Infof("支付成功: %v", paymentInfo)
		}
	}()

	// 10. 返回结果
	return result, protocol.Success
}

func (s *PaymentService) GetCheckoutID(payment *models.Payment) string {
	cfg := s.config
	checkoutID := utils.GenerateCheckoutID()
	checkout := &protocol.Checkout{
		UserID:    payment.GetUserID(),
		PaymentID: payment.PaymentID,
		OrderID:   payment.GetOrderID(),
		ExpiredAt: time.Now().Add(time.Duration(cfg.PaymentTimeout) * time.Second).UnixMilli(),
	}
	err := models.SetObjectCache(checkoutID, checkout, time.Duration(checkout.ExpiredAt)+time.Minute) // 多加一分钟缓存
	if err != nil {
		log.Get().Errorf("设置支付缓存失败: checkout_id=%s, error=%v", checkoutID, err)
	}
	return checkoutID
}

// GetCheckoutStatus 根据checkout_id查询支付状态
func (s *PaymentService) GetCheckoutStatus(req *protocol.CheckoutStatusRequest) (*protocol.Checkout, protocol.ErrorCode) {
	// 从Redis获取checkout信息
	checkout := &protocol.Checkout{}
	err := models.GetObjectCache(req.CheckoutID, checkout)
	if err != nil {
		log.Get().Warnf("Checkout not found in Redis: checkout_id=%s, error=%v", req.CheckoutID, err)
		return nil, protocol.CheckoutNotFound
	}

	// 验证用户ID
	if checkout.UserID != req.UserID {
		log.Get().Warnf("User ID mismatch for checkout: checkout_id=%s, request_user_id=%s, checkout_user_id=%s",
			req.CheckoutID, req.UserID, checkout.UserID)
		return nil, protocol.CheckoutNotFound
	}

	// 检查是否过期
	if checkout.ExpiredAt > 0 && time.Now().UnixMilli() > checkout.ExpiredAt {
		log.Get().Warnf("Checkout expired: checkout_id=%s, expired_at=%d", req.CheckoutID, checkout.ExpiredAt)
		return nil, protocol.CheckoutExpired
	}

	// 查询最新的订单状态
	if checkout.OrderID != "" {
		order := models.GetOrderByID(checkout.OrderID)
		if order != nil {
			checkout.OrderStatus = order.GetStatus()
		}
	}

	// 查询最新的支付状态
	if checkout.PaymentID != "" {
		payment := models.GetPaymentByID(checkout.PaymentID)
		if payment != nil {
			checkout.PaymentStatus = payment.GetStatus()
		}
	}

	log.Get().Infof("Checkout status retrieved successfully: checkout_id=%s, user_id=%s, order_status=%s, payment_status=%s",
		req.CheckoutID, req.UserID, checkout.OrderStatus, checkout.PaymentStatus)

	return checkout, protocol.Success
}

// GetAvailablePaymentMethods 获取可用的支付方式列表
// 根据币种和金额查询支持的支付方式
// 任何错误情况下都返回空数组，不返回错误
func (s *PaymentService) GetAvailablePaymentMethods(req *protocol.PaymentMethodsRequest) []string {
	order := models.GetOrderByID(req.OrderID)
	if order == nil {
		return []string{}
	}
	currency := order.GetCurrency()
	amount := order.GetPaymentAmount() // 使用getter方法获取decimal.Decimal类型

	if currency == "" || amount.LessThanOrEqual(decimal.Zero) {
		return []string{}
	}

	// 调用模型层查询
	methods, err := models.GetAvailablePaymentMethods(currency, amount)
	if err != nil {
		log.Get().Errorf("查询可用支付方式失败: %v", err)
		return []string{}
	}

	if methods == nil {
		return []string{}
	}

	return methods
}

// GetPaymentRouter 根据支付参数获取合适的支付服务
// 根据支付方式、币种、金额、商户ID等参数，从路由表获取适合的支付账户ID，
// 从PaymentChannels获取相应的服务并返回路由信息
func (s *PaymentService) GetPaymentRouter(req *protocol.PaymentRouteRequest) (*RouterInfo, protocol.ErrorCode) {
	// 参数验证
	if req == nil {
		log.Get().Errorf("请求参数不能为空")
		return nil, protocol.InvalidParams
	}

	if err := req.Validate(); err != nil {
		log.Get().Errorf("参数验证失败: %v", err)
		return nil, protocol.InvalidParams
	}

	// 转换金额
	amountDecimal, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Get().Errorf("金额格式错误: %v", err)
		return nil, protocol.InvalidParams
	}

	if amountDecimal.IsNegative() {
		log.Get().Errorf("金额不能为负数: %s", req.Amount)
		return nil, protocol.InvalidParams
	}

	// 获取所有激活状态的支付路由
	routers, err := models.GetPaymentRouterWithRegion()
	if err != nil {
		log.Get().Errorf("查询支付路由失败: error=%v", err)
		return nil, protocol.NoAvailablePaymentService
	}

	if len(routers) == 0 {
		log.Get().Errorf("未找到任何激活的支付路由")
		return nil, protocol.NoAvailablePaymentService
	}

	// 在代码层进行业务逻辑筛选
	var matchedRouters []*models.PaymentRouters
	for _, router := range routers {
		// 检查支付方式：支持通配符 * 或精确匹配
		routerPaymentMethod := router.GetPaymentMethod()
		if routerPaymentMethod != "" && routerPaymentMethod != "*" && routerPaymentMethod != req.PaymentMethod {
			continue
		}

		// 如果支付方式是通配符，需要检查渠道是否支持该支付方式
		if routerPaymentMethod == "*" || routerPaymentMethod == "" {
			if channelAccountID := router.GetChannelAccountID(); channelAccountID != "" {
				// 获取渠道支持的支付方式列表
				channelMethods, err := models.GetPaymentMethodsFromChannel(channelAccountID)
				if err != nil {
					log.Get().Warnf("获取渠道支付方式失败，跳过路由: channel_account_id=%s, router_id=%s, error=%v",
						channelAccountID, router.RouterID, err)
					continue
				}
				// 检查请求的支付方式是否在渠道支持的列表中
				if !slices.Contains(channelMethods, req.PaymentMethod) {
					continue // 渠道不支持该支付方式，跳过
				}
			}
		}

		// 检查地区：支持通配符 * 或为空或精确匹配
		routerRegion := router.GetRegion()
		if routerRegion != "" && routerRegion != "*" && routerRegion != req.Region {
			continue
		}
		// 检查金额范围
		if !router.IsInAmountRange(req.Currency, amountDecimal) {
			continue
		}

		// 通过所有检查，添加到匹配列表
		matchedRouters = append(matchedRouters, router)
	}

	if len(matchedRouters) == 0 {
		log.Get().Errorf("未找到匹配的支付路由: method=%s, currency=%s, region=%s, amount=%s",
			req.PaymentMethod, req.Currency, req.Region, req.Amount)
		return nil, protocol.NoAvailablePaymentService
	}

	// 按优先级逐个验证路由可用性（已经按优先级排序）
	for _, router := range matchedRouters {
		// 获取渠道账户ID
		channelAccountID := router.GetChannelAccountID()
		if channelAccountID == "" {
			log.Get().Warnf("路由缺少渠道账户ID，跳过: router_id=%s", router.RouterID)
			continue
		}

		// 从PaymentChannels获取支付服务
		channel, exists := PaymentChannels[channelAccountID]
		if !exists || channel == nil {
			log.Get().Warnf("未找到支付渠道服务，跳过: channel_account_id=%s, router_id=%s", channelAccountID, router.RouterID)
			continue
		}

		// 构建路由信息
		routerInfo := &RouterInfo{
			ChannelService:   channel,
			ChannelCode:      router.GetChannelCode(),
			ChannelAccountID: channelAccountID,
			RouterID:         router.RouterID,
			Priority:         router.GetPriority(),
			PaymentMethod:    router.GetPaymentMethod(),
			Currency:         router.GetCurrency(),
			Region:           router.GetRegion(),
		}

		// 找到可用的支付服务，返回路由信息
		log.Get().Infof("成功获取支付路由: method=%s, currency=%s, region=%s, amount=%s, channel_account_id=%s, router_id=%s, priority=%d",
			req.PaymentMethod, req.Currency, req.Region, req.Amount, channelAccountID, router.RouterID, router.GetPriority())

		return routerInfo, protocol.Success
	}

	// 所有路由都不可用
	log.Get().Errorf("所有支付路由都不可用: method=%s, currency=%s, region=%s, amount=%s, matched_routers=%d",
		req.PaymentMethod, req.Currency, req.Region, req.Amount, len(matchedRouters))
	return nil, protocol.NoAvailablePaymentService
}

func (s *PaymentService) CancelPayment(req *protocol.CancelPaymentRequest) protocol.ErrorCode {
	order := models.GetOrderByID(req.OrderID)
	// 如果提供了用户ID，检查订单归属
	if order.GetUserID() != req.UserID {
		return protocol.OrderNotFound
	}
	if order.IsCompleted() {
		return protocol.OrderAlreadyCompleted
	}
	if !order.IsFinished() {
		return protocol.OrderCannotCancel
	}
	clearOrderPayment := false
	defer func() {
		if clearOrderPayment &&
			(order.GetPaymentMethod() != "" ||
				order.GetPaymentResult() != "" ||
				order.GetPaymentID() != "" ||
				order.GetPaymentRedirectURL() != "" ||
				order.GetPaymentStatus() != protocol.StatusPending) {
			values := models.OrderValues{}
			values.SetPaymentMethod("").
				SetPaymentStatus(protocol.StatusPending).
				SetPaymentID("").
				SetPaymentResult("").
				SetPaymentRedirectURL("")
			if err := models.UpdateOrder(models.DB, order, &values); err != nil {
				log.Get().Errorf("清除订单支付信息失败: order_id=%s, error=%v", order.OrderID, err)
			} else {
				log.Get().Infof("成功清除订单支付信息: order_id=%s", order.OrderID)
			}
		}
	}()
	// 查询非失败的支付记录
	payment := models.GetPaymentByID(order.GetPaymentID())
	if payment == nil {
		payment = models.GetNotFailedPaymentByOrderID(order.OrderID)
	}
	if payment == nil || payment.GetStatus() == protocol.StatusFailed || payment.GetStatus() == protocol.StatusCancelled {
		clearOrderPayment = true
		return protocol.Success
	}
	if payment.GetPaymentMethod() == protocol.PaymentMethodCash {
		return protocol.OrderCannotCancel
	}
	if payment.GetStatus() == protocol.StatusSuccess {
		return protocol.OrderCannotCancel
	}
	values := models.PaymentValues{}
	values.SetStatus(protocol.StatusCancelled).
		SetResMsg("Cancelled by user")
	if req.Reason != "" {
		values.SetResMsg(fmt.Sprintf("%v: %v", values.GetResMsg(), req.Reason))
	}
	if err := models.UpdatePaymentValues(models.DB, payment, &values); err != nil {
		log.Get().Errorf("取消支付记录失败: payment_id=%s, error=%v", payment.PaymentID, err)
		return protocol.DatabaseError
	}
	clearOrderPayment = true
	return protocol.Success
}
