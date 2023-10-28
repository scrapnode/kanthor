#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";
/bin/sh "$SCRIPT_DIR/docker_build.sh"

docker compose down || true
docker compose up -d
docker compose ps