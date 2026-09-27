import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { useDeferredValue } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, ParentWidth } from "../../../../common/gcsim";
import { BarChart, BarChartLegend } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	names?: string[];
};

export default ({ data, names }: Props) => {
	const { t } = useTranslation();
	const deferred = useDeferredValue(data);
	const energy = deferred?.statistics?.total_source_energy;

	return (
		<Card className="flex flex-col col-span-full h-auto p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle
						title={t("result.per_source", {
							s: t("result.total_energy"),
						})}
					/>
				</div>
				<div className="flex flex-grow justify-start sm:justify-center pb-5 sm:pb-0 items-center">
					<BarChartLegend names={names} />
				</div>
			</div>
			<ParentWidth>
				{(width) => <BarChart width={width} energy={energy} names={names} />}
			</ParentWidth>
		</Card>
	);
};
