
gen-go:
	go generate infrastructure/ioc/generate.go
	go generate migration/ioc/generate.go
	go generate dataplane/ioc/generate.go
