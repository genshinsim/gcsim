import { Editor, ExecutorProvider } from "@gcsim/components";
import type {
	Executor,
	ExecutorSupplier,
	model,
	ParsedResult,
	Sample,
} from "@gcsim/types";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { sampleTeam } from "../samples";
import { sampleConfig } from "./sampleConfig";

const emptyResult: ParsedResult = {
	characters: [],
	errors: [],
	player_initial_pos: { x: 0, y: 0, r: 0 },
};

const fakeExecutor: Executor = {
	ready: () => Promise.resolve(true),
	running: () => false,
	validate: () => Promise.resolve(emptyResult),
	sample: () => Promise.resolve({} as Sample),
	run: () => Promise.resolve(true),
	cancel: () => {},
	buildInfo: () => ({ hash: "", date: "" }),
};
const supplier: ExecutorSupplier<Executor> = () => fakeExecutor;

const TOGGLES_KEY = "gcsim-config-editor-tools";

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
	decorators: [
		(Story, context) => {
			const toggles = context.parameters.editorToggles ?? {
				team: true,
				nameSearch: true,
				tips: true,
			};
			localStorage.setItem(TOGGLES_KEY, JSON.stringify(toggles));
			return (
				<ExecutorProvider exec={supplier}>
					<Story />
				</ExecutorProvider>
			);
		},
	],
	args: {
		config: sampleConfig,
		setConfig: () => {},
		isValid: true,
		error: null,
		parsedTeam: sampleTeam,
		teamCharacters,
		showThemeSelector: true,
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
	parameters: { editorToggles: { team: true, nameSearch: false, tips: false } },
	args: {
		parsedTeam: sampleTeam,
		teamCharacters,
	},
};

export const WithTeamError: Story = {
	parameters: { editorToggles: { team: true, nameSearch: false, tips: false } },
	args: {
		isValid: false,
		error: "invalid action: unknown key 'foo'",
		parsedTeam: sampleTeam,
		teamCharacters,
	},
};

export const AllTools: Story = {
	parameters: { editorToggles: { team: true, nameSearch: true, tips: true } },
	args: {
		parsedTeam: sampleTeam,
		teamCharacters,
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
