import { dynamicKey } from "@gcsim/localization";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useTranslation } from "react-i18next";
import {
	CardTitle,
	useDataColors,
	useRefreshWithTimer,
} from "../../../../common/gcsim";
import { BarChart, BarChartLegend } from "./BarChart";

function canonicalActionOrder(
	actions: model.SourceStats[],
	actionKeys: string[],
): string[] {
	const names = new Set(
		actions.flatMap((a) => (a.sources ? Object.keys(a.sources) : [])),
	);
	return [...names].sort(
		(a, b) => actionKeys.indexOf(a) - actionKeys.indexOf(b),
	);
}

type Props = {
	data: model.SimulationResult | null;
	running: boolean;
	names?: string[];
};

export default ({ data, running, names }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();

	const [actions, timer] = useRefreshWithTimer(
		(d) =>
			d?.statistics?.character_actions?.map((s) =>
				s.sources
					? {
							sources: Object.fromEntries(
								Object.entries(s.sources).map(([k, v]) => [
									t(dynamicKey(`actions.${k}`)),
									v,
								]),
							),
						}
					: {},
			),
		5000,
		data,
		running,
	);

	const actionNames = actions
		? canonicalActionOrder(actions, DataColors.actionKeys)
		: null;

	return (
		<Card className="flex flex-col col-span-3 min-h-[384px] p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle title={t("simple.actions")} timer={timer} />
				</div>
				<div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
					<BarChartLegend actionNames={actionNames} />
				</div>
			</div>
			<ParentSize className="flex-grow">
				{({ width, height }) => (
					<BarChart
						width={width}
						height={height}
						actions={actions}
						names={names}
						actionNames={actionNames}
					/>
				)}
			</ParentSize>
		</Card>
	);
};
