import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
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
import { ExecutorProvider } from "./ExecutorProvider";
import { makeExecutor } from "./testExecutor";

const baseProps = {
	config: "cfg text",
	setConfig: () => {},
	isValid: true,
	error: null,
	parsedTeam: [],
};

function renderEditor(
	props: Partial<React.ComponentProps<typeof Editor>> = {},
	fake = makeExecutor(),
) {
	const utils = render(
		<ExecutorProvider exec={fake.supplier}>
			<Editor {...baseProps} {...props} />
		</ExecutorProvider>,
	);
	return { ...utils, fake };
}

beforeEach(() => {
	lastRun = undefined;
	localStorage.clear();
});

describe("Editor", () => {
	it("renders the controlled config and calls setConfig on edits, holding no state", async () => {
		const setConfig = vi.fn();
		const { rerender } = renderEditor({ config: "one", setConfig });
		const ace = screen.getByTestId<HTMLTextAreaElement>("ace");
		expect(ace.value).toBe("one");

		await userEvent.type(ace, "!");
		expect(setConfig).toHaveBeenCalled();

		const fake = makeExecutor();
		rerender(
			<ExecutorProvider exec={fake.supplier}>
				<Editor {...baseProps} config="two" setConfig={setConfig} />
			</ExecutorProvider>,
		);
		expect(screen.getByTestId<HTMLTextAreaElement>("ace").value).toBe("two");
	});

	it("hides gated content until its toggle is persisted", () => {
		renderEditor();
		expect(screen.queryByTestId("editor-team-view")).toBeNull();
	});

	it("shows gated content when the persisted toggle is on", () => {
		localStorage.setItem(
			"gcsim-config-editor-tools",
			JSON.stringify({ team: true, nameSearch: false, tips: false }),
		);
		renderEditor();
		expect(screen.getByTestId("editor-team-view")).toBeTruthy();
	});

	it("disables Run unless the executor is ready and the config is valid", async () => {
		const { rerender } = renderEditor({ isValid: false });
		const run = () => screen.getByRole("button", { name: "simple.run" });
		await waitFor(() => expect(run()).toBeTruthy());
		expect(run()).toBeDisabled();

		const fake = makeExecutor();
		rerender(
			<ExecutorProvider exec={fake.supplier}>
				<Editor {...baseProps} isValid={true} />
			</ExecutorProvider>,
		);
		await waitFor(() => expect(run()).toBeEnabled());
	});

	it("invokes the provider run with the current config from the Run button", async () => {
		const { fake } = renderEditor({ config: "run me" });
		const run = await screen.findByRole("button", { name: "simple.run" });
		await waitFor(() => expect(run).toBeEnabled());
		await userEvent.click(run);
		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("run me", expect.any(Function)),
		);
	});

	it("binds the provider run into the editor surface", async () => {
		const { fake } = renderEditor({ config: "hotkey run" });
		expect(lastRun).toBeTypeOf("function");
		lastRun?.();
		await waitFor(() =>
			expect(fake.run).toHaveBeenCalledWith("hotkey run", expect.any(Function)),
		);
	});

	it("renders the theme selector and settings slot only when asked", () => {
		const { rerender } = renderEditor();
		expect(screen.queryByRole("combobox")).toBeNull();

		const fake = makeExecutor();
		rerender(
			<ExecutorProvider exec={fake.supplier}>
				<Editor
					{...baseProps}
					showThemeSelector
					settings={<div data-testid="host-settings" />}
				/>
			</ExecutorProvider>,
		);
		expect(screen.getByRole("combobox")).toBeTruthy();
		expect(screen.getByTestId("host-settings")).toBeTruthy();
	});
});
