"use client";

import { useEffect, useState } from "react";
import { apiRequest } from "../../../lib/api";

const periods = [
  { label: "近 7 天", value: "week" },
  { label: "近 30 天", value: "month" },
  { label: "全量", value: "all" }
];

export default function AdminAnalyticsPage() {
  const [period, setPeriod] = useState("week");
  const [funnel, setFunnel] = useState(null);
  const [overview, setOverview] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadData(period);
  }, [period]);

  async function loadData(p) {
    try {
      setLoading(true);
      setError("");
      const [funnelData, overviewData] = await Promise.all([
        apiRequest(`/analytics/funnel?period=${p}`),
        apiRequest(`/analytics/overview?period=${p}`)
      ]);
      setFunnel(funnelData.funnel || {});
      setOverview(overviewData);
    } catch (err) {
      setError(err.message || "加载失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={{ maxWidth: 1080, margin: "24px auto", padding: 16 }}>
      <h1>运营埋点看板（管理员）</h1>
      {error ? <p style={{ color: "#dc2626" }}>{error}</p> : null}
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 16 }}>
        <h2>周期</h2>
        <div style={{ display: "flex", gap: 8 }}>
          {periods.map((item) => (
            <button
              key={item.value}
              onClick={() => setPeriod(item.value)}
              style={{
                padding: "8px 12px",
                borderRadius: 8,
                border: item.value === period ? "1px solid #0369a1" : "1px solid #cbd5e1",
                background: item.value === period ? "#e0f2fe" : "#fff"
              }}
            >
              {item.label}
            </button>
          ))}
        </div>
      </section>

      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 16 }}>
        <h2>漏斗（任务创建 → 报告查看）</h2>
        {loading ? <p>加载中...</p> : null}
        <ul>
          <li>任务创建：{funnel?.task_created || 0}</li>
          <li>文档上传：{funnel?.document_uploaded || 0}</li>
          <li>开始分析：{funnel?.task_run || 0}</li>
          <li>分析完成：{funnel?.task_completed || 0}</li>
          <li>报告查看：{funnel?.report_viewed || 0}</li>
          <li>报告导出：{funnel?.report_exported || 0}</li>
        </ul>
      </section>

      <section style={{ background: "#fff", borderRadius: 12, padding: 16 }}>
        <h2>核心指标</h2>
        {loading ? <p>加载中...</p> : null}
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, minmax(0, 1fr))", gap: 12 }}>
          <Metric title="任务创建" value={overview?.task_created || 0} />
          <Metric title="分析完成" value={overview?.task_completed || 0} />
          <Metric title="报告查看" value={overview?.report_viewed || 0} />
          <Metric title="完成率" value={toPercent(overview?.completion_rate)} />
          <Metric title="查看率" value={toPercent(overview?.view_rate)} />
          <Metric title="导出率" value={toPercent(overview?.export_rate)} />
        </div>
      </section>
    </main>
  );
}

function Metric({ title, value }) {
  return (
    <div style={{ border: "1px solid #e2e8f0", borderRadius: 10, padding: 12, background: "#f8fafc" }}>
      <div style={{ color: "#475569", fontSize: 13 }}>{title}</div>
      <div style={{ marginTop: 6, fontSize: 22, fontWeight: 700 }}>{value}</div>
    </div>
  );
}

function toPercent(v) {
  if (typeof v !== "number") {
    return "0%";
  }
  return `${(v * 100).toFixed(1)}%`;
}
