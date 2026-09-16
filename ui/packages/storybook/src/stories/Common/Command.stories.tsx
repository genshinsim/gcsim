import {
	Command,
	CommandDialog,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
	CommandSeparator,
	CommandShortcut,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { FileText, Settings, Sparkles, Zap } from "lucide-react";

const meta: Meta<typeof Command> = {
	title: "Common/Command",
	component: Command,
	tags: ["autodocs"],
} as Meta<typeof Command>;

export default meta;
type Story = StoryObj<typeof meta>;

function Items() {
	return (
		<CommandList>
			<CommandEmpty>No results found.</CommandEmpty>
			<CommandGroup heading="Characters">
				<CommandItem>
					<Sparkles />
					<span>Raiden Shogun</span>
				</CommandItem>
				<CommandItem>
					<Sparkles />
					<span>Hu Tao</span>
				</CommandItem>
				<CommandItem>
					<Sparkles />
					<span>Nahida</span>
				</CommandItem>
			</CommandGroup>
			<CommandSeparator />
			<CommandGroup heading="Actions">
				<CommandItem>
					<Zap />
					<span>Run simulation</span>
					<CommandShortcut>⌘R</CommandShortcut>
				</CommandItem>
				<CommandItem>
					<FileText />
					<span>Open config</span>
					<CommandShortcut>⌘O</CommandShortcut>
				</CommandItem>
				<CommandItem disabled>
					<Settings />
					<span>Settings (disabled)</span>
				</CommandItem>
			</CommandGroup>
		</CommandList>
	);
}

export const Default: Story = {
	render: () => (
		<Command className="max-w-md rounded-lg border shadow-md">
			<CommandInput placeholder="Search characters and actions..." />
			<Items />
		</Command>
	),
};

export const Palette: Story = {
	render: () => (
		<CommandDialog defaultOpen>
			<CommandInput placeholder="Type a command or search..." />
			<Items />
		</CommandDialog>
	),
};
