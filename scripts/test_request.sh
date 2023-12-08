#!/bin/sh
set -e

API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8180"}
REQUEST_WAIT_TIME=${REQUEST_WAIT_TIME:-1}

go run main.go migrate database up && go run main.go migrate datastore up
go run main.go setup account kanthor_root_key --data=scripts/k6/httpbin.json --output="$API_CREDS_PATH/sdk.json"

TEST_APP_ID=$(cat $API_CREDS_PATH/sdk.json | jq -r '.applications[0]')
TEST_USERNAME=$(cat $API_CREDS_PATH/sdk.json | jq -r '.credentials.username')
TEST_PASSWORD=$(cat $API_CREDS_PATH/sdk.json | jq -r '.credentials.password')
AUTH_TOKEN=$(echo -n "$TEST_USERNAME:$TEST_PASSWORD" | base64 -w 0)
REQUEST_ID=$(uuidgen)

echo "sleep $REQUEST_WAIT_TIME second(s)"
sleep $REQUEST_WAIT_TIME

curl --verbose -X PUT "$API_ENDPOINT/api/application/$TEST_APP_ID/message" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $REQUEST_ID" \
    -H "X-Authorization-Engine: sdk.internal" \
    -H "Authorization: Basic $AUTH_TOKEN" \
    --data-raw '{"type":"testing.traffic.request","body":{"hello":"world"},"headers":{"x-client":"curl"}}'