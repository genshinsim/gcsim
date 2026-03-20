import type { Meta, StoryObj } from "@storybook/react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@gcsim/primitives";

const meta = {
  title: "Primitives/Tabs",
  component: Tabs,
  tags: ["autodocs"],
} satisfies Meta<typeof Tabs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Tabs defaultValue="results" className="w-[400px]">
      <TabsList>
        <TabsTrigger value="results">Results</TabsTrigger>
        <TabsTrigger value="config">Config</TabsTrigger>
        <TabsTrigger value="sample">Sample</TabsTrigger>
      </TabsList>
      <TabsContent value="results">
        <p className="p-4 text-sm text-muted-foreground">Results content goes here.</p>
      </TabsContent>
      <TabsContent value="config">
        <p className="p-4 text-sm text-muted-foreground">Config content goes here.</p>
      </TabsContent>
      <TabsContent value="sample">
        <p className="p-4 text-sm text-muted-foreground">Sample content goes here.</p>
      </TabsContent>
    </Tabs>
  ),
};
