package models

import (
	"fmt"

	"greenride/internal/config"
	"greenride/internal/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	// 使用配置中的DSN
	dsn := cfg.Database.DSN
	if dsn == "" {
		return fmt.Errorf("database DSN is empty")
	}

	var logLevel logger.LogLevel
	if cfg.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	DB = db
	return nil
}

// AutoMigrate 运行数据库迁移
func AutoMigrate() error {
	log.Get().Info("Running database migrations...")

	// 定义所有需要迁移的表模型
	tables := []any{
		// 用户相关
		&User{},
		&UserAccount{},
		&UserAddress{},
		&UserPaymentMethod{},
		&UserPromotion{},
		&UserLocationHistory{},

		// 管理员
		&Admin{},

		// 身份验证
		&Identity{},

		// 订单相关
		&Order{},
		&OrderRating{},
		&OrderHistoryLog{},
		// &OrderStats{}, // Comented out if it doesn't exist
		&RideOrder{},

		// 派单相关
		&DispatchRecord{},

		// 车辆相关
		&Vehicle{},
		&VehicleType{},

		// 价格相关
		&PriceRule{},
		&PriceSnapshot{},

		// 支付相关
		&Payment{},
		&PaymentMethod{},
		&PaymentChannels{},

		// 钱包相关
		&Wallet{},
		&WalletTransaction{},
		&Withdrawal{},

		// 促销相关
		&Promotion{},

		// 消息相关
		&Message{},
		&MessageTemplate{},
		&Notification{},

		// FCM相关
		&FCMToken{},
		&FCMMessageLog{},

		// 服务区域
		&ServiceArea{},

		// 公告
		&Announcement{},

		// 反馈
		&Feedback{},

		// 支持配置
		&SupportConfig{},

		// 任务
		&Task{},
	}

	for _, table := range tables {
		if err := DB.AutoMigrate(table); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	log.Get().Info("Database migrations completed successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
