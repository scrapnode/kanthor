package streaming

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/validation"
)

func NewPublisher(conf *PublisherConfig, logger logging.Logger, validator validation.Validator) (Publisher, error) {
	return NewNatsPublisher(conf, logger, validator), nil
}

type Publisher interface {
	patterns.Connectable
	Pub(ctx context.Context, event *Event) error
}

type PublisherConfig struct {
	Connection ConnectionConfig `json:"connection" yaml:"connection" mapstructure:"connection" validate:"required"`
}

func (conf *PublisherConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
