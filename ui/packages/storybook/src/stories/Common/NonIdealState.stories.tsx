import { Button, NonIdealState } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inbox } from "lucide-react";

const meta: Meta<typeof NonIdealState> = {
	title: "Common/NonIdealState",
	component: NonIdealState,
	tags: ["autodocs"],
} as Meta<typeof NonIdealState>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		icon: <Inbox />,
		title: "No results found",
		description: "Try adjusting your filters to find what you're looking for.",
		action: <Button variant="outline">Reset filters</Button>,
	},
};

export const Loading: Story = {
	args: {
		loading: true,
		title: "Loading",
	},
};

export const Horizontal: Story = {
	args: {
		layout: "horizontal",
		icon: <Inbox />,
		title: "No results found",
		description: "Try adjusting your filters.",
	},
};
