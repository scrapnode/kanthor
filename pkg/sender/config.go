package sender

type Config struct {
	Timeout int   `json:"timeout" mapstructure:"timeout"`
	Retry   Retry `json:"retry" mapstructure:"retry"`
}

type Retry struct {
	Count       int `json:"count" mapstructure:"count"`
	WaitTime    int `json:"wait_time" mapstructure:"wait_time"`
	WaitTimeMax int `json:"wait_time_max" mapstructure:"wait_time_max"`
}
