package entrypoint

import (
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/storage/config"
	"github.com/scrapnode/kanthor/services/storage/entrypoint/consumer"
	"github.com/scrapnode/kanthor/services/storage/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	ds datastore.Datastore,
	uc usecase.Storage,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, ds, uc)
}
