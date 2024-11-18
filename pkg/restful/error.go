package restful

type RestError interface {
	error
	Code() string
	Message() string
	Status() int
}
