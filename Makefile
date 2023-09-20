all: ioc docs

ioc:
	./scripts/gen_ioc.sh

docs:
	./scripts/gen_docs.sh
