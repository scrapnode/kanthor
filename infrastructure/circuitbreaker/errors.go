package circuitbreaker

import "errors"

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests = errors.New("CIRCUIT_BREAKER.STAGE_HALF_OPEN.TOO_MANY_REQUESTS")
	// ErrOpenState is returned when the CB state is open
	ErrOpenState   = errors.New("CIRCUIT_BREAKER.STAGE_OPEN.OPENED")
	ErrStageChange = errors.New("CIRCUIT_BREAKER.STAGE.CHANGE")
)
