package handlers

import (
	"encoding/json"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
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

// MoMoWebhook handles MTN MoMo payment callbacks
// @Summary Handle MTN MoMo payment webhook
// @Description Receives MoMo payment status callbacks and updates payment record
// @Tags Api,Payment,Webhook
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Param webhook_data body protocol.MapData true "MoMo Webhook data"
// @Success 200 {object} map[string]string "Success response"
// @Router /webhook/momo/{payment_id} [post]
func (a *Api) MoMoWebhook(c *gin.Context) {
	// Get payment_id from path
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		log.Get().Errorf("MoMo Webhook: missing payment_id in path")
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Missing payment_id"})
		return
	}

	// Parse request body
	var webhookData protocol.MapData
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		log.Get().Errorf("MoMo Webhook: invalid JSON data for payment_id=%s, error=%v", paymentID, err)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid JSON data"})
		return
	}

	log.Get().Infof("MoMo Webhook received: payment_id=%s, data=%s", paymentID, webhookData.ToJson())

	// Find payment record
	payment := models.GetPaymentByID(paymentID)
	if payment == nil {
		log.Get().Errorf("MoMo Webhook: payment not found for payment_id=%s", paymentID)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Payment not found"})
		return
	}

	// Validate channel type
	if payment.GetChannelCode() != protocol.PaymentChannelMoMo {
		log.Get().Errorf("MoMo Webhook: invalid channel for payment_id=%s, expected=%s, actual=%s",
			paymentID, protocol.PaymentChannelMoMo, payment.GetChannelCode())
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid payment channel"})
		return
	}

	// Get MoMo service instance
	channelAccountID := payment.GetChannelAccountID()
	if channelAccountID == "" {
		log.Get().Errorf("MoMo Webhook: missing channel_account_id for payment_id=%s", paymentID)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Missing channel account"})
		return
	}

	momoService, exists := services.PaymentChannels[channelAccountID]
	if !exists {
		log.Get().Errorf("MoMo Webhook: channel service not found for account_id=%s, payment_id=%s",
			channelAccountID, paymentID)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Channel service not found"})
		return
	}

	momo, ok := momoService.(*services.MoMoService)
	if !ok {
		log.Get().Errorf("MoMo Webhook: invalid service type for account_id=%s, payment_id=%s",
			channelAccountID, paymentID)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Invalid service type"})
		return
	}

	// Process webhook data
	result := momo.ResolveResponse(webhookData)
	if result == nil {
		log.Get().Errorf("MoMo Webhook: failed to process webhook for payment_id=%s", paymentID)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Failed to process webhook"})
		return
	}

	// Update payment record
	values := &models.PaymentValues{}
	values.SetStatus(result.Status).
		SetChannelStatus(result.ChannelStatus).
		SetResCode(result.ResCode).
		SetResMsg(result.ResMsg)

	if result.ChannelPaymentID != "" {
		values.SetChannelPaymentID(result.ChannelPaymentID)
	}

	if result.Status == protocol.StatusSuccess || result.Status == protocol.StatusFailed {
		values.SetCompletedAt(utils.TimeNowMilli())
	}

	if err := models.UpdatePaymentValues(models.DB, payment, values); err != nil {
		log.Get().Errorf("MoMo Webhook: failed to update payment for payment_id=%s, error=%v", paymentID, err)
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Failed to update payment"})
		return
	}

	// Trigger order status check
	go func() {
		services.GetOrderService().CheckOrderPayment(payment.GetOrderID(), payment.PaymentID)
	}()

	log.Get().Infof("MoMo Webhook processed successfully: payment_id=%s, status=%s", paymentID, result.Status)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// StripeWebhook handles Stripe webhook events
// @Summary Handle Stripe webhook events
// @Description Receives Stripe webhook events (payment_intent.succeeded, etc.) and updates payment records
// @Tags Api,Payment,Webhook
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe webhook signature"
// @Success 200 {object} map[string]bool "Success response"
// @Router /webhook/stripe [post]
func (a *Api) StripeWebhook(c *gin.Context) {
	// Read raw body for signature verification
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Get().Errorf("Stripe Webhook: cannot read body, error=%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
		return
	}

	// Get Stripe signature header
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		log.Get().Errorf("Stripe Webhook: missing Stripe-Signature header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	// Find a Stripe service to verify the webhook
	var event *stripe.Event
	var stripeService *services.StripeService

	for _, svc := range services.PaymentChannels {
		if ss, ok := svc.(*services.StripeService); ok {
			evt, err := ss.VerifyWebhookSignature(payload, signature)
			if err == nil {
				event = evt
				stripeService = ss
				break
			}
		}
	}

	if event == nil {
		log.Get().Errorf("Stripe Webhook: invalid signature or no Stripe service configured")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	log.Get().Infof("Stripe Webhook received: event_type=%s, event_id=%s", event.Type, event.ID)

	// Handle specific event types
	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed", "payment_intent.canceled":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Get().Errorf("Stripe Webhook: invalid PaymentIntent payload, error=%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		// Get payment_id from metadata
		paymentID, ok := paymentIntent.Metadata["payment_id"]
		if !ok || paymentID == "" {
			log.Get().Warnf("Stripe Webhook: no payment_id in metadata for PaymentIntent=%s", paymentIntent.ID)
			c.JSON(http.StatusOK, gin.H{"received": true, "status": "ignored", "reason": "No payment_id in metadata"})
			return
		}

		// Find payment record
		payment := models.GetPaymentByID(paymentID)
		if payment == nil {
			log.Get().Warnf("Stripe Webhook: payment not found for payment_id=%s", paymentID)
			c.JSON(http.StatusOK, gin.H{"received": true, "status": "ignored", "reason": "Payment not found"})
			return
		}

		// Validate channel
		if payment.GetChannelCode() != protocol.PaymentChannelStripe {
			log.Get().Warnf("Stripe Webhook: channel mismatch for payment_id=%s, expected=stripe, got=%s",
				paymentID, payment.GetChannelCode())
			c.JSON(http.StatusOK, gin.H{"received": true, "status": "ignored", "reason": "Channel mismatch"})
			return
		}

		// Process via service
		result := stripeService.ResolvePaymentIntentEvent(&paymentIntent)

		// Update payment record
		values := &models.PaymentValues{}
		values.SetStatus(result.Status).
			SetChannelStatus(result.ChannelStatus).
			SetResCode(result.ResCode).
			SetResMsg(result.ResMsg)

		if result.ChannelPaymentID != "" {
			values.SetChannelPaymentID(result.ChannelPaymentID)
		}

		if result.Status == protocol.StatusSuccess || result.Status == protocol.StatusFailed {
			values.SetCompletedAt(utils.TimeNowMilli())
		}

		if err := models.UpdatePaymentValues(models.DB, payment, values); err != nil {
			log.Get().Errorf("Stripe Webhook: failed to update payment for payment_id=%s, error=%v", paymentID, err)
			// Still return 200 to acknowledge receipt
		}

		// Trigger order status check
		go func() {
			services.GetOrderService().CheckOrderPayment(payment.GetOrderID(), payment.PaymentID)
		}()

		log.Get().Infof("Stripe Webhook processed: payment_id=%s, event=%s, status=%s",
			paymentID, event.Type, result.Status)

	default:
		log.Get().Infof("Stripe Webhook: unhandled event type %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
