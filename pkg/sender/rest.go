package sender

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
	"time"
)

func Rest(conf *Config, logger logging.Logger) Send {
	logger = logger.With("pkg", "sender.rest")

	client := resty.New().
		SetLogger(logger).
		SetTimeout(time.Millisecond * time.Duration(conf.Timeout)).
		SetRetryCount(conf.Retry.Count).
		SetRetryWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTime)).
		SetRetryMaxWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTimeMax))

	return func(req *Request) (*Response, error) {
		r := client.R().
			SetHeaderMultiValues(req.Headers).
			SetBody(req.Body)

		logger.Debugw("sending", "uri", req.Uri)
		err := fmt.Errorf("sender.rest: unsupported method [%s]", req.Method)
		var rp *resty.Response

		if req.Method == http.MethodPost {
			rp, err = r.Post(req.Uri)
		}
		if req.Method == http.MethodPut {
			rp, err = r.Put(req.Uri)
		}
		if req.Method == http.MethodPatch {
			rp, err = r.Patch(req.Uri)
		}

		if err != nil {
			return nil, err
		}

		res := &Response{
			Status:  rp.StatusCode(),
			Headers: rp.Header(),
			// follow redirect url and got final url
			// most time the response url is same as request url
			Uri:  rp.RawResponse.Request.URL.String(),
			Body: string(rp.Body()),
		}
		logger.Debugw("sent", "uri", res.Uri)

		return res, nil
	}
}
