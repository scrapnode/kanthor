#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

MONITORING=${MONITORING:-""}
if [ "$MONITORING" != "" ];
then
  docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml down
  docker volume prune -f
  
  docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml up -d cache warehouse streaming migration-database migration-datastore
  
  docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml up -d clickhouse uptrace
else  
  docker compose -f docker-compose.yaml down
  docker volume prune -f

  docker compose up --build -d cache warehouse streaming migration-database migration-datastore
  docker compose ps
fi