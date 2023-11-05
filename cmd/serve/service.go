package serve

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/ioc"
)

func Service(name string, provider configuration.Provider) (patterns.Runnable, error) {
	if name == services.PORTAL {
		return ioc.Portal(provider)
	}
	if name == services.SDK {
		return ioc.Sdk(provider)
	}
	if name == services.SCHEDULER {
		return ioc.Scheduler(provider)
	}
	if name == services.DISPATCHER {
		return ioc.Dispatcher(provider)
	}
	if name == services.STORAGE {
		return ioc.Storage(provider)
	}
	if name == services.ATTEMPT_TRIGGER_PLANNER {
		return ioc.AttemptTriggerPlanner(provider)
	}
	if name == services.ATTEMPT_TRIGGER_EXECUTOR {
		return ioc.AttemptTriggerExecutor(provider)
	}
	if name == services.ATTEMPT_ENDEAVOR_PLANNER {
		return ioc.AttemptEndeavorPlanner(provider)
	}
	if name == services.ATTEMPT_ENDEAVOR_EXECUTOR {
		return ioc.AttemptEndeavorExecutor(provider)
	}

	return nil, fmt.Errorf("serve.service: unknown service [%s]", name)
}
