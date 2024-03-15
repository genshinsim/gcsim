import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { merge } from "lodash-es";
import { initReactI18next } from "react-i18next";
import Chinese from "./locales/Chinese.json";
import English from "./locales/English.json";
import German from "./locales/German.json";
import Japanese from "./locales/Japanese.json";
import Korean from "./locales/Korean.json";
import names from "./locales/names.generated.json";
import traveler_names from "./locales/names.traveler.json";
import Russian from "./locales/Russian.json";
import Spanish from "./locales/Spanish.json";

const resources = {
  en: {
    translation: English,
    game: merge(names.English, traveler_names.English),
  },
  zh: {
    translation: Chinese,
    game: merge(names.Chinese, traveler_names.Chinese),
  },
  ja: {
    translation: Japanese,
    game: merge(names.Japanese, traveler_names.Japanese),
  },
  ko: {
    translation: Korean,
    game: merge(names.Korean, traveler_names.Korean),
  },
  es: {
    translation: Spanish,
    game: merge(names.Spanish, traveler_names.Spanish),
  },
  ru: {
    translation: Russian,
    game: merge(names.Russian, traveler_names.Russian),
  },
  de: {
    translation: German,
    game: merge(names.German, traveler_names.German),
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
