import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import { resources } from "./resources";

export function initI18n(lng = "en") {
  return i18next.use(initReactI18next).init({
    resources,
    lng,
    fallbackLng: "en",
    ns: ["translation", "game"],
    defaultNS: "translation",
    interpolation: { escapeValue: false },
  });
}
