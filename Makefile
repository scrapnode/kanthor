gen: wire migration docs

migration:
	go run main.go migrate

wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go

docs:
	 swag init -d services/sdkapi -o services/sdkapi/docs -g swagger.go
