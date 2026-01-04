package middleware

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"greenride/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// AuthorizationHeader 认证请求头
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer前缀
	BearerPrefix = "Bearer "
	// TokenKey 请求参数中的token键名
	TokenKey = "token"
)

// JWTClaims JWT载荷
type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"` // passenger, driver
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

// 辅助函数

func ValidToken(c *gin.Context, jwt_secrets []byte) *jwt.Token {
	// 从Header获取Token
	tokenString := GetTokenFromRequest(c)
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt_secrets, nil
	})

	if err != nil || !token.Valid {
		return nil
	}
	return token
}

// GenerateUserToken 生成支付页面访问Token
// 该Token包含payment_id和user_id，用于支付页面的身份验证和数据获取
func GenerateUserToken(userID string, expiresAt time.Time, jwtSecret string) string {
	// 获取JWT配置
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ""
	}

	return tokenString
}

// 从请求获取Token
func GetTokenFromRequest(c *gin.Context) string {
	// 优先从请求头获取token
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader != "" {
		// 检查 Authorization 格式是否为 Bearer <token>
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString != authHeader {
			return tokenString
		}
	}

	// 从请求参数获取token
	return c.Query(TokenKey)
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User type not found",
			})
			c.Abort()
			return
		}

		// 检查角色权限
		userTypeStr := userType.(string)
		hasPermission := slices.Contains(allowedRoles, userTypeStr)

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 简单的速率限制中间件（生产环境建议使用Redis实现）
func RateLimitMiddleware() gin.HandlerFunc {
	// 这里可以集成第三方限流库，如golang.org/x/time/rate
	return func(c *gin.Context) {
		// 简单示例，实际应该根据IP或用户进行限制
		c.Next()
	}
}

// RequireUserType 用户类型验证中间件
func RequireUserType(userTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User type not found",
			})
			c.Abort()
			return
		}

		userTypeStr := userType.(string)
		for _, allowedType := range userTypes {
			if userTypeStr == allowedType {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Access denied for user type",
		})
		c.Abort()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 计算请求处理时间
		latency := time.Since(start)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 记录请求完成信息
		log.Get().Infof("[REQUEST_END] %s %s | %d | %v | %s",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency,
			c.ClientIP())
	}
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理在处理请求过程中可能出现的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		}
	}
}
