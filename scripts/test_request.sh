#!/bin/sh
set -e

STORAGE_PATH=${STORAGE_PATH:-"/tmp"}
PORTAL_AUTH_CREDENTIALS=${PORTAL_AUTH_CREDENTIALS:-"YWRtaW5Aa2FudGhvcmxhYnMuY29tOmNoYW5nZW1lbm93"}
PORTAL_API_ENDPOINT=${PORTAL_API_ENDPOINT:-"http://localhost:8280/api"}
SDK_API_ENDPOINT=${SDK_API_ENDPOINT:-"http://localhost:8180/api"}
REQUEST_WAIT_TIME=${REQUEST_WAIT_TIME:-1}
REQUEST_COUNT=${REQUEST_COUNT:-1}
TEST_WORKSPACE_SNAPSHOT_PATH=${TEST_WORKSPACE_SNAPSHOT_PATH:-"scripts/k6/httpbin.json"}

go run main.go migrate database up && go run main.go migrate datastore up

# prepare new workspace with new application

IDEMPTOTENCY_KEY_WORKSPACE_CREATE=$(uuidgen)
curl -s -X POST "$PORTAL_API_ENDPOINT/workspace" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_CREATE" \
    -H "Authorization: basic $PORTAL_AUTH_CREDENTIALS" \
    -H 'Content-Type: application/json' \
    -d '{"name": "test request"}' > "$STORAGE_PATH/workspace.json"

TEST_WORKSPACE_ID=$(cat $STORAGE_PATH/workspace.json | jq -r '.id')
jq '{snapshot: .}' $TEST_WORKSPACE_SNAPSHOT_PATH > "$STORAGE_PATH/workspace.snapshot.json"

IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER=$(uuidgen)
curl -s -X POST "$PORTAL_API_ENDPOINT/workspace/$TEST_WORKSPACE_ID/transfer" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER" \
    -H "Authorization: basic $PORTAL_AUTH_CREDENTIALS" \
    -d @$STORAGE_PATH/workspace.snapshot.json > "$STORAGE_PATH/workspace.transfer.json"

jq '{id: .app_id[0]}' "$STORAGE_PATH/workspace.transfer.json" > "$STORAGE_PATH/application.json"

# only retrive the app id from trust source
TEST_APP_ID=$(cat "$STORAGE_PATH/application.json" | jq -r '.id')
echo "App ID: $TEST_APP_ID"

for SEQ in $( seq 1 $REQUEST_COUNT )
do
    IDEMPTOTENCY_KEY=$(uuidgen)
    echo -n "$IDEMPTOTENCY_KEY/$SEQ -> $TEST_APP_ID\n"
    curl -s -X POST "$SDK_API_ENDPOINT/message" \
        -H "Content-Type: application/json" \
        -H "Idempotency-Key: $IDEMPTOTENCY_KEY" \
        -H "X-Authorization-Engine: ask" \
        -H "X-Authorization-Workspace: $TEST_WORKSPACE_ID" \
        -H "Authorization: Basic $PORTAL_AUTH_CREDENTIALS" \
        -d "{\"app_id\": \"$TEST_APP_ID\",\"type\":\"testing.traffic.request\",\"body\":{\"hello\":\"world\",\"seq\":$SEQ},\"headers\":{\"x-client\":\"curl\"}}" > /dev/null
done

echo "App ID: $TEST_APP_ID | $REQUEST_COUNT messages"
