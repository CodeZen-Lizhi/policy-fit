"use client";

import Link from "next/link";
import LanguageSwitcher from "./language-switcher";
import { useI18n } from "./i18n-provider";

export default function TopNav() {
  const { t } = useI18n();

  return (
    <header
      style={{
        position: "sticky",
        top: 0,
        zIndex: 10,
        background: "#ffffffdd",
        backdropFilter: "blur(6px)",
        borderBottom: "1px solid #e2e8f0"
      }}
    >
      <div style={{ maxWidth: 1080, margin: "0 auto", padding: "10px 16px", display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12 }}>
        <nav style={{ display: "flex", gap: 14, flexWrap: "wrap" }}>
          <Link href="/">{t("nav_home")}</Link>
          <Link href="/tasks/new">{t("nav_new_task")}</Link>
          <Link href="/history">{t("nav_history")}</Link>
          <Link href="/admin/rules">{t("nav_admin_rules")}</Link>
          <Link href="/admin/analytics">{t("nav_admin_analytics")}</Link>
        </nav>
        <LanguageSwitcher />
      </div>
    </header>
  );
}
