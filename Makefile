gen: gen-wire

gen-wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go
