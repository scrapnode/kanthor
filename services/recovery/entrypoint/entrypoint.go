package entrypoint

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/recovery/config"
	"github.com/scrapnode/kanthor/services/recovery/entrypoint/scanner"
	"github.com/scrapnode/kanthor/services/recovery/usecase"
)

func Scanner(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Recovery,
) patterns.Runnable {
	return scanner.New(conf, logger, infra, db, ds, uc)
}
