import { Label, NumberInput, type NumberInputProps } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";

const meta: Meta<typeof NumberInput> = {
	title: "Common/NumberInput",
	component: NumberInput,
	tags: ["autodocs"],
} as Meta<typeof NumberInput>;

export default meta;
type Story = StoryObj<typeof meta>;

function Controlled({
	value: initial,
	...props
}: Omit<NumberInputProps, "onValueChange">) {
	const [value, setValue] = useState(initial);
	return (
		<div className="flex w-48 flex-col gap-2">
			<NumberInput value={value} onValueChange={setValue} {...props} />
			<span className="text-sm text-g-ink-mute">value: {value}</span>
		</div>
	);
}

export const Default: Story = {
	render: () => <Controlled value={5} />,
};

export const MinMax: Story = {
	render: () => (
		<Label className="flex flex-col items-start gap-2">
			Level (1–90)
			<Controlled value={90} min={1} max={90} />
		</Label>
	),
};

export const Step: Story = {
	render: () => <Controlled value={1000} step={500} min={0} />,
};
