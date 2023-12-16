package authenticator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/logging"
)

var EngineForward = "forward"

func NewForward(conf *Forward, logger logging.Logger, send sender.Send, cb circuitbreaker.CircuitBreaker) Verifier {
	return &forward{conf: conf, send: send, cb: cb}
}

type forward struct {
	conf *Forward
	send sender.Send
	cb   circuitbreaker.CircuitBreaker
}

func (verifier *forward) Verify(ctx context.Context, request *Request) (*Account, error) {
	req := &sender.Request{
		Method:  http.MethodGet,
		Headers: http.Header{},
		Uri:     verifier.conf.Uri,
	}

	// add authorization header
	req.Headers.Add(HeaderAuthCredentials, request.Credentials)
	// then add others headers as well
	for _, key := range verifier.conf.Headers {
		if value, has := request.Metadata[key]; has {
			req.Headers.Add(key, value)
		}
	}

	res, err := circuitbreaker.Do[sender.Response](
		verifier.cb,
		verifier.conf.Uri,
		func() (interface{}, error) {
			return verifier.send(ctx, req)
		},
		func(err error) error {
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	var acc Account
	if err := json.Unmarshal(res.Body, &acc); err != nil {
		return nil, err
	}

	if err := acc.Validate(); err != nil {
		return nil, err
	}

	return &acc, nil
}
