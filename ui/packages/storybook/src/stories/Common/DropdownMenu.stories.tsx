import {
	Button,
	DropdownMenu,
	DropdownMenuCheckboxItem,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuShortcut,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof DropdownMenu> = {
	title: "Common/DropdownMenu",
	component: DropdownMenu,
	tags: ["autodocs"],
} as Meta<typeof DropdownMenu>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<DropdownMenu defaultOpen>
			<DropdownMenuTrigger asChild>
				<Button variant="outline">Actions</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-52">
				<DropdownMenuLabel>Config</DropdownMenuLabel>
				<DropdownMenuItem>
					Duplicate
					<DropdownMenuShortcut>⌘D</DropdownMenuShortcut>
				</DropdownMenuItem>
				<DropdownMenuCheckboxItem checked>
					Show tooltips
				</DropdownMenuCheckboxItem>
				<DropdownMenuSub>
					<DropdownMenuSubTrigger>Export as</DropdownMenuSubTrigger>
					<DropdownMenuSubContent>
						<DropdownMenuItem>JSON</DropdownMenuItem>
						<DropdownMenuItem>Image</DropdownMenuItem>
					</DropdownMenuSubContent>
				</DropdownMenuSub>
				<DropdownMenuSeparator />
				<DropdownMenuItem variant="destructive">Delete</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	),
};
