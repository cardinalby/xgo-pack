package logging

import "sync"

type syncedWrapper struct {
	mu     *sync.Mutex
	logger Logger
}

func NewSyncedWrapper(logger Logger) SyncedLogger {
	return &syncedWrapper{
		logger: logger,
		mu:     &sync.Mutex{},
	}
}

func (w *syncedWrapper) AcquireLock() (logger Logger, unlock func()) {
	w.mu.Lock()
	return w.logger, w.mu.Unlock
}

func (w *syncedWrapper) Print(v ...any) {
	w.executeSynced(func(logger Logger) {
		logger.Print(v...)
	})
}

func (w *syncedWrapper) Printf(format string, v ...any) {
	w.executeSynced(func(logger Logger) {
		logger.Printf(format, v...)
	})
}

func (w *syncedWrapper) Println(v ...any) {
	w.executeSynced(func(logger Logger) {
		logger.Println(v...)
	})
}

func (w *syncedWrapper) WithPrefix(prefix string) Logger {
	return &syncedWrapper{
		logger: w.logger.WithPrefix(prefix),
		mu:     w.mu,
	}
}

func (w *syncedWrapper) executeSynced(f func(logger Logger)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	f(w.logger)
}
