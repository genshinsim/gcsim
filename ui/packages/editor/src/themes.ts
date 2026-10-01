import { HighlightStyle, syntaxHighlighting } from "@codemirror/language";
import type { Extension } from "@codemirror/state";
import { EditorView } from "@codemirror/view";
import { tags as t } from "@lezer/highlight";

interface Palette {
	background: string;
	foreground: string;
	gutter: string;
	border: string;
	selection: string;
	cursor: string;
	raised: string;
	scrollbar: string;
	danger: string;
	comment: string;
	keyword: string;
	string: string;
	number: string;
	character: string;
	action: string;
	stat: string;
	element: string;
	punctuation: string;
}

export interface EditorTheme {
	id: string;
	label: string;
	dark?: boolean;
	palette: Palette;
}

export const EDITOR_THEMES = [
	{
		id: "app",
		label: "gcsim",
		palette: {
			background: "var(--g-surface)",
			foreground: "var(--g-text)",
			gutter: "var(--g-text-mute)",
			border: "var(--g-border)",
			selection: "var(--g-accent-weak)",
			cursor: "var(--g-accent)",
			raised: "var(--g-surface-2)",
			scrollbar: "var(--g-surface-3)",
			danger: "var(--g-danger)",
			comment: "var(--g-text-mute)",
			keyword: "var(--g-accent)",
			string: "var(--g-dendro)",
			number: "var(--g-geo)",
			character: "var(--g-pyro)",
			action: "var(--g-hydro)",
			stat: "var(--g-electro)",
			element: "var(--g-anemo)",
			punctuation: "var(--g-text-dim)",
		},
	},
	{
		id: "one_dark",
		label: "One Dark",
		dark: true,
		palette: {
			background: "#282c34",
			foreground: "#abb2bf",
			gutter: "#636d83",
			border: "#181a1f",
			selection: "#3e4451",
			cursor: "#528bff",
			raised: "#21252b",
			scrollbar: "#4b5263",
			danger: "#e06c75",
			comment: "#7f848e",
			keyword: "#c678dd",
			string: "#98c379",
			number: "#d19a66",
			character: "#e5c07b",
			action: "#61afef",
			stat: "#e06c75",
			element: "#56b6c2",
			punctuation: "#abb2bf",
		},
	},
	{
		id: "tomorrow_night",
		label: "Tomorrow Night",
		dark: true,
		palette: {
			background: "#1d1f21",
			foreground: "#c5c8c6",
			gutter: "#6f7270",
			border: "#282a2e",
			selection: "#373b41",
			cursor: "#aeafad",
			raised: "#282a2e",
			scrollbar: "#4d5057",
			danger: "#cc6666",
			comment: "#969896",
			keyword: "#b294bb",
			string: "#b5bd68",
			number: "#de935f",
			character: "#f0c674",
			action: "#81a2be",
			stat: "#cc6666",
			element: "#8abeb7",
			punctuation: "#c5c8c6",
		},
	},
	{
		id: "monokai",
		label: "Monokai",
		dark: true,
		palette: {
			background: "#272822",
			foreground: "#f8f8f2",
			gutter: "#90908a",
			border: "#1e1f1c",
			selection: "#49483e",
			cursor: "#f8f8f0",
			raised: "#3e3d32",
			scrollbar: "#5b5a4e",
			danger: "#f92672",
			comment: "#75715e",
			keyword: "#f92672",
			string: "#e6db74",
			number: "#ae81ff",
			character: "#a6e22e",
			action: "#66d9ef",
			stat: "#fd971f",
			element: "#ae81ff",
			punctuation: "#f8f8f2",
		},
	},
	{
		id: "dracula",
		label: "Dracula",
		dark: true,
		palette: {
			background: "#282a36",
			foreground: "#f8f8f2",
			gutter: "#6272a4",
			border: "#191a21",
			selection: "#44475a",
			cursor: "#f8f8f0",
			raised: "#343746",
			scrollbar: "#565a6e",
			danger: "#ff5555",
			comment: "#6272a4",
			keyword: "#ff79c6",
			string: "#f1fa8c",
			number: "#bd93f9",
			character: "#50fa7b",
			action: "#8be9fd",
			stat: "#ffb86c",
			element: "#bd93f9",
			punctuation: "#f8f8f2",
		},
	},
	{
		id: "solarized_dark",
		label: "Solarized Dark",
		dark: true,
		palette: {
			background: "#002b36",
			foreground: "#839496",
			gutter: "#586e75",
			border: "#073642",
			selection: "#073642",
			cursor: "#93a1a1",
			raised: "#073642",
			scrollbar: "#0d4b5c",
			danger: "#dc322f",
			comment: "#586e75",
			keyword: "#859900",
			string: "#2aa198",
			number: "#d33682",
			character: "#b58900",
			action: "#268bd2",
			stat: "#cb4b16",
			element: "#6c71c4",
			punctuation: "#839496",
		},
	},
	{
		id: "solarized_light",
		label: "Solarized Light",
		palette: {
			background: "#fdf6e3",
			foreground: "#657b83",
			gutter: "#93a1a1",
			border: "#eee8d5",
			selection: "#eee8d5",
			cursor: "#586e75",
			raised: "#eee8d5",
			scrollbar: "#d6ceb5",
			danger: "#dc322f",
			comment: "#93a1a1",
			keyword: "#859900",
			string: "#2aa198",
			number: "#d33682",
			character: "#b58900",
			action: "#268bd2",
			stat: "#cb4b16",
			element: "#6c71c4",
			punctuation: "#657b83",
		},
	},
	{
		id: "github_light",
		label: "GitHub Light",
		palette: {
			background: "#ffffff",
			foreground: "#24292f",
			gutter: "#8c959f",
			border: "#d0d7de",
			selection: "#b6e3ff",
			cursor: "#0969da",
			raised: "#f6f8fa",
			scrollbar: "#d0d7de",
			danger: "#cf222e",
			comment: "#6e7781",
			keyword: "#cf222e",
			string: "#0a3069",
			number: "#0550ae",
			character: "#8250df",
			action: "#953800",
			stat: "#116329",
			element: "#0969da",
			punctuation: "#24292f",
		},
	},
] as const satisfies readonly EditorTheme[];

