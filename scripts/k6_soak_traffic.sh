#!/bin/sh
set -e

export TEST_VUS=${K6_VUS:-"50"}
export TEST_DURATION_START=${K6_START_DURATION:-"5m"}
export TEST_DURATION_MID=${K6_MID_DURATION:-"30m"}
export TEST_DURATION_END=${K6_MID_DURATION:-"5m"}
export API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}
export API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8180"}

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

K6_VUS=$K6_VUS K6_START_DURATION=$K6_START_DURATION K6_MID_DURATION=$K6_MID_DURATION K6_END_DURATION=$K6_END_DURATION API_CREDS_PATH=$API_CREDS_PATH API_ENDPOINT=$API_ENDPOINT k6 run scripts/k6/stability.js