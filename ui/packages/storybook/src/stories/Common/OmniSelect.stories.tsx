import { CharacterSelect, OmniSelect } from "@gcsim/components";
import { Button, CommandItem } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";

const fruits = [
	{ key: "apple", label: "Apple" },
	{ key: "banana", label: "Banana" },
	{ key: "cherry", label: "Cherry" },
	{ key: "durian", label: "Durian" },
	{ key: "elderberry", label: "Elderberry" },
];

const meta: Meta<typeof OmniSelect> = {
	title: "Common/OmniSelect",
	component: OmniSelect,
	tags: ["autodocs"],
} as Meta<typeof OmniSelect>;

export default meta;
type Story = StoryObj<typeof meta>;

function DemoOmniSelect() {
	const [isOpen, setIsOpen] = useState(false);
	const [value, setValue] = useState<(typeof fruits)[number]>();

	return (
		<div className="flex flex-col items-start gap-2">
			<Button onClick={() => setIsOpen(true)}>
				{value ? `Selected: ${value.label}` : "Pick a fruit"}
			</Button>
			<OmniSelect
				isOpen={isOpen}
				onClose={() => setIsOpen(false)}
				items={fruits}
				itemKey={(item) => item.key}
				value={value}
				itemPredicate={(item, query) =>
					item.label.toLowerCase().includes(query.toLowerCase())
				}
				onSelect={setValue}
				title="Select a fruit"
				placeholder="Search fruits..."
				itemRenderer={(item, { selected, onSelect }) => (
					<CommandItem key={item.key} value={item.key} onSelect={onSelect}>
						{item.label}
						{selected && <span className="ml-auto text-xs">current</span>}
					</CommandItem>
				)}
			/>
		</div>
	);
}

export const Default: Story = {
	render: () => <DemoOmniSelect />,
};

function DemoCharacterSelect() {
	const [isOpen, setIsOpen] = useState(false);
	const [value, setValue] = useState<string>();

	return (
		<div className="flex flex-col items-start gap-2">
			<Button onClick={() => setIsOpen(true)}>
				{value ? `Selected: ${value}` : "Pick a character"}
			</Button>
			<CharacterSelect
				isOpen={isOpen}
				onClose={() => setIsOpen(false)}
				value={value}
				onSelect={(character) => setValue(character)}
			/>
		</div>
	);
}

export const TypedCharacterSelect: Story = {
	render: () => <DemoCharacterSelect />,
};
