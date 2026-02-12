package middleware

import (
	"net/http"
	"strings"

	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
)

// Paths that are always allowed, even during maintenance.
// These include public endpoints apps need for startup, auth, active-ride
// completion, payment webhooks, and the admin panel itself.
var maintenanceExemptPrefixes = []string{
	"/health",
	"/login",
	"/admin/login",
	"/register",
	"/support/config",
	"/system/config",
	"/webhook/",
	"/swagger/",
	// Active ride completion endpoints - do not interrupt rides in progress
	"/order/arrived",
	"/order/start",
	"/order/finish",
	"/order/rating",
	"/order/contact",
	"/order/eta",
	"/order/cash/received",
	"/order/payment",
	"/payment/",
	"/location/update",
	// Allow drivers to go offline during maintenance
	"/offline",
	"/logout",
	"/profile",
}

// MaintenanceMiddleware blocks non-exempt API requests when maintenance mode is enabled.
func MaintenanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Check if this path is exempt from maintenance blocking
		for _, prefix := range maintenanceExemptPrefixes {
			if strings.HasPrefix(path, prefix) || strings.HasSuffix(path, prefix) {
				c.Next()
				return
			}
		}

		// Check maintenance status (Redis-cached, fast)
		if services.GetSystemConfigService().IsMaintenanceMode() {
			config := services.GetSystemConfigService().GetConfig()
			c.JSON(http.StatusServiceUnavailable, &protocol.Result{
				Code: protocol.MaintenanceMode.GetCode(),
				Msg:  config.MaintenanceMessage,
				Data: gin.H{
					"maintenance":       true,
					"message":           config.MaintenanceMessage,
					"support_phone":     config.MaintenancePhone,
					"maintenance_started_at": config.MaintenanceStartAt,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
