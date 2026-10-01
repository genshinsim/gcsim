import type { model } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { TooltipList, TooltipRow } from "../../common/gcsim";

type Props = {
	idx: number;
	delta: number;
	data: model.OverviewStats;
};

export const HistogramTooltipContent = ({ idx, delta, data }: Props) => {
	const { i18n, t } = useTranslation();
	const min = data.min ?? 0;
	const max = data.max ?? 0;
	const count = data.histogram?.[idx] ?? 0;

	const lower = delta === 0 ? min : min + idx / delta;
	const upper = delta === 0 ? max : min + (idx + 1) / delta;
	const binOf = (v?: number) =>
		v == null ? null : Math.floor(delta * (v - min));

	const stats = [
		{ name: t("result.stat_mean"), value: data.mean },
		{ name: t("result.stat_p25"), value: data.q1 },
		{ name: t("result.stat_p50"), value: data.q2 },
		{ name: t("result.stat_p75"), value: data.q3 },
	].filter((s) => binOf(s.value) === idx);

	return (
		<div className="flex flex-col font-g-mono text-g-xs">
			<TooltipList>
				{stats.map((s) => (
					<TooltipRow key={s.name} name={s.name} value={s.value} />
				))}
				<TooltipRow name={t("result.lower")} value={lower} />
				<TooltipRow name={t("result.upper")} value={upper} />
				<TooltipRow
					name={t("result.iterations_short")}
					value={count}
					format={(n) => n?.toLocaleString(i18n.language)}
				/>
			</TooltipList>
		</div>
	);
};
