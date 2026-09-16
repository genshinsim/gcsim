import { Checkbox, Input, Label } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Label> = {
	title: "Common/Label",
	component: Label,
	tags: ["autodocs"],
} as Meta<typeof Label>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithInput: Story = {
	render: () => (
		<div className="flex flex-col gap-2">
			<Label htmlFor="iterations">Iterations</Label>
			<Input id="iterations" placeholder="1000" />
		</div>
	),
};

export const WithCheckbox: Story = {
	render: () => (
		<Label>
			<Checkbox defaultChecked />
			Enable optimizer
		</Label>
	),
};
