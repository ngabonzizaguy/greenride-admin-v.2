package handlers

import (
	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 用户信息接口
// =============================================================================

// Profile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags Api,用户
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.User} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /profile [get]
func (a *Api) Profile(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户（已在中间件中验证过）
	user := GetUserFromContext(c)
	userInfo := user.Protocol()
	if userInfo.UserType == protocol.UserTypeDriver {
		vehicle := models.GetVehicleByDriverID(userInfo.UserID)
		if vehicle != nil {
			userInfo.Vehicle = vehicle.Protocol()
		}
	}
	// 返回用户信息
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(userInfo, lang))
}

// UserOnline 司机上线
// @Summary 司机上线
// @Description 司机上线，只有司机类型用户可以调用此接口。上线后司机状态变为在线，车辆状态变为可用
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.UserOnlineRequest true "上线请求"
// @Security BearerAuth
// @Router /online [post]
func (a *Api) UserOnline(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)

	var req protocol.UserOnlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	// 检查用户类型
	if !user.IsDriver() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	if err := services.GetUserService().UserOnline(req); err != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(err, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// UserOffline 司机下线
// @Summary 司机下线
// @Description 司机下线，只有司机类型用户可以调用此接口。下线后司机状态变为离线，车辆状态变为不可用
// @Tags Api,司机
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result "下线成功"
// @Security BearerAuth
// @Router /offline [post]
func (a *Api) UserOffline(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
		return
	}

	// 检查用户类型
	if !user.IsDriver() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	if err := services.GetUserService().UserOffline(user.UserID); err != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(err, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Username    *string `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Birthday    *int64  `json:"birthday,omitempty"`
	Language    *string `json:"language,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	State       *string `json:"state,omitempty"`
	Country     *string `json:"country,omitempty"`
	PostalCode  *string `json:"postal_code,omitempty"`
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Description 按需更新用户个人信息，只更新提供的字段
// @Tags Api,用户
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "更新请求"
// @Success 200 {object} protocol.Result{data=protocol.User} "更新成功"
// @Failure 200 {object} protocol.Result "更新失败"
// @Security BearerAuth
// @Router /profile/update [post]
func (a *Api) UpdateProfile(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 按需更新字段
	values := &models.UserValues{}

	if req.FirstName != nil {
		values.SetFirstName(*req.FirstName)
	}
	if req.LastName != nil {
		values.SetLastName(*req.LastName)
	}
	if req.Username != nil {
		values.SetUsername(*req.Username)
	}
	if req.DisplayName != nil {
		values.DisplayName = req.DisplayName
	}
	if req.Avatar != nil {
		values.Avatar = req.Avatar
	}
	if req.Gender != nil {
		values.Gender = req.Gender
	}
	if req.Birthday != nil {
		values.Birthday = req.Birthday
	}
	if req.Language != nil {
		values.Language = req.Language
	}
	if req.Timezone != nil {
		values.Timezone = req.Timezone
	}
	if req.Address != nil {
		values.Address = req.Address
	}
	if req.City != nil {
		values.City = req.City
	}
	if req.State != nil {
		values.State = req.State
	}
	if req.Country != nil {
		values.Country = req.Country
	}
	if req.PostalCode != nil {
		values.PostalCode = req.PostalCode
	}

	// 更新用户信息
	if errCode := services.GetUserService().UpdateUser(user, values); errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	// 重新获取更新后的用户信息
	updatedUser := services.GetUserService().GetUserByID(user.UserID)
	if updatedUser == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(updatedUser.Protocol(), lang))
}

// UpdateAvatarResponse 更新头像响应
type UpdateAvatarResponse struct {
	AvatarURL string `json:"avatar_url"` // 新的头像URL
}

// UpdateAvatar 更新用户头像
// @Summary 更新用户头像
// @Description 上传并更新用户头像到AWS S3
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "头像文件 (jpg, jpeg, png, gif, webp, 最大5MB)"
// @Success 200 {object} protocol.Result{data=UpdateAvatarResponse} "上传成功"
// @Failure 400 {object} protocol.Result "请求参数错误 (2002:缺少参数, 2007:不支持的文件类型, 2005:文件过大)"
// @Failure 500 {object} protocol.Result "服务器错误 (1004:服务不可用, 1007:文件操作错误, 1009:第三方服务错误)"
// @Security ApiKeyAuth
// @Router /profile/update/avatar [put]
func (a *Api) UpdateAvatar(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	log.Infof("UpdateAvatar started - user_id: %s, user_type: %s", user.UserID, user.GetUserType())

	// 获取上传的头像文件
	avatarFile, err := c.FormFile("avatar")
	if err != nil {
		log.Errorf("UpdateAvatar failed to get avatar file - user_id: %s, error: %s", user.UserID, err.Error())
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang))
		return
	}

	log.Infof("UpdateAvatar file received - user_id: %s, filename: %s, size: %d", user.UserID, avatarFile.Filename, avatarFile.Size)

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(avatarFile.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		log.Warnf("UpdateAvatar invalid file extension - user_id: %s, filename: %s, extension: %s", user.UserID, avatarFile.Filename, ext)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidContentType, lang))
		return
	}

	log.Infof("UpdateAvatar file extension validated - user_id: %s, extension: %s", user.UserID, ext)

	// 检查文件大小（限制为5MB）
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if avatarFile.Size > maxSize {
		log.Warnf("UpdateAvatar file size too large - user_id: %s, size: %d, max_size: %d", user.UserID, avatarFile.Size, maxSize)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.RequestTooLarge, lang))
		return
	}

	log.Infof("UpdateAvatar file size validated - user_id: %s, size: %d", user.UserID, avatarFile.Size)

	// 打开上传的文件
	file, err := avatarFile.Open()
	if err != nil {
		log.Errorf("UpdateAvatar failed to open uploaded file - user_id: %s, error: %s", user.UserID, err.Error())
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.FileError, lang))
		return
	}
	defer file.Close()

	log.Infof("UpdateAvatar file opened successfully - user_id: %s", user.UserID)

	// 使用 UserService 统一的头像上传方法
	log.Infof("UpdateAvatar starting avatar upload - user_id: %s, filename: %s", user.UserID, avatarFile.Filename)
	uploadReq := &services.UploadAvatarRequest{
		UserID:    user.UserID,
		Reader:    file,
		Extension: ext,
		Filename:  avatarFile.Filename,
		Size:      avatarFile.Size,
	}

	uploadResp := services.GetUserService().UploadAvatar(c.Request.Context(), uploadReq)
	if uploadResp.ErrorCode != protocol.Success {
		log.Errorf("UpdateAvatar failed to upload avatar - user_id: %s, error_code: %s", user.UserID, uploadResp.ErrorCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(uploadResp.ErrorCode, lang))
		return
	}

	log.Infof("UpdateAvatar completed successfully - user_id: %s, avatar_url: %s", user.UserID, uploadResp.AvatarURL)

	// 返回新的头像URL
	response := UpdateAvatarResponse{
		AvatarURL: uploadResp.AvatarURL,
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

// DeleteAccountRequest 删除账户请求
type DeleteAccountRequest struct {
	Reason string `json:"reason,omitempty"` // 删除原因（可选）
}

// DeleteAccount 删除账户
// @Summary 删除账户
// @Description 硬删除用户账户，彻底移除用户主数据
// @Tags Api,用户
// @Accept json
// @Produce json
// @Param request body DeleteAccountRequest true "删除请求"
// @Success 200 {object} protocol.Result "删除成功"
// @Failure 200 {object} protocol.Result "删除失败"
// @Security BearerAuth
// @Router /account/delete [post]
func (a *Api) DeleteAccount(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)

	// 检查用户类型，只有乘客能删除账户
	if !user.IsPassenger() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}

	var req DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 调用 UserService 的删除函数
	errCode := services.GetUserService().DeleteUserByID(user.UserID, req.Reason, true)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// =============================================================================
