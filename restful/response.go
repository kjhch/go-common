package restful

import "github.com/kjhch/go-common/restful/errno"

type Response struct {
	Code          string `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	ReferenceLink string `json:"reference_link,omitempty"`
	Data          any    `json:"data,omitempty"`
}

func NewResponse(errcode string) (status int, response *Response) {
	err, ok := errno.ErrMap[errcode]
	if ok {
		status = err.Status
		response = &Response{
			Code:    err.Code,
			Message: err.Message,
		}
		return
	}
	status = 500
	response = &Response{
		Code:    "-1",
		Message: "Unknown error",
	}
	return
}
