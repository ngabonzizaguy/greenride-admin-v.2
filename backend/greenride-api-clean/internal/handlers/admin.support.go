package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Admin Support Configuration Endpoints
// ============================================================================

// GetSupportConfig 获取支持配置
// @Summary 获取支持配置
// @Description 获取当前的客服支持配置信息
// @Tags Admin,Support
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=protocol.SupportConfigResponse}
// @Router /support/config [get]
func (t *Admin) GetSupportConfig(c *gin.Context) {
	config, err := services.GetSupportService().GetConfig()
	if err != nil {
		log.Printf("Error getting support config: %v", err)
		// Return default config on error
		defaultConfig := &protocol.SupportConfigResponse{
			SupportEmail:       "support@greenride.rw",
			SupportPhone:       "+250 788 000 000",
			SupportHours:       "Mon-Fri 8:00 AM - 6:00 PM",
			EmergencyPhone:     "+250 788 000 001",
			WhatsAppNumber:     "+250 788 000 000",
			ResponseTimeTarget: 24,
			AutoReplyEnabled:   true,
			EscalationEnabled:  true,
			EscalationTimeout:  48,
		}
		c.JSON(http.StatusOK, protocol.NewSuccessResult(defaultConfig))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(config))
}

// UpdateSupportConfig 更新支持配置
// @Summary 更新支持配置
// @Description 更新客服支持配置信息
// @Tags Admin,Support
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.SupportConfigUpdateRequest true "配置信息"
// @Success 200 {object} protocol.Result
// @Router /support/config [post]
func (t *Admin) UpdateSupportConfig(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.SupportConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// Get admin info
	admin := t.GetUserFromContext(c)
	if admin == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
		return
	}

	// Update config
	err := services.GetSupportService().UpdateConfig(&req, admin.AdminID)
	if err != nil {
		log.Printf("Error updating support config: %v", err)
		c.JSON(http.StatusInternalServerError, protocol.NewBusinessErrorResult("Failed to update support configuration"))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

