package handlers

import (
	"net/http"

	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
)

// GetSystemConfig 获取系统配置（公共接口，用于移动端启动检查）
// @Summary 获取系统配置
// @Description 获取系统全局配置，包括维护模式状态。此端点无需认证。
// @Tags System
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.SystemConfigResponse}
// @Router /system/config [get]
func (a *Api) GetSystemConfig(c *gin.Context) {
	config := services.GetSystemConfigService().GetConfig()
	c.JSON(http.StatusOK, protocol.NewSuccessResult(config))
}

// AdminGetSystemConfig 管理员获取系统配置
// @Summary 获取系统配置（管理员）
// @Tags Admin,System
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=protocol.SystemConfigResponse}
// @Router /system/config [get]
func (a *Admin) AdminGetSystemConfig(c *gin.Context) {
	config := services.GetSystemConfigService().GetConfig()
	c.JSON(http.StatusOK, protocol.NewSuccessResult(config))
}

// AdminUpdateSystemConfig 管理员更新系统配置
// @Summary 更新系统配置（管理员）
// @Description 更新系统全局配置，例如启用/禁用维护模式
// @Tags Admin,System
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.SystemConfigUpdateRequest true "系统配置更新"
// @Success 200 {object} protocol.Result{data=protocol.SystemConfigResponse}
// @Router /system/config [post]
func (a *Admin) AdminUpdateSystemConfig(c *gin.Context) {
	admin := a.GetUserFromContext(c)
	if admin == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
		return
	}

	var req protocol.SystemConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, "", err.Error()))
		return
	}

	if err := services.GetSystemConfigService().UpdateConfig(&req, admin.AdminID); err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult("Failed to update system config"))
		return
	}

	// Return updated config
	config := services.GetSystemConfigService().GetConfig()
	c.JSON(http.StatusOK, protocol.NewSuccessResult(config))
}

type HardDeleteCleanupRequest struct {
	Confirm string `json:"confirm"` // must be PURGE_LEGACY_DELETED
	DryRun  bool   `json:"dry_run"`
}

// AdminPurgeLegacyDeleted performs hard-delete cleanup for legacy soft-deleted rows.
// Admin auth is required; current deployment uses admin role only (no super_admin split).
func (a *Admin) AdminPurgeLegacyDeleted(c *gin.Context) {
	admin := a.GetUserFromContext(c)
	if admin == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
		return
	}

	var req HardDeleteCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, "", err.Error()))
		return
	}
	if req.Confirm != "PURGE_LEGACY_DELETED" {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, "", "invalid confirmation token"))
		return
	}

	summary, errCode := services.RunHardDeleteCleanupWithOptions(req.DryRun)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, ""))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(summary))
}
