import type { model } from "@gcsim/types";
import classNames from "classnames";
import type { ReactNode } from "react";
import {
	CharacterActionsBarChart,
	CharacterDPSBarChart,
	CharacterDPSCard,
	CumulativeDamageCard,
	DamageTimelineCard,
	ElementDPSCard,
	EndingEnergyBarChart,
	FieldTimeCard,
	RollupCards,
	SourceDPSBarChart,
	SourceReactionsBarChart,
	TargetAuraUptimeBarChart,
	TargetDPSCard,
	TotalSourceEnergyBarChart,
} from "../Cards";
import DistributionCard from "./DistributionCard";
import Metadata from "./Metadata";
import TargetInfo from "./TargetInfo";
import TeamHeader from "./TeamHeader";

export type ResultsViewProps = {
	model: model.SimulationResult | null;
	names?: string[];
};

export const ResultsView = (props: ResultsViewProps) => (
	<div className="w-full 2xl:mx-auto 2xl:container px-2">
		<SingleGroup {...props} />
	</div>
);

const SingleGroup = ({ model: data, names }: ResultsViewProps) => (
	<Group>
		<TeamHeader characters={data?.character_details} />
		<Metadata data={data} />
		<RollupCards data={data} />
		<TargetInfo enemies={data?.target_details} player={data?.player_position} />
		<DistributionCard data={data} />

		<DamageTimelineCard data={data} names={names} />
		<CumulativeDamageCard data={data} />

		<CharacterDPSCard data={data} names={names} />
		<ElementDPSCard data={data} />
		<TargetDPSCard data={data} />

		<CharacterDPSBarChart data={data} names={names} />

		<SourceDPSBarChart data={data} names={names} />

		<CharacterActionsBarChart data={data} names={names} />

		<FieldTimeCard data={data} names={names} />

		<TotalSourceEnergyBarChart data={data} names={names} />

		<EndingEnergyBarChart data={data} names={names} />

		<SourceReactionsBarChart data={data} names={names} />

		<TargetAuraUptimeBarChart data={data} />
	</Group>
);

type GroupProps = {
	children: ReactNode;
	className?: string;
};

const Group = ({ children, className }: GroupProps) => {
	const cls = classNames(
		className,
		"grid overflow-hidden",
		"grid-cols-2 sm:grid-cols-6",
		"gap-y-2",
		"sm:gap-2",
	);

	return <div className={cls}>{children}</div>;
};
