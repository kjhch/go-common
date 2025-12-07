package space

import "github.com/google/wire"

var BaseProviderSet = wire.NewSet(NewApp, NewConfigLoader, NewLogger, StartersProvider, NewGrpcServer, NewHttpServer)
var KafkaProvider = wire.NewSet(NewKafkaListener)
var NoKafkaProvider = wire.NewSet(wire.Value((*KafkaListener)(nil)))

func StartersProvider(grpcServer *GrpcServer, httpServer *HttpServer, listener *KafkaListener) []Starter {
	var starters []Starter
	if grpcServer != nil {
		starters = append(starters, grpcServer)
	}
	if httpServer != nil {
		starters = append(starters, httpServer)
	}
	if listener != nil {
		starters = append(starters, listener)
	}
	return starters
}
