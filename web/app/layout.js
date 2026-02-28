import { I18nProvider } from "../components/i18n-provider";
import TopNav from "../components/top-nav";

export const metadata = {
  title: "Policy Fit",
  description: "Policy x Health Fit"
};

export default function RootLayout({ children }) {
  return (
    <html lang="zh-CN">
      <body style={{ margin: 0, fontFamily: "'Noto Sans SC', 'PingFang SC', sans-serif", background: "#f5f7fb", color: "#0f172a" }}>
        <I18nProvider>
          <TopNav />
          {children}
        </I18nProvider>
      </body>
    </html>
  );
}
