"use client";

import { useI18n } from "./i18n-provider";

export default function LanguageSwitcher() {
  const { locale, setLocale, t } = useI18n();

  return (
    <div style={{ display: "flex", gap: 8 }}>
      <button
        onClick={() => setLocale("zh-CN")}
        style={{
          border: locale === "zh-CN" ? "1px solid #0369a1" : "1px solid #cbd5e1",
          borderRadius: 6,
          padding: "6px 10px",
          background: locale === "zh-CN" ? "#e0f2fe" : "#fff"
        }}
      >
        {t("lang_zh")}
      </button>
      <button
        onClick={() => setLocale("en-US")}
        style={{
          border: locale === "en-US" ? "1px solid #0369a1" : "1px solid #cbd5e1",
          borderRadius: 6,
          padding: "6px 10px",
          background: locale === "en-US" ? "#e0f2fe" : "#fff"
        }}
      >
        {t("lang_en")}
      </button>
    </div>
  );
}
