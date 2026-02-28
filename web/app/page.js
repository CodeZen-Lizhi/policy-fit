"use client";

import { useI18n } from "../components/i18n-provider";

export default function HomePage() {
  const { t } = useI18n();

  return (
    <main style={{ maxWidth: 960, margin: "40px auto", padding: 16 }}>
      <h1>{t("home_title")}</h1>
      <p>{t("home_desc")}</p>
    </main>
  );
}
