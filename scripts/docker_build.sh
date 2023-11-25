#!/bin/sh
set -e

docker build --platform linux/amd64 --progress=plain -t kanthorlabs/base:latest docker/base
docker build --platform linux/amd64 --progress=plain -t kanthorlabs/kanthor:latest .