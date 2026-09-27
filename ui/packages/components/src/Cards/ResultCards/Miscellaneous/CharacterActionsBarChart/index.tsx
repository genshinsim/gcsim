import { dynamicKey } from "@gcsim/localization";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useDeferredValue, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, useDataColors } from "../../../../common/gcsim";
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
	names?: string[];
};

export default ({ data, names }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();

	const deferred = useDeferredValue(data);
	const characterActions = deferred?.statistics?.character_actions;
	const actions = useMemo(
		() =>
			characterActions?.map((s) =>
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
		[characterActions, t],
	);

	const actionNames = useMemo(
		() =>
			actions ? canonicalActionOrder(actions, DataColors.actionKeys) : null,
		[actions, DataColors.actionKeys],
	);

	return (
		<Card className="flex flex-col col-span-3 min-h-[384px] p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle title={t("simple.actions")} stale={deferred !== data} />
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
