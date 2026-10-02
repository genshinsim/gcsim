// @vitest-environment jsdom
import "@gcsim/components/vitest.setup";
import { cleanup, render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k, i18n: { language: "en" } }),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import { SimulatorLayout, type SimulatorLayoutProps } from "./SimulatorLayout";

const baseProps: SimulatorLayoutProps = {
	config: "cfg text",
	setConfig: () => {},
	error: null,
	parsedTeam: [],
	onRun: () => {},
	canRun: true,
	busy: false,
};

function renderLayout(props: Partial<SimulatorLayoutProps> = {}) {
	return render(<SimulatorLayout {...baseProps} {...props} />);
}

const runButton = () => screen.getByRole("button", { name: "simple.run" });

afterEach(cleanup);

describe("SimulatorLayout", () => {
	it("always shows the team and action list, with no tools menu", () => {
		renderLayout();
		expect(screen.getByTestId("editor-team-composer")).toBeInTheDocument();
		expect(screen.getByText("simple.action_list")).toBeInTheDocument();
		expect(screen.queryByTestId("editor-helper-tools")).toBeNull();
	});

	it("shows the parsed team with no add or remove controls", () => {
		renderLayout({
			parsedTeam: [
				{
					name: "amber",
					level: 80,
					max_level: 90,
					element: "pyro",
					cons: 0,
					weapon: { name: "dullblade", refine: 1, level: 1, max_level: 20 },
					talents: { attack: 6, skill: 6, burst: 6 },
					sets: {},
					stats: [],
					snapshot: [],
				},
			],
		});
		const team = within(screen.getByTestId("editor-team-composer"));
		expect(team.getByText(/game:character_names.amber/)).toBeInTheDocument();
		expect(team.getAllByRole("button")).toHaveLength(1);
	});

	it("shows the config error once, under the editor", () => {
		const { container } = renderLayout({ error: "bad line 3" });
		expect(screen.getAllByText("bad line 3")).toHaveLength(1);
		expect(
			container.querySelector("#config_editor [role=alert]"),
		).toHaveTextContent("bad line 3");
	});

	it("disables Run while canRun is false", () => {
		const { rerender } = renderLayout({ canRun: false });
		expect(runButton()).toBeDisabled();
		rerender(<SimulatorLayout {...baseProps} canRun={true} />);
		expect(runButton()).toBeEnabled();
	});

	it("calls onRun from the Run button", async () => {
		const onRun = vi.fn();
		renderLayout({ onRun });
		await userEvent.click(runButton());
		expect(onRun).toHaveBeenCalledTimes(1);
	});

	it("renders the settings slot", () => {
		renderLayout({ settings: <div data-testid="host-settings" /> });
		expect(screen.getByTestId("host-settings")).toBeInTheDocument();
	});
});
