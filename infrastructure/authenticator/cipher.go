package authenticator

import (
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewCipher(conf *Config, logger logging.Logger) (Authenticator, error) {
	logger = logger.With("authenticator", "cipher")
	symmetric, err := cryptography.NewSymmetric(conf.Cipher)
	if err != nil {
		return nil, err
	}
	return &cipher{conf: conf, logger: logger, symmetric: symmetric}, nil
}

// short of Access Secret Key
type cipher struct {
	conf   *Config
	logger logging.Logger

	symmetric cryptography.Symmetric
}

func (authenticator *cipher) Scheme() string {
	return "basic"
}

func (authenticator *cipher) Verify(token string) (*Account, error) {
	sub, err := authenticator.symmetric.StringDecrypt(token)
	if err != nil {
		authenticator.logger.Error(err)
		return nil, ErrMalformedToken
	}

	account := &Account{
		Sub:  sub,
		Iss:  "kanthor.authenticator.cipher",
		Aud:  "kanthor",
		Name: "Kanthor Cipher",
	}
	return account, nil
}
