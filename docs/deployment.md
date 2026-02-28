# Deployment Guide

## 1. Prepare environment

```bash
cp .env.prod.example .env.prod
# fill secrets
```

## 2. Validate configuration

```bash
APP_ENV=prod ENV_FILE=.env.prod make env-check
```

## 3. Start dependencies

```bash
docker-compose up -d
```

## 4. Run migrations

```bash
APP_ENV=prod ENV_FILE=.env.prod make migrate-up
```

## 5. Deploy services

```bash
APP_ENV=prod ENV_FILE=.env.prod make build
APP_ENV=prod ENV_FILE=.env.prod ./bin/api
APP_ENV=prod ENV_FILE=.env.prod ./bin/worker
```

## 6. Post-deploy checks

- `GET /health` should be `ok`
- `GET /ready` should be `ready`
- `GET /metrics` should return pipeline counters

## 7. Backup and recovery scripts

```bash
scripts/release/backup_db.sh
scripts/release/backup_storage.sh
scripts/release/rollback.sh api
scripts/release/rollback.sh worker
scripts/release/rollback.sh rule
scripts/release/rollback.sh db
```

## 8. Gray release

```bash
scripts/release/gray_release.sh 10 v1.0.0-rc1
scripts/release/gray_release.sh 100 v1.0.0
```
