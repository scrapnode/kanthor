gen: wire migration docs

migration:
	go run main.go migrate

wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go

docs:
	 swag init -q -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal
