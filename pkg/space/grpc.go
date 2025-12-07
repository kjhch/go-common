package space

import (
	"google.golang.org/grpc"
	"net"
)

type GrpcRegistrant interface {
	Register(s grpc.ServiceRegistrar)
	ServerAddr() string
}

type GrpcServer struct {
	logger     *Logger
	server     *grpc.Server
	registrant GrpcRegistrant
}

func NewGrpcServer(logger *Logger, registrant GrpcRegistrant) *GrpcServer {
	return &GrpcServer{
		logger:     logger,
		registrant: registrant,
		server:     grpc.NewServer(),
	}
}

func (s *GrpcServer) Start() {
	addr := s.registrant.ServerAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("Grpc服务监听失败", "addr", addr, "err", err)
		panic(err)
	}

	s.registrant.Register(s.server)

	s.logger.Info("Grpc服务已启动", "addr", addr)
	if err = s.server.Serve(lis); err != nil {
		s.logger.Error("Grpc服务运行失败", "addr", addr, "err", err)
		panic(err)
	}
	s.logger.Info("Grpc服务已关闭")
}

func (s *GrpcServer) Stop() {
	s.logger.Info("Grpc服务关闭中...")
	s.server.GracefulStop()
}
