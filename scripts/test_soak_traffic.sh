#!/bin/sh
set -e

TEST_VUS=${TEST_VUS:-"50"}
TEST_DURATION_START=${TEST_DURATION_START:-"5s"}
TEST_DURATION_MID=${TEST_DURATION_MID:-"5m"}
TEST_DURATION_END=${K6_MID_DURATION:-"5s"}
API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8180"}

docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml -f docker-compose.debugging.yaml down
docker volume prune -f

docker compose -f docker-compose.yaml up -d streaming cache warehouse
echo "#1 sleep 5s"
sleep 5

go run main.go migrate database up && go run main.go migrate datastore up

docker compose -f docker-compose.debugging.yaml up -d
docker compose -f docker-compose.yaml up -d sdk portal scheduler storage

ENABLE_ATTEMPT=${ENABLE_ATTEMPT:-"0"}
if [ $ENABLE_ATTEMPT = "1" ]; then
    docker compose -f docker-compose.yaml up -d attempt-trigger-planner attempt-trigger-executor attempt-endeavor-planner attempt-endeavor-executor
fi

echo "#2 sleep 5s"
sleep 5

go run main.go setup account kanthor_root_key --data=scripts/k6/httpbin.json --output="$API_CREDS_PATH/sdk.json"

TEST_VUS=$TEST_VUS TEST_DURATION_START=$TEST_DURATION_START TEST_DURATION_MID=$TEST_DURATION_MID TEST_DURATION_END=$TEST_DURATION_END API_CREDS_PATH=$API_CREDS_PATH API_ENDPOINT=$API_ENDPOINT k6 run scripts/k6/stability.js