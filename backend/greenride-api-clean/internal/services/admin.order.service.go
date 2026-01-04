package services

import (
	"sync"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

type AdminOrderService struct {
}

var (
	adminOrderInstance *AdminOrderService
	adminOrderOnce     sync.Once
)

func GetAdminOrderService() *AdminOrderService {
	adminOrderOnce.Do(func() {
		SetupAdminOrderService()
	})
	return adminOrderInstance
}
func SetupAdminOrderService() {
	adminOrderInstance = &AdminOrderService{}
}

// ============================================================================
// Admin功能 - 管理员专用方法
// ============================================================================

// SearchOrders 搜索订单（管理员用）
func (s *AdminOrderService) SearchOrders(req *protocol.OrderSearchRequest) ([]*models.Order, int64, protocol.ErrorCode) {
	query := models.GetDB().Model(&models.Order{})

	// 关键词搜索（按用户名和司机名搜索）
	if req.Keyword != "" {
		// 构建子查询来搜索用户名和司机名
		userSubquery := models.GetDB().Model(&models.User{}).Select("user_id").Where("username LIKE ?", "%"+req.Keyword+"%")
		driverSubquery := models.GetDB().Model(&models.User{}).Select("user_id").Where("username LIKE ? AND user_type = ?", "%"+req.Keyword+"%", "driver")

		query = query.Where("user_id IN (?) OR provider_id IN (?)", userSubquery, driverSubquery)
	}

	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}

	// 订单类型过滤
	if req.OrderType != "" {
		query = query.Where("order_type = ?", req.OrderType)
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 支付状态过滤
	if req.PaymentStatus != "" {
		query = query.Where("payment_status = ?", req.PaymentStatus)
	}

	// 用户ID过滤
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	// 服务提供者ID过滤
	if req.ProviderID != "" {
		query = query.Where("provider_id = ?", req.ProviderID)
	}

	// 日期范围过滤 (使用时间戳毫秒)
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", *req.EndDate)
	}

	// 金额范围过滤
	if req.MinAmount != nil {
		query = query.Where("total_amount >= ?", *req.MinAmount)
	}
	if req.MaxAmount != nil {
		query = query.Where("total_amount <= ?", *req.MaxAmount)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.OrderSearchFailed
	}

	// 获取订单列表
	var orders []*models.Order
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, protocol.OrderSearchFailed
	}

	return orders, total, protocol.Success
}

// CancelOrderByAdmin 管理员取消订单
func (s *AdminOrderService) CancelOrderByAdmin(orderID, adminID, reason string) protocol.ErrorCode {
	return GetOrderService().CancelOrder(orderID, adminID, reason)
}

func (s *AdminOrderService) GetOrderInfo(order *models.Order) *protocol.Order {
	return GetOrderService().GetOrderInfo(order)
}

func (s *AdminOrderService) GetOrderInfoByID(orderID string) *protocol.Order {
	order := models.GetOrderByID(orderID)
	if order == nil {
		return nil
	}
	return GetOrderService().GetOrderInfo(order)
}

// EstimateOrder 管理员订单价格预估
func (s *AdminOrderService) EstimateOrder(req *protocol.AdminOrderEstimateRequest) (*protocol.OrderPrice, protocol.ErrorCode) {
	// 直接使用嵌入的EstimateRequest，只需要设置管理员特有的默认值
	estimateReq := req.EstimateRequest

	if estimateReq.PassengerCount == 0 {
		estimateReq.PassengerCount = 1
	}
	if estimateReq.Currency == "" {
		estimateReq.Currency = protocol.CurrencyRWF // 使用默认币种
	}
	if estimateReq.ScheduledAt != 0 {
		estimateReq.ScheduledAt = utils.TimeNowMilli()
	} else {
		estimateReq.ScheduledAt = 0 // 使用当前时间
	}

	// 调用OrderService的价格预估方法
	orderPrice, errCode := GetOrderService().EstimateOrder(estimateReq)
	if errCode != protocol.Success {
		return nil, errCode
	}

	return orderPrice, protocol.Success
}

// CreateOrderForUser 管理员为用户创建订单
func (s *AdminOrderService) CreateOrderForUser(req *protocol.AdminCreateOrderRequest) (*protocol.Order, protocol.ErrorCode) {
	log.Get().Infof("CreateOrderForUser: 开始创建订单，UserID=%s, PriceID=%s", req.UserID, req.PriceID)

	// 验证用户是否存在
	user := GetUserService().GetUserByID(req.UserID)
	if user == nil {
		log.Get().Errorf("CreateOrderForUser: 用户不存在，UserID=%s", req.UserID)
		return nil, protocol.UserNotFound
	}

	// 验证用户状态
	if user.GetStatus() != protocol.StatusActive {
		log.Get().Warnf("CreateOrderForUser: 用户状态不活跃，UserID=%s, Status=%s", req.UserID, user.GetStatus())
		return nil, protocol.UserNotActive
	}

	// 验证价格快照
	pricingService := GetPriceRuleService()
	priceSnapshot, errCode := pricingService.ValidateAndLockPriceID(req.PriceID)
	if errCode != protocol.Success {
		log.Get().Errorf("CreateOrderForUser: 价格快照验证失败，PriceID=%s, ErrorCode=%s", req.PriceID, errCode)
		return nil, errCode
	}

	log.Get().Infof("CreateOrderForUser: 价格快照验证成功，快照中UserID=%s", priceSnapshot.GetUserID())

	// 构建标准的创建订单请求
	createOrderReq := &protocol.CreateOrderRequest{
		UserID:  req.UserID,
		PriceID: req.PriceID,
		Notes:   req.Notes,
	}

	log.Get().Infof("CreateOrderForUser: 调用OrderService.CreateOrder，UserID=%s, PriceID=%s", req.UserID, req.PriceID)

	// 调用标准的订单创建流程
	order, errCode := GetOrderService().CreateOrder(createOrderReq)
	if errCode != protocol.Success {
		log.Get().Errorf("CreateOrderForUser: OrderService.CreateOrder失败，ErrorCode=%s", errCode)
		return nil, errCode
	}

	// 记录管理员操作历史（异步）
	go func() {
		log.Get().Infof("Admin %s created order %s for user %s, reason: %s",
			req.AdminID, order.OrderID, req.UserID, req.AdminReason)
	}()

	log.Get().Infof("CreateOrderForUser: 订单创建成功，OrderID=%s", order.OrderID)
	return order, protocol.Success
}
