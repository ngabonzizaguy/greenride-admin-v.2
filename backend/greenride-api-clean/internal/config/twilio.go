package config

type TwilioConfig struct {
	Accounts []TwilioAccountConfig `mapstructure:"accounts"`
}

// TwilioAccountConfig 代表单个Twilio账号的配置
type TwilioAccountConfig struct {
	AccountSID string   `mapstructure:"account_sid"`
	AuthToken  string   `mapstructure:"auth_token"`
	Phones     []string `mapstructure:"phones"`      // 关联的电话号码列表
	ServiceSID string   `mapstructure:"service_sid"` // Verify Service SID (可选)
}

func (c *TwilioConfig) Validate() {
	// 确保账号列表不为空
	if len(c.Accounts) > 0 {
		for _, item := range c.Accounts {
			item.Validate()
		}
	}
}

func (c *TwilioAccountConfig) Validate() {
	// Let the service layer handle missing credentials instead of crashing the whole app on load
}
