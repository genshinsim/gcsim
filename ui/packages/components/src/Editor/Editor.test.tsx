import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({ t: (k: string) => k }),
}));

vi.mock("./AceEditorWrapper", () => ({
	themes: ["github", "tomorrow_night"],
	AceEditorWrapper: ({
		value,
		onChange,
		theme,
		fontSize,
		maxLines,
	}: {
		value: string;
		onChange: (v: string) => void;
		theme?: string;
		fontSize?: number;
		maxLines?: number;
	}) => (
		<textarea
			data-testid="ace"
			data-theme={theme}
			data-font-size={fontSize}
			data-max-lines={maxLines}
			value={value}
			onChange={(e) => onChange(e.currentTarget.value)}
		/>
	),
}));

import { Editor } from "./Editor";

const baseProps: React.ComponentProps<typeof Editor> = {
	value: "cfg text",
	onChange: () => {},
	theme: "tomorrow_night",
	fontSize: 14,
};

const ace = () => screen.getByTestId<HTMLTextAreaElement>("ace");

describe("Editor", () => {
	it("renders the controlled value and calls onChange on edits, holding no state", async () => {
		const onChange = vi.fn();
		const { rerender } = render(
			<Editor {...baseProps} value="one" onChange={onChange} />,
		);
		expect(ace().value).toBe("one");

		await userEvent.type(ace(), "!");
		expect(onChange).toHaveBeenCalledWith("one!");
		expect(ace().value).toBe("one");

		rerender(<Editor {...baseProps} value="two" onChange={onChange} />);
		expect(ace().value).toBe("two");
	});

	it("applies the theme and font size to the text area", () => {
		render(<Editor {...baseProps} theme="github" fontSize={18} />);
		expect(ace().dataset.theme).toBe("github");
		expect(ace().dataset.fontSize).toBe("18");
	});

	it("passes the line limit to the text area", () => {
		render(<Editor {...baseProps} maxLines={Infinity} />);
		expect(ace().dataset.maxLines).toBe("Infinity");
	});

	it("hides the appearance toolbar without a change handler", () => {
		render(<Editor {...baseProps} />);
		expect(screen.queryByRole("combobox")).toBeNull();
		expect(screen.queryByRole("spinbutton")).toBeNull();
	});

	it("reports theme changes through onAppearanceChange", async () => {
		const onAppearanceChange = vi.fn();
		render(<Editor {...baseProps} onAppearanceChange={onAppearanceChange} />);
		await userEvent.selectOptions(screen.getByRole("combobox"), "github");
		expect(onAppearanceChange).toHaveBeenLastCalledWith({
			theme: "github",
			fontSize: 14,
		});
	});

	it("reports font size changes through onAppearanceChange", async () => {
		const onAppearanceChange = vi.fn();
		render(<Editor {...baseProps} onAppearanceChange={onAppearanceChange} />);
		const input = screen.getByRole("spinbutton");
		await userEvent.type(input, "6");
		expect(onAppearanceChange).toHaveBeenLastCalledWith({
			theme: "tomorrow_night",
			fontSize: 146,
		});
	});
});
