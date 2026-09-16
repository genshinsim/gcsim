import type { Character } from "@gcsim/types";
import type React from "react";

export interface EditorProps {
	config: string;
	setConfig: (v: string) => void;
	isValid: boolean;
	error: string | null;
	parsedTeam: Character[];
	run: () => void;
	settings?: React.ReactNode;
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
