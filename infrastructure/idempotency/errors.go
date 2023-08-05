package idempotency

import "errors"

var (
	ErrAlreadyConnected = errors.New("IDEMPOTENCY.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("IDEMPOTENCY.CONNECTION.NOT_CONNECTED")
)
