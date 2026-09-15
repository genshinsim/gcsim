import { CharacterActionsBarChart } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const names: string[] = sampleResult.character_details?.map(
	(c: { name: string }) => c.name,
);

const meta: Meta<typeof CharacterActionsBarChart> = {
	title: "Cards/CharacterActionsBarChart",
	component: CharacterActionsBarChart,
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
		names,
	},
};

export const NoData: Story = {
	args: {
		data: null,
		running: false,
		names,
	},
};
