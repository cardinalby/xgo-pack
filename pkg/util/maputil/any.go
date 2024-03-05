package maputil

func AnyValue[K comparable, V any](m map[K]V, predicate func(V) bool) bool {
	for _, v := range m {
		if predicate(v) {
			return true
		}
	}
	return false
}
