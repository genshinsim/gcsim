import { PreviewCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react";
import { cloneDeep, merge } from "lodash";
import { sampleResult } from "../samples";

// More on how to set up stories at: https://storybook.js.org/docs/writing-stories#default-export
const meta: Meta<typeof PreviewCard> = {
  title: "Cards/PreviewCard",
  component: PreviewCard,
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
    data: sampleResult,
  },
};

let incomplete = cloneDeep(sampleResult);
export const WithIncomplete: Story = {
  args: {
    data: merge(incomplete, {
      incomplete_characters: ["xingqiu"],
    }),
  },
};
