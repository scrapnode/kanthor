package authenticator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(conf *Config, logger logging.Logger) (Authenticator, error) {
	if conf.Engine == EngineAsk {
		return NewAsk(conf, logger)
	}

	return nil, fmt.Errorf("authenticator: unknow engine")
}

type Authenticator interface {
	Verify(credentials string) (*Account, error)
}
