package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 状态常量
const (
	StatusOn  = "on"
	StatusOff = "off"
)

// 全局配置实例
var (
	config *Config
)

type Config struct {
	Debug      bool              `mapstructure:"debug"`
	Env        string            `mapstructure:"env"`
	Server     *ServerConfig     `mapstructure:"server"`
	Database   *DatabaseConfig   `mapstructure:"database"`
	Redis      *RedisConfig      `mapstructure:"redis"`
	Log        *LogConfig        `mapstructure:"log"`
	I18n       *I18nConfig       `mapstructure:"i18n"` // 国际化配置
	FCM        *FCMConfig        `mapstructure:"fcm"`
	Firebase   map[string]string `mapstructure:"firebase"` // Firebase服务账户配置 (使用map格式)
	Twilio     *TwilioConfig     `mapstructure:"twilio"`
	Google     *GoogleConfig     `mapstructure:"google"`
	AWS        *AWSConfig        `mapstructure:"aws"`
	Email      *EmailConfig      `mapstructure:"email"`       // 邮件服务配置
	SMS        *SMSConfig        `mapstructure:"sms"`         // SMS服务配置
	VerifyCode *VerifyCodeConfig `mapstructure:"verify_code"` // 验证码配置
	Dispatch   *DispatchConfig   `mapstructure:"dispatch"`    // 派单系统配置
	Task       *TaskConfig       `mapstructure:"task"`        // 任务调度配置
	Promotion  *PromotionConfig  `mapstructure:"promotion"`   // 优惠券配置
	Payment    *PaymentConfig       `mapstructure:"payment"`     // 支付配置
	KPay       *KPayConfig          `mapstructure:"kpay"`        // KPay支付配置
	MoMo       *MoMoGlobalConfig    `mapstructure:"momo"`        // MTN MoMo支付配置
	Stripe     *StripeGlobalConfig  `mapstructure:"stripe"`      // Stripe支付配置
	Order      *OrderConfig         `mapstructure:"order"`       // 订单配置
	InnoPaaS   *InnoPaaSConfig   `mapstructure:"innopaas"`    // InnoPaaS SMS配置
}

func (c *Config) IsSandbox() bool {
	return c.Env != ProdEnv
}

// Validate 验证并设置日志配置默认值
func (l *LogConfig) Validate() {
	if l.Level == "" {
		l.Level = "info"
	}
	if l.Format == "" {
		l.Format = "json"
	}
	if l.Output == "" {
		l.Output = "both"
	}
	// 验证 Loki 配置
	if l.Loki != nil {
		l.Loki.Validate()
	}
}

type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn"` // DSN连接字符串
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	Debug           bool          `mapstructure:"debug"`
}

type FCMConfig struct {
	ServerKey   string `mapstructure:"server_key"`
	ProjectID   string `mapstructure:"project_id"`
	DatabaseURL string `mapstructure:"database_url"`
}

type AWSConfig struct {
	// S3 存储配置
	S3Region      string `mapstructure:"s3_region"`      // S3区域
	S3Bucket      string `mapstructure:"s3_bucket"`      // S3存储桶名称
	CloudFrontURL string `mapstructure:"cloudfront_url"` // CloudFront CDN URL（可选）
}

type WebSocketConfig struct {
	Port            string `mapstructure:"port"`
	ReadBufferSize  int    `mapstructure:"read_buffer_size"`
	WriteBufferSize int    `mapstructure:"write_buffer_size"`
	CheckOrigin     bool   `mapstructure:"check_origin"`
}

