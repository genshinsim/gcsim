import type { Meta, StoryObj } from "@storybook/react";
import { ScrollArea } from "@gcsim/primitives";

const meta = {
  title: "Primitives/ScrollArea",
  component: ScrollArea,
  tags: ["autodocs"],
} satisfies Meta<typeof ScrollArea>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <ScrollArea className="h-[200px] w-[350px] rounded-md border p-4">
      {Array.from({ length: 20 }, (_, i) => (
        <div key={i} className="py-1 text-sm">
          Item {i + 1}
        </div>
      ))}
    </ScrollArea>
  ),
};
