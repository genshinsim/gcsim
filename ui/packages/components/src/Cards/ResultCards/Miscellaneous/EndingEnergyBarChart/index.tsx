import { Card } from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useTranslation } from "react-i18next";
import { CardTitle, useRefreshWithTimer } from "../../../../common/gcsim";
import { BarChart } from "./BarChart";

type Props = {
	data: model.SimulationResult | null;
	running: boolean;
	names?: string[];
};

export default ({ data, running, names }: Props) => {
	const { t } = useTranslation();
	const [stats, timer] = useRefreshWithTimer(
		(d) => {
			return {
				data: d?.statistics?.end_stats,
			};
		},
		5000,
		data,
		running,
	);

	return (
		<Card className="flex flex-col col-span-full h-auto p-5">
			<div className="flex flex-col sm:flex-row justify-start gap-5">
				<div className="flex flex-col gap-2">
					<CardTitle title={t("result.ending_energy")} timer={timer} />
				</div>
			</div>
			<ParentSize>
				{({ width, height }) => (
					<BarChart
						width={width}
						height={height}
						end_stats={stats.data}
						names={names}
					/>
				)}
			</ParentSize>
		</Card>
	);
};
