package authenticator

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var SchemeBearer = "bearer"

func BearerStrategy(conf *Bearer) (Authenticate, error) {
	return func(ctx context.Context, request *Request) (*Account, error) {
		token, err := jwt.ParseWithClaims(
			request.Credentials,
			&JwtClaims{},
			func(token *jwt.Token) (interface{}, error) {
				// IMPORANT: force use HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("AUTHENTICATOR.STRATEGY.BEARER.MALFORMED_ALGO: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(conf.Secret), nil
			},
			jwt.WithIssuedAt(),
		)
		if err != nil {
			return nil, err
		}

		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			return nil, fmt.Errorf("AUTHENTICATOR.STRATEGY.BEARER.MALFORMED_CLAIMS")
		}

		account := &Account{
			Sub:     claims.Subject,
			Name:    claims.Name,
			Picture: claims.Picture,
		}
		return account, nil
	}, nil
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Name    string `json:"name" yaml:"name"`
	Picture string `json:"picture" yaml:"picture"`
}
