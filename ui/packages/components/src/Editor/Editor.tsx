import React from "react";
import { useTranslation } from "react-i18next";
import { AceEditorWrapper } from "./AceEditorWrapper";
import { HelperTools } from "./HelperTools";
import { TeamView } from "./TeamView";
import type { EditorProps } from "./types";
import { type Theme, themes } from "./types";

const LOCALSTORAGE_THEME_KEY = "gcsim-config-editor-theme";
const LOCALSTORAGE_FONT_SIZE_KEY = "gcsim-config-editor-font-size";

export const Editor = ({
	config,
	setConfig,
	isValid,
	error,
	parsedTeam,
	run,
	settings,
	showTeam = false,
	showTools = false,
	showThemeSelector = false,
}: EditorProps) => {
	const { t } = useTranslation();
	const [theme, setTheme] = React.useState<Theme>(() => {
		return localStorage.getItem(LOCALSTORAGE_THEME_KEY) ?? "tomorrow_night";
	});
	const [fontSize, setFontSize] = React.useState(() => {
		return localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY)
			? Number(localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY))
			: 14;
	});
	React.useEffect(() => {
		localStorage.setItem(LOCALSTORAGE_THEME_KEY, theme);
		localStorage.setItem(LOCALSTORAGE_FONT_SIZE_KEY, fontSize.toString());
	}, [theme, fontSize]);

	return (
		<div className="flex flex-col gap-2">
			{showThemeSelector || settings ? (
				<div className="flex flex-wrap items-center gap-4">
					{showThemeSelector ? (
						<>
							<label className="flex items-center gap-2">
								{t("simple.editor_theme")}
								<select
									value={theme}
									onChange={(e) => setTheme(e.currentTarget.value)}
								>
									{themes.map((th) => (
										<option key={th} value={th}>
											{th}
										</option>
									))}
								</select>
							</label>
							<label className="flex items-center gap-2">
								{t("simple.font_size")}
								<input
									type="number"
									value={fontSize}
									onChange={(e) => setFontSize(Number(e.currentTarget.value))}
								/>
							</label>
						</>
					) : null}
					{settings}
				</div>
			) : null}
			<AceEditorWrapper
				cfg={config}
				onChange={setConfig}
				onRun={run}
				theme={theme}
				fontSize={fontSize}
			/>
			{showTeam ? (
				<TeamView parsedTeam={parsedTeam} isValid={isValid} error={error} />
			) : null}
			{showTools ? <HelperTools /> : null}
		</div>
	);
};
