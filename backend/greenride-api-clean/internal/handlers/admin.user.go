package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户列表
// @Description 管理员获取用户列表，支持分页和过滤
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.SearchRequest true "搜索条件"
// @Success 200 {object} protocol.PageResult
// @Router /admin/users/search [post]
func (t *Admin) GetUserList(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 解析请求体
	var req protocol.SearchRequest
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

	// 根据是否有关键字决定使用哪个方法
	// 统一使用 SearchUsers 方法，传入完整的请求对象
	users, total, errCode := services.GetUserService().SearchUsers(&req)

	if errCode != protocol.Success {
		log.Printf("Error getting users: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	// 转换为响应格式
	userResponses := make([]*protocol.User, 0, len(users))

	// If searching for drivers, fetch their vehicles
	var vehicleMap map[string]*models.Vehicle
	if req.UserType == protocol.UserTypeDriver || req.UserType == "" {
		driverIDs := make([]string, 0)
		for _, user := range users {
			if user.GetUserType() == protocol.UserTypeDriver {
				driverIDs = append(driverIDs, user.UserID)
			}
		}

		if len(driverIDs) > 0 {
			var vehicles []*models.Vehicle
			if err := models.GetDB().Where("driver_id IN ?", driverIDs).Find(&vehicles).Error; err == nil {
				vehicleMap = make(map[string]*models.Vehicle)
				for _, vehicle := range vehicles {
					if vehicle.DriverID != nil {
						vehicleMap[*vehicle.DriverID] = vehicle
					}
				}
			}
		}
	}

	for _, user := range users {
		userResp := user.Protocol()

		// Populate vehicle if driver and vehicle exists
		if user.GetUserType() == protocol.UserTypeDriver && vehicleMap != nil {
			if vehicle, ok := vehicleMap[user.UserID]; ok {
				userResp.Vehicle = vehicle.Protocol()
			}
		}

		userResponses = append(userResponses, userResp)
	}

	// 返回结果
	result := protocol.NewPageResult(userResponses, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// @Summary 获取用户详情
// @Description 管理员获取单个用户的详细信息
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserIDRequest true "用户ID"
// @Router /admin/users/detail [post]
func (t *Admin) GetUserDetail(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 获取用户信息
	user := services.GetUserService().GetUserByID(req.UserID)
	if user == nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 构建用户详情响应
	info := user.Protocol()

	// 如果是司机，添加车辆信息
	if user.GetUserType() == protocol.UserTypeDriver {
		// 获取司机关联的车辆
		vehicle := models.GetVehicleByDriverID(user.UserID)
		if vehicle != nil {
			info.Vehicle = vehicle.Protocol()
		}
	}

	// 返回用户信息
	c.JSON(http.StatusOK, protocol.NewSuccessResult(info))
}

// @Summary 获取用户行程记录
// @Description 管理员获取用户的行程历史记录
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserRidesRequest true "用户行程查询条件"
// @Router /admin/users/rides [post]
func (t *Admin) GetUserRides(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserRidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// 获取用户订单历史（包括行程订单）
	orders, total := services.GetOrderService().GetOrdersByUser(&req)
	// 返回结果
	result := protocol.NewPageResult(orders, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	result.AddAttach("params", req)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// @Summary 获取用户派单记录
// @Description 管理员获取用户的派单历史记录，包括司机接收到的派单和乘客发起的派单
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserDispatchsRequest true "用户派单记录查询条件"
// @Success 200 {object} protocol.PageResult
// @Router /admin/users/dispatches [post]
func (t *Admin) GetUserDispatchs(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserDispatchsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 参数验证和默认值设置
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// 验证日期范围（如果提供）
	if req.StartDate != nil && req.EndDate != nil {
		if *req.StartDate > *req.EndDate {
			c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang, "start_date cannot be later than end_date"))
			return
		}
	}

	// 获取用户历史派单记录
	records, total := services.GetDispatchService().GetDispatchRecordsByUser(&req)

	// 检查服务调用是否出错
	if records == nil {
		log.Printf("Error getting dispatch records for user %s", req.UserID)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.SystemError, lang))
		return
	}

	// 返回结果
	result := protocol.NewPageResult(records, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	result.AddAttach("params", req)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// @Summary 更新用户状态
// @Description 管理员更新用户的状态和活跃状态
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserStatusUpdateRequest true "用户状态更新信息"
// @Router /admin/users/status [post]
func (t *Admin) UpdateUserStatus(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 更新用户状态
	errCode := services.GetUserService().UpdateUserStatus(req.UserID, req.Status, *req.IsActive)
	if errCode != protocol.Success {
		log.Printf("Error updating user status: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 更新用户信息
// @Description 管理员更新用户的基本信息
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserUpdateRequest true "用户信息更新"
// @Router /admin/users/update [post]
func (t *Admin) UpdateUser(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 获取现有用户信息
	user := services.GetUserService().GetUserByID(req.UserID)
	if user == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 构建更新值
	values := &models.UserValues{
		UserType:    req.UserType,
		Email:       req.Email,
		Phone:       req.Phone,
		CountryCode: req.CountryCode,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Avatar:      req.Avatar,
		Gender:      req.Gender,
		Birthday:    req.Birthday,
		Language:    req.Language,
		Timezone:    req.Timezone,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Status:      req.Status,
	}

	// 如果更新了位置信息，设置位置更新时间
	if req.Latitude != nil && req.Longitude != nil {
		now := utils.TimeNowMilli()
		values.LocationUpdatedAt = &now
	}

	// 验证邮箱唯一性（如果更新了邮箱）
	if req.Email != nil && *req.Email != user.GetEmail() {
		if services.GetUserService().IsEmailExists(*req.Email) {
			c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.EmailAlreadyExists, lang))
			return
		}
	}

	// 验证手机号唯一性（如果更新了手机号）
	if req.Phone != nil && *req.Phone != user.GetPhone() {
		if services.GetUserService().IsPhoneExists(*req.Phone) {
			c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.PhoneAlreadyExists, lang))
			return
		}
	}

	// 更新用户信息
	errCode := services.GetUserService().UpdateUser(user, values)
	if errCode != protocol.Success {
		log.Printf("Error updating user: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}
	// 返回更新后的用户信息
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 验证用户
// @Description 管理员验证用户的邮箱或手机号
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserVerifyRequest true "用户验证信息"
// @Router /admin/users/verify [post]
func (t *Admin) VerifyUser(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.UserVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 获取操作者ID
	operatorID := c.GetString("user_id")

	// 验证用户认证
	errCode := services.GetUserService().VerifyUser(req.UserID, req.IsEmailVerified, req.IsPhoneVerified, operatorID)
	if errCode != protocol.Success {
		log.Printf("Error verifying user: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// CreateUser 管理员创建用户
// @Summary 管理员创建用户
// @Description 管理员创建新用户，密码为空，直接发放优惠券
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminCreateUserRequest true "创建用户请求"
// @Success 200 {object} protocol.Result{data=protocol.User}
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /admin/users/create [post]
func (t *Admin) CreateUser(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求体
	var req protocol.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 设置默认用户类型
	if req.UserType == "" {
		req.UserType = protocol.UserTypePassenger
	}

	// 获取管理员信息
	admin := t.GetUserFromContext(c)

	// 创建用户
	user, errCode := services.GetUserService().CreateUserByAdmin(&req, admin.AdminID)
	if errCode != protocol.Success {
		log.Printf("Admin CreateUser failed with error code: %s", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(user.Protocol()))
}

// AdminDeleteUserRequest 管理员删除用户请求
type AdminDeleteUserRequest struct {
	UserID string `json:"user_id" binding:"required"` // 用户ID
	Reason string `json:"reason,omitempty"`           // 删除原因（可选）
}

// DeleteUser 管理员删除用户
// @Summary 管理员删除用户
// @Description 管理员删除用户账户，支持删除所有类型的用户（乘客和司机）
// @Tags Admin,管理员-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body AdminDeleteUserRequest true "删除用户请求"
// @Success 200 {object} protocol.Result "删除成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /admin/users/delete [post]
func (t *Admin) DeleteUser(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 解析请求体
	var req AdminDeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 获取管理员信息（用于日志记录）
	admin := t.GetUserFromContext(c)

	// 调用 UserService 的删除函数（管理员可以删除所有类型用户，所以 isPassengerOnly=false）
	errCode := services.GetUserService().DeleteUserByID(req.UserID, req.Reason, false)
	if errCode != protocol.Success {
		log.Printf("Admin DeleteUser failed - admin_id: %s, target_user_id: %s, error_code: %s",
			admin.AdminID, req.UserID, errCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	log.Printf("Admin DeleteUser success - admin_id: %s, target_user_id: %s, reason: %s",
		admin.AdminID, req.UserID, req.Reason)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 获取附近司机列表
// @Description 管理员获取附近在线司机列表，用于Live Map显示
// @Tags Admin,管理员-司机
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param latitude query number true "中心纬度"
// @Param longitude query number true "中心经度"
// @Param radius_km query number false "搜索半径（公里），默认50km"
// @Param limit query int false "返回数量限制，默认100"
// @Success 200 {object} protocol.Result "成功返回司机列表"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /admin/drivers/nearby [get]
func (t *Admin) GetNearbyDrivers(c *gin.Context) {
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

	// 设置默认值 - Admin需要更大的范围和更多结果
	radiusKm := req.RadiusKm
	if radiusKm <= 0 {
		radiusKm = 50.0 // 默认50km半径（比乘客端大）
	}
	limit := req.Limit
	if limit <= 0 || limit > 200 {
		limit = 100 // 默认100个，最多200个
	}

	etaMode := req.EtaMode
	if etaMode == "" {
		etaMode = "none" // Admin map doesn't need ETA by default
	}
	// 调用服务获取附近司机
	drivers, err := services.GetUserService().GetNearbyDrivers(req.Latitude, req.Longitude, radiusKm, limit, etaMode)
	if err != nil {
		log.Printf("Admin GetNearbyDrivers failed: %v", err)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(protocol.SystemError, lang))
		return
	}

	// Count total driver statuses for fleet tracker
	db := models.GetDB()
	var totalDrivers, onlineCount, busyCount int64
	db.Model(&models.User{}).Where("user_type = ?", protocol.UserTypeDriver).Count(&totalDrivers)
	db.Model(&models.User{}).Where("user_type = ? AND online_status = ?", protocol.UserTypeDriver, "online").Count(&onlineCount)
	db.Model(&models.User{}).Where("user_type = ? AND online_status = ?", protocol.UserTypeDriver, "busy").Count(&busyCount)
	offlineCount := totalDrivers - onlineCount - busyCount
	if offlineCount < 0 {
		offlineCount = 0
	}

	// 构造响应
	response := &protocol.GetNearbyDriversResponse{
		Drivers:      drivers,
		Count:        len(drivers),
		TotalDrivers: totalDrivers,
		OnlineCount:  onlineCount,
		BusyCount:    busyCount,
		OfflineCount: offlineCount,
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(response))
}
