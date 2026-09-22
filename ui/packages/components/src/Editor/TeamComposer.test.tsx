import type { model } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

type PickerItem = { key: string; source: string; text: string; label: string };

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
}));

vi.mock("../Cards", () => ({
	TeamCard: ({
		team,
		handleRemove,
		handleAdd,
	}: {
		team: model.Character[];
		handleRemove: (index: number) => () => void;
		handleAdd?: () => void;
	}) => (
		<div data-testid="team-card">
			{team.map((c, index) => (
				<button
					key={c.name ?? index}
					type="button"
					onClick={handleRemove(index)}
				>
					delete-{c.name}
				</button>
			))}
			{handleAdd ? (
				<button type="button" onClick={handleAdd}>
					add
				</button>
			) : null}
		</div>
	),
}));

vi.mock("../common/gcsim", async (orig) => {
	const actual = await orig<typeof import("../common/gcsim")>();
	return {
		...actual,
		characters: ["klee", "amber", "bennett"],
		characterLabel: (k: string) => k,
		OmniSelect: ({
			isOpen,
			items,
			onSelect,
		}: {
			isOpen: boolean;
			items: PickerItem[];
			onSelect: (item: PickerItem) => void;
		}) =>
			isOpen ? (
				<div data-testid="picker">
					{items.map((it) => (
						<button
							type="button"
							key={`${it.source}-${it.key}`}
							onClick={() => onSelect(it)}
						>
							{`${it.source}:${it.key}`}
						</button>
					))}
				</div>
			) : null,
	};
});

import { TeamComposer } from "./TeamComposer";

function char(name: string): model.Character {
	return {
		name,
		level: 80,
		max_level: 90,
		element: "pyro",
		cons: 0,
		weapon: { name: "dullblade", refine: 1, level: 1, max_level: 20 },
		talents: { attack: 6, skill: 6, burst: 6 },
		sets: {},
		stats: new Array(22).fill(0),
		snapshot: new Array(22).fill(0),
	};
}

const source = {
	createCharacter: (key: string) => char(key),
	imported: [],
};

beforeEach(() => {
	vi.clearAllMocks();
});

describe("TeamComposer", () => {
	it("renders the team through TeamCard", () => {
		render(
			<TeamComposer
				parsedTeam={[char("amber"), char("bennett")]}
				error={null}
				config=""
				setConfig={() => {}}
			/>,
		);
		expect(screen.getByText("delete-amber")).toBeTruthy();
		expect(screen.getByText("delete-bennett")).toBeTruthy();
	});

	it("surfaces the validation error through a destructive alert", () => {
		render(
			<TeamComposer
				parsedTeam={[]}
				error="bad action list"
				config=""
				setConfig={() => {}}
			/>,
		);
		expect(screen.getByText("bad action list")).toBeTruthy();
	});

	it("removes a card by rewriting the config", async () => {
		const setConfig = vi.fn();
		render(
			<TeamComposer
				parsedTeam={[char("amber"), char("bennett")]}
				error={null}
				config="amber char lvl=1/1 cons=0 talent=1,1,1;\ntarget lvl=100;"
				setConfig={setConfig}
				characters={source}
			/>,
		);
		await userEvent.click(screen.getByText("delete-amber"));
		expect(setConfig).toHaveBeenCalledTimes(1);
		const written = setConfig.mock.calls[0][0] as string;
		expect(written).toContain("bennett char");
		expect(written).not.toContain("amber char");
	});

	it("adds a character through the picker by rewriting the config", async () => {
		const setConfig = vi.fn();
		render(
			<TeamComposer
				parsedTeam={[char("amber")]}
				error={null}
				config=""
				setConfig={setConfig}
				characters={source}
			/>,
		);
		await userEvent.click(screen.getByRole("button", { name: "add" }));
		await userEvent.click(screen.getByText("default:klee"));
		const written = setConfig.mock.calls[0][0] as string;
		expect(written).toContain("amber char");
		expect(written).toContain("klee char");
	});

	it("hides the add affordance when no character source is injected", () => {
		render(
			<TeamComposer
				parsedTeam={[char("amber")]}
				error={null}
				config=""
				setConfig={() => {}}
			/>,
		);
		expect(screen.queryByRole("button", { name: "add" })).toBeNull();
	});
});
