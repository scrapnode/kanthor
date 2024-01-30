#!/bin/sh
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )";

sh "$SCRIPT_DIR/gen_swagger.sh"
sh "$SCRIPT_DIR/gen_ioc.sh"
sh "$SCRIPT_DIR/ci_vet.sh"

git add **/checksum

NOW=$(date +%Y.%-m%d.%-H%M)
HAHSTAGS=${HAHSTAGS:-"#sdk"}

git commit --allow-empty -m "ci($NOW): ✨🐛🚨 $HAHSTAGS"

TARGET=${1:-origin}
echo "---------------------------"
printf "Pushing... $NOW --> %s\n" "$TARGET"
echo "---------------------------"
git push "$TARGET"  