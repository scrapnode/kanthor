package serve

import (
	"fmt"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/ioc"
)

func Service(name string, conf *config.Config, logger logging.Logger) (services.Service, error) {
	if name == config.SERVICE_PORTAL_API {
		return ioc.InitializePortalApi(conf, logger)
	}
	if name == config.SERVICE_SDK_API {
		return ioc.InitializeSdkApi(conf, logger)
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
