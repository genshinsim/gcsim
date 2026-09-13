import { SatoriPreviewCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { cloneDeep, merge } from "lodash-es";
import { sampleResult } from "../samples";

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
