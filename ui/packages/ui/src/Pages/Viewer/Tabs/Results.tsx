import { Card, Colors } from "@blueprintjs/core";
import {
	CharacterActionsBarChart,
	CharacterDPSBarChart,
	CharacterDPSCard,
	CumulativeDamageCard,
	DamageTimelineCard,
	ElementDPSCard,
	EndingEnergyBarChart,
	FieldTimeCard,
	SourceDPSBarChart,
	SourceReactionsBarChart,
	TargetAuraUptimeBarChart,
	TargetDPSCard,
	TotalSourceEnergyBarChart,
} from "@gcsim/components";
import type { model, SimResults } from "@gcsim/types";
import classNames from "classnames";
import { type ReactNode, useEffect, useRef } from "react";
import { FiLink2 } from "react-icons/fi";
import { useLocation } from "react-router";
import {
	DistributionCard,
	RollupCards,
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
			{/* <Overview {...props} />
      <Damage {...props} /> */}
			{/* <Energy {...props} />
      <Reactions {...props} />
      <Healing {...props} />
      <Shields {...props} />
      <SimDetails {...props} /> */}
		</div>
	);
};

const SingleGroup = ({ data, modelData, running, names }: Props) => (
	<Group>
		<TeamHeader characters={data?.character_details} />
		<Metadata modelData={modelData} />
		<RollupCards data={data} />
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

const Overview = ({ data, modelData }: Props) => (
	<Group>
		<TeamHeader characters={data?.character_details} />
		<Metadata modelData={modelData} />
		<RollupCards data={data} />
		<TargetInfo
			enemies={modelData?.target_details}
			player={modelData?.player_position}
		/>
		<DistributionCard modelData={modelData} />
	</Group>
);

const Damage = ({ modelData, running, names }: Props) => (
	<Group>
		<Heading text="Damage" target="damage" color={Colors.VERMILION5} />
		<DamageTimelineCard data={modelData} running={running} names={names} />

		<CharacterDPSCard data={modelData} running={running} names={names} />
		<ElementDPSCard data={modelData} running={running} />
		<TargetDPSCard data={modelData} running={running} />

		<CharacterDPSBarChart data={modelData} running={running} names={names} />

		{/* <Card className="flex col-span-full h-64 min-h-full">
      Damage breakdown table(s)
    </Card> */}
	</Group>
);

const Energy = (_props: Props) => (
	<Group>
		<Heading text="Energy" target="energy" color={Colors.CERULEAN5} />
		<Card className="flex col-span-full h-64 min-h-full">
			Energy over time + cumu gained + cumu wasted
		</Card>
		<Card className="flex col-span-full h-64 min-h-full">
			Energy produced by source
		</Card>
		<Card className="flex col-span-full h-64 min-h-full">
			Incoming energy per character breakdown
		</Card>
	</Group>
);

const Reactions = (_props: Props) => (
	<Group>
		<Heading
			text="Reactions & Auras"
			target="reactions"
			color={Colors.VIOLET5}
		/>
		<Card className="flex col-span-4 h-64 min-h-full">
			Aura uptime timeline (worst, best, heatmap)
		</Card>
		<Card className="flex col-span-2 h-64 min-h-full">
			Aura uptime (pie? vs bar?)
		</Card>
		<Card className="flex col-span-3 h-64 min-h-full">
			Reactions triggered bar chart
		</Card>
		<Card className="flex col-span-3 h-64 min-h-full">
			Reactions by source?
		</Card>
	</Group>
);

const Healing = (_props: Props) => (
	<Group>
		<Heading text="Healing" target="healing" color={Colors.FOREST5} />
		<Card className="flex col-span-full h-64 min-h-full">
			Effective Healing Timeline (+ HP per char?)
		</Card>
		<Card className="flex col-span-3 h-64 min-h-full">Healing by Src</Card>
		<Card className="flex col-span-3 h-64 min-h-full">Healing by target</Card>
	</Group>
);

const Shields = (_props: Props) => (
	<Group>
		<Heading text="Shields" target="shields" color={Colors.GOLD5} />
		<Card className="flex col-span-full h-64 min-h-full">
			Shield timeline (worst, best, heatmap)
		</Card>
		<Card className="flex col-span-3 h-[512px] min-h-full">
			Shield uptime bar chart + pie
		</Card>
		<Card className="flex col-span-3 h-[512px] min-h-full">
			Shield hp bar chart
		</Card>
		<Card className="flex col-span-full h-64 min-h-full">Shield table</Card>
	</Group>
);

const SimDetails = (_props: Props) => (
	<Group>
		<Heading text="Simulation Details" target="sim" color={Colors.TURQUOISE5} />
		<Card className="flex col-span-2 h-64 min-h-full">
			Character uptime (pie?)
		</Card>
		<Card className="flex col-span-4 h-64 min-h-full">
			Character Uptime timeline (worst, best, heatmap)
		</Card>
		<Card className="flex col-span-4 h-64 min-h-full">
			Failed actions bar graph
		</Card>
		<Card className="flex col-span-2 h-64 min-h-full">
			Faied actions timeline (worst, best, heatmap)
		</Card>
		{/* tables? */}
	</Group>
);

type HeadingProps = {
	text: string;
	target: string;
	color?: string;
};

const Heading = ({ text, target, color }: HeadingProps) => {
	const linkClass = classNames(
		"ml-3 mt-1",
		"text-blue-500",
		"opacity-0 group-hover:opacity-100 transition-opacity",
		"flex justify-center items-center",
	);

	return (
		<h2 className="group flex whitespace-pre-wrap col-span-full text-xl font-semibold mt-12 mb-1">
			<span style={{ color: color }}>{text}</span>
			<a href={"#" + target} id={target} className={linkClass}>
				<FiLink2 />
			</a>
		</h2>
	);
};

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
