package services

import (
	"errors"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"
	"slices"
	"time"

	"gorm.io/gorm"
)

// ============================================================================
// 优惠码管理相关方法
// ============================================================================

// SearchPromotions 搜索优惠码（统一的搜索接口）
func (s *AdminService) SearchPromotions(req *protocol.PromotionSearchRequest) (result []*protocol.Promotion, total int64) {
	var promotions models.Promotions

	// 构建查询
	query := s.db.Model(&models.Promotion{})

	// 关键词搜索（优惠码、标题、描述）
	if req.Keyword != "" {
		query = query.Where("code LIKE ? OR title LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 优惠码类型筛选
	if req.Type != "" {
		query = query.Where("promo_type = ?", req.Type)
	}

	// 状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.Limit
	err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&promotions).Error
	if err != nil {
		return
	}
	result = promotions.Protocol()
	return
}

// GetPromoCodeList 获取优惠码列表（管理端）- 保持向后兼容
func (s *AdminService) GetPromoCodeList(keyword, promoType, status, approvalStatus string, page, limit int) ([]*models.Promotion, int64, error) {
	var promoCodes []*models.Promotion
	var total int64

	// 构建查询
	query := s.db.Model(&models.Promotion{})

	// 关键词搜索（优惠码、标题、描述）
	if keyword != "" {
		query = query.Where("code LIKE ? OR title LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 优惠码类型筛选
	if promoType != "" {
		query = query.Where("promo_type = ?", promoType)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 审批状态筛选
	if approvalStatus != "" {
		query = query.Where("approval_status = ?", approvalStatus)
	}

	// 统计总数
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&promoCodes).Error
	if err != nil {
		return nil, 0, err
	}

	return promoCodes, total, nil
}

// CreatePromotion 创建优惠码
func (s *AdminService) CreatePromotion(req *protocol.CreatePromotionRequest) (*models.Promotion, protocol.ErrorCode) {
	// 业务验证
	if errCode := s.validatePromotionRequest(req); errCode != protocol.Success {
		return nil, errCode
	}

	// 检查优惠码是否已存在
	var existingPromo models.Promotion
	if err := s.db.Where("code = ?", req.Code).First(&existingPromo).Error; err == nil {
		return nil, protocol.PromotionAlreadyExists
	}

	// 创建优惠码
	promotion := &models.Promotion{
		PromotionID:     utils.GeneratePromotionID(),
		PromotionValues: &models.PromotionValues{},
	}

	// 设置基本信息
	promotion.SetCode(req.Code).
		SetTitle(req.Title).
		SetDiscountType(req.DiscountType).
		SetDescription(req.Description).
		SetCreatedBy(req.UserID).
		SetStatus(protocol.StatusActive)

	// 设置可选信息
	if req.MaxDiscountAmount != nil {
		promotion.MaxDiscountAmount = req.MaxDiscountAmount
	}
	if req.MinOrderAmount != nil {
		promotion.MinOrderAmount = req.MinOrderAmount
	}
	if req.UsageLimit != nil {
		promotion.UsageLimit = req.UsageLimit
	}
	if req.UserUsageLimit != nil {
		promotion.UserUsageLimit = req.UserUsageLimit
	}
	// 保存到数据库
	if err := s.db.Create(promotion).Error; err != nil {
		return nil, protocol.PromotionCreationFailed
	}

	return promotion, protocol.Success
}

// GetPromotionDetail 根据ID获取优惠码详情
func (s *AdminService) GetPromotionDetail(promotionID string) *protocol.Promotion {
	var promotion models.Promotion
	if err := s.db.Where("promotion_id = ?", promotionID).First(&promotion).Error; err != nil {
		return nil
	}

	return promotion.Protocol()
}

// UpdatePromotion 更新优惠码
func (s *AdminService) UpdatePromotion(req *protocol.UpdatePromotionRequest) protocol.ErrorCode {
	// 业务验证
	if errCode := s.validatePromotionUpdateRequest(req); errCode != protocol.Success {
		return errCode
	}

	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", req.PromotionID).First(&promoCode).Error; err != nil {
		return protocol.PromotionNotFound
	}

	// 更新字段使用 PromotionValues 的 setter 方法
	if req.Title != nil {
		promoCode.SetTitle(*req.Title)
	}
	if req.Description != nil {
		promoCode.SetDescription(*req.Description)
	}
	if req.DiscountValue != nil {
		promoCode.SetDiscountValue(*req.DiscountValue)
	}
	if req.MaxDiscountAmount != nil {
		promoCode.SetMaxDiscountAmount(*req.MaxDiscountAmount)
	}
	if req.MinOrderAmount != nil {
		promoCode.SetMinOrderAmount(*req.MinOrderAmount)
	}
	if req.UsageLimit != nil {
		promoCode.SetUsageLimit(*req.UsageLimit)
	}
	if req.UserUsageLimit != nil {
		promoCode.SetUserUsageLimit(*req.UserUsageLimit)
	}
	if req.StartDate != nil {
		promoCode.SetStartDate(*req.StartDate)
	}
	if req.EndDate != nil {
		promoCode.SetEndDate(*req.EndDate)
	}
	if req.ValidCities != nil {
		promoCode.SetValidCities(*req.ValidCities)
	}
	if req.ValidVehicleTypes != nil {
		promoCode.SetValidVehicleTypes(*req.ValidVehicleTypes)
	}
	if req.Priority != nil {
		promoCode.SetPriority(*req.Priority)
	}
	if req.Tags != nil {
		promoCode.SetTags(*req.Tags)
	}

	// 更新数据
	if err := s.db.Save(&promoCode).Error; err != nil {
		return protocol.PromotionUpdateFailed
	}

	return protocol.Success
}

// UpdatePromoCodeStatus 更新优惠码状态
func (s *AdminService) UpdatePromoCodeStatus(req *protocol.UpdatePromotionStatusRequest) protocol.ErrorCode {
	// 验证状态值
	validStatuses := []string{protocol.StatusActive, protocol.StatusInactive, protocol.StatusExpired, protocol.StatusSuspended, protocol.StatusDeleted}
	validStatus := slices.Contains(validStatuses, req.Status)
	if !validStatus {
		return protocol.PromotionStatusInvalid
	}

	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", req.PromotionID).First(&promoCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.PromotionNotFound
		}
		return protocol.DatabaseError
	}

	// 更新状态和相关字段
	promoCode.SetStatus(req.Status).SetUpdatedBy(req.UserID)

	if req.Reason != "" {
		promoCode.SetStatusReason(req.Reason)
	}

	// 根据状态设置相应的时间戳
	currentTime := time.Now().Unix()
	switch req.Status {
	case protocol.StatusActive:
		promoCode.SetActivatedAt(currentTime)
	case protocol.StatusSuspended:
		promoCode.SetSuspendedAt(currentTime)
	case protocol.StatusDeleted:
		promoCode.SetDeletedAt(currentTime)
	}

	// 更新数据
	if err := s.db.Save(&promoCode).Error; err != nil {
		return protocol.PromotionUpdateFailed
	}

	return protocol.Success
}

