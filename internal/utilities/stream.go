package utilities

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func FlattenArray[V any](arr [][]V) []V {
	var newArr []V
	for _, a := range arr {
		newArr = append(newArr, a...)
	}

	return newArr
}
