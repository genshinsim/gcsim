// Self-contained palette for the Satori preview card. Values mirror the
// blueprint-derived colors used by the live PreviewCard (see common/gcsim), but
// are inlined here so this component pulls in no @visx/* or i18next code.

// Character slice colors, matching DataColors.character (qualitative3).
const qualitative3 = [
	"#147EB3", // CERULEAN3
	"#29A634", // FOREST3
	"#D1980B", // GOLD3
	"#D33D17", // VERMILION3
	"#9D3F9D", // VIOLET3
	"#00A396", // TURQUOISE3
	"#DB2C6F", // ROSE3
	"#8EB125", // LIME3
	"#946638", // SEPIA3
	"#7961DB", // INDIGO3
];

// Level-label colors, matching DataColorsConst.qualitative5.
const qualitative5 = [
	"#68C1EE", // CERULEAN5
	"#62D96B", // FOREST5
	"#FBD065", // GOLD5
	"#FF9980", // VERMILION5
	"#D69FD6", // VIOLET5
	"#7AE1D8", // TURQUOISE5
	"#FF66A1", // ROSE5
	"#D4F17E", // LIME5
	"#D0B090", // SEPIA5
	"#BDADFF", // INDIGO5
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

// Metadata pill palette, matching the live badge variants (common/ui/badge).
export const PRIMARY_BG = "#0f172a"; // bg-primary (slate-900)
export const PRIMARY_FG = "#f8fafc"; // text-primary-foreground
export const AMBER_700 = "#b45309"; // warning bg
export const ROSE_700 = "#be123c"; // danger bg

// Portrait background images (pre-blended per element) live in ./backgrounds.
