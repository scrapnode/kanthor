package auth

import "github.com/scrapnode/kanthor/infrastructure/logging"

func New(conf *Config, logger logging.Logger) Auth {
	return NewInternal(conf, logger)
}

type Auth interface {
	User() *User
	Tier() string
}

type User struct {
	Id string
}

type Config struct {
}
