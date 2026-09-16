import { Button } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const variants = [
	"default",
	"destructive",
	"outline",
	"secondary",
	"ghost",
	"link",
] as const;

const sizes = [
	"default",
	"xs",
	"sm",
	"lg",
	"icon",
	"icon-xs",
	"icon-sm",
	"icon-lg",
] as const;

const meta: Meta<typeof Button> = {
	title: "Common/Button",
	component: Button,
	tags: ["autodocs"],
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		children: "Button",
	},
};

export const AllVariants: Story = {
	render: () => (
		<div style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
			{variants.map((variant) => (
				<Button key={variant} variant={variant}>
					{variant}
				</Button>
			))}
		</div>
	),
};

export const AllSizes: Story = {
	render: () => (
		<div
			style={{
				display: "flex",
				gap: "0.5rem",
				flexWrap: "wrap",
				alignItems: "center",
			}}
		>
			{sizes.map((size) => (
				<Button key={size} size={size}>
					{size.startsWith("icon") ? "•" : size}
				</Button>
			))}
		</div>
	),
};
