package authenticator

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewAsk(conf *Config, logger logging.Logger) (Authenticator, error) {
	logger = logger.With("authenticator", "ask")
	return &ask{conf: conf, logger: logger}, nil
}

// short of Access Secret Key
type ask struct {
	conf   *Config
	logger logging.Logger
}

func (authenticator *ask) Verify(credentials string) (*Account, error) {
	ak, sk, err := ParseBasicCredentials(credentials)
	if err != nil {
		authenticator.logger.Error(err)
		return nil, ErrMalformedToken
	}

	accessOK := ak == authenticator.conf.Ask.AccessKey
	secretOk := sk == authenticator.conf.Ask.SecretKey
	if !accessOK || !secretOk {
		return nil, ErrInvalidCredentials
	}

	account := &Account{
		Sub:  ak,
		Iss:  "kanthor.authenticator.ask",
		Aud:  "kanthor",
		Name: "Kanthor Ask",
	}
	return account, nil
}
