package config

const (
	// 验证码默认配置
	DefaultVerifyCodeLength     = 6  // 默认验证码长度
	DefaultVerifyCodeExpiration = 5  // 默认验证码有效期(分钟)
	DefaultVerifyCodeInterval   = 60 // 默认发送间隔(秒)
	DefaultVerifyCodeMaxTimes   = 10 // 默认每天最大发送次数
)

// VerifyCodeConfig 验证码配置
type VerifyCodeConfig struct {
	Length        int  `mapstructure:"length" yaml:"length" json:"length"`                         // 验证码长度
	Expiration    int  `mapstructure:"expiration" yaml:"expiration" json:"expiration"`             // 验证码有效期(分钟)
	SendInterval  int  `mapstructure:"send_interval" yaml:"send_interval" json:"send_interval"`    // 发送间隔(秒)
	MaxSendTimes  int  `mapstructure:"max_send_times" yaml:"max_send_times" json:"max_send_times"` // 每天最大发送次数
	LocalTemplate bool `mapstructure:"local_template" yaml:"local_template" json:"local_template"` // 是否使用本地模板
	BypassOTP     bool `mapstructure:"bypass_otp" yaml:"bypass_otp" json:"bypass_otp"`             // 是否绕过OTP验证
}

// Validate 验证验证码配置
func (c *VerifyCodeConfig) Validate() {
	if c.Length <= 0 {
		c.Length = DefaultVerifyCodeLength
	}
	if c.Expiration <= 0 {
		c.Expiration = DefaultVerifyCodeExpiration
	}
	if c.SendInterval <= 0 {
		c.SendInterval = DefaultVerifyCodeInterval
	}
	if c.MaxSendTimes <= 0 {
		c.MaxSendTimes = DefaultVerifyCodeMaxTimes
	}
}
