import { FloatStatTooltipContent } from "@gcsim/components";
import { Tooltip, TooltipContent, TooltipTrigger } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Tooltip> = {
	title: "Common/Tooltip",
	component: Tooltip,
	tags: ["autodocs"],
} as Meta<typeof Tooltip>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {
	render: () => (
		<Tooltip defaultOpen>
			<TooltipTrigger>Hover for tooltip</TooltipTrigger>
			<TooltipContent>Tooltip content</TooltipContent>
		</Tooltip>
	),
};

export const FloatStat: Story = {
	render: () => (
		<Tooltip defaultOpen>
			<TooltipTrigger>Hover for stats</TooltipTrigger>
			<TooltipContent>
				<FloatStatTooltipContent
					title="Furina DPS"
					data={{ mean: 12345.6, min: 9876.5, max: 15432.1, sd: 1234.5 }}
					percent={0.42}
				/>
			</TooltipContent>
		</Tooltip>
	),
};
