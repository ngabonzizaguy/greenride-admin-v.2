package handlers

import (
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取优惠券列表
// @Description 获取优惠券列表，支持分页和过滤
// @Tags Api,优惠券
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.UserPromotionsRequest true "优惠券搜索条件"
// @Router /promotions [post]
func (t *Api) Promotions(c *gin.Context) {
	lang := middleware.GetLanguageFromContext(c)
	// 解析请求体
	var req protocol.UserPromotionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResult(protocol.InvalidJSON, lang, err.Error()))
		return
	}
	user := GetUserFromContext(c)
	if !user.IsPassenger() {
		c.JSON(http.StatusForbidden, protocol.NewErrorResult(protocol.PermissionDenied, lang))
		return
	}
	req.UserID = user.UserID
	req.UserType = user.GetUserType()

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Status == "" {
		req.Status = protocol.StatusAvailable
	}

	// 获取优惠码列表
	promotions, total := services.GetPromotionService().GetPromotionByUser(&req)
	// 返回结果
	result := protocol.NewPageResult(promotions, total, &protocol.Pagination{
		Page: req.Page,
		Size: req.Limit,
	})
	result.AddAttach("params", req)
	c.JSON(http.StatusOK, protocol.NewSuccessResult(result))
}
