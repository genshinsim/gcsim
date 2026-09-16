import {
	Button,
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { ChevronsUpDown } from "lucide-react";

const meta: Meta<typeof Collapsible> = {
	title: "Common/Collapsible",
	component: Collapsible,
	tags: ["autodocs"],
} as Meta<typeof Collapsible>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Collapsible defaultOpen className="flex w-72 flex-col gap-2">
			<div className="flex items-center justify-between gap-4">
				<span className="text-sm font-medium">Advanced options</span>
				<CollapsibleTrigger asChild>
					<Button variant="ghost" size="icon-sm">
						<ChevronsUpDown />
					</Button>
				</CollapsibleTrigger>
			</div>
			<CollapsibleContent className="flex flex-col gap-2 text-sm">
				<div className="rounded-md border px-3 py-2">Fixed RNG seed</div>
				<div className="rounded-md border px-3 py-2">Debug logging</div>
			</CollapsibleContent>
		</Collapsible>
	),
};
