#!/usr/bin/env bash
set -euo pipefail

# Example gray release helper (platform-agnostic)
PERCENT=${1:-10}
VERSION=${2:-v1.0.0-rc1}

if ! [[ "$PERCENT" =~ ^[0-9]+$ ]]; then
  echo "percent must be numeric" >&2
  exit 1
fi
if [ "$PERCENT" -le 0 ] || [ "$PERCENT" -gt 100 ]; then
  echo "percent must be in (0,100]" >&2
  exit 1
fi

echo "[gray] deploy version=${VERSION} traffic=${PERCENT}%"
echo "[gray] observe key metrics: error_rate, task_success_rate, p95_latency"
echo "[gray] if healthy, gradually increase to 100%"
