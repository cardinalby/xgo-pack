package fnutil

func Once(f func()) func() {
	var called bool
	return func() {
		if called {
			return
		}
		called = true
		f()
	}
}
