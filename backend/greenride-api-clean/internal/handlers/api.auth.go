package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"greenride/internal/config"
	"greenride/internal/i18n"
	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"greenride/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

// =============================================================================
// 认证接口
// =============================================================================

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	UserType    string `json:"user_type,omitempty"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code,omitempty"` // 可选，默认为卢旺达(RW)
	Password    string `json:"password"`
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	InviteCode  string `json:"invite_code,omitempty"`                    // 可选邀请码
	VerifyCode  string `json:"verify_code,omitempty" binding:"required"` // 验证码 (可选，如果提供则优先验证手机，然后邮箱)
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	UserType string `json:"user_type,omitempty"`
}

// Register 用户注册 (乘客/司机)
// @Summary 用户注册
// @Description 用户注册接口，支持乘客和司机注册。country_code为可选参数，默认为卢旺达(RW)。如果提供verify_code，会自动按优先级验证：先验证手机验证码，如果失败再验证邮箱验证码
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} protocol.Result{data=RegisterResponse} "注册成功"
// @Failure 200 {object} protocol.Result "参数错误或用户已存在"
// @Router /register [post]
func (a *Api) Register(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	//默认创建乘客用户
	req.UserType = protocol.UserTypePassenger

	// 验证验证码(如果提供)
	var verifiedMethod string
	isVerified := false
	// 优先验证手机验证码
	if req.Phone != "" {
		if services.GetVerifyCodeService().VerifySMSCode(protocol.VerifyCodeTypeRegister, req.UserType, req.Phone, req.VerifyCode) {
			// 检查手机号和用户类型组合是否已存在
			if services.GetUserService().IsPhoneExistsWithType(req.Phone, req.UserType) {
				c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserTypeExists, lang))
				return
			}
			isVerified = true
			verifiedMethod = "phone"
		}
	}

	// 如果手机验证码验证失败，尝试邮箱验证码
	if !isVerified && req.Email != "" {
		if services.GetVerifyCodeService().VerifyEmailCode(protocol.VerifyCodeTypeRegister, req.UserType, req.Email, req.VerifyCode) {
			// 检查邮箱和用户类型组合是否已存在
			if services.GetUserService().IsEmailExistsWithType(req.Email, req.UserType) {
				c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserTypeExists, lang))
				return
			}
			isVerified = true
			verifiedMethod = "email"
		}
	}

	// 如果都验证失败，返回错误
	if !isVerified {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidVerificationCode, lang))
		return
	}
	var user *models.User
	switch verifiedMethod {
	case "email":
		user = models.GetUserByEmailAndType(fmt.Sprintf("deleted_%s", req.Email), req.UserType)
	case "phone":
		user = models.GetUserByPhoneAndType(fmt.Sprintf("deleted_%s", req.Phone), req.UserType)
	}
	if user != nil {
		// 用户存在且已删除，恢复用户
		user.SetStatus(protocol.StatusActive).
			SetOnlineStatus(protocol.StatusOnline).
			SetDeletedAt(0).
			SetUsername(req.Username).
			SetPhone(strings.TrimPrefix(user.GetPhone(), "deleted_")).
			SetEmail(strings.TrimPrefix(user.GetEmail(), "deleted_"))
		if err := services.GetUserService().UpdateUser(user, user.UserValues); err != protocol.Success {
			c.JSON(http.StatusOK, protocol.NewErrorResult(err, lang))
			return
		}

		c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(RegisterResponse{
			UserID:   user.UserID,
			Email:    user.GetEmail(),
			Phone:    user.GetPhone(),
			UserType: user.GetUserType(),
		}, lang))
		return
	}
	// 创建用户
	user = models.NewUser()
	values := &models.UserValues{}
	values.SetUserType(req.UserType).
		SetUsername(req.Username).
		SetOnlineStatus(protocol.StatusOnline)

	// 如果验证码验证成功，直接标记为已验证
	switch verifiedMethod {
	case "email":
		values.SetEmail(req.Email).
			MarkEmailAsVerified()
	case "phone":
		values.SetPhone(req.Phone).
			MarkPhoneAsVerified()
	}

	// 默认使用卢旺达国家代码
	values.SetCountryCode("RW")
	if req.CountryCode != "" {
		values.SetCountryCode(req.CountryCode)
	}

	if req.InviteCode != "" {
		inviter := models.GetUserByInviteCode(req.InviteCode)
		if inviter != nil {
			values.SetInvitedBy(inviter.UserID)
		}
	}

	// 根据当前环境设置sandbox标记
	// 如果不是生产环境，将sandbox设置为1，表示测试/沙盒用户
	env := config.Get().Env
	if strings.HasPrefix(req.Phone, "+86") || env != config.ProdEnv {
		values.SetSandbox(1)
	} else {
		values.SetSandbox(0)
	}

	user.SetValues(values)

	if req.Password != "" {
		// 加密密码
		hashedPassword, err := utils.HashPassword(req.Password, user.Salt)
		if err != nil {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InternalError, lang))
			return
		}
		user.Password = &hashedPassword
	}

	// 保存用户
	if errCode := services.GetUserService().CreateUser(user); errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	// 注册成功后发送用户注册信号（异步处理，不影响注册流程）
	if user.IsPassenger() {
		if err := services.SendUserRegisteredSignal(user.UserID); err != nil {
			log.Get().WithField("user_id", user.UserID).
				Errorf("Failed to send user registered signal: %v", err)
			// 错误不影响注册流程，只记录日志
		}
	}

	response := RegisterResponse{
		UserID:   user.UserID,
		Email:    req.Email,
		Phone:    req.Phone,
		UserType: req.UserType,
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email      string `json:"email,omitempty"` // login_type=email时必填
	Phone      string `json:"phone,omitempty"` // login_type=phone时必填
	Password   string `json:"password,omitempty"`
	VerifyCode string `json:"verify_code,omitempty"` // login_type=phone时必填
	UserType   string `json:"user_type"`             // 必填：passenger, driver，用于区分同一手机号/邮箱的不同身份

	// FCM推送相关字段（可选）
	FCMToken string `json:"fcm_token,omitempty"` // FCM推送令牌
	DeviceID string `json:"device_id,omitempty"` // 设备ID
	Platform string `json:"platform,omitempty"`  // 平台：ios, android, web
	AppID    string `json:"app_id,omitempty"`    // 应用标识
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token       string         `json:"token"`
	ExpiresAt   int64          `json:"expires_at"`
	RefreshAt   int64          `json:"refresh_at"`
	User        *protocol.User `json:"user"`
	IsVerified  bool           `json:"is_verified"`
	NeedsVerify []string       `json:"needs_verify,omitempty"` // ["email", "phone"]
}

