import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectSeparator,
	SelectTrigger,
	SelectValue,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Select> = {
	title: "Common/Select",
	component: Select,
	tags: ["autodocs"],
} as Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

function ElementSelect({
	size,
	disabled,
}: {
	size?: "sm" | "default";
	disabled?: boolean;
}) {
	return (
		<Select disabled={disabled}>
			<SelectTrigger size={size} className="w-48">
				<SelectValue placeholder="Pick an element" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectLabel>Elements</SelectLabel>
					<SelectItem value="pyro">Pyro</SelectItem>
					<SelectItem value="hydro">Hydro</SelectItem>
					<SelectItem value="electro">Electro</SelectItem>
					<SelectItem value="cryo">Cryo</SelectItem>
				</SelectGroup>
				<SelectSeparator />
				<SelectItem value="none" disabled>
					Physical (disabled)
				</SelectItem>
			</SelectContent>
		</Select>
	);
}

export const Default: Story = {
	render: () => <ElementSelect />,
};

export const Small: Story = {
	render: () => <ElementSelect size="sm" />,
};

export const Disabled: Story = {
	render: () => <ElementSelect disabled />,
};
