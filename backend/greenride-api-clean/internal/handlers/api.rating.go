package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 订单评价相关接口
// =============================================================================

// CreateOrderRating 创建订单评价
// @Summary 创建评价
// @Description 对订单进行评价
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.CreateOrderRatingRequest true "评价请求"
// @Success 200 {object} protocol.Result{data=protocol.Rating} "评价成功"
// @Failure 200 {object} protocol.Result "评价失败"
// @Security BearerAuth
// @Router /order/rating [post]
func (a *Api) CreateOrderRating(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	user := GetUserFromContext(c)

	var req protocol.CreateOrderRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateOrderRating: Invalid request parameters: %v", err)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	rating, errCode := services.GetOrderRatingService().CreateRating(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResult(rating.Protocol()))
}

// GetOrderRatings 获取订单评价列表
// @Summary 获取评价列表
// @Description 获取订单的评价列表
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.OrderIDRequest true "订单ID请求"
// @Success 200 {object} protocol.Result{data=[]protocol.Rating} "获取成功"
// @Failure 200 {object} protocol.Result "获取失败"
// @Security BearerAuth
// @Router /order/ratings [post]
func (a *Api) GetOrderRatings(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.OrderIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("GetOrderRatings: Invalid request parameters: %v", err)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetUserFromContext(c)
	req.UserID = user.UserID
	list := services.GetOrderRatingService().GetRatingsByOrder(&req)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(list))
}

// UpdateOrderRating 更新订单评价
// @Summary 更新评价
// @Description 更新订单评价
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.UpdateOrderRatingRequest true "评价更新请求"
// @Success 200 {object} protocol.Result "更新成功"
// @Failure 200 {object} protocol.Result "更新失败"
// @Security BearerAuth
// @Router /rating/update [post]
func (a *Api) UpdateOrderRating(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.UpdateOrderRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UpdateOrderRating: Invalid request parameters: %v", err)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	if len(req.Comment) > 500 {
		log.Printf("UpdateOrderRating: Comment too long: %d characters", len(req.Comment))
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, "Comment too long (max 500 characters)"))
		return
	}
	user := GetUserFromContext(c)
	req.UserID = user.UserID
	errCode := services.GetOrderRatingService().UpdateRating(&req)
	if errCode != protocol.Success {
		log.Printf("UpdateOrderRating: Failed to update rating %d with error code: %s", req.RatingID, errCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// DeleteOrderRating 删除订单评价
// @Summary 删除评价
// @Description 删除订单评价
// @Tags Api,订单
// @Accept json
// @Produce json
// @Param request body protocol.RatingIDRequest true "评价ID请求"
// @Success 200 {object} protocol.Result "删除成功"
// @Failure 200 {object} protocol.Result "删除失败"
// @Security BearerAuth
// @Router /rating/delete [post]
func (a *Api) DeleteOrderRating(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	user := GetUserFromContext(c)
	var req protocol.RatingIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("DeleteOrderRating: Invalid request parameters: %v", err)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	req.UserID = user.UserID
	errCode := services.GetOrderRatingService().DeleteRating(&req)
	if errCode != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}

// ReplyToRating 回复评价
// @Summary 回复评价
// @Description 服务提供者回复用户评价
// @Tags Api,司机
// @Accept json
// @Produce json
// @Param request body protocol.ReplyToRatingRequest true "回复请求"
// @Success 200 {object} protocol.Result "回复成功"
// @Failure 200 {object} protocol.Result "回复失败"
// @Security BearerAuth
// @Router /rating/reply [post]
func (a *Api) ReplyToRating(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)

	var req protocol.ReplyToRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ReplyToRating: Invalid request parameters: %v", err)
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}

	if len(req.Reply) > 300 {
		log.Printf("ReplyToRating: Reply too long: %d characters", len(req.Reply))
		c.JSON(http.StatusOK, protocol.NewErrorResult(protocol.InvalidParams, lang, "Reply too long (max 300 characters)"))
		return
	}
	user := GetUserFromContext(c)
	req.UserID = user.UserID
	errCode := services.GetOrderRatingService().ReplyToRating(&req)
	if errCode != protocol.Success {
		log.Printf("ReplyToRating: Failed to reply to rating %s with error code: %s", req.RatingID, errCode)
		c.JSON(http.StatusOK, protocol.NewErrorResult(errCode, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(""))
}
