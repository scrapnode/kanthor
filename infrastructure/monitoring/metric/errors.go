package metric

import "errors"

var (
	ErrNotConnected     = errors.New("MONITORING.METRIC.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("MONITORING.METRIC.CONNECTION.ALREADY_CONNECTED")
)
