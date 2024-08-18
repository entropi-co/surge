package utilities

func StringDefault(value string, def string) string {
	if value != "" {
		return value
	} else {
		return def
	}
}

func OrDefault[T any](value *T, def *T) *T {
	if value != nil {
		return value
	}

	return def
}

func OrDefaultFn[T any](value *T, def func() *T) *T {
	if value != nil {
		return value
	}

	return def()
}

func Coalesce[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}

	return nil
}

func CoalesceFn[T any](args ...func() *T) *T {
	for _, arg := range args {
		if evaluated := arg(); evaluated != nil {
			return evaluated
		}
	}

	return nil
}

func CoalesceString(args ...string) string {
	for _, arg := range args {
		if arg != "" {
			return arg
		}
	}

	return ""
}
