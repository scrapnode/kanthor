package sdkapi

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.SDK_API.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.SDK_API.ALREAD_STARTED")
)
