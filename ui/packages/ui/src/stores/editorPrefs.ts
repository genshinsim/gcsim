import { defaultEditorPrefs, type EditorPrefs } from "@gcsim/components";
import React from "react";

const THEME_KEY = "gcsim-config-editor-theme";
const FONT_SIZE_KEY = "gcsim-config-editor-font-size";
const TOGGLES_KEY = "gcsim-config-editor-tools";

type PrefsStorage = Pick<Storage, "getItem" | "setItem">;

function loadToggles(storage: PrefsStorage): EditorPrefs["toggles"] {
	try {
		const raw = storage.getItem(TOGGLES_KEY);
		return raw
			? { ...defaultEditorPrefs.toggles, ...JSON.parse(raw) }
			: defaultEditorPrefs.toggles;
	} catch {
		return defaultEditorPrefs.toggles;
	}
}

export function loadEditorPrefs(storage: PrefsStorage): EditorPrefs {
	const fontSize = storage.getItem(FONT_SIZE_KEY);
	return {
		theme: storage.getItem(THEME_KEY) ?? defaultEditorPrefs.theme,
		fontSize: fontSize ? Number(fontSize) : defaultEditorPrefs.fontSize,
		toggles: loadToggles(storage),
	};
}

export function saveEditorPrefs(storage: PrefsStorage, prefs: EditorPrefs) {
	storage.setItem(THEME_KEY, prefs.theme);
	storage.setItem(FONT_SIZE_KEY, prefs.fontSize.toString());
	storage.setItem(TOGGLES_KEY, JSON.stringify(prefs.toggles));
}

export function useEditorPrefs() {
	const [prefs, setPrefs] = React.useState(() => loadEditorPrefs(localStorage));
	React.useEffect(() => saveEditorPrefs(localStorage, prefs), [prefs]);
	return [prefs, setPrefs] as const;
}
