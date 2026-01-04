package config

import (
	"path"
	"strings"
)

const (
	DefaultKPayCallbackURL = "/webhook/kpay"
)

type KPayConfig struct {
	LogoURL     string `mapstructure:"logo_url" json:"logo_url"` // 结账页面Logo URL
	CallbackUrl string `mapstructure:"callback_url" json:"callback_url"`
	ReturnURL   string `mapstructure:"return_url" json:"return_url"`
	Timeout     int    `mapstructure:"timeout" json:"timeout"` // 请求超时时间，单位秒
}

func (c *KPayConfig) Validate() error {
	if c.Timeout <= 0 {
		c.Timeout = DefaultPaymentTimeout
	}
	cfg := Get().Payment
	if c.LogoURL == "" {
		c.LogoURL = cfg.LogoURL
	}
	if c.CallbackUrl == "" {
		c.CallbackUrl = DefaultKPayCallbackURL
	}
	if !strings.HasPrefix(c.CallbackUrl, "http") {
		c.CallbackUrl = path.Join(cfg.CallbackHost, c.CallbackUrl)
	}
	if c.ReturnURL == "" {
		c.ReturnURL = cfg.ReturnURL
	}
	return nil
}
