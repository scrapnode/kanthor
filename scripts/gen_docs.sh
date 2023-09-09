#!/usr/bin/env bash
set -e

CHECKSUM_NEW=$(find services/portalapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat services/portalapi/docs/checksum)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating portalapi docs ...";
  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal;
  find services/portalapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/portalapi/docs/checksum;
fi

CHECKSUM_NEW=$(find services/sdkapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat services/sdkapi/docs/checksum)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating sdkapi docs ...";
  swag init -q --instanceName Sdk -d services/sdkapi -o services/sdkapi/docs -g swagger.go --parseDependency --parseInternal;
  find services/sdkapi -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/sdkapi/docs/checksum;
fi