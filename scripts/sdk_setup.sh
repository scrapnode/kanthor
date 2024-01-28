#!/bin/sh
set -e

STORAGE_PATH=${STORAGE_PATH:-"/tmp"}
TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS=${TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS:-"YWRtaW5Aa2FudGhvcmxhYnMuY29tOmNoYW5nZW1lbm93"}
TEST_KANTHOR_PORTAL_API_ENDPOINT=${TEST_KANTHOR_PORTAL_API_ENDPOINT:-"http://localhost:8280/api"}
TEST_WORKSPACE_SNAPSHOT_PATH=${TEST_WORKSPACE_SNAPSHOT_PATH:-"scripts/data/httpbin.json"}

NOW=$(date '+%Y-%m-%d %H:%M:%S')

# prepare new workspace with new application
IDEMPTOTENCY_KEY_WORKSPACE_CREATE=$(uuidgen)
curl -s -X POST "$TEST_KANTHOR_PORTAL_API_ENDPOINT/workspace" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_CREATE" \
    -H "Authorization: basic $TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS" \
    -H 'Content-Type: application/json' \
    -d "{\"name\": \"sdk test at $NOW\"}" > "$STORAGE_PATH/workspace.json"

TEST_WORKSPACE_ID=$(cat $STORAGE_PATH/workspace.json | jq -r '.id')



TEST_WORKSPACE_ID=$(cat $STORAGE_PATH/workspace.json | jq -r '.id')
jq '{snapshot: .}' $TEST_WORKSPACE_SNAPSHOT_PATH > "$STORAGE_PATH/workspace.snapshot.json"

IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER=$(uuidgen)
curl -s -X POST "$TEST_KANTHOR_PORTAL_API_ENDPOINT/workspace/$TEST_WORKSPACE_ID/transfer" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_TRANSFER" \
    -H "Authorization: basic $TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS" \
    -d @$STORAGE_PATH/workspace.snapshot.json > "$STORAGE_PATH/workspace.transfer.json"

jq '{id: .app_id[0]}' "$STORAGE_PATH/workspace.transfer.json" > "$STORAGE_PATH/application.json"
TEST_APPLICATION_ID=$(cat $STORAGE_PATH/application.json | jq -r '.id')
echo -n "$TEST_APPLICATION_ID"  > "$STORAGE_PATH/application.plain"



IDEMPTOTENCY_KEY_WORKSPACE_CREDENTIALS_GENERATE=$(uuidgen)
WORKSPACE_CREDENTIALS_EXPIRED_AT=$(date -d '+1 hour' '+%s%N' | cut -b1-13)
curl -s -X POST "$TEST_KANTHOR_PORTAL_API_ENDPOINT/credentials" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $IDEMPTOTENCY_KEY_WORKSPACE_CREDENTIALS_GENERATE" \
    -H "X-Authorization-Engine: ask" \
    -H "X-Authorization-Workspace: $TEST_WORKSPACE_ID" \
    -H "Authorization: basic $TEST_KANTHOR_PORTAL_AUTH_CREDENTIALS" \
    -d "{\"name\": \"sdk test at $NOW\", \"expired_at\": $WORKSPACE_CREDENTIALS_EXPIRED_AT}" > "$STORAGE_PATH/workspace.credentials.json"

TEST_WORKSPACE_CREDENITIALS_USER=$(cat $STORAGE_PATH/workspace.credentials.json | jq -r '.user')
TEST_WORKSPACE_CREDENITIALS_PASS=$(cat $STORAGE_PATH/workspace.credentials.json | jq -r '.password')
echo -n "$TEST_WORKSPACE_CREDENITIALS_USER:$TEST_WORKSPACE_CREDENITIALS_PASS"  > "$STORAGE_PATH/workspace.credentials.plain"
