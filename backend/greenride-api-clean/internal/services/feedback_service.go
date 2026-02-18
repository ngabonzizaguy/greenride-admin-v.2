package services

import (
	"database/sql"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
)

// FeedbackService 反馈服务
type FeedbackService struct{}

// FeedbackWithUser 带用户信息的反馈
type FeedbackWithUser struct {
	*models.Feedback
	UserFullName string `json:"user_full_name"`
	UserEmail    string `json:"user_email"`
	UserPhone    string `json:"user_phone"`
}

// CreateFeedback 创建反馈记录
func (s *FeedbackService) CreateFeedback(title, content, email, name, phone, userID, category, feedbackType, severity string) (*models.Feedback, error) {
	feedback := models.NewFeedback()
	feedback.SetContent(title, content, "")
	feedback.SetContact(name, phone, email)
	if userID != "" {
		feedback.SetUserID(userID)
	}
	if feedbackType != "" {
		feedback.SetFeedbackType(feedbackType)
	} else {
		feedback.SetFeedbackType(protocol.FeedbackTypeSuggestion)
	}
	if category != "" {
		feedback.SetCategory(category)
	} else {
		feedback.SetCategory(protocol.FeedbackCategoryOther)
	}
	if severity != "" {
		feedback.SetSeverity(severity)
	}

	db := models.GetDB()
	if db == nil {
		return nil, sql.ErrConnDone
	}

	err := db.Create(feedback).Error
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// SearchFeedback 搜索反馈列表
func (s *FeedbackService) SearchFeedback(req *protocol.FeedbackSearchRequest) ([]*FeedbackWithUser, int64, protocol.ErrorCode) {
	db := models.GetDB()
	if db == nil {
		return nil, 0, protocol.SystemError
	}

	var feedbacks []*models.Feedback
	var total int64

	query := db.Model(&models.Feedback{})

	// Apply filters
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ? OR contact_email LIKE ?", keyword, keyword, keyword)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.FeedbackType != "" {
		query = query.Where("feedback_type = ?", req.FeedbackType)
	}

	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	if req.Severity != "" {
		query = query.Where("severity = ?", req.Severity)
	}

	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}

	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.StartDate != nil {
		query = query.Where("created_at >= ?", *req.StartDate)
	}

	if req.EndDate != nil {
		query = query.Where("created_at <= ?", *req.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	// Apply pagination and sorting
	offset := (req.Page - 1) * req.Limit
	query = query.Order("created_at DESC").Offset(offset).Limit(req.Limit)

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, 0, protocol.SystemError
	}

	// Build response with user info
	result := make([]*FeedbackWithUser, 0, len(feedbacks))
	for _, fb := range feedbacks {
		item := &FeedbackWithUser{Feedback: fb}

		// Try to get user info if user_id exists
		if fb.GetUserID() != "" {
			user := GetUserService().GetUserByID(fb.GetUserID())
			if user != nil {
				item.UserFullName = user.GetFullName()
				item.UserEmail = user.GetEmail()
				item.UserPhone = user.GetPhone()
			}
		}

		result = append(result, item)
	}

	return result, total, protocol.Success
}

// GetFeedbackByID 根据ID获取反馈
func (s *FeedbackService) GetFeedbackByID(feedbackID string) (*FeedbackWithUser, protocol.ErrorCode) {
	db := models.GetDB()
	if db == nil {
		return nil, protocol.SystemError
	}

	var feedback models.Feedback
	if err := db.Where("feedback_id = ?", feedbackID).First(&feedback).Error; err != nil {
		return nil, protocol.InvalidParams // feedback not found
	}

	result := &FeedbackWithUser{Feedback: &feedback}

	// Get user info if available
	if feedback.GetUserID() != "" {
		user := GetUserService().GetUserByID(feedback.GetUserID())
		if user != nil {
			result.UserFullName = user.GetFullName()
			result.UserEmail = user.GetEmail()
			result.UserPhone = user.GetPhone()
		}
	}

	return result, protocol.Success
}

