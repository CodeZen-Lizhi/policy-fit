# API

Base path: `/api/v1`

All `/api/v1/*` endpoints require `Authorization: Bearer <jwt>`.

## Task APIs

### POST `/tasks`
Create task.

Request:
```json
{"request_id":"optional-id"}
```

Response:
```json
{"data":{"task_id":1,"status":"pending"}}
```

### GET `/tasks/{id}`
Get task detail.

### POST `/tasks/{id}/documents`
Upload a PDF document.

Form fields:
- `docType`: `report` | `policy` | `disclosure`
- `file`: PDF file

### POST `/tasks/{id}/run`
Enqueue task for worker processing.

### GET `/tasks/{id}/findings`
Get findings and risk summary.

### GET `/tasks/{id}/export?format=md|pdf&lang=zh-CN|en-US`
Export report in Markdown or PDF, with optional language.

### GET `/tasks/compare?from={taskId}&to={taskId}`
Compare two historical tasks.

### DELETE `/tasks/{id}`
Delete task and related data.

## Rule Admin APIs

Admin routes:
- `GET /admin/rules/active`
- `GET /admin/rules/versions?limit=20`
- `POST /admin/rules/publish`
- `POST /admin/rules/rollback`
- `POST /admin/rules/gray`

### POST `/admin/rules/publish`
Request:
```json
{
  "changelog": "adjust thresholds",
  "content": {
    "topics": ["hypertension"],
    "policy_types": ["pre_existing"]
  }
}
```

### POST `/admin/rules/rollback`
Request:
```json
{"version":"v20260228123000"}
```

### POST `/admin/rules/gray`
Request:
```json
{"version":"v20260228123000","enabled":true}
```

## Analytics APIs

### POST `/analytics/events`
Track frontend/user analytics event.

Request:
```json
{
  "event_name":"report_viewed",
  "task_id":123,
  "properties":{"source":"web"}
}
```

### GET `/analytics/funnel?period=week|month|all`
Get funnel counts:
- `task_created`
- `document_uploaded`
- `task_run`
- `task_completed`
- `report_viewed`
- `report_exported`

### GET `/analytics/overview?period=week|month|all`
Get core metrics:
- `task_created`
- `task_completed`
- `report_viewed`
- `report_exported`
- `task_deleted`
- `completion_rate`
- `view_rate`
- `export_rate`
