package authenticator

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/go-jsonnet"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/logging"
)

var EngineExternal = "external"

func NewExternal(conf *External, logger logging.Logger, send sender.Send, cb circuitbreaker.CircuitBreaker) (Verifier, error) {
	return &external{conf: conf, send: send, cb: cb}, nil
}

type external struct {
	conf *External
	send sender.Send
	cb   circuitbreaker.CircuitBreaker
}

func (verifier *external) Verify(ctx context.Context, request *Request) (*Account, error) {
	req := &sender.Request{
		Method:  http.MethodGet,
		Headers: http.Header{},
		Uri:     verifier.conf.Uri,
	}

	// add authorization header
	req.Headers.Add(HeaderAuthnCredentials, request.Credentials)
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

	return verifier.parse(res.Body)
}

func (verifier *external) parse(response []byte) (*Account, error) {
	acc := &Account{}

	// use mapper if we configured it
	if verifier.conf.Mapper != nil {
		vm := jsonnet.MakeVM()
		vm.ExtCode("session", string(response))

		// @TODO: parse this block once
		r, err := url.Parse(verifier.conf.Mapper.Uri)
		if err != nil {
			return nil, err
		}
		snippet, err := base64.StdEncoding.DecodeString(r.Host)
		if err != nil {
			return nil, err
		}

		j, err := vm.EvaluateAnonymousSnippet("default.jsonnet", string(snippet))
		if err != nil {
			return nil, err
		}

		// try to map the account directly
		if err := json.Unmarshal([]byte(j), acc); err != nil {
			return nil, err
		}
	} else {
		// try to map the account directly
		if err := json.Unmarshal(response, acc); err != nil {
			return nil, err
		}
	}

	if err := acc.Validate(); err != nil {
		return nil, err
	}
	return acc, nil
}
