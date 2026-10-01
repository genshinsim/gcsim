import type { model } from "@gcsim/types";
import { LegendOrdinal } from "@visx/legend";
import { scaleOrdinal } from "@visx/scale";
import { range } from "lodash-es";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import {
	FloatStatTooltipContent,
	HorizontalBarStack,
	NoData,
	useDataColors,
} from "../../../../common/gcsim";

type Props = {
	width: number;
	height: number;
	names?: string[];
	dps?: model.DescriptiveStats[];
};

export const ByCharacterLegend = ({ names }: { names?: string[] }) => {
	const { DataColors } = useDataColors();
	if (names == null) {
		return null;
	}

	const scale = scaleOrdinal({
		domain: names,
		range: names.map((_, i) => DataColors.character(i)),
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

export const ByCharacterChart = ({ width, height, names, dps }: Props) => {
	const { DataColors } = useDataColors();
	const { t } = useTranslation();
	const { data, keys, xMax } = useMemo(
		() => characterChartData(dps, names),
		[dps, names],
	);

	if (dps == null || names == null || keys.length === 0) {
		return <NoData />;
	}

	return (
		<HorizontalBarStack<CharacterData, number>
			width={width}
			height={height}
			xDomain={[0, xMax]}
			yDomain={names}
			y={(d) => d.name}
			data={data}
			keys={keys}
			value={(d, k) => (d.index === k ? (d.data.mean ?? 0) : 0)}
			stat={(d) => d.data}
			barColor={(k) => DataColors.character(k)}
			hoverColor={(k) => DataColors.characterLabel(k)}
			tooltipContent={(d, k) => (
				<FloatStatTooltipContent
					title={t("result.dps_title", { name: d.name })}
					data={d.data}
					color={DataColors.characterLabel(k)}
					percent={d.share}
				/>
			)}
		/>
	);
};

type CharacterData = {
	name: string;
	data: model.DescriptiveStats;
	index: number;
	share: number;
};

type ChartData = {
	data: CharacterData[];
	keys: number[];
	xMax: number;
};

export function characterChartData(
	dps?: model.DescriptiveStats[],
	names?: string[],
): ChartData {
	if (dps == null || names == null) {
		return { data: [], keys: [], xMax: 0 };
	}

	const total = dps.reduce((sum, v) => sum + (v.mean ?? 0), 0);
	let maxDPS = 0;
	const data: CharacterData[] = dps.map((v, i) => {
		const charMax = Math.max(v.max ?? 0, (v.mean ?? 0) + (v.sd ?? 0));
		maxDPS = Math.max(maxDPS, charMax);
		return { name: names[i], data: v, index: i, share: (v.mean ?? 0) / total };
	});

	return { data: data, keys: range(names.length), xMax: maxDPS };
}
