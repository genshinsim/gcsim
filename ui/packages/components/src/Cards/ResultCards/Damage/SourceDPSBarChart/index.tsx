import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, useRefreshWithTimer } from "../../../../common/gcsim";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	running: boolean;
	names?: string[];
};

type Graphs = Map<string, string>;

export default ({ data, running, names }: Props) => {
	const { t } = useTranslation();
	const graphs: Graphs = new Map([
		["dps", "DPS"],
		["damage_instances", t("result.dmg_instances")],
	]);
	const [graph, setGraph] = useState("dps");

	const all_filter = t("result.all");
	const filters: string[] = [all_filter, ...(names ?? [])];
	const [filter, setFilter] = useState(all_filter);
	useResetFilterToAllOnLanguageChange(setFilter, all_filter);

	const [stats, timer] = useRefreshWithTimer(
		(d) => {
			return {
				dps: d?.statistics?.source_dps,
				damage_instances: d?.statistics?.source_damage_instances,
			};
		},
		5000,
		data,
		running,
	);

	const chart_data = graph === "dps" ? stats.dps : stats.damage_instances;

	return (
		<Card className="flex flex-col col-span-full min-h-96 p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle
						title={t("result.source", { s: graphs.get(graph) })}
						timer={timer}
					/>
					<div className="flex flex-row gap-4">
						<Options graph={graph} setGraph={setGraph} graphs={graphs} />
						<Filters filter={filter} setFilter={setFilter} filters={filters} />
					</div>
				</div>
				<div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
					<BarChartLegend names={names} />
				</div>
			</div>
			<ParentSize className="overflow-x-auto">
				{({ width }) => (
					<BarChart
						width={width}
						dps={chart_data}
						names={names}
						all_filter={all_filter}
						filter={filter}
					/>
				)}
			</ParentSize>
		</Card>
	);
};

function useResetFilterToAllOnLanguageChange(
	setFilter: (v: string) => void,
	all_filter: string,
) {
	useEffect(() => {
		setFilter(all_filter);
	}, [setFilter, all_filter]);
}

const Options = ({
	graph,
	setGraph,
	graphs,
}: {
	graph: string;
	setGraph: (v: string) => void;
	graphs: Graphs;
}) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-row items-center gap-2 mb-2">
			<span className="text-g-xs font-g-mono text-g-ink-mute">
				{t("result.type")}
			</span>
			<select
				className="rounded-g-md border border-g-line bg-transparent px-2 py-1 text-g-sm"
				value={graph}
				onChange={(e) => setGraph(e.target.value)}
			>
				{[...graphs.keys()].map((key) => (
					<option key={key} value={key}>
						{graphs.get(key)}
					</option>
				))}
			</select>
		</div>
	);
};

const Filters = ({
	filter,
	setFilter,
	filters,
}: {
	filter: string;
	setFilter: (v: string) => void;
	filters: string[];
}) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-row items-center gap-2 mb-2">
			<span className="text-g-xs font-g-mono text-g-ink-mute">
				{t("db.character")}
			</span>
			<select
				className="rounded-g-md border border-g-line bg-transparent px-2 py-1 text-g-sm"
				value={filter}
				onChange={(e) => setFilter(e.target.value)}
			>
				{[...filters].map((key) => (
					<option key={key} value={key}>
						{key}
					</option>
				))}
			</select>
		</div>
	);
};
