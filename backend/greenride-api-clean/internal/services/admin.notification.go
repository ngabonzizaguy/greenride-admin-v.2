package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"log"

	"gorm.io/gorm"
)

// AdminNotificationService 管理员通知服务
type AdminNotificationService struct {
	db *gorm.DB
}

// GetAdminNotificationService 获取管理员通知服务实例
func GetAdminNotificationService() *AdminNotificationService {
	return &AdminNotificationService{
		db: models.GetDB(),
	}
}

// SendNotification 发送通知（广播或指定用户）
func (s *AdminNotificationService) SendNotification(req *protocol.AdminSendNotificationRequest) protocol.ErrorCode {
	// 确定目标用户
	var userIDs []string
	var userType string

	switch req.Audience {
	case "all":
		// 获取所有活跃用户和司机
		var users []models.User
		if err := s.db.Where("status = ? AND deleted_at IS NULL", protocol.StatusActive).Find(&users).Error; err != nil {
			log.Printf("Failed to get users for broadcast: %v", err)
			return protocol.DatabaseError
		}
		for _, u := range users {
			userIDs = append(userIDs, u.UserID)
		}
		userType = "" // 所有类型
	case "drivers":
		// 获取所有活跃司机
		var drivers []models.User
		if err := s.db.Where("user_type = ? AND status = ? AND deleted_at IS NULL", protocol.UserTypeDriver, protocol.StatusActive).Find(&drivers).Error; err != nil {
			log.Printf("Failed to get drivers for broadcast: %v", err)
			return protocol.DatabaseError
		}
		for _, d := range drivers {
			userIDs = append(userIDs, d.UserID)
		}
		userType = protocol.UserTypeDriver
	case "users":
		// 获取所有活跃乘客
		var passengers []models.User
		if err := s.db.Where("user_type = ? AND status = ? AND deleted_at IS NULL", protocol.UserTypePassenger, protocol.StatusActive).Find(&passengers).Error; err != nil {
			log.Printf("Failed to get passengers for broadcast: %v", err)
			return protocol.DatabaseError
		}
		for _, p := range passengers {
			userIDs = append(userIDs, p.UserID)
		}
		userType = protocol.UserTypePassenger
	default:
		return protocol.InvalidParams
	}

	// 创建通知记录
	now := utils.TimeNowMilli()

	// 批量创建通知
	for _, userID := range userIDs {
		notification := models.NewNotificationV2()

		// Use setter methods for fields that have them, direct assignment for others
		notification.SetUserID(userID)
		if userType != "" {
			notification.UserType = utils.StringPtr(userType)
		}
		notification.SetType(req.Type)
		notification.Category = utils.StringPtr(req.Category)
		notification.SetTitle(req.Title)
		notification.SetContent(req.Content)
		if req.Summary != "" {
			notification.Summary = utils.StringPtr(req.Summary)
		}
		notification.SetStatus(models.NotificationStatusPending)
		notification.Priority = utils.StringPtr(models.NotificationPriorityNormal)
		// IsRead is already set to false in NewNotificationV2()
		notification.SetChannels([]string{"push"})
		notification.PushTitle = utils.StringPtr(req.Title)
		notification.PushBody = utils.StringPtr(req.Content)

		// 如果指定了计划发送时间
		if req.ScheduledAt != nil && *req.ScheduledAt > now {
			notification.SetScheduledTime(*req.ScheduledAt)
		} else {
			notification.SetScheduledTime(now)
		}

		// 保存通知
		if err := s.db.Create(notification).Error; err != nil {
			log.Printf("Failed to create notification for user %s: %v", userID, err)
			continue // 继续处理其他用户
		}

		// TODO: 发送推送通知（可以异步处理）
		// 这里可以调用 FCM 服务发送推送
	}

	return protocol.Success
}

// SearchNotifications 搜索通知（管理员查看通知历史）
func (s *AdminNotificationService) SearchNotifications(req *protocol.AdminNotificationSearchRequest) ([]*models.Notification, int64, protocol.ErrorCode) {
	var notifications []*models.Notification
	var total int64

	query := s.db.Model(&models.Notification{})

	// 用户类型过滤
	if req.UserType != "" {
		query = query.Where("user_type = ?", req.UserType)
	}

	// 通知类型过滤
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 关键字搜索
	if req.Keyword != "" {
		searchTerm := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", searchTerm, searchTerm)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.DatabaseError
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&notifications).Error; err != nil {
		return nil, 0, protocol.DatabaseError
	}

	return notifications, total, protocol.Success
}

// GetUnreadCount 获取未读通知数量（管理员）
func (s *AdminNotificationService) GetUnreadCount(adminID string) (int64, protocol.ErrorCode) {
	var count int64

	// 管理员通知：user_type = 'admin' 或 user_id = adminID
	if err := s.db.Model(&models.Notification{}).
		Where("(user_type = ? OR user_id = ?) AND is_read = ?", "admin", adminID, false).
		Count(&count).Error; err != nil {
		return 0, protocol.DatabaseError
	}

	return count, protocol.Success
}

// MarkAsRead 标记通知为已读
func (s *AdminNotificationService) MarkAsRead(notificationID string) protocol.ErrorCode {
	now := utils.TimeNowMilli()
	isRead := true

	if err := s.db.Model(&models.Notification{}).
		Where("notification_id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": isRead,
			"read_at": now,
			"status":  models.NotificationStatusRead,
		}).Error; err != nil {
		return protocol.DatabaseError
	}

	return protocol.Success
}

// MarkAllAsRead 标记所有通知为已读（管理员）
func (s *AdminNotificationService) MarkAllAsRead(adminID string) protocol.ErrorCode {
	now := utils.TimeNowMilli()
	isRead := true

	if err := s.db.Model(&models.Notification{}).
		Where("(user_type = ? OR user_id = ?) AND is_read = ?", "admin", adminID, false).
		Updates(map[string]interface{}{
			"is_read": isRead,
			"read_at": now,
			"status":  models.NotificationStatusRead,
		}).Error; err != nil {
		return protocol.DatabaseError
	}

	return protocol.Success
}
