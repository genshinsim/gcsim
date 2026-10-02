import {
	Editor,
	type EditorPrefs,
	loadEditorPrefs,
	saveEditorPrefs,
} from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { useArgs } from "storybook/preview-api";
import { fn } from "storybook/test";
import { sampleConfig } from "./sampleConfig";

const meta: Meta<typeof Editor> = {
	title: "Editor/Editor",
	component: Editor,
	parameters: {
		layout: "padded",
	},
	tags: ["autodocs"],
	render: function Render(args) {
		const [, updateArgs] = useArgs();
		return (
			<Editor
				{...args}
				onChange={(value) => {
					args.onChange(value);
					updateArgs({ value });
				}}
			/>
		);
	},
	args: {
		value: sampleConfig,
		onChange: fn(),
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

const seedPrefs = (prefs: Partial<EditorPrefs>) => () => {
	const saved = loadEditorPrefs(localStorage);
	saveEditorPrefs(localStorage, { ...saved, ...prefs });
	return () => saveEditorPrefs(localStorage, saved);
};

export const Primary: Story = {};

export const ClassicTheme: Story = {
	beforeEach: seedPrefs({ theme: "monokai" }),
};

export const WithErrors: Story = {
	args: {
		error:
			"ln4:9: unexpected token\nln12: invalid action\n\tconfig does not contain any targets",
	},
};

export const LargeFont: Story = {
	beforeEach: seedPrefs({ fontSize: 20 }),
};

export const GrowToFit: Story = {
	args: {
		maxLines: Infinity,
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
