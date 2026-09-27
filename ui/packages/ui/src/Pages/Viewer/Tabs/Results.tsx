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
} from "@gcsim/components";
import type { model } from "@gcsim/types";
import classNames from "classnames";
import { type ReactNode, useEffect, useRef } from "react";
import { useLocation } from "react-router";
import {
	DistributionCard,
	TargetInfo,
	TeamHeader,
} from "../Components/Overview";
import Metadata from "../Components/Overview/Metadata";

type Props = {
	data: model.SimulationResult | null;
	running: boolean;
	names?: string[];
};

export default (props: Props) => {
	useScrollToLocation();

	return (
		<div className="w-full 2xl:mx-auto 2xl:container px-2">
			<SingleGroup {...props} />
		</div>
	);
};

const SingleGroup = ({ data, running, names }: Props) => (
	<Group>
		<TeamHeader characters={data?.character_details} />
		<Metadata data={data} />
		<RollupCards data={data} />
		<TargetInfo enemies={data?.target_details} player={data?.player_position} />
		<DistributionCard data={data} />

		<DamageTimelineCard data={data} running={running} names={names} />
		<CumulativeDamageCard data={data} running={running} />

		<CharacterDPSCard data={data} running={running} names={names} />
		<ElementDPSCard data={data} running={running} />
		<TargetDPSCard data={data} running={running} />

		<CharacterDPSBarChart data={data} running={running} names={names} />

		<SourceDPSBarChart data={data} running={running} names={names} />

		<CharacterActionsBarChart data={data} running={running} names={names} />

		<FieldTimeCard data={data} running={running} names={names} />

		<TotalSourceEnergyBarChart data={data} running={running} names={names} />

		<EndingEnergyBarChart data={data} running={running} names={names} />

		<SourceReactionsBarChart data={data} running={running} names={names} />

		<TargetAuraUptimeBarChart data={data} running={running} />
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

function useScrollToLocation() {
	const scrolled = useRef(false);
	const { key, hash } = useLocation();
	const prevKey = useRef(key);

	useEffect(() => {
		if (hash == null) {
			return;
		}

		if (prevKey.current !== key) {
			prevKey.current = key;
			scrolled.current = false;
		}

		if (scrolled.current) {
			return;
		}
		const id = hash.replace("#", "");
		if (!id) {
			return;
		}
		const element = document.getElementById(id);
		if (element) {
			element.scrollIntoView({ behavior: "smooth" });
			scrolled.current = true;
		}
	});
}
