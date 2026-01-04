package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// =============================================================================
// 工具方法
// =============================================================================

// getUserClaim 从上下文获取用户信息
func GetUserClaim(c *gin.Context) *middleware.JWTClaims {
	// 首先尝试从 user_claim 键获取（如果其他地方设置了）
	value, exists := c.Get("user_claim")
	if exists {
		if claim, ok := value.(*middleware.JWTClaims); ok {
			return claim
		}
	}

	// 如果没有，从单独的键构建 JWTClaims
	userID, userIDExists := c.Get("user_id")
	if !userIDExists {
		return nil
	}

	email, _ := c.Get("email")
	userType, _ := c.Get("user_type")

	// 构建并返回 JWTClaims
	return &middleware.JWTClaims{
		UserID:   cast.ToString(userID),
		Email:    cast.ToString(email),
		UserType: cast.ToString(userType),
	}
}

// GetUserFromContext 从上下文获取完整的用户对象
func GetUserFromContext(c *gin.Context) *models.User {
	value, exists := c.Get("user")
	if !exists {
		return nil
	}
	if user, ok := value.(*models.User); ok {
		return user
	}
	return nil
}

// GetAdminFromContext 从上下文获取完整的管理员对象
func GetAdminFromContext(c *gin.Context) *models.Admin {
	value, exists := c.Get("user")
	if !exists {
		return nil
	}
	if user, ok := value.(*models.Admin); ok {
		return user
	}
	return nil
}
