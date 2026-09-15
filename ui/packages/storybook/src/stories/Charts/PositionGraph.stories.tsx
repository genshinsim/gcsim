import { PositionGraph } from "@gcsim/components";
import type { model } from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, waitFor } from "storybook/test";
import { sampleResult } from "../samples";

const meta: Meta<typeof PositionGraph> = {
	title: "Charts/PositionGraph",
	component: PositionGraph,
	tags: ["autodocs"],
	args: {
		width: 320,
		height: 320,
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
	args: {
		enemies: sampleResult.target_details as model.Enemy[],
		player: sampleResult.player_position as model.Coord,
	},
};

export const NoTargets: Story = {
	args: {
		enemies: undefined,
		player: undefined,
	},
};

// Hovering a target opens the popover-primitive tooltip.
export const WithTooltip: Story = {
	args: {
		enemies: sampleResult.target_details as model.Enemy[],
		player: sampleResult.player_position as model.Coord,
	},
	play: async ({ canvasElement }) => {
		const circle = canvasElement.querySelector("circle");
		await expect(circle).toBeTruthy();

		circle!.dispatchEvent(
			new MouseEvent("mousemove", { bubbles: true, clientX: 0, clientY: 0 }),
		);

		await waitFor(() =>
			expect(
				document.querySelector('[data-slot="popover-content"]'),
			).toBeInTheDocument(),
		);
	},
};
