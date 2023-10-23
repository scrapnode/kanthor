package portalapi

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.PORTAL_API.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.PORTAL_API.ALREAD_STARTED")
)
