import {
  INITIAL_VIEWPORTS,
  MINIMAL_VIEWPORTS,
} from "@storybook/addon-viewport";
import type { Preview } from "@storybook/react";
import "../src/index.css";

const customViewports = {
  desktop1024: {
    name: "desktop-1024",
    styles: {
      width: "1024",
      height: "768",
    },
  },
  desktop1280: {
    name: "desktop-1280",
    styles: {
      width: "1280",
      height: "1024",
    },
  },
  desktop1366: {
    name: "desktop-1366",
    styles: {
      width: "1366",
      height: "768",
    },
  },
  desktop1920: {
    name: "desktop-1920",
    styles: {
      width: "1920",
      height: "1080",
    },
  },
};

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    viewport: {
      viewports: {
        ...INITIAL_VIEWPORTS,
        ...MINIMAL_VIEWPORTS,
        ...customViewports,
      },
      defaultViewport: "desktop",
    },
  },
};

export default preview;
