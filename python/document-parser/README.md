# document-parser

Python OCR/Document parsing service for Policy Fit.

## Features

- `POST /parse/document`: OCR-aware plain text extraction for PDF/image payloads.
- `POST /parse/report`: extract health-report facts from text.
- `POST /parse/policy`: extract policy clauses and metadata from text.
- Returns `quality_score` and `hints` for low-quality parsing output.

## Quick start

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8081
```

## API examples

### Parse document

```bash
curl -X POST http://localhost:8081/parse/document \
  -H 'Content-Type: application/json' \
  -d '{"filename":"sample.pdf","mime_type":"application/pdf","content_base64":"..."}'
```

### Parse report

```bash
curl -X POST http://localhost:8081/parse/report \
  -H 'Content-Type: application/json' \
  -d '{"text":"血压 150/95，建议复查"}'
```

### Parse policy

```bash
curl -X POST http://localhost:8081/parse/policy \
  -H 'Content-Type: application/json' \
  -d '{"text":"第1条 既往症定义..."}'
```

## Test

```bash
pytest -q
python scripts/benchmark.py --url http://localhost:8081 --requests 200
```
