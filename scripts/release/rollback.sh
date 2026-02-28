#!/usr/bin/env bash
set -euo pipefail

TARGET=${1:-}
if [ -z "$TARGET" ]; then
  echo "usage: $0 <api|worker|rule|db>" >&2
  exit 1
fi

case "$TARGET" in
  api)
    echo "rollback api deployment to previous stable image tag"
    ;;
  worker)
    echo "rollback worker deployment to previous stable image tag"
    ;;
  rule)
    echo "call /api/v1/admin/rules/rollback with previous version"
    ;;
  db)
    echo "run migration rollback in test env: make migrate-down"
    ;;
  *)
    echo "unknown rollback target: $TARGET" >&2
    exit 1
    ;;
esac
