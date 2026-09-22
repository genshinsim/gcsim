import { SampleEventDetails } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleEvent } from "../samples";

const meta: Meta<typeof SampleEventDetails> = {
	title: "Sample/SampleEventDetails",
	component: SampleEventDetails,
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<div className="w-[768px] max-w-full bg-g-surface p-4">
				<Story />
			</div>
		),
	],
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Nested: Story = {
	args: {
		title: sampleEvent.msg,
		data: sampleEvent,
	},
};

export const RawFallback: Story = {
	args: {
		title: sampleEvent.msg,
		raw: JSON.stringify(sampleEvent, null, 2),
	},
};

export const Empty: Story = {
	args: {
		title: "No event data",
	},
};
