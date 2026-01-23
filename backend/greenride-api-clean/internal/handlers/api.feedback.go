package handlers

import (
	"fmt"
	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SubmitFeedback 提交用户反馈
// @Summary 提交用户反馈
// @Description 用户提交反馈信息，无需登录，邮箱作为联系方式
// @Tags Feedback
// @Accept json
// @Produce json
// @Param feedback body protocol.FeedbackRequest true "反馈信息"
// @Success 200 {object} protocol.Result "成功响应"
// @Failure 400 {object} protocol.Result "请求错误"
// @Failure 429 {object} protocol.Result "提交频率过高"
// @Router /feedback/submit [post]
func (a *Api) SubmitFeedback(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 验证邮箱格式
	if !strings.Contains(req.Email, "@") {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidEmail, lang))
		return
	}

	// IP限流检查 - 一分钟内只能提交一次
	clientIP := c.ClientIP()
	rateKey := fmt.Sprintf("feedback:ratelimit:%s", clientIP)

	// 检查是否存在缓存，简化错误处理
	if models.Exists(rateKey) {
		c.JSON(http.StatusTooManyRequests, protocol.NewErrorResult(protocol.RateLimitExceeded, lang))
		return
	}

	// 创建反馈记录
	feedbackService := services.GetFeedbackService()
	feedback, err := feedbackService.CreateFeedback(req.Title, req.Content, req.Email)
	if err != nil {
		log.Errorf("保存反馈失败: %v", err)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.InternalError, lang))
		return
	}

	// 设置限流缓存 - 1分钟过期
	_ = models.Set(rateKey, "1", time.Minute)

	// 返回成功响应
	c.JSON(http.StatusOK, protocol.NewSuccessResult(protocol.FeedbackResponse{
		FeedbackID: feedback.FeedbackID,
	}))
}

// GetSupportConfig 获取支持配置（移动端API - 无需认证）
// @Summary 获取支持配置
// @Description 获取当前的客服支持配置信息，供移动端使用
// @Tags Feedback,Support
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.SupportConfigResponse}
// @Failure 500 {object} protocol.Result
// @Router /support/config [get]
func (a *Api) GetSupportConfig(c *gin.Context) {
	config, err := services.GetSupportService().GetConfig()
	if err != nil {
		log.Errorf("Error getting support config: %v", err)
		// Return default config on error
		defaultConfig := &protocol.SupportConfigResponse{
			SupportEmail:       "support@greenride.rw",
			SupportPhone:       "+250 788 000 000",
			SupportHours:       "Mon-Fri 8:00 AM - 6:00 PM",
			EmergencyPhone:     "+250 788 000 001",
			WhatsAppNumber:      "+250 788 000 000",
			ResponseTimeTarget:  24,
			AutoReplyEnabled:    true,
			EscalationEnabled:   true,
			EscalationTimeout:   48,
		}
		c.JSON(http.StatusOK, protocol.NewSuccessResult(defaultConfig))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(config))
}
