import type { Meta, StoryObj } from "@storybook/react";
import { Input } from "@gcsim/primitives";

const meta = {
  title: "Primitives/Input",
  component: Input,
  tags: ["autodocs"],
  argTypes: {
    type: { control: "select", options: ["text", "password", "email", "number", "search"] },
    disabled: { control: "boolean" },
    placeholder: { control: "text" },
  },
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { placeholder: "Type something..." },
};

export const Password: Story = {
  args: { type: "password", placeholder: "Enter password..." },
};

export const Disabled: Story = {
  args: { placeholder: "Disabled input", disabled: true },
};

export const WithValue: Story = {
  args: { defaultValue: "Hello world" },
};
