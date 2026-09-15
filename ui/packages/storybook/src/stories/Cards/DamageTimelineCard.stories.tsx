import { DamageTimelineCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const names: string[] = sampleResult.character_details?.map(
	(c: { name: string }) => c.name,
);

const meta: Meta<typeof DamageTimelineCard> = {
	title: "Cards/DamageTimelineCard",
	component: DamageTimelineCard,
	tags: ["autodocs"],
	args: {},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
	args: {
		data: sampleResult,
		running: false,
		names,
	},
};

export const Running: Story = {
	args: {
		data: null,
		running: true,
		names,
	},
};
