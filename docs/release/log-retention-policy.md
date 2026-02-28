# Log Retention Policy

- API and worker logs are retained for 30 days by default.
- Sensitive fields are masked in request logs:
  - phone numbers
  - storage keys / file paths
  - long text fragments
- Request IDs are included in every response and log line.
- Production deployment should set log rotation policy via platform (e.g. Loki/ELK/Cloud logging).
