import { Label, Switch } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Switch> = {
	title: "Common/Switch",
	component: Switch,
	tags: ["autodocs"],
} as Meta<typeof Switch>;

export default meta;
type Story = StoryObj<typeof meta>;

export const States: Story = {
	render: () => (
		<div className="flex flex-col gap-4">
			<Label>
				<Switch />
				Off
			</Label>
			<Label>
				<Switch defaultChecked />
				On
			</Label>
			<Label>
				<Switch disabled />
				Disabled
			</Label>
			<Label>
				<Switch defaultChecked disabled />
				Disabled on
			</Label>
		</div>
	),
};
