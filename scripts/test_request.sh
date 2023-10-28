#!/bin/sh
set -e

export API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}
export API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8180"}

docker compose -f docker-compose.debugging.yaml up -d

go run main.go setup account kanthor_root_key --data=scripts/k6/httpbin.json --output="$API_CREDS_PATH/sdkapi.json"

export TEST_APP_ID=$(cat $API_CREDS_PATH/sdkapi.json | jq -r '.applications[0]')
export TEST_USERNAME=$(cat $API_CREDS_PATH/sdkapi.json | jq -r '.credentials.username')
export TEST_PASSWORD=$(cat $API_CREDS_PATH/sdkapi.json | jq -r '.credentials.password')
export AUTH_TOKEN=$(echo -n "$TEST_USERNAME:$TEST_PASSWORD" | base64 -w 0)
export REQUEST_ID=$(uuidgen)

curl -X PUT "$API_ENDPOINT/api/application/$TEST_APP_ID/message" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $REQUEST_ID" \
    -H "Authorization: Basic $AUTH_TOKEN" \
    --data-raw '{"type":"testing.traffic.request","body":{"hello":"world"},"headers":{"x-client":"curl"}}'