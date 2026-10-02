import type { EditorThemeId } from "@gcsim/editor";
import type { model } from "@gcsim/types";
import type { EditorPrefs } from "./editorPrefs";

export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
}

export interface EditorProps {
	value: string;
	onChange: (v: string) => void;
	theme?: EditorThemeId;
	fontSize?: number;
	onAppearanceChange?: (next: EditorPrefs) => void;
	maxLines?: number;
	error?: string | null;
	importedCharacters?: ImportedCharacterOption[];
}
