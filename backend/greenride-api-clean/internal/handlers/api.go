package handlers

import (
	"net/http"

	apidocs "greenride/docs/api"
	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/middleware"
	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// AuthMiddleware JWT认证中间件
func (a *Api) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.ValidToken(c, []byte(a.Jwt.Secret))
		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}

		// 提取用户信息
		if claims, ok := token.Claims.(*middleware.JWTClaims); ok {
			// 直接从数据库获取完整的用户对象
			user := services.GetUserService().GetUserByID(claims.UserID)
			if user == nil {
				c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
				c.Abort()
				return
			}

			// 检查用户状态
			if user.GetStatus() != protocol.StatusActive {
				c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
				c.Abort()
				return
			}

			// 将完整的用户对象存储到上下文中
			c.Set("user", user)

			// 保持向后兼容，也设置单独的键（可选）
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("user_type", claims.UserType)
		}

		c.Next()
	}
}

// SetupRouter sets up the router for the API service
func (a *Api) SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	cfg := config.Get()
	if cfg.Env == protocol.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger()) // 添加自定义请求日志中间件

	// CORS 已移至 nginx 统一处理

	// Initialize Swagger Info
	apidocs.SwaggerInfoapi.Title = "GreenRide API"
	apidocs.SwaggerInfoapi.Description = "GreenRide API Documentation"
	apidocs.SwaggerInfoapi.Version = "1.0"
	apidocs.SwaggerInfoapi.Host = ""
	apidocs.SwaggerInfoapi.BasePath = "/"
	apidocs.SwaggerInfoapi.Schemes = []string{"http", "https"}

	// Add Swagger documentation route
	if cfg.Env != "prod" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("api")))
	}

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": a.Name,
			"version": a.Version,
			"port":    a.Port,
		})
	})

	api := router.Group("/")
	{
		// 无需认证的路由
		api.POST("/register", a.Register)
		api.POST("/login", a.Login)
		api.POST("/send-verify-code", a.SendVerifyCode)
		api.POST("/verify-code", a.VerifyCode)
		api.POST("/reset-password", a.ResetPassword)
		api.POST("/feedback/submit", a.SubmitFeedback) // 提交反馈 - 无需认证
		api.GET("/support/config", a.GetSupportConfig)  // 获取支持配置 - 无需认证（公共信息）

		// Checkout 状态查询接口
		api.POST("/checkout/status", a.GetCheckoutStatus) // 查询checkout状态

		// Webhook 回调接口 - 无需认证（第三方支付回调）
		api.POST("/webhook/kpay/:payment_id", a.KPayWebhook)   // KPay 支付回调
		api.POST("/webhook/momo/:payment_id", a.MoMoWebhook)   // MTN MoMo 支付回调
		api.POST("/webhook/stripe", a.StripeWebhook)           // Stripe 支付回调
		api.POST("/webhook/innopaas", a.InnoPaaSWebhook)       // InnoPaaS OTP/消息状态回调

	}

	// 需要认证的路由
	authRequired := api.Group("")
	authRequired.Use(a.AuthMiddleware()) // 后续添加认证中间件
	{
		authRequired.GET("/profile", a.Profile)
		authRequired.POST("/logout", a.Logout)
		authRequired.POST("/change-password", a.ChangePassword)

		// 新增接口
		authRequired.POST("/online", a.UserOnline)                  // 司机上线
		authRequired.POST("/offline", a.UserOffline)                // 司机下线
		authRequired.POST("/profile/update", a.UpdateProfile)       // 更新个人信息
		authRequired.POST("/profile/update/avatar", a.UpdateAvatar) // 更新用户头像
		authRequired.POST("/account/delete", a.DeleteAccount)       // 删除账户

		// 订单相关接口 (通用订单系统，支持网约车等多种订单类型)
		// 订单预估接口
		authRequired.POST("/order/estimate", a.EstimateOrder) // 预估订单

		// 通用订单接口
		authRequired.POST("/order/create", a.CreateOrder)       // 创建订单
		authRequired.POST("/orders", a.GetOrders)               // 获取订单列表
		authRequired.POST("/order/detail", a.GetOrderDetail)    // 获取订单详情
		authRequired.POST("/order/accept", a.AcceptOrder)       // 接受订单
		authRequired.POST("/order/reject", a.RejectOrder)       // 拒绝订单
		authRequired.POST("/order/arrived", a.ArrivedOrder)     // 到达订单
		authRequired.POST("/order/start", a.StartOrder)         // 开始订单
		authRequired.POST("/order/finish", a.FinishOrder)       // 完成订单
		authRequired.POST("/order/cancel", a.CancelOrder)       // 取消订单
		authRequired.POST("/order/rating", a.CreateOrderRating) // 创建订单评价
		authRequired.POST("/order/ratings", a.GetOrderRatings)  // 获取订单评价
		authRequired.POST("/order/contact", a.GetOrderContact)  // 获取订单联系方式（通话权限）
		authRequired.POST("/order/eta", a.GetOrderETA)          // 获取订单实时ETA

		// 服务提供者接口 (司机、外卖员等)
		authRequired.POST("/nearby", a.GetNearbyOrders)                // 获取附近订单
		authRequired.POST("/order/nearby", a.GetNearbyOrders)          // 获取附近订单
		authRequired.POST("/order/cash/received", a.OrderCashReceived) // 确认现金收款
		authRequired.POST("/order/payment", a.OrderPayment)            // 处理订单支付

		// 支付方式接口
		authRequired.POST("/payment/methods", a.GetPaymentMethods) // 获取支付方式列表
		authRequired.POST("/payment/cancel", a.CancelPayment)      // 取消支付

		// 车辆信息接口
		authRequired.POST("/vehicle", a.GetUserVehicle) // 获取用户车辆信息
		authRequired.POST("/vehicles", a.GetVehicles)   // 获取车辆列表

		// 位置管理接口
		authRequired.POST("/location/update", a.UpdateLocation)  // 更新司机位置
		authRequired.GET("/location/current", a.CurrentLocation) // 获取当前位置
		authRequired.GET("/drivers/nearby", a.GetNearbyDrivers)  // 获取附近司机（乘客用）

		authRequired.POST("/rating/update", a.UpdateOrderRating) // 更新评价
		authRequired.POST("/rating/delete", a.DeleteOrderRating) // 删除评价
		authRequired.POST("/rating/reply", a.ReplyToRating)      // 回复评价

		authRequired.POST("/promotions", a.Promotions) // 获取优惠码列表

		// 本地广告统计接口
		authRequired.POST("/ads/list", a.GetLocalAdvertisements)      // 获取本地广告列表
		authRequired.POST("/ads/detail", a.GetLocalAdvertisementByID) // 获取单个广告详情
		authRequired.POST("/ads/stats", a.UpdateAdvertisementStats)   // 更新广告统计
	}

	log.Infof("API router setup completed for service: %s on port %s", a.ServiceConfig.Name, a.ServiceConfig.Port)
	return router
}

// Api V2版本API处理器
type Api struct {
	*config.ServiceConfig
}

// NewApi 创建V2 API处理器
func NewApi() *Api {
	cfg := config.Get()
	if cfg == nil {
		return nil
	}
	apiCfg := cfg.Server.Api
	if apiCfg == nil {
		return nil
	}

	return &Api{
		ServiceConfig: apiCfg,
	}
}
