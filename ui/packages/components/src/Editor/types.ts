import { DEFAULT_EDITOR_THEME, type EditorThemeId } from "@gcsim/editor";
import type { model } from "@gcsim/types";

// TEMPORARY: types for the TeamComposer add/remove crutch.
export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
}

export interface TeamComposerCharacterSource {
	createCharacter: (key: string) => model.Character;
	imported?: ImportedCharacterOption[];
}

export interface EditorAppearance {
	theme: EditorThemeId;
	fontSize: number;
}

export interface EditorProps {
	value: string;
	onChange: (v: string) => void;
	theme?: EditorThemeId;
	fontSize?: number;
	onAppearanceChange?: (next: EditorAppearance) => void;
	maxLines?: number;
	error?: string | null;
}

export interface EditorToggles {
	team: boolean;
	nameSearch: boolean;
	tips: boolean;
}

export interface EditorPrefs {
	toggles: EditorToggles;
	theme: EditorThemeId;
	fontSize: number;
}

export const defaultEditorPrefs: EditorPrefs = {
	toggles: { team: true, nameSearch: true, tips: true },
	theme: DEFAULT_EDITOR_THEME,
	fontSize: 14,
};
