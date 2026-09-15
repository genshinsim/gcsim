import { CharacterDPSBarChart } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, userEvent, within } from "storybook/test";
import { sampleResult } from "../samples";

const meta: Meta<typeof CharacterDPSBarChart> = {
	title: "Cards/CharacterDPSBarChart",
	component: CharacterDPSBarChart,
	tags: ["autodocs"],
	decorators: [
		(Story) => (
			<div className="w-[800px]">
				<Story />
			</div>
		),
	],
};

export default meta;
type Story = StoryObj<typeof meta>;

const names = sampleResult.character_details?.map((c) => c.name ?? "");

const selectGrouping = async (canvasElement: HTMLElement, value: string) => {
	const canvas = within(canvasElement);
	const select = await canvas.findByRole("combobox");
	await userEvent.selectOptions(select, value);
	await expect(select).toHaveValue(value);
};

// Default grouping renders the by-element breakdown.
export const ByElement: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
};

export const ByCharacter: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
	play: async ({ canvasElement }) => {
		await selectGrouping(canvasElement, "character");
	},
};

export const ByTarget: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
	play: async ({ canvasElement }) => {
		await selectGrouping(canvasElement, "target");
	},
};

export const NoData: Story = {
	args: {
		data: null,
		running: false,
	},
};
