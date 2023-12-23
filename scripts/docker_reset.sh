#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml down
docker volume prune -f

docker compose up -d cache warehouse streaming migration-database migration-datastore
docker compose ps