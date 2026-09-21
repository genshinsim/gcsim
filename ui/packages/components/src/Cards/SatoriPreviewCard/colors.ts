// Self-contained palette for the Satori preview card. Values mirror the
// blueprint-derived colors used by the live PreviewCard (see common/gcsim), but
// are inlined here so this component pulls in no @visx/* or i18next code.

// Satori can't read the --g-sN CSS vars, so the series palette is inlined here
// and must be kept in sync with theme.css by hand.
const qualitative3 = [
	"#cc79a7",
	"#f0e442",
	"#d55e00",
	"#56b4e9",
	"#e69f00",
	"#009e73",
	"#0072b2",
];

// Label tint of each qualitative3 entry: 0.6*base + 0.4*white (this card's
// fixed dark background makes --g-text light, so q5 mixes toward white).
const qualitative5 = [
	"#e0afca",
	"#f6ef8e",
	"#e69e66",
	"#9ad2f2",
	"#f0c566",
	"#66c5ab",
	"#66aad1",
];

export function characterColor(i: number): string {
	return qualitative3[i % qualitative3.length];
}

export function levelColor(i: number): string {
	return qualitative5[i % qualitative5.length];
}

// Element slice colors keyed by the raw element name, matching the `value`
// colors in DataColors.element.
const elementColors: Record<string, string> = {
	electro: "#9D3F9D",
	pyro: "#D33D17",
	cryo: "#4B8DAA",
	hydro: "#147EB3",
	dendro: "#29A634",
	anemo: "#00A396",
	geo: "#D1980B",
	physical: "#946638",
	frozen: "#00A396",
	quicken: "#8EB125",
};

export function elementColor(key: string): string {
	return elementColors[key] ?? "#9ca3af"; // tailwind gray-400 fallback
}

// Timeline area color: DataColorsConst.qualitative3(8) === SEPIA3.
export const TIMELINE_COLOR = "#946638";

// Histogram bar colors, matching the live Graphs/Histogram (Vermilion).
export const HISTOGRAM_BAR = "#D33D17"; // VERMILION3
export const HISTOGRAM_ACCENT = "#96290D"; // VERMILION1

// Card chrome.
export const SLATE_800 = "#1e293b";
export const SLATE_700 = "#334155";
export const GRAY_600 = "#4b5563";
export const GRAY_700 = "#374151";
export const GRAY_400 = "#9ca3af";

// Metadata pill palette, matching the live badge variants (@gcsim/primitives badge).
export const PRIMARY_BG = "#0f172a";
export const PRIMARY_FG = "#f8fafc";
export const AMBER_700 = "#b45309"; // warning bg
export const ROSE_700 = "#be123c"; // danger bg

// Portrait background images (pre-blended per element) live in ./backgrounds.
