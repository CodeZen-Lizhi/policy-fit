#!/usr/bin/env bash
set -euo pipefail

: "${STORAGE_TYPE:?STORAGE_TYPE is required}"

BACKUP_DIR=${BACKUP_DIR:-./backups/storage}
mkdir -p "${BACKUP_DIR}"
TS=$(date +%Y%m%d-%H%M%S)

if [ "${STORAGE_TYPE}" = "local" ]; then
  : "${STORAGE_PATH:?STORAGE_PATH is required for local storage}"
  tar -czf "${BACKUP_DIR}/local-storage-${TS}.tar.gz" -C "${STORAGE_PATH}" .
  echo "local storage backup created: ${BACKUP_DIR}/local-storage-${TS}.tar.gz"
  exit 0
fi

if [ "${STORAGE_TYPE}" = "s3" ]; then
  : "${S3_ENDPOINT:?S3_ENDPOINT is required}"
  : "${S3_BUCKET:?S3_BUCKET is required}"
  : "${S3_ACCESS_KEY:?S3_ACCESS_KEY is required}"
  : "${S3_SECRET_KEY:?S3_SECRET_KEY is required}"
  echo "for s3 backup, configure your object storage sync command (aws s3 sync / mc mirror)"
  exit 0
fi

echo "unsupported STORAGE_TYPE=${STORAGE_TYPE}" >&2
exit 1
