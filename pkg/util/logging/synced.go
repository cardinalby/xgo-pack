package logging

import "sync"

type syncedWrapper struct {
	sync.Mutex
	Logger
}

func NewSyncedWrapper(logger Logger) SyncedLogger {
	return &syncedWrapper{Logger: logger}
}
