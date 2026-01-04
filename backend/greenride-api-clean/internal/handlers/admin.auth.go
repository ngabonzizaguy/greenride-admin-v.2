package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login 管理员登录
// @Summary 管理员登录
// @Description 管理员用户名密码登录
// @Tags Admin,管理员-认证
// @Accept json
// @Produce json
// @Param request body protocol.AdminLoginRequest true "登录信息"
// @Success 200 {object} protocol.Result{data=models.Admin}
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /login [post]
func (t *Admin) Login(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 查找管理员
	user := services.GetAdminAdminService().GetAdminByUsername(req.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}
	// 验证密码
	if !services.GetAdminAdminService().VerifyPassword(user, req.Password) {
		// 记录失败登录
		services.GetAdminAdminService().RecordFailedLogin(user)
		c.JSON(http.StatusUnauthorized, protocol.NewErrorResult(protocol.InvalidCredentials, lang))
		return
	}

	// 检查登录权限
	clientIP := c.ClientIP()
	errorCode := services.GetAdminAdminService().CheckLoginPermission(user, clientIP)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	tokenString, err := t.GenerateAuthToken(user, expiresAt)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InternalError, lang))
		return
	}

	// 记录登录成功
	sessionID := utils.GenerateUUID()
	if errorCode := services.GetAdminAdminService().RecordLogin(user, clientIP, sessionID); errorCode != protocol.Success {
		log.Printf("Error recording login: %s", errorCode)
	}

	// 返回结果
	c.JSON(http.StatusOK, protocol.NewSuccessResult(map[string]interface{}{
		"token": tokenString,
		"user":  user.Protocol(),
	}))
}

func (t *Admin) GenerateAuthToken(user *models.Admin, expiresAt time.Time) (string, error) {
	claims := &middleware.JWTClaims{
		UserID: user.AdminID,
		Email:  user.GetEmail(),
		Role:   user.GetRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.Jwt.Secret)) // 使用配置中的JWT密钥
	return tokenString, err
}

// ChangePassword 管理员修改密码
// @Summary 管理员修改密码
// @Description 管理员修改自己的密码
// @Tags Admin,管理员-认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminChangePasswordRequest true "密码修改信息"
// @Success 200 {object} protocol.Result
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /change-password [post]
func (t *Admin) ChangePassword(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.AdminChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 从JWT中获取管理员ID
	adminID := c.GetString("user_id")

	// 查找管理员
	admin := services.GetAdminAdminService().GetAdminByID(adminID)
	if admin == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 验证旧密码
	if !services.GetAdminAdminService().VerifyPassword(admin, req.OldPassword) {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidPassword, lang, "Current password is incorrect"))
		return
	}

	// 验证新密码强度
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.WeakPassword, lang, "Password must be at least 8 characters long"))
		return
	}

	// 更新密码
	errorCode := services.GetAdminAdminService().UpdatePassword(adminID, req.NewPassword)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// ResetPassword 重置管理员密码
// @Summary 重置管理员密码
// @Description 重置指定管理员的密码
// @Tags Admin,管理员-管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminResetPasswordRequest true "密码重置信息"
// @Success 200 {object} protocol.Result
// @Failure 400 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /reset-password [post]
func (t *Admin) ResetPassword(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.AdminResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// 从JWT中获取操作者管理员ID
	operatorAdminID := c.GetString("user_id")

	// 查找操作者，验证权限
	operator := services.GetAdminAdminService().GetAdminByID(operatorAdminID)
	if operator == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewErrorResult(protocol.InvalidCredentials, lang))
		return
	}

	// 检查操作者权限（这里简单检查，实际应该有更复杂的权限系统）
	operatorRole := operator.GetRole()
	if operatorRole != "super_admin" {
		c.JSON(http.StatusForbidden, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}

	// 验证新密码强度
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.WeakPassword, lang))
		return
	}

	// 重置密码
	errorCode := services.GetAdminAdminService().ResetPassword(
		req.TargetAdminID,
		req.NewPassword,
		operatorAdminID,
	)
	if errorCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errorCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// Logout 管理员登出
// @Summary 管理员登出
// @Description 管理员退出登录
// @Tags Admin,管理员-认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /logout [post]
func (t *Admin) Logout(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 从JWT中获取管理员ID
	adminID := c.GetString("user_id")

	// 查找管理员
	admin := services.GetAdminAdminService().GetAdminByID(adminID)
	if admin == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 执行登出
	if errorCode := services.GetAdminAdminService().Logout(admin); errorCode != protocol.Success {
		log.Printf("Error during logout: %s", errorCode)
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// Info 获取管理员信息
// @Summary 获取当前管理员信息
// @Description 获取当前登录管理员的详细信息
// @Tags Admin,管理员-认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=models.Admin}
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /info [get]
func (t *Admin) Info(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 从JWT中获取管理员ID
	userID := c.GetString("user_id")

	// 查找管理员
	user := services.GetAdminAdminService().GetAdminByID(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	// 获取管理员信息
	c.JSON(http.StatusOK, protocol.NewSuccessResult(user.Protocol()))
}
