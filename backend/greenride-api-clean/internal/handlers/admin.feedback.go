package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ============================================================================
// Admin Feedback Management Endpoints
// ============================================================================

// SearchFeedback 搜索反馈列表
// @Summary 搜索反馈列表
// @Description 管理员搜索反馈/投诉列表，支持分页和过滤
// @Tags Admin,Feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.FeedbackSearchRequest true "搜索条件"
// @Success 200 {object} protocol.PageResult
// @Router /feedback/search [post]
func (t *Admin) SearchFeedback(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.FeedbackSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	// Search feedback
	feedbacks, total, errCode := services.GetFeedbackService().SearchFeedback(&req)
	if errCode != protocol.Success {
		log.Printf("Error searching feedback: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	// Convert to response format
	responses := make([]*protocol.FeedbackListItem, 0, len(feedbacks))
	for _, fb := range feedbacks {
		responses = append(responses, toFeedbackListItem(fb))
	}

	// Return paginated result
	result := protocol.NewPageResult(responses, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}

// GetFeedbackDetail 获取反馈详情
// @Summary 获取反馈详情
// @Description 管理员获取单个反馈的详细信息
// @Tags Admin,Feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.FeedbackIDRequest true "反馈ID"
// @Success 200 {object} protocol.Result{data=protocol.FeedbackDetail}
// @Router /feedback/detail [post]
func (t *Admin) GetFeedbackDetail(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.FeedbackIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// Get feedback
	feedback, errCode := services.GetFeedbackService().GetFeedbackByID(req.FeedbackID)
	if errCode != protocol.Success {
		log.Printf("Error getting feedback detail: %v", errCode)
		c.JSON(http.StatusNotFound, protocol.NewErrorResult(errCode, lang))
		return
	}

	// Increment view count
	_ = services.GetFeedbackService().IncrementViewCount(req.FeedbackID)

	// Get user info if available
	var userInfo *protocol.FeedbackUserInfo
	if feedback.GetUserID() != "" {
		user := services.GetUserService().GetUserByID(feedback.GetUserID())
		if user != nil {
			userInfo = &protocol.FeedbackUserInfo{
				UserID:   user.UserID,
				FullName: user.GetFullName(),
				Email:    user.GetEmail(),
				Phone:    user.GetPhone(),
				Avatar:   user.GetAvatar(),
			}
		}
	}

	// Convert to detailed response
	detail := toFeedbackDetail(feedback, userInfo)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(detail))
}

// UpdateFeedback 更新反馈
// @Summary 更新反馈
// @Description 管理员更新反馈状态、分配、解决方案等
// @Tags Admin,Feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.FeedbackUpdateRequest true "更新信息"
// @Success 200 {object} protocol.Result
// @Router /feedback/update [post]
func (t *Admin) UpdateFeedback(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.FeedbackUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	// Get admin info
	admin := t.GetUserFromContext(c)
	if admin == nil {
		c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
		return
	}

	// Update feedback
	errCode := services.GetFeedbackService().UpdateFeedback(&req, admin.AdminID)
	if errCode != protocol.Success {
		log.Printf("Error updating feedback: %v", errCode)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// GetFeedbackStats 获取反馈统计
// @Summary 获取反馈统计
// @Description 获取反馈统计数据用于仪表盘展示
// @Tags Admin,Feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=protocol.FeedbackStats}
// @Router /feedback/stats [get]
func (t *Admin) GetFeedbackStats(c *gin.Context) {
	stats, err := services.GetFeedbackService().GetFeedbackStats()
	if err != nil {
		log.Printf("Error getting feedback stats: %v", err)
		// Return default stats on error
		defaultStats := &protocol.FeedbackStats{
			TotalFeedback:    0,
			PendingCount:     0,
			InProgressCount:  0,
			ResolvedCount:    0,
			ComplaintCount:   0,
			SuggestionCount:  0,
			AvgResponseTime:  0,
			AvgResolutionTime: 0,
		}
		c.JSON(http.StatusOK, protocol.NewSuccessResult(defaultStats))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(stats))
}

// ============================================================================
// Helper Functions
// ============================================================================

func toFeedbackListItem(fb *services.FeedbackWithUser) *protocol.FeedbackListItem {
	item := &protocol.FeedbackListItem{
		FeedbackID:   fb.FeedbackID,
		Title:        fb.GetTitle(),
		FeedbackType: fb.GetFeedbackType(),
		Category:     fb.GetCategory(),
		Status:       fb.GetStatus(),
		Severity:     fb.GetSeverity(),
		Priority:     fb.GetPriority(),
		Rating:       fb.GetRating(),
		CreatedAt:    fb.CreatedAt,
		UpdatedAt:    fb.UpdatedAt,
	}

	// Add user info if available
	if fb.UserFullName != "" {
		item.UserName = fb.UserFullName
		item.UserEmail = fb.UserEmail
	} else if fb.ContactEmail != nil && *fb.ContactEmail != "" {
		item.UserEmail = *fb.ContactEmail
		if fb.ContactName != nil {
			item.UserName = *fb.ContactName
		}
	}

	return item
}

func toFeedbackDetail(fb *services.FeedbackWithUser, userInfo *protocol.FeedbackUserInfo) *protocol.FeedbackDetail {
	detail := &protocol.FeedbackDetail{
		FeedbackID:   fb.FeedbackID,
		Title:        fb.GetTitle(),
		Content:      fb.GetContent(),
		Description:  getStringPtr(fb.Description),
		FeedbackType: fb.GetFeedbackType(),
		Category:     fb.GetCategory(),
		Subcategory:  getStringPtr(fb.Subcategory),
		Status:       fb.GetStatus(),
		Severity:     fb.GetSeverity(),
		Priority:     fb.GetPriority(),
		Rating:       fb.GetRating(),
		Resolution:   fb.GetResolution(),
		CreatedAt:    fb.CreatedAt,
		UpdatedAt:    fb.UpdatedAt,
		User:         userInfo,
	}

	// Add contact info
	if fb.ContactName != nil {
		detail.ContactName = *fb.ContactName
	}
	if fb.ContactEmail != nil {
		detail.ContactEmail = *fb.ContactEmail
	}
	if fb.ContactPhone != nil {
		detail.ContactPhone = *fb.ContactPhone
	}

	// Add timestamps
	if fb.ResolvedAt != nil {
		detail.ResolvedAt = fb.ResolvedAt
	}
	if fb.AssignedTo != nil {
		detail.AssignedTo = *fb.AssignedTo
	}
	if fb.HandledBy != nil {
		detail.HandledBy = *fb.HandledBy
	}

	// Add order info
	if fb.OrderID != nil {
		detail.OrderID = *fb.OrderID
	}

	// Add location
	if fb.Latitude != nil && fb.Longitude != nil {
		detail.Latitude = fb.Latitude
		detail.Longitude = fb.Longitude
	}
	if fb.Location != nil {
		detail.Location = *fb.Location
	}

	// Add device info
	if fb.DeviceType != nil {
		detail.DeviceType = *fb.DeviceType
	}
	if fb.AppVersion != nil {
		detail.AppVersion = *fb.AppVersion
	}

	// Add attachments
	detail.Attachments = fb.GetAttachments()

	// Add internal notes
	if fb.InternalNotes != nil {
		detail.InternalNotes = *fb.InternalNotes
	}
	if fb.PublicResponse != nil {
		detail.PublicResponse = *fb.PublicResponse
	}

	return detail
}

func getStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

