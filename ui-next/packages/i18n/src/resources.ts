import Chinese from "./locales/Chinese.json";
import English from "./locales/English.json";
import German from "./locales/German.json";
import Japanese from "./locales/Japanese.json";
import Korean from "./locales/Korean.json";
import names from "./locales/names.generated.json";
import travelerNames from "./locales/names.traveler.json";
import Russian from "./locales/Russian.json";
import Spanish from "./locales/Spanish.json";

export const resources = {
  en: {
    translation: English,
    game: { ...names.English, ...travelerNames.English },
  },
  zh: {
    translation: Chinese,
    game: { ...names.Chinese, ...travelerNames.Chinese },
  },
  ja: {
    translation: Japanese,
    game: { ...names.Japanese, ...travelerNames.Japanese },
  },
  ko: {
    translation: Korean,
    game: { ...names.Korean, ...travelerNames.Korean },
  },
  es: {
    translation: Spanish,
    game: { ...names.Spanish, ...travelerNames.Spanish },
  },
  ru: {
    translation: Russian,
    game: { ...names.Russian, ...travelerNames.Russian },
  },
  de: {
    translation: German,
    game: { ...names.German, ...travelerNames.German },
  },
};

export const specialLocales = ["zh", "ja", "ko"];
