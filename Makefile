all: ioc swagger vet

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh
	
vet:
	./scripts/ci_vet.sh

