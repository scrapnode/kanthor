package authenticator

import (
	"encoding/base64"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"strings"
)

func NewASK(conf *Config, logger logging.Logger) Authenticator {
	return &ask{conf: conf, logger: logger}
}

// short of Access Secret Key
type ask struct {
	conf   *Config
	logger logging.Logger
}

func (authenticator *ask) Scheme() string {
	return "basic"
}

func (authenticator *ask) Verify(token string) (*Account, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, ErrMalformedToken
	}

	as := strings.Split(string(bytes), ":")
	if len(as) != 2 {
		return nil, ErrMalformedPayload
	}

	accessOK := as[0] == authenticator.conf.AccessSecretKey.AccessKey
	secretOk := as[1] == authenticator.conf.AccessSecretKey.SecretKey
	if !accessOK || !secretOk {
		return nil, ErrInvalidCredentials
	}

	account := &Account{
		Sub:  authenticator.conf.AccessSecretKey.AccessKey,
		Iss:  "kanthor.system",
		Aud:  "kanthor",
		Name: "Kanthor",
	}
	return account, nil
}
