package streaming

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewPublisher(conf *PublisherConfig, logger logging.Logger) Publisher {
	return NewNatsPublisher(conf, logger)
}

type Publisher interface {
	patterns.Connectable
	Pub(ctx context.Context, event *Event) error
}

type PublisherConfig struct {
	*ConnectionConfig
}

func (conf PublisherConfig) Validate() error {
	return validator.New().Struct(conf)
}
