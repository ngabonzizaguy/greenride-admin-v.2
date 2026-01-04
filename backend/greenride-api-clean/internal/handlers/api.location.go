package handlers

import (
	"net/http"

	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
)

// 位置管理相关API处理器
// 包含司机位置上报功能

// @Summary 更新司机位置
// @Description 司机上报当前位置信息
// @Tags Api,位置
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body protocol.UpdateLocationRequest true "位置信息"
// @Router /location/update [post]
func (a *Api) UpdateLocation(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求
	var req protocol.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}
	// 获取当前用户
	user := GetUserFromContext(c)
	req.UserID = user.UserID

	// 调用服务更新位置
	errCode := services.GetUserService().UpdateUserLocation(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 获取当前用户位置
// @Description 获取当前用户的位置信息
// @Tags Api,位置
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Router /location/current [get]
func (a *Api) CurrentLocation(c *gin.Context) {
	// 获取当前用户
	user := GetUserFromContext(c)

	// 构造位置响应
	loc := &protocol.Location{
		Latitude:  user.GetLatitude(),
		Longitude: user.GetLongitude(),
		Address:   user.GetAddress(),
		UpdatedAt: user.GetLocationUpdatedAt(),
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(loc))
}

// @Summary 获取附近司机列表
// @Description 乘客获取附近在线司机列表，用于选择特定司机
// @Tags Api,位置
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param latitude query number true "乘客当前纬度"
// @Param longitude query number true "乘客当前经度"
// @Param radius_km query number false "搜索半径（公里），默认5km"
// @Param limit query int false "返回数量限制，默认20"
// @Router /drivers/nearby [get]
func (a *Api) GetNearbyDrivers(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求参数
	var req protocol.GetNearbyDriversRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 验证经纬度范围
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 调用服务获取附近司机
	drivers, err := services.GetUserService().GetNearbyDrivers(req.Latitude, req.Longitude, req.RadiusKm, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.SystemError, lang))
		return
	}

	// 构造响应
	response := &protocol.GetNearbyDriversResponse{
		Drivers: drivers,
		Count:   len(drivers),
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(response))
}
