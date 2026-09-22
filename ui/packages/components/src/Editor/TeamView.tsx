import {
	Alert,
	AlertDescription,
	AlertTitle,
	CommandItem,
} from "@gcsim/primitives";
import { dynamicKey } from "@gcsim/localization";
import type { model } from "@gcsim/types";
import { Plus } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";
import { CharacterCard } from "../Cards";
import { OmniSelect, characterLabel, characters } from "../common/gcsim";
import { ConsolidateCharStats } from "./charStats";
import { cfgFromTeam } from "./teamConfig";
import type { TeamViewCharacterSource } from "./types";

export interface TeamViewProps {
	parsedTeam: model.Character[];
	error: string | null;
	config: string;
	setConfig: (v: string) => void;
	characters?: TeamViewCharacterSource;
}

type PickerItem = {
	key: string;
	source: "default" | "user";
	text: string;
	label: string;
};

const itemKey = (item: PickerItem) => `${item.source}-${item.key}`;

const itemPredicate = (item: PickerItem, query: string) => {
	const normalized = query.trim().toLowerCase();
	if (normalized.length === 0) {
		return true;
	}
	return `${item.label} ${item.key} ${item.text}`
		.toLowerCase()
		.includes(normalized);
};

export function TeamView({
	parsedTeam,
	error,
	config,
	setConfig,
	characters: source,
}: TeamViewProps) {
	const { t } = useTranslation();
	const [pickerOpen, setPickerOpen] = React.useState(false);
	const [showDetails, setShowDetails] = React.useState(false);
	const [showSnapshot, setShowSnapshot] = React.useState(false);

	const teamStats = ConsolidateCharStats(t, parsedTeam);
	const onTeam = new Set(parsedTeam.map((c) => c.name ?? ""));

	const writeTeam = (team: model.Character[]) => {
		setConfig(cfgFromTeam(team, config));
	};

	const handleRemove = (index: number) => () => {
		writeTeam(parsedTeam.filter((_, i) => i !== index));
	};

	const handleAdd = (item: PickerItem) => {
		if (!source) {
			return;
		}
		setPickerOpen(false);
		const character =
			item.source === "user"
				? structuredClone(
						source.imported?.find((c) => c.key === item.key)?.character,
					)
				: source.createCharacter(item.key);
		if (character) {
			writeTeam([...parsedTeam, character]);
		}
	};

	const items: PickerItem[] = [];
	if (source) {
		characters.forEach((k) => {
			if (!onTeam.has(k)) {
				items.push({ key: k, source: "default", text: characterLabel(k), label: "" });
			}
		});
		source.imported?.forEach((option) => {
			if (!onTeam.has(option.character.name ?? "")) {
				items.push({
					key: option.key,
					source: "user",
					text: characterLabel(option.character.name ?? ""),
					label: option.label ?? "",
				});
			}
		});
	}

	return (
		<div data-testid="editor-team-view" className="flex flex-col gap-2">
			{error ? (
				<Alert variant="destructive">
					<AlertTitle>
						{t("viewer.error_encountered") + t("viewer.config_invalid")}
					</AlertTitle>
					<AlertDescription>
						<pre className="whitespace-pre-wrap">{error}</pre>
					</AlertDescription>
				</Alert>
			) : null}

			<div className="flex flex-row flex-wrap">
				{parsedTeam.map((c, index) => (
					<CharacterCard
						key={c.name ?? index}
						char={c}
						stats={teamStats.stats[c.name ?? ""]}
						snapshot={teamStats.snapshot[c.name ?? ""]}
						statsRows={teamStats.maxRows}
						name={t(dynamicKey(`game:character_names.${c.name ?? ""}`))}
						constellationLabel={`${t("character.c_pre")}${c.cons ?? 0}${t("character.c_post")}`}
						levelLabel={t("character.lvl")}
						talentsLabel={t("character.talents")}
						artifactStatsLabel={t("character.artifact_stats")}
						totalStatsLabel={t("character.total_stats")}
						weaponName={t(dynamicKey(`game:weapon_names.${c.weapon?.name ?? ""}`))}
						handleToggleDetail={() => setShowDetails((v) => !v)}
						handleToggleSnapshot={() => setShowSnapshot((v) => !v)}
						showDetails={showDetails}
						showSnapshot={showSnapshot}
						handleDelete={handleRemove(index)}
						className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2 pb-2"
					/>
				))}

				{source && parsedTeam.length < 4 ? (
					<div className="basis-full sm:basis-1/2 hd:basis-1/4 pr-2 pb-2 pt-2">
						<button
							type="button"
							aria-label={t("db.characters")}
							className="bg-g-surface-2 rounded-g-md hover:bg-g-surface-3 flex items-center justify-center min-h-[226px] h-full w-full"
							onClick={() => setPickerOpen(true)}
						>
							<Plus size={30} color="var(--g-text-mute)" />
						</button>
					</div>
				) : null}
			</div>

			{source ? (
				<OmniSelect<PickerItem>
					isOpen={pickerOpen}
					onClose={() => setPickerOpen(false)}
					items={items}
					itemKey={itemKey}
					itemPredicate={itemPredicate}
					itemRenderer={(item, state) => (
						<CommandItem value={itemKey(item)} onSelect={state.onSelect}>
							<span className="flex-1">{item.text}</span>
							{item.label ? (
								<span className="text-g-ink-mute text-g-xs">{item.label}</span>
							) : null}
						</CommandItem>
					)}
					onSelect={handleAdd}
					title={t("db.characters")}
					placeholder={t("db.type_to_search")}
				/>
			) : null}
		</div>
	);
}
