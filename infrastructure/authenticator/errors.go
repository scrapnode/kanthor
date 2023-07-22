package authenticator

import "errors"

var (
	ErrMalformedToken        = errors.New("AUTHENTICATOR.MALFORMED_TOKEN")
	ErrMalformedTokenPayload = errors.New("AUTHENTICATOR.MALFORMED_TOKEN_PAYLOAD")
	ErrInvalidCredentials    = errors.New("AUTHENTICATOR.INVALID_CREDENTIALS")
)
