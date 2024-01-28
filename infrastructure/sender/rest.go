package sender

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/scrapnode/kanthor/logging"
)

func Rest(conf *Config, logger logging.Logger) Send {
	logger = logger.With("sender", "rest")

	client := resty.New().
		SetLogger(logger).
		SetTimeout(time.Millisecond * time.Duration(conf.Timeout)).
		SetRetryCount(conf.Retry.Count).
		SetRetryWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTime)).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			status := r.StatusCode()
			url := r.Request.URL
			if status >= http.StatusInternalServerError {
				logger.Warnw("INFRASTRUCTURE.SENDER.REST.RETRYING", "status", status, "url", url)
				return true
			}
			return false
		})
	if conf.Trace {
		client = client.EnableTrace()
	}

	return func(ctx context.Context, req *Request) (*Response, error) {
		r := client.R().
			SetContext(ctx).
			SetHeaderMultiValues(req.Headers)

		var rp *resty.Response
		err := fmt.Errorf("INFRASTRUCTURE.SENDER.REST.METHOD.%s.NOT_SUPPORT.ERROR", strings.ToUpper(req.Method))

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

		logger.Debugw("INFRASTRUCTURE.SENDER.REST.SENT", "traces", rp.Request.TraceInfo())
		if err != nil {
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

		return res, nil
	}
}

type RestError struct {
	Status  string
	Headers http.Header
	Uri     string
	Body    []byte
}

func (err *RestError) Error() string {
	return err.Status
}
