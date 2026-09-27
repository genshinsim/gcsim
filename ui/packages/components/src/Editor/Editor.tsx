import { Button, Spinner } from "@gcsim/primitives";
import { Play } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { AceEditorWrapper } from "./AceEditorWrapper";
import { type EditorToggles, HelperTools } from "./HelperTools";
import { NameSearch } from "./NameSearch";
import { SectionDivider } from "./SectionDivider";
import { TeamComposer } from "./TeamComposer";
import { ActionListTip, TeamTip } from "./Tips";
import type { EditorProps } from "./types";
import { themes } from "./types";

export const Editor = ({
	config,
	setConfig,
	error,
	parsedTeam,
	settings,
	teamCharacters,
	showThemeSelector = false,
	onRun,
	canRun,
	busy = false,
	prefs,
	onPrefsChange,
}: EditorProps) => {
	const { t } = useTranslation();
	const { toggles, theme, fontSize } = prefs;

	// Ace binds commands once on mount, so the hotkey reads the latest props via a ref.
	const runRef = React.useRef({ onRun, canRun });
	runRef.current = { onRun, canRun };
	const runFromHotkey = React.useCallback(() => {
		if (runRef.current.canRun) {
			runRef.current.onRun();
		}
	}, []);

	const toggle = (key: keyof EditorToggles) =>
		onPrefsChange({ ...prefs, toggles: { ...toggles, [key]: !toggles[key] } });

	return (
		<div className="flex flex-col">
			{toggles.team ? (
				<>
					<SectionDivider>{t("simple.team")}</SectionDivider>
					{toggles.tips ? <TeamTip onHide={() => toggle("tips")} /> : null}
					<TeamComposer
						parsedTeam={parsedTeam}
						error={error}
						config={config}
						setConfig={setConfig}
						characters={teamCharacters}
					/>
				</>
			) : null}

			{toggles.nameSearch ? (
				<>
					<SectionDivider>{t("simple.name_search")}</SectionDivider>
					<NameSearch />
				</>
			) : null}

			<SectionDivider>{t("simple.action_list")}</SectionDivider>
			{toggles.tips ? <ActionListTip onHide={() => toggle("tips")} /> : null}

			{showThemeSelector ? (
				<div className="flex flex-wrap items-center justify-end gap-4">
					<label className="flex items-center gap-2">
						{t("simple.font_size")}
						<input
							type="number"
							value={fontSize}
							onChange={(e) =>
								onPrefsChange({
									...prefs,
									fontSize: Number(e.currentTarget.value),
								})
							}
						/>
					</label>
					<label className="flex items-center gap-2">
						{t("simple.editor_theme")}
						<select
							value={theme}
							onChange={(e) =>
								onPrefsChange({ ...prefs, theme: e.currentTarget.value })
							}
						>
							{themes.map((th) => (
								<option key={th} value={th}>
									{th}
								</option>
							))}
						</select>
					</label>
				</div>
			) : null}

			<AceEditorWrapper
				cfg={config}
				onChange={setConfig}
				onRun={runFromHotkey}
				theme={theme}
				fontSize={fontSize}
			/>

			<div className="sticky bottom-0 z-10 mt-1 flex flex-row flex-wrap place-items-center gap-1 bg-g-canvas p-2">
				<div className="flex flex-grow basis-full items-center p-1 sm:basis-0">
					{settings}
				</div>
				<div className="flex basis-full flex-row flex-wrap gap-1 p-1 sm:basis-2/3">
					<HelperTools toggles={toggles} onToggle={toggle} className="flex-1" />
					<Button className="flex-1" onClick={onRun} disabled={!canRun}>
						{busy ? <Spinner /> : <Play />}
						{t("simple.run")}
					</Button>
				</div>
			</div>
		</div>
	);
};
