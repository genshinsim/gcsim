import { Button, DBCard } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react";
import _ from "lodash";
import { dbEntries } from "../samples";

// More on how to set up stories at: https://storybook.js.org/docs/writing-stories#default-export
const meta: Meta<typeof DBCard> = {
  title: "Cards/DBCard",
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
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
};

export const PrimaryTablet: Story = {
  args: {
    entry: dbEntries.data[0],
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
  parameters: {
    viewport: {
      defaultViewport: "tablet",
    },
  },
};

export const PrimaryMobile: Story = {
  args: {
    entry: dbEntries.data[0],
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
  parameters: {
    viewport: {
      defaultViewport: "mobile1",
    },
  },
};

export const BGOverride: Story = {
  args: {
    className: "bg-orange-600",
    entry: dbEntries.data[0],
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
};

const longDesc = _.merge(_.cloneDeep(dbEntries.data[0]), {
  description: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ullamcorper eu felis vitae suscipit. Sed suscipit commodo lectus at rutrum. Integer dictum laoreet augue a tincidunt. Praesent rutrum nunc non sodales vulputate. Nunc eget purus tincidunt, euismod ligula nec, feugiat diam. Sed lobortis scelerisque nulla, at ultrices diam viverra quis. Pellentesque leo justo, venenatis id dapibus sit amet, mollis quis sem. Proin leo nunc, commodo a tempus non, vehicula id felis. Donec accumsan non odio at laoreet.

  Vivamus ut tortor lacus. Pellentesque fringilla diam id justo accumsan, eget rutrum eros efficitur. Morbi vel pharetra tellus. Nullam velit libero, efficitur et ultricies vel, auctor sed eros. Fusce facilisis turpis a lacus convallis congue. Aenean in hendrerit diam. Cras nunc magna, efficitur quis mauris quis, lobortis bibendum enim. Duis feugiat tellus id urna commodo varius. In nec tellus augue. Mauris ac pharetra libero.
  
  Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Maecenas sagittis augue id quam lacinia, at dignissim ex porta. Nulla mi neque, scelerisque ut libero nec, tristique imperdiet leo. Aliquam tempus nulla vel maximus dictum. Nullam nec lectus mauris. Pellentesque egestas lacus tortor, ac hendrerit elit pulvinar eget. Aenean finibus, nisl eget lacinia egestas, tellus nulla venenatis enim, vitae maximus nisl massa id sem. Mauris euismod ex nisl, aliquet dignissim enim elementum nec. Sed ac neque vehicula, dapibus urna vel, varius nunc. Sed nulla arcu, luctus non vehicula eu, volutpat ut elit. Nunc justo lacus, varius et gravida non, rhoncus nec odio. Nulla vel tincidunt arcu, quis malesuada arcu. Nullam sed imperdiet dolor, sed tempor tellus. Sed dictum consectetur mauris.
  
  Proin efficitur elit ut ornare pulvinar. Vivamus iaculis diam elit, vitae fermentum risus aliquam sit amet. Vestibulum leo diam, tincidunt ut massa at, rhoncus fermentum mi. Maecenas a maximus sem. Cras quis odio leo. Suspendisse potenti. Integer nulla sem, facilisis posuere gravida a, scelerisque sed ligula. Ut et lorem nec nisl porttitor pretium. Aliquam arcu eros, venenatis vitae condimentum in, pulvinar ac augue.
  
  Suspendisse venenatis magna lectus, a rhoncus eros dapibus vel. Vivamus ullamcorper ligula in justo cursus tempus. Suspendisse potenti. Proin quam leo, gravida vitae finibus a, gravida eget velit. Integer sagittis lectus non porta sagittis. Sed pretium mi vitae velit venenatis ullamcorper. Nullam a lacinia nulla. Sed ultricies sem libero, in fringilla tellus bibendum sed. Aenean in bibendum neque. Nullam sed congue ex, non rhoncus mauris. Nam finibus tellus dictum volutpat tincidunt.`,
});

export const LongDesc: Story = {
  args: {
    entry: longDesc,
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
};

export const LongDescTablet: Story = {
  args: {
    entry: longDesc,
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
  parameters: {
    viewport: {
      defaultViewport: "tablet",
    },
  },
};

export const LongDescMobile: Story = {
  args: {
    entry: longDesc,
    footer: (
      <a href="#">
        <Button className="bg-blue-700 opacity-90">Show Results Page</Button>
      </a>
    ),
  },
  parameters: {
    viewport: {
      defaultViewport: "mobile1",
    },
  },
};
