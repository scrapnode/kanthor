#!/bin/sh
set -e

CHECKSUM_NEW=$(find services -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat services/ioc/checksum)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating services ioc ...";
  go generate services/ioc/generate.go;
  find services -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/ioc/checksum;
fi