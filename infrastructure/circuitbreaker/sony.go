package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/logging"
	"github.com/sony/gobreaker"
)

func NewSony(conf *Config, logger logging.Logger) (CircuitBreaker, error) {
	logger = logger.With("circuitbreaker", "gobreaker")
	return &sonycb{
		conf:     conf,
		logger:   logger,
		breakers: map[string]*gobreaker.CircuitBreaker{},
	}, nil
}

type sonycb struct {
	conf   *Config
	logger logging.Logger

	breakers map[string]*gobreaker.CircuitBreaker
	mu       sync.RWMutex
}

func (cb *sonycb) Do(cmd string, handler Handler, onError ErrorHandler) (interface{}, error) {
	breaker := cb.get(cmd, onError)
	data, err := breaker.Execute(handler)
	// convert error
	if err != nil {
		return nil, cb.error(err)
	}

	return data, nil
}

func (cb *sonycb) get(cmd string, onError ErrorHandler) *gobreaker.CircuitBreaker {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if breaker, ok := cb.breakers[cmd]; ok {
		return breaker
	}

	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: cmd,
		// the maximum number of requests allowed to pass through when the CircuitBreaker is half-open
		MaxRequests: cb.conf.Half.PassthroughRequests,
		// the cyclic period of the closed state for CircuitBreaker to clear the internal Counts
		Interval: time.Millisecond * time.Duration(cb.conf.Close.CleanupInterval),
		// the period of the open state, after which the state of CircuitBreaker becomes half-open
		Timeout: time.Millisecond * time.Duration(cb.conf.Open.Duration),
		// if ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.ConsecutiveFailures >= cb.conf.Open.Conditions.ErrorConsecutive {
				cb.logger.Warnw("open because of error consecutively", "consecutive", counts.ConsecutiveFailures, "threshold", cb.conf.Open.Conditions.ErrorConsecutive)
				return true
			}

			reached := counts.TotalFailures > cb.conf.Half.PassthroughRequests
			ratio := float32(counts.TotalFailures) / float32(counts.Requests)
			if reached && ratio >= cb.conf.Open.Conditions.ErrorRatio {
				cb.logger.Warnw("open because of error ratio", "ratio", ratio, "threshold", cb.conf.Open.Conditions.ErrorRatio)
				return true
			}

			return false
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			cb.logger.Warnw(ErrStageChange.Error(), "name", name, "from", from.String(), "to", to.String())
		},
		IsSuccessful: func(err error) bool {
			return onError(err) == nil
		},
	})

	cb.breakers[cmd] = breaker
	return cb.breakers[cmd]
}

func (cb *sonycb) error(err error) error {
	if errors.Is(err, gobreaker.ErrTooManyRequests) {
		return ErrTooManyRequests
	}
	if errors.Is(err, gobreaker.ErrOpenState) {
		return ErrOpenState
	}
	return err
}
