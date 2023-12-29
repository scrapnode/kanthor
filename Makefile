all: ioc swagger swagger_clients vet

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

swagger_clients:
	./scripts/gen_swagger_clients.sh
	
vet:
	./scripts/ci_vet.sh

