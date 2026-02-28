# Rollback Drill Report (2026-02-28)

## Targets

- API rollback
- Worker rollback
- Rule version rollback
- Database rollback (test environment)

## Executed Commands / Procedures

1. API rollback procedure script path: `scripts/release/rollback.sh api`.
2. Worker rollback procedure script path: `scripts/release/rollback.sh worker`.
3. Rule rollback procedure uses admin API: `POST /api/v1/admin/rules/rollback`.
4. DB rollback procedure in test environment:
   - `make migrate-down`
   - `make migrate-up`

## Outcome

- Rollback playbooks are executable and documented.
- Rule rollback endpoint available and covered by handler tests.
- DB migration rollback/restore path available.

## Follow-up

- Execute the same drill on staging with deployment platform hooks.
