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
import type { model, SimResults } from "@gcsim/types";
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
	data: SimResults | null;
	modelData: model.SimulationResult | null;
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

const SingleGroup = ({ data, modelData, running, names }: Props) => (
	<Group>
		<TeamHeader characters={data?.character_details} />
		<Metadata modelData={modelData} />
		<RollupCards data={modelData} />
		<TargetInfo
			enemies={modelData?.target_details}
			player={modelData?.player_position}
		/>
		<DistributionCard modelData={modelData} />

		<DamageTimelineCard data={modelData} running={running} names={names} />
		<CumulativeDamageCard data={modelData} running={running} />

		<CharacterDPSCard data={modelData} running={running} names={names} />
		<ElementDPSCard data={modelData} running={running} />
		<TargetDPSCard data={modelData} running={running} />

		<CharacterDPSBarChart data={modelData} running={running} names={names} />

		<SourceDPSBarChart data={modelData} running={running} names={names} />

		<CharacterActionsBarChart
			data={modelData}
			running={running}
			names={names}
		/>

		<FieldTimeCard data={modelData} running={running} names={names} />

		<TotalSourceEnergyBarChart
			data={modelData}
			running={running}
			names={names}
		/>

		<EndingEnergyBarChart data={modelData} running={running} names={names} />

		<SourceReactionsBarChart data={modelData} running={running} names={names} />

		<TargetAuraUptimeBarChart data={modelData} running={running} />
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