// UpdateFeedback 更新反馈
func (s *FeedbackService) UpdateFeedback(req *protocol.FeedbackUpdateRequest, adminID string) protocol.ErrorCode {
	db := models.GetDB()
	if db == nil {
		return protocol.SystemError
	}

	var feedback models.Feedback
	if err := db.Where("feedback_id = ?", req.FeedbackID).First(&feedback).Error; err != nil {
		return protocol.InvalidParams // feedback not found
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Status != nil {
		updates["status"] = *req.Status

		// Set timestamps based on status
		now := utils.TimeNowMilli()
		switch *req.Status {
		case protocol.StatusReviewing:
			updates["handled_by"] = adminID
			updates["handled_at"] = now
		case protocol.StatusResolved:
			updates["resolved_at"] = now
		case protocol.StatusClosed:
			updates["closed_at"] = now
		}
	}

	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}

	if req.Severity != nil {
		updates["severity"] = *req.Severity
	}

	if req.Resolution != nil {
		updates["resolution"] = *req.Resolution
	}

	if req.ResolutionNote != nil {
		updates["resolution_note"] = *req.ResolutionNote
	}

	if req.AssignedTo != nil {
		updates["assigned_to"] = *req.AssignedTo
		updates["assigned_at"] = utils.TimeNowMilli()
	}

	if req.InternalNotes != nil {
		updates["internal_notes"] = *req.InternalNotes
	}

	if req.PublicResponse != nil {
		updates["public_response"] = *req.PublicResponse
	}

	// Increment update count
	currentCount := 0
	if feedback.UpdateCount != nil {
		currentCount = *feedback.UpdateCount
	}
	updates["update_count"] = currentCount + 1

	if err := db.Model(&feedback).Updates(updates).Error; err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// IncrementViewCount 增加查看次数
func (s *FeedbackService) IncrementViewCount(feedbackID string) error {
	db := models.GetDB()
	if db == nil {
		return sql.ErrConnDone
	}

	return db.Model(&models.Feedback{}).
		Where("feedback_id = ?", feedbackID).
		UpdateColumn("view_count", db.Raw("COALESCE(view_count, 0) + 1")).Error
}

// GetFeedbackStats 获取反馈统计
func (s *FeedbackService) GetFeedbackStats() (*protocol.FeedbackStats, error) {
	db := models.GetDB()
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stats := &protocol.FeedbackStats{}

	// Total count
	if err := db.Model(&models.Feedback{}).Count(&stats.TotalFeedback).Error; err != nil {
		return nil, err
	}

	// Status counts
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := db.Model(&models.Feedback{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case protocol.StatusPending:
			stats.PendingCount = sc.Count
		case protocol.StatusInProgress, protocol.StatusReviewing:
			stats.InProgressCount += sc.Count
		case protocol.StatusResolved:
			stats.ResolvedCount = sc.Count
		case protocol.StatusClosed:
			stats.ClosedCount = sc.Count
		}
	}

	// Type counts
	var typeCounts []struct {
		FeedbackType string
		Count        int64
	}
	if err := db.Model(&models.Feedback{}).
		Select("feedback_type, count(*) as count").
		Group("feedback_type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		switch tc.FeedbackType {
		case protocol.FeedbackTypeComplaint:
			stats.ComplaintCount = tc.Count
		case protocol.FeedbackTypeSuggestion:
			stats.SuggestionCount = tc.Count
		case protocol.FeedbackTypeCompliment:
			stats.ComplimentCount = tc.Count
		case protocol.FeedbackTypeBugReport:
			stats.BugReportCount = tc.Count
		}
	}

	// Average response time (in hours)
	var avgResponse struct {
		AvgTime float64
	}
	if err := db.Model(&models.Feedback{}).
		Select("AVG(response_time) as avg_time").
		Where("response_time > 0").
		Scan(&avgResponse).Error; err == nil {
		stats.AvgResponseTime = avgResponse.AvgTime
	}

	// Average resolution time (in hours)
	var avgResolution struct {
		AvgTime float64
	}
	if err := db.Model(&models.Feedback{}).
		Select("AVG(resolution_time) as avg_time").
		Where("resolution_time > 0").
		Scan(&avgResolution).Error; err == nil {
		stats.AvgResolutionTime = avgResolution.AvgTime
	}

	return stats, nil
}

// DeleteFeedback 删除单个反馈
func (s *FeedbackService) DeleteFeedback(feedbackID string) protocol.ErrorCode {
	db := models.GetDB()
	if db == nil {
		return protocol.SystemError
	}

	result := db.Where("feedback_id = ?", feedbackID).Delete(&models.Feedback{})
	if result.Error != nil {
		return protocol.SystemError
	}

	if result.RowsAffected == 0 {
		return protocol.InvalidParams // feedback not found
	}

	return protocol.Success
}

// BulkDeleteFeedback 批量删除反馈
func (s *FeedbackService) BulkDeleteFeedback(feedbackIDs []string) (int64, protocol.ErrorCode) {
	if len(feedbackIDs) == 0 {
		return 0, protocol.InvalidParams
	}

	db := models.GetDB()
	if db == nil {
		return 0, protocol.SystemError
	}

	result := db.Where("feedback_id IN ?", feedbackIDs).Delete(&models.Feedback{})
	if result.Error != nil {
		return 0, protocol.SystemError
	}

	return result.RowsAffected, protocol.Success
}

// GetFeedbackService 获取反馈服务实例
func GetFeedbackService() *FeedbackService {
	return &FeedbackService{}
}
