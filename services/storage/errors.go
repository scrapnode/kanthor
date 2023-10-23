package storage

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.STORAGE.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.STORAGE.ALREAD_STARTED")
)
