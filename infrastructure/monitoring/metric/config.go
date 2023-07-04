package metric

type Config struct {
	Enable bool `json:"enable" mapstructure:"enable"`
	Server struct {
		Addr string `json:"addr" mapstructure:"addr"`
	} `json:"server" mapstructure:"server"`
}
