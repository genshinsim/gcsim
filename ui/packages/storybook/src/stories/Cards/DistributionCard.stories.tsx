import { HistogramGraph } from "@gcsim/components";
import type { model } from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, waitFor } from "storybook/test";
import { sampleResult } from "../samples";

const meta: Meta<typeof HistogramGraph> = {
	title: "Cards/DistributionCard",
	component: HistogramGraph,
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof meta>;

const dps = sampleResult.statistics.dps as model.OverviewStats;

export const Primary: Story = {
	args: {
		width: 520,
		height: 250,
		data: dps,
	},
};

// Hovering the chart opens the popover-primitive tooltip.
export const WithTooltip: Story = {
	args: {
		width: 520,
		height: 250,
		data: dps,
	},
	play: async ({ canvasElement }) => {
		const svg = canvasElement.querySelector("svg");
		await expect(svg).toBeTruthy();

		// Move the pointer past the left axis margin, near the centre of the
		// distribution where bins have a non-zero count.
		const rect = svg!.getBoundingClientRect();
		svg!.dispatchEvent(
			new MouseEvent("mousemove", {
				bubbles: true,
				clientX: rect.left + rect.width * 0.5,
				clientY: rect.top + rect.height * 0.5,
			}),
		);

		await waitFor(() =>
			expect(
				document.querySelector('[data-slot="popover-content"]'),
			).toBeInTheDocument(),
		);
	},
};
