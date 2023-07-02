package circuitbreaker

type Config struct {
	Timeout               int `json:"timeout" mapstructure:"timeout"`
	SleepWindow           int `json:"sleep_window" mapstructure:"sleep_window"`
	ErrorPercentThreshold int `json:"error_percent_threshold" mapstructure:"error_percent_threshold"`
}
