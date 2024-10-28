package restful

type Response struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Succeeded(data any) (status int, response *Response) {
	return Failed(OK, data)
}

func Failed(err *Error, data any) (status int, response *Response) {
	status = err.Status
	response = &Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    data,
	}
	return
}
