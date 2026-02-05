package config

import (
	"path"
	"strings"
)

const (
	// Default callback URL path for Stripe webhooks
	DefaultStripeCallbackURL = "/webhook/stripe"
)

// StripeGlobalConfig holds global Stripe configuration from yaml files
type StripeGlobalConfig struct {
	CallbackUrl string `mapstructure:"callback_url" json:"callback_url"`
	Timeout     int    `mapstructure:"timeout" json:"timeout"`
}

// Validate validates and sets default values for Stripe global config
func (c *StripeGlobalConfig) Validate() error {
	if c.Timeout <= 0 {
		c.Timeout = DefaultPaymentTimeout
	}
	cfg := Get().Payment
	if c.CallbackUrl == "" {
		c.CallbackUrl = DefaultStripeCallbackURL
	}
	if !strings.HasPrefix(c.CallbackUrl, "http") {
		c.CallbackUrl = path.Join(cfg.CallbackHost, c.CallbackUrl)
	}
	return nil
}
