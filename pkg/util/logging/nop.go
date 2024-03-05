package logging

type nopLogger struct {
}

func NewNopLogger() Logger {
	return nopLogger{}
}

func (l nopLogger) Print(v ...any) {

}

func (l nopLogger) Printf(_ string, _ ...any) {

}

func (l nopLogger) Println(v ...any) {

}

func (l nopLogger) WithPrefix(_ string) Logger {
	return nopLogger{}
}
