import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { useDeferredValue } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle, ParentWidth } from "../../../../common/gcsim";
import { BarChart } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	names?: string[];
};

export default ({ data, names }: Props) => {
	const { t } = useTranslation();
	const deferred = useDeferredValue(data);
	const endStats = deferred?.statistics?.end_stats;

	return (
		<Card className="flex flex-col col-span-full h-auto p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle
						title={t("result.ending_energy")}
						stale={deferred !== data}
					/>
				</div>
			</div>
			<ParentWidth>
				{(width) => (
					<BarChart width={width} end_stats={endStats} names={names} />
				)}
			</ParentWidth>
		</Card>
	);
};
