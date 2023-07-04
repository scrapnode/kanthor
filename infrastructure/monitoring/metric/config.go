package metric

type Config struct {
	Server struct {
		Addr string `json:"addr" mapstructure:"addr"`
	} `json:"server" mapstructure:"server"`
}
