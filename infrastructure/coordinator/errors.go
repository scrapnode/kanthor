package coordinator

import (
	"errors"
)

var (
	ErrNotReady         = errors.New("COORDINATOR.NOT_READY")
	ErrNotLive          = errors.New("COORDINATOR.NOT_LIVE")
	ErrAlreadyConnected = errors.New("COORDINATOR.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("COORDINATOR.CONNECTION.NOT_CONNECTED")
)