// Login 用户登录 (乘客/司机)
// @Summary 用户登录
// @Description 用户登录接口，支持邮箱和手机号登录
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} protocol.Result{data=LoginResponse} "登录成功"
// @Failure 200 {object} protocol.Result "参数错误或认证失败"
// @Router /login [post]
func (a *Api) Login(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	log.Get().Infof("[Login] %v", utils.ToJsonString(req))

	// 查找用户 - 必须包含用户类型
	var user *models.User
	isRefresh := false
	if req.Email != "" {
		user = services.GetUserService().GetUserByEmailAndType(req.Email, req.UserType)
	} else if req.Phone != "" {
		user = services.GetUserService().GetUserByPhoneAndType(req.Phone, req.UserType)
	} else {
		token := middleware.ValidToken(c, []byte(a.Jwt.Secret))
		// 提取用户信息
		if claims, ok := token.Claims.(*middleware.JWTClaims); ok {
			user = services.GetUserService().GetUserByID(cast.ToString(claims.UserID))
		}
		// 如果是刷新token失败，返回token相关错误
		if user == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidToken, lang))
			return
		}
		isRefresh = true
	}
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserNotFound, lang))
		return
	}

	if !isRefresh {
		// 验证密码
		if req.Password != "" && !services.GetUserService().VerifyPassword(user, req.Password) {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidPassword, lang))
			return
		}
		// 验证验证码
		if req.VerifyCode != "" && !services.GetVerifyCodeService().VerifySMSCode(protocol.VerifyCodeTypeLogin, req.UserType, req.Phone, req.VerifyCode) {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidVerificationCode, lang))
			return
		}
	}

	// 检查用户状态
	if user.GetStatus() != protocol.StatusActive {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AccountDisabled, lang))
		return
	}

	// 更新最后登录时间和在线状态
	if !isRefresh {
		values := &models.UserValues{}
		values.UpdateLastLogin()
		values.SetActiveStatus(protocol.StatusOnline) // 登录时设置为在线状态
		services.GetUserService().UpdateUser(user, values)
	}

	// 注册FCM Token（如果提供了）
	if req.FCMToken != "" && services.GetFirebaseService() != nil {
		services.GetFirebaseService().RegisterFCMToken(
			user.UserID,
			req.FCMToken,
			req.DeviceID,
			req.Platform,
			req.AppID,
		)
	}
	expiresAt := time.Now().Add(a.Jwt.ExpiresIn)
	tokenString, err := a.GenerateAuthToken(user, expiresAt)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InternalError, lang))
		return
	}

	// 检查验证状态
	needsVerify := []string{}
	if !user.GetIsEmailVerified() && user.GetEmail() != "" {
		needsVerify = append(needsVerify, "email")
	}
	if !user.GetIsPhoneVerified() && user.GetPhone() != "" {
		needsVerify = append(needsVerify, "phone")
	}
	userInfo := user.Protocol()
	if userInfo.UserType == protocol.UserTypeDriver {
		vehicle := models.GetVehicleByDriverID(userInfo.UserID)
		if vehicle != nil {
			userInfo.Vehicle = vehicle.Protocol()
		}
	}
	switch user.GetUserType() {
	case protocol.UserTypeDriver: //先刷新下用户订单序列
		go services.GetUserService().RefreshDriverOrderQueue(user.UserID)
		go services.GetUserService().RefreshDriverRuntimeCache(user.UserID)
	case protocol.UserTypePassenger:
		go services.GetUserService().RefreshUserOrderQueue(user.UserID)
	}
	response := LoginResponse{
		Token:       tokenString,
		ExpiresAt:   expiresAt.UnixMilli(),
		RefreshAt:   utils.TimeNowMilli(),
		User:        userInfo,
		IsVerified:  len(needsVerify) == 0,
		NeedsVerify: needsVerify,
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

func (a *Api) GenerateAuthToken(user *models.User, expiresAt time.Time) (string, error) {
	claims := &middleware.JWTClaims{
		UserID:   user.UserID,
		UserType: user.GetUserType(),
		Email:    user.GetEmail(),
		Phone:    user.GetPhone(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.Jwt.Secret)) // 使用配置中的JWT密钥
	return tokenString, err
}

// SendVerifyCodeRequest 发送验证码请求
type SendVerifyCodeRequest struct {
	Type        string `json:"type" binding:"required,oneof=register login reset_password change_phone change_email"` // 验证码类型
	Email       string `json:"email,omitempty"`                                                                       // 邮箱地址，有则发送邮箱验证码
	Phone       string `json:"phone,omitempty"`                                                                       // 手机号码，有则发送短信验证码
	CountryCode string `json:"country_code,omitempty"`                                                                // 手机号国家代码，默认为卢旺达(RW)
	UserType    string `json:"user_type" binding:"required,oneof=passenger driver"`                                   // 用户类型，用于区分同一联系方式的不同身份
}

// SendVerifyCode 发送验证码
// @Summary 发送验证码
// @Description 发送验证码到邮箱或手机。系统会自动根据提供的email和phone字段判断发送方式，优先使用phone。country_code为可选参数，默认为卢旺达(RW)
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body SendVerifyCodeRequest true "发送验证码请求"
// @Success 200 {object} protocol.Result "发送成功"
// @Failure 200 {object} protocol.Result "参数错误或发送失败"
// @Router /send-verify-code [post]
func (a *Api) SendVerifyCode(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req SendVerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}
	if req.Phone == "" && req.Email == "" {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "email or phone is required"))
		return
	}
	if req.Phone != "" && req.Email != "" {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, "email and phone cannot be used together"))
		return
	}

	// 设置默认国家代码
	if req.CountryCode == "" {
		req.CountryCode = "RW" // 默认为卢旺达
	}
	// 根据类型进行业务验证
	switch req.Type {
	case protocol.VerifyCodeTypeRegister:
		// 注册验证码：检查用户和用户类型组合是否已存在
		if req.Email != "" && services.GetUserService().IsEmailExistsWithType(req.Email, req.UserType) {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserTypeExists, lang))
			return
		}
		if req.Phone != "" && services.GetUserService().IsPhoneExistsWithType(req.Phone, req.UserType) {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserTypeExists, lang))
			return
		}
	case protocol.VerifyCodeTypeLogin, protocol.VerifyCodeTypeResetPassword:
		// 登录或重置密码验证码：检查用户是否存在
		var userExists bool
		if req.Email != "" {
			userExists = services.GetUserService().IsEmailExistsWithType(req.Email, req.UserType)
		} else {
			userExists = services.GetUserService().IsPhoneExistsWithType(req.Phone, req.UserType)
		}
		if !userExists {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.UserNotFound, lang))
			return
		}
	}
	var target, contentType string
	// 发送验证码
	if req.Email != "" {
		target = req.Email
		contentType = protocol.MsgChannelEmail
	} else {
		target = req.Phone
		contentType = protocol.MsgChannelSms
	}
	errCode, remainingSeconds := services.GetVerifyCodeService().SendVerifyCode(contentType, target, req.UserType, req.Type, lang)
	if errCode != protocol.Success {
		if errCode == protocol.VerificationCooldown {
			result := protocol.NewErrorResult(protocol.VerificationCooldown, lang, strconv.Itoa(remainingSeconds))
			result.Data = map[string]any{
				"remaining_seconds": strconv.Itoa(remainingSeconds),
			}
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		}
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Type        string `json:"type" binding:"required"`   // 验证码类型
	Method      string `json:"method" binding:"required"` // email, phone
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Code        string `json:"code" binding:"required"`                             // 验证码
	UserType    string `json:"user_type" binding:"required,oneof=passenger driver"` // 用户类型，用于区分同一联系方式的不同身份
	UserID      string `json:"user_id,omitempty"`                                   // 某些情况下需要用户ID
	NewPassword string `json:"new_password,omitempty"`                              // 重置密码时需要
	NewEmail    string `json:"new_email,omitempty"`                                 // 更换邮箱时需要
	NewPhone    string `json:"new_phone,omitempty"`                                 // 更换手机时需要
}

