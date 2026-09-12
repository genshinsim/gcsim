import type { model } from "@gcsim/types";
import { HistogramChart } from "./charts/HistogramChart";
import { CharacterPie, ElementPie } from "./charts/PieChart";
import { TimelineChart } from "./charts/TimelineChart";
import { SLATE_700, SLATE_800 } from "./colors";
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

// Fixed card geometry (540x250), all in pixels so each chart knows its exact
// render size — no DOM measurement, no responsive wrappers.
const CARD_W = 540;
const CARD_H = 250;
const PAD = 4;
const ROW_GAP = 4;

const INNER_W = CARD_W - 2 * PAD;
const PORTRAIT_GAP = 4;
const PORTRAIT_W = (INNER_W - 3 * PORTRAIT_GAP) / 4;
const PORTRAIT_H = 104;

const GRAPH_GAP = 4;
const GRAPH_H = 96;
const WIDE_W = 148; // timeline + histogram
const PIE_W = 106;

// A standalone, fixed-size (540x250) preview card that renders deterministically
// from a single model.SimulationResult. Hook-free, flexbox-only, no DOM
// measurement — shaped so `satori(<SatoriPreviewCard data={...} />, { fonts })`
// produces valid SVG.
export const SatoriPreviewCard = ({
	data,
	assetBase = DEFAULT_ASSET_BASE,
}: SatoriPreviewCardProps) => {
	const cell = {
		display: "flex" as const,
		alignItems: "center" as const,
		justifyContent: "center" as const,
		borderRadius: 4,
		backgroundColor: SLATE_700,
	};

	return (
		<div
			style={{
				display: "flex",
				flexDirection: "column",
				width: CARD_W,
				height: CARD_H,
				padding: PAD,
				gap: ROW_GAP,
				backgroundColor: SLATE_800,
				fontFamily: "Inter",
			}}
		>
			<Portraits
				data={data}
				width={PORTRAIT_W}
				height={PORTRAIT_H}
				gap={PORTRAIT_GAP}
				assetBase={assetBase}
			/>

			<Metadata data={data} />

			<div
				style={{
					display: "flex",
					flexDirection: "row",
					gap: GRAPH_GAP,
					justifyContent: "center",
					margin: "0 4px",
				}}
			>
				<div style={{ ...cell, width: WIDE_W, height: GRAPH_H }}>
					<TimelineChart
						data={data.statistics?.damage_buckets}
						width={WIDE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: PIE_W, height: GRAPH_H }}>
					<CharacterPie
						dps={data.statistics?.character_dps}
						width={PIE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: PIE_W, height: GRAPH_H }}>
					<ElementPie
						dps={data.statistics?.element_dps}
						width={PIE_W}
						height={GRAPH_H}
					/>
				</div>
				<div style={{ ...cell, width: WIDE_W, height: GRAPH_H }}>
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
