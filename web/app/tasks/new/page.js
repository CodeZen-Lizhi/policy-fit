"use client";

import { useState } from "react";
import { trackEvent } from "../../../lib/analytics";
import { useI18n } from "../../../components/i18n-provider";

export default function NewTaskPage() {
  const { t } = useI18n();
  const [report, setReport] = useState(null);
  const [policy, setPolicy] = useState(null);
  const [disclosure, setDisclosure] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function onStartAnalyze() {
    setIsSubmitting(true);
    await trackEvent("frontend_start_clicked", {
      properties: {
        has_report: !!report,
        has_policy: !!policy,
        has_disclosure: !!disclosure
      }
    });
    setIsSubmitting(false);
  }

  return (
    <main style={{ maxWidth: 960, margin: "24px auto", padding: 16 }}>
      <h1>{t("new_task_title")}</h1>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, boxShadow: "0 8px 24px rgba(15,23,42,0.06)" }}>
        <Upload label={t("new_task_report")} onChange={setReport} />
        <Upload label={t("new_task_policy")} onChange={setPolicy} />
        <Upload label={t("new_task_disclosure")} onChange={setDisclosure} />
        <button
          style={{ marginTop: 12, padding: "10px 14px", borderRadius: 8, border: 0, background: "#0ea5e9", color: "#fff" }}
          disabled={!report || !policy}
          onClick={onStartAnalyze}
        >
          {isSubmitting ? t("new_task_submitting") : t("new_task_start")}
        </button>
      </section>
      <p style={{ marginTop: 16, fontSize: 13, color: "#475569" }}>
        {t("disclaimer")}
      </p>
    </main>
  );
}

function Upload({ label, onChange }) {
  return (
    <label style={{ display: "block", marginBottom: 12 }}>
      <div style={{ marginBottom: 6 }}>{label}</div>
      <input type="file" accept="application/pdf" onChange={(e) => onChange(e.target.files?.[0] || null)} />
    </label>
  );
}
