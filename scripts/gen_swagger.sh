#!/bin/sh
set -e

OPENAPI_DIR=openapi
CHECKSUM_FILE="$OPENAPI_DIR/checksum"

PORTAL_DIR=services/portal/entrypoint/rest
SDK_DIR=services/sdk/entrypoint/rest

CHECKSUM_NEW=$(find $OPENAPI_DIR -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  rm -rf "$OPENAPI_DIR/docs"

  echo "generating portal ...";
  swag init -q --instanceName Portal -d $PORTAL_DIR -g entrypoint_swagger.go -o $OPENAPI_DIR --parseDependency --parseInternal;

  echo "generating sdk ...";
  swag init -q --instanceName Sdk -d $SDK_DIR -g entrypoint_swagger.go -o $OPENAPI_DIR --parseDependency --parseInternal;

  echo "generating checksum ...";
  find $OPENAPI_DIR -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $CHECKSUM_FILE;
fi
 