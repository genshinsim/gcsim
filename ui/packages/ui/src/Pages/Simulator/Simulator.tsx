import { Editor, ExecutorProvider, useValidation } from "@gcsim/components";
import type { Executor, ExecutorSupplier } from "@gcsim/executors";
import type { model, SimResults } from "@gcsim/types";
import { throttle } from "lodash-es";
import React from "react";
import { useNavigate } from "react-router";
import { Viewport } from "../../Components";
import { CharMap } from "../../Data";
import { appActions, defaultStats } from "../../Stores/appSlice";
import {
	type RootState,
	useAppDispatch,
	useAppSelector,
} from "../../Stores/store";
import { viewerActions } from "../../Stores/viewerSlice";
import { VIEWER_THROTTLE } from "../Viewer";
import { EditorSettings } from "./EditorSettings";

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

function SimulatorEditor({ cfg }: { cfg: string }) {
	const dispatch = useAppDispatch();
	const imported = useAppSelector(
		(state: RootState) => state.user_data.GOODImport,
	);
	const { isValid, error, parsedTeam } = useValidation(cfg);

	const setConfig = (newCfg: string) => {
		dispatch(appActions.setCfg({ cfg: newCfg, keepTeam: false }));
	};

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
			setConfig={setConfig}
			isValid={isValid}
			error={error}
			parsedTeam={parsedTeam}
			teamCharacters={teamCharacters}
			settings={<EditorSettings />}
			showThemeSelector
		/>
	);
}

export function Simulator({ exec }: { exec: ExecutorSupplier<Executor> }) {
	const dispatch = useAppDispatch();
	const navigate = useNavigate();
	const cfg = useAppSelector((state: RootState) => state.app.cfg);

	const onResult = React.useMemo(
		() =>
			throttle(
				(res: SimResults, hash: string) => {
					dispatch(viewerActions.setResult({ data: res, hash }));
				},
				VIEWER_THROTTLE,
				{ leading: true, trailing: true },
			),
		[dispatch],
	);

	const navigateOnRun = () => {
		dispatch(viewerActions.start());
		navigate("/web");
	};

	return (
		<Viewport className="flex flex-col gap-2">
			<ExecutorProvider
				exec={exec}
				onResult={onResult}
				navigateOnRun={navigateOnRun}
			>
				<SimulatorEditor cfg={cfg} />
			</ExecutorProvider>
		</Viewport>
	);
}
