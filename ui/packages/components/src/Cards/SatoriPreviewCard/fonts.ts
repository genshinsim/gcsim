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
// satori dependency). `data` is typed `ArrayBuffer | Uint8Array` (not the Node
// `Buffer` global) so this module needs no `@types/node`; a Node `Buffer` still
// satisfies `Uint8Array`.
export type SatoriFont = {
	name: string;
	data: ArrayBuffer | Uint8Array;
	weight: 400 | 700;
	style: "normal";
};

// The Node TTF loader (monoFontUrls / loadCardFonts) lives in ./fontsNode so
// this module stays free of any `node:*` reference — it is imported by the card
// tree and must bundle cleanly for the browser and the Cloudflare Worker.
