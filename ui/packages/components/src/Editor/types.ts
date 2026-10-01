import { DEFAULT_EDITOR_THEME, type EditorThemeId } from "@gcsim/editor";
import type { model } from "@gcsim/types";

export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
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
	importedCharacters?: ImportedCharacterOption[];
}

export interface EditorPrefs {
	theme: EditorThemeId;
	fontSize: number;
}

export const defaultEditorPrefs: EditorPrefs = {
	theme: DEFAULT_EDITOR_THEME,
	fontSize: 14,
};
