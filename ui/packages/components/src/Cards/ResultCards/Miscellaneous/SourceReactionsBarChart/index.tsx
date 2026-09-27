import { dynamicKey } from "@gcsim/localization";
import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useDeferredValue, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle } from "../../../../common/gcsim";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	names?: string[];
};

export default ({ data, names }: Props) => {
	const { t } = useTranslation();
	const deferred = useDeferredValue(data);
	const sourceReactions = deferred?.statistics?.source_reactions;
	const reactions = useMemo(
		() =>
			sourceReactions?.map((s) =>
				s.sources
					? {
							sources: Object.fromEntries(
								Object.entries(s.sources).map(([k, v]) => [
									t(dynamicKey("reactions." + k)),
									v,
								]),
							),
						}
					: {},
			),
		[sourceReactions, t],
	);

	return (
		<Card className="flex flex-col col-span-3 min-h-[384px] p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle
						title={t("result.per_source", {
							s: t("result.reactions"),
						})}
					/>
				</div>
				<div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
					<BarChartLegend names={names} />
				</div>
			</div>
			<ParentSize className="flex-grow">
				{({ width, height }) => (
					<BarChart
						width={width}
						height={height}
						reactions={reactions}
						names={names}
					/>
				)}
			</ParentSize>
		</Card>
	);
};
