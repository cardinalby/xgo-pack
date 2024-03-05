package logging

import "sync"

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	WithPrefix(prefix string) Logger
}

type SyncedLogger interface {
	Logger
	sync.Locker
}
