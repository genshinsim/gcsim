import {
	Button,
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Dialog> = {
	title: "Common/Dialog",
	component: Dialog,
	tags: ["autodocs"],
} as Meta<typeof Dialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Dialog defaultOpen>
			<DialogTrigger asChild>
				<Button variant="outline">Share config</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Share config</DialogTitle>
					<DialogDescription>
						Anyone with the link can view this simulation.
					</DialogDescription>
				</DialogHeader>
				<DialogFooter>
					<DialogClose asChild>
						<Button variant="outline">Cancel</Button>
					</DialogClose>
					<Button>Copy link</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	),
};
