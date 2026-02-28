"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";
import enUS from "../locales/en-US.json";
import zhCN from "../locales/zh-CN.json";

const I18nContext = createContext({
  locale: "zh-CN",
  setLocale: () => {},
  t: (key) => key
});

const dictionaries = {
  "zh-CN": zhCN,
  "en-US": enUS
};

export function I18nProvider({ children }) {
  const [locale, setLocale] = useState("zh-CN");

  useEffect(() => {
    const stored = window.localStorage.getItem("locale");
    if (stored && dictionaries[stored]) {
      setLocale(stored);
    }
  }, []);

  const value = useMemo(() => {
    const dict = dictionaries[locale] || dictionaries["zh-CN"];
    return {
      locale,
      setLocale: (next) => {
        if (!dictionaries[next]) {
          return;
        }
        setLocale(next);
        window.localStorage.setItem("locale", next);
      },
      t: (key) => dict[key] || key
    };
  }, [locale]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  return useContext(I18nContext);
}