// VerifyCode 验证验证码
// @Summary 验证验证码
// @Description 验证邮箱或手机验证码
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "验证请求"
// @Success 200 {object} protocol.Result "验证成功"
// @Failure 200 {object} protocol.Result "验证失败"
// @Router /verify-code [post]
func (a *Api) VerifyCode(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 验证验证码
	var isValid bool

	if req.Method == "email" {
		isValid = services.GetVerifyCodeService().VerifyEmailCode(req.Type, req.UserType, req.Email, req.Code)
	} else {
		isValid = services.GetVerifyCodeService().VerifySMSCode(req.Type, req.UserType, req.Phone, req.Code)
	}

	if !isValid {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidVerificationCode, lang))
		return
	}

	// 根据类型处理业务逻辑
	response := map[string]interface{}{
		"message": i18n.TranslateMessage("verification_success", lang),
		"type":    req.Type,
		"method":  req.Method,
	}

	switch req.Type {
	case "register":
		// 注册验证码验证成功，标记用户为已验证
		var user *models.User
		if req.Method == "email" {
			user = services.GetUserService().GetUserByEmailAndType(req.Email, req.UserType)
		} else {
			user = services.GetUserService().GetUserByPhoneAndType(req.Phone, req.UserType)
		}

		if user != nil {
			values := &models.UserValues{}
			if req.Method == "email" {
				values.MarkEmailAsVerified()
			} else {
				values.MarkPhoneAsVerified()
			}
			services.GetUserService().UpdateUser(user, values)
			response["user_verified"] = true
		}

	case "reset_password":
		// 重置密码
		if req.NewPassword == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "new_password is required"))
			return
		}

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

		response["password_reset"] = true

	case "change_email":
		// 更换邮箱 (需要登录状态)
		user := GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
			return
		}

		if req.NewEmail == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "new_email is required"))
			return
		}

		values := &models.UserValues{}
		values.SetEmail(req.NewEmail).MarkEmailAsVerified()
		if errCode := services.GetUserService().UpdateUser(user, values); errCode != protocol.Success {
			c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
			return
		}

		response["email_changed"] = true

	case "change_phone":
		// 更换手机号 (需要登录状态)
		user := GetUserFromContext(c)
		if user == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.AuthenticationFailed, lang))
			return
		}

		if req.NewPhone == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.MissingParams, lang, "new_phone is required"))
			return
		}

		values := &models.UserValues{}
		values.SetPhone(req.NewPhone).MarkPhoneAsVerified()
		if errCode := services.GetUserService().UpdateUser(user, values); errCode != protocol.Success {
			c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
			return
		}

		response["phone_changed"] = true
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}

