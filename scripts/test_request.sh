#!/bin/sh
set -e

API_CREDS_PATH=${API_CREDS_PATH:-"/tmp"}
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8180"}
REQUEST_WAIT_TIME=${REQUEST_WAIT_TIME:-1}
REQUEST_COUNT=${REQUEST_COUNT:-1}

go run main.go migrate database up && go run main.go migrate datastore up
go run main.go setup account admin@kanthorlabs.com --data=scripts/k6/httpbin.json --output="$API_CREDS_PATH/sdk.json"

TEST_APP_ID=$(cat $API_CREDS_PATH/sdk.json | jq -r '.applications[0]')
TEST_USERNAME=$(cat $API_CREDS_PATH/sdk.json | jq -r '.credentials.username')
TEST_PASSWORD=$(cat $API_CREDS_PATH/sdk.json | jq -r '.credentials.password')
AUTH_TOKEN=$(echo -n "$TEST_USERNAME:$TEST_PASSWORD" | base64 -w 0)

echo "sleep $REQUEST_WAIT_TIME second(s)"
sleep $REQUEST_WAIT_TIME

echo "sending $REQUEST_COUNT messages ---> $TEST_APP_ID"

for SEQ in $( seq 1 $REQUEST_COUNT )
do
    IDEMPTOTENCY_KEY=$(uuidgen)
    echo -n "\n--> $SEQ -> $IDEMPTOTENCY_KEY\n"
    curl "$API_ENDPOINT/api/message" \
        -H "Content-Type: application/json" \
        -H "Idempotency-Key: $IDEMPTOTENCY_KEY" \
        -H "X-Authorization-Engine: sdk.internal" \
        -H "Authorization: Basic $AUTH_TOKEN" \
        --data-raw "{\"app_id\": \"$TEST_APP_ID\",\"type\":\"testing.traffic.request\",\"body\":{\"hello\":\"world\",\"seq\":$SEQ},\"headers\":{\"x-client\":\"curl\"}}"
done

echo ""
echo "sent $REQUEST_COUNT messages ---> $TEST_APP_ID"
