package streaming

import "errors"

var (
	ErrSubAckFail = errors.New("STREAMING.SUBSCRIBER.ACK_FAIL")
	ErrSubNakFail = errors.New("STREAMING.SUBSCRIBER.NAK_FAIL")
)
