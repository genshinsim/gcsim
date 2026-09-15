import { SourceDPSBarChart } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { expect, userEvent, within } from "storybook/test";
import { sampleResult } from "../samples";

const meta: Meta<typeof SourceDPSBarChart> = {
	title: "Cards/SourceDPSBarChart",
	component: SourceDPSBarChart,
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

// Default renders per-source DPS across all characters.
export const DPS: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
};

export const DamageInstances: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const [typeSelect] = await canvas.findAllByRole("combobox");
		await userEvent.selectOptions(typeSelect, "damage_instances");
		await expect(typeSelect).toHaveValue("damage_instances");
	},
};

export const FilteredByCharacter: Story = {
	args: {
		data: sampleResult,
		running: false,
		names: names,
	},
	play: async ({ canvasElement }) => {
		const canvas = within(canvasElement);
		const [, filterSelect] = await canvas.findAllByRole("combobox");
		await userEvent.selectOptions(filterSelect, names?.[0] ?? "");
		await expect(filterSelect).toHaveValue(names?.[0] ?? "");
	},
};

export const NoData: Story = {
	args: {
		data: null,
		running: false,
	},
};
