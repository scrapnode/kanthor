package auth

import "github.com/scrapnode/kanthor/infrastructure/logging"

func NewInternal(conf *Config, logger logging.Logger) Auth {
	return &internal{}
}

type internal struct {
}

func (auth *internal) User() *User {
	return &User{Id: "user_0abcdefghijklmnopqrstuvwxyz"}
}

func (auth *internal) Tier() string {
	return "default"
}
