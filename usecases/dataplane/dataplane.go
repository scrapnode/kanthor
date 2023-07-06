package dataplane

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
	cache cache.Cache,
	meter metric.Meter,
) Dataplane {
	return &dataplane{
		conf:      conf,
		logger:    logger,
		timer:     timer,
		publisher: publisher,
		repos:     repos,
		cache:     cache,
		meter:     meter,
	}
}

type dataplane struct {
	conf      *config.Config
	logger    logging.Logger
	timer     timer.Timer
	publisher streaming.Publisher
	repos     repositories.Repositories
	cache     cache.Cache
	meter     metric.Meter
}
