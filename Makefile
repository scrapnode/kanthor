gen: gen-go gen-proto

gen-proto:
	protoc --proto_path=dataplane/servers/grpc/protos --go_out=dataplane/servers/grpc/protos --go_opt=paths=source_relative --go-grpc_out=dataplane/servers/grpc/protos  --go-grpc_opt=paths=source_relative dataplane/servers/grpc/protos/message.proto

gen-go:
	go generate infrastructure/ioc/generate.go
	go generate migration/ioc/generate.go
	go generate dataplane/ioc/generate.go
