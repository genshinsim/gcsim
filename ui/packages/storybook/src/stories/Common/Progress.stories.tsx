import { Progress } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Progress> = {
	title: "Common/Progress",
	component: Progress,
	tags: ["autodocs"],
} as Meta<typeof Progress>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Values: Story = {
	render: () => (
		<div className="flex w-80 flex-col gap-4">
			<Progress value={0} />
			<Progress value={40} />
			<Progress value={100} />
		</div>
	),
};
