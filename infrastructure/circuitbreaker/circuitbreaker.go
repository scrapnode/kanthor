package circuitbreaker

import "github.com/scrapnode/kanthor/logging"

func New(conf *Config, logger logging.Logger) (CircuitBreaker, error) {
	return NewSony(conf, logger)
}

type CircuitBreaker interface {
	Do(cmd string, onHandle Handler, onError ErrorHandler) (any, error)
}

type Handler func() (any, error)

type ErrorHandler func(err error) error

func Do[T any](cb CircuitBreaker, cmd string, onHandle Handler, onError ErrorHandler) (*T, error) {
	out, err := cb.Do(cmd, onHandle, onError)

	if out == nil {
		return nil, err
	}

	return out.(*T), nil
}
