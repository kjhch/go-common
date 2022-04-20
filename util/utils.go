package util

func NewPtr[a any](val a) *a {
	return &val
}
