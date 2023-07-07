package circuitbreaker

import "errors"

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests = errors.New("circuitbreaker: too many requests")
	// ErrOpenState is returned when the CB state is open
	ErrOpenState = errors.New("circuitbreaker: circuit breaker is open")
)
