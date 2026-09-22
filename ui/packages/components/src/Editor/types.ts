import type { model } from "@gcsim/types";
import type React from "react";

// A GOOD/Enka-imported character offered in the TeamView add picker alongside
// the base roster. TEMPORARY: exists only for the add/remove crutch.
export interface ImportedCharacterOption {
	key: string;
	label?: string;
	character: model.Character;
}

// App-injected source for the TeamView add crutch. The base roster ships with
// @gcsim/components; only default-character construction (needs element data)
// and imported characters are app-specific, so they arrive through here.
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
	run: () => void;
	settings?: React.ReactNode;
	teamCharacters?: TeamViewCharacterSource;
	showTeam?: boolean;
	showTools?: boolean;
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
