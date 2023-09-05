all: wire docs

wire:
	go generate infrastructure/ioc/generate.go
	go generate services/ioc/generate.go

docs: docs-portalapi docs-sdkapi

docs-portalapi:
	CHECKSUM_NEW=$$(find services/portalapi -type f -name '*.go' -exec md5sum {} \; | sort -k 2 | md5sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat services/portalapi/docs/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating portalapi docs ..."; \
	  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal; \
	  find services/portalapi -type f -name '*.go' -exec md5sum {} \; | sort -k 2 | md5sum | cut -d  ' ' -f1 > services/portalapi/docs/checksum; \
	fi

docs-sdkapi:
	CHECKSUM_NEW=$$(find services/sdkapi -type f -name '*.go' -exec md5sum {} \; | sort -k 2 | md5sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat services/sdkapi/docs/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating sdkapi docs ..."; \
	  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal; \
	  find services/sdkapi -type f -name '*.go' -exec md5sum {} \; | sort -k 2 | md5sum | cut -d  ' ' -f1 > services/sdkapi/docs/checksum; \
	fi