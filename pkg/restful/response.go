package restful

import (
	"errors"
	"net/http"
)

type Response struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Succeeded(data any) (status int, response *Response) {
	status = http.StatusOK
	response = &Response{
		Code: "0",
		Data: data,
	}
	return
}

func Failed(err error, data any) (status int, response *Response) {
	var restError RestError
	// 默认http状态码保持200，用业务码区分
	status = http.StatusOK
	response = &Response{
		Code:    "-1",
		Message: err.Error(),
		Data:    data,
	}
	if errors.As(err, &restError) {
		status = restError.Status()
		response = &Response{
			Code:    restError.Code(),
			Message: restError.Message(),
			Data:    data,
		}
	}
	return
}
