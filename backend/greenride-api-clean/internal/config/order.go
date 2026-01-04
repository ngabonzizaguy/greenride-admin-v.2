package config

const (
	DefaultOrderExpireMinutes     = 60
	DefaultRideOrderExpireMinutes = 60
)

type OrderConfig struct {
	ExpireMinutes int              `mapstructure:"expire_minutes"`
	RideOrder     *RideOrderConfig `mapstructure:"ride_order"`
}

type RideOrderConfig struct {
	ExpireMinutes int `mapstructure:"expire_minutes"`
}

func (c *OrderConfig) Validate() error {
	if c.ExpireMinutes <= 0 {
		c.ExpireMinutes = DefaultOrderExpireMinutes
	}
	if c.RideOrder == nil {
		c.RideOrder = &RideOrderConfig{}
	}
	if err := c.RideOrder.Validate(); err != nil {
		return err
	}
	return nil
}
func (c *RideOrderConfig) Validate() error {
	if c.ExpireMinutes <= 0 {
		c.ExpireMinutes = DefaultRideOrderExpireMinutes
	}
	return nil
}