// Validate 验证并设置所有配置默认值
func (c *Config) Validate() {
	c.Env = strings.ToLower(c.Env)
	if c.Env == "" || (c.Env != DevEnv && c.Env != ProdEnv) {
		c.Env = DefaultEnv
	}
	if c.Server != nil {
		c.Server.Validate()
	}

	if c.Log == nil {
		c.Log = &LogConfig{}
	}
	c.Log.Validate()

	// Validate other configs
	c.validateDatabaseConfig()
	c.validateRedisConfig()
	if c.I18n == nil {
		c.I18n = &I18nConfig{}
	}
	c.I18n.Validate()
	if c.Email != nil {
		c.Email.Validate()
	}
	if c.SMS != nil {
		c.SMS.Validate()
	}
	if c.VerifyCode != nil {
		c.VerifyCode.Validate()
	}
	if c.Dispatch != nil {
		c.Dispatch.Validate()
	}
	if c.Task != nil {
		c.Task.Validate()
	}
	if c.Google != nil {
		c.Google.Validate()
	}
	if c.Twilio != nil {
		c.Twilio.Validate()
	}
	// 初始化优惠券配置
	if c.Promotion != nil {
		c.Promotion.Validate()
	}
	if c.Payment != nil {
		c.Payment.Validate()
	}
	if c.KPay != nil {
		c.KPay.Validate()
	}
	if c.MoMo != nil {
		c.MoMo.Validate()
	}
	if c.Stripe != nil {
		c.Stripe.Validate()
	}
	if c.Order == nil {
		c.Order = &OrderConfig{}
	}
	if err := c.Order.Validate(); err != nil {
		fmt.Printf("Order config validation error: %v\n", err)
	}
}

func (c *Config) validateDatabaseConfig() {
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 5
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.ConnMaxLifetime == 0 {
		c.Database.ConnMaxLifetime = 300 * time.Second
	}
}

func (c *Config) validateRedisConfig() {
	// Redis配置验证
	// DSN格式验证可以在连接时进行
}

type I18nConfig struct {
	LocalesDir      string   `mapstructure:"locales_dir"`      // locales文件目录
	DefaultLanguage string   `mapstructure:"default_language"` // 默认语言
	SupportedLangs  []string `mapstructure:"supported_langs"`  // 支持的语言列表
}

// Validate 验证并设置I18n配置默认值
func (i *I18nConfig) Validate() {
	// 设置默认路径为容器中的绝对路径
	if i.LocalesDir == "" {
		i.LocalesDir = "/app/internal/locales" // 容器中的绝对路径
	}

	if i.DefaultLanguage == "" {
		i.DefaultLanguage = "en"
	}

	if len(i.SupportedLangs) == 0 {
		i.SupportedLangs = []string{"en", "zh"}
	}
}

// generateJWTSecret 生成安全的JWT密钥
func generateJWTSecret() (string, error) {
	// 生成64字节（512位）的随机密钥
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// LoadConfig 加载配置
func LoadConfig() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 支持从环境变量读取
	viper.SetEnvPrefix("GREENRIDE")
	viper.AutomaticEnv()
	// 支持嵌套键映射，例如 DATABASE_DSN 映射到 database.dsn
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 尝试加载 .env 文件 (可选)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		// 如果文件不存在，忽略错误
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error loading .env file: %v\n", err)
		}
	} else {
		fmt.Println(".env file loaded successfully")
	}

	// 重新设置回 yaml 配置文件名并加载
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err = viper.MergeInConfig(); err != nil {
		return
	}

	// 根据环境加载具体的配置文件 (prod.yaml, dev.yaml 等)
	env := viper.GetString("env")
	if env != "" {
		viper.SetConfigName(env)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Printf("Warning: environment specific config file '%s.yaml' not found, using default\n", env)
		} else {
			fmt.Printf("Loaded environment specific config: %s.yaml\n", env)
		}
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	fmt.Println("Configuration loaded successfully")
	// 验证并设置默认值
	config.Validate()
	return
}

// Get 获取配置单例
func Get() *Config {
	return config
}

// Set 设置配置单例
func Set(cfg *Config) {
	config = cfg
}

// AWS 配置相关方法
func (c *AWSConfig) IsConfigured() bool {
	// 由于我们使用IAM角色，只需要检查S3相关配置
	return c.S3Bucket != "" && c.S3Region != ""
}

// GetS3Bucket 获取 S3 桶名
func (c *AWSConfig) GetS3Bucket() string {
	return c.S3Bucket
}

