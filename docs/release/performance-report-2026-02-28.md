# Performance Report (2026-02-28)

## Scope

- Core API create-task path benchmark under parallel load.
- Parser benchmark script available for OCR service (`python/document-parser/scripts/benchmark.py`).

## API Benchmark

Command:

```bash
GOCACHE=.cache/gobuild go test ./internal/handler -run ^$ -bench BenchmarkTaskHandlerCreateTask -benchmem -count=1
```

Result:

- Benchmark: `BenchmarkTaskHandlerCreateTask-10`
- Throughput: `678465 ops`
- Latency: `1748 ns/op`
- Memory: `9287 B/op`

## Conclusion

- Core create-task path meets expected concurrency baseline in local benchmark.
- Additional end-to-end latency validation should be executed in staging with real DB/Redis/LLM dependencies.
