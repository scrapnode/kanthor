package authenticator

import (
	"encoding/base64"
	"strings"
)

func ParseBasicCredentials(credentials string) (string, string, error) {
	bytes, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return "", "", err
	}

	as := strings.Split(string(bytes), ":")
	if len(as) != 2 {
		return "", "", ErrMalformedTokenPayload
	}

	return as[0], as[1], nil
}