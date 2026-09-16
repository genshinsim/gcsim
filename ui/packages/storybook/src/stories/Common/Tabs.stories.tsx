import { Tabs, TabsContent, TabsList, TabsTrigger } from "@gcsim/primitives";
import type { Meta, StoryObj } from "@storybook/react-vite";

const meta: Meta<typeof Tabs> = {
	title: "Common/Tabs",
	component: Tabs,
	tags: ["autodocs"],
} as Meta<typeof Tabs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Tabs defaultValue="overview" className="w-80">
			<TabsList>
				<TabsTrigger value="overview">Overview</TabsTrigger>
				<TabsTrigger value="damage">Damage</TabsTrigger>
				<TabsTrigger value="energy" disabled>
					Energy
				</TabsTrigger>
			</TabsList>
			<TabsContent value="overview">Team overview and summary.</TabsContent>
			<TabsContent value="damage">Per-character damage breakdown.</TabsContent>
		</Tabs>
	),
};
