import type { StorybookConfig } from "@storybook/react-vite";

const config: StorybookConfig = {
	stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
	addons: [
		"@chromatic-com/storybook",
		"@storybook/addon-docs",
		"storybook-react-i18next",
	],
	framework: {
		name: "@storybook/react-vite",
		options: {},
	},
	typescript: {
		// Overrides the default Typescript configuration to allow multi-package components to be documented via Autodocs.
		reactDocgen: "react-docgen",
		check: false,
	},
};

export default config;
