package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	subscriber streaming.Subscriber,
	uc usecase.Dispatcher,
	meter metric.Meter,
) services.Service {
	logger.With("service", "dispatcher")
	return &dispatcher{conf: conf, logger: logger, subscriber: subscriber, uc: uc, meter: meter}
}

type dispatcher struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	uc         usecase.Dispatcher
	meter      metric.Meter
}

func (service *dispatcher) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	if err := service.subscriber.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("started")
	return nil
}

func (service *dispatcher) Stop(ctx context.Context) error {
	service.logger.Info("stopped")

	if err := service.subscriber.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *dispatcher) Run(ctx context.Context) error {
	return service.subscriber.Sub(ctx, Consumer(service))
}
