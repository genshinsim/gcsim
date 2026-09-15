import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { useMemo } from "react";
import {
	FloatStatTooltipContent,
	HorizontalBarStack,
	NoData,
	useDataColors,
} from "../../../../common/gcsim";
import { type StatMap, type StatMapDatum, useStatMapData } from "./StatMap";

type Props = {
	width: number;
	height: number;
	names?: string[];
	dps?: StatMap[];
};

export const ByElementLegend = ({ dps }: { dps?: StatMap[] }) => {
	const { DataColors } = useDataColors();
	const keys = useMemo(() => {
		if (dps == null) {
			return [];
		}

		const elements = new Set<string>();
		for (let i = 0; i < dps.length; i++) {
			for (const key in dps[i]) {
				elements.add(key);
			}
		}
		return Array.from(elements);
	}, [dps]);

	const scale = scaleOrdinal({
		domain: keys,
		range: keys.map((v) => DataColors.element(v)),
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

export const ByElementChart = ({ width, height, names, dps }: Props) => {
	const { DataColors } = useDataColors();
	const { data, keys, xMax } = useStatMapData(dps, names);

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
			barColor={DataColors.element}
			hoverColor={DataColors.elementLabel}
			tooltipContent={(d, k) => (
				<FloatStatTooltipContent
					title={`${d.name} ${k} DPS`}
					data={d.data[k]}
					color={DataColors.elementLabel(k)}
					percent={(d.data[k].mean ?? 0) / d.total}
				/>
			)}
		/>
	);
};
