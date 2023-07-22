package authenticator

import (
	"github.com/scrapnode/kanthor/infrastructure/crypto"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewCipher(conf *Config, logger logging.Logger) Authenticator {
	logger = logger.With("authenticator", "cipher")
	return &cipher{conf: conf, logger: logger, ase: crypto.NewAES(conf.Cipher.Key)}
}

// short of Access Secret Key
type cipher struct {
	conf   *Config
	logger logging.Logger

	ase *crypto.AES
}

func (authenticator *cipher) Scheme() string {
	return "basic"
}

func (authenticator *cipher) Verify(token string) (*Account, error) {
	sub, err := authenticator.ase.DecryptString(token)
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
