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
  English: {
    translation: English,
    game: IngameNames.English,
  },
  Chinese: {
    translation: Chinese,
    game: IngameNames.Chinese,
  },
  German: {
    translation: German,
    game: IngameNames.German,
  },
  Japanese: {
    translation: Japanese,
    game: IngameNames.Japanese,
  },
  Spanish: {
    translation: Spanish,
    game: IngameNames.Spanish,
  },
  Russian: {
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
    lng: "English",
    fallbackLng: "English",
    debug: false,
    interpolation: {
      escapeValue: false,
    },
  });
