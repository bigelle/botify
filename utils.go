package botify

func notEmptyString(str string) func() bool {
	return func() bool {
		return len(str) != 0
	}
}

func notEmptyInt(i int) func() bool {
	return func() bool {
		return i != 0
	}
}

func notEmptySlice[T any](sl []T) func() bool {
	return func() bool {
		return len(sl) != 0 && sl != nil
	}
}
