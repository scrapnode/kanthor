package sender

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func Rest(conf *Config, logger logging.Logger) Send {
	logger = logger.With("pkg", "sender.rest")

	client := resty.New().
		SetLogger(logger).
		SetTimeout(time.Millisecond * time.Duration(conf.Timeout)).
		SetRetryCount(conf.Retry.Count).
		SetRetryWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTime)).
		SetRetryMaxWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTimeMax)).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			status := r.StatusCode()
			url := r.Request.URL
			if status >= http.StatusInternalServerError {
				logger.Warnw("retrying", "status", status, "url", url)
				return true
			}
			return false
		})
	if conf.EnableTrace {
		client = client.EnableTrace()
	}

	return func(ctx context.Context, req *Request) (*Response, error) {
		r := client.R().
			SetContext(ctx).
			SetHeaderMultiValues(req.Headers)
		logger.Debugw("sending", "uri", req.Uri)
		err := fmt.Errorf("sender.rest: unsupported method [%s]", req.Method)
		var rp *resty.Response

		if req.Method == http.MethodGet {
			rp, err = r.Get(req.Uri)
		}
		if req.Method == http.MethodPost {
			rp, err = r.SetBody(req.Body).Post(req.Uri)
		}
		if req.Method == http.MethodPut {
			rp, err = r.SetBody(req.Body).Put(req.Uri)
		}
		if req.Method == http.MethodPatch {
			rp, err = r.SetBody(req.Body).Patch(req.Uri)
		}

		if err != nil {
			logger.Debugw("sent", "trace_info", rp.Request.TraceInfo())
			return nil, err
		}

		res := &Response{
			Status:  rp.StatusCode(),
			Headers: rp.Header(),
			// follow redirect url and got final url
			// most time the response url is same as request url
			Uri:  rp.RawResponse.Request.URL.String(),
			Body: rp.Body(),
		}
		logger.Debugw("sent", "uri", res.Uri, "trace_info", rp.Request.TraceInfo())

		return res, nil
	}
}
