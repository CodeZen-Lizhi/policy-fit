import base64

from app.parser_core import parse_document, parse_policy_sections, parse_report_facts


def test_parse_document_raw_text():
    output = parse_document(raw_text="第一段\n\n第二段")
    assert output.quality_score > 0
    assert len(output.paragraphs) == 2


def test_parse_document_empty_payload():
    output = parse_document(raw_text="")
    assert output.quality_score == 0
    assert "Empty payload" in output.hints


def test_parse_report_facts():
    facts, quality, hints = parse_report_facts("血压 155/95，血糖 7.2")
    assert len(facts) >= 1
    assert quality > 0
    assert isinstance(hints, list)


def test_parse_policy_sections():
    sections, quality, _ = parse_policy_sections("第1条 既往症定义\n第2条 等待期说明")
    assert len(sections) >= 1
    assert quality > 0


def test_parse_document_base64_text_payload():
    raw = "sample text payload".encode("utf-8")
    b64 = base64.b64encode(raw).decode("utf-8")
    output = parse_document(content_base64=b64, mime_type="text/plain")
    assert "sample text payload" in output.text
