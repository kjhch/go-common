package utils

func New[T any](val T) *T {
	return &val
}
