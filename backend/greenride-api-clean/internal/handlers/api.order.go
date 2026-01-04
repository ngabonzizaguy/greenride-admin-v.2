package handlers

import (
	"net/http"

	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 通用订单管理接口 (支持行程、外卖、网购等多种订单类型)
// =============================================================================

// =============================================================================
// 网约车订单管理接口
// =============================================================================

// EstimateOrder 预估订单费用
// @Summary 预估订单费用
// @Description 根据起点终点和订单类型预估费用、距离和时间，支持价格快照和规则引擎。统一接口支持网约车、外卖、网购等多种订单类型
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.EstimateRequest true "预估请求，包含订单类型、用户信息、起终点等"
// @Success 200 {object} protocol.Result{data=protocol.OrderPrice} "预估成功，返回价格信息"
// @Failure 200 {object} protocol.Result "预估失败"
// @Security BearerAuth
// @Router /order/estimate [post]
func (a *Api) EstimateOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	user := GetUserFromContext(c)
	req.UserID = user.UserID
	// 设置默认订单类型为网约车
	if req.OrderType == "" {
		req.OrderType = protocol.RideOrder
	}
	// 调用整合后的预估服务
	response, errCode := services.GetOrderService().EstimateOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 创建新的订单，支持网约车等类型
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.CreateOrderRequest true "创建订单请求"
// @Success 200 {object} protocol.Result{data=protocol.Order} "创建成功"
// @Failure 200 {object} protocol.Result "创建失败"
// @Security BearerAuth
// @Router /order/create [post]
func (a *Api) CreateOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	var req protocol.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	user := GetUserFromContext(c)
	if !user.IsPassenger() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	req.UserID = user.UserID
	// 创建网约车订单（主表+详情表）
	protocolOrder, errCode := services.GetOrderService().CreateOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(protocolOrder, lang))
}

// GetOrders 获取订单列表
// @Summary 获取订单列表
// @Description 获取用户的订单列表，支持分页和类型过滤
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.UserRidesRequest true "获取订单列表请求"
// @Success 200 {object} protocol.PageResult "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /orders [post]
func (a *Api) GetOrders(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 获取当前用户
	user := GetUserFromContext(c)
	var req protocol.UserRidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	req.UserID = user.UserID
	req.UserType = user.GetUserType()
	// 获取订单列表
	orders, total := services.GetOrderService().GetOrdersByUser(&req)

	result := protocol.NewPageResult(orders, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// GetOrderDetail 获取订单详情
// @Summary 获取订单详情
// @Description 获取指定订单的详细信息
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.OrderIDRequest true "订单ID请求"
// @Success 200 {object} protocol.Result{data=protocol.Order} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /order/detail [post]
func (a *Api) GetOrderDetail(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)

	var req protocol.OrderIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 获取主订单
	order := services.GetOrderService().GetOrderInfoByID(req.OrderID)
	if order == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.DatabaseError, lang))
		return
	}

	// 权限检查
	userType := user.GetUserType()
	if userType == protocol.UserTypeDriver && order.ProviderID != "" && order.ProviderID != user.UserID {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	if userType == protocol.UserTypePassenger && order.UserID != user.UserID {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(order, lang))
}

// AcceptOrder 接受订单
// @Summary 接受订单
// @Description 服务提供者接受指定的订单
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderAcceptRequest true "接受订单请求"
// @Success 200 {object} protocol.Result "接受成功"
// @Failure 200 {object} protocol.Result "接受失败"
// @Security BearerAuth
// @Router /order/accept [post]
func (a *Api) AcceptOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	// 获取当前用户
	user := GetUserFromContext(c)
	req.UserID = user.UserID
	// 接受订单
	errCode := services.GetOrderService().AcceptOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// RejectOrder 拒绝订单
// @Summary 拒绝订单
// @Description 司机拒绝派单，支持枚举原因和自定义原因
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderActionRequest true "拒绝订单请求"
// @Success 200 {object} protocol.Result "拒绝成功"
// @Failure 200 {object} protocol.Result "拒绝失败"
// @Security BearerAuth
// @Router /order/reject [post]
func (a *Api) RejectOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 获取当前用户
	user := GetUserFromContext(c)
	req.UserID = user.UserID

	// 验证拒绝原因
	if req.RejectReason == "" {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, "reject_reason is required"))
		return
	}

	// 调用派单服务处理拒绝逻辑
	errCode := services.GetOrderService().RejectOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// ArrivedOrder 司机到达
