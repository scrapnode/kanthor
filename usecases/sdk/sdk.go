package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/sdk/repos"
	"sync"
)

type SDK interface {
	patterns.Connectable
	Application() Application
}

func New(
	conf *config.Config,
	logger logging.Logger,
	cryptography cryptography.Cryptography,
	timer timer.Timer,
	cache cache.Cache,
	meter metric.Meter,
	repos repos.Repositories,
) SDK {
	return &sdk{
		conf:         conf,
		logger:       logger,
		cryptography: cryptography,
		timer:        timer,
		cache:        cache,
		meter:        meter,
		repos:        repos,
	}
}

type sdk struct {
	conf         *config.Config
	logger       logging.Logger
	cryptography cryptography.Cryptography
	timer        timer.Timer
	cache        cache.Cache
	meter        metric.Meter
	repos        repos.Repositories

	mu          sync.RWMutex
	application *application
}

func (uc *sdk) Connect(ctx context.Context) error {
	if err := uc.cache.Connect(ctx); err != nil {
		return err
	}

	if err := uc.repos.Connect(ctx); err != nil {
		return err
	}

	uc.logger.Info("connected")
	return nil
}

func (uc *sdk) Disconnect(ctx context.Context) error {
	uc.logger.Info("disconnected")

	if err := uc.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := uc.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (uc *sdk) Application() Application {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.application == nil {
		uc.application = &application{}
	}
	return uc.application
}
