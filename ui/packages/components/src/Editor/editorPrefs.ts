import {
	DEFAULT_EDITOR_THEME,
	type EditorThemeId,
	isEditorThemeId,
} from "@gcsim/editor";
import React from "react";

const THEME_KEY = "gcsim-config-editor-color-theme";
const FONT_SIZE_KEY = "gcsim-config-editor-font-size";

export const MIN_FONT_SIZE = 10;
export const MAX_FONT_SIZE = 28;

export interface EditorPrefs {
	theme: EditorThemeId;
	fontSize: number;
}

export const defaultEditorPrefs: EditorPrefs = {
	theme: DEFAULT_EDITOR_THEME,
	fontSize: 14,
};

export function clampFontSize(size: number) {
	if (!Number.isFinite(size)) return defaultEditorPrefs.fontSize;
	return Math.min(MAX_FONT_SIZE, Math.max(MIN_FONT_SIZE, Math.round(size)));
}

type PrefsStorage = Pick<Storage, "getItem" | "setItem">;

export function loadEditorPrefs(storage: PrefsStorage): EditorPrefs {
	const fontSize = storage.getItem(FONT_SIZE_KEY);
	const theme = storage.getItem(THEME_KEY);
	return {
		theme: isEditorThemeId(theme) ? theme : defaultEditorPrefs.theme,
		fontSize: fontSize
			? clampFontSize(Number(fontSize))
			: defaultEditorPrefs.fontSize,
	};
}

export function saveEditorPrefs(storage: PrefsStorage, prefs: EditorPrefs) {
	storage.setItem(THEME_KEY, prefs.theme);
	storage.setItem(FONT_SIZE_KEY, prefs.fontSize.toString());
}

export function useEditorPrefs() {
	const [prefs, setPrefs] = React.useState(() => loadEditorPrefs(localStorage));
	React.useEffect(() => saveEditorPrefs(localStorage, prefs), [prefs]);
	return [prefs, setPrefs] as const;
}
