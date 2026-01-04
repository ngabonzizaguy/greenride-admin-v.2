package handlers

import (
	"net/http"

	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 本地广告接口
// =============================================================================

// GetLocalAdvertisements 获取本地广告列表
// @Summary 获取本地广告列表
// @Description 获取本地商家广告列表，支持按城市、地区、类别筛选
// @Tags Api 广告
// @Accept json
// @Produce json
// @Param request body protocol.LocalAdvertisementListRequest true "查询请求"
// @Success 200 {object} protocol.Result{data=[]protocol.LocalAdvertisement} "成功返回广告列表"
// @Failure 400 {object} protocol.Result "参数错误"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /ads/list [post]
func (a *Api) GetLocalAdvertisements(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 绑定JSON参数
	var req protocol.LocalAdvertisementListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Get().Errorf("参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 获取本地广告服务
	service := services.GetLocalAdvertisementService()

	// 调用服务获取广告列表
	advertisements, errorCode := service.GetLocalAdvertisements(&req)
	if errorCode != protocol.Success {
		log.Get().Errorf("获取本地广告列表失败: error_code=%s, city=%s, region=%s, category=%s",
			errorCode, req.City, req.Region, req.Category)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errorCode, lang))
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(advertisements, lang))

	log.Get().Infof("获取本地广告列表成功: city=%s, region=%s, category=%s, count=%d",
		req.City, req.Region, req.Category, len(advertisements))
}

// UpdateAdvertisementStats 更新广告统计信息
// @Summary 更新广告统计信息
// @Description 记录广告的浏览、点击、电话等统计信息
// @Tags Api 广告
// @Accept json
// @Produce json
// @Param request body protocol.LocalAdvertisementStatsRequest true "统计请求"
// @Success 200 {object} protocol.Result "更新成功"
// @Failure 400 {object} protocol.Result "参数错误"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /ads/stats [post]
func (a *Api) UpdateAdvertisementStats(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 绑定请求参数
	var req protocol.LocalAdvertisementStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Get().Errorf("参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 获取本地广告服务
	service := services.GetLocalAdvertisementService()

	// 调用服务更新统计信息
	errorCode := service.UpdateAdvertisementStats(&req)
	if errorCode != protocol.Success {
		log.Get().Errorf("更新广告统计失败: error_code=%s, ad_id=%s, stats_type=%s",
			errorCode, req.AdID, req.StatsType)
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResult(errorCode, lang))
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))

	log.Get().Infof("更新广告统计成功: ad_id=%s, stats_type=%s, user_id=%s",
		req.AdID, req.StatsType, req.UserID)
}

// GetLocalAdvertisementByID 获取单个广告详情
// @Summary 获取广告详情
// @Description 根据广告ID获取详细信息
// @Tags Api 广告
// @Accept json
// @Produce json
// @Param request body protocol.LocalAdvertisementDetailRequest true "详情请求"
// @Success 200 {object} protocol.Result{data=protocol.LocalAdvertisement} "成功返回广告详情"
// @Failure 400 {object} protocol.Result "参数错误"
// @Failure 404 {object} protocol.Result "广告不存在"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /ads/detail [post]
func (a *Api) GetLocalAdvertisementByID(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	// 绑定JSON参数
	var req protocol.LocalAdvertisementDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Get().Errorf("参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidParams, lang))
		return
	}

	// 获取本地广告服务
	service := services.GetLocalAdvertisementService()

	// 调用服务获取广告详情
	advertisement, errorCode := service.GetLocalAdvertisementByID(req.AdID)
	if errorCode != protocol.Success {
		log.Get().Errorf("获取广告详情失败: error_code=%s, ad_id=%s", errorCode, req.AdID)
		statusCode := http.StatusInternalServerError
		if errorCode == protocol.UserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, protocol.NewErrorResult(errorCode, lang))
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(advertisement, lang))

	log.Get().Infof("获取广告详情成功: ad_id=%s", req.AdID)
}
