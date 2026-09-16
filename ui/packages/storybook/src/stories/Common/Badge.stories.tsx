import { Badge } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Badge> = {
	title: "Common/Badge",
	component: Badge,
	tags: ["autodocs"],
} as Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

const variants = [
	"default",
	"secondary",
	"destructive",
	"success",
	"warning",
	"outline",
] as const;

export const AllVariants: Story = {
	render: () => (
		<div className="flex flex-wrap gap-2">
			{variants.map((variant) => (
				<Badge key={variant} variant={variant}>
					{variant}
				</Badge>
			))}
		</div>
	),
};