// 车辆信息接口
// =============================================================================

// UserVehicleRequest 用户车辆查询请求
type UserVehicleRequest struct {
	// 无需参数，直接根据当前用户ID查询车辆信息
}

// GetUserVehicle 获取用户车辆信息
// @Summary 获取用户车辆信息
// @Description 获取当前用户的车辆信息，司机返回绑定车辆，车主返回拥有的车辆
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body UserVehicleRequest true "查询请求"
// @Success 200 {object} protocol.Result{data=protocol.Vehicle} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /vehicle [post]
func (a *Api) GetUserVehicle(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
		return
	}
	if !user.IsDriver() {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	// 司机查询绑定的车辆
	vehicle := services.GetVehicleService().GetVehicleByDriverID(user.UserID)
	if vehicle == nil {
		c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
		return
	}
	// 返回车辆信息
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(vehicle.Protocol(), lang))
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Method      string `json:"method" binding:"required,oneof=email phone"` // email, phone
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Code        string `json:"code" binding:"required"`                             // 验证码
	NewPassword string `json:"new_password" binding:"required"`                     // 新密码
	UserType    string `json:"user_type" binding:"required,oneof=passenger driver"` // 用户类型
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过验证码重置密码
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} protocol.Result "重置成功"
// @Failure 200 {object} protocol.Result "重置失败"
// @Router /reset-password [post]
func (a *Api) ResetPassword(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 验证验证码
	var isValid bool
	if req.Method == "email" {
		if req.Email == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "email is required"))
			return
		}
		isValid = services.GetVerifyCodeService().VerifyEmailCode(protocol.VerifyCodeTypeResetPassword, req.UserType, req.Email, req.Code)
	} else {
		if req.Phone == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "phone is required"))
			return
		}
		isValid = services.GetVerifyCodeService().VerifySMSCode(protocol.VerifyCodeTypeResetPassword, req.UserType, req.Phone, req.Code)
	}

	if !isValid {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidVerificationCode, lang))
		return
	}

	// 查找用户
	var user *models.User
	if req.Method == "email" {
		user = services.GetUserService().GetUserByEmailAndType(req.Email, req.UserType)
	} else {
		user = services.GetUserService().GetUserByPhoneAndType(req.Phone, req.UserType)
	}

	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 更新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword, user.Salt)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InternalError, lang))
		return
	}

	values := &models.UserValues{}
	values.Password = &hashedPassword
	if errCode := services.GetUserService().UpdateUser(user, values); errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// ChangePassword 修改密码 (需要登录)
// @Summary 修改密码
// @Description 用户修改密码，需要提供旧密码
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} protocol.Result "修改成功"
// @Failure 200 {object} protocol.Result "修改失败"
// @Security BearerAuth
// @Router /change-password [post]
func (a *Api) ChangePassword(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 获取当前用户
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 验证旧密码
	if !services.GetUserService().VerifyPassword(user, req.OldPassword) {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidPassword, lang))
		return
	}

	// 更新新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword, user.Salt)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InternalError, lang))
		return
	}

	values := &models.UserValues{}
	values.Password = &hashedPassword
	if errCode := services.GetUserService().UpdateUser(user, values); errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}
