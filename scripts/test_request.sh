#!/bin/sh
set -e

TEST_STORAGE_PATH=${TEST_STORAGE_PATH:-"/tmp"}
KANTHORPORTAL_API_ENDPOINT=${KANTHORPORTAL_API_ENDPOINT:-"http://localhost:8280/api"}
KANTHOR_PORTAL_AUTH_CREDENTIALS=${KANTHOR_PORTAL_AUTH_CREDENTIALS:-"YWRtaW5Aa2FudGhvcmxhYnMuY29tOmNoYW5nZW1lbm93"}
KANTHOR_SDK_API_ENDPOINT=${KANTHOR_SDK_API_ENDPOINT:-"http://localhost:8180/api"}
TEST_WORKSPACE_SNAPSHOT_PATH=${TEST_WORKSPACE_SNAPSHOT_PATH:-"data/snapshot.json"}
TEST_REQUEST_WAIT_TIME=${TEST_REQUEST_WAIT_TIME:-1}
TEST_REQUEST_COUNT=${TEST_REQUEST_COUNT:-1}

go run main.go migrate database up && go run main.go migrate datastore up

NOW=$(date '+%Y-%m-%d %H:%M:%S')
# prepare new workspace with new application
IDEMPTOTENCY_KEY_WORKSPACE_CREATE=$(uuidgen)
curl -s -X POST "$KANTHORPORTAL_API_ENDPOINT/workspace" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_CREATE" \
  -H "Authorization: basic $KANTHOR_PORTAL_AUTH_CREDENTIALS" \
  -H 'Content-Type: application/json' \
  -d "{\"name\": \"test workspace at $NOW\"}" > "$TEST_STORAGE_PATH/workspace.json"

TEST_WORKSPACE_ID=$(cat $TEST_STORAGE_PATH/workspace.json | jq -r '.id')
jq '{snapshot: .}' $TEST_WORKSPACE_SNAPSHOT_PATH > "$TEST_STORAGE_PATH/workspace.snapshot.json"
echo "Ws ID: $TEST_WORKSPACE_ID"

IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER=$(uuidgen)
curl -s -X POST "$KANTHORPORTAL_API_ENDPOINT/workspace/$TEST_WORKSPACE_ID/transfer" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER" \
  -H "Authorization: basic $KANTHOR_PORTAL_AUTH_CREDENTIALS" \
  -d @$TEST_STORAGE_PATH/workspace.snapshot.json > "$TEST_STORAGE_PATH/workspace.transfer.json"

jq '{id: .app_id[0]}' "$TEST_STORAGE_PATH/workspace.transfer.json" > "$TEST_STORAGE_PATH/application.json"

# only retrive the app id from trust source
TEST_APP_ID=$(cat "$TEST_STORAGE_PATH/application.json" | jq -r '.id')
echo "App ID: $TEST_APP_ID"

for SEQ in $( seq 1 $TEST_REQUEST_COUNT )
do
  sleep $TEST_REQUEST_WAIT_TIME
  IDEMPTOTENCY_KEY=$(uuidgen)
  echo -n "$IDEMPTOTENCY_KEY/$SEQ -> $TEST_APP_ID\n"
  curl -s -X POST "$KANTHOR_SDK_API_ENDPOINT/message" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY" \
    -H "X-Authorization-Engine: ask" \
    -H "X-Authorization-Workspace: $TEST_WORKSPACE_ID" \
    -H "Authorization: Basic $KANTHOR_PORTAL_AUTH_CREDENTIALS" \
    -d "{\"app_id\": \"$TEST_APP_ID\",\"type\":\"testing.traffic.request\",\"body\":{\"hello\":\"world\",\"seq\":$SEQ},\"headers\":{\"x-client\":\"curl\"}}"  
done

echo "App ID: $TEST_APP_ID | $TEST_REQUEST_COUNT messages"
