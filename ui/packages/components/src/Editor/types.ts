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

export type Theme = string;

export interface EditorAppearance {
	theme: Theme;
	fontSize: number;
}

export interface EditorProps {
	value: string;
	onChange: (v: string) => void;
	theme?: Theme;
	fontSize?: number;
	onAppearanceChange?: (next: EditorAppearance) => void;
	maxLines?: number;
}

export interface EditorToggles {
	team: boolean;
	nameSearch: boolean;
	tips: boolean;
}

export interface EditorPrefs {
	toggles: EditorToggles;
	theme: Theme;
	fontSize: number;
}

export const defaultEditorPrefs: EditorPrefs = {
	toggles: { team: true, nameSearch: true, tips: true },
	theme: "tomorrow_night",
	fontSize: 14,
};
