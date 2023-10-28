package infrastructure

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(conf *config.Config, logger logging.Logger) (*Infrastructure, error) {
	t := timer.New()
	crypt, err := cryptography.New(&conf.Cryptography)
	if err != nil {
		return nil, err
	}
	send, err := sender.New(&conf.Sender, logger)
	if err != nil {
		return nil, err
	}
	idemp, err := idempotency.New(&conf.Idempotency, logger)
	if err != nil {
		return nil, err
	}
	cb, err := circuitbreaker.New(&conf.CircuitBreaker, logger)
	if err != nil {
		return nil, err
	}
	authz, err := authorizator.New(&conf.Authorizator, logger)
	if err != nil {
		return nil, err
	}
	lock, err := dlm.New(&conf.DistributedLockManager)
	if err != nil {
		return nil, err
	}
	m, err := metric.New(&conf.Metric, logger)
	if err != nil {
		return nil, err
	}
	c, err := cache.New(&conf.Cache, logger)
	if err != nil {
		return nil, err
	}
	s, err := streaming.New(&conf.Stream, logger)
	if err != nil {
		return nil, err
	}

	infra := &Infrastructure{
		conf:   conf,
		logger: logger.With("infrastructure", "default"),

		Timer:                  t,
		Cryptography:           crypt,
		Send:                   send,
		Idempotency:            idemp,
		CircuitBreaker:         cb,
		Authorizator:           authz,
		DistributedLockManager: lock,
		Stream:                 s,
		Cache:                  c,
		Metric:                 m,
	}
	return infra, nil
}

type Infrastructure struct {
	conf   *config.Config
	logger logging.Logger

	Timer                  timer.Timer
	Cryptography           cryptography.Cryptography
	Send                   sender.Send
	Idempotency            idempotency.Idempotency
	CircuitBreaker         circuitbreaker.CircuitBreaker
	DistributedLockManager dlm.Factory
	Authorizator           authorizator.Authorizator
	Stream                 streaming.Stream
	Cache                  cache.Cache
	Metric                 metric.Metric
}

func (infra *Infrastructure) Readiness() error {
	if err := infra.Idempotency.Readiness(); err != nil {
		return err
	}
	if err := infra.Authorizator.Readiness(); err != nil {
		return err
	}
	if err := infra.Stream.Readiness(); err != nil {
		return err
	}
	if err := infra.Cache.Readiness(); err != nil {
		return err
	}
	if err := infra.Metric.Readiness(); err != nil {
		return err
	}
	return nil
}

func (infra *Infrastructure) Liveness() error {
	if err := infra.Idempotency.Liveness(); err != nil {
		return err
	}
	if err := infra.Authorizator.Liveness(); err != nil {
		return err
	}
	if err := infra.Stream.Liveness(); err != nil {
		return err
	}
	if err := infra.Cache.Liveness(); err != nil {
		return err
	}
	if err := infra.Metric.Liveness(); err != nil {
		return err
	}
	return nil
}

func (infra *Infrastructure) Connect(ctx context.Context) error {
	if err := infra.Idempotency.Connect(ctx); err != nil {
		return err
	}
	if err := infra.Authorizator.Connect(ctx); err != nil {
		return err
	}
	if err := infra.Stream.Connect(ctx); err != nil {
		return err
	}
	if err := infra.Cache.Connect(ctx); err != nil {
		return err
	}
	if err := infra.Metric.Connect(ctx); err != nil {
		return err
	}

	infra.logger.Infow("connected")
	return nil
}

func (infra *Infrastructure) Disconnect(ctx context.Context) error {
	infra.logger.Infow("disconnected")
	var returning error

	if err := infra.Idempotency.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}
	if err := infra.Authorizator.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}
	if err := infra.Stream.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}
	if err := infra.Cache.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}
	if err := infra.Metric.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}
