from fastapi import FastAPI, HTTPException

from .parser_core import parse_document, parse_policy_sections, parse_report_facts
from .schemas import (
    ErrorResponse,
    ParseDocumentRequest,
    ParseDocumentResponse,
    ParsePolicyRequest,
    ParsePolicyResponse,
    ParseReportRequest,
    ParseReportResponse,
)

app = FastAPI(title="document-parser", version="0.1.0")


@app.get("/health")
def health() -> dict:
    return {"status": "ok"}


@app.post("/parse/document", response_model=ParseDocumentResponse, responses={400: {"model": ErrorResponse}})
def parse_document_api(request: ParseDocumentRequest) -> ParseDocumentResponse:
    output = parse_document(
        raw_text=request.raw_text or "",
        content_base64=request.content_base64 or "",
        mime_type=request.mime_type,
        enable_ocr=request.enable_ocr,
    )
    if not output.text.strip() and len(output.paragraphs) == 0:
        raise HTTPException(status_code=400, detail={"error": "parse failed", "hints": output.hints})
    return ParseDocumentResponse(
        text=output.text,
        paragraphs=output.paragraphs,
        quality_score=output.quality_score,
        hints=output.hints,
    )


@app.post("/parse/report", response_model=ParseReportResponse)
def parse_report_api(request: ParseReportRequest) -> ParseReportResponse:
    facts, quality_score, hints = parse_report_facts(request.text)
    return ParseReportResponse(facts=facts, quality_score=quality_score, hints=hints)


@app.post("/parse/policy", response_model=ParsePolicyResponse)
def parse_policy_api(request: ParsePolicyRequest) -> ParsePolicyResponse:
    sections, quality_score, hints = parse_policy_sections(request.text)
    return ParsePolicyResponse(sections=sections, quality_score=quality_score, hints=hints)
