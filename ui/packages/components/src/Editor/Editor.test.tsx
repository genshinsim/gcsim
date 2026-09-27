import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

let mountedRun: (() => void) | undefined;
vi.mock("./AceEditorWrapper", () => ({
	AceEditorWrapper: ({
		cfg,
		onChange,
		onRun,
	}: {
		cfg: string;
		onChange: (v: string) => void;
		onRun?: () => void;
	}) => {
		mountedRun ??= onRun;
		return (
			<textarea
				data-testid="ace"
				value={cfg}
				onChange={(e) => onChange(e.currentTarget.value)}
			/>
		);
	},
}));

import { Editor } from "./Editor";
import { defaultEditorPrefs } from "./types";

const baseProps: React.ComponentProps<typeof Editor> = {
	config: "cfg text",
	setConfig: () => {},
	error: null,
	parsedTeam: [],
	onRun: () => {},
	canRun: true,
	prefs: defaultEditorPrefs,
	onPrefsChange: () => {},
};

function renderEditor(
	props: Partial<React.ComponentProps<typeof Editor>> = {},
) {
	return render(<Editor {...baseProps} {...props} />);
}

const runButton = () => screen.getByRole("button", { name: "simple.run" });

beforeEach(() => {
	mountedRun = undefined;
});

describe("Editor", () => {
	it("renders the controlled config and calls setConfig on edits, holding no state", async () => {
		const setConfig = vi.fn();
		const { rerender } = renderEditor({ config: "one", setConfig });
		const ace = screen.getByTestId<HTMLTextAreaElement>("ace");
		expect(ace.value).toBe("one");

		await userEvent.type(ace, "!");
		expect(setConfig).toHaveBeenCalled();

		rerender(<Editor {...baseProps} config="two" setConfig={setConfig} />);
		expect(screen.getByTestId<HTMLTextAreaElement>("ace").value).toBe("two");
	});

	it("reveals the team by default to match the live route", () => {
		renderEditor();
		expect(screen.getByTestId("editor-team-composer")).toBeTruthy();
	});

	it("hides gated content when its toggle pref is off", () => {
		renderEditor({
			prefs: {
				...defaultEditorPrefs,
				toggles: { team: false, nameSearch: false, tips: false },
			},
		});
		expect(screen.queryByTestId("editor-team-composer")).toBeNull();
	});

	it("reports a toggled tool through onPrefsChange", async () => {
		const onPrefsChange = vi.fn();
		renderEditor({ onPrefsChange });
		await userEvent.click(
			screen.getAllByRole("button", { name: "simple.hide_all_tips" })[0],
		);
		expect(onPrefsChange).toHaveBeenCalledWith({
			...defaultEditorPrefs,
			toggles: { ...defaultEditorPrefs.toggles, tips: false },
		});
	});

	it("reports theme and font size changes through onPrefsChange", async () => {
		const onPrefsChange = vi.fn();
		renderEditor({ onPrefsChange, showThemeSelector: true });
		await userEvent.selectOptions(screen.getByRole("combobox"), "github");
		expect(onPrefsChange).toHaveBeenLastCalledWith({
			...defaultEditorPrefs,
			theme: "github",
		});
	});

	it("disables Run while canRun is false", () => {
		const { rerender } = renderEditor({ canRun: false });
		expect(runButton()).toBeDisabled();
		rerender(<Editor {...baseProps} canRun={true} />);
		expect(runButton()).toBeEnabled();
	});

	it("calls onRun from the Run button", async () => {
		const onRun = vi.fn();
		renderEditor({ onRun });
		await userEvent.click(runButton());
		expect(onRun).toHaveBeenCalledTimes(1);
	});

	it("routes the hotkey to the latest onRun, only while canRun", () => {
		const first = vi.fn();
		const latest = vi.fn();
		const { rerender } = renderEditor({ onRun: first, canRun: false });
		mountedRun?.();
		expect(first).not.toHaveBeenCalled();

		rerender(<Editor {...baseProps} onRun={latest} canRun={true} />);
		mountedRun?.();
		expect(first).not.toHaveBeenCalled();
		expect(latest).toHaveBeenCalledTimes(1);
	});

	it("renders the theme selector and settings slot only when asked", () => {
		const { rerender } = renderEditor();
		expect(screen.queryByRole("combobox")).toBeNull();

		rerender(
			<Editor
				{...baseProps}
				showThemeSelector
				settings={<div data-testid="host-settings" />}
			/>,
		);
		expect(screen.getByRole("combobox")).toBeTruthy();
		expect(screen.getByTestId("host-settings")).toBeTruthy();
	});
});
