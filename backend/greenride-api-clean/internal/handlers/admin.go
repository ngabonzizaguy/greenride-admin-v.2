package handlers

import (
	"log"
	"net/http"

	admindocs "greenride/docs/admin"
	"greenride/internal/config"
	"greenride/internal/middleware"
	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Admin 管理员处理器
type Admin struct {
	*config.ServiceConfig
}

// NewAdmin 创建管理员处理器
func NewAdmin() *Admin {
	cfg := config.Get()
	if cfg == nil {
		return nil
	}
	adminCfg := cfg.Server.Admin
	if adminCfg == nil {
		return nil
	}

	return &Admin{
		ServiceConfig: adminCfg,
	}
}

// SetupRouter sets up the router for the Admin service
// 在 SetupRouter 函数中添加管理员管理路由
func (t *Admin) SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	cfg := config.Get()
	if cfg.Env == protocol.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Local Development CORS
	if cfg.Env != protocol.EnvProduction {
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})
	}
	router.Use(middleware.RequestLogger()) // 添加自定义请求日志中间件

	// CORS 已移至 nginx 统一处理

	// Initialize Swagger Info
	admindocs.SwaggerInfoadmin.Title = "GreenRide Admin API"
	admindocs.SwaggerInfoadmin.Description = "GreenRide Admin API Documentation"
	admindocs.SwaggerInfoadmin.Version = "1.0"
	admindocs.SwaggerInfoadmin.Host = ""
	admindocs.SwaggerInfoadmin.BasePath = "/"
	admindocs.SwaggerInfoadmin.Schemes = []string{"http", "https"}

	// Add Swagger documentation route
	if cfg.Env != "prod" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("admin")))
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "admin"})
	})

	// 管理员认证相关（无前缀，简化路径）
	router.POST("/login", t.Login)

	// Admin routes - 移除 /admin 前缀，简化为根路径
	adminAPI := router.Group("/")
	{

		// 需要JWT认证的端点
		adminAPI.Use(t.AuthMiddleware())
		{
			adminAPI.POST("/logout", t.Logout)
			adminAPI.GET("/info", t.Info)
			adminAPI.POST("/change-password", t.ChangePassword)
			adminAPI.POST("/reset-password", t.ResetPassword)

			// Dashboard 统计相关
			dashboardAPI := adminAPI.Group("/dashboard")
			{
				dashboardAPI.GET("/stats", t.GetDashboardStats)        // 获取仪表盘统计数据
				dashboardAPI.GET("/revenue", t.GetRevenueChart)        // 获取收入图表数据
				dashboardAPI.GET("/user-growth", t.GetUserGrowthChart) // 获取用户增长图表数据
			}

			// 用户管理相关
			userAPI := adminAPI.Group("/users")
			{
				userAPI.POST("/search", t.GetUserList)      // 搜索用户（支持分页）
				userAPI.POST("/detail", t.GetUserDetail)    // 获取用户详情
				userAPI.POST("/create", t.CreateUser)       // 管理员创建用户
				userAPI.POST("/update", t.UpdateUser)       // 更新用户信息
				userAPI.POST("/status", t.UpdateUserStatus) // 统一的状态更新接口
				userAPI.POST("/verify", t.VerifyUser)       // 审核用户认证（替代VerifyDriver）
				userAPI.POST("/rides", t.GetUserRides)      // 获取用户行程历史
			}

			// 车辆管理相关
			vehicleAPI := adminAPI.Group("/vehicles")
			{
				vehicleAPI.POST("/search", t.SearchVehicles)      // 搜索车辆
				vehicleAPI.POST("/detail", t.GetVehicleDetail)    // 获取车辆详情
				vehicleAPI.POST("/update", t.UpdateVehicle)       // 更新车辆信息
				vehicleAPI.POST("/status", t.UpdateVehicleStatus) // 更新车辆状态
				vehicleAPI.POST("/delete", t.DeleteVehicle)       // 删除车辆
				vehicleAPI.POST("/create", t.CreateVehicle)       // 创建车辆
			}
			// 订单管理相关
			ordersAPI := adminAPI.Group("/orders")
			{
				ordersAPI.POST("/search", t.SearchOrders)    // 搜索订单
				ordersAPI.POST("/detail", t.GetOrderDetail)  // 获取订单详情
				ordersAPI.POST("/estimate", t.EstimateOrder) // 管理员订单预估
				ordersAPI.POST("/create", t.CreateOrder)     // 管理员创建订单
				ordersAPI.POST("/cancel", t.CancelOrder)     // 取消订单
			}

			// 反馈/投诉管理相关
			feedbackAPI := adminAPI.Group("/feedback")
			{
				feedbackAPI.POST("/search", t.SearchFeedback)     // 搜索反馈列表
				feedbackAPI.POST("/detail", t.GetFeedbackDetail)  // 获取反馈详情
				feedbackAPI.POST("/update", t.UpdateFeedback)     // 更新反馈
				feedbackAPI.GET("/stats", t.GetFeedbackStats)     // 获取反馈统计
			}

			// 支持配置相关
			supportAPI := adminAPI.Group("/support")
			{
				supportAPI.GET("/config", t.GetSupportConfig)     // 获取支持配置
				supportAPI.POST("/config", t.UpdateSupportConfig) // 更新支持配置
			}
		}
	}

	return router
}

