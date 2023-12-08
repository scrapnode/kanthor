all: ioc swagger swagger_clients

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

swagger_clients:
	./scripts/gen_swagger_client.sh
