package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 发送通知
// @Description 管理员发送广播通知给用户/司机/所有人
// @Tags Admin,管理员-通知
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminSendNotificationRequest true "发送通知请求"
// @Router /admin/notifications/send [post]
func (t *Admin) SendNotification(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.AdminSendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	errCode := services.GetAdminNotificationService().SendNotification(&req)
	if errCode != protocol.Success {
		log.Printf("Error sending notification: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 搜索通知
// @Description 管理员搜索通知历史记录
// @Tags Admin,管理员-通知
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.AdminNotificationSearchRequest true "搜索请求"
// @Router /admin/notifications/search [post]
func (t *Admin) SearchNotifications(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.AdminNotificationSearchRequest
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

	notifications, total, errCode := services.GetAdminNotificationService().SearchNotifications(&req)
	if errCode != protocol.Success {
		log.Printf("Error searching notifications: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	// 转换为响应格式
	notificationResponses := make([]map[string]interface{}, 0, len(notifications))
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, notification.ToMap())
	}

	result := protocol.NewPageResult(notificationResponses, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// @Summary 获取未读通知数量
// @Description 获取管理员的未读通知数量
// @Tags Admin,管理员-通知
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /admin/notifications/unread-count [get]
func (t *Admin) GetUnreadCount(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	admin := t.GetUserFromContext(c)

	count, errCode := services.GetAdminNotificationService().GetUnreadCount(admin.AdminID)
	if errCode != protocol.Success {
		log.Printf("Error getting unread count: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(map[string]interface{}{
		"count": count,
	}))
}

// @Summary 标记通知为已读
// @Description 标记单个通知为已读
// @Tags Admin,管理员-通知
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.NotificationIDRequest true "通知ID"
// @Router /admin/notifications/mark-read [post]
func (t *Admin) MarkAsRead(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	var req protocol.NotificationIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	errCode := services.GetAdminNotificationService().MarkAsRead(req.NotificationID)
	if errCode != protocol.Success {
		log.Printf("Error marking notification as read: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// @Summary 标记所有通知为已读
// @Description 标记管理员的所有通知为已读
// @Tags Admin,管理员-通知
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /admin/notifications/mark-all-read [post]
func (t *Admin) MarkAllAsRead(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	admin := t.GetUserFromContext(c)

	errCode := services.GetAdminNotificationService().MarkAllAsRead(admin.AdminID)
	if errCode != protocol.Success {
		log.Printf("Error marking all notifications as read: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}
