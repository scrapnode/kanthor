package idempotency

import "errors"

var (
	ErrNotConnected     = errors.New("INFRASTRUCTURE.IDEMPOTENCY.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("INFRASTRUCTURE.IDEMPOTENCY.ALREADY_CONNECTED.ERROR")
)
