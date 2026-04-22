package space

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
)

type HttpRegistrant interface {
	Register(server *http.Server)
}

type HttpServer struct {
	cl         *ConfigLoader
	logger     *Logger
	server     *http.Server
	registrant HttpRegistrant
}

func NewHttpServer(cl *ConfigLoader, logger *Logger, registrant HttpRegistrant) *HttpServer {
	if cl.injectConf.Server.Http.Addr == "" {
		return nil
	}
	return &HttpServer{
		cl:         cl,
		logger:     logger,
		server:     new(http.Server),
		registrant: registrant,
	}
}

func (s *HttpServer) Start() {
	addr := s.cl.injectConf.Server.Http.Addr
	s.registrant.Register(s.server)
	s.server.Addr = addr
	s.server.Handler = s.loggingMiddleware(s.server.Handler)
	s.logger.Info("Http服务已启动" + addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Http服务运行失败"+addr, "err", err)
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

func (s *HttpServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 请求前逻辑：记录请求URI
		if reqId := r.Header.Get("x-request-id"); reqId != "" {
			r = r.WithContext(context.WithValue(r.Context(), KeyRequestID, reqId))
		}

		if slices.ContainsFunc(s.cl.injectConf.Log.Http.EnabledRoutes, func(s string) bool {
			splits := strings.Split(s, " ")
			return splits[0] == r.Method && splits[1] == r.URL.Path
		}) {
			s.logger.InfoContext(r.Context(), fmt.Sprintf("[http]method:%s, uri:%s", r.Method, r.RequestURI))
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		// 请求后逻辑
	})
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
