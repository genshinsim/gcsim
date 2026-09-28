import type { model } from "@gcsim/types";
import type React from "react";

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

export interface EditorProps {
	config: string;
	setConfig: (v: string) => void;
	error: string | null;
	parsedTeam: model.Character[];
	settings?: React.ReactNode;
	teamCharacters?: TeamComposerCharacterSource;
	showThemeSelector?: boolean;
	showHelpers?: boolean;
	onRun: () => void;
	canRun: boolean;
	busy?: boolean;
	prefs: EditorPrefs;
	onPrefsChange: (next: EditorPrefs) => void;
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

export interface AceEditorWrapperProps {
	cfg: string;
	onChange: (v: string) => void;
	onRun?: () => void;
	maxLines?: number;
	fontSize?: number;
	theme?: Theme;
}

export const themes = [
	"monokai",
	"github",
	"tomorrow",
	"tomorrow_night",
	"kuroir",
	"twilight",
	"xcode",
	"textmate",
	"solarized_dark",
	"solarized_light",
	"terminal",
];
export type Theme = (typeof themes)[number];
