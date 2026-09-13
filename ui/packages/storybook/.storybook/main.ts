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
	viteFinal: async (config) => {
		// The portrait compositor imports `@cf-wasm/photon` bare, which resolves to
		// the auto-initialising `workerd` build (its `.wasm` sync-import doesn't
		// work under Vite). In the browser we want the `others` build, which is
		// hand-initialised from a `?url` wasm — see the CompositedLive story. Alias
		// only the exact bare specifier (not the `./photon.wasm` subpath).
		config.resolve ??= {};
		const alias = config.resolve.alias;
		const entry = {
			find: /^@cf-wasm\/photon$/,
			replacement: "@cf-wasm/photon/others",
		};
		if (Array.isArray(alias)) {
			alias.push(entry);
		} else {
			config.resolve.alias = [
				...Object.entries(alias ?? {}).map(([find, replacement]) => ({
					find,
					replacement: replacement as string,
				})),
				entry,
			];
		}
		return config;
	},
};

export default config;
