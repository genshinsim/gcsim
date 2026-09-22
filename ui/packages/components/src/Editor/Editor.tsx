import { Button, Spinner } from "@gcsim/primitives";
import { Play } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { AceEditorWrapper } from "./AceEditorWrapper";
import { useExecutor } from "./ExecutorProvider";
import { type EditorToggles, HelperTools } from "./HelperTools";
import { NameSearch } from "./NameSearch";
import { TeamView } from "./TeamView";
import { Tips } from "./Tips";
import type { EditorProps } from "./types";
import { type Theme, themes } from "./types";

const LOCALSTORAGE_THEME_KEY = "gcsim-config-editor-theme";
const LOCALSTORAGE_FONT_SIZE_KEY = "gcsim-config-editor-font-size";
const LOCALSTORAGE_TOGGLES_KEY = "gcsim-config-editor-tools";

const defaultToggles: EditorToggles = {
	team: false,
	nameSearch: false,
	tips: false,
};

function loadToggles(): EditorToggles {
	try {
		const raw = localStorage.getItem(LOCALSTORAGE_TOGGLES_KEY);
		return raw ? { ...defaultToggles, ...JSON.parse(raw) } : defaultToggles;
	} catch {
		return defaultToggles;
	}
}

export const Editor = ({
	config,
	setConfig,
	isValid,
	error,
	parsedTeam,
	settings,
	teamCharacters,
	showThemeSelector = false,
}: EditorProps) => {
	const { t } = useTranslation();
	const { run, isReady } = useExecutor();
	const [theme, setTheme] = React.useState<Theme>(() => {
		return localStorage.getItem(LOCALSTORAGE_THEME_KEY) ?? "tomorrow_night";
	});
	const [fontSize, setFontSize] = React.useState(() => {
		return localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY)
			? Number(localStorage.getItem(LOCALSTORAGE_FONT_SIZE_KEY))
			: 14;
	});
	const [toggles, setToggles] = React.useState<EditorToggles>(loadToggles);

	React.useEffect(() => {
		localStorage.setItem(LOCALSTORAGE_THEME_KEY, theme);
		localStorage.setItem(LOCALSTORAGE_FONT_SIZE_KEY, fontSize.toString());
	}, [theme, fontSize]);
	React.useEffect(() => {
		localStorage.setItem(LOCALSTORAGE_TOGGLES_KEY, JSON.stringify(toggles));
	}, [toggles]);

	const toggle = (key: keyof EditorToggles) =>
		setToggles((prev) => ({ ...prev, [key]: !prev[key] }));

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
				onRun={() => run(config)}
				theme={theme}
				fontSize={fontSize}
			/>
			<div className="flex flex-row flex-wrap items-center gap-2">
				<Button onClick={() => run(config)} disabled={!isReady || !isValid}>
					{isReady ? <Play /> : <Spinner />}
					{t("simple.run")}
				</Button>
				<HelperTools toggles={toggles} onToggle={toggle} />
			</div>
			{toggles.team ? (
				<TeamView
					parsedTeam={parsedTeam}
					error={error}
					config={config}
					setConfig={setConfig}
					characters={teamCharacters}
				/>
			) : null}
			{toggles.nameSearch ? <NameSearch /> : null}
			{toggles.tips ? <Tips onHide={() => toggle("tips")} /> : null}
		</div>
	);
};
