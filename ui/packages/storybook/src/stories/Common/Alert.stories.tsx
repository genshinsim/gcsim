import { Alert, AlertDescription, AlertTitle } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { CircleAlert, CircleCheck, Info, TriangleAlert } from "lucide-react";

const meta: Meta<typeof Alert> = {
	title: "Common/Alert",
	component: Alert,
	tags: ["autodocs"],
} as Meta<typeof Alert>;

export default meta;
type Story = StoryObj<typeof meta>;

const variants = [
	{ variant: "default", icon: <Info /> },
	{ variant: "destructive", icon: <CircleAlert /> },
	{ variant: "warning", icon: <TriangleAlert /> },
	{ variant: "success", icon: <CircleCheck /> },
] as const;

export const AllVariants: Story = {
	render: () => (
		<div className="flex max-w-lg flex-col gap-4">
			{variants.map(({ variant, icon }) => (
				<Alert key={variant} variant={variant}>
					{icon}
					<AlertTitle>{variant}</AlertTitle>
					<AlertDescription>
						This is a {variant} alert with a title and a longer body message.
					</AlertDescription>
				</Alert>
			))}
		</div>
	),
};
