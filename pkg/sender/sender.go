package sender

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
	"net/url"
	"strings"
)

func New(conf *Config, logger logging.Logger) Send {
	rest := Rest(conf, logger)

	return func(req *Request) (*Response, error) {
		uri, err := url.Parse(req.Uri)
		if err != nil {
			return nil, fmt.Errorf("sender: %v", err)
		}

		// http & https
		if strings.HasPrefix(uri.Scheme, "http") {
			return rest(req)
		}

		return nil, fmt.Errorf("sender: unsupported scheme [%s]", uri.Scheme)
	}
}

type Sender func(conf *Config) Send

type Send func(req *Request) (*Response, error)

type Request struct {
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    string      `json:"body"`
}

type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    string      `json:"body"`
}
