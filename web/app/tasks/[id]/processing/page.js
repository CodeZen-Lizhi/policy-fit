"use client";

import { useEffect, useState } from "react";

const phases = ["pending", "parsing", "extracting", "matching", "success"];

export default function ProcessingPage() {
  const [idx, setIdx] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setIdx((prev) => (prev < phases.length - 1 ? prev + 1 : prev));
    }, 1500);
    return () => clearInterval(timer);
  }, []);

  return (
    <main style={{ maxWidth: 760, margin: "24px auto", padding: 16 }}>
      <h1>任务处理中</h1>
      <p>当前阶段：{phases[idx]}</p>
      <progress value={(idx + 1) * 20} max={100} style={{ width: "100%" }} />
      <p style={{ color: "#64748b" }}>轮询策略：前 30s 每 3s 一次，之后每 10s 一次。</p>
    </main>
  );
}
