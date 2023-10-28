#!/bin/sh
set -e

PORTAL_FOLDER=services/portal/entrypoint/rest
PORTAL_FILE_CHECKSUM="$PORTAL_FOLDER/docs/checksum"

SDK_FOLDER=services/sdk/entrypoint/rest
SDK_FILE_CHECKSUM="$SDK_FOLDER/docs/checksum"


CHECKSUM_NEW=$(find $PORTAL_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $PORTAL_FILE_CHECKSUM)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating portal docs ...";
  swag init -q --instanceName Sdk -d $PORTAL_FOLDER -o $PORTAL_FOLDER/docs -g swagger.go --parseDependency --parseInternal;
  find $PORTAL_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $PORTAL_FILE_CHECKSUM;
fi

CHECKSUM_NEW=$(find $SDK_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $SDK_FILE_CHECKSUM)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating sdk docs ...";
  swag init -q --instanceName Sdk -d $SDK_FOLDER -o $SDK_FOLDER/docs -g swagger.go --parseDependency --parseInternal;
  find $SDK_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $SDK_FILE_CHECKSUM;
fi