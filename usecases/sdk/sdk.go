package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
	"sync"
)

type Sdk interface {
	patterns.Connectable
	Workspace() Workspace
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
	Message() Message
}

func New(
	conf *config.Config,
	logger logging.Logger,
	cryptography cryptography.Cryptography,
	timer timer.Timer,
	cache cache.Cache,
	publisher streaming.Publisher,
	repos repos.Repositories,
) Sdk {
	return &sdk{
		conf:         conf,
		logger:       logger,
		cryptography: cryptography,
		timer:        timer,
		cache:        cache,
		publisher:    publisher,
		repos:        repos,
	}
}

type sdk struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	publisher    streaming.Publisher
	repos        repos.Repositories

	mu           sync.RWMutex
	workspace    *workspace
	application  *application
	endpoint     *endpoint
	endpointRule *endpointRule
	message      *message
}

func (uc *sdk) Connect(ctx context.Context) error {
	if err := uc.cache.Connect(ctx); err != nil {
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
	uc.logger.Info("disconnected")

	if err := uc.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *sdk) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:         uc.conf,
			logger:       uc.logger,
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
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
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
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
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
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
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			repos:        uc.repos,
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
			cryptography: uc.cryptography,
			timer:        uc.timer,
			cache:        uc.cache,
			publisher:    uc.publisher,
			repos:        uc.repos,
		}
	}
	return uc.message
}
