import {
	FONT_NAME,
	type SatoriFont,
} from "@gcsim/components/src/Cards/SatoriPreviewCard/fonts";
// Bundle the TTF bytes at build time via the wrangler `Data` module rule for
// `**/*.ttf` (see wrangler.jsonc); each import is the file's raw ArrayBuffer.
// This is the Worker-side replacement for the components package's Node
// `loadCardFonts()` (which uses `node:fs` and cannot run on a Worker).
import boldFont from "@gcsim/components/src/Cards/SatoriPreviewCard/fonts/DejaVuSansMono-Bold.ttf";
import regularFont from "@gcsim/components/src/Cards/SatoriPreviewCard/fonts/DejaVuSansMono-Regular.ttf";

// Satori font entries for the card, in the `{ name, data, weight, style }`
// shape the card's FONT_NAME / SatoriFont already define. Reuses FONT_NAME so
// the family stays in sync with the card tree's `fontFamily`.
export const cardFonts: SatoriFont[] = [
	{ name: FONT_NAME, data: regularFont, weight: 400, style: "normal" },
	{ name: FONT_NAME, data: boldFont, weight: 700, style: "normal" },
];
