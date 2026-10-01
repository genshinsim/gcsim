import { useTranslation } from "react-i18next";
import { DataColorsConst } from "../DataColors";
import { TooltipList, TooltipRow, TooltipTitle } from "../TooltipRows";

export interface TooltipData {
	player: boolean;
	index: number;
	x: number;
	y: number;
	r: number;
}

export const PositionTooltipContent = ({ data }: { data: TooltipData }) => {
	const { t } = useTranslation();

	const title = data.player
		? t("result.player")
		: `${t("viewer.target")} ${data.index + 1}`;
	const titleColor = data.player
		? DataColorsConst.gray
		: DataColorsConst.qualitative3(data.index);

	return (
		<div className="flex flex-col font-g-mono text-g-xs">
			<TooltipTitle title={title} color={titleColor} />
			<TooltipList>
				<TooltipRow name="x" value={data.x} />
				<TooltipRow name="y" value={data.y} />
				<TooltipRow name="r" value={data.r} />
			</TooltipList>
		</div>
	);
};
