package cache

type Config struct {
	Uri                 string `json:"uri" mapstructure:"uri"`
	TimeToLiveInSeconds int    `json:"time_to_live_in_seconds" mapstructure:"time_to_live_in_seconds"`
}
