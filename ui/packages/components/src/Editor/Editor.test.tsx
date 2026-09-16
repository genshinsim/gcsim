import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
}));

let lastRun: (() => void) | undefined;
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
		lastRun = onRun;
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

const baseProps = {
	config: "cfg text",
	setConfig: () => {},
	isValid: true,
	error: null,
	parsedTeam: [],
	run: () => {},
};

beforeEach(() => {
	lastRun = undefined;
	localStorage.clear();
});

describe("Editor", () => {
	it("renders the controlled config and calls setConfig on edits, holding no state", async () => {
		const setConfig = vi.fn();
		const { rerender } = render(
			<Editor {...baseProps} config="one" setConfig={setConfig} />,
		);
		const ace = screen.getByTestId<HTMLTextAreaElement>("ace");
		expect(ace.value).toBe("one");

		await userEvent.type(ace, "!");
		expect(setConfig).toHaveBeenCalled();
		expect(screen.getByTestId<HTMLTextAreaElement>("ace").value).toBe("one");

		rerender(<Editor {...baseProps} config="two" setConfig={setConfig} />);
		expect(screen.getByTestId<HTMLTextAreaElement>("ace").value).toBe("two");
	});

	it("is a bare editor with all tool flags off", () => {
		render(<Editor {...baseProps} />);
		expect(screen.getByTestId("ace")).toBeTruthy();
		expect(screen.queryByTestId("editor-team-view")).toBeNull();
		expect(screen.queryByTestId("editor-helper-tools")).toBeNull();
		expect(screen.queryByRole("combobox")).toBeNull();
	});

	it("renders the team view only when showTeam is set", () => {
		const { rerender } = render(<Editor {...baseProps} />);
		expect(screen.queryByTestId("editor-team-view")).toBeNull();
		rerender(<Editor {...baseProps} showTeam />);
		expect(screen.getByTestId("editor-team-view")).toBeTruthy();
	});

	it("renders the helper tools only when showTools is set", () => {
		const { rerender } = render(<Editor {...baseProps} />);
		expect(screen.queryByTestId("editor-helper-tools")).toBeNull();
		rerender(<Editor {...baseProps} showTools />);
		expect(screen.getByTestId("editor-helper-tools")).toBeTruthy();
	});

	it("renders the theme selector and settings slot only when asked", () => {
		const { rerender } = render(<Editor {...baseProps} />);
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

	it("binds the host run callback into the editor surface", () => {
		const run = vi.fn();
		render(<Editor {...baseProps} run={run} />);
		expect(lastRun).toBe(run);
	});
});
