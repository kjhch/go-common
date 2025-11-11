package space

import (
	"errors"
	"net/http"
)

type Response struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"msg,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessfulRestResp(data any) (status int, response Response) {
	status = http.StatusOK
	response = Response{
		Code: "0",
		Data: data,
	}
	return
}

type restInfo interface {
	Code() string
	Message() string
	HttpStatus() int
}

func FailedRestResp(err error, data any) (status int, response *Response) {
	var restError restInfo
	// 默认http状态码保持200，用业务码区分
	status = http.StatusOK
	response = &Response{
		Code:    "-1",
		Message: err.Error(),
		Data:    data,
	}
	if errors.As(err, &restError) {
		status = restError.HttpStatus()
		response = &Response{
			Code:    restError.Code(),
			Message: restError.Message(),
			Data:    data,
		}
	}
	return
}
