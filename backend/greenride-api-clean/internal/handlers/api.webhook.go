package handlers

import (
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// KPayWebhook 处理 KPay 支付 Webhook 通知
// @Summary 处理 KPay 支付 Webhook 通知
// @Description 接收 KPay 的支付状态回调通知，更新支付记录状态
// @Tags Api,支付,Webhook
// @Accept json
// @Produce json
// @Param payment_id path string true "支付ID"
// @Param webhook_data body protocol.MapData true "KPay Webhook 数据"
// @Success 200 {object} protocol.Result "处理成功"
// @Failure 400 {object} protocol.Result "请求数据错误"
// @Failure 404 {object} protocol.Result "支付记录不存在"
// @Failure 500 {object} protocol.Result "处理失败"
// @Router /webhook/{payment_id} [post]
func (a *Api) KPayWebhook(c *gin.Context) {
	// 获取路径参数中的 payment_id
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		log.Get().Errorf("KPay Webhook: missing payment_id in path")
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   "",
			"refid": "",
			"reply": "Missing payment_id",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 解析请求体
	var webhookData protocol.MapData
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		log.Get().Errorf("KPay Webhook: invalid JSON data for payment_id=%s, error=%v", paymentID, err)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   "",
			"refid": paymentID,
			"reply": "Invalid JSON data",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 记录 Webhook 接收日志
	log.Get().Infof("KPay Webhook received: payment_id=%s, data=%s", paymentID, webhookData.ToJson())

	// 查找支付记录
	payment := models.GetPaymentByID(paymentID)
	if payment == nil {
		log.Get().Errorf("KPay Webhook: payment not found for payment_id=%s", paymentID)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   "",
			"refid": paymentID,
			"reply": "Payment not found",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 验证支付记录的渠道类型
	if payment.GetChannelCode() != protocol.PaymentChannelKPay {
		log.Get().Errorf("KPay Webhook: invalid channel for payment_id=%s, expected=%s, actual=%s",
			paymentID, protocol.PaymentChannelKPay, payment.GetChannelCode())
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   payment.GetChannelPaymentID(),
			"refid": paymentID,
			"reply": "Invalid payment channel",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 获取 KPay 服务实例
	channelAccountID := payment.GetChannelAccountID()
	if channelAccountID == "" {
		log.Get().Errorf("KPay Webhook: missing channel_account_id for payment_id=%s", paymentID)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   payment.GetChannelPaymentID(),
			"refid": paymentID,
			"reply": "Missing channel account",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	kpayService, exists := services.PaymentChannels[channelAccountID]
	if !exists {
		log.Get().Errorf("KPay Webhook: channel service not found for account_id=%s, payment_id=%s",
			channelAccountID, paymentID)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   payment.GetChannelPaymentID(),
			"refid": paymentID,
			"reply": "Channel service not found",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	kpay, ok := kpayService.(*services.KPayService)
	if !ok {
		log.Get().Errorf("KPay Webhook: invalid service type for account_id=%s, payment_id=%s",
			channelAccountID, paymentID)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   payment.GetChannelPaymentID(),
			"refid": paymentID,
			"reply": "Invalid service type",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 调用 KPayService 的 ResolveResponse 处理 webhook 数据，只获取结果
	result := kpay.ResolveResponse(webhookData)
	if result == nil {
		log.Get().Errorf("KPay Webhook: failed to process webhook for payment_id=%s", paymentID)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   payment.GetChannelPaymentID(),
			"refid": paymentID,
			"reply": "Failed to process webhook",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 在 API 层处理更新逻辑
	values := &models.PaymentValues{}
	values.SetStatus(result.Status).
		SetChannelStatus(result.ChannelStatus).
		SetResCode(result.ResCode).
		SetResMsg(result.ResMsg).
		SetRedirectURL("")

	// 如果有渠道支付ID，更新它
	if result.ChannelPaymentID != "" {
		values.SetChannelPaymentID(result.ChannelPaymentID)
	}

	// 如果支付成功/失败，设置完成时间
	if result.Status == protocol.StatusSuccess || result.Status == protocol.StatusFailed {
		values.SetCompletedAt(utils.TimeNowMilli())
	}

	// 更新支付记录
	if err := models.UpdatePaymentValues(models.DB, payment, values); err != nil {
		log.Get().Errorf("KPay Webhook: failed to update payment for payment_id=%s, error=%v", paymentID, err)
		// 返回符合 KPay 要求的失败响应
		response := protocol.MapData{
			"tid":   result.ChannelPaymentID,
			"refid": paymentID,
			"reply": "Failed to update payment record",
		}
		c.JSON(http.StatusOK, response)
		return
	}
	go func() {
		services.GetOrderService().CheckOrderPayment(payment.GetOrderID(), payment.PaymentID)
	}()

	// 记录处理成功日志
	log.Get().Infof("KPay Webhook processed successfully: payment_id=%s, status=%s, channel_status=%s",
		paymentID, result.Status, result.ChannelStatus)

	// 异步处理后续操作
	go func() {
		logger := log.GetServiceLogger("webhook")
		logger.Infof("KPay Webhook 处理完成 - 支付ID: %s, 订单ID: %s, 状态: %s -> %s",
			paymentID, payment.GetOrderID(), payment.GetStatus(), result.Status)

		// 如果支付成功，记录成功信息
		if result.Status == protocol.StatusSuccess {
			successInfo := map[string]any{
				"payment_id":           paymentID,
				"order_id":             payment.GetOrderID(),
				"amount":               payment.GetAmount().String(),
				"currency":             payment.GetCurrency(),
				"channel_payment_id":   result.ChannelPaymentID,
				"webhook_processed_at": utils.TimeNowMilli(),
			}
			logger.Infof("KPay 支付成功回调: %v", successInfo)
		}
	}()

	// 返回成功响应（符合 KPay 要求的响应格式）
	response := protocol.MapData{
		"tid":   result.ChannelPaymentID,
		"refid": paymentID,
		"reply": "OK",
	}

	c.JSON(http.StatusOK, response)
}

// InnoPaaSWebhook handles InnoPaaS message status callbacks (OTP delivery status, etc.).
// Must return HTTP 200 within 3 seconds; InnoPaaS validates connectivity on save.
func (a *Api) InnoPaaSWebhook(c *gin.Context) {
	var body protocol.MapData
	_ = c.ShouldBindJSON(&body)
	if body != nil {
		log.Get().Infof("InnoPaaS webhook: %s", body.ToJson())
	}
	c.Status(http.StatusOK)
}
