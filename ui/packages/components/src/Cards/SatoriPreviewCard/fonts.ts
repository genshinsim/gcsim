// Inter is bundled with this package (src/Cards/SatoriPreviewCard/fonts). The
// card tree uses `fontFamily: "Inter"`; the caller passes the loaded buffers via
// Satori's `fonts` option, e.g. `satori(<SatoriPreviewCard .../>, { fonts })`.

// Shape of a single Satori font entry (kept local so this package needs no
// satori dependency).
export type SatoriFont = {
	name: string;
	data: ArrayBuffer | Buffer;
	weight: 400 | 700;
	style: "normal";
};

// URLs of the bundled TTF files, resolved relative to this module.
export const interFontUrls = {
	regular: new URL("./fonts/Inter-Regular.ttf", import.meta.url),
	bold: new URL("./fonts/Inter-Bold.ttf", import.meta.url),
};

// Node helper: read the bundled Inter TTFs and return them as Satori font
// entries. `node:fs` is imported dynamically so this module stays safe to bundle
// for the browser (Storybook renders the card without calling this).
export async function loadInterFonts(): Promise<SatoriFont[]> {
	const { readFile } = await import("node:fs/promises");
	const [regular, bold] = await Promise.all([
		readFile(interFontUrls.regular),
		readFile(interFontUrls.bold),
	]);
	return [
		{ name: "Inter", data: regular, weight: 400, style: "normal" },
		{ name: "Inter", data: bold, weight: 700, style: "normal" },
	];
}
