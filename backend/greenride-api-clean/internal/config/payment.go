package config

const (
	DefaultPaymentTimeout      = 30 * 60 // 默认支付请求超时时间，单位秒
	DefaultRequestTimeout      = 30
	DefaultPaymentCallbackHost = "https://api.greenrideafrica.com"
	DefaultPaymentReturnURL    = "https://www.greenrideafrica.com/payment_result"
)

type PaymentConfig struct {
	Sandbox        int    `mapstructure:"sandbox" json:"sandbox"`
	LogoURL        string `mapstructure:"logo_url" json:"logo_url"` // 结账页面Logo URL
	ReturnURL      string `mapstructure:"return_url" json:"return_url"`
	CallbackHost   string `mapstructure:"callback_host" json:"callback_host"`
	RequestTimeout int    `mapstructure:"request_timeout" json:"request_timeout"` // 请求超时时间，单位秒
	PaymentTimeout int    `mapstructure:"payment_timeout" json:"payment_timeout"` // 请求超时时间，单位秒
}

func (c *PaymentConfig) IsSandbox() bool {
	return c.Sandbox == 1
}

func (c *PaymentConfig) Validate() error {
	if c.RequestTimeout <= 0 {
		c.RequestTimeout = DefaultPaymentTimeout
	}
	if c.PaymentTimeout <= 0 {
		c.PaymentTimeout = DefaultPaymentTimeout
	}
	if c.ReturnURL == "" {
		c.ReturnURL = DefaultPaymentReturnURL
	}
	if c.CallbackHost == "" {
		c.CallbackHost = DefaultPaymentCallbackHost
	}
	if c.Sandbox != 1 {
		c.Sandbox = 0
	}

	return nil
}
