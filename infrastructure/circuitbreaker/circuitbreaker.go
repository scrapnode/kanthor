package circuitbreaker

import "github.com/scrapnode/kanthor/infrastructure/logging"

func New(conf *Config, logger logging.Logger) CircuitBreaker {
	return NewSony(conf, logger)
}

type CircuitBreaker interface {
	Do(cmd string, onHandle Handler, onError ErrorHandler) (interface{}, error)
}

type Handler func() (interface{}, error)

type ErrorHandler func(err error) error

func Do[T any](cb CircuitBreaker, cmd string, onHandle Handler, onError ErrorHandler) (*T, error) {
	out, err := cb.Do(cmd, onHandle, onError)
	if err != nil {
		return nil, err
	}

	return out.(*T), nil
}
