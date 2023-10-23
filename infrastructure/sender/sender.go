package sender

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func New(conf *Config, logger logging.Logger) (Send, error) {
	rest := Rest(conf, logger)

	return func(ctx context.Context, req *Request) (*Response, error) {
		uri, err := url.Parse(req.Uri)
		if err != nil {
			return nil, fmt.Errorf("sender: %v", err)
		}

		// http & https
		if strings.HasPrefix(uri.Scheme, "http") {
			return rest(ctx, req)
		}

		return nil, fmt.Errorf("sender: unsupported scheme [%s]", uri.Scheme)
	}, nil
}

type Sender func(conf *Config) Send

type Send func(context.Context, *Request) (*Response, error)

type Request struct {
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}

type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}
