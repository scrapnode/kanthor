package serve

import (
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/ioc"
)

func Service(name string, conf *config.Config, logger logging.Logger) (services.Service, error) {
	if name == services.CONTROLPLANE {
		return ioc.InitializeControlplane(conf, logger)
	}
	if name == services.DATAPLANE {
		return ioc.InitializeDataplane(conf, logger)
	}
	if name == services.SCHEDULER {
		return ioc.InitializeScheduler(conf, logger)
	}
	if name == services.DISPATCHER {
		return ioc.InitializeDispatcher(conf, logger)
	}

	return nil, fmt.Errorf("serve.service: unknow service [%s]", name)
}