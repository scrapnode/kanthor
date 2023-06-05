package services

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(logger logging.Logger, message Message) Services {
	logger = logger.With("component", "dataplane.services")
	return &services{logger: logger, message: message}
}

type Services interface {
	patterns.Connectable
	Message() Message
}

type services struct {
	logger  logging.Logger
	message Message
}

func (service *services) Connect(ctx context.Context) error {
	if err := service.message.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("connected")
	return nil
}

func (service *services) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.message.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (service *services) Message() Message {
	return service.message
}
