package authenticator

import "errors"

var (
	ErrMalformedToken     = errors.New("AUTHENTICATOR.MALFORMED_TOKEN")
	ErrMalformedPayload   = errors.New("AUTHENTICATOR.MALFORMED_PAYLOAD")
	ErrInvalidCredentials = errors.New("AUTHENTICATOR.INVALID_CREDENTIALS")
)
