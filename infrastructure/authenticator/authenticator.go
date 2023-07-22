package authenticator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(conf *Config, logger logging.Logger) (Authenticator, error) {
	if conf.Engine == EngineAsk {
		return NewAsk(conf, logger), nil
	}
	if conf.Engine == EngineCipher {
		return NewCipher(conf, logger), nil
	}

	return nil, fmt.Errorf("authenticator: unknow engine")
}

type Authenticator interface {
	Scheme() string
	Verify(token string) (*Account, error)
}