// @Summary 司机到达
// @Description 司机标记已到达接客地点，更新订单状态为司机已到达
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderActionRequest true "司机到达请求"
// @Success 200 {object} protocol.Result "到达成功"
// @Failure 200 {object} protocol.Result "到达失败"
// @Security BearerAuth
// @Router /order/arrived [post]
func (a *Api) ArrivedOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	var req protocol.OrderActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	errCode := services.GetOrderService().ArrivedOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// StartOrder 开始订单
// @Summary 开始订单
// @Description 司机开始执行订单，标记订单状态为进行中
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderActionRequest true "开始订单请求"
// @Success 200 {object} protocol.Result "开始成功"
// @Failure 200 {object} protocol.Result "开始失败"
// @Security BearerAuth
// @Router /order/start [post]
func (a *Api) StartOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	var req protocol.OrderActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	errCode := services.GetOrderService().StartOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// FinishOrder 完成订单
// @Summary 完成订单
// @Description 司机完成订单执行，标记订单状态为已完成
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderActionRequest true "完成订单请求"
// @Success 200 {object} protocol.Result "完成成功"
// @Failure 200 {object} protocol.Result "完成失败"
// @Security BearerAuth
// @Router /order/finish [post]
func (a *Api) FinishOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	var req protocol.OrderActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	errCode := services.GetOrderService().FinishOrder(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 取消指定的订单
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.OrderCancelAPIRequest true "取消请求"
// @Success 200 {object} protocol.Result "取消成功"
// @Failure 200 {object} protocol.Result "取消失败"
// @Security BearerAuth
// @Router /order/cancel [post]
func (a *Api) CancelOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	// 获取当前用户
	user := GetUserFromContext(c)
	// 设置用户ID
	req.UserID = user.UserID

	// 取消订单
	errCode := services.GetOrderService().CancelOrderRequest(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// =============================================================================
// 附近订单和支付确认接口
// =============================================================================

// GetNearbyOrders 获取附近订单
// @Summary 获取附近的待接订单
// @Description 司机或服务提供者根据当前地理位置获取附近的待接订单列表，支持按距离、订单类型等条件进行筛选
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.GetNearbyOrdersRequest true "获取附近订单请求"
// @Success 200 {object} protocol.Result{data=[]protocol.Order} "获取成功，返回附近订单列表"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "用户未认证"
// @Failure 403 {object} protocol.Result "权限不足，仅司机或服务提供者可访问"
// @Failure 500 {object} protocol.Result "服务器内部错误"
// @Security BearerAuth
// @Router /order/nearby [post]
func (a *Api) GetNearbyOrders(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.GetNearbyOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	// 获取当前用户
	user := GetUserFromContext(c)
	// 检查用户类型
	if !user.IsDriver() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	req.OrderType = protocol.RideOrder
	// 使用 OrderService 获取附近订单
	response, errCode := services.GetOrderService().GetNearbyOrders(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

// OrderCashReceived 确认现金收款
// @Summary 确认现金收款
// @Description 司机确认已收到现金支付
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.OrderPaymentRequest true "确认现金收款请求"
// @Success 200 {object} protocol.Result "确认成功"
// @Failure 200 {object} protocol.Result "确认失败"
// @Security BearerAuth
// @Router /order/cash/received [post]
func (a *Api) OrderCashReceived(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetUserFromContext(c)
	req.UserID = user.UserID
	req.PaymentMethod = protocol.PaymentMethodCash

	// 使用 OrderService 确认收款
	_, errCode := services.GetOrderService().OrderPayment(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// OrderPayment 处理订单支付
// @Summary 处理订单支付
// @Description 处理用户订单支付，支持多种支付方式。沙盒模式下不实际请求支付渠道，直接返回成功。
// @Tags Api,订单,支付
// @Accept json
// @Produce json
// @Param request body protocol.OrderPaymentRequest true "订单支付请求"
// @Success 200 {object} protocol.Result{data=protocol.ChannelResult} "支付处理成功"
// @Failure 200 {object} protocol.Result "支付处理失败"
// @Security BearerAuth
// @Router /order/payment [post]
func (a *Api) OrderPayment(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetUserFromContext(c)
	req.UserID = user.UserID

	// 使用 OrderService 处理订单支付
	result, errCode := services.GetOrderService().OrderPayment(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}
