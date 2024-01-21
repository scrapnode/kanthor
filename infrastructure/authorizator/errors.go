package authorizator

import "errors"

var (
	ErrNotReady         = errors.New("AUTHORIZATOR.NOT_READY")
	ErrNotLive          = errors.New("AUTHORIZATOR.NOT_LIVE")
	ErrAlreadyConnected = errors.New("AUTHORIZATOR.ADAPTER.CONNECTION.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("AUTHORIZATOR.ADAPTER.CONNECTION.NOT_CONNECTED.ERROR")

	ErrWatcherAlreadyConnected = errors.New("AUTHORIZATOR.WATCHER.CONNECTION.ALREADY_CONNECTED.ERROR")
	ErrWatcherNotConnected     = errors.New("AUTHORIZATOR.WATCHER.CONNECTION.NOT_CONNECTED.ERROR")
)
