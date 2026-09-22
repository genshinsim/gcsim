import type { model } from "@gcsim/types";
import type React from "react";

// TEMPORARY: types for the TeamView add/remove crutch.
export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
}

export interface TeamViewCharacterSource {
	createCharacter: (key: string) => model.Character;
	imported?: ImportedCharacterOption[];
}

export interface EditorProps {
	config: string;
	setConfig: (v: string) => void;
	isValid: boolean;
	error: string | null;
	parsedTeam: model.Character[];
	settings?: React.ReactNode;
	teamCharacters?: TeamViewCharacterSource;
	showThemeSelector?: boolean;
}

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
