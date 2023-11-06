all: ioc swagger

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh
