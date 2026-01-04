package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrders 获取车辆列表
// @Summary 获取车辆列表
// @Description 获取车辆列表，支持分页和类型过滤
// @Tags Api,车辆
// @Accept json
// @Produce json
// @Param request body protocol.GetVehiclesRequest true "获取车辆列表请求"
// @Success 200 {object} protocol.PageResult "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /vehicles [post]
func (a *Api) GetVehicles(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 获取当前用户
	user := GetUserFromContext(c)
	var req protocol.GetVehiclesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	if !user.IsDriver() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	// 获取车辆列表
	vehicles, total := services.GetVehicleService().GetVehicles(&req)
	result := protocol.NewPageResult(vehicles, total, &protocol.Pagination{
		Page: 1,
		Size: 2000,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}
