import base64
import io
import re
from dataclasses import dataclass
from typing import List, Tuple

from .schemas import HealthFact, Paragraph, PolicyClause


@dataclass
class ParseOutput:
    text: str
    paragraphs: List[Paragraph]
    quality_score: float
    hints: List[str]


def parse_document_bytes(payload: bytes, mime_type: str, enable_ocr: bool) -> ParseOutput:
    text = ""
    hints: List[str] = []

    if mime_type == "application/pdf":
        text = extract_text_from_pdf(payload)
        if len(text.strip()) < 30 and enable_ocr:
            hints.append("PDF text sparse; OCR fallback may be required")
    elif mime_type.startswith("image/"):
        text = extract_text_from_image(payload, enable_ocr)
        if not text.strip():
            hints.append("OCR could not detect enough text")
    else:
        hints.append("Unknown mime_type, treating as plain text")
        text = payload.decode("utf-8", errors="ignore")

    paragraphs = to_paragraphs(text)
    quality_score = calc_quality_score(text, paragraphs)
    hints.extend(quality_hints(quality_score, len(paragraphs)))

    return ParseOutput(text=text, paragraphs=paragraphs, quality_score=quality_score, hints=dedup(hints))


def parse_document_base64(content_base64: str, mime_type: str, enable_ocr: bool) -> ParseOutput:
    raw = base64.b64decode(content_base64)
    return parse_document_bytes(raw, mime_type, enable_ocr)


def extract_text_from_pdf(payload: bytes) -> str:
    try:
        from pypdf import PdfReader
    except Exception:
        return ""

    try:
        reader = PdfReader(io.BytesIO(payload))
        texts = []
        for page in reader.pages:
            texts.append(page.extract_text() or "")
        return "\n\n".join(texts)
    except Exception:
        return ""


def extract_text_from_image(payload: bytes, enable_ocr: bool) -> str:
    if not enable_ocr:
        return ""
    try:
        from PIL import Image
        import pytesseract
    except Exception:
        return ""

    try:
        image = Image.open(io.BytesIO(payload))
        return pytesseract.image_to_string(image, lang="chi_sim+eng")
    except Exception:
        return ""


def to_paragraphs(text: str) -> List[Paragraph]:
    blocks = [b.strip() for b in re.split(r"\n\s*\n", text.replace("\r\n", "\n")) if b.strip()]
    result: List[Paragraph] = []
    for idx, block in enumerate(blocks, start=1):
        result.append(
            Paragraph(
                loc=f"para_{idx}",
                page=1,
                index=idx,
                text=block,
            )
        )
    return result


def calc_quality_score(text: str, paragraphs: List[Paragraph]) -> float:
    length_score = min(len(text.strip()) / 2400.0, 1.0)
    para_score = min(len(paragraphs) / 20.0, 1.0)
    score = (length_score * 0.7) + (para_score * 0.3)
    return round(max(0.0, min(score, 1.0)), 3)


def quality_hints(score: float, paragraph_count: int) -> List[str]:
    hints: List[str] = []
    if score < 0.35:
        hints.append("Low OCR quality, recommend re-uploading a clearer file")
    if paragraph_count == 0:
        hints.append("No readable paragraphs extracted")
    if paragraph_count < 3:
        hints.append("Very few paragraphs found; result may be incomplete")
    return hints


def parse_report_facts(text: str) -> Tuple[List[HealthFact], float, List[str]]:
    facts: List[HealthFact] = []

    bp_match = re.search(r"(\d{2,3})\s*/\s*(\d{2,3})", text)
    if bp_match:
        facts.append(
            HealthFact(
                category="blood_pressure",
                label="blood_pressure",
                value=f"{bp_match.group(1)}/{bp_match.group(2)}",
                confidence=0.82,
                evidence=bp_match.group(0),
            )
        )

    glucose_match = re.search(r"(血糖|glucose)[^\d]*(\d+(?:\.\d+)?)", text, re.IGNORECASE)
    if glucose_match:
        facts.append(
            HealthFact(
                category="blood_glucose",
                label="blood_glucose",
                value=glucose_match.group(2),
                confidence=0.78,
                evidence=glucose_match.group(0),
            )
        )

    hints: List[str] = []
    if not facts:
        hints.append("No structured health facts matched; please check text quality")
    quality = round(min(1.0, 0.45 + len(facts) * 0.2), 3)
    return facts, quality, hints


def parse_policy_sections(text: str) -> Tuple[List[PolicyClause], float, List[str]]:
    clauses: List[PolicyClause] = []
    clause_patterns = [
        ("pre_existing", r"既往症"),
        ("waiting_period", r"等待期"),
        ("exclusion", r"免责|责任免除"),
        ("disclosure", r"告知"),
        ("disease_definition", r"疾病定义|特定疾病"),
    ]

    lines = [line.strip() for line in text.splitlines() if line.strip()]
    for line in lines:
        for clause_type, pattern in clause_patterns:
            if re.search(pattern, line, re.IGNORECASE):
                clauses.append(
                    PolicyClause(
                        clause_type=clause_type,
                        title=line[:24],
                        content=line,
                        confidence=0.75,
                    )
                )
                break

    hints: List[str] = []
    if not clauses:
        hints.append("No known policy clauses matched")
    quality = round(min(1.0, 0.4 + len(clauses) * 0.08), 3)
    return clauses, quality, hints


def parse_document(raw_text: str = "", content_base64: str = "", mime_type: str = "application/pdf", enable_ocr: bool = True) -> ParseOutput:
    if raw_text.strip():
        paragraphs = to_paragraphs(raw_text)
        score = calc_quality_score(raw_text, paragraphs)
        return ParseOutput(raw_text, paragraphs, score, quality_hints(score, len(paragraphs)))
    if content_base64.strip():
        return parse_document_base64(content_base64, mime_type, enable_ocr)
    return ParseOutput("", [], 0.0, ["Empty payload"])


def dedup(items: List[str]) -> List[str]:
    seen = set()
    out = []
    for item in items:
        if item not in seen:
            seen.add(item)
            out.append(item)
    return out
