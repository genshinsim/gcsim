import { Editor } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleConfig } from "./sampleConfig";

const meta: Meta<typeof Editor> = {
	title: "Editor/Editor",
	component: Editor,
	parameters: {
		layout: "padded",
	},
	tags: ["autodocs"],
	args: {
		config: sampleConfig,
		setConfig: () => {},
		run: () => {},
		isValid: true,
		error: null,
		parsedTeam: [],
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {};

export const WithThemeSelector: Story = {
	args: {
		showThemeSelector: true,
	},
};

export const WithTools: Story = {
	args: {
		showTeam: true,
		showTools: true,
	},
};

export const PrimaryMobile: Story = {
	parameters: {
		viewport: {
			defaultViewport: "mobile1",
		},
	},
};

export const PrimaryTablet: Story = {
	parameters: {
		viewport: {
			defaultViewport: "tablet",
		},
	},
};
