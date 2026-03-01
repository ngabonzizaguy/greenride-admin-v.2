package services

import (
	"fmt"
	"sync"

	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"gorm.io/gorm"
)

// OrderRatingService 订单评价服务
type OrderRatingService struct {
	db           *gorm.DB
	orderService *OrderService
}

var (
	orderRatingServiceInstance *OrderRatingService
	orderRatingServiceOnce     sync.Once
)

// GetOrderRatingService 获取订单评价服务单例
func GetOrderRatingService() *OrderRatingService {
	orderRatingServiceOnce.Do(func() {
		SetupOrderRatingService()
	})
	return orderRatingServiceInstance
}

// SetupOrderRatingService 设置订单评价服务
func SetupOrderRatingService() {
	orderRatingServiceInstance = &OrderRatingService{
		db:           models.GetDB(),
		orderService: GetOrderService(),
	}
}

// NewOrderRatingService 创建订单评价服务实例
func NewOrderRatingService() *OrderRatingService {
	return &OrderRatingService{
		db:           models.GetDB(),
		orderService: GetOrderService(),
	}
}

// CreateRating 创建评价
func (s *OrderRatingService) CreateRating(req *protocol.CreateOrderRatingRequest) (*models.OrderRating, protocol.ErrorCode) {

	if len(req.Comment) > 500 {
		log.Get().Warnf("CreateOrderRating: Comment too long: %d characters", len(req.Comment))
		return nil, protocol.InvalidParams
	}
	rater := models.GetUserByID(req.UserID)
	if rater == nil {
		return nil, protocol.UserNotFound
	}
	// 获取订单信息以确定被评价者
	order := models.GetOrderByID(req.OrderID)
	if order == nil {
		log.Get().Errorf("CreateOrderRating: Order not found: %s", req.OrderID)
		return nil, protocol.OrderNotFound
	}

	if order.GetStatus() != protocol.StatusCompleted {
		return nil, protocol.OrderCannotUpdate
	}

	// 检查评价者权限
	if order.GetUserID() != req.UserID && order.GetProviderID() != req.UserID {
		return nil, protocol.RatingNotAllowed
	}

	// 检查是否已经评价过
	existingRating := s.GetRatingByOrderAndRater(order.OrderID, rater.UserID, rater.GetUserType())
	if existingRating != nil {
		return nil, protocol.RatingAlreadyExists
	}
	// 创建评价
	rating := &models.OrderRating{
		OrderID:     req.OrderID,
		RaterID:     req.UserID,
		RaterType:   rater.GetUserType(),
		Rating:      req.Rating,
		Comment:     &req.Comment,
		Tags:        req.Tags,
		IsAnonymous: req.IsAnonymous,
		CreatedAt:   utils.TimeNowMilli(),
		UpdatedAt:   utils.TimeNowMilli(),
	}

	if order.GetUserID() == rater.UserID {
		rating.RateeID = order.GetProviderID()
		provider := models.GetUserByID(order.GetProviderID())
		if provider != nil {
			rating.RateeType = provider.GetUserType()
		}
	} else if order.GetProviderID() == rater.UserID {
		rating.RateeID = order.GetUserID()
		user := models.GetUserByID(order.GetUserID())
		if user != nil {
			rating.RateeType = user.GetUserType()
		}
	}

	if err := s.db.Create(rating).Error; err != nil {
		return nil, protocol.DatabaseError
	}

	// 更新用户平均评分
	go s.updateUserAverageRating(rating.RateeID, rating.RateeType)

	return rating, protocol.Success
}

// GetRatingByOrderAndRater 根据订单和评价者获取评价
func (s *OrderRatingService) GetRatingByOrderAndRater(orderID, raterID, raterType string) *models.OrderRating {
	var rating models.OrderRating
	err := s.db.Where("order_id = ? AND rater_id = ? AND rater_type = ?", orderID, raterID, raterType).First(&rating).Error
	if err != nil {
		return nil
	}
	return &rating
}

// GetRatingsByOrder 获取订单的所有评价
func (s *OrderRatingService) GetRatingsByOrder(req *protocol.OrderIDRequest) []*protocol.Rating {

	order := models.GetOrderByID(req.OrderID)
	if order == nil {
		return nil
	}
	if order.GetUserID() != req.UserID && order.GetProviderID() != req.UserID {
		return nil
	}
	var ratings models.OrderRatings

	// 基础查询
	query := s.db.Where("order_id = ?", req.OrderID)
	err := query.Order("created_at DESC").Find(&ratings).Error
	if err != nil {
		return nil
	}
	return ratings.Protocol()
}

// GetRatingsByRatee 获取用户收到的评价
func (s *OrderRatingService) GetRatingsByRatee(rateeID, rateeType string, page, limit int) ([]*models.OrderRating, int64, error) {
	var ratings []*models.OrderRating
	var total int64

	query := s.db.Where("ratee_id = ? AND ratee_type = ?", rateeID, rateeType)

	// 计算总数
	if err := query.Model(&models.OrderRating{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&ratings).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings: %w", err)
	}

	return ratings, total, nil
}

