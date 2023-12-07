package authenticator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/logging"
)

var SchemeForward = "forward"

func ForwardStrategy(conf *Forward) (Authenticate, error) {
	send, err := sender.New(sender.DefaultConfig, logging.NewNoop())
	if err != nil {
		return nil, err
	}

	cb, err := circuitbreaker.New(circuitbreaker.DefaultConfig, logging.NewNoop())
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, request *Request) (*Account, error) {
		return circuitbreaker.Do[Account](
			cb,
			conf.Endpoint,
			func() (interface{}, error) {
				in := &sender.Request{
					Method:  http.MethodGet,
					Headers: http.Header{},
					Uri:     conf.Endpoint,
				}
				// add authentication header first
				in.Headers.Add(HeaderAuth, request.Credentials)
				// then add more other allow header
				for _, key := range conf.RequestHeaders {
					if value, ok := request.Metadata[key]; ok {
						in.Headers.Add(key, value)
					}
				}

				res, err := send(context.Background(), in)
				if err != nil {
					return nil, err
				}

				var account Account
				if err := json.Unmarshal(res.Body, &account); err != nil {
					return nil, err
				}

				return &account, nil
			},
			func(err error) error {
				return err
			},
		)
	}, nil
}
