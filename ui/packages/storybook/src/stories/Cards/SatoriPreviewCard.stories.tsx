import { SatoriPreviewCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { cloneDeep, merge } from "lodash-es";
import { sampleResult } from "../samples";
import { compositePortraitsInBrowser } from "../samples/compositePortraitsBrowser";
import portrait0 from "../samples/satoriPortraits/portrait0.png";
import portrait1 from "../samples/satoriPortraits/portrait1.png";
import portrait2 from "../samples/satoriPortraits/portrait2.png";
import portrait3 from "../samples/satoriPortraits/portrait3.png";
import { satoriSample } from "../samples/satoriSample";

// The Satori preview card is a fixed-size (540x250), hook-free tree. It renders
// as ordinary React in the browser, so these stories double as a visual check of
// the same tree Satori rasterizes server-side.
const meta: Meta<typeof SatoriPreviewCard> = {
	title: "Cards/SatoriPreviewCard",
	component: SatoriPreviewCard,
	tags: ["autodocs"],
	argTypes: {},
	args: {},
	parameters: {
		viewport: {
			defaultViewport: "discord",
		},
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
	args: {
		data: sampleResult,
	},
};

const incomplete = cloneDeep(sampleResult);
export const WithIncomplete: Story = {
	args: {
		data: merge(incomplete, {
			incomplete_characters: ["xingqiu"],
		}),
	},
};

const missingImages = cloneDeep(sampleResult);
missingImages.character_details[0].sets = { fake: 4 };
missingImages.character_details[0].name = "fake";
missingImages.character_details[0].weapon.name = "fake";

export const MissingImages: Story = {
	args: {
		data: missingImages,
	},
};

// The production OG render (edge Worker) pre-composites each portrait with
// Photon — element bg + avatar + weapon + artifact set(s), with the white icon
// outline and two-set slice baked in — then passes them as `portraits`. It's the
// visual check of the outline + slice: bennett (2pc + 2pc → half/half join),
// xiangling (lone 2pc → left-half slice), xingqiu/zhongli (single sets → full
// flower).
//
// Photon is a wasm step, so two variants:
//   - Composited: portraits baked by scripts/genSatoriPortraits.ts. Works in any
//     build (incl. static/Chromatic); regenerate when the compositor changes.
//   - CompositedLive: runs the compositor in the browser on every render, so
//     edits to portraitCompositor.ts show up on hot-reload. Dev only (needs the
//     `/api` proxy), so it's disabled in Chromatic.
export const Composited: Story = {
	args: {
		data: satoriSample,
		portraits: [portrait0, portrait1, portrait2, portrait3],
	},
};

export const CompositedLive: Story = {
	parameters: { chromatic: { disable: true } },
	loaders: [
		async () => ({
			portraits: await compositePortraitsInBrowser(satoriSample),
		}),
	],
	render: (_args, { loaded }) => (
		<SatoriPreviewCard
			data={satoriSample}
			portraits={(loaded as { portraits: string[] }).portraits}
		/>
	),
};
