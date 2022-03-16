import LanguageDetector from "i18next-browser-languagedetector";
import i18n from "i18next";
import English from "../public/locales/English.json";
import Chinese from "../public/locales/Chinese.json";
import Japanese from "../public/locales/Japanese.json";
import Spanish from "../public/locales/Spanish.json";
import { initReactI18next } from "react-i18next";

i18n
  .use(initReactI18next)
  .use(LanguageDetector)
  .init({
    resources: {
      English: {
        translation: English,
      },
      Chinese: {
        translation: Chinese,
      },
      Japanese: {
        translation: Japanese,
      },
      Spanish: {
        translation: Spanish,
      },
    },
    fallbackLng: "English",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });
