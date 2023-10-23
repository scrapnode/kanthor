package repos

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SDK.REPOS.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SDK.REPOS.CONNECTION.ALREADY_CONNECTED")
)
