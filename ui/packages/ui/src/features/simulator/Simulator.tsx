import { Editor, useExecutor, useValidation } from "@gcsim/components";
import type { model } from "@gcsim/types";
import React from "react";
import { Viewport } from "../../components";
import { CharMap } from "../../data";
import { useDraft } from "../../stores/AppState";
import { useEditorPrefs } from "../../stores/editorPrefs";
import { EditorSettings } from "./EditorSettings";
import {
	ImportedCharactersProvider,
	useImportedCharacters,
} from "./ImportedCharacters";

const defaultStats = [
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

function newCharFromKey(key: string): model.Character {
	return {
		name: key,
		level: 80,
		max_level: 90,
		element: CharMap[key].element,
		cons: 0,
		weapon: { name: "dullblade", refine: 1, level: 1, max_level: 20 },
		talents: { attack: 6, skill: 6, burst: 6 },
		stats: [...defaultStats],
		snapshot: [...defaultStats],
		sets: {},
	};
}

function SimulatorEditor() {
	const { cfg, setCfg } = useDraft();
	const { imported } = useImportedCharacters();
	const { isValid, error, parsedTeam } = useValidation(cfg);
	const { run, isReady } = useExecutor();
	const [prefs, setPrefs] = useEditorPrefs();

	const teamCharacters = React.useMemo(
		() => ({
			createCharacter: newCharFromKey,
			imported: Object.entries(imported).map(([key, character]) => ({
				key,
				label: character.name,
				character: character as model.Character,
			})),
		}),
		[imported],
	);

	return (
		<Editor
			config={cfg}
			setConfig={setCfg}
			error={error}
			parsedTeam={parsedTeam}
			teamCharacters={teamCharacters}
			settings={<EditorSettings />}
			showThemeSelector
			onRun={() => run(cfg)}
			canRun={isReady && isValid}
			busy={!isReady}
			prefs={prefs}
			onPrefsChange={setPrefs}
		/>
	);
}

export function Simulator() {
	return (
		<Viewport className="flex flex-col gap-2">
			<ImportedCharactersProvider>
				<SimulatorEditor />
			</ImportedCharactersProvider>
		</Viewport>
	);
}
