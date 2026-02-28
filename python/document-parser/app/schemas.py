from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class ParseDocumentRequest(BaseModel):
    filename: str = Field(default="unknown")
    mime_type: str = Field(default="application/pdf")
    content_base64: Optional[str] = None
    raw_text: Optional[str] = None
    enable_ocr: bool = Field(default=True)


class Paragraph(BaseModel):
    loc: str
    page: int
    index: int
    text: str


class ParseDocumentResponse(BaseModel):
    text: str
    paragraphs: List[Paragraph]
    quality_score: float
    hints: List[str]


class ParseReportRequest(BaseModel):
    text: str


class HealthFact(BaseModel):
    category: str
    label: str
    value: Optional[str] = None
    confidence: float
    evidence: str


class ParseReportResponse(BaseModel):
    facts: List[HealthFact]
    quality_score: float
    hints: List[str]


class ParsePolicyRequest(BaseModel):
    text: str


class PolicyClause(BaseModel):
    clause_type: str
    title: str
    content: str
    confidence: float


class ParsePolicyResponse(BaseModel):
    sections: List[PolicyClause]
    quality_score: float
    hints: List[str]


class ErrorResponse(BaseModel):
    error: str
    hints: List[str] = Field(default_factory=list)
    detail: Dict[str, Any] = Field(default_factory=dict)
