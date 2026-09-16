import {
	Navbar,
	NavbarDivider,
	NavbarGroup,
	NavbarHeading,
} from "@gcsim/components";
import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { Calculator, Database, FileText, Globe, Menu } from "lucide-react";

const meta: Meta<typeof Navbar> = {
	title: "Common/Navbar",
	component: Navbar,
	tags: ["autodocs"],
	parameters: { layout: "fullscreen" },
} as Meta<typeof Navbar>;

export default meta;
type Story = StoryObj<typeof meta>;

function DemoNavbar() {
	return (
		<Navbar>
			<NavbarGroup>
				<NavbarHeading className="font-mono">gcsim</NavbarHeading>
				<NavbarDivider />
				<Button variant="ghost" size="sm">
					<Calculator />
					Simulator
				</Button>
				<Button variant="ghost" size="sm">
					<Database />
					Teams DB
				</Button>
				<Button variant="ghost" size="sm">
					<FileText />
					Documentation
				</Button>
			</NavbarGroup>
			<NavbarGroup align="end">
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<Button variant="ghost" size="sm">
							<Globe />
							English
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuItem>English</DropdownMenuItem>
						<DropdownMenuItem>中文</DropdownMenuItem>
						<DropdownMenuItem>日本語</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
				<Button variant="ghost" size="icon-sm" aria-label="Menu">
					<Menu />
				</Button>
			</NavbarGroup>
		</Navbar>
	);
}

export const Default: Story = {
	render: () => <DemoNavbar />,
};
