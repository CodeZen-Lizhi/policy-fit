# Architecture

## Runtime
- API: Gin
- Worker: Redis queue consumer
- DB: PostgreSQL
- Object Storage: Local or S3/MinIO
- Config: Viper
- Logging: Zap

## Main flow
1. Client creates task
2. Client uploads report/policy documents
3. Client triggers run
4. Worker consumes queue and runs parsing/extracting/matching phases
5. Findings persisted and task status updated

## Key packages
- `internal/repository`: data access
- `internal/service`: business logic
- `internal/handler`: HTTP layer
- `internal/jobs`: queue payload + worker
- `internal/parser`: PDF parser
- `internal/llm`: provider and schema validation
- `internal/ruleengine`: matching/scoring
