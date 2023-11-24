#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

sh "$SCRIPT_DIR/gen_swagger.sh"
sh "$SCRIPT_DIR/gen_ioc.sh"
sh "$SCRIPT_DIR/ci_vet.sh"

find . -type f -name 'checksum' -exec git add {} \;

NOW=$(date +%Y.%-m%d.%-H%M)

git commit --allow-empty -m "ci($NOW): ✨🐛🚨"

TARGET=${1:-origin}
echo "---------------------------"
printf "Pushing... $NOW --> %s\n" "$TARGET"
echo "---------------------------"
git push "$TARGET"