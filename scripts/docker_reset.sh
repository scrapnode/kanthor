#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";
/bin/sh "$SCRIPT_DIR/docker_build.sh"

docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml down
docker volume prune -f

docker compose up -d
docker compose ps