import { TargetAuraUptimeBarChart } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const meta: Meta<typeof TargetAuraUptimeBarChart> = {
	title: "Cards/TargetAuraUptimeBarChart",
	component: TargetAuraUptimeBarChart,
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<div className="w-[800px]">
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
		running: false,
	},
};

export const NoData: Story = {
	args: {
		data: null,
		running: false,
	},
};
