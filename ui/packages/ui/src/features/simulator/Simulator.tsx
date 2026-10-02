import {
	type ImportedCharacterOption,
	useEditorPrefs,
	useExecutor,
	useValidation,
} from "@gcsim/components";
import type { model } from "@gcsim/types";
import React from "react";
import { Viewport } from "../../components";
import { useDraft } from "../../stores/AppState";
import { EditorSettings } from "./EditorSettings";
import {
	ImportedCharactersProvider,
	useImportedCharacters,
} from "./ImportedCharacters";
import { SharedConfigDialog } from "./SharedConfigDialog";
import { SimulatorLayout } from "./SimulatorLayout";

function SimulatorEditor() {
	const { cfg, setCfg } = useDraft();
	const { imported } = useImportedCharacters();
	const { isValid, error, parsedTeam } = useValidation(cfg);
	const { run, isReady } = useExecutor();
	const [prefs, setPrefs] = useEditorPrefs();

	const importedCharacters: ImportedCharacterOption[] = React.useMemo(
		() =>
			Object.entries(imported).map(([key, character]) => ({
				key,
				label: character.name,
				character: character as model.Character,
			})),
		[imported],
	);

	return (
		<>
			<SharedConfigDialog onLoad={setCfg} />
			<SimulatorLayout
				config={cfg}
				setConfig={setCfg}
				error={error}
				parsedTeam={parsedTeam}
				importedCharacters={importedCharacters}
				settings={<EditorSettings />}
				onRun={() => run(cfg)}
				canRun={isReady && isValid}
				busy={!isReady}
				prefs={prefs}
				onPrefsChange={setPrefs}
			/>
		</>
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
