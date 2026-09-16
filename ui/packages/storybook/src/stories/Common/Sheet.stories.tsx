import {
	Button,
	Sheet,
	SheetClose,
	SheetContent,
	SheetDescription,
	SheetFooter,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Sheet> = {
	title: "Common/Sheet",
	component: Sheet,
	tags: ["autodocs"],
} as Meta<typeof Sheet>;

export default meta;
type Story = StoryObj<typeof meta>;

function DemoSheet({ side }: { side?: "top" | "right" | "bottom" | "left" }) {
	return (
		<Sheet defaultOpen>
			<SheetTrigger asChild>
				<Button variant="outline">Open {side ?? "right"} sheet</Button>
			</SheetTrigger>
			<SheetContent side={side}>
				<SheetHeader>
					<SheetTitle>Settings</SheetTitle>
					<SheetDescription>Adjust simulation options.</SheetDescription>
				</SheetHeader>
				<SheetFooter>
					<SheetClose asChild>
						<Button>Done</Button>
					</SheetClose>
				</SheetFooter>
			</SheetContent>
		</Sheet>
	);
}

export const Right: Story = {
	render: () => <DemoSheet side="right" />,
};

export const Left: Story = {
	render: () => <DemoSheet side="left" />,
};

export const Bottom: Story = {
	render: () => <DemoSheet side="bottom" />,
};
