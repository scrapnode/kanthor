package entrypoint

import (
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services/dispatcher/config"
	"github.com/scrapnode/kanthor/services/dispatcher/entrypoint/consumer"
	"github.com/scrapnode/kanthor/services/dispatcher/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Dispatcher,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, uc)
}
