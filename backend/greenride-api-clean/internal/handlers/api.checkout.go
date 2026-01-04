package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCheckoutStatus 根据checkout_id查询支付状态接口
// @Summary 查询checkout状态
// @Description 根据checkout_id查询订单和支付状态
// @Tags Api,支付
// @Accept json
// @Produce json
// @Param request body protocol.CheckoutStatusRequest true "checkout状态查询请求"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "查询成功，返回checkout详情"
// @Failure 200 {object} protocol.Result "查询失败"
// @Router /checkout/status [post]
func (a *Api) GetCheckoutStatus(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析 JSON 请求体
	var req protocol.CheckoutStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 调用服务层
	paymentService := services.GetPaymentService()
	checkout, errorCode := paymentService.GetCheckoutStatus(&req)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(checkout, lang))
}
