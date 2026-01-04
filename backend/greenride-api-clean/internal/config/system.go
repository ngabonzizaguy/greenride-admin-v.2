package config

import (
	"log"
	"net/http"
	"time"
)

const (
	DefaultReadTimeout  = 60
	DefaultWriteTimeout = 60
)

// ServerConfig 服务配置
type ServerConfig struct {
	Port  string         `mapstructure:"port"`  // 服务端口
	Api   *ServiceConfig `mapstructure:"api"`   // API服务配置
	Admin *ServiceConfig `mapstructure:"admin"` // 管理后台服务配置
}

type ServiceConfig struct {
	Prefix       string     `mapstructure:"prefix"`
	Name         string     `mapstructure:"name"`
	Port         string     `mapstructure:"port"`
	Version      string     `mapstructure:"version"`
	ReadTimeout  int        `mapstructure:"read_timeout"`  // 读取超时时间(秒)
	WriteTimeout int        `mapstructure:"write_timeout"` // 写入超时时间(秒)
	Jwt          *JWTConfig `mapstructure:"jwt"`           // JWT配置
}

func (s *ServiceConfig) ToServer() *http.Server {
	return &http.Server{
		Addr:         ":" + s.Port,
		ReadTimeout:  DefaultReadTimeout * time.Second,
		WriteTimeout: DefaultWriteTimeout * time.Second,
	}
}

func (s *ServiceConfig) Validate() {
	if s.Port == "" {
		panic("Service port cannot be empty")
	}
	if s.ReadTimeout <= 0 {
		s.ReadTimeout = DefaultReadTimeout
	}
	if s.WriteTimeout <= 0 {
		s.WriteTimeout = DefaultWriteTimeout
	}
	if s.Jwt == nil {
		s.Jwt = &JWTConfig{}
	}
	s.Jwt.Validate()
}

const (
	ProdEnv                 = "prod"
	DevEnv                  = "dev"
	DefaultEnv              = ProdEnv
	DefaultTaskScanInterval = 3  // 默认任务扫描间隔(秒)
	DefaultCacheDuration    = 10 // 默认缓存时间(分钟)
)

// Validate 验证并设置服务配置默认值
func (c *ServerConfig) Validate() {
	if c.Api != nil {
		c.Api.Validate()
	}
	if c.Admin != nil {
		c.Admin.Validate()
	}
}

// TaskConfig 任务配置
type TaskConfig struct {
	ScanInterval int  `mapstructure:"scan_interval"`          // 扫描间隔(秒)
	Enabled      bool `mapstructure:"enabled" default:"true"` // 定时任务开关，默认开启
}

// Validate 验证并设置任务配置默认值
func (c *TaskConfig) Validate() {
	if c.ScanInterval <= 0 {
		c.ScanInterval = DefaultTaskScanInterval
	}
}

// DBConfig 数据库配置
type DBConfig struct {
	Dsn             string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	Debug           bool          `mapstructure:"debug"`
}

type RedisConfig struct {
	Dsn string `mapstructure:"dsn"`
}

// LogConfig 日志配置
type LogConfig struct {
	Path     string                `mapstructure:"path"`     // 日志路径
	Level    string                `mapstructure:"level"`    // 日志级别
	Format   string                `mapstructure:"format"`   // 日志格式 (json, text)
	Output   string                `mapstructure:"output"`   // 输出方式 (stdout, file, both)
	Services map[string]*LogConfig `mapstructure:"services"` // 服务日志配置
	Loki     *LokiConfig           `mapstructure:"loki"`     // Grafana Loki 配置
}

// LokiConfig Grafana Loki 配置
type LokiConfig struct {
	Enabled   bool              `mapstructure:"enabled"`    // 是否启用 Loki Hook
	URL       string            `mapstructure:"url"`        // Loki push endpoint
	Username  string            `mapstructure:"username"`   // Grafana Cloud user ID
	Password  string            `mapstructure:"password"`   // Grafana Cloud API key
	Labels    map[string]string `mapstructure:"labels"`     // 默认标签
	BatchWait time.Duration     `mapstructure:"batch_wait"` // 批量等待时间
	BatchSize int               `mapstructure:"batch_size"` // 批量大小 (bytes)
	Hook      bool              `mapstructure:"hook"`       // hook是否启用
	Promtail  bool              `mapstructure:"promtail"`   // promtail是否启用
}

// Validate 验证并设置 Loki 配置默认值
func (l *LokiConfig) Validate() {
	if l.BatchWait == 0 {
		l.BatchWait = 1 * time.Second
	}
	if l.BatchSize == 0 {
		l.BatchSize = 1024 * 1024 // 1MB
	}
	if l.Labels == nil {
		l.Labels = make(map[string]string)
	}
}

const (
	JWTDefaultExpiration = "336h" // 2周
	JWTDetailIssuer      = "Greenride"
)

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration string        `mapstructure:"expiration"`
	Issuer     string        `mapstructure:"issuer"`
	Audience   string        `mapstructure:"audience"`
	ExpiresIn  time.Duration `mapstructure:"expires_in"`
}

func (c *JWTConfig) Validate() {
	// 如果JWT密钥为空，自动生成一个安全的随机密钥
	if c.Secret == "" {
		secret, err := generateJWTSecret()
		if err != nil {
			log.Printf("Failed to generate JWT secret, using fallback: %v", err)
			c.Secret = "greenride-fallback-jwt-secret-please-configure-properly"
		} else {
			c.Secret = secret
			log.Printf("Auto-generated JWT secret for security (length: %d)", len(secret))
		}
	}

	if c.Expiration == "" {
		c.Expiration = JWTDefaultExpiration
	}
	if c.Issuer == "" {
		c.Issuer = JWTDetailIssuer
	}

	// 解析Expiration字符串并设置ExpiresIn duration
	duration, err := time.ParseDuration(c.Expiration)
	if err != nil {
		log.Printf("Failed to parse JWT expiration '%s', using default: %v", c.Expiration, err)
		duration, _ = time.ParseDuration(JWTDefaultExpiration) // 使用默认值 "336h"
	}
	c.ExpiresIn = duration
}
