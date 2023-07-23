package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
	"sync"
)

type Dataplane interface {
	patterns.Connectable
	Message() Message
	Application() Application
	AppCreds() AppCreds
}

func New(
	conf *config.Config,
	logger logging.Logger,
	symmetric cryptography.Symmetric,
	timer timer.Timer,
	publisher streaming.Publisher,
	cache cache.Cache,
	meter metric.Meter,
	authorizator authorizator.Authorizator,
	repos repos.Repositories,
) Dataplane {
	return &dataplane{
		conf:         conf,
		logger:       logger,
		symmetric:    symmetric,
		timer:        timer,
		publisher:    publisher,
		cache:        cache,
		meter:        meter,
		authorizator: authorizator,
		repos:        repos,
	}
}

type dataplane struct {
	conf         *config.Config
	logger       logging.Logger
	symmetric    cryptography.Symmetric
	timer        timer.Timer
	publisher    streaming.Publisher
	cache        cache.Cache
	meter        metric.Meter
	authorizator authorizator.Authorizator
	repos        repos.Repositories

	mu          sync.RWMutex
	message     *message
	application *application
	appcreds    *appcreds
}

func (usecase *dataplane) Connect(ctx context.Context) error {
	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.authorizator.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *dataplane) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.authorizator.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (usecase *dataplane) Message() Message {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.message == nil {
		usecase.message = &message{
			conf:      usecase.conf,
			logger:    usecase.logger,
			timer:     usecase.timer,
			publisher: usecase.publisher,
			repos:     usecase.repos,
			cache:     usecase.cache,
			meter:     usecase.meter,
		}
	}

	return usecase.message
}

func (usecase *dataplane) Application() Application {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.application == nil {
		usecase.application = &application{
			conf:         usecase.conf,
			logger:       usecase.logger,
			symmetric:    usecase.symmetric,
			timer:        usecase.timer,
			repos:        usecase.repos,
			cache:        usecase.cache,
			meter:        usecase.meter,
			authorizator: usecase.authorizator,
		}
	}

	return usecase.application
}

func (usecase *dataplane) AppCreds() AppCreds {
	usecase.mu.Lock()
	defer usecase.mu.Unlock()

	if usecase.appcreds == nil {
		usecase.appcreds = &appcreds{
			conf:         usecase.conf,
			logger:       usecase.logger,
			symmetric:    usecase.symmetric,
			timer:        usecase.timer,
			repos:        usecase.repos,
			cache:        usecase.cache,
			meter:        usecase.meter,
			authorizator: usecase.authorizator,
		}
	}

	return usecase.appcreds
}
