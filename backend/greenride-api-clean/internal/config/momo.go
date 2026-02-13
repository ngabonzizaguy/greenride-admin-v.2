package config

import (
	"path"
	"strings"
)

const (
	// MoMo API Base URLs
	MoMoSandboxBaseURL    = "https://sandbox.momodeveloper.mtn.com"
	MoMoProductionBaseURL = "https://proxy.momoapi.mtn.co.rw"

	// Default callback URL path
	DefaultMoMoCallbackURL = "/webhook/momo"
)

// MoMoGlobalConfig holds global MoMo configuration from yaml files
type MoMoGlobalConfig struct {
	CallbackUrl string `mapstructure:"callback_url" json:"callback_url"`
	Timeout     int    `mapstructure:"timeout" json:"timeout"`
}

// Validate validates and sets default values for MoMo global config
func (c *MoMoGlobalConfig) Validate() error {
	if c.Timeout <= 0 {
		c.Timeout = DefaultPaymentTimeout
	}
	cfg := Get().Payment
	if c.CallbackUrl == "" {
		c.CallbackUrl = DefaultMoMoCallbackURL
	}
	if !strings.HasPrefix(c.CallbackUrl, "http") {
		c.CallbackUrl = path.Join(cfg.CallbackHost, c.CallbackUrl)
	}
	return nil
}

// GetBaseURL returns the appropriate base URL based on environment
func GetMoMoBaseURL(environment string) string {
	if strings.ToLower(environment) == "production" {
		return MoMoProductionBaseURL
	}
	return MoMoSandboxBaseURL
}
