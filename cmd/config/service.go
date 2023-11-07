package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services"
	attempt "github.com/scrapnode/kanthor/services/attempt/config"
	dispatcher "github.com/scrapnode/kanthor/services/dispatcher/config"
	portal "github.com/scrapnode/kanthor/services/portal/config"
	scheduler "github.com/scrapnode/kanthor/services/scheduler/config"
	sdk "github.com/scrapnode/kanthor/services/sdk/config"
	storage "github.com/scrapnode/kanthor/services/storage/config"
)

func Service(provider configuration.Provider, name string) (validator.Validator, error) {
	if name == services.SDK {
		return sdk.New(provider)
	}
	if name == services.PORTAL {
		return portal.New(provider)
	}
	if name == services.SCHEDULER {
		return scheduler.New(provider)
	}
	if name == services.DISPATCHER {
		return dispatcher.New(provider)
	}
	if name == services.STORAGE {
		return storage.New(provider)
	}
	if name == services.ATTEMPT {
		return attempt.New(provider)
	}

	return nil, fmt.Errorf("SYSTEM.SERVICE.UNKNOWN: %s", name)
}