// ApprovePromotion 审批通过优惠码
func (s *AdminService) ApprovePromotion(req *protocol.ApprovePromotionRequest) protocol.ErrorCode {
	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", req.PromotionID).First(&promoCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.PromotionNotFound
		}
		return protocol.DatabaseError
	}

	// 批准：使用 PromotionValues 的方法
	currentTime := time.Now().Unix()
	promoCode.Approve(req.UserID, req.Notes).
		SetStatus(protocol.StatusActive).
		SetActivatedAt(currentTime)

	if err := s.db.Save(&promoCode).Error; err != nil {
		return protocol.PromotionApprovalFailed
	}

	return protocol.Success
}

// RejectPromotion 审批拒绝优惠码
func (s *AdminService) RejectPromotion(req *protocol.ApprovePromotionRequest) protocol.ErrorCode {
	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", req.PromotionID).First(&promoCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.PromotionNotFound
		}
		return protocol.DatabaseError
	}

	// 拒绝：使用 PromotionValues 的方法
	promoCode.Reject(req.UserID, req.Notes).
		SetStatus(protocol.StatusInactive)

	if err := s.db.Save(&promoCode).Error; err != nil {
		return protocol.PromotionApprovalFailed
	}

	return protocol.Success
} // GetPromotionUsage 获取优惠码使用统计
func (s *AdminService) GetPromotionUsage(promoCodeID string, startDate, endDate *int64) (map[string]interface{}, protocol.ErrorCode) {
	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", promoCodeID).First(&promoCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, protocol.PromotionNotFound
		}
		return nil, protocol.DatabaseError
	}

	// 获取基本使用统计
	usage := map[string]interface{}{
		"promotion_id":   promoCodeID,
		"code":           promoCode.GetCode(),
		"title":          promoCode.GetTitle(),
		"usage_count":    promoCode.GetUsageCount(),
		"usage_limit":    promoCode.GetUsageLimit(),
		"total_discount": 0.0,
		"order_count":    0,
		"revenue_impact": 0.0,
	}

	// 查询该优惠码相关的订单使用情况
	orderQuery := s.db.Model(&models.Order{}).Where("promo_code = ?", promoCode.GetCode())

	if startDate != nil {
		orderQuery = orderQuery.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		orderQuery = orderQuery.Where("created_at <= ?", *endDate)
	}

	var orderCount int64
	var totalDiscount float64
	var revenueImpact float64

	// 统计订单数量
	if err := orderQuery.Count(&orderCount).Error; err == nil {
		usage["order_count"] = orderCount
	}

	// 统计总折扣金额
	if err := orderQuery.Select("COALESCE(SUM(promo_discount), 0)").Scan(&totalDiscount).Error; err == nil {
		usage["total_discount"] = totalDiscount
	}

	// 统计收入影响（总订单金额）
	if err := orderQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&revenueImpact).Error; err == nil {
		usage["revenue_impact"] = revenueImpact
	}

	// 按时间分组的使用统计（如果提供了时间范围）
	if startDate != nil && endDate != nil {
		var dailyUsage []map[string]any

		// 这里可以添加按日期分组的统计查询
		// 暂时简化处理
		usage["daily_usage"] = dailyUsage
	}

	return usage, protocol.Success
}

