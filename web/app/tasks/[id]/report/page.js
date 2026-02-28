"use client";

import Link from "next/link";
import { useEffect } from "react";
import { trackEvent } from "../../../../lib/analytics";

export default function ReportPage() {
  useEffect(() => {
    trackEvent("report_viewed", { properties: { source: "web_report_page" } });
  }, []);

  return (
    <main style={{ maxWidth: 960, margin: "24px auto", padding: 16 }}>
      <h1>风险总览</h1>
      <div style={{ display: "flex", gap: 12, marginBottom: 16 }}>
        <Card color="#ef4444" title="高风险" value="2" />
        <Card color="#f59e0b" title="待确认" value="3" />
        <Card color="#22c55e" title="暂无冲突" value="5" />
      </div>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16 }}>
        <h2>风险列表</h2>
        <p>🔴 高血压风险 · 可能触发既往症告知冲突</p>
        <p>体检证据：血压 155/95 mmHg（2026-01-10）</p>
        <p>条款证据：既往症定义第 3 条</p>
        <Link href="/tasks/1/findings/hypertension">查看详情</Link>
      </section>
      <p style={{ marginTop: 16, fontSize: 13, color: "#475569" }}>
        免责声明：本报告仅为辅助解读，不构成理赔结论。
      </p>
    </main>
  );
}

function Card({ color, title, value }) {
  return (
    <div style={{ flex: 1, background: "#fff", borderRadius: 12, padding: 12, borderTop: `4px solid ${color}` }}>
      <div style={{ color }}>{title}</div>
      <strong style={{ fontSize: 28 }}>{value}</strong>
    </div>
  );
}
