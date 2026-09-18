import type { Decorator, Preview } from "@storybook/react-vite";
import React from "react";
import { INITIAL_VIEWPORTS, MINIMAL_VIEWPORTS } from "storybook/viewport";
import "../src/index.css";
import i18n from "./i18n";

const withTheme: Decorator = (Story, context) => {
	const dark = context.globals.theme !== "light";
	return React.createElement(
		"div",
		{
			className: dark ? "dark" : undefined,
			style: {
				backgroundColor: "hsl(var(--deprecated-background))",
				color: "hsl(var(--deprecated-foreground))",
				minHeight: "100vh",
				padding: "1rem",
			},
		},
		React.createElement(Story),
	);
};

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
	discord: {
		name: "discord",
		styles: {
			width: "520",
			height: "250",
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
			options: {
				...INITIAL_VIEWPORTS,
				...MINIMAL_VIEWPORTS,
				...customViewports,
			},
		},
		i18n,
	},
	decorators: [withTheme],
	globalTypes: {
		theme: {
			description: "shadcn theme",
			toolbar: {
				title: "Theme",
				icon: "circlehollow",
				items: [
					{ value: "dark", title: "Dark" },
					{ value: "light", title: "Light" },
				],
				dynamicTitle: true,
			},
		},
	},
	globals: {
		locale: "en",
		theme: "dark",
		locales: {
			en: "English",
			zh: "中文",
			ja: "日本語",
			ko: "한국어",
			es: "Español",
			ru: "Русский",
			de: "Deutsch",
		},
	},
};

export default preview;
