import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import Backend from "i18next-http-backend";
import { initReactI18next } from "react-i18next";

import { resources } from "@gcsim/localization";

// used in result tab for graph y axis text direction/offset handling
export const specialLocales = ["zh", "ja", "ko"];

i18n
  .use(initReactI18next)
  .use(LanguageDetector)
  .use(Backend)
  .init({
    resources,
    defaultNS: "translation",
    fallbackLng: "en",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
    returnNull: false,
  });

export default i18n;
