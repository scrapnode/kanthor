gen: gen-go gen-proto

gen-proto:
	protoc --proto_path=infrastructure/gateway/grpc/protos --go_out=infrastructure/gateway/grpc/protos --go_opt=paths=source_relative --go-grpc_out=infrastructure/gateway/grpc/protos --go-grpc_opt=paths=source_relative infrastructure/gateway/grpc/protos/*.proto
	protoc --proto_path=services/controlplane/grpc/protos --go_out=services/controlplane/grpc/protos --go_opt=paths=source_relative --go-grpc_out=services/controlplane/grpc/protos --go-grpc_opt=paths=source_relative services/controlplane/grpc/protos/*.proto
	protoc --proto_path=services/dataplane/grpc/protos --go_out=services/dataplane/grpc/protos --go_opt=paths=source_relative --go-grpc_out=services/dataplane/grpc/protos --go-grpc_opt=paths=source_relative services/dataplane/grpc/protos/*.proto

gen-go:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go
