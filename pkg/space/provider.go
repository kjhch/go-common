package space

import "github.com/google/wire"

type ProviderConf struct {
	ExcludeGrpc  bool
	ExcludeHttp  bool
	ExcludeKafka bool
}

func Providers(conf ProviderConf) wire.ProviderSet {
	providers := []any{NewApp, NewConfigLoader, NewLogger, startersProvider}

	if conf.ExcludeGrpc {
		providers = append(providers, wire.Value((*GrpcServer)(nil)))
	} else {
		providers = append(providers, NewGrpcServer)

	}

	if conf.ExcludeHttp {
		providers = append(providers, wire.Value((*HttpServer)(nil)))
	} else {
		providers = append(providers, NewHttpServer)

	}

	if conf.ExcludeKafka {
		providers = append(providers, wire.Value((*KafkaListener)(nil)))
	} else {
		providers = append(providers, NewKafkaListener)

	}

	return wire.NewSet(providers...)
}

func startersProvider(grpcServer *GrpcServer, httpServer *HttpServer, listener *KafkaListener) []Starter {
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
