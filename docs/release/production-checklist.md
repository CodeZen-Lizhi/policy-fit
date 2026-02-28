# Production Checklist

## 1. Environment Configuration (T-0701)

- [x] `.env.example` contains all required keys.
- [x] Added `ADMIN_USER_IDS` for admin route governance.
- [x] Added parser and i18n related runtime options.
- [ ] Fill real production values in `.env.prod` (outside repo secret management).

## 2. Backup Strategy (T-0702)

- [x] Database backup script prepared: `scripts/release/backup_db.sh`
- [x] Object storage backup script prepared: `scripts/release/backup_storage.sh`
- [ ] Connect scripts to production scheduler (cron/Argo/CI).

## 3. Log Retention & Masking (T-0703)

- [x] Sensitive field masking in middleware logs.
- [x] Request ID propagated for traceability.
- [x] Logging policy documented in `docs/release/log-retention-policy.md`.
