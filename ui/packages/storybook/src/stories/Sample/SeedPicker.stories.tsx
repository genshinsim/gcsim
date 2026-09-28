import { SeedPicker } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";

const meta: Meta<typeof SeedPicker> = {
	title: "Sample/SeedPicker",
	component: SeedPicker,
	decorators: [
		(Story) => (
			<div className="w-[480px] max-w-full bg-g-surface p-4">
				<Story />
			</div>
		),
	],
	args: {
		seeds: {
			sample: "14923829596786871748",
			min: "1203948576012398475",
			max: "9834750192837465012",
			p25: "5519283746501928374",
			p50: "7710293847561029384",
			p75: "3301928475610293847",
		},
		value: null,
		onPick: fn(),
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const MaxSelected: Story = {
	args: { value: "9834750192837465012" },
};

export const CustomSeed: Story = {
	args: { value: "18446744073709551615" },
};

export const Running: Story = {
	args: { value: "1203948576012398475", running: true },
};
