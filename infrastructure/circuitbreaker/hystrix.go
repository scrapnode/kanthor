package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"sync"
)

func NewHystrix(conf *Config, logger logging.Logger) CircuitBreaker {
	hystrix.SetLogger(&hystrixlog{logger: logger.With("component", "circuitbreaker.hystrix")})
	return &hystrixgo{conf: conf}
}

type hystrixgo struct {
	conf       *Config
	configured sync.Map
}

func (cb *hystrixgo) Do(cmd string, onHandle Handler, onError ErrorHandler) (interface{}, error) {
	if _, set := cb.configured.Load(cmd); !set {
		hystrix.ConfigureCommand(cmd, hystrix.CommandConfig{
			Timeout:               cb.conf.Timeout,
			SleepWindow:           cb.conf.SleepWindow,
			ErrorPercentThreshold: cb.conf.ErrorPercentThreshold,
		})
	}

	output := make(chan interface{}, 1)
	errors := hystrix.Go(cmd,
		func() error {
			out, err := onHandle()
			if err != nil {
				return err
			}
			output <- out
			return nil
		},
		func(err error) error {
			return onError(err)
		},
	)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}

type hystrixlog struct {
	logger logging.Logger
}

func (log *hystrixlog) Printf(format string, items ...interface{}) {
	log.logger.Debugf(format, items...)
}
