package serve

import (
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/services"
)

func MetricExporter(name string, conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	if name == services.CONTROLPLANE {
		return metric.NewExporter(&conf.Controlplane.Metrics, logger), nil
	}
	if name == services.DATAPLANE {
		return metric.NewExporter(&conf.Dataplane.Metrics, logger), nil
	}
	if name == services.SCHEDULER {
		return metric.NewExporter(&conf.Scheduler.Metrics, logger), nil
	}
	if name == services.DISPATCHER {
		return metric.NewExporter(&conf.Dispatcher.Metrics, logger), nil
	}

	return nil, fmt.Errorf("serve.metric.exporter: unknow service [%s]", name)
}
