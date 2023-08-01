package authenticator

import (
	"encoding/base64"
	"strings"
)

func ParseBasicCredentials(credentials string) (string, string, error) {
	segments := strings.Split(credentials, " ")
	if len(segments) != 2 {
		return "", "", ErrMalformedCredentials
	}
	if !strings.EqualFold(segments[0], "basic") {
		return "", "", ErrMalformedScheme
	}

	if segments[1] == "" {
		return "", "", ErrMalformedToken
	}
	bytes, err := base64.StdEncoding.DecodeString(segments[1])
	if err != nil {
		return "", "", err
	}

	as := strings.Split(string(bytes), ":")
	if len(as) != 2 {
		return "", "", ErrMalformedTokenPayload
	}

	return as[0], as[1], nil
}
