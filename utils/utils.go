package utils

func Contains[T comparable](value T, arr []T) bool {
	for _, k := range arr {
		if k == value {
			return true
		}
	}
	return false
}
