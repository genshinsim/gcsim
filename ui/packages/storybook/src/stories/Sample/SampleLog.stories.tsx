import {
	loadSampleFilter,
	SampleLog,
	saveSampleFilter,
} from "@gcsim/components";
import { AllSampleOptions } from "@gcsim/components/src/SampleLog/SampleOptions";
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
		onDownload: fn(),
		onGenerate: fn(),
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

const seedFilter = (filter: string[]) => () => {
	const saved = loadSampleFilter(localStorage);
	saveSampleFilter(localStorage, filter);
	return () => saveSampleFilter(localStorage, saved);
};

export const Default: Story = {
	beforeEach: seedFilter(AllSampleOptions),
};

export const DamageOnly: Story = {
	beforeEach: seedFilter(["damage"]),
};

export const WithoutDownload: Story = {
	args: { onDownload: undefined },
	beforeEach: seedFilter(AllSampleOptions),
};

export const NoSample: Story = {
	args: { sample: null },
};