// AuthMiddleware JWT认证中间件
func (t *Admin) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.ValidToken(c, []byte(t.Jwt.Secret))
		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, protocol.NewAuthErrorResult())
			c.Abort()
			return
		}
		// 提取claims
		if claims, ok := token.Claims.(*middleware.JWTClaims); ok {
			// 直接从数据库获取完整的用户对象
			user := services.GetAdminAdminService().GetAdminByID(claims.UserID)
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

// GetDashboardStats 获取仪表盘统计数据
// @Summary 获取仪表盘统计数据
// @Description 获取用户、司机、车辆、行程等核心统计指标
// @Tags Admin,Dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=services.DashboardStats}
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /dashboard/stats [get]
func (a *Admin) GetDashboardStats(c *gin.Context) {
	stats, err := services.GetDashboardService().GetDashboardStats()
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		// 返回默认数据而不是错误
		defaultStats := &services.DashboardStats{
			TotalUsers:         0,
			TotalDrivers:       0,
			TotalVehicles:      0,
			TotalTrips:         0,
			TotalRevenue:       0.0,
			ActiveTrips:        0,
			MonthlyGrowth:      services.MonthlyGrowth{},
			RecentTrips:        []services.RecentTrip{},
			TopDrivers:         []services.TopDriver{},
			VehicleUtilization: []services.VehicleUtilization{},
		}
		c.JSON(http.StatusOK, protocol.NewSuccessResult(defaultStats))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(stats))
}

// GetRevenueChart 获取收入图表数据
// @Summary 获取收入图表数据
// @Description 获取指定时间段的收入趋势图表数据
// @Tags Admin,Dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param period query string false "时间段" Enums(7d,30d,12m) default(12m)
// @Param timezone query string false "时区" default(Africa/Kigali)
// @Success 200 {object} protocol.Result{data=[]services.RevenueData}
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /dashboard/revenue [get]
func (a *Admin) GetRevenueChart(c *gin.Context) {
	period := c.DefaultQuery("period", "12m")
	timezone := c.DefaultQuery("timezone", "")

	data, err := services.GetDashboardService().GetRevenueChart(period, timezone)
	if err != nil {
		log.Printf("Error getting revenue chart data: %v", err)
		c.JSON(http.StatusInternalServerError, protocol.NewBusinessErrorResult("Failed to get revenue chart data"))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(data))
}

// GetUserGrowthChart 获取用户增长图表数据
// @Summary 获取用户增长图表数据
// @Description 获取指定时间段的用户和司机增长趋势图表数据
// @Tags Admin,Dashboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param period query string false "时间段" Enums(7d,30d,12m) default(12m)
// @Param timezone query string false "时区" default(Africa/Kigali)
// @Success 200 {object} protocol.Result{data=[]services.UserGrowthData}
// @Failure 401 {object} protocol.Result
// @Failure 500 {object} protocol.Result
// @Router /dashboard/user-growth [get]
func (a *Admin) GetUserGrowthChart(c *gin.Context) {
	period := c.DefaultQuery("period", "12m")
	timezone := c.DefaultQuery("timezone", "")

	data, err := services.GetDashboardService().GetUserGrowthChart(period, timezone)
	if err != nil {
		log.Printf("Error getting user growth chart data: %v", err)
		c.JSON(http.StatusInternalServerError, protocol.NewBusinessErrorResult("Failed to get user growth chart data"))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(data))
}

// GetUserFromContext 从上下文获取完整的用户对象
func (s *Admin) GetUserFromContext(c *gin.Context) *models.Admin {
	value, exists := c.Get("user")
	if !exists {
		return nil
	}
	if user, ok := value.(*models.Admin); ok {
		return user
	}
	return nil
}
