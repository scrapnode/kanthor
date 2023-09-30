package idempotency

import "errors"

var (
	ErrNotConnected     = errors.New("IDEMPOTENCY.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("IDEMPOTENCY.CONNECTION.ALREADY_CONNECTED")
)
