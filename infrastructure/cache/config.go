package cache

type Config struct {
	Uri        string `json:"uri" mapstructure:"uri"`
	TimeToLive int    `json:"time_to_live" mapstructure:"time_to_live"`
}
