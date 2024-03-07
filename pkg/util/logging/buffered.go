package logging

import (
	"slices"
	"sync"
)

type BufferedLogger struct {
	mu          sync.Mutex
	callsBuffer []func(logger Logger)
	logger      Logger
}

func NewBufferedLogger(logger Logger) *BufferedLogger {
	return &BufferedLogger{logger: logger}
}

func (l *BufferedLogger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.callsBuffer = append(l.callsBuffer, func(logger Logger) {
		logger.Print(v...)
	})
}

func (l *BufferedLogger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.callsBuffer = append(l.callsBuffer, func(logger Logger) {
		logger.Printf(format, v...)
	})
}

func (l *BufferedLogger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.callsBuffer = append(l.callsBuffer, func(logger Logger) {
		logger.Println(v...)
	})
}

func (l *BufferedLogger) WithPrefix(prefix string) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &BufferedLogger{
		logger:      l.logger.WithPrefix(prefix),
		callsBuffer: slices.Clone(l.callsBuffer),
	}
}

func (l *BufferedLogger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger := l.logger
	if syncedLogger, ok := l.logger.(SyncedLogger); ok {
		var unlock func()
		logger, unlock = syncedLogger.AcquireLock()
		defer unlock()
	}

	for _, call := range l.callsBuffer {
		call(logger)
	}
	l.callsBuffer = nil
}
