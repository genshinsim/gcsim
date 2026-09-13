import type { model } from "@gcsim/types";
import { HistogramChart } from "./charts/HistogramChart";
import { CharacterPie, ElementPie } from "./charts/PieChart";
import { TimelineChart } from "./charts/TimelineChart";
import { SLATE_700, SLATE_800 } from "./colors";
import { FONT_FAMILY } from "./fonts";
import { Metadata } from "./Metadata";
import { Portraits } from "./Portraits";

// Default asset host, matching the live asset URL scheme. Absolute so Satori
// (which has no dev-server proxy) can fetch the images at render time.
export const DEFAULT_ASSET_BASE = "https://gcsim.app/api/assets";

export type SatoriPreviewCardProps = {
	data: model.SimulationResult;
	/** Injectable asset-base seam. Default: the live gcsim asset host. */
	assetBase?: string;
};

// Fixed card geometry (540x250), mirroring the live PreviewCard's measured
// layout pixel-for-pixel: a 4px inset around each row (its `m-1`), 4-across
// portraits, a metadata row, then the four graph cells.
const CARD_W = 540;
const CARD_H = 250;
const INSET = 4; // the live card's m-1 / ml-1 / mr-1 / mb-1

// Portrait row: 4 portraits, each 127x106 with a 4px margin (grid-cols-4 + m-1).
const PORTRAIT_W = 127;
const PORTRAIT_H = 106;
const PORTRAIT_MARGIN = 4;

// Graph row: timeline + histogram are wide (w-48 shrunk to 154), the two pies
// are 106; all 84 tall with a 4px gap (matches the live flex row exactly).
const GRAPH_GAP = 4;
const GRAPH_H = 84;
const WIDE_W = 154; // timeline + histogram
const PIE_W = 106;

const cell = {
	display: "flex" as const,
	alignItems: "center" as const,
	justifyContent: "center" as const,
	borderRadius: 4,
	backgroundColor: SLATE_700,
	height: GRAPH_H,
};

// A standalone, fixed-size (540x250) preview card that renders deterministically
// from a single model.SimulationResult. Hook-free, flexbox-only, no DOM
// measurement — shaped so `satori(<SatoriPreviewCard data={...} />, { fonts })`
// produces valid SVG.
export const SatoriPreviewCard = ({
	data,
	assetBase = DEFAULT_ASSET_BASE,
}: SatoriPreviewCardProps) => {
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "column",
				width: CARD_W,
				height: CARD_H,
				backgroundColor: SLATE_800,
				fontFamily: FONT_FAMILY,
			}}
		>
			<Portraits
				data={data}
				width={PORTRAIT_W}
				height={PORTRAIT_H}
				margin={PORTRAIT_MARGIN}
				assetBase={assetBase}
			/>

			<Metadata data={data} />

			<div
				style={{
					display: "flex",
					flexDirection: "row",
					gap: GRAPH_GAP,
					margin: `0 ${INSET}px ${INSET}px`,
				}}
			>
				<div style={{ ...cell, width: WIDE_W }}>
					<TimelineChart
						data={data.statistics?.damage_buckets}
						width={WIDE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: PIE_W }}>
					<CharacterPie
						dps={data.statistics?.character_dps}
						width={PIE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: PIE_W }}>
					<ElementPie
						dps={data.statistics?.element_dps}
						width={PIE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: WIDE_W }}>
					<HistogramChart
						data={data.statistics?.dps}
						width={WIDE_W}
						height={GRAPH_H}
					/>
				</div>
			</div>
		</div>
	);
};
