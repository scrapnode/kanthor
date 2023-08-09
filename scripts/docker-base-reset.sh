#!/usr/bin/env bash
set -e

docker compose down && docker volume prune -f
docker compose up -d streaming cache warehouse migration