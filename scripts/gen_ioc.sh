#!/bin/sh
set -e

SERVICES_FOLDER=services
SERVICES_FILE_CHECKSUM="$SERVICES_FOLDER/ioc/checksum"

CHECKSUM_NEW=$(find $SERVICES_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $SERVICES_FILE_CHECKSUM || true)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating services ioc ...";
  go generate $SERVICES_FOLDER/ioc/generate.go;
  find $SERVICES_FOLDER -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $SERVICES_FILE_CHECKSUM;
fi