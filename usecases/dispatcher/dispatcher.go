package dispatcher

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
	dispatch sender.Send,
	cache cache.Cache,
	cb circuitbreaker.CircuitBreaker,
	meter metric.Meter,
) Dispatcher {
	return &dispatcher{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		repos:     repos,
		dispatch:  dispatch,
		cache:     cache,
		cb:        cb,
		meter:     meter,
	}
}

type dispatcher struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repositories.Repositories
	dispatch  sender.Send
	cache     cache.Cache
	cb        circuitbreaker.CircuitBreaker
	meter     metric.Meter
}
