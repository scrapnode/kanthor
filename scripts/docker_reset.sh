#!/bin/sh
set -e

docker compose down && docker volume prune -f
docker compose up -d streaming cache warehouse
go run main.go migrate database up && go run main.go migrate datastore up