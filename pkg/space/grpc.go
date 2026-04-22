package space

import (
	"context"
	"fmt"
	"net"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcRegistrant interface {
	Register(s grpc.ServiceRegistrar)
}

type GrpcServer struct {
	cl         *ConfigLoader
	logger     *Logger
	server     *grpc.Server
	registrant GrpcRegistrant
}

func NewGrpcServer(cl *ConfigLoader, logger *Logger, registrant GrpcRegistrant) *GrpcServer {
	if cl.injectConf.Server.Grpc.Addr == "" {
		return nil
	}
	s := &GrpcServer{
		cl:         cl,
		logger:     logger,
		registrant: registrant,
	}
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.unaryInterceptor), grpc.StreamInterceptor(s.streamInterceptor))
	return s
}

func (s *GrpcServer) Start() {
	addr := s.cl.injectConf.Server.Grpc.Addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("Grpc服务监听失败"+addr, "err", err)
		panic(err)
	}

	s.registrant.Register(s.server)

	s.logger.Info("Grpc服务已启动" + addr)
	if err = s.server.Serve(lis); err != nil {
		s.logger.Error("Grpc服务运行失败"+addr, "err", err)
		panic(err)
	}
	s.logger.Info("Grpc服务已关闭")
}

func (s *GrpcServer) Stop() {
	s.logger.Info("Grpc服务关闭中...")
	s.server.GracefulStop()
}

func (s *GrpcServer) unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	value := md.Get(KeyRequestID)
	if len(value) > 0 {
		ctx = context.WithValue(ctx, KeyRequestID, value[0])
	}

	log := slices.Contains(s.cl.injectConf.Log.Grpc.EnabledMethods, info.FullMethod)
	if log {
		s.logger.InfoContext(ctx, fmt.Sprintf("[grpc]接口:%s, 参数:%s", info.FullMethod, req))
	}

	resp, err := handler(ctx, req) // 调用真正的业务逻辑

	if log || err != nil {
		s.logger.InfoContext(ctx, fmt.Sprintf("[grpc]接口:%s, 响应:%s", info.FullMethod, resp), "err", err)
	}

	return resp, err
}

func (s *GrpcServer) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	value := md.Get(KeyRequestID)
	if len(value) > 0 {
		ctx = context.WithValue(ctx, KeyRequestID, value[0])
	}

	// 包装 stream
	wrapped := &wrappedStream{
		ServerStream: ss,
		ctx:          ctx,
	}

	err := handler(srv, wrapped)
	s.logger.InfoContext(ctx, fmt.Sprintf("[grpc]接口:%s", info.FullMethod), "err", err)
	return err
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
