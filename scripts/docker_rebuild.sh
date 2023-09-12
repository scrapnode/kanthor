#!/bin/sh
set -e

docker compose build --progress=plain --no-cache sdkapi
docker compose build --progress=plain
docker compose stop sdkapi portalapi scheduler dispatcher storage
docker compose rm -f sdkapi portalapi scheduler dispatcher storage
docker compose up -d sdkapi portalapi scheduler dispatcher storage