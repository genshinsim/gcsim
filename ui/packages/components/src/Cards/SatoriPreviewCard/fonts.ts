// DejaVu Sans Mono is bundled with this package (src/Cards/SatoriPreviewCard/
// fonts). It stands in for the CSS `font-mono` stack the live PreviewCard renders
// with (ui-monospace, SFMono-Regular, Menlo, …, "Liberation Mono", "Courier New",
// monospace) — on the headless-Chrome preview that resolves to a system monospace,
// so the Satori card bundles one to stay deterministic with no system fonts.
//
// The card tree uses `fontFamily: FONT_FAMILY`; the caller passes the loaded
// buffers via Satori's `fonts` option, e.g. `satori(<SatoriPreviewCard .../>,
// { fonts })`.

// Font family used across the card. The bundled family comes first so Satori
// (which matches by the loaded font `name`) picks it; the generic `monospace`
// fallback lets the browser (Storybook, which loads no font file) render the
// card in the same default monospace the live PreviewCard's `font-mono` uses.
export const FONT_FAMILY = "DejaVu Sans Mono, monospace";

// The Satori font `name` — must match the leading family in FONT_FAMILY.
export const FONT_NAME = "DejaVu Sans Mono";

// Shape of a single Satori font entry (kept local so this package needs no
// satori dependency).
export type SatoriFont = {
	name: string;
	data: ArrayBuffer | Buffer;
	weight: 400 | 700;
	style: "normal";
};

// URLs of the bundled TTF files, resolved relative to this module.
export const monoFontUrls = {
	regular: new URL("./fonts/DejaVuSansMono-Regular.ttf", import.meta.url),
	bold: new URL("./fonts/DejaVuSansMono-Bold.ttf", import.meta.url),
};

// Node helper: read the bundled TTFs and return them as Satori font entries.
// `node:fs` is imported dynamically so this module stays safe to bundle for the
// browser (Storybook renders the card without calling this).
export async function loadCardFonts(): Promise<SatoriFont[]> {
	const { readFile } = await import("node:fs/promises");
	const [regular, bold] = await Promise.all([
		readFile(monoFontUrls.regular),
		readFile(monoFontUrls.bold),
	]);
	return [
		{ name: FONT_NAME, data: regular, weight: 400, style: "normal" },
		{ name: FONT_NAME, data: bold, weight: 700, style: "normal" },
	];
}
