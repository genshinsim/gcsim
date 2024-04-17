import { Editor } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react";
import { sampleConfig } from "./sampleConfig";

// More on how to set up stories at: https://storybook.js.org/docs/writing-stories#default-export
const meta: Meta<typeof Editor> = {
  title: "Editor/Editor",
  component: Editor,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: "padded",
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
    cfg: sampleConfig,
  },
};

export const PrimaryMobile: Story = {
  args: {
    cfg: sampleConfig,
  },
  parameters: {
    viewport: {
      defaultViewport: "mobile1",
    },
  },
};

export const PrimaryTablet: Story = {
  args: {
    cfg: sampleConfig,
  },
  parameters: {
    viewport: {
      defaultViewport: "tablet",
    },
  },
};
