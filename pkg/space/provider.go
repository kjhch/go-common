package space

import (
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ProviderSet = wire.NewSet(NewApp, NewConfigLoader, NewLogger, NewPgxDB, NewRDB, NewMinioClient,
	StartersProvider, NewGrpcServer, NewHttpServer, NewKafkaListener, NewKafkaWriter)

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

func NewMinioClient(cl *ConfigLoader) *minio.Client {
	endpoint := cl.injectConf.Oss.Endpoint
	accessKeyID := cl.injectConf.Oss.Key
	secretAccessKey := cl.injectConf.Oss.Secret
	useSSL := false
	if endpoint == "" {
		return nil
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}
