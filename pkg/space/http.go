package space

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type HttpRegistrant interface {
	Register(server *http.Server)
	ServerAddr() string
}

type HttpServer struct {
	logger     *Logger
	server     *http.Server
	registrant HttpRegistrant
}

func NewHttpServer(logger *Logger, registrant HttpRegistrant) *HttpServer {
	return &HttpServer{
		logger:     logger,
		server:     new(http.Server),
		registrant: registrant,
	}
}

func (s *HttpServer) Start() {
	addr := s.registrant.ServerAddr()
	s.registrant.Register(s.server)
	s.server.Addr = addr
	s.logger.Info("Http服务已启动", "addr", addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Http服务运行失败", "addr", addr, "err", err)
		panic(err)
	}
	s.logger.Info("Http服务已关闭")
}

func (s *HttpServer) Stop() {
	s.logger.Info("Http服务关闭中...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Http服务关闭失败", "err", err)
	}
}

//------------------------------------------------------------------------------

type Response struct {
	Code    string `json:"code,omitempty"`
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

func FailedRestResp(err error, data any) (status int, response *Response) {
	var restError interface {
		Code() string
		Message() string
		HttpStatus() int
	}
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
