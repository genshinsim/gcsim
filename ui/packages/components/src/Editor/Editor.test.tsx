import { diagnosticCount } from "@codemirror/lint";
import { EditorView } from "@codemirror/view";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
}));

import { Editor } from "./Editor";
import { MAX_FONT_SIZE } from "./EditorSettings";

const baseProps: React.ComponentProps<typeof Editor> = {
	value: "cfg text",
	onChange: () => {},
	fontSize: 14,
	theme: "app",
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

	it("hides the settings menu without a change handler", () => {
		render(<Editor {...baseProps} />);
		expect(
			screen.queryByRole("button", { name: "simple.editor_settings" }),
		).toBeNull();
	});

	it("steps the font size within bounds", async () => {
		const onAppearanceChange = vi.fn();
		const { rerender } = render(
			<Editor {...baseProps} onAppearanceChange={onAppearanceChange} />,
		);
		await openSettings();
		await userEvent.click(
			screen.getByRole("button", { name: "simple.decrease_font_size" }),
		);
		expect(onAppearanceChange).toHaveBeenLastCalledWith({
			fontSize: 13,
			theme: "app",
		});

		rerender(
			<Editor
				{...baseProps}
				fontSize={MAX_FONT_SIZE}
				onAppearanceChange={onAppearanceChange}
			/>,
		);
		expect(
			screen.getByRole("button", { name: "simple.increase_font_size" }),
		).toBeDisabled();
	});

	it("picks a theme from the menu", async () => {
		const onAppearanceChange = vi.fn();
		render(<Editor {...baseProps} onAppearanceChange={onAppearanceChange} />);
		await openSettings();
		expect(screen.getByRole("radio", { name: "gcsim" })).toBeChecked();
		await userEvent.click(screen.getByRole("radio", { name: "Dracula" }));
		expect(onAppearanceChange).toHaveBeenLastCalledWith({
			fontSize: 14,
			theme: "dracula",
		});
	});
});
