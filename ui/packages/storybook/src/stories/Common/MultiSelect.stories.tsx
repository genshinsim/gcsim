import { MultiSelect } from "@gcsim/components";
import { Badge, CommandItem } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { X } from "lucide-react";
import { useState } from "react";

type Fruit = { key: string; label: string };

const fruits: Fruit[] = [
	{ key: "apple", label: "Apple" },
	{ key: "banana", label: "Banana" },
	{ key: "cherry", label: "Cherry" },
	{ key: "durian", label: "Durian" },
	{ key: "elderberry", label: "Elderberry" },
];

const meta: Meta<typeof MultiSelect> = {
	title: "Common/MultiSelect",
	component: MultiSelect,
	tags: ["autodocs"],
} as Meta<typeof MultiSelect>;

export default meta;
type Story = StoryObj<typeof meta>;

function DemoMultiSelect() {
	const [value, setValue] = useState<Fruit[]>([fruits[0]]);

	return (
		<div className="w-96">
			<MultiSelect<Fruit>
				items={fruits}
				itemKey={(item) => item.key}
				value={value}
				onChange={setValue}
				placeholder="Select fruits..."
				itemPredicate={(item, query) =>
					item.label.toLowerCase().includes(query.toLowerCase())
				}
				itemRenderer={(item, { selected, onSelect }) => (
					<CommandItem key={item.key} value={item.key} onSelect={onSelect}>
						<span className="flex-1">{item.label}</span>
						{selected && <span className="text-xs">✓</span>}
					</CommandItem>
				)}
				tagRenderer={(item, { onRemove }) => (
					<Badge key={item.key} variant="secondary" className="gap-1">
						{item.label}
						<button
							type="button"
							aria-label={`Remove ${item.label}`}
							onClick={onRemove}
						>
							<X className="size-3" />
						</button>
					</Badge>
				)}
			/>
		</div>
	);
}

export const Default: Story = {
	render: () => <DemoMultiSelect />,
};
