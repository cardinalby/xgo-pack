package logging

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	WithPrefix(prefix string) Logger
}

type SyncedLogger interface {
	Logger
	AcquireLock() (logger Logger, unlock func())
}
