import { GRAY_400, SLATE_700 } from "../colors";

// Serialize an SVG markup string into a data URI. Satori lays out <img> but not
// nested SVG primitives, so every chart is embedded as an <img> of its own SVG.
export function svgDataUri(svg: string): string {
	return `data:image/svg+xml,${encodeURIComponent(svg)}`;
}

type ChartImgProps = {
	svg: string;
	width: number;
	height: number;
};

// A rendered chart, embedded as a fixed-size image of its SVG string.
export const ChartImg = ({ svg, width, height }: ChartImgProps) => {
	return (
		<img
			src={svgDataUri(svg)}
			width={width}
			height={height}
			alt=""
			style={{ width, height }}
		/>
	);
};

type NoDataProps = {
	width: number;
	height: number;
};

// Placeholder shown when a chart has no data to render.
export const NoData = ({ width, height }: NoDataProps) => {
	return (
		<div
			style={{
				display: "flex",
				width,
				height,
				alignItems: "center",
				justifyContent: "center",
				color: GRAY_400,
				backgroundColor: SLATE_700,
				fontSize: 11,
			}}
		>
			no data
		</div>
	);
};
