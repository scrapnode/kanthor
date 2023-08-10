all: migration wire docs

migration:
	go run main.go migrate

wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go

docs:
	 swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal
	 swag init -q --instanceName Portal -d services/portalapi -o services/portalapi/docs -g swagger.go --parseDependency --parseInternal
