package authenticator

import "errors"

var (
	ErrMalformedCredentials  = errors.New("AUTHENTICATOR.MALFORMED_CREDENTIALS")
	ErrMalformedScheme       = errors.New("AUTHENTICATOR.MALFORMED_SCHEME")
	ErrMalformedToken        = errors.New("AUTHENTICATOR.MALFORMED_TOKEN")
	ErrMalformedTokenPayload = errors.New("AUTHENTICATOR.MALFORMED_TOKEN_PAYLOAD")
	ErrInvalidCredentials    = errors.New("AUTHENTICATOR.INVALID_CREDENTIALS")
)
