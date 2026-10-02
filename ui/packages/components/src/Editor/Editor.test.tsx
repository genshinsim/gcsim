import { diagnosticCount } from "@codemirror/lint";
import { EditorView } from "@codemirror/view";
import type { model } from "@gcsim/types";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k, i18n: { language: "en" } }),
}));

import { Editor } from "./Editor";
import { MAX_FONT_SIZE, saveEditorPrefs } from "./editorPrefs";

const baseProps: React.ComponentProps<typeof Editor> = {
	value: "cfg text",
	onChange: () => {},
};

function view(container: HTMLElement) {
	const dom = container.querySelector<HTMLElement>(".cm-editor");
	const v = dom && EditorView.findFromDOM(dom);
	if (!v) throw new Error("no editor");
	return v;
}

const openSettings = () =>
	userEvent.click(
		screen.getByRole("button", { name: "simple.editor_settings" }),
	);

describe("Editor", () => {
	beforeEach(() => localStorage.clear());

	it("renders the controlled value and calls onChange on edits", () => {
		const onChange = vi.fn();
		const { container, rerender } = render(
			<Editor {...baseProps} value="one" onChange={onChange} />,
		);
		expect(view(container).state.doc.toString()).toBe("one");

		view(container).dispatch({ changes: { from: 3, insert: "!" } });
		expect(onChange).toHaveBeenCalledWith("one!");

		rerender(<Editor {...baseProps} value="two" onChange={onChange} />);
		expect(view(container).state.doc.toString()).toBe("two");
	});

	it("marks positioned validation errors in the editor", () => {
		const { container, rerender } = render(
			<Editor
				{...baseProps}
				value={"a\nb\nc"}
				error={"ln2:1: bad\n\tconfig does not contain any targets"}
			/>,
		);
		expect(diagnosticCount(view(container).state)).toBe(1);

		rerender(<Editor {...baseProps} value={"a\nb\nc"} error={null} />);
		expect(diagnosticCount(view(container).state)).toBe(0);
	});

	it("shows the settings menu with the saved prefs", async () => {
		saveEditorPrefs(localStorage, { theme: "dracula", fontSize: 18 });
		render(<Editor {...baseProps} />);
		await openSettings();
		expect(screen.getByRole("radio", { name: "Dracula" })).toBeChecked();
		expect(screen.getByRole("status")).toHaveTextContent("18");
	});

	it("steps the font size within bounds", async () => {
		render(<Editor {...baseProps} />);
		await openSettings();
		await userEvent.click(
			screen.getByRole("button", { name: "simple.decrease_font_size" }),
		);
		expect(screen.getByRole("status")).toHaveTextContent("13");
	});

	it("disables increasing past the largest font size", async () => {
		saveEditorPrefs(localStorage, { theme: "app", fontSize: MAX_FONT_SIZE });
		render(<Editor {...baseProps} />);
		await openSettings();
		expect(
			screen.getByRole("button", { name: "simple.increase_font_size" }),
		).toBeDisabled();
	});

	it("picks a theme from the menu", async () => {
		render(<Editor {...baseProps} />);
		await openSettings();
		expect(screen.getByRole("radio", { name: "gcsim" })).toBeChecked();
		await userEvent.click(screen.getByRole("radio", { name: "Dracula" }));
		expect(screen.getByRole("radio", { name: "Dracula" })).toBeChecked();
	});

	it("keeps settings for an editor mounted later", async () => {
		const { unmount } = render(<Editor {...baseProps} />);
		await openSettings();
		await userEvent.click(screen.getByRole("radio", { name: "Dracula" }));
		await userEvent.click(
			screen.getByRole("button", { name: "simple.increase_font_size" }),
		);
		unmount();

		render(<Editor {...baseProps} />);
		await openSettings();
		expect(screen.getByRole("radio", { name: "Dracula" })).toBeChecked();
		expect(screen.getByRole("status")).toHaveTextContent("15");
	});

	it("formats the config from the toolbar", async () => {
		const user = userEvent.setup();
		const onChange = vi.fn();
		render(
			<Editor
				{...baseProps}
				value={"while true {\nhutao burst;\n}"}
				onChange={onChange}
			/>,
		);
		await user.click(
			screen.getByRole("button", { name: "simple.format_config" }),
		);
		expect(onChange).toHaveBeenLastCalledWith(
			"while true {\n\thutao burst;\n}\n",
		);
	});

	it("searches every name by shortcut and copies the key", async () => {
		const user = userEvent.setup();
		const writeText = vi
			.spyOn(navigator.clipboard, "writeText")
			.mockResolvedValue();
		render(<Editor {...baseProps} />);
		await user.click(
			screen.getByRole("button", { name: "simple.shortcut_search" }),
		);
		await user.type(screen.getByRole("combobox"), "ht");
		await user.click(screen.getByRole("option", { name: /^hutao ht/ }));
		expect(writeText).toHaveBeenCalledWith("hutao");
		expect(screen.queryByRole("dialog")).toBeNull();
	});

	it("inserts an imported character block at the top", async () => {
		const user = userEvent.setup();
		const onChange = vi.fn();
		const character = {
			name: "hutao",
			level: 90,
			max_level: 90,
			cons: 1,
			talents: { attack: 10, skill: 10, burst: 10 },
			weapon: { name: "homa", refine: 1, level: 90, max_level: 90 },
			sets: {},
			stats: [],
		} as unknown as model.Character;
		render(
			<Editor
				{...baseProps}
				value="active hutao;"
				onChange={onChange}
				importedCharacters={[{ key: "hutao", character }]}
			/>,
		);
		await user.click(
			screen.getByRole("button", { name: "simple.imported_characters" }),
		);
		await user.click(screen.getByRole("option", { name: /hutao/ }));
		expect(onChange).toHaveBeenCalledWith(
			'hutao char lvl=90/90 cons=1 talent=10,10,10;\nhutao add weapon="homa" refine=1 lvl=90/90;\n\nactive hutao;',
		);
	});
});
