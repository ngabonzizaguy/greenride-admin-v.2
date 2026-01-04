package config

// SMS配置结构
type SMSConfig struct {
	ServiceName    string `mapstructure:"service_name"`    // 使用的短信服务提供商名称
	DefaultAccount string `mapstructure:"default_account"` // 默认账号ID
	DefaultNumber  string `mapstructure:"default_number"`  // 默认发送号码
}

// Validate 验证并设置SMS配置默认值
func (s *SMSConfig) Validate() {
	if s.ServiceName == "" {
		s.ServiceName = "twilio" // 默认使用Twilio服务
	}
}
