package usecases

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane"
	dataplanerepos "github.com/scrapnode/kanthor/usecases/dataplane/repos"
	"github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/scrapnode/kanthor/usecases/portal"
	portalrepos "github.com/scrapnode/kanthor/usecases/portal/repos"
	"github.com/scrapnode/kanthor/usecases/scheduler"
	schedulerrepos "github.com/scrapnode/kanthor/usecases/scheduler/repos"
)

func NewPortal(
	conf *config.Config,
	logger logging.Logger,
	cryptography cryptography.Cryptography,
	timer timer.Timer,
	cache cache.Cache,
	meter metric.Meter,
	repos portalrepos.Repositories,
) portal.Portal {
	logger = logger.With("usecase", "dataplane")
	return portal.New(conf, logger, cryptography, timer, cache, meter, repos)
}

func NewDataplane(
	conf *config.Config,
	logger logging.Logger,
	symmetric cryptography.Symmetric,
	timer timer.Timer,
	publisher streaming.Publisher,
	cache cache.Cache,
	meter metric.Meter,
	repos dataplanerepos.Repositories,
) dataplane.Dataplane {
	logger = logger.With("usecase", "dataplane")
	return dataplane.New(conf, logger, symmetric, timer, publisher, cache, meter, repos)
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
