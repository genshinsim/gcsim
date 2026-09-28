import { DefaultSampleOptions, SampleLog } from "@gcsim/components";
import { sampleFixture } from "@gcsim/components/src/SampleLog/testdata";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";

const meta: Meta<typeof SampleLog> = {
	title: "Sample/SampleLog",
	component: SampleLog,
	decorators: [
		(Story) => (
			<div className="flex flex-col gap-[15px] bg-g-surface p-4">
				<Story />
			</div>
		),
	],
	args: {
		sample: sampleFixture,
		settings: DefaultSampleOptions,
		onSettingsChange: fn(),
		onDownload: fn(),
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const DamageOnly: Story = {
	args: { settings: ["damage"] },
};

export const WithoutDownload: Story = {
	args: { onDownload: undefined },
};
