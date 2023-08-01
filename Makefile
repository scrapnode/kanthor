gen: gen-wire migrate

migrate:
	go run main.go migrate

gen-wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go
