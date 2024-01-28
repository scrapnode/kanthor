#!/bin/sh
set -e

STORAGE_PATH=${STORAGE_PATH:-"/tmp"}
TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS=${TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS:-"YWRtaW5Aa2FudGhvcmxhYnMuY29tOmNoYW5nZW1lbm93"}
TEST_KANTHOR_PORTAL_API_ENDPOINT=${TEST_KANTHOR_PORTAL_API_ENDPOINT:-"http://localhost:8280/api"}
TEST_WORKSPACE_SNAPSHOT_PATH=${TEST_WORKSPACE_SNAPSHOT_PATH:-"scripts/data/httpbin.json"}

NOW=$(date '+%Y-%m-%d %H:%M:%S')

# prepare new workspace with new application
IDEMPTOTENCY_KEY_WORKSPACE_CREATE=$(cat /proc/sys/kernel/random/uuid)
curl -X POST "$PORTAL_API_ENDPOINT/workspace" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_CREATE" \
    -H "Authorization: basic $PORTAL_AUTH_CREDENTIALS" \
    -H 'Content-Type: application/json' \
    -d "{\"name\": \"test workspace of $NOW\"}" > "$STORAGE_PATH/workspace.json"

TEST_WORKSPACE_ID=$(cat $STORAGE_PATH/workspace.json | jq -r '.id')
echo "Ws ID: $TEST_WORKSPACE_ID"

jq '{snapshot: .}' $TEST_WORKSPACE_SNAPSHOT_PATH > "$STORAGE_PATH/workspace.snapshot.json"

echo "----snapshot---"
cat "$STORAGE_PATH/workspace.snapshot.json"
echo "----snapshot---"

IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER=$(cat /proc/sys/kernel/random/uuid)
curl -X POST "$PORTAL_API_ENDPOINT/workspace/$TEST_WORKSPACE_ID/transfer" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER" \
    -H "Authorization: basic $PORTAL_AUTH_CREDENTIALS" \
    -d @$STORAGE_PATH/workspace.snapshot.json > "$STORAGE_PATH/workspace.transfer.json"

jq '{id: .app_id[0]}' "$STORAGE_PATH/workspace.transfer.json" > "$STORAGE_PATH/application.json"

# only retrive the app id from trust source
TEST_APP_ID=$(cat "$STORAGE_PATH/application.json" | jq -r '.id')

if [ -z "${TEST_APP_ID}" ]; then
    echo "App ID is empty"
    exit 1
fi

if [ $TEST_APP_ID = "null" ]; then
    echo "App ID is null"
    exit 2
fi

echo "App ID: $TEST_APP_ID"