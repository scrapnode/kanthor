#!/bin/sh
set -e

CHECKSUM_FILE=./checksum

CHECKSUM_NEW=$(find . -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)

if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "vetting ..."
  # only vet on our packages
  go list ./... | grep kanthor | go vet -v
  find . -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > ./checksum
fi