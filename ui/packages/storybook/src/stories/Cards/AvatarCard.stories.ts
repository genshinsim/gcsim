import { AvatarCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { sampleTeam } from "../samples";

// More on how to set up stories at: https://storybook.js.org/docs/writing-stories#default-export
const meta: Meta<typeof AvatarCard> = {
  title: "Cards/AvatarCard",
  component: AvatarCard,
  parameters: {
    // Optional parameter to center the component in the Canvas. More info: https://storybook.js.org/docs/configure/story-layout
    layout: "fullscreen",
  },
  // This component will have an automatically generated Autodocs entry: https://storybook.js.org/docs/writing-docs/autodocs
  tags: ["autodocs"],
  // More on argTypes: https://storybook.js.org/docs/api/argtypes
  argTypes: {},
  // Use `fn` to spy on the onClick arg, which will appear in the actions panel once invoked: https://storybook.js.org/docs/essentials/actions#action-args
  args: {
    onImageLoaded: fn(),
  },
};

export default meta;
type Story = StoryObj<typeof meta>;

// More on writing stories with args: https://storybook.js.org/docs/writing-stories/args
export const Primary: Story = {
  args: {
    chars: sampleTeam,
    invalid: ["yaemiko"],
  },
};

export const WithNull: Story = {
  args: {
    chars: [sampleTeam[0], sampleTeam[1], sampleTeam[2], null],
  },
};
