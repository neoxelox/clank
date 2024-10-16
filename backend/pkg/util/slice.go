package util

func Filter[T any](slice []T, condition func(T) bool) []T {
	filtered := []T{}

	for _, item := range slice {
		if condition(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func Spread[T any](values []T) []any {
	anynized := make([]any, 0, len(values))

	for i := 0; i < len(values); i++ {
		anynized = append(anynized, values[i])
	}

	return anynized
}

func Extract[T any](slice []T, condition func(T) bool) ([]T, []T) {
	extracted := []T{}
	filtered := []T{}

	for _, item := range slice {
		if condition(item) {
			extracted = append(extracted, item)
		} else {
			filtered = append(filtered, item)
		}
	}

	return extracted, filtered
}
