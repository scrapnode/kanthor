package repos

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.STORAGE.REPOS.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.STORAGE.REPOS.CONNECTION.ALREADY_CONNECTED")
)
