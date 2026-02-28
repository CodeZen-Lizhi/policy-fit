# Gray Release Drill Report (2026-02-28)

## Scope

- Target traffic: 10%
- Services: API + Worker
- Observed metrics:
  - task success rate
  - API p95 latency
  - worker fail ratio

## Steps

1. Execute `scripts/release/gray_release.sh 10 v1.0.0-rc1`.
2. Verify core metrics collection endpoint `/metrics` and analytics dashboard endpoint.
3. Simulate 24h check by replaying sample workload and collecting hourly snapshots.
4. Validate no regression in `make test` and no critical alert triggered.
5. Promote to full release by rerunning script with `100`.

## Result

- Gray flow script validated.
- Metrics snapshot path confirmed.
- Promotion path documented and executable.

## Notes

- Real production 24h observation must be executed with real traffic and alerting platform.
