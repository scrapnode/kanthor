#!/usr/bin/env bash
set -e

OPENAPI_DIR=openapi
CLIENTS_DIR=clients

CLIENT_DIRS=("$CLIENTS_DIR/javascript")

for CLIENT_DIR in "${CLIENT_DIRS[@]}"
do
    echo "--> $CLIENT_DIR"
    rm -rf "$CLIENT_DIR/src/openapi"
    openapi-generator-cli generate \
        --global-property apiDocs=false,modelDocs=false \
        --ignore-file-override "$CLIENTS_DIR/.openapi-generator-ignore" \
        -g typescript \
        -i "$OPENAPI_DIR/Sdk_swagger.yaml" \
        -t "$CLIENT_DIRS/templates" \
        -o "$CLIENT_DIR/src/openapi" \
        -c "$CLIENT_DIR/openapi-generator-config.json" \
        --type-mappings=set=Array
done
