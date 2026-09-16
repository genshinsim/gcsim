import {
	Button,
	ButtonGroup,
	ButtonGroupSeparator,
	ButtonGroupText,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof ButtonGroup> = {
	title: "Common/ButtonGroup",
	component: ButtonGroup,
	tags: ["autodocs"],
} as Meta<typeof ButtonGroup>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Horizontal: Story = {
	render: () => (
		<ButtonGroup>
			<Button variant="outline">Run</Button>
			<Button variant="outline">Share</Button>
			<ButtonGroupSeparator />
			<Button variant="outline">Reset</Button>
		</ButtonGroup>
	),
};

export const WithText: Story = {
	render: () => (
		<ButtonGroup>
			<ButtonGroupText>Iterations</ButtonGroupText>
			<Button variant="outline">1000</Button>
			<Button variant="outline">5000</Button>
		</ButtonGroup>
	),
};

export const Vertical: Story = {
	render: () => (
		<ButtonGroup orientation="vertical">
			<Button variant="outline">Top</Button>
			<Button variant="outline">Middle</Button>
			<Button variant="outline">Bottom</Button>
		</ButtonGroup>
	),
};
