package authenticator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

var (
	HeaderAuth  = "authorization"
	SchemeBasic = "basic"
)

func New(conf *Config, logger logging.Logger) (Authenticator, error) {
	if conf.Engine == EngineAsk {
		return NewAsk(conf, logger)
	}

	return nil, fmt.Errorf("authenticator: unknown engine")
}

type Authenticator interface {
	Verify(credentials string) (*Account, error)
}
