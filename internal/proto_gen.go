//go:generate sh -c  "protoc --proto_path=../apis --go_out=../pkg --go-grpc_out=../pkg ../apis/*.proto"

package main
