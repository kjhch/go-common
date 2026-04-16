package space

import (
	"net"

	"google.golang.org/grpc"
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
	return &GrpcServer{
		cl:         cl,
		logger:     logger,
		registrant: registrant,
		server:     grpc.NewServer(),
	}
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
