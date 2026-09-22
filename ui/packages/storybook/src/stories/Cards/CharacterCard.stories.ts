import { CharacterCard, type CharStatBlock } from "@gcsim/components";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";
import { sampleTeam } from "../samples";

const sampleStats: CharStatBlock[] = [
	{ key: "hp", name: "HP", t: "both", flat: 30112, percent: 0.466 },
	{ key: "atk", name: "ATK", t: "both", flat: 1201, percent: 0.331 },
	{ key: "def", name: "DEF", t: "f", flat: 687, percent: 0 },
	{ key: "em", name: "Elemental Mastery", t: "f", flat: 40, percent: 0 },
	{ key: "er", name: "Energy Recharge", t: "%", flat: 0, percent: 1.678 },
	{ key: "cr", name: "CRIT Rate", t: "%", flat: 0, percent: 0.626 },
	{ key: "cd", name: "CRIT DMG", t: "%", flat: 0, percent: 1.214 },
	{ key: "hydro", name: "Hydro DMG Bonus", t: "%", flat: 0, percent: 0.466 },
];

const sampleSnapshot: CharStatBlock[] = [
	{ key: "hp", name: "HP", t: "f", flat: 31204, percent: 0 },
	{ key: "atk", name: "ATK", t: "f", flat: 1407, percent: 0 },
	{ key: "em", name: "Elemental Mastery", t: "f", flat: 40, percent: 0 },
	{ key: "er", name: "Energy Recharge", t: "%", flat: 0, percent: 1.678 },
	{ key: "cr", name: "CRIT Rate", t: "%", flat: 0, percent: 0.626 },
	{ key: "cd", name: "CRIT DMG", t: "%", flat: 0, percent: 1.214 },
	{ key: "hydro", name: "Hydro DMG Bonus", t: "%", flat: 0, percent: 0.466 },
];

const meta: Meta<typeof CharacterCard> = {
	title: "Cards/CharacterCard",
	component: CharacterCard,
	parameters: {
		layout: "fullscreen",
	},
	tags: ["autodocs"],
	argTypes: {},
	args: {
		char: sampleTeam[0],
		stats: sampleStats,
		snapshot: sampleSnapshot,
		statsRows: sampleStats.length,
		name: "Neuvillette",
		constellationLabel: "C0",
		levelLabel: "Lvl",
		talentsLabel: "Talents",
		artifactStatsLabel: "Artifact Stats (w/o Set Effects)",
		totalStatsLabel: "Total Stats (at 0 seconds, best effort basis)",
		weaponName: "Tome of the Eternal Flow",
		className: "max-w-sm",
		handleDelete: fn(),
		handleToggleDetail: fn(),
		handleToggleSnapshot: fn(),
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {},
};

export const Skeleton: Story = {
	args: {
		isSkeleton: true,
	},
};

export const ViewerMode: Story = {
	args: {
		viewerMode: true,
		showSnapshot: true,
	},
};
