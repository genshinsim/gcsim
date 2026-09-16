import { Checkbox, Label } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Checkbox> = {
	title: "Common/Checkbox",
	component: Checkbox,
	tags: ["autodocs"],
} as Meta<typeof Checkbox>;

export default meta;
type Story = StoryObj<typeof meta>;

export const States: Story = {
	render: () => (
		<div className="flex flex-col gap-4">
			<Label>
				<Checkbox />
				Unchecked
			</Label>
			<Label>
				<Checkbox defaultChecked />
				Checked
			</Label>
			<Label>
				<Checkbox disabled />
				Disabled
			</Label>
			<Label>
				<Checkbox defaultChecked disabled />
				Disabled checked
			</Label>
		</div>
	),
};
