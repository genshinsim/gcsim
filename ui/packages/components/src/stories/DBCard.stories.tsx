import { DBCard } from "@/Cards/DBCard/DBCard";
import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { dbEntries } from "./samples";
import { Button } from "@/common/ui/button";

// More on how to set up stories at: https://storybook.js.org/docs/writing-stories#default-export
const meta: Meta<typeof DBCard> = {
  title: "DBCard",
  component: DBCard,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: "fullscreen",
  },
  // This component will have an automatically generated Autodocs entry: https://storybook.js.org/docs/writing-docs/autodocs
  tags: ["autodocs"],
  // More on argTypes: https://storybook.js.org/docs/api/argtypes
  argTypes: {},
  // Use `fn` to spy on the onClick arg, which will appear in the actions panel once invoked: https://storybook.js.org/docs/essentials/actions#action-args
  args: {},
};

export default meta;
type Story = StoryObj<typeof meta>;

// More on writing stories with args: https://storybook.js.org/docs/writing-stories/args
export const Primary: Story = {
  args: {
    entry: dbEntries.data[0],
    footer: <Button>Show Results Page</Button>
  },
};
