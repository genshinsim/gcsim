// Node-only font loading for the Satori preview card. Kept in a separate module
// from ./fonts (which the card tree imports) so that ./fonts stays free of any
// `node:*` reference and bundles cleanly for the browser and Cloudflare Workers.
// On Workers, bundle the TTF bytes at build time instead of calling this (see
// @gcsim/workers src/preview/cardFonts.ts).

import { FONT_NAME, type SatoriFont } from "./fonts";

// URLs of the bundled TTF files, resolved relative to this module.
export const monoFontUrls = {
	regular: new URL("./fonts/DejaVuSansMono-Regular.ttf", import.meta.url),
	bold: new URL("./fonts/DejaVuSansMono-Bold.ttf", import.meta.url),
};

// Node helper: read the bundled TTFs and return them as Satori font entries.
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
