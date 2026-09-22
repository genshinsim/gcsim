import {
	Alert,
	AlertDescription,
	AlertTitle,
	CommandItem,
} from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import React from "react";
import { useTranslation } from "react-i18next";
import { TeamCard } from "../Cards";
import { characterLabel, characters, OmniSelect } from "../common/gcsim";
import { cfgFromTeam } from "./teamConfig";
import type { TeamComposerCharacterSource } from "./types";

export interface TeamComposerProps {
	parsedTeam: model.Character[];
	error: string | null;
	config: string;
	setConfig: (v: string) => void;
	characters?: TeamComposerCharacterSource;
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

export function TeamComposer({
	parsedTeam,
	error,
	config,
	setConfig,
	characters: source,
}: TeamComposerProps) {
	const { t } = useTranslation();
	const [pickerOpen, setPickerOpen] = React.useState(false);

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
				items.push({
					key: k,
					source: "default",
					text: characterLabel(k),
					label: "",
				});
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
		<div data-testid="editor-team-composer" className="flex flex-col gap-2">
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

			<TeamCard
				team={parsedTeam}
				handleRemove={handleRemove}
				handleAdd={source ? () => setPickerOpen(true) : undefined}
			/>

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
