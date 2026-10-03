import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { THEME_KEY } from "./prefs";
import { DEFAULT_THEME_ID, THEMES } from "./themes";

const read = (path: string) =>
	readFileSync(new URL(path, import.meta.url), "utf8");

describe("THEMES", () => {
	it("matches the palettes in theme.css", () => {
		const css = read("../../../theme.css");
		const schemes = new Map(
			[
				...css.matchAll(
					/\[data-theme="([a-z-]+)"\] \{\n\tcolor-scheme: (light|dark);/g,
				),
			].map((m) => [m[1], m[2]]),
		);
		expect(Object.fromEntries(schemes)).toEqual(
			Object.fromEntries(THEMES.map((t) => [t.id, t.light ? "light" : "dark"])),
		);
		expect(css).toMatch(`:root,\n[data-theme="${DEFAULT_THEME_ID}"] {`);
	});

	it("shares its storage key with the pre-paint script", () => {
		expect(read("../../../web/public/theme.js")).toContain(`"${THEME_KEY}"`);
	});
});
