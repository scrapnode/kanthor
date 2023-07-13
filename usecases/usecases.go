package usecases

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane"
	controlplanerepos "github.com/scrapnode/kanthor/usecases/controlplane/repos"
	"github.com/scrapnode/kanthor/usecases/dataplane"
	dataplanerepos "github.com/scrapnode/kanthor/usecases/dataplane/repos"
	"github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/scrapnode/kanthor/usecases/scheduler"
	schedulerrepos "github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func NewControlplane(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	repos controlplanerepos.Repositories,
	cache cache.Cache,
	meter metric.Meter,
) controlplane.Controlplane {
	logger = logger.With("usecase", "controlplane")
	return controlplane.New(conf, logger, timer, cache, meter, repos)
}

func NewDataplane(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	cache cache.Cache,
	meter metric.Meter,
	repos dataplanerepos.Repositories,
) dataplane.Dataplane {
	logger = logger.With("usecase", "dataplane")
	return dataplane.New(conf, logger, timer, publisher, cache, meter, repos)
}

func NewScheduler(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	cache cache.Cache,
	meter metric.Meter,
	repos schedulerrepos.Repositories,
) scheduler.Scheduler {
	logger = logger.With("usecase", "scheduler")
	return scheduler.New(conf, logger, timer, publisher, cache, meter, repos)
}

func NewDispatcher(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	dispatch sender.Send,
	cache cache.Cache,
	cb circuitbreaker.CircuitBreaker,
	meter metric.Meter,
) dispatcher.Dispatcher {
	logger = logger.With("usecase", "scheduler")
	return dispatcher.New(conf, logger, timer, publisher, dispatch, cache, cb, meter)
}
