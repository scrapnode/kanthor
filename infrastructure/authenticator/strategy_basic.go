package authenticator

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
)

var SchemeBasic = "basic"

func BasicStrategy(verify func(ctx context.Context, username, password string) (*Account, error)) (Authenticate, error) {
	return func(ctx context.Context, request *Request) (*Account, error) {
		bytes, err := base64.StdEncoding.DecodeString(request.Credentials)
		if err != nil {
			return nil, err
		}

		up := strings.Split(string(bytes), ":")
		if len(up) != 2 {
			return nil, errors.New("AUTHENTICATOR.STRATEGY.BASIC.MALFORMED_CREDENTIALS")
		}
		return verify(ctx, up[0], up[1])
	}, nil
}
