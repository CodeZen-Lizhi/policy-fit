from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"


def test_parse_document_with_raw_text():
    response = client.post(
        "/parse/document",
        json={"raw_text": "第1段\n\n第2段", "mime_type": "text/plain"},
    )
    assert response.status_code == 200
    payload = response.json()
    assert payload["quality_score"] > 0
    assert len(payload["paragraphs"]) == 2


def test_parse_report():
    response = client.post("/parse/report", json={"text": "血压 150/95"})
    assert response.status_code == 200
    assert "facts" in response.json()


def test_parse_policy():
    response = client.post("/parse/policy", json={"text": "既往症定义\n等待期 90 天"})
    assert response.status_code == 200
    assert "sections" in response.json()
