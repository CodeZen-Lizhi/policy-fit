"use client";

import { useEffect, useState } from "react";
import { apiRequest } from "../../../lib/api";

const defaultContent = JSON.stringify(
  {
    topics: ["hypertension", "blood-sugar", "liver"],
    policy_types: ["pre_existing", "waiting_period"]
  },
  null,
  2
);

export default function AdminRulesPage() {
  const [active, setActive] = useState(null);
  const [versions, setVersions] = useState([]);
  const [changelog, setChangelog] = useState("");
  const [contentText, setContentText] = useState(defaultContent);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    refresh();
  }, []);

  async function refresh() {
    try {
      setLoading(true);
      setError("");
      const [activeData, listData] = await Promise.all([
        apiRequest("/admin/rules/active"),
        apiRequest("/admin/rules/versions?limit=20")
      ]);
      setActive(activeData);
      setVersions(listData.items || []);
    } catch (err) {
      setError(err.message || "加载失败");
    } finally {
      setLoading(false);
    }
  }

  async function publish() {
    try {
      const content = JSON.parse(contentText);
      if (!changelog.trim()) {
        setError("请填写 changelog");
        return;
      }
      setLoading(true);
      setError("");
      await apiRequest("/admin/rules/publish", {
        method: "POST",
        body: { changelog: changelog.trim(), content }
      });
      setChangelog("");
      await refresh();
    } catch (err) {
      setError(err.message || "发布失败");
    } finally {
      setLoading(false);
    }
  }

  async function rollback(version) {
    try {
      setLoading(true);
      setError("");
      await apiRequest("/admin/rules/rollback", {
        method: "POST",
        body: { version }
      });
      await refresh();
    } catch (err) {
      setError(err.message || "回滚失败");
    } finally {
      setLoading(false);
    }
  }

  async function toggleGray(version, enabled) {
    try {
      setLoading(true);
      setError("");
      await apiRequest("/admin/rules/gray", {
        method: "POST",
        body: { version, enabled }
      });
      await refresh();
    } catch (err) {
      setError(err.message || "灰度设置失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={{ maxWidth: 1080, margin: "24px auto", padding: 16 }}>
      <h1>规则管理（管理员）</h1>
      {error ? <p style={{ color: "#dc2626" }}>{error}</p> : null}
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 16 }}>
        <h2>当前生效版本</h2>
        {active ? (
          <p>
            <strong>{active.version}</strong> · {active.changelog}
          </p>
        ) : (
          <p style={{ color: "#64748b" }}>暂无生效版本</p>
        )}
      </section>

      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 16 }}>
        <h2>发布新版本</h2>
        <div style={{ display: "grid", gap: 10 }}>
          <input
            placeholder="changelog"
            value={changelog}
            onChange={(e) => setChangelog(e.target.value)}
            style={{ padding: 8, borderRadius: 8, border: "1px solid #cbd5e1" }}
          />
          <textarea
            rows={8}
            value={contentText}
            onChange={(e) => setContentText(e.target.value)}
            style={{ padding: 8, borderRadius: 8, border: "1px solid #cbd5e1", fontFamily: "monospace" }}
          />
          <button
            onClick={publish}
            disabled={loading}
            style={{ width: 140, padding: "10px 12px", borderRadius: 8, border: 0, background: "#0ea5e9", color: "#fff" }}
          >
            发布版本
          </button>
        </div>
      </section>

      <section style={{ background: "#fff", borderRadius: 12, padding: 16 }}>
        <h2>历史版本</h2>
        {loading ? <p>加载中...</p> : null}
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead>
            <tr>
              <th style={thStyle}>版本</th>
              <th style={thStyle}>变更说明</th>
              <th style={thStyle}>状态</th>
              <th style={thStyle}>操作</th>
            </tr>
          </thead>
          <tbody>
            {versions.map((item) => (
              <tr key={item.version}>
                <td style={tdStyle}>{item.version}</td>
                <td style={tdStyle}>{item.changelog}</td>
                <td style={tdStyle}>
                  {item.is_active ? "生效中" : "未生效"}
                  {item.is_gray ? " / 灰度" : ""}
                </td>
                <td style={tdStyle}>
                  <button style={actionBtn} onClick={() => rollback(item.version)}>
                    回滚到此
                  </button>
                  <button style={actionBtn} onClick={() => toggleGray(item.version, !item.is_gray)}>
                    {item.is_gray ? "关闭灰度" : "开启灰度"}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </main>
  );
}

const thStyle = { textAlign: "left", padding: 8, borderBottom: "1px solid #e2e8f0" };
const tdStyle = { padding: 8, borderBottom: "1px solid #f1f5f9", verticalAlign: "top" };
const actionBtn = {
  marginRight: 8,
  padding: "6px 10px",
  borderRadius: 6,
  border: "1px solid #cbd5e1",
  background: "#fff"
};
