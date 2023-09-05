all: wire docs

wire: wire-infrastructure wire-services

wire-infrastructure:
	CHECKSUM_NEW=$$(find infrastructure/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat infrastructure/ioc/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating infrastructure ioc ..."; \
	  go generate infrastructure/ioc/generate.go; \
	  find infrastructure/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > infrastructure/ioc/checksum; \
	fi

wire-services:
	CHECKSUM_NEW=$$(find services/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat services/ioc/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating services ioc ..."; \
	  go generate services/ioc/generate.go; \
	  find services/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/ioc/checksum; \
	fi

docs: docs-portalapi docs-sdkapi

docs-portalapi:
	CHECKSUM_NEW=$$(find services/portalapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat services/portalapi/docs/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating portalapi docs ..."; \
	  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal; \
	  find services/portalapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/portalapi/docs/checksum; \
	fi

docs-sdkapi:
	CHECKSUM_NEW=$$(find services/sdkapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1) && \
	CHECKSUM_OLD=$$(cat services/sdkapi/docs/checksum) && \
	if [ "$$CHECKSUM_NEW" != "$$CHECKSUM_OLD" ]; \
	then \
	  echo "generating sdkapi docs ..."; \
	  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal; \
	  find services/sdkapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/sdkapi/docs/checksum; \
	fi