import { OmniSelect } from "@gcsim/components";
import { dynamicKey } from "@gcsim/localization";
import { CommandItem } from "@gcsim/primitives";
import type { Character } from "@gcsim/types";
import React from "react";
import { useTranslation } from "react-i18next";
import { CharMap } from "../../Data";
import { appActions } from "../../Stores/appSlice";
import {
	type RootState,
	useAppDispatch,
	useAppSelector,
} from "../../Stores/store";
import { Builder } from "./Components/TeamBuilder/Builder";

type CharSource = "user" | "default";

interface TeamCharacterItem {
	key: string;
	charKey: string;
	source: CharSource;
	text: string;
	label: string;
}

const itemKey = (item: TeamCharacterItem) => `${item.source}-${item.key}`;

const itemPredicate = (item: TeamCharacterItem, query: string) => {
	const normalized = query.trim().toLowerCase();
	if (normalized.length === 0) {
		return true;
	}
	return `${item.label} ${item.key} ${item.text}`
		.toLowerCase()
		.includes(normalized);
};

function newCharFromKey(k: string): Character {
	return {
		name: k,
		level: 80,
		max_level: 90,
		element: CharMap[k].element,
		cons: 0,
		weapon: {
			name: "dullblade",
			refine: 1,
			level: 1,
			max_level: 20,
		},
		talents: {
			attack: 6,
			skill: 6,
			burst: 6,
		},
		stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
		snapshot: [
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		],
		sets: {},
	};
}

export function Team() {
	const { t } = useTranslation();

	const { imported, team } = useAppSelector((state: RootState) => {
		return {
			imported: state.user_data.GOODImport,
			team: state.app.team,
		};
	});
	const [open, setOpen] = React.useState<boolean>(false);
	const dispatch = useAppDispatch();

	const handleRemove = (index: number) => {
		return () => {
			dispatch(appActions.deleteCharacter({ index }));
		};
	};

	const handleAdd = (item: TeamCharacterItem) => {
		setOpen(false);
		if (item.source === "user") {
			const character: Character = JSON.parse(
				JSON.stringify(imported[item.key]),
			);
			dispatch(appActions.addCharacter({ character }));
		} else {
			dispatch(
				appActions.addCharacter({ character: newCharFromKey(item.key) }),
			);
		}
	};

	const onTeam = new Set(team.map((c) => c.name));

	const items: TeamCharacterItem[] = [];
	Object.keys(CharMap).forEach((k) => {
		if (onTeam.has(k)) {
			return;
		}
		items.push({
			key: k,
			charKey: k,
			source: "default",
			text: t(dynamicKey("game:character_names." + k)),
			label: "",
		});
	});
	Object.keys(imported).forEach((k) => {
		const e = imported[k];
		if (onTeam.has(e.name)) {
			return;
		}
		let label = e.enka_build_name !== undefined ? ` ${e.enka_build_name}` : "";
		label += e.date_added !== undefined ? ` (Imported on ${e.date_added})` : "";
		items.push({
			key: k,
			charKey: e.name,
			source: "user",
			text: t(dynamicKey("game:character_names." + e.name)),
			label,
		});
	});

	return (
		<div className="flex flex-col">
			<Builder
				team={team}
				handleAdd={() => setOpen(true)}
				handleRemove={handleRemove}
			/>

			<OmniSelect<TeamCharacterItem>
				isOpen={open}
				onClose={() => setOpen(false)}
				items={items}
				itemKey={itemKey}
				itemPredicate={itemPredicate}
				itemRenderer={(item, state) => (
					<CommandItem value={itemKey(item)} onSelect={state.onSelect}>
						<span className="flex-1">{item.text}</span>
						{item.label && (
							<span className="text-muted-foreground text-xs">
								{item.label}
							</span>
						)}
					</CommandItem>
				)}
				onSelect={handleAdd}
				title={t("db.characters")}
				placeholder={t("db.type_to_search")}
			/>
		</div>
	);
}
