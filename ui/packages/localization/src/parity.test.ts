import { describe, expect, it } from "vitest";
import { resources } from "./index";

type LocaleKey = keyof typeof resources;

// zh/ja/ko still carry untranslated `translation` copy; folding them in is
// tracked in genshinsim/gcsim#2934. (This is unrelated to `specialLocales`,
// which is a CJK chart-axis rendering flag, not a translation-status list.)
const deferredLocales: LocaleKey[] = ["zh", "ja", "ko"];

function flattenKeys(obj: Record<string, unknown>, prefix = ""): string[] {
	const keys: string[] = [];
	for (const [k, v] of Object.entries(obj)) {
		const key = prefix ? `${prefix}.${k}` : k;
		if (v !== null && typeof v === "object" && !Array.isArray(v)) {
			keys.push(...flattenKeys(v as Record<string, unknown>, key));
		} else {
			keys.push(key);
		}
	}
	return keys;
}

const englishKeys = new Set(flattenKeys(resources.en.translation));

const checkedLocales = (Object.keys(resources) as LocaleKey[]).filter(
	(lng) => lng !== "en" && !deferredLocales.includes(lng),
);

describe("translation locale key parity with English", () => {
	it.each(checkedLocales)("%s has exactly the English keys", (lng) => {
		const localeKeys = new Set(flattenKeys(resources[lng].translation));
		const missing = [...englishKeys].filter((k) => !localeKeys.has(k)).sort();
		const orphan = [...localeKeys].filter((k) => !englishKeys.has(k)).sort();
		expect({ missing, orphan }).toEqual({ missing: [], orphan: [] });
	});
});
