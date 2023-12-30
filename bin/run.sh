#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

schemadoc gen --in ./schemas --out ./out -vv
markdownlint --fix out/

echo "[run] ✅"
