import { Button, Spinner } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Spinner> = {
	title: "Common/Spinner",
	component: Spinner,
	tags: ["autodocs"],
} as Meta<typeof Spinner>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Sizes: Story = {
	render: () => (
		<div className="flex items-center gap-4">
			<Spinner className="size-4" />
			<Spinner className="size-6" />
			<Spinner className="size-8" />
		</div>
	),
};

export const InButton: Story = {
	render: () => (
		<Button disabled>
			<Spinner />
			Running…
		</Button>
	),
};
