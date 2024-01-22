#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

docker compose -f docker-compose.yaml down
docker volume prune -f

docker compose up -d cache warehouse streaming migration-database migration-datastore
docker compose ps