// GetRatingsByRater 获取用户给出的评价
func (s *OrderRatingService) GetRatingsByRater(raterID, raterType string, page, limit int) ([]*models.OrderRating, int64, error) {
	var ratings []*models.OrderRating
	var total int64

	query := s.db.Where("rater_id = ? AND rater_type = ?", raterID, raterType)

	// 计算总数
	if err := query.Model(&models.OrderRating{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&ratings).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings: %w", err)
	}

	return ratings, total, nil
}

// GetUserAverageRating 获取用户平均评分
func (s *OrderRatingService) GetUserAverageRating(userID, userType string) (float64, int64, error) {
	var result struct {
		AvgRating float64
		Count     int64
	}

	err := s.db.Model(&models.OrderRating{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("ratee_id = ? AND ratee_type = ?", userID, userType).
		Scan(&result).Error

	if err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return result.AvgRating, result.Count, nil
}

// UpdateRating 更新评价
func (s *OrderRatingService) UpdateRating(req *protocol.UpdateOrderRatingRequest) protocol.ErrorCode {
	// 验证评价者权限
	var existingRating models.OrderRating
	err := s.db.Where("id = ? AND rater_id = ? ", req.RatingID, req.UserID).First(&existingRating).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return protocol.RatingNotFound
		}
		return protocol.DatabaseError
	}

	// 更新评价
	updates := map[string]interface{}{
		"rating":     req.Rating,
		"comment":    req.Comment,
		"updated_at": utils.TimeNowMilli(),
	}

	if err := s.db.Model(&existingRating).Updates(updates).Error; err != nil {
		return protocol.DatabaseError
	}

	// 更新用户平均评分
	go s.updateUserAverageRating(existingRating.RateeID, existingRating.RateeType)

	return protocol.Success
}

// DeleteRating 删除评价（硬删除）
func (s *OrderRatingService) DeleteRating(req *protocol.RatingIDRequest) protocol.ErrorCode {
	// 验证评价者权限
	var rating models.OrderRating
	err := s.db.Where("id = ? AND rater_id = ?", req.RatingID, req.UserID).First(&rating).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return protocol.RatingNotFound
		}
		return protocol.DatabaseError
	}

	// 硬删除
	if err := s.db.Unscoped().Delete(&rating).Error; err != nil {
		return protocol.DatabaseError
	}

	// 更新用户平均评分
	go s.updateUserAverageRating(rating.RateeID, rating.RateeType)

	return protocol.Success
}

// ReplyToRating 回复评价
func (s *OrderRatingService) ReplyToRating(req *protocol.ReplyToRatingRequest) protocol.ErrorCode {
	// 验证回复者权限（只有被评价者或管理员可以回复）
	rating := models.GetOrderRatingsByID(req.RatingID)
	if rating == nil {
		return protocol.RatingNotFound
	}
	admin := models.GetAdminByID(req.UserID)
	// 检查权限：被评价者或管理员
	if rating.RateeID != req.UserID && admin == nil {
		return protocol.PermissionDenied
	}

	// 更新回复
	updates := map[string]any{
		"reply":      req.Reply,
		"replied_at": utils.TimeNowMilli(),
		"updated_at": utils.TimeNowMilli(),
	}

	if err := s.db.Model(&rating).Updates(updates).Error; err != nil {
		return protocol.DatabaseError
	}

	return protocol.Success
}

// GetRatingStatistics 获取评价统计信息
func (s *OrderRatingService) GetRatingStatistics(userID, userType string) (map[string]interface{}, error) {
	// 获取评分分布
	var distribution []struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}

	err := s.db.Model(&models.OrderRating{}).
		Select("FLOOR(rating) as rating, COUNT(*) as count").
		Where("ratee_id = ? AND ratee_type = ?", userID, userType).
		Group("FLOOR(rating)").
		Order("rating").
		Scan(&distribution).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get rating distribution: %w", err)
	}

	// 获取总体统计
	avgRating, totalCount, err := s.GetUserAverageRating(userID, userType)
	if err != nil {
		return nil, err
	}

	// 计算最近30天的评价数量
	thirtyDaysAgo := utils.TimeNowMilli() - (30 * 24 * 60 * 60 * 1000)
	var recentCount int64
	err = s.db.Model(&models.OrderRating{}).
		Where("ratee_id = ? AND ratee_type = ? AND created_at >= ?", userID, userType, thirtyDaysAgo).
		Count(&recentCount).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent ratings count: %w", err)
	}

	statistics := map[string]interface{}{
		"average_rating":      avgRating,
		"total_ratings":       totalCount,
		"recent_ratings":      recentCount,
		"rating_distribution": distribution,
	}

	return statistics, nil
}

// updateUserAverageRating 更新用户平均评分（异步操作）
func (s *OrderRatingService) updateUserAverageRating(userID, userType string) {
	avgRating, count, err := s.GetUserAverageRating(userID, userType)
	if err != nil {
		return // 静默忽略错误，不影响主流程
	}

	// TODO: 更新用户表中的平均评分字段
	// 这里需要根据实际的用户表结构来实现
	_ = avgRating
	_ = count
}
