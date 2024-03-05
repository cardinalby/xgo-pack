package typeutil

func Ptr[T any](value T) *T {
	return &value
}

func PtrValueOrDefault[T any](ptr *T) T {
	if ptr == nil {
		var empty T
		return empty
	}
	return *ptr
}

func PtrValueOr[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
