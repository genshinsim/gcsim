import { defaultEditorPrefs, Editor } from "@gcsim/components";
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
				onAppearanceChange={
					args.onAppearanceChange
						? (appearance) => {
								args.onAppearanceChange?.(appearance);
								updateArgs(appearance);
							}
						: undefined
				}
			/>
		);
	},
	args: {
		value: sampleConfig,
		onChange: fn(),
		theme: defaultEditorPrefs.theme,
		fontSize: defaultEditorPrefs.fontSize,
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {};

export const WithToolbar: Story = {
	args: {
		onAppearanceChange: fn(),
	},
};

export const LightTheme: Story = {
	args: {
		theme: "github",
	},
};

export const LargeFont: Story = {
	args: {
		fontSize: 20,
	},
};

export const PrimaryMobile: Story = {
	args: {
		onAppearanceChange: fn(),
	},
	parameters: {
		viewport: {
			defaultViewport: "mobile1",
		},
	},
};
