#!/usr/bin/env python3
import argparse
import statistics
import time

import requests


SAMPLE_TEXT = """
体检报告摘要：
血压 155/95 mmHg，建议复查。
空腹血糖 7.2 mmol/L。
""".strip()


def run(url: str, requests_count: int) -> None:
    latencies = []
    for _ in range(requests_count):
        started = time.time()
        response = requests.post(
            f"{url}/parse/document",
            json={"raw_text": SAMPLE_TEXT, "mime_type": "text/plain"},
            timeout=5,
        )
        response.raise_for_status()
        latencies.append((time.time() - started) * 1000)

    p50 = statistics.median(latencies)
    p95 = sorted(latencies)[int(len(latencies) * 0.95) - 1]
    print(f"requests={requests_count} p50_ms={p50:.2f} p95_ms={p95:.2f}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Simple parser service benchmark")
    parser.add_argument("--url", default="http://localhost:8081", help="service base url")
    parser.add_argument("--requests", type=int, default=200, help="number of requests")
    args = parser.parse_args()
    run(args.url.rstrip("/"), args.requests)
