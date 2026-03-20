import type { Preview } from "@storybook/react";
import "../src/storybook.css";

const preview: Preview = {
  parameters: {
    backgrounds: {
      default: "dark",
      values: [
        { name: "dark", value: "#1c2127" },
        { name: "light", value: "#ffffff" },
      ],
    },
  },
};

export default preview;
