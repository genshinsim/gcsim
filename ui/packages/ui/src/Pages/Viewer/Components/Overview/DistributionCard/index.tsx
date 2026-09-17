import { Colors, HistogramGraph } from "@gcsim/components";
import {
	Card,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import type { model } from "@gcsim/types";
import { ParentSize } from "@visx/responsive";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { CardTitle } from "../../Util";

type Props = {
	modelData: model.SimulationResult | null;
};

type GraphConfig = {
	label: string;
	pick: (stats?: model.SimulationStatistics) => model.OverviewStats | undefined;
	barColor: string;
	accentColor: string;
	hoverColor: string;
};

const GRAPHS: Record<string, GraphConfig> = {
	dps: {
		label: "DPS",
		pick: (s) => s?.dps,
		barColor: Colors.VERMILION3,
		accentColor: Colors.VERMILION1,
		hoverColor: Colors.VERMILION5,
	},
	eps: {
		label: "EPS",
		pick: (s) => s?.eps,
		barColor: Colors.CERULEAN3,
		accentColor: Colors.CERULEAN1,
		hoverColor: Colors.CERULEAN5,
	},
	rps: {
		label: "RPS",
		pick: (s) => s?.rps,
		barColor: Colors.VIOLET3,
		accentColor: Colors.VIOLET1,
		hoverColor: Colors.VIOLET5,
	},
	hps: {
		label: "HPS",
		pick: (s) => s?.hps,
		barColor: Colors.FOREST3,
		accentColor: Colors.FOREST1,
		hoverColor: Colors.FOREST5,
	},
	shp: {
		label: "SHP",
		pick: (s) => s?.shp,
		barColor: Colors.GOLD3,
		accentColor: Colors.GOLD1,
		hoverColor: Colors.GOLD5,
	},
	dur: {
		label: "Dur",
		pick: (s) => s?.duration,
		barColor: Colors.TURQUOISE3,
		accentColor: Colors.TURQUOISE1,
		hoverColor: Colors.TURQUOISE5,
	},
};

export default ({ modelData }: Props) => {
	const { t } = useTranslation();
	const [graph, setGraph] = useState("dps");
	const cfg = GRAPHS[graph];
	const stats = modelData?.statistics;

	return (
		<Card className="col-span-3 min-h-full h-72 min-w-[280px] flex flex-col justify-start gap-2 p-5">
			<div className="flex flex-row justify-start">
				<Select value={graph} onValueChange={setGraph}>
					<SelectTrigger className="w-24">
						<SelectValue />
					</SelectTrigger>
					<SelectContent>
						{Object.entries(GRAPHS).map(([key, c]) => (
							<SelectItem key={key} value={key}>
								{c.label}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<div className="flex flex-grow justify-center items-end">
					{cfg != null && (
						<CardTitle
							title={t("result.dist", { d: cfg.label })}
							tooltip="test"
						/>
					)}
				</div>
			</div>
			<ParentSize>
				{({ width, height }) => (
					<HistogramGraph
						width={width}
						height={height}
						data={cfg?.pick(stats)}
						barColor={cfg?.barColor}
						accentColor={cfg?.accentColor}
						hoverColor={cfg?.hoverColor}
					/>
				)}
			</ParentSize>
		</Card>
	);
};
