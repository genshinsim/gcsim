import {
	clampFontSize,
	defaultEditorPrefs,
	type EditorPrefs,
	isEditorThemeId,
} from "@gcsim/components";
import React from "react";

const THEME_KEY = "gcsim-config-editor-color-theme";
const FONT_SIZE_KEY = "gcsim-config-editor-font-size";

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
