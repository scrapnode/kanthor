package sdk

import (
	"context"
	"sync"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
)

type Sdk interface {
	patterns.Connectable
	WorkspaceCredentials() WorkspaceCredentials
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
	Message() Message
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	publisher streaming.Publisher,
	repos repos.Repositories,
) Sdk {
	logger = logger.With("usecase", "sdk")

	return &sdk{
		conf:      conf,
		logger:    logger,
		infra:     infra,
		publisher: publisher,
		repos:     repos,
	}
}

type sdk struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
	repos     repos.Repositories

	mu                   sync.RWMutex
	workspaceCredentials *workspaceCredentials
	application          *application
	endpoint             *endpoint
	endpointRule         *endpointRule
	message              *message
}

func (uc *sdk) Readiness() error {
	if err := uc.infra.Readiness(); err != nil {
		return err
	}
	if err := uc.repos.Readiness(); err != nil {
		return err
	}
	if err := uc.publisher.Readiness(); err != nil {
		return err
	}
	return nil
}

func (uc *sdk) Liveness() error {
	if err := uc.infra.Liveness(); err != nil {
		return err
	}
	if err := uc.repos.Liveness(); err != nil {
		return err
	}
	if err := uc.publisher.Liveness(); err != nil {
		return err
	}
	return nil
}

func (uc *sdk) Connect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.infra.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	if err := uc.publisher.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *sdk) Disconnect(ctx context.Context) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.logger.Info("disconnected")

	if err := uc.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.infra.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *sdk) WorkspaceCredentials() WorkspaceCredentials {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspaceCredentials == nil {
		uc.workspaceCredentials = &workspaceCredentials{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.workspaceCredentials
}

func (uc *sdk) Application() Application {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.application == nil {
		uc.application = &application{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.application
}

func (uc *sdk) Endpoint() Endpoint {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endpoint == nil {
		uc.endpoint = &endpoint{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.endpoint
}

func (uc *sdk) EndpointRule() EndpointRule {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endpointRule == nil {
		uc.endpointRule = &endpointRule{
			conf:   uc.conf,
			logger: uc.logger,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}
	return uc.endpointRule
}

func (uc *sdk) Message() Message {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.message == nil {
		uc.message = &message{
			conf:      uc.conf,
			logger:    uc.logger,
			infra:     uc.infra,
			publisher: uc.publisher,
			repos:     uc.repos,
		}
	}
	return uc.message
}
