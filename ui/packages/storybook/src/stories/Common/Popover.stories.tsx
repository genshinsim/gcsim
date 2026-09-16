import { Popover, PopoverContent, PopoverTrigger } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Popover> = {
	title: "Common/Popover",
	component: Popover,
	tags: ["autodocs"],
} as Meta<typeof Popover>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {
	render: () => (
		<Popover defaultOpen>
			<PopoverTrigger>Open popover</PopoverTrigger>
			<PopoverContent>Popover content</PopoverContent>
		</Popover>
	),
};
