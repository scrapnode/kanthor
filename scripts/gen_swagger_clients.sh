#!/bin/sh
set -e

OPENAPI_DIR=openapi

GO_CLIENT_DIR=clients/sdk-go

rm -rf $GO_CLIENT_DIR/internal/openapi
openapi-generator-cli generate -i $OPENAPI_DIR/Sdk_swagger.json \
    --ignore-file-override .openapi-generator-ignore \
    -g go \
    -o $GO_CLIENT_DIR/internal/openapi \
    -c $GO_CLIENT_DIR/openapi-generator-config.json \
    -t $GO_CLIENT_DIR/templates