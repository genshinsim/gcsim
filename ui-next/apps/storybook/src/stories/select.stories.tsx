import type { Meta, StoryObj } from "@storybook/react";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@gcsim/primitives";

const meta = {
  title: "Primitives/Select",
  component: Select,
  tags: ["autodocs"],
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Select>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select mode" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Execution Mode</SelectLabel>
          <SelectItem value="wasm">WASM</SelectItem>
          <SelectItem value="server">Server</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  ),
};

export const WithGroups: Story = {
  render: () => (
    <Select>
      <SelectTrigger className="w-[220px]">
        <SelectValue placeholder="Select element" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Elements</SelectLabel>
          <SelectItem value="pyro">Pyro</SelectItem>
          <SelectItem value="hydro">Hydro</SelectItem>
          <SelectItem value="electro">Electro</SelectItem>
          <SelectItem value="cryo">Cryo</SelectItem>
          <SelectItem value="anemo">Anemo</SelectItem>
          <SelectItem value="geo">Geo</SelectItem>
          <SelectItem value="dendro">Dendro</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  ),
};
