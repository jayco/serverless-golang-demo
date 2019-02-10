#!/bin/bash
set -eu

touch go.mod

PROJECT_NAME=$(basename $(pwd | xargs dirname))
CURRENT_DIR=$(basename $(pwd))

CONTENT=$(cat <<-EOD
module github.com/${PROJECT_NAME}/${CURRENT_DIR}

EOD
)

echo "$CONTENT" > go.mod
