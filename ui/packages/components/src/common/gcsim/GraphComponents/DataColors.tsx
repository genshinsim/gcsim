import { scaleOrdinal } from "@visx/scale";
import i18next from "i18next";

// Consume the raw `--g-*` palette layer, not the coined `--color-g-*` tokens:
// the latter are `@theme inline` values that only exist as CSS vars where a
// Tailwind utility references them, so they resolve to nothing inside SVG fills.
function series(i: number) {
	return `var(--g-s${(i % 6) + 1})`;
}

function shade(base: string, towards: "text" | "bg", pct: number) {
	return `color-mix(in srgb, ${base} ${pct}%, var(--g-${towards}))`;
}

const q1 = (i: number) => shade(series(i), "bg", 55);
const q2 = (i: number) => shade(series(i), "bg", 78);
const q3 = (i: number) => series(i);
const q4 = (i: number) => shade(series(i), "text", 80);
const q5 = (i: number) => shade(series(i), "text", 60);

function element(token: string) {
	const base = `var(--g-${token})`;
	return {
		value: base,
		label: shade(base, "text", 80),
		highlight: shade(base, "text", 60),
	};
}

function seriesTriple(i: number) {
	return { value: q3(i), label: q4(i), highlight: q5(i) };
}

type ActionColor = {
	highlight: string;
	label: string;
	value: string;
};

type ElementColor = {
	label: string;
	highlight: string;
	value: string;
};

export const DataColorsConst = {
	gray: "var(--g-border)",

	qualitative1: (i: number) => q1(i),
	qualitative2: (i: number) => q2(i),
	qualitative3: (i: number) => q3(i),
	qualitative4: (i: number) => q4(i),
	qualitative5: (i: number) => q5(i),
};

export function useDataColors() {
	const actionKeys = [
		"actions.attack",
		"actions.charge",
		"actions.aim",
		"actions.skill",
		"actions.burst",
		"actions.low_plunge",
		"actions.high_plunge",
		"actions.dash",
		"actions.jump",
		"actions.walk",
		"actions.swap",
	] as const;
	const actions: Map<string, ActionColor> = new Map(
		actionKeys.map((key, i) => [i18next.t(key), seriesTriple(i)]),
	);

	const actionColor = scaleOrdinal<string, string>({
		domain: Array.from(actions.keys()),
		range: Array.from(actions.values()).map((e) => e.value),
	});

	const actionLabelColor = scaleOrdinal<string, string>({
		domain: Array.from(actions.keys()),
		range: Array.from(actions.values()).map((e) => e.label),
	});

	const actionHighlightColor = scaleOrdinal<string, string>({
		domain: Array.from(actions.keys()),
		range: Array.from(actions.values()).map((e) => e.highlight),
	});

	const elements: Map<string, ElementColor> = new Map([
		[i18next.t("elements.electro"), element("electro")],
		[i18next.t("elements.pyro"), element("pyro")],
		[i18next.t("elements.cryo"), element("cryo")],
		[i18next.t("elements.hydro"), element("hydro")],
		[i18next.t("elements.dendro"), element("dendro")],
		[i18next.t("elements.anemo"), element("anemo")],
		[i18next.t("elements.geo"), element("geo")],
		[i18next.t("elements.physical"), element("text-dim")],

		// not possible, but defined in attributes/element.go so here just in case
		[i18next.t("elements.frozen"), element("cryo")],
		[i18next.t("elements.quicken"), element("dendro")],
	]);

	const elementColor = scaleOrdinal<string, string>({
		domain: Array.from(elements.keys()),
		range: Array.from(elements.values()).map((e) => e.value),
	});

	const elementLabelColor = scaleOrdinal<string, string>({
		domain: Array.from(elements.keys()),
		range: Array.from(elements.values()).map((e) => e.label),
	});

	const elementHighlightColor = scaleOrdinal<string, string>({
		domain: Array.from(elements.keys()),
		range: Array.from(elements.values()).map((e) => e.highlight),
	});

	const reactableModifiers: Map<string, ElementColor> = new Map([
		[i18next.t("elements.electro"), element("electro")],
		[i18next.t("elements.pyro"), element("pyro")],
		[i18next.t("elements.cryo"), element("cryo")],
		[i18next.t("elements.hydro"), element("hydro")],
		[i18next.t("elements.dendro"), element("dendro")],
		[i18next.t("elements.anemo"), element("anemo")],
		[i18next.t("elements.geo"), element("geo")],
		[i18next.t("elements.frozen"), element("cryo")],
		[i18next.t("elements.quicken"), element("dendro")],
		[i18next.t("elements.dendro-fuel"), element("dendro")],
		[i18next.t("elements.burning"), element("pyro")],
	]);

	const reactableModifierColor = scaleOrdinal<string, string>({
		domain: Array.from(reactableModifiers.keys()),
		range: Array.from(reactableModifiers.values()).map((e) => e.value),
	});

	const reactableModifierLabelColor = scaleOrdinal<string, string>({
		domain: Array.from(reactableModifiers.keys()),
		range: Array.from(reactableModifiers.values()).map((e) => e.label),
	});

	const reactableModifierHighlightColor = scaleOrdinal<string, string>({
		domain: Array.from(reactableModifiers.keys()),
		range: Array.from(reactableModifiers.values()).map((e) => e.highlight),
	});

	return {
		DataColors: {
			reactableModifierKeys: [...reactableModifiers.keys()],
			reactableModifier: reactableModifierColor,
			reactableModifierLabel: reactableModifierLabelColor,
			reactableModifierHighlight: reactableModifierHighlightColor,

			actionKeys: [...actions.keys()],
			action: actionColor,
			actionLabel: actionLabelColor,
			actionHighlight: actionHighlightColor,

			element: elementColor,
			elementLabel: elementLabelColor,
			elementHighlight: elementHighlightColor,

			// TODO: better colors for characters?
			character: (i: number) => q3(i),
			characterLabel: (i: number) => q4(i),

			target: (k: string) => q3(Number(k) - 1),
			targetLabel: (k: string) => q4(Number(k) - 1),
		},
	};
}
