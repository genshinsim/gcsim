import type { model } from "@gcsim/types";
import type { JSX } from "react";
import { useTranslation } from "react-i18next";
import {
	type TooltipFormat,
	TooltipList,
	TooltipRow,
	TooltipTitle,
} from "./TooltipRows";

type Props = {
	title: string | JSX.Element;
	data: model.DescriptiveStats;
	color?: string;
	percent?: number;
	format?: TooltipFormat;
};

export default ({ title, data, color, percent, format }: Props) => {
	const { t } = useTranslation();
	return (
		<div className="flex flex-col font-g-mono text-g-xs">
			<TooltipTitle title={title} color={color} percent={percent} />
			<TooltipList>
				<TooltipRow
					format={format}
					color={color}
					name={t("result.stat_mean")}
					value={data.mean}
				/>
				<TooltipRow
					format={format}
					color={color}
					name={t("result.stat_min")}
					value={data.min}
				/>
				<TooltipRow
					format={format}
					color={color}
					name={t("result.stat_max")}
					value={data.max}
				/>
				<TooltipRow
					format={format}
					color={color}
					name={t("result.stat_std")}
					value={data.sd}
				/>
			</TooltipList>
		</div>
	);
};
