export type Theme = {
	id: string;
	name: string;
	light: boolean;
};

// Must match the palettes in packages/theme.css. Cryo is bare :root there.
export const THEMES: Theme[] = [
	{ id: "cryo", name: "Cryo", light: false },
	{ id: "abyss-d", name: "Abyss", light: false },
	{ id: "abyss-l", name: "Abyss", light: true },
	{ id: "ember-d", name: "Ember", light: false },
	{ id: "ember-l", name: "Ember", light: true },
	{ id: "pulse-d", name: "Pulse", light: false },
	{ id: "pulse-l", name: "Pulse", light: true },
	{ id: "twilight", name: "Twilight", light: false },
	{ id: "qilin", name: "Qilin", light: false },
	{ id: "glacier", name: "Glacier", light: false },
	{ id: "azure", name: "Azure", light: false },
	{ id: "aqua", name: "Aqua", light: false },
	{ id: "frostfall", name: "Frostfall", light: true },
	{ id: "blossom", name: "Blossom", light: true },
	{ id: "lantern", name: "Lantern", light: true },
];

export const DEFAULT_THEME_ID = "cryo";

export function getTheme(id: string): Theme {
	return (
		THEMES.find((t) => t.id === id) ??
		(THEMES.find((t) => t.id === DEFAULT_THEME_ID) as Theme)
	);
}
