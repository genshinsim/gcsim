import {
	Editor,
	type ImportedCharacterOption,
	SectionDivider,
	TeamCard,
} from "@gcsim/components";
import { Button, Spinner } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { Play } from "lucide-react";
import type React from "react";
import { useTranslation } from "react-i18next";

export interface SimulatorLayoutProps {
	config: string;
	setConfig: (v: string) => void;
	error: string | null;
	parsedTeam: model.Character[];
	importedCharacters?: ImportedCharacterOption[];
	settings?: React.ReactNode;
	onRun: () => void;
	canRun: boolean;
	busy: boolean;
}

export function SimulatorLayout({
	config,
	setConfig,
	error,
	parsedTeam,
	importedCharacters,
	settings,
	onRun,
	canRun,
	busy,
}: SimulatorLayoutProps) {
	const { t } = useTranslation();

	return (
		<div className="flex flex-col">
			<SectionDivider>{t("simple.team")}</SectionDivider>
			<div data-testid="editor-team-composer">
				<TeamCard team={parsedTeam} />
			</div>

			<SectionDivider>{t("simple.action_list")}</SectionDivider>

			<Editor
				value={config}
				onChange={setConfig}
				error={error}
				importedCharacters={importedCharacters}
			/>

			<div className="sticky bottom-0 z-10 mt-1 flex flex-row flex-wrap place-items-center gap-1 bg-g-canvas p-2">
				<div className="flex flex-grow basis-full items-center p-1 sm:basis-0">
					{settings}
				</div>
				<div className="flex basis-full flex-row flex-wrap gap-1 p-1 sm:basis-2/3">
					<Button className="flex-1" onClick={onRun} disabled={!canRun}>
						{busy ? <Spinner /> : <Play />}
						{t("simple.run")}
					</Button>
				</div>
			</div>
		</div>
	);
}
