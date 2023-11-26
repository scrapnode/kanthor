package entrypoint

import (
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/entrypoint/rest"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

func Rest(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Sdk,
) patterns.Runnable {
	return rest.New(conf, logger, infra, db, uc)
}
