// @vitest-environment jsdom
import "@testing-library/jest-dom/vitest";
import { defaultEditorPrefs } from "@gcsim/components";
import { cleanup, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
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
	prefs: defaultEditorPrefs,
	onPrefsChange: () => {},
};

function renderLayout(props: Partial<SimulatorLayoutProps> = {}) {
	return render(<SimulatorLayout {...baseProps} {...props} />);
}

const runButton = () => screen.getByRole("button", { name: "simple.run" });

afterEach(cleanup);

describe("SimulatorLayout", () => {
	it("reveals the team, name search, tips and tools by default", () => {
		renderLayout();
		expect(screen.getByTestId("editor-team-composer")).toBeInTheDocument();
		expect(screen.getByText("simple.name_search")).toBeInTheDocument();
		expect(
			screen.getAllByRole("button", { name: "simple.hide_all_tips" }),
		).toHaveLength(2);
		expect(screen.getByTestId("editor-helper-tools")).toBeInTheDocument();
		expect(screen.getByText("simple.action_list")).toBeInTheDocument();
	});

	it("hides gated content when its toggle pref is off", () => {
		renderLayout({
			prefs: {
				...defaultEditorPrefs,
				toggles: { team: false, nameSearch: false, tips: false },
			},
		});
		expect(screen.queryByTestId("editor-team-composer")).toBeNull();
		expect(screen.queryByText("simple.name_search")).toBeNull();
		expect(
			screen.queryByRole("button", { name: "simple.hide_all_tips" }),
		).toBeNull();
		expect(screen.getByTestId("editor-helper-tools")).toBeInTheDocument();
	});

	it("shows the config error when the team is hidden", () => {
		renderLayout({
			error: "bad line 3",
			prefs: {
				...defaultEditorPrefs,
				toggles: { ...defaultEditorPrefs.toggles, team: false },
			},
		});
		expect(screen.getByText("bad line 3")).toBeInTheDocument();
	});

	it("shows the config error once when the team is visible", () => {
		renderLayout({ error: "bad line 3" });
		expect(screen.getAllByText("bad line 3")).toHaveLength(1);
	});

	it("reports hidden tips through onPrefsChange", async () => {
		const onPrefsChange = vi.fn();
		renderLayout({ onPrefsChange });
		await userEvent.click(
			screen.getAllByRole("button", { name: "simple.hide_all_tips" })[0],
		);
		expect(onPrefsChange).toHaveBeenCalledWith({
			...defaultEditorPrefs,
			toggles: { ...defaultEditorPrefs.toggles, tips: false },
		});
	});

	it("shows the editor toolbar and stores theme changes in prefs", async () => {
		const onPrefsChange = vi.fn();
		renderLayout({ onPrefsChange });
		await userEvent.selectOptions(screen.getByRole("combobox"), "github");
		expect(onPrefsChange).toHaveBeenLastCalledWith({
			...defaultEditorPrefs,
			theme: "github",
		});
	});

	it("stores font size changes in prefs", async () => {
		const onPrefsChange = vi.fn();
		renderLayout({ onPrefsChange });
		await userEvent.type(screen.getByRole("spinbutton"), "6");
		expect(onPrefsChange).toHaveBeenLastCalledWith({
			...defaultEditorPrefs,
			fontSize: 146,
		});
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
