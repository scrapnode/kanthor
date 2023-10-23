package authorizator

import "errors"

var (
	ErrNotReady         = errors.New("AUTHORIZATOR.NOT_READY")
	ErrNotLive          = errors.New("AUTHORIZATOR.NOT_LIVE")
	ErrAlreadyConnected = errors.New("AUTHORIZATOR.ADAPTER.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("AUTHORIZATOR.ADAPTER.CONNECTION.NOT_CONNECTED")

	ErrWatcherAlreadyConnected = errors.New("AUTHORIZATOR.WATCHER.CONNECTION.ALREADY_CONNECTED")
	ErrWatcherNotConnected     = errors.New("AUTHORIZATOR.WATCHER.CONNECTION.NOT_CONNECTED")
)
