"use client";

import { useEffect } from "react";
import { trackEvent } from "../../../lib/analytics";

export default function ComparePage() {
  useEffect(() => {
    trackEvent("history_compare_viewed", { properties: { source: "web_compare_page" } });
  }, []);

  return (
    <main style={{ maxWidth: 960, margin: "24px auto", padding: 16 }}>
      <h1>历史报告对比</h1>
      <p>新增风险、等级变化与消失风险对比视图。</p>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16 }}>
        <h2>差异摘要</h2>
        <ul>
          <li>新增风险：1</li>
          <li>等级变化：2</li>
          <li>消失风险：1</li>
        </ul>
      </section>
    </main>
  );
}
