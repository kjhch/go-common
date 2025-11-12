package space

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var ErrDomain = ""

type BizError interface {
	error
	Join(err error) error
	CustomizeMsg(msg string) BizError
}

type bizErr struct {
	code       string
	msg        string
	grpcStatus codes.Code
	httpStatus int
}

func (b *bizErr) Code() string {
	return b.code
}

func (b *bizErr) Message() string {
	return b.msg
}

func (b *bizErr) HttpStatus() int {
	return b.httpStatus
}

func (b *bizErr) GRPCStatus() *status.Status {
	return status.New(b.grpcStatus, b.Error())
}

func (b *bizErr) Error() string {
	return fmt.Sprintf("%v(%v:%v)", b.msg, ErrDomain, b.code)
}

func (b *bizErr) Join(err error) error {
	return errors.Join(b, err)
}

func (b *bizErr) CustomizeMsg(msg string) BizError {
	return &bizErr{
		code:       b.code,
		msg:        msg,
		grpcStatus: b.grpcStatus,
		httpStatus: b.httpStatus,
	}
}

func NewErr(code, msg string, grpcStatus, httpStatus int) BizError {
	r := &bizErr{
		code:       code,
		msg:        msg,
		httpStatus: httpStatus,
		grpcStatus: codes.Code(grpcStatus),
	}
	return r
}

//type ErrOpt func(err *bizErr)
//
//func WithErrGrpcStatus(status int) ErrOpt {
//	return func(err *bizErr) {
//		err.grpcStatus = codes.Code(status)
//	}
//}
//func WithErrHttpStatus(status int) ErrOpt {
//	return func(err *bizErr) {
//		err.httpStatus = status
//	}
//}

var (
	ErrUnauthenticated  = NewErr("UNAUTHENTICATED", "用户未认证，请登录", int(codes.Unauthenticated), http.StatusUnauthorized)
	ErrPermissionDenied = NewErr("NO_PERMISSION", "用户无权限", int(codes.PermissionDenied), http.StatusForbidden)

	ErrArgsRequired = NewErr("ARGS_REQUIRED", "必填参数为空，请检查后重试", int(codes.InvalidArgument), http.StatusBadRequest)
	ErrFileTooLarge = NewErr("FILE_TOO_LARGE", "文件大小超出限制，请更换后再试", int(codes.InvalidArgument), http.StatusRequestEntityTooLarge)
	ErrFileType     = NewErr("FILE_TYPE_ERR", "文件类型不支持，请更换后再试", int(codes.InvalidArgument), http.StatusBadRequest)

	ErrTooManyRequests = NewErr("TOO_MANY_REQS", "请求过于频繁，请稍后再试", int(codes.ResourceExhausted), http.StatusTooManyRequests)

	ErrParseJson         = NewErr("PARSE_JSON_ERR", "服务端数据解析失败，请稍后再试", int(codes.Internal), http.StatusInternalServerError)
	ErrDatabaseService   = NewErr("DB_SVC_ERR", "数据库服务异常，请稍后再试", int(codes.Internal), http.StatusInternalServerError)
	ErrCacheService      = NewErr("CACHE_SVC_ERR", "缓存服务异常，请稍后再试", int(codes.Internal), http.StatusInternalServerError)
	ErrStorageService    = NewErr("STORAGE_SVC_ERR", "存储服务异常，请稍后再试", int(codes.Internal), http.StatusInternalServerError)
	ErrDownstreamService = NewErr("DOWNSTREAM_SVC_ERR", "下游服务异常，请稍后再试", int(codes.Unavailable), http.StatusServiceUnavailable)
)
