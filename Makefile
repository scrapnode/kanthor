gen: gen-go gen-proto

gen-proto:
	protoc --proto_path=services/dataplane/grpc/protos --go_out=services/dataplane/grpc/protos --go_opt=paths=source_relative --go-grpc_out=services/dataplane/grpc/protos --go-grpc_opt=paths=source_relative services/dataplane/grpc/protos/dataplane.proto

gen-go:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go
