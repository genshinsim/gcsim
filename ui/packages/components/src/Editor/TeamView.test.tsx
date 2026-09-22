import type { model } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
}));

// biome-ignore lint/suspicious/noExplicitAny: test doubles
vi.mock("../Cards", () => ({
	CharacterCard: ({ char, handleDelete }: any) => (
		<div data-testid="card">
			<span>{char.name}</span>
			<button type="button" onClick={handleDelete}>
				delete-{char.name}
			</button>
		</div>
	),
}));

vi.mock("../common/gcsim", async (orig) => {
	const actual = await orig<typeof import("../common/gcsim")>();
	return {
		...actual,
		characters: ["klee", "amber", "bennett"],
		characterLabel: (k: string) => k,
		// biome-ignore lint/suspicious/noExplicitAny: test double
		OmniSelect: ({ isOpen, items, onSelect }: any) =>
			isOpen ? (
				<div data-testid="picker">
					{/* biome-ignore lint/suspicious/noExplicitAny: test double */}
					{items.map((it: any) => (
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

import { TeamView } from "./TeamView";

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

describe("TeamView", () => {
	it("renders a card per parsed team character", () => {
		render(
			<TeamView
				parsedTeam={[char("amber"), char("bennett")]}
				error={null}
				config=""
				setConfig={() => {}}
			/>,
		);
		expect(screen.getAllByTestId("card")).toHaveLength(2);
	});

	it("surfaces the validation error through a destructive alert", () => {
		render(
			<TeamView
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
			<TeamView
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
			<TeamView
				parsedTeam={[char("amber")]}
				error={null}
				config=""
				setConfig={setConfig}
				characters={source}
			/>,
		);
		await userEvent.click(
			screen.getByRole("button", { name: "db.characters" }),
		);
		await userEvent.click(screen.getByText("default:klee"));
		const written = setConfig.mock.calls[0][0] as string;
		expect(written).toContain("amber char");
		expect(written).toContain("klee char");
	});

	it("hides the add affordance when no character source is injected", () => {
		render(
			<TeamView
				parsedTeam={[char("amber")]}
				error={null}
				config=""
				setConfig={() => {}}
			/>,
		);
		expect(screen.queryByRole("button", { name: "db.characters" })).toBeNull();
	});
});
