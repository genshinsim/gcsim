import LanguageDetector from "i18next-browser-languagedetector";
import i18n from "i18next";
import English from "../public/locales/English.json";
import Chinese from "../public/locales/Chinese.json";
import German from "../public/locales/German.json";
import Japanese from "../public/locales/Japanese.json";
import Spanish from "../public/locales/Spanish.json";
import Russian from "../public/locales/Russian.json";
import IngameNames from "../public/locales/IngameNames.json";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: English,
    game: IngameNames.English,
  },
  zh: {
    translation: Chinese,
    game: IngameNames.Chinese,
  },
  de: {
    translation: German,
    game: IngameNames.German,
  },
  ja: {
    translation: Japanese,
    game: IngameNames.Japanese,
  },
  es: {
    translation: Spanish,
    game: IngameNames.Spanish,
  },
  ru: {
    translation: Russian,
    game: IngameNames.Russian,
  },
};

i18n
  .use(initReactI18next)
  .use(LanguageDetector)
  .init({
    resources,
    defaultNS: "translation",
    lng: "en",
    fallbackLng: "en",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });
