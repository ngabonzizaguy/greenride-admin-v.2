package handlers

import (
	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// 订单管理相关接口
// ============================================================================

// SearchOrders 搜索订单
// SearchOrders 搜索订单
// @Summary 搜索订单
// @Description 管理员搜索订单，支持多种条件过滤
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.OrderSearchRequest true "订单搜索条件"
// @Success 200 {object} protocol.PageResult
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/orders/search [post]
func (t *Admin) SearchOrders(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 解析请求体
	var req protocol.OrderSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// 执行搜索
	orders, total, errCode := services.GetAdminOrderService().SearchOrders(&req)
	if errCode != protocol.Success {
		log.Get().Errorf("Error searching orders with error code: %s", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.OrderSearchFailed, lang))
		return
	}

	// 构建响应数据 - 使用 protocol.Order
	list := make([]*protocol.Order, len(orders))
	for i, order := range orders {
		// 使用 GetOrderInfo 方法获取完整的 protocol.Order 对象
		orderInfo := services.GetAdminOrderService().GetOrderInfo(order)
		list[i] = orderInfo
	}

	// Calculate ride stats for dashboard cards
	db := models.GetDB()
	var activeCount int64
	db.Model(&models.Order{}).
		Where("status IN ?", []string{
			protocol.StatusRequested,
			protocol.StatusAccepted,
			protocol.StatusDriverComing,
			protocol.StatusDriverArrived,
			protocol.StatusInProgress,
		}).
		Count(&activeCount)

	// Today's completed and cancelled (Rwanda timezone, UTC+2)
	loc, _ := time.LoadLocation("Africa/Kigali")
	if loc == nil {
		loc = time.FixedZone("CAT", 2*3600)
	}
	now := time.Now().In(loc)
	todayStartMs := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).UnixMilli()

	var completedToday, cancelledToday int64
	db.Model(&models.Order{}).Where("status = ? AND completed_at >= ?", protocol.StatusCompleted, todayStartMs).Count(&completedToday)
	db.Model(&models.Order{}).Where("status = ? AND cancelled_at >= ?", protocol.StatusCancelled, todayStartMs).Count(&cancelledToday)

	// 返回结果
	result := protocol.NewPageResult(list, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	result.AddAttach("params", req)
	result.AddAttach("total_count", total)
	result.AddAttach("active_count", activeCount)
	result.AddAttach("completed_today", completedToday)
	result.AddAttach("cancelled_today", cancelledToday)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// GetOrderDetail 获取订单详情
// GetOrderDetail 获取订单详情
// @Summary 获取订单详情
// @Description 管理员获取单个订单的详细信息
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.IDRequest true "订单ID"
// @Success 200 {object} protocol.Result{data=protocol.Order}
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 404 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/orders/detail [post]
func (t *Admin) GetOrderDetail(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.OrderIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	// 获取订单详情 - 使用 GetOrderInfoByID 返回 protocol.Order
	orderDetail := services.GetAdminOrderService().GetOrderInfoByID(req.OrderID)
	if orderDetail == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.OrderNotFound, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(orderDetail))
}

// CancelOrder 取消订单
// @Summary 取消订单
// @Description 管理员取消指定订单
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminOrderCancelRequest true "订单取消信息"
// @Success 200 {object} protocol.Result
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 404 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/orders/cancel [post]
func (t *Admin) CancelOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 解析请求体 - 使用专门的管理员取消订单请求结构体
	var req protocol.AdminOrderCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetAdminFromContext(c)
	// 取消订单
	errCode := services.GetAdminOrderService().CancelOrderByAdmin(req.OrderID, user.AdminID, req.Reason)
	if errCode != protocol.Success {
		log.Get().Errorf("Error cancelling order with error code: %s", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.OrderCancelFailed, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// EstimateOrder 管理员订单价格预估
// @Summary 管理员订单价格预估
// @Description 管理员进行订单价格预估，无需用户ID
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminOrderEstimateRequest true "订单预估信息"
// @Success 200 {object} protocol.Result{data=protocol.OrderPrice}
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/orders/estimate [post]
func (t *Admin) EstimateOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求体
	var req protocol.AdminOrderEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	if req.EstimateRequest == nil {
		req.EstimateRequest = &protocol.EstimateRequest{}
	}
	if req.EstimateRequest.OrderType == "" {
		req.EstimateRequest.OrderType = protocol.RideOrder // 使用正确的订单类型常量
	}
	// 与 app 同一套定价逻辑；业务错误返回 200 + 错误体，便于前端展示
	estimate, errCode := services.GetAdminOrderService().EstimateOrder(&req)
	if errCode != protocol.Success {
		log.Get().Warnf("Admin EstimateOrder error: %s", errCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(estimate))
}

// CreateOrder 管理员创建订单
// @Summary 管理员代客户创建订单
// @Description 管理员为指定用户创建订单
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminCreateOrderRequest true "创建订单请求"
// @Success 200 {object} protocol.Result{data=protocol.Order}
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/orders/create [post]
func (t *Admin) CreateOrder(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求体
	var req protocol.AdminCreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 获取管理员信息
	admin := t.GetUserFromContext(c)
	req.AdminID = admin.AdminID

	// 创建订单（业务错误如 6007 用 200 + code/msg 返回，便于前端展示）
	order, errCode := services.GetAdminOrderService().CreateOrderForUser(&req)
	if errCode != protocol.Success {
		log.Get().Warnf("Admin CreateOrder business error: %s", errCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(order))
}

// GetOrderETA 管理员获取订单 ETA（与 app 使用同一套逻辑）
// @Summary 获取订单 ETA
// @Description 获取订单的司机到达时间预估，逻辑与公众 API /order/eta 一致
// @Tags Admin,管理员-订单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.OrderETARequest true "order_id"
// @Success 200 {object} protocol.Result{data=protocol.OrderETAResponse}
// @Router /order/eta [post]
func (t *Admin) GetOrderETA(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderETARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	order := models.GetOrderByID(req.OrderID)
	if order == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.OrderNotFound, lang))
		return
	}
	// 以乘客身份调用同一套 ETA 逻辑，保证与 app 一致
	req.UserID = order.GetUserID()
	response, errCode := services.GetOrderService().GetOrderETA(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(response))
}
