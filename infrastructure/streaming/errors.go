package streaming

import "errors"

var (
	ErrNotConnected = errors.New("STREAMING.CONNECTION.NOT_CONNECTED")
	ErrSubAckFail   = errors.New("STREAMING.SUBSCRIBER.ACK_FAIL")
	ErrSubNakFail   = errors.New("STREAMING.SUBSCRIBER.NAK_FAIL")
)
