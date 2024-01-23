package entrypoint

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint/consumer"
	"github.com/scrapnode/kanthor/services/attempt/entrypoint/cronjob"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
)

func Cronjob(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return cronjob.New(conf, logger, infra, db, ds, uc)
}

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, db, ds, uc)
}
