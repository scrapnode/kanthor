package healthcheck

type Config struct {
	Dest    string `json:"dest" yaml:"dest" mapstructure:"dest" validate:"required"`
	Timeout int    `json:"timeout" yaml:"timeout" mapstructure:"timeout" validate:"required,number,gte=0"`
	MaxTry  int    `json:"max_try" yaml:"max_try" mapstructure:"max_try" validate:"required,number,gte=0"`
}

func DefaultConfig(dest string) *Config {
	return &Config{Dest: dest, Timeout: 3000, MaxTry: 1}
}

type Server interface {
	Readiness(check func() error) error
	Liveness(check func() error) error
}

type Client interface {
	Readiness() error
	Liveness() error
}
