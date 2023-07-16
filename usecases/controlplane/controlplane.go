package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/controlplane/repos"
	"sync"
)

type Controlplane interface {
	patterns.Connectable
	Workspace() Workspace
}

func New(
	conf *config.Config,
	logger logging.Logger,
	timer timer.Timer,
	cache cache.Cache,
	meter metric.Meter,
	repos repos.Repositories,
) Controlplane {
	return &controlplane{
		conf:   conf,
		logger: logger,
		timer:  timer,
		cache:  cache,
		meter:  meter,
		repos:  repos,
	}
}

type controlplane struct {
	conf   *config.Config
	logger logging.Logger
	timer  timer.Timer
	cache  cache.Cache
	meter  metric.Meter
	repos  repos.Repositories

	mu        sync.RWMutex
	worksapce *worksapce
}

func (usecase *controlplane) Connect(ctx context.Context) error {
	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *controlplane) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *controlplane) Workspace() Workspace {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.worksapce == nil {
		usecase.worksapce = &worksapce{
			conf:   usecase.conf,
			logger: usecase.logger,
			timer:  usecase.timer,
			cache:  usecase.cache,
			meter:  usecase.meter,
			repos:  usecase.repos,
		}
	}

	return usecase.worksapce
}
