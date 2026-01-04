package services

import (
	"greenride/internal/models"
	"greenride/internal/protocol"
	"sync"
)

type PromotionService struct {
}

var (
	promotionServiceInstance     *PromotionService
	promotionServiceInstanceOnce = sync.Once{}
)

func GetPromotionService() *PromotionService {
	if promotionServiceInstance == nil {
		SetupPromotionService()
	}
	return promotionServiceInstance
}

func SetupPromotionService() {
	promotionServiceInstanceOnce.Do(func() {
		promotionServiceInstance = &PromotionService{}
	})
}

func (s *PromotionService) GetPromotionByUser(req *protocol.UserPromotionsRequest) (list []*protocol.UserPromotion, total int64) {
	query := models.GetDB().Model(&models.UserPromotion{})
	query.Where("user_id = ?", req.UserID)
	if req.Status != "" {
		query = query.Where("status=?", req.Status)
	}
	// 计算总数
	if _err := query.Count(&total).Error; _err != nil {
		return
	}

	// 获取订单列表
	var results []*models.UserPromotion
	var result_list []int64
	offset := (req.Page - 1) * req.Limit
	if _err := query.Select([]string{"id"}).Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&result_list).Error; _err != nil {
		return
	}
	// 根据订单ID列表获取完整的订单信息
	if len(result_list) > 0 {
		if _err := models.GetDB().Where("id IN ?", result_list).Order("created_at DESC").Find(&results).Error; _err != nil {
			return
		}
	}
	for _, item := range results {
		list = append(list, item.Protocol())
	}

	return
}
