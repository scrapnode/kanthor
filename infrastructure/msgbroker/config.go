package msgbroker

import (
	"github.com/scrapnode/kanthor/infrastructure/config"
)

type Config struct {
	Uri    string `json:"uri"`
	Stream struct {
		Name     string `json:"name"`
		Replicas int    `json:"replicas"`
		Subject  string `json:"subject"`
		Limits   struct {
			Msgs     int64 `json:"msgs"`
			MsgBytes int32 `json:"msg_bytes"`
			Bytes    int64 `json:"bytes"`
			Age      int64 `json:"age"`
		} `json:"limits"`
	} `json:"stream"`
	Consumer ConfigConsumer `json:"consumer"`
}

type ConfigConsumer struct {
	Name      string `json:"name"`
	Temporary bool   `json:"temporary"`
	MaxRetry  int    `json:"max_retry"`
}

const ConfigName = "msgbroker"

func GetConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.UnmarshalKey(ConfigName, &cfg)
	return &cfg, err
}
