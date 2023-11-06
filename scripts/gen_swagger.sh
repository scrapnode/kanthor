#!/bin/sh
set -e

OPENAPI_FOLDER=openapi
CHECKSUM_FILE="$OPENAPI_FOLDER/checksum"

PORTAL_FOLDER=services/portal/entrypoint/rest
SDK_FOLDER=services/sdk/entrypoint/rest

CHECKSUM_NEW=$(find $OPENAPI_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  rm -rf "$OPENAPI_FOLDER/docs"

  echo "generating portal ...";
  swag init -q --instanceName Portal -d $PORTAL_FOLDER -g swagger.go -o $OPENAPI_FOLDER --parseDependency --parseInternal;

  echo "generating sdk ...";
  swag init -q --instanceName Sdk -d $SDK_FOLDER -g swagger.go -o $OPENAPI_FOLDER --parseDependency --parseInternal;

  echo "generating checksum ...";
  find $OPENAPI_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $CHECKSUM_FILE;
fi
 