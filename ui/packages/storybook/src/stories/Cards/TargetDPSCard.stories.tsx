import { TargetDPSCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const meta: Meta<typeof TargetDPSCard> = {
	title: "Cards/TargetDPSCard",
	component: TargetDPSCard,
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<div className="w-[400px]">
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
