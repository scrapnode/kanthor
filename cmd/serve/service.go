package serve

import (
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/ioc"
)

func Service(name string, conf *config.Config, logger logging.Logger) (services.Service, error) {
	if name == services.PORTAL_API {
		return ioc.InitializePortalApi(conf, logger)
	}
	if name == services.SDK_API {
		return ioc.InitializeSdkApi(conf, logger)
	}
	if name == services.SCHEDULER {
		return ioc.InitializeScheduler(conf, logger)
	}
	if name == services.DISPATCHER {
		return ioc.InitializeDispatcher(conf, logger)
	}
	if name == services.STORAGE {
		return ioc.InitializeStorage(conf, logger)
	}

	return nil, fmt.Errorf("serve.service: unknown service [%s]", name)
}
