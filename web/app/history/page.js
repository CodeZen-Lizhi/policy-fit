"use client";

import { useI18n } from "../../components/i18n-provider";

export default function HistoryPage() {
  const { t } = useI18n();

  return (
    <main style={{ maxWidth: 960, margin: "24px auto", padding: 16 }}>
      <h1>{t("history_title")}</h1>
      <table style={{ width: "100%", borderCollapse: "collapse", background: "#fff" }}>
        <thead>
          <tr>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #e2e8f0" }}>{t("history_task_id")}</th>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #e2e8f0" }}>{t("history_status")}</th>
            <th style={{ textAlign: "left", padding: 8, borderBottom: "1px solid #e2e8f0" }}>{t("history_created_at")}</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td style={{ padding: 8, borderBottom: "1px solid #f1f5f9" }}>1</td>
            <td style={{ padding: 8, borderBottom: "1px solid #f1f5f9" }}>success</td>
            <td style={{ padding: 8, borderBottom: "1px solid #f1f5f9" }}>2026-02-28 22:00</td>
          </tr>
        </tbody>
      </table>
    </main>
  );
}
