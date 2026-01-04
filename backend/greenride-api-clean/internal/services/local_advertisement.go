package services

import (
	"context"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"sync"
)

// LocalAdvertisementService 本地广告服务
type LocalAdvertisementService struct {
}

var (
	localAdServiceInstance *LocalAdvertisementService
	localAdServiceOnce     sync.Once
)

// GetLocalAdvertisementService 获取本地广告服务单例
func GetLocalAdvertisementService() *LocalAdvertisementService {
	localAdServiceOnce.Do(func() {
		SetupLocalAdvertisementService()
	})
	return localAdServiceInstance
}

// SetupLocalAdvertisementService 设置本地广告服务
func SetupLocalAdvertisementService() {
	localAdServiceInstance = &LocalAdvertisementService{}
}

// GetLocalAdvertisements 获取本地广告列表
func (s *LocalAdvertisementService) GetLocalAdvertisements(req *protocol.LocalAdvertisementListRequest) ([]protocol.LocalAdvertisement, protocol.ErrorCode) {
	if req == nil {
		req = &protocol.LocalAdvertisementListRequest{}
	}

	// 查询广告，现在支持category参数
	ads, err := models.GetAllActiveLocalAdvertisements(req.City, req.Region, req.Category)
	if err != nil {
		log.Get().Errorf("获取本地广告列表失败: %v", err)
		return nil, protocol.DatabaseError
	}

	// 转换为响应格式
	var responses []protocol.LocalAdvertisement
	for _, ad := range ads {
		response := ad.Protocol()
		responses = append(responses, response)
	}

	log.Get().Infof("获取本地广告列表成功: city=%s, region=%s, category=%s, total=%d",
		req.City, req.Region, req.Category, len(responses))

	return responses, protocol.Success
}

// UpdateAdvertisementStats 更新广告统计信息
func (s *LocalAdvertisementService) UpdateAdvertisementStats(req *protocol.LocalAdvertisementStatsRequest) protocol.ErrorCode {
	if req == nil || req.AdID == "" || req.StatsType == "" {
		return protocol.InvalidParams
	}

	// 验证统计类型
	if req.StatsType != "view" && req.StatsType != "click" && req.StatsType != "call" {
		return protocol.InvalidParams
	}

	// 更新统计信息
	err := models.UpdateLocalAdvertisementStats(models.DB, req.AdID, req.StatsType)
	if err != nil {
		log.Get().Errorf("更新广告统计失败: ad_id=%s, stats_type=%s, error=%v", req.AdID, req.StatsType, err)
		return protocol.DatabaseError
	}

	log.Get().Infof("更新广告统计成功: ad_id=%s, stats_type=%s, user_id=%s", req.AdID, req.StatsType, req.UserID)
	return protocol.Success
}

// GetLocalAdvertisementByID 根据ID获取广告详情
func (s *LocalAdvertisementService) GetLocalAdvertisementByID(adID string) (*protocol.LocalAdvertisement, protocol.ErrorCode) {
	if adID == "" {
		return nil, protocol.InvalidParams
	}

	ad, err := models.GetLocalAdvertisementByID(adID)
	if err != nil {
		log.Get().Errorf("获取广告详情失败: ad_id=%s, error=%v", adID, err)
		return nil, protocol.UserNotFound // 使用已存在的错误码
	}

	if !ad.IsVisible() {
		return nil, protocol.UserNotFound // 使用已存在的错误码
	}

	response := ad.Protocol()
	return &response, protocol.Success
}

// UpdateLocalAdvertisementWithGoogleInfo 根据地点名称更新Google信息
func (s *LocalAdvertisementService) UpdateLocalAdvertisementWithGoogleInfo(ctx context.Context, placeName string, city string, country string) error {
	// 获取Google服务
	googleService := GetGoogleService()
	if googleService == nil {
		return protocol.NewError(protocol.SystemError, "Google service not available")
	}

	// 根据名称查找本地广告
	ads, err := models.GetLocalAdvertisementsByName(placeName)
	if err != nil {
		log.Get().Errorf("查找本地广告失败: name=%s, error=%v", placeName, err)
		return err
	}

	if len(ads) == 0 {
		log.Get().Warnf("未找到名称匹配的本地广告: name=%s", placeName)
		return protocol.NewError(protocol.UserNotFound, "No local advertisements found with name: "+placeName)
	}

	// 对每个匹配的广告进行处理
	for _, ad := range ads {
		log.Get().Infof("处理广告: ad_id=%s, name=%s", ad.AdID, ad.GetName())

		// 1. 搜索Google Place
		placeResult, err := googleService.SearchPlacesByName(ctx, placeName, city, country)
		if err != nil {
			log.Get().Errorf("搜索Google Place失败: name=%s, error=%v", placeName, err)
			continue // 继续处理下一个广告
		}

		log.Get().Infof("找到Google Place: place_id=%s, name=%s", placeResult.PlaceID, placeResult.Name)

		// 2. 获取Place详情
		placeDetails, err := googleService.GetPlaceDetails(ctx, placeResult.PlaceID)
		if err != nil {
			log.Get().Errorf("获取Google Place详情失败: place_id=%s, error=%v", placeResult.PlaceID, err)
			continue // 继续处理下一个广告
		}

		// 3. 提取完整的Google数据
		googleData := googleService.ExtractCompleteGoogleData(placeDetails)

		// 4. 使用完整数据更新本地广告信息
		err = models.UpdateLocalAdvertisementWithCompleteGoogleInfo(ad.AdID, googleData)
		if err != nil {
			log.Get().Errorf("更新本地广告完整Google信息失败: ad_id=%s, error=%v", ad.AdID, err)
			continue
		}

		log.Get().Infof("成功更新广告完整Google信息: ad_id=%s, place_id=%s, rating=%.1f, phone=%s, website=%s, address=%s",
			ad.AdID, placeDetails.PlaceID, placeDetails.Rating,
			placeDetails.InternationalPhoneNumber, placeDetails.Website, placeDetails.FormattedAddress)
	}

	return nil
}
