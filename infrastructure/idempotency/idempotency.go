package idempotency

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"net/url"
	"strings"
)

func New(conf *Config, logger logging.Logger) (Idempotency, error) {
	uri, err := url.Parse(conf.Uri)
	if err != nil {
		logger.Warnw("unable to parse conf.Uri", "uri", conf.Uri)
		return nil, err
	}

	if strings.HasPrefix(uri.Scheme, "redis") {
		return NewRedis(conf, logger), nil
	}

	return nil, fmt.Errorf("idempotency: unknown engine")

}

type Idempotency interface {
	patterns.Connectable
	Validate(ctx context.Context, key string) (bool, error)
}
