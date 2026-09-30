import { Editor, useExecutor, useValidation } from "@gcsim/components";
import type { model } from "@gcsim/types";
import { getRouteApi } from "@tanstack/react-router";
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

const routeApi = getRouteApi("/simulator");

// Injects a `?cfg=` search param straight into the editor's text, once, on
// arrival - e.g. from an external tool linking "open in gcsim" with an
// already-built config. Calls the plain setCfg so this just prefills the
// textarea, letting the user review/edit/run it themselves.
function useCfgFromSearch(setCfg: (cfg: string) => void) {
	const { cfg } = routeApi.useSearch();
	const navigate = routeApi.useNavigate();
	// Consumed at most once per link: without stripping it from the URL
	// afterwards, a refresh (or reopeqning the same bookmarked link) would
	// keep clobbering any edits the user made since.
	const consumed = React.useRef(false);

	React.useEffect(() => {
		if (!cfg || consumed.current) return;
		consumed.current = true;
		setCfg(cfg);
		void navigate({ to: "/simulator", search: {}, replace: true });
	}, [cfg, navigate, setCfg]);
}

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
	useCfgFromSearch(setCfg);
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
