import { ResultsView } from "@gcsim/components";
import type { model } from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleResult } from "../samples";

const result = sampleResult as model.SimulationResult;

const meta: Meta<typeof ResultsView> = {
	title: "Views/ResultsView",
	component: ResultsView,
	parameters: { layout: "fullscreen" },
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
	args: {
		model: result,
		names: result.character_details?.map((c) => c.name ?? ""),
	},
};

export const Loading: Story = {
	args: {
		model: null,
	},
};
