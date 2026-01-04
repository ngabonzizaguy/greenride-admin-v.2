package main

import (
	"fmt"
	"greenride/internal/config"
	"greenride/internal/handlers"
	"greenride/internal/i18n"
	"greenride/internal/log"
	"greenride/internal/models"
	"greenride/internal/services"
	"greenride/internal/task"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @title        GreenRide API
// @version      1.0
// @description  GreenRide API Documentation
// @description ## Response Codes:
// @description - 0000: Success - 操作成功
// @description - 9999: System Error - 系统错误
// @description - 1000: Parameter Error - 参数错误
// @description - 1001: Authorization Error - 授权失败
// @description - 1002: Authentication Expired - 认证过期
// @description - 2001: User Error - 用户相关错误
// @description - 4001: Business Error - 业务错误
// @description - 4002: Data Expired - 数据过期

// @BasePath    /
// @schemes     http

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization

func InitialConfig() error {
	cfg := config.Get()

	// 初始化数据库
	if err := models.InitDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	// 运行数据库迁移
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
		return err
	}
	// 初始化Redis
	if err := models.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return err
	}

	// 初始化国际化
	translator := i18n.NewFileTranslator()
	i18n.SetGlobalTranslator(translator)

	// 加载翻译资源，仅在配置了路径时加载
	if cfg.I18n != nil && cfg.I18n.LocalesDir != "" {
		localesDir := cfg.I18n.LocalesDir
		if err := translator.LoadTranslations(localesDir); err != nil {
			log.Warnf("I18n: Failed to load translations from '%s': %v", localesDir, err)
		} else {
			log.Infof("I18n: Translations loaded successfully from '%s'", localesDir)
		}
	} else {
		log.Info("I18n: locales_dir not configured, using default English messages")
	}
	services.SetupService()
	services.GetAdminAdminService().EnsureDefaultAdmin()
	if err := task.Run(); err != nil {
		log.Fatalf("Failed to start task service: %v", err)
	}
	return nil
}

var (
	g errgroup.Group
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		panic(err)
	}

	// 初始化应用配置
	err = InitialConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.Get()

	// 设置Gin模式
	if cfg.Env == config.ProdEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	// 启动API服务
	g.Go(func() error {
		apiService := handlers.NewApi()
		if apiService == nil {
			log.Fatal("Failed to create API service - configuration may be invalid")
			return fmt.Errorf("failed to create API service")
		}
		server := apiService.ToServer()
		server.Handler = apiService.SetupRouter()

		log.Infof("Starting API Service on port %s", apiService.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("API Service error: %v", err)
		} else {
			log.Info("API Service started successfully")
		}
		return err
	})

	// 启动Admin服务
	g.Go(func() error {
		adminService := handlers.NewAdmin()
		if adminService == nil {
			log.Fatal("Failed to create Admin service - configuration may be invalid")
			return fmt.Errorf("failed to create Admin service")
		}
		server := adminService.ToServer()
		server.Handler = adminService.SetupRouter()

		log.Infof("Starting Admin Service on port %s", adminService.Port)
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("Admin Service error: %v", err)
		} else {
			log.Info("Admin Service started successfully")
		}
		return err
	})

	// 等待所有服务
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
