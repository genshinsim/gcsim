import type { model } from "@gcsim/types";
import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Colors, useRefresh } from "../../../common/gcsim";
import RollupCard from "./RollupCard";

const ROLLUPS = [
	{
		statKey: "dps",
		titleKey: "result.dps_long",
		abbrev: "DPS",
		color: Colors.VERMILION3,
		fractionDigits: 0,
	},
	{
		statKey: "eps",
		titleKey: "result.eps_long",
		abbrev: "EPS",
		color: Colors.CERULEAN3,
		fractionDigits: 2,
	},
	{
		statKey: "rps",
		titleKey: "result.rps_long",
		abbrev: "RPS",
		color: Colors.VIOLET3,
		fractionDigits: 2,
	},
	{
		statKey: "hps",
		titleKey: "result.hps_long",
		abbrev: "HPS",
		color: Colors.FOREST3,
		fractionDigits: 0,
	},
	{
		statKey: "shp",
		titleKey: "result.shp_long",
		abbrev: "SHP",
		color: Colors.GOLD3,
		fractionDigits: 0,
	},
	{
		statKey: "duration",
		titleKey: "result.dur_long",
		abbrev: "Dur",
		color: Colors.TURQUOISE3,
		fractionDigits: 2,
		labelKey: "result.seconds_short",
	},
] as const;

type RollupConfig = (typeof ROLLUPS)[number];

type Props = {
	data: model.SimulationResult | null;
};

export default ({ data }: Props) => (
	<div className="col-span-full flex flex-row flex-wrap gap-2 justify-center">
		{ROLLUPS.map((config) => (
			<StatRollupCard key={config.statKey} data={data} config={config} />
		))}
	</div>
);

const StatRollupCard = ({
	data,
	config,
}: {
	data: model.SimulationResult | null;
	config: RollupConfig;
}) => {
	const { i18n, t } = useTranslation();
	const fmt = useCallback(
		(val?: number) =>
			val?.toLocaleString(i18n.language, {
				maximumFractionDigits: config.fractionDigits,
			}),
		[i18n, config.fractionDigits],
	);

	const stat = useRefresh((d) => d?.statistics?.[config.statKey], 200, data);
	const auxStats = useMemo(
		() => [
			{ title: "min", value: fmt(stat?.min) },
			{ title: "max", value: fmt(stat?.max) },
			{ title: "std", value: fmt(stat?.sd) },
			{ title: "p25", value: fmt(stat?.q1) },
			{ title: "p50", value: fmt(stat?.q2) },
			{ title: "p75", value: fmt(stat?.q3) },
		],
		[stat, fmt],
	);

	return (
		<RollupCard
			color={config.color}
			title={`${t(config.titleKey)} (${config.abbrev})`}
			label={"labelKey" in config ? t(config.labelKey) : undefined}
			value={fmt(stat?.mean)}
			auxStats={auxStats}
		/>
	);
};
