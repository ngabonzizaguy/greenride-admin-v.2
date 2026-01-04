package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaymentMethods 获取支付方式列表接口
// @Summary 获取支付方式列表
// @Description 根据货币类型和金额查询支持的支付方式列表
// @Tags Api,支付
// @Accept json
// @Produce json
// @Param request body protocol.PaymentMethodsRequest true "支付方式查询请求"
// @Success 200 {object} protocol.Result{data=object} "获取成功，返回支付方式列表"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /payment/methods [post]
func (a *Api) GetPaymentMethods(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析 JSON 请求体
	var req protocol.PaymentMethodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	user := GetUserFromContext(c)
	if !user.IsPassenger() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	// 调用服务层
	paymentService := services.GetPaymentService()
	methods := paymentService.GetAvailablePaymentMethods(&req)

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(methods, lang))
}

// CancelPayment 取消支付接口
// @Summary 取消支付
// @Description 用户主动取消订单支付
// @Tags Api,支付
// @Accept json
// @Produce json
// @Param request body protocol.CancelPaymentRequest true "取消支付请求"
// @Success 200 {object} protocol.Result "取消成功"
// @Failure 200 {object} protocol.Result "取消失败"
// @Security BearerAuth
// @Router /payment/cancel [post]
func (a *Api) CancelPayment(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析 JSON 请求体
	var req protocol.CancelPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetUserFromContext(c)
	if !user.IsPassenger() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	req.UserID = user.UserID
	// 调用服务层
	paymentService := services.GetPaymentService()
	errCode := paymentService.CancelPayment(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
