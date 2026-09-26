import { dynamicKey } from "@gcsim/localization";
import type { model } from "@gcsim/types";
import { Plus } from "lucide-react";
import React, { type JSX } from "react";
import { useTranslation } from "react-i18next";
import { CharacterCard } from "../CharacterCard/CharacterCard";
import { ConsolidateCharStats } from "./charStats";

type Props = {
	team: model.Character[];
	handleRemove: (index: number) => () => void;
	handleAdd?: () => void;
};

export const TeamCard = (props: Props) => {
	const { t } = useTranslation();
	const [showDetails, setShowDetails] = React.useState(false);
	const [showSnapshot, setShowSnapshot] = React.useState(false);
	const teamStats = ConsolidateCharStats(t, props.team);

	const handleToggleDetail = () => {
		setShowDetails(!showDetails);
	};
	const handleToggleSnapshot = () => {
		setShowSnapshot(!showSnapshot);
	};

	const cards: JSX.Element[] = props.team.map((c, index) => {
		const name = c.name ?? "";
		return (
			<CharacterCard
				key={name || index}
				char={c}
				stats={teamStats.stats[name]}
				snapshot={teamStats.snapshot[name]}
				statsRows={teamStats.maxRows}
				name={t(dynamicKey(`game:character_names.${name}`))}
				constellationLabel={`${t("character.c_pre")}${c.cons ?? 0}${t("character.c_post")}`}
				levelLabel={t("character.lvl")}
				talentsLabel={t("character.talents")}
				artifactStatsLabel={t("character.artifact_stats")}
				totalStatsLabel={t("character.total_stats")}
				weaponName={t(dynamicKey(`game:weapon_names.${c.weapon?.name ?? ""}`))}
				handleToggleDetail={handleToggleDetail}
				handleToggleSnapshot={handleToggleSnapshot}
				showDetails={showDetails}
				showSnapshot={showSnapshot}
				handleDelete={props.handleRemove(index)}
				className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2 pb-2"
			/>
		);
	});

	if (props.handleAdd && cards.length < 4) {
		cards.push(
			<div
				className="basis-full sm:basis-1/2 hd:basis-1/4 pr-2 pb-2 pt-2"
				key="_blank"
			>
				<button
					type="button"
					aria-label={t("db.characters")}
					className="bg-g-surface-2 rounded-g-md hover:bg-g-surface-3 flex items-center justify-center min-h-[226px] h-full w-full"
					onClick={props.handleAdd}
				>
					<Plus size={30} color="var(--g-text-mute)" />
				</button>
			</div>,
		);
	}

	return <div className="flex flex-row flex-wrap pl-2">{cards}</div>;
};
