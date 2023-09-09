all: wire docs

wire:
	./scripts/gen_ioc.sh

docs:
	./scripts/gen_docs.sh
