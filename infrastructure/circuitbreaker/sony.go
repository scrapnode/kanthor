package circuitbreaker

import (
	"errors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/sony/gobreaker"
	"sync"
	"time"
)

func NewSony(conf *Config, logger logging.Logger) CircuitBreaker {
	return &sonycb{
		conf:     conf,
		logger:   logger,
		breakers: map[string]*gobreaker.CircuitBreaker{},
	}
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
		MaxRequests: uint32(cb.conf.HalfOpenMaxPassThroughRequests),
		// the cyclic period of the closed state for CircuitBreaker to clear the internal Counts
		Interval: time.Millisecond * time.Duration(cb.conf.CloseStateClearInterval),
		// the period of the open state, after which the state of CircuitBreaker becomes half-open
		Timeout: time.Millisecond * time.Duration(cb.conf.OpenStateDuration),
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			ratio := float64(counts.TotalFailures) / float64(counts.Requests)
			reachedRequestCount := int(counts.Requests) >= cb.conf.HalfOpenTriggerMinimumRequests
			reachedErrorRatio := ratio >= cb.conf.HalfOpenTriggerErrorThresholdRatio
			return reachedRequestCount && reachedErrorRatio
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			cb.logger.Warnw("Circuit Breaker stage change", "cb_name", name, "from", from.String(), "to", to.String(), "should_alert", from == gobreaker.StateOpen)
		},
		IsSuccessful: func(err error) bool {
			return onError(err) == nil
		},
	})

	cb.breakers[cmd] = breaker
	return breaker
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
