package usecase

import (
	"sync"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/sdk/config"
	"github.com/scrapnode/kanthor/services/sdk/repositories"
)

type Sdk interface {
	Workspace() Workspace
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
	Message() Message
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repositories repositories.Repositories,
) Sdk {
	logger = logger.With("usecase", "sdk")

	return &sdk{
		conf:         conf,
		logger:       logger,
		infra:        infra,
		repositories: repositories,
	}
}

type sdk struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	workspace    *workspace
	application  *application
	endpoint     *endpoint
	endpointRule *endpointRule
	message      *message

	mu sync.Mutex
}

func (uc *sdk) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.workspace
}

func (uc *sdk) Application() Application {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.application == nil {
		uc.application = &application{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.application
}

func (uc *sdk) Endpoint() Endpoint {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endpoint == nil {
		uc.endpoint = &endpoint{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.endpoint
}

func (uc *sdk) EndpointRule() EndpointRule {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endpointRule == nil {
		uc.endpointRule = &endpointRule{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			repositories: uc.repositories,
		}
	}
	return uc.endpointRule
}

func (uc *sdk) Message() Message {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.message == nil {
		uc.message = &message{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			publisher:    uc.infra.Stream.Publisher("sdk.message"),
			repositories: uc.repositories,
		}
	}
	return uc.message
}
