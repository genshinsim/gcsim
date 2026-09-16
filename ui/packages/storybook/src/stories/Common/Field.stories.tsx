import {
	Field,
	FieldContent,
	FieldDescription,
	FieldError,
	FieldGroup,
	FieldLabel,
	Input,
	Switch,
} from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Field> = {
	title: "Common/Field",
	component: Field,
	tags: ["autodocs"],
} as Meta<typeof Field>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Vertical: Story = {
	render: () => (
		<FieldGroup className="max-w-sm">
			<Field>
				<FieldLabel htmlFor="name">Config name</FieldLabel>
				<Input id="name" placeholder="my-team" />
				<FieldDescription>Shown in the results list.</FieldDescription>
			</Field>
			<Field>
				<FieldLabel htmlFor="iters">Iterations</FieldLabel>
				<Input id="iters" defaultValue="abc" aria-invalid />
				<FieldError errors={[{ message: "Must be a number." }]} />
			</Field>
		</FieldGroup>
	),
};

export const Horizontal: Story = {
	render: () => (
		<Field orientation="horizontal" className="max-w-sm">
			<FieldContent>
				<FieldLabel htmlFor="optimize">Optimize stats</FieldLabel>
				<FieldDescription>Run the substat optimizer.</FieldDescription>
			</FieldContent>
			<Switch id="optimize" />
		</Field>
	),
};
