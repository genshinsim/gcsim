import LanguageDetector from "i18next-browser-languagedetector";
import i18n from "i18next";
import English from "../../locales/English.json";
import Chinese from "../../locales/Chinese.json";
import German from "../../locales/German.json";
import Japanese from "../../locales/Japanese.json";
import Spanish from "../../locales/Spanish.json";
import Russian from "../../locales/Russian.json";
import IngameNames from "../../locales/IngameNames.json";
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
