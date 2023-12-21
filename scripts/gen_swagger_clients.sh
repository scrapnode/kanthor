#!/bin/sh
set -e

OPENAPI_DIR=openapi
CLIENTS_DIR=clients

# delimited strings is used as array
CLIENT_DIRS="$CLIENTS_DIR/golang"

for CLIENT_DIR in $CLIENT_DIRS
do
    echo "--> $CLIENT_DIR"
done
