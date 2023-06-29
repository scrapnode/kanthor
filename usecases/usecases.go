package usecases

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane"
	"github.com/scrapnode/kanthor/usecases/scheduler"
)

func NewDataplane(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
) dataplane.Dataplane {
	logger = logger.With("usecase", "dataplane")
	return dataplane.New(conf, logger, timer, publisher, repos)
}

func NewScheduler(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	publisher streaming.Publisher,
	repos repositories.Repositories,
) scheduler.Scheduler {
	logger = logger.With("usecase", "scheduler")
	return scheduler.New(conf, logger, timer, publisher, repos)
}
