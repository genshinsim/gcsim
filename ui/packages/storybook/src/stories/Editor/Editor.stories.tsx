import { defaultEditorPrefs, Editor } from "@gcsim/components";
import type { model } from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { useArgs } from "storybook/preview-api";
import { fn } from "storybook/test";
import { sampleTeam } from "../samples";
import { sampleConfig } from "./sampleConfig";

const teamCharacters = {
	createCharacter: (key: string): model.Character => ({
		name: key,
		element: "pyro",
		level: 80,
		max_level: 90,
		cons: 0,
		weapon: { name: "dullblade", refine: 1, level: 1, max_level: 20 },
		talents: { attack: 6, skill: 6, burst: 6 },
		sets: {},
		stats: [],
		snapshot: [],
	}),
	imported: [],
};

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
				onPrefsChange={(prefs) => {
					args.onPrefsChange(prefs);
					updateArgs({ prefs });
				}}
			/>
		);
	},
	args: {
		config: sampleConfig,
		setConfig: () => {},
		error: null,
		parsedTeam: sampleTeam,
		teamCharacters,
		showThemeSelector: true,
		onRun: fn(),
		canRun: true,
		prefs: defaultEditorPrefs,
		onPrefsChange: fn(),
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

export const WithTeam: Story = {
	args: {
		prefs: {
			...defaultEditorPrefs,
			toggles: { team: true, nameSearch: false, tips: false },
		},
		parsedTeam: sampleTeam,
		teamCharacters,
	},
};

export const WithTeamError: Story = {
	args: {
		prefs: {
			...defaultEditorPrefs,
			toggles: { team: true, nameSearch: false, tips: false },
		},
		canRun: false,
		error: "invalid action: unknown key 'foo'",
		parsedTeam: sampleTeam,
		teamCharacters,
	},
};

export const AllTools: Story = {
	args: {
		parsedTeam: sampleTeam,
		teamCharacters,
	},
};

export const Busy: Story = {
	args: {
		canRun: false,
		busy: true,
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
