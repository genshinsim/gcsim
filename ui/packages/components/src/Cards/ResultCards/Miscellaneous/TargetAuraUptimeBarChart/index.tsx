import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
	CardTitle,
	useDataColors,
	useRefreshWithTimer,
} from "../../../../common/gcsim";
import { Card } from "../../../../common/ui/card";
import { BarChart } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	running: boolean;
};

export default ({ data, running }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();
	const [stats, timer] = useRefreshWithTimer(
		(d) => {
			return {
				data: d?.statistics?.target_aura_uptime
					? d?.statistics?.target_aura_uptime.map((s) =>
							s.sources
								? {
										sources: Object.fromEntries(
											Object.entries(s.sources).map(([k, v]) => [
												t("elements." + k),
												v,
											]),
										),
									}
								: {},
						)
					: undefined,
			};
		},
		5000,
		data,
		running,
	);

	const targets = useMemo(() => {
		if (stats.data == null) {
			return [];
		}

		const targets = new Set<string>();
		for (let i = 0; i < stats.data.length; i++) {
			targets.add(i.toString());
		}
		return Array.from(targets);
	}, [stats.data]);
	const [target, setTarget] = useState("0");

	const auras = useMemo(() => {
		if (stats.data == null) {
			return [];
		}

		const auras = new Set<string>();
		for (const key in stats.data[target]?.sources) {
			auras.add(key);
		}
		return Array.from(auras).sort(
			(a, b) =>
				DataColors.reactableModifierKeys.indexOf(a) -
				DataColors.reactableModifierKeys.indexOf(b),
		);
	}, [stats.data, DataColors.reactableModifierKeys, target]);

	return (
		<Card className="flex flex-col col-span-3 min-h-[384px] p-5">
			<div className="flex flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle title={t("result.target_aura_uptime")} timer={timer} />
					<Options target={target} setTarget={setTarget} targets={targets} />
				</div>
			</div>
			<ParentSize className="flex-grow">
				{({ width, height }) => (
					<BarChart
						width={width}
						height={height}
						auraUptime={stats.data}
						auras={auras}
						target={target}
					/>
				)}
			</ParentSize>
		</Card>
	);
};

const Options = ({
	target,
	setTarget,
	targets,
}: {
	target: string;
	setTarget: (v: string) => void;
	targets: string[];
}) => {
	const { t } = useTranslation();

	return (
		<div className="flex flex-row items-center gap-2 mb-2">
			<span className="text-xs font-mono text-gray-400">
				{t("viewer.target")}
			</span>
			<select
				className="rounded border border-gray-500 bg-transparent px-2 py-1 text-sm"
				value={target}
				onChange={(e) => setTarget(e.target.value)}
			>
				{targets.map((target) => (
					<option key={target} value={target}>
						{Number(target) + 1}
					</option>
				))}
			</select>
		</div>
	);
};
