package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// GetVehicles 获取车辆列表
// @Summary 获取车辆列表
// @Description 管理员获取车辆列表，支持分页和过滤
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleListRequest true "查询请求"
// @Success 200 {object} protocol.Result{data=protocol.PageResult} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /admin/vehicles [post]
func (t *Admin) GetVehicles(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleListRequest
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

	// 获取车辆列表
	vehicles, total, errorCode := services.GetAdminVehicleService().GetVehicleList(&req)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	// 转换为protocol.Vehicle
	responseData := make([]*protocol.Vehicle, len(vehicles))
	for i, vehicle := range vehicles {
		// Ensure driver details are populated (vehicle.Protocol only sets DriverID)
		responseData[i] = services.GetVehicleService().GetVehicleInfo(vehicle)
	}

	// 返回结果
	result := protocol.NewPageResult(responseData, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// GetVehicleDetail 获取车辆详情
// @Summary 获取车辆详情
// @Description 管理员获取单个车辆的详细信息
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleDetailRequest true "查询请求"
// @Success 200 {object} protocol.Result{data=protocol.Vehicle} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /admin/vehicle/detail [post]
func (t *Admin) GetVehicleDetail(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	vehicle := services.GetAdminVehicleService().GetVehicleByID(req.VehicleID)
	if vehicle == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.VehicleNotFound, lang))
		return
	}

	// Ensure driver details are populated (vehicle.Protocol only sets DriverID)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(services.GetVehicleService().GetVehicleInfo(vehicle)))
}

// UpdateVehicle 更新车辆信息
// @Summary 更新车辆信息
// @Description 管理员更新车辆信息
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleUpdateRequest true "更新请求"
// @Success 200 {object} protocol.Result "更新成功"
// @Failure 200 {object} protocol.Result "更新失败"
// @Security BearerAuth
// @Router /admin/vehicle/update [post]
func (t *Admin) UpdateVehicle(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	errorCode := services.GetAdminVehicleService().UpdateVehicle(&req)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// UpdateVehicleStatus 更新车辆状态
// @Summary 更新车辆状态
// @Description 管理员更新车辆状态（激活、验证等）
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleStatusUpdateRequest true "状态更新请求"
// @Success 200 {object} protocol.Result "更新成功"
// @Failure 200 {object} protocol.Result "更新失败"
// @Security BearerAuth
// @Router /admin/vehicle/status [post]
func (t *Admin) UpdateVehicleStatus(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 验证状态值
	if req.Status != nil {
		validStatuses := []string{protocol.StatusActive, protocol.StatusInactive, protocol.StatusMaintenance, protocol.StatusRetired}
		if !slices.Contains(validStatuses, *req.Status) {
			c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.VehicleInvalidStatus, lang))
			return
		}
	}

	errorCode := services.GetAdminVehicleService().UpdateVehicleStatus(req.VehicleID, req.Status, req.VerifyStatus)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// SearchVehicles 搜索车辆
// @Summary 搜索车辆
// @Description 管理员搜索车辆，支持多种条件过滤
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleSearchRequest true "搜索请求"
// @Success 200 {object} protocol.Result{data=protocol.PageResult} "搜索成功"
// @Failure 200 {object} protocol.Result "搜索失败"
// @Security BearerAuth
// @Router /admin/vehicle/search [post]
func (t *Admin) SearchVehicles(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleSearchRequest
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
	vehicles, total, errorCode := services.GetAdminVehicleService().SearchVehicles(&req)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	// 返回结果
	result := protocol.NewPageResult(vehicles, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	result.AddAttach("params", req)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// CreateVehicle 创建车辆
// @Summary 创建车辆
// @Description 管理员创建新车辆
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleCreateRequest true "创建请求"
// @Success 200 {object} protocol.Result{data=protocol.Vehicle} "创建成功"
// @Failure 200 {object} protocol.Result "创建失败"
// @Security BearerAuth
// @Router /admin/vehicle/create [post]
func (t *Admin) CreateVehicle(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 验证必填字段
	if req.Brand == "" || req.Model == "" || req.PlateNumber == "" {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.MissingParams, lang))
		return
	}

	// 创建车辆
	vehicle, errorCode := services.GetAdminVehicleService().CreateVehicle(&req)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(vehicle.Protocol()))
}

// DeleteVehicle 删除车辆
// @Summary 删除车辆
// @Description 管理员删除车辆（硬删除）
// @Tags Admin,管理员-车辆
// @Accept json
// @Produce json
// @Param request body protocol.VehicleDeleteRequest true "删除请求"
// @Success 200 {object} protocol.Result "删除成功"
// @Failure 200 {object} protocol.Result "删除失败"
// @Security BearerAuth
// @Router /admin/vehicle/delete [post]
func (t *Admin) DeleteVehicle(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.VehicleDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	if req.VehicleID == "" {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.MissingParams, lang))
		return
	}

	log.Printf("Admin delete vehicle request - ID: %s, Reason: %s", req.VehicleID, req.Reason)

	// 删除车辆
	errorCode := services.GetAdminVehicleService().DeleteVehicle(req.VehicleID)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}
