import type { model } from "@gcsim/types";

export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
}

export interface EditorProps {
	value: string;
	onChange: (v: string) => void;
	maxLines?: number;
	error?: string | null;
	importedCharacters?: ImportedCharacterOption[];
}
