# Operations Runbook

## Restart services
- API: restart process `cmd/api`
- Worker: restart process `cmd/worker`

## Queue backlog handling
1. Check queue depth in Redis key `analysis_tasks`
2. Scale worker concurrency by `WORKER_CONCURRENCY`
3. Check dead-letter queue `analysis_tasks_dead_letter`

## Failure rate alert (>15%)
1. Check recent worker logs for parser/llm/rule failures
2. Verify database and redis readiness (`/ready`)
3. If persistent, disable run endpoint temporarily and drain queue

## Migration rollback
- One step down: `make migrate-down`
- Re-apply: `make migrate-up`
