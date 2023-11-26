package entrypoint

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/scheduler/config"
	"github.com/scrapnode/kanthor/services/scheduler/entrypoint/consumer"
	"github.com/scrapnode/kanthor/services/scheduler/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Scheduler,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, db, uc)
}
