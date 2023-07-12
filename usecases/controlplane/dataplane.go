package dataplane

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	repos repos.Repositories,
	cache cache.Cache,
	meter metric.Meter,
) Dataplane {
	return &dataplane{
		conf:   conf,
		logger: logger,
		timer:  timer,
		repos:  repos,
		cache:  cache,
		meter:  meter,
	}
}

type dataplane struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	repos  repos.Repositories
	cache  cache.Cache
	meter  metric.Meter
}