// DeletePromotion 硬删除优惠码
func (s *AdminService) DeletePromotion(promotionID, operatorID string) protocol.ErrorCode {
	// 检查优惠码是否存在
	var promoCode models.Promotion
	if err := s.db.Where("promotion_id = ?", promotionID).First(&promoCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return protocol.PromotionNotFound
		}
		return protocol.DatabaseError
	}

	// 检查优惠码是否已被使用（防止删除已使用的优惠码）
	var usageCount int64
	if err := s.db.Model(&models.UserPromotion{}).Where("promotion_id = ?", promotionID).Count(&usageCount).Error; err != nil {
		return protocol.DatabaseError
	}

	if usageCount > 0 {
		return protocol.PromotionInUse
	}

	// 执行硬删除
	if err := s.db.Where("promotion_id = ?", promotionID).Delete(&models.Promotion{}).Error; err != nil {
		return protocol.PromotionDeleteFailed
	}

	return protocol.Success
}

// validatePromotionRequest 验证优惠券请求的业务逻辑
func (s *AdminService) validatePromotionRequest(req *protocol.CreatePromotionRequest) protocol.ErrorCode {
	// 验证折扣值合理性
	if req.DiscountType == protocol.PromoDiscountTypePercentage {
		if req.DiscountValue <= 0 || req.DiscountValue > 100 {
			return protocol.PromotionValueInvalid
		}
	} else if req.DiscountType == protocol.PromoDiscountTypeFixedAmount {
		if req.DiscountValue <= 0 {
			return protocol.PromotionValueInvalid
		}
	}

	// 验证最大折扣金额
	if req.MaxDiscountAmount != nil && *req.MaxDiscountAmount <= 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证最小订单金额
	if req.MinOrderAmount != nil && *req.MinOrderAmount < 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证时间逻辑
	if req.StartDate != nil && req.EndDate != nil {
		if *req.StartDate >= *req.EndDate {
			return protocol.PromotionTimeInvalid
		}
		// 检查开始时间是否在过去
		if *req.StartDate < time.Now().Unix() {
			return protocol.PromotionTimeInvalid
		}
	}

	// 验证使用限制
	if req.UsageLimit != nil && *req.UsageLimit <= 0 {
		return protocol.PromotionValueInvalid
	}

	if req.UserUsageLimit != nil && *req.UserUsageLimit <= 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证优惠码格式
	if len(req.Code) < 3 || len(req.Code) > 50 {
		return protocol.PromotionValueInvalid
	}

	// 验证标题长度
	if len(req.Title) == 0 || len(req.Title) > 255 {
		return protocol.PromotionValueInvalid
	}

	return protocol.Success
}

// validatePromotionUpdateRequest 验证优惠券更新请求的业务逻辑
func (s *AdminService) validatePromotionUpdateRequest(req *protocol.UpdatePromotionRequest) protocol.ErrorCode {
	// 验证折扣值合理性
	if req.DiscountValue != nil && *req.DiscountValue <= 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证最大折扣金额
	if req.MaxDiscountAmount != nil && *req.MaxDiscountAmount <= 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证最小订单金额
	if req.MinOrderAmount != nil && *req.MinOrderAmount < 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证时间逻辑
	if req.StartDate != nil && req.EndDate != nil {
		if *req.StartDate >= *req.EndDate {
			return protocol.PromotionTimeInvalid
		}
	}

	// 验证使用限制
	if req.UsageLimit != nil && *req.UsageLimit <= 0 {
		return protocol.PromotionValueInvalid
	}

	if req.UserUsageLimit != nil && *req.UserUsageLimit <= 0 {
		return protocol.PromotionValueInvalid
	}

	// 验证标题长度
	if req.Title != nil && (len(*req.Title) == 0 || len(*req.Title) > 255) {
		return protocol.PromotionValueInvalid
	}

	return protocol.Success
}