export type EditorThemeId = (typeof EDITOR_THEMES)[number]["id"];

export const DEFAULT_EDITOR_THEME: EditorThemeId = "app";

export function isEditorThemeId(id: unknown): id is EditorThemeId {
	return EDITOR_THEMES.some((theme) => theme.id === id);
}

function build({ palette: p, dark }: EditorTheme): Extension {
	const view = EditorView.theme(
		{
			"&": {
				backgroundColor: p.background,
				color: p.foreground,
				border: `1px solid ${p.border}`,
				borderRadius: "var(--g-r-md)",
				// CodeMirror stacks its layers up to z-index 300 (panels); keep
				// them inside the editor so app menus and dialogs draw above.
				isolation: "isolate",
			},
			"&.cm-focused": { outline: `1px solid ${p.cursor}` },
			// Chrome keeps a stale "//" ligature glyph (blank first slash) when
			// CodeMirror edits a text node in place.
			".cm-scroller": {
				fontFamily: "var(--g-font-mono)",
				fontVariantLigatures: "none",
				lineHeight: "1.5",
			},
			// scrollbar-color covers Firefox and Chromium; the pseudo-elements
			// cover Safari.
			".cm-scroller, .cm-tooltip-autocomplete > ul, .cm-gcsim-errors": {
				scrollbarColor: `${p.scrollbar} transparent`,
			},
			".cm-scroller::-webkit-scrollbar, .cm-tooltip-autocomplete > ul::-webkit-scrollbar, .cm-gcsim-errors::-webkit-scrollbar":
				{ width: "8px", height: "8px" },
			".cm-scroller::-webkit-scrollbar-track, .cm-tooltip-autocomplete > ul::-webkit-scrollbar-track, .cm-gcsim-errors::-webkit-scrollbar-track":
				{ background: "transparent" },
			".cm-scroller::-webkit-scrollbar-thumb, .cm-tooltip-autocomplete > ul::-webkit-scrollbar-thumb, .cm-gcsim-errors::-webkit-scrollbar-thumb":
				{ background: p.scrollbar, borderRadius: "4px" },
			".cm-scroller::-webkit-scrollbar-corner": { background: "transparent" },
			".cm-content": { caretColor: p.cursor },
			".cm-cursor, .cm-dropCursor": { borderLeftColor: p.cursor },
			"&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, ::selection":
				{ backgroundColor: p.selection },
			".cm-gutters": {
				backgroundColor: p.background,
				color: p.gutter,
				border: "none",
				borderRight: `1px solid ${p.border}`,
			},
			".cm-activeLineGutter": {
				backgroundColor: p.raised,
				color: p.foreground,
			},
			"&.cm-focused .cm-matchingBracket": {
				backgroundColor: p.selection,
				outline: `1px solid ${p.cursor}`,
			},
			".cm-tooltip, .cm-panels": {
				backgroundColor: p.raised,
				border: `1px solid ${p.border}`,
				color: p.foreground,
			},
			".cm-tooltip-autocomplete > ul > li[aria-selected]": {
				backgroundColor: p.selection,
				color: p.foreground,
			},
			".cm-completionDetail": { color: p.comment },
			".cm-searchMatch": {
				backgroundColor: `color-mix(in srgb, ${p.number} 30%, transparent)`,
			},
			".cm-foldPlaceholder": {
				backgroundColor: p.raised,
				color: p.comment,
				border: "none",
			},
			".cm-panels-bottom": {
				position: "static",
				border: "none",
				borderTop: `1px solid ${p.border}`,
				borderRadius: "0 0 var(--g-r-md) var(--g-r-md)",
			},
			".cm-gcsim-errors": {
				maxHeight: "9em",
				overflowY: "auto",
				padding: "6px 10px",
				fontFamily: "var(--g-font-mono)",
				fontSize: "0.9em",
				lineHeight: "1.5",
				borderLeft: `3px solid ${p.danger}`,
				backgroundColor: `color-mix(in srgb, ${p.danger} 8%, ${p.background})`,
			},
			".cm-gcsim-errors-title": { color: p.danger, fontWeight: "600" },
			".cm-gcsim-errors-list": { margin: "0", padding: "0", listStyle: "none" },
			".cm-gcsim-errors-list > li": { whiteSpace: "pre-wrap" },
			".cm-gcsim-errors-pos": {
				font: "inherit",
				color: p.danger,
				background: "none",
				border: "none",
				padding: "0",
				marginRight: "1ch",
				cursor: "pointer",
				textDecoration: "underline dotted",
			},
			".cm-lintRange-error": {
				backgroundImage: "none",
				textDecoration: `underline wavy ${p.danger}`,
			},
		},
		{ dark },
	);
	const highlight = HighlightStyle.define([
		{ tag: t.keyword, color: p.keyword, fontWeight: "600" },
		{ tag: t.lineComment, color: p.comment, fontStyle: "italic" },
		{ tag: t.string, color: p.string },
		{ tag: [t.number, t.bool], color: p.number },
		{ tag: t.className, color: p.character, fontWeight: "600" },
		{ tag: t.function(t.variableName), color: p.action },
		{ tag: t.attributeName, color: p.stat },
		{ tag: t.atom, color: p.element },
		{ tag: t.variableName, color: p.foreground },
		{ tag: [t.operator, t.punctuation], color: p.punctuation },
	]);
	return [view, syntaxHighlighting(highlight)];
}

const built = new Map<string, Extension>();

// Extensions are cached so switching back to a theme reuses the same
// instance and CodeMirror doesn't mount a duplicate style sheet.
export function editorTheme(id: EditorThemeId): Extension {
	let ext = built.get(id);
	if (!ext) {
		const theme = EDITOR_THEMES.find((th) => th.id === id) ?? EDITOR_THEMES[0];
		ext = build(theme);
		built.set(id, ext);
	}
	return ext;
}
