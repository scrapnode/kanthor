#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";
REBUILD=${REBUILD:-"0"}

if [ "$REBUILD" == "1" ];
then
    /bin/sh "$SCRIPT_DIR/docker_build.sh"
fi


docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml down
docker volume prune -f

docker compose up -d cache warehouse streaming
docker compose ps