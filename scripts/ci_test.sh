#!/bin/sh
set -e

go test --count=1 -cover -coverprofile cover.out  \
    ./configuration/... \
    ./database/... \
    ./datastore/... \
    ./gateway/... \
    ./infrastructure/... \
    ./internal/... \
    ./pkg/... \
    ./services/...

# to view coverage percentage on default browser, uncomment the line bellow
# go tool cover -html=cover.out