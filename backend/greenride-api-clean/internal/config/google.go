package config

type GoogleConfig struct {
	MapsAPIKey       string `mapstructure:"maps_api_key"`
	CloudAPIKey      string `mapstructure:"cloud_api_key"`
	ServiceAccountID string `mapstructure:"service_account_id"`
	ProjectID        string `mapstructure:"project_id"`
	// Maps API specific settings
	RequestTimeout       int `mapstructure:"request_timeout"`         // API请求超时时间(秒)
	MaxRetries           int `mapstructure:"max_retries"`             // 最大重试次数
	RateLimitWindow      int `mapstructure:"rate_limit_window"`       // 限流窗口(秒)
	MaxRequestsPerWindow int `mapstructure:"max_requests_per_window"` // 窗口内最大请求数
}

// Validate 验证并设置Google配置默认值
func (g *GoogleConfig) Validate() {
	if g.RequestTimeout <= 0 {
		g.RequestTimeout = 30 // 默认30秒超时
	}
	if g.MaxRetries <= 0 {
		g.MaxRetries = 3 // 默认重试3次
	}
	if g.RateLimitWindow <= 0 {
		g.RateLimitWindow = 60 // 默认60秒窗口
	}
	if g.MaxRequestsPerWindow <= 0 {
		g.MaxRequestsPerWindow = 50 // 默认每分钟50个请求
	}
}
