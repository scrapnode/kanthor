package serve

import (
	"fmt"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services/ioc"
)

func Service(name string, conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	if name == config.SERVICE_PORTAL {
		return ioc.InitializePortal(conf, logger)
	}
	if name == config.SERVICE_SDK {
		return ioc.InitializeSdk(conf, logger)
	}
	if name == config.SERVICE_SCHEDULER {
		return ioc.InitializeScheduler(conf, logger)
	}
	if name == config.SERVICE_DISPATCHER {
		return ioc.InitializeDispatcher(conf, logger)
	}
	if name == config.SERVICE_STORAGE {
		return ioc.InitializeStorage(conf, logger)
	}
	if name == config.SERVICE_ATTEMPT_TRIGGER_PLANNER {
		return ioc.InitializeAttemptTriggerPlanner(conf, logger)
	}
	if name == config.SERVICE_ATTEMPT_TRIGGER_EXECUTOR {
		return ioc.InitializeAttemptTriggerExecutor(conf, logger)
	}

	return nil, fmt.Errorf("serve.service: unknown service [%s]", name)
}
