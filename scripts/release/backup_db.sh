#!/usr/bin/env bash
set -euo pipefail

: "${DB_HOST:?DB_HOST is required}"
: "${DB_PORT:?DB_PORT is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_NAME:?DB_NAME is required}"
: "${PGPASSWORD:?PGPASSWORD is required}"

BACKUP_DIR=${BACKUP_DIR:-./backups/db}
mkdir -p "${BACKUP_DIR}"

TS=$(date +%Y%m%d-%H%M%S)
FILE="${BACKUP_DIR}/policyfit-${TS}.sql.gz"

pg_dump -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" "${DB_NAME}" | gzip > "${FILE}"
echo "db backup created: ${FILE}"
