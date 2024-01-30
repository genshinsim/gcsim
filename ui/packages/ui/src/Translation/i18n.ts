import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { initReactI18next } from "react-i18next";
import Chinese from "./locales/Chinese.json";
import English from "./locales/English.json";
import IngameNames from "./locales/IngameNames.json";
import Japanese from "./locales/Japanese.json";
import Korean from "./locales/Korean.json";
import Russian from "./locales/Russian.json";
import Spanish from "./locales/Spanish.json";
import German from "./locales/German.json";

const resources = {
  en: {
    translation: English,
    game: IngameNames.English,
  },
  zh: {
    translation: Chinese,
    game: IngameNames.Chinese,
  },
  ja: {
    translation: Japanese,
    game: IngameNames.Japanese,
  },
  ko: {
    translation: Korean,
    game: IngameNames.Korean,
  },
  es: {
    translation: Spanish,
    game: IngameNames.Spanish,
  },
  ru: {
    translation: Russian,
    game: IngameNames.Russian,
  },
  de: {
    translation: German,
    game: IngameNames.German,
  },
};

// used in result tab for graph y axis text direction/offset handling
export const specialLocales = ["zh", "ja", "ko"];

i18n
  .use(initReactI18next)
  .use(LanguageDetector)
  .init({
    resources,
    defaultNS: "translation",
    fallbackLng: "en",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });
