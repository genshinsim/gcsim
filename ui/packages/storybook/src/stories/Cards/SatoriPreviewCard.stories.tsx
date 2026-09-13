import { SatoriPreviewCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { cloneDeep, merge } from "lodash-es";
import { sampleResult } from "../samples";
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
// outline and two-set slice baked in — then passes them as `portraits`. Photon
// is a Node/wasm step, so this story consumes portraits baked by
// scripts/genSatoriPortraits.ts against the same sample. It's the visual check
// of the outline + slice: bennett (2pc + 2pc → half/half join), xiangling (lone
// 2pc → left-half slice), xingqiu/zhongli (single sets → full flower).
export const Composited: Story = {
	args: {
		data: satoriSample,
		portraits: [portrait0, portrait1, portrait2, portrait3],
	},
};
