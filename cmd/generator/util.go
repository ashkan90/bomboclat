package generator

func Keys[K comparable, V any](m map[K]V) []K {
	var keys = make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
