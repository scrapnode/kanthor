package config

import (
	"github.com/scrapnode/kanthor/infrastructure/auth"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Config struct {
	Database  database.Config            `json:"database"`
	Datastore datastore.Config           `json:"datastore"`
	Streaming streaming.ConnectionConfig `json:"streaming_publisher"`

	Dataplane *Dataplane `json:"dataplane"`
}

type Dataplane struct {
	Logger logging.Config `json:"logger"`
	Auth   auth.Config    `json:"auth"`
	Server struct {
		Addr string `json:"addr"`
	} `json:"server"`
	Message struct {
		BucketLayout string `json:"bucket_layout"`
	} `json:"message"`
}

func New(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Unmarshal(&cfg)
	return &cfg, err
}
