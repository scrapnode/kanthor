package circuitbreaker

import "errors"

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests = errors.New("INFRASTRUCTURE.CIRCUIT_BREAKER.STAGE_HALF_OPEN.TOO_MANY_REQUESTS.ERROR")
	// ErrOpenState is returned when the CB state is open
	ErrOpenState   = errors.New("INFRASTRUCTURE.CIRCUIT_BREAKER.STAGE_OPEN.OPENED.ERROR")
	ErrStageChange = errors.New("INFRASTRUCTURE.CIRCUIT_BREAKER.STAGE.CHANGE.ERROR")
)
