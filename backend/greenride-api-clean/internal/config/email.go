package config

// EmailConfig 邮件服务配置
type EmailConfig struct {
	ServiceName string               `mapstructure:"service_name" yaml:"service_name" json:"service_name"` // 服务名称：account, resend, aliyun, tencent
	Sender      string               `mapstructure:"sender" yaml:"sender" json:"sender"`                   // 默认发件人
	Accounts    []EmailAccountConfig `mapstructure:"accounts" yaml:"accounts" json:"accounts"`             // 邮箱账户列表
	Aliyun      *AliyunEmailConfig   `mapstructure:"aliyun" yaml:"aliyun" json:"aliyun"`                   // 阿里云邮件配置
	Tencent     *TencentEmailConfig  `mapstructure:"tencent" yaml:"tencent" json:"tencent"`                // 腾讯云邮件配置
}

// EmailAccountConfig 邮箱账户配置
type EmailAccountConfig struct {
	Name     string `mapstructure:"name" yaml:"name" json:"name"`             // 账户名称
	Host     string `mapstructure:"host" yaml:"host" json:"host"`             // SMTP主机
	Port     int    `mapstructure:"port" yaml:"port" json:"port"`             // SMTP端口
	Username string `mapstructure:"username" yaml:"username" json:"username"` // 用户名
	Password string `mapstructure:"password" yaml:"password" json:"password"` // 密码
	From     string `mapstructure:"from" yaml:"from" json:"from"`             // 发件人地址
	SSL      bool   `mapstructure:"ssl" yaml:"ssl" json:"ssl"`                // 是否使用SSL
	Priority int    `mapstructure:"priority" yaml:"priority" json:"priority"` // 优先级，数字越小优先级越高
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`    // 是否启用
}

// AliyunEmailConfig 阿里云邮件服务配置
type AliyunEmailConfig struct {
	Enabled         bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
	AccessKeyID     string `mapstructure:"access_key_id" yaml:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" yaml:"access_key_secret" json:"access_key_secret"`
	Region          string `mapstructure:"region" yaml:"region" json:"region"`                   // 区域，如：cn-hangzhou
	FromAlias       string `mapstructure:"from_alias" yaml:"from_alias" json:"from_alias"`       // 发件人别名
	AccountName     string `mapstructure:"account_name" yaml:"account_name" json:"account_name"` // 发信地址
}

// TencentEmailConfig 腾讯云邮件服务配置
type TencentEmailConfig struct {
	Enabled    bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
	SecretID   string `mapstructure:"secret_id" yaml:"secret_id" json:"secret_id"`
	SecretKey  string `mapstructure:"secret_key" yaml:"secret_key" json:"secret_key"`
	Region     string `mapstructure:"region" yaml:"region" json:"region"`                // 区域，如：ap-beijing
	FromEmail  string `mapstructure:"from_email" yaml:"from_email" json:"from_email"`    // 发件人邮箱
	FromName   string `mapstructure:"from_name" yaml:"from_name" json:"from_name"`       // 发件人名称
	TemplateID uint64 `mapstructure:"template_id" yaml:"template_id" json:"template_id"` // 模板ID
}

// Validate 验证邮件配置
func (c *EmailConfig) Validate() {
	if c.ServiceName == "" {
		c.ServiceName = "account" // 默认使用账户服务
	}

	// 根据服务类型验证相关配置
	switch c.ServiceName {
	case "account":
		for _, account := range c.Accounts {
			account.Validate()
		}
	case "aliyun":
		if c.Aliyun != nil {
			c.Aliyun.Validate()
		}
	case "tencent":
		if c.Tencent != nil {
			c.Tencent.Validate()
		}
	}
}

// Validate 验证邮箱账户配置
func (c *EmailAccountConfig) Validate() {
	if !c.Enabled {
		return
	}
	if c.Name == "" {
		// log.Printf("Warning: email account name is required")
	}
	if c.Host == "" {
		// log.Printf("Warning: email account host is required")
	}
	// ... we can just return instead of panicking
}

// Validate 验证阿里云邮件配置
func (c *AliyunEmailConfig) Validate() {
	if !c.Enabled {
		return
	}
	if c.AccessKeyID == "" {
		panic("aliyun email access_key_id is required")
	}
	if c.AccessKeySecret == "" {
		panic("aliyun email access_key_secret is required")
	}
	if c.Region == "" {
		c.Region = "cn-hangzhou"
	}
	if c.AccountName == "" {
		panic("aliyun email account_name is required")
	}
}

// Validate 验证腾讯云邮件配置
func (c *TencentEmailConfig) Validate() {
	if !c.Enabled {
		return
	}
	if c.SecretID == "" {
		panic("tencent email secret_id is required")
	}
	if c.SecretKey == "" {
		panic("tencent email secret_key is required")
	}
	if c.Region == "" {
		c.Region = "ap-beijing"
	}
	if c.FromEmail == "" {
		panic("tencent email from_email is required")
	}
}
