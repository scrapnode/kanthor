#!/bin/sh
set -e

K6_VUS=${K6_VUS:-"100"}
K6_DURATION=${K6_DURATION:-"15m"}
API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}

docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml stop 
docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml down
docker volume prune -f

docker compose -f docker-compose.yaml up -d streaming cache warehouse
echo "#1 sleep 5s"
sleep 5

go run main.go migrate database up && go run main.go migrate datastore up

docker compose -f docker-compose.debugging.yaml up -d
docker compose -f docker-compose.yaml up -d sdkapi portalapi scheduler dispatcher storage

echo "#2 sleep 5s"
sleep 5

go run main.go setup account kanthor_root_key --data=scripts/k6/httpbin.json --generate-credentials --output="$API_CREDS_PATH/sdkapi.json"
K6_VUS=$K6_VUS K6_DURATION=$K6_DURATION API_CREDS_PATH=$API_CREDS_PATH API_ENDPOINT=http://localhost:8180 k6 run scripts/k6/stability.js