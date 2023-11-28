package entrypoint

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/attempt/config"
	endeavore "github.com/scrapnode/kanthor/services/attempt/entrypoint/endeavor/executor"
	endeavorp "github.com/scrapnode/kanthor/services/attempt/entrypoint/endeavor/planner"
	triggercli "github.com/scrapnode/kanthor/services/attempt/entrypoint/trigger/cli"
	triggere "github.com/scrapnode/kanthor/services/attempt/entrypoint/trigger/executor"
	triggerp "github.com/scrapnode/kanthor/services/attempt/entrypoint/trigger/planner"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func TriggerPlanner(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return triggerp.New(conf, logger, infra, db, ds, uc)
}

func TriggerExecutor(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return triggere.New(conf, logger, infra, db, ds, uc)
}

func TriggerCli(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.CommandLine {
	return triggercli.New(conf, logger, infra, db, ds, uc)
}

func EndeavorPlanner(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return endeavorp.New(conf, logger, infra, db, ds, uc)
}

func EndeavorExecutor(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return endeavore.New(conf, logger, infra, db, ds, uc)
}
