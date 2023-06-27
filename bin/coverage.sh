#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

./bin/test.sh
gocovsh --profile coverage.out
