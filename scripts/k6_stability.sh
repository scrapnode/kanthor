#!/bin/sh
set -e

K6_VUS=${K6_VUS:-"100"}
K6_DURATION=${K6_DURATION:-"15m"}
API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}

docker compose up -d
docker compose -f docker-compose.debugging.yaml up -d

echo "sleep 10s"
sleep 10

go run main.go setup account kanthor_root_key --data=scripts/k6/httpbin.json --generate-credentials --output="$API_CREDS_PATH/sdkapi.json"
K6_VUS=$K6_VUS K6_DURATION=$K6_DURATION API_CREDS_PATH=$API_CREDS_PATH API_ENDPOINT=http://localhost:8180 k6 run scripts/k6/stability.js