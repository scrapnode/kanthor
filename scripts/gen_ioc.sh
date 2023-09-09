#!/usr/bin/env bash
set -e

CHECKSUM_NEW=$(find infrastructure/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat infrastructure/ioc/checksum)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating infrastructure ioc ...";
  go generate infrastructure/ioc/generate.go;
  find infrastructure/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > infrastructure/ioc/checksum;
fi

CHECKSUM_NEW=$(find services/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat services/ioc/checksum)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating services ioc ...";
  go generate services/ioc/generate.go;
  find services/ioc -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > services/ioc/checksum;
fi