// GetS3Region 获取 S3 区域
func (c *AWSConfig) GetS3Region() string {
	return c.S3Region
}

// ValidateAWS 验证 AWS 配置
func (c *AWSConfig) ValidateAWS() error {
	if c.S3Region == "" {
		return fmt.Errorf("AWS S3 Region is required")
	}
	if c.S3Bucket == "" {
		return fmt.Errorf("AWS S3 Bucket is required")
	}
	return nil
}

type InnoPaaSConfig struct {
	Endpoint      string `mapstructure:"endpoint" yaml:"endpoint" json:"endpoint"`
	AppKey        string `mapstructure:"app_key" yaml:"app_key" json:"app_key"`
	AppSecret     string `mapstructure:"app_secret" yaml:"app_secret" json:"app_secret"`           // For OTP v3 x-signature (MD5)
	Authorization string `mapstructure:"authorization" yaml:"authorization" json:"authorization"`   // Legacy token auth (if no app_secret)
	SenderID      string `mapstructure:"sender_id" yaml:"sender_id" json:"sender_id"`               // Optional sender (e.g. WhatsApp)
}

// Validate 验证InnoPaaS配置
func (c *InnoPaaSConfig) Validate() {
	if c.Endpoint == "" {
		c.Endpoint = "https://api.innopaas.com/api/otp/v3/msg/send/verify"
	}
}

// PromotionConfig 优惠券配置
type PromotionConfig struct {
	EnableWelcomeCoupon    string `mapstructure:"enable_welcome_coupon"`
	WelcomeCouponCode      string `mapstructure:"welcome_coupon_code"`
	WelcomeCouponValidDays int    `mapstructure:"welcome_coupon_valid_days"`

	// 推荐优惠券配置
	EnableReferralCoupon    string `mapstructure:"enable_referral_coupon"`
	ReferralCouponCode      string `mapstructure:"referral_coupon_code"`
	ReferralCouponValidDays int    `mapstructure:"referral_coupon_valid_days"`
}

// Validate 验证并设置优惠券配置默认值
func (p *PromotionConfig) Validate() {
	// 默认开启欢迎优惠券
	if p.EnableWelcomeCoupon == "" {
		p.EnableWelcomeCoupon = StatusOn
	}
	if p.WelcomeCouponCode == "" {
		p.WelcomeCouponCode = "WELCOME_NEW_USER"
	}
	if p.WelcomeCouponValidDays <= 0 {
		p.WelcomeCouponValidDays = 30
	}

	// 默认开启推荐优惠券
	if p.EnableReferralCoupon == "" {
		p.EnableReferralCoupon = StatusOn
	}
	if p.ReferralCouponCode == "" {
		p.ReferralCouponCode = "REFERRAL_REWARD"
	}
	if p.ReferralCouponValidDays <= 0 {
		p.ReferralCouponValidDays = 30
	}
}

// GetPromotionConfig 获取优惠券配置（带默认值）
func GetPromotionConfig() PromotionConfig {
	cfg := Get().Promotion
	if cfg == nil {
		cfg = &PromotionConfig{}
	}

	// 设置默认值
	result := *cfg

	// 默认开启优惠券功能
	// 欢迎优惠券默认配置
	if result.EnableWelcomeCoupon == "" {
		result.EnableWelcomeCoupon = StatusOn
	}
	if result.WelcomeCouponCode == "" {
		result.WelcomeCouponCode = "WELCOME_NEW_USER"
	}
	if result.WelcomeCouponValidDays <= 0 {
		result.WelcomeCouponValidDays = 30
	}

	// 推荐优惠券默认配置
	if result.EnableReferralCoupon == "" {
		result.EnableReferralCoupon = StatusOn
	}
	if result.ReferralCouponCode == "" {
		result.ReferralCouponCode = "INVITE_NEW_USER"
	}
	if result.ReferralCouponValidDays <= 0 {
		result.ReferralCouponValidDays = 30
	}

	return result
}