// LogoutRequest 登出请求（可选参数）
type LogoutRequest struct {
	FcmToken string `json:"fcm_token,omitempty"` // FCM Token（可选，如果提供则只停用该设备的token）
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，清除服务端相关信息。如果提供fcm_token，只停用该设备的FCM token；否则停用所有设备的token
// @Tags Api,认证
// @Accept json
// @Produce json
// @Param request body LogoutRequest false "登出请求（可选参数）"
// @Success 200 {object} protocol.Result "登出成功"
// @Security BearerAuth
// @Router /logout [post]
func (a *Api) Logout(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	user := GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusOK, protocol.NewSuccessResult("")) // 未登录状态直接返回成功
		return
	}
	// 尝试解析请求体获取fcm_token（可选）
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, err.Error()))
		return
	}

	// 清除用户的FCM Token
	if services.GetFirebaseService() != nil {
		// 如果提供了fcm_token，只停用该设备的token
		err := services.GetFirebaseService().DeactivateToken(user.UserID, req.FcmToken)
		if err != nil {
			log.Get().Warnf("Failed to deactivate FCM token for user %s token %s: %v", user.UserID, req.FcmToken, err)
		} else {
			log.Get().Infof("FCM token deactivated successfully for user %s token %s", user.UserID, req.FcmToken)
		}
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}
