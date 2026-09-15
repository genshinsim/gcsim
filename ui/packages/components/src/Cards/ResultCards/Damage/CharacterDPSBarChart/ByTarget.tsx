import type { model } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
	FloatStatTooltipContent,
	HorizontalBarStack,
	NoData,
	useDataColors,
} from "../../../../common/gcsim";
import { type StatMapDatum, useStatMapData } from "./StatMap";

type Props = {
	width: number;
	height: number;
	names?: string[];
	dps?: model.TargetStats[];
};

export const ByTargetLegend = ({ dps }: { dps?: model.TargetStats[] }) => {
	const { DataColors } = useDataColors();
	const keys = useMemo(() => {
		if (dps == null) {
			return [];
		}

		const targets = new Set<string>();
		for (let i = 0; i < dps.length; i++) {
			for (const key in dps[i].targets) {
				targets.add(key);
			}
		}
		return Array.from(targets);
	}, [dps]);

	const scale = scaleOrdinal({
		domain: keys,
		range: keys.map((v) => DataColors.target(v)),
	});

	return (
		<LegendOrdinal
			scale={scale}
			direction="row"
			labelMargin="0 15px 0 0"
			className="flex-wrap"
		/>
	);
};

export const ByTargetChart = ({ width, height, names, dps }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();
	const { data, keys, xMax } = useStatMapData(
		dps?.map((d) => d.targets),
		names,
	);

	if (dps == null || names == null || keys.length === 0) {
		return <NoData />;
	}

	return (
		<HorizontalBarStack<StatMapDatum, string>
			width={width}
			height={height}
			xDomain={[0, xMax]}
			yDomain={names}
			y={(d) => d.name}
			data={data}
			keys={keys}
			value={(d, k) => (k in d.data ? (d.data[k].mean ?? 0) : 0)}
			stat={(d, k) => d.data[k]}
			barColor={(k) => DataColors.target(k)}
			hoverColor={(k) => DataColors.targetLabel(k)}
			tooltipContent={(d, k) => (
				<FloatStatTooltipContent
					title={`${d.name} ${t("viewer.target")} ${k} DPS`}
					data={d.data[k]}
					color={DataColors.targetLabel(k)}
					percent={(d.data[k].mean ?? 0) / d.total}
				/>
			)}
		/>
	);
};
