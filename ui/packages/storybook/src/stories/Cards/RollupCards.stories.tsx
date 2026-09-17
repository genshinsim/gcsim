import { RollupCards } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const meta: Meta<typeof RollupCards> = {
	title: "Cards/RollupCards",
	component: RollupCards,
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<div className="grid grid-cols-6 w-[1200px]">
				<Story />
			</div>
		),
	],
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
	args: {
		data: sampleResult,
	},
};

export const NoData: Story = {
	args: {
		data: null,
	},
};
