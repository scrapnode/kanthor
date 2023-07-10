package authenticator

import "github.com/scrapnode/kanthor/infrastructure/logging"

type Authenticator interface {
	Verify(token string) (*Account, error)
}

func New(conf *Config, logger logging.Logger) Authenticator {
	return NewASK(conf, logger)
}
