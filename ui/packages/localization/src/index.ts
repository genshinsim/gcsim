import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { merge } from "lodash-es";
import { initReactI18next } from "react-i18next";
import Chinese from "./locales/Chinese.json";
import English from "./locales/English.json";
import German from "./locales/German.json";
import Japanese from "./locales/Japanese.json";
import Korean from "./locales/Korean.json";
import names from "./locales/names.dm.json";
import names_override from "./locales/names.override.json";
import Russian from "./locales/Russian.json";
import Spanish from "./locales/Spanish.json";

export const resources = {
	en: {
		translation: English,
		game: merge(names.English, names_override.English),
	},
	zh: {
		translation: Chinese,
		game: merge(names.Chinese, names_override.Chinese),
	},
	ja: {
		translation: Japanese,
		game: merge(names.Japanese, names_override.Japanese),
	},
	ko: {
		translation: Korean,
		game: merge(names.Korean, names_override.Korean),
	},
	es: {
		translation: Spanish,
		game: merge(names.Spanish, names_override.Spanish),
	},
	ru: {
		translation: Russian,
		game: merge(names.Russian, names_override.Russian),
	},
	de: {
		translation: German,
		game: merge(names.German, names_override.German),
	},
};

export const specialLocales = ["zh", "ja", "ko"];

type I18nModule = Parameters<typeof i18n.use>[0];

export function initI18n(options?: {
	use?: I18nModule[];
	init?: Parameters<typeof i18n.init>[0];
}) {
	const instance = i18n.use(initReactI18next).use(LanguageDetector);
	for (const mod of options?.use ?? []) {
		instance.use(mod);
	}
	instance.init({
		resources,
		defaultNS: "translation",
		fallbackLng: "en",
		debug: false,
		interpolation: {
			escapeValue: false,
		},
		...options?.init,
	});
	return i18n;
}
