package authenticator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/logging"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

	return verifier.parse(res.Body)
}

func (verifier *external) parse(response []byte) (*Account, error) {
	acc := &Account{}

	// use mapper if we configured it
	if len(verifier.conf.Mapper) > 0 {
		var strresp = string(response)
		var snippet = "{}"
		var err error

		for target, src := range verifier.conf.Mapper {
			value := gjson.Get(strresp, src)
			snippet, err = sjson.Set(snippet, target, value.String())
			if err != nil {
				return nil, err
			}
		}

		if err := json.Unmarshal([]byte(snippet), acc); err != nil {
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
