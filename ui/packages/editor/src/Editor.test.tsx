import { diagnosticCount } from "@codemirror/lint";
import { EditorView } from "@codemirror/view";
import { fireEvent, render, screen } from "@testing-library/react";
import { expect, test, vi } from "vitest";
import { Editor } from "./Editor";

function viewOf(container: HTMLElement) {
	const dom = container.querySelector<HTMLElement>(".cm-editor");
	if (!dom) throw new Error("no editor");
	const view = EditorView.findFromDOM(dom);
	if (!view) throw new Error("no view");
	return view;
}

test("renders the value", () => {
	const { container } = render(<Editor value="bennett char lvl=90;" />);
	expect(viewOf(container).state.doc.toString()).toBe("bennett char lvl=90;");
});

test("edits call onChange", () => {
	const onChange = vi.fn();
	const { container } = render(<Editor value="a" onChange={onChange} />);
	const view = viewOf(container);
	view.dispatch({ changes: { from: 1, insert: "b" } });
	expect(onChange).toHaveBeenCalledWith("ab");
});

test("follows value changes without echoing them through onChange", () => {
	const onChange = vi.fn();
	const { container, rerender } = render(
		<Editor value="a" onChange={onChange} />,
	);
	rerender(<Editor value="xyz" onChange={onChange} />);
	expect(viewOf(container).state.doc.toString()).toBe("xyz");
	expect(onChange).not.toHaveBeenCalled();
});

test("value changes edit only the text that differs, keeping the cursor", () => {
	const { container, rerender } = render(<Editor value="active hutao;" />);
	const view = viewOf(container);
	view.dispatch({ selection: { anchor: 7 } });
	rerender(<Editor value={"xq char;\n\nactive hutao;"} />);
	expect(view.state.selection.main.head).toBe(17);
});

test("readOnly blocks editing and can be toggled", () => {
	const { container, rerender } = render(<Editor value="a" readOnly />);
	expect(viewOf(container).state.readOnly).toBe(true);
	rerender(<Editor value="a" />);
	expect(viewOf(container).state.readOnly).toBe(false);
});

test("an error marks the code and lists the whole message in a panel", () => {
	const { container, rerender } = render(
		<Editor
			value={"a\nb"}
			error={"ln2:1: bad\n\tconfig has no targets"}
			errorTitle="Invalid Config"
		/>,
	);
	const view = viewOf(container);
	expect(diagnosticCount(view.state)).toBe(1);
	const panel = screen.getByRole("alert");
	expect(panel).toHaveTextContent("Invalid Config");
	expect(panel).toHaveTextContent("2:1bad");
	expect(panel).toHaveTextContent("config has no targets");

	rerender(<Editor value={"a\nb"} error={null} />);
	expect(diagnosticCount(view.state)).toBe(0);
	expect(screen.queryByRole("alert")).toBeNull();
});

test("an error's position jumps to its line", () => {
	const { container } = render(
		<Editor value={"a\nbcd"} error={"ln2:3: bad"} />,
	);
	fireEvent.click(screen.getByRole("button", { name: "2:3" }));
	const view = viewOf(container);
	expect(view.state.selection.main.head).toBe(4);
	expect(view.hasFocus).toBe(true);
});

test("switches theme in place", () => {
	const { container, rerender } = render(<Editor value="a" theme="app" />);
	const view = viewOf(container);
	const before = view.dom.className;
	rerender(<Editor value="a" theme="github_light" />);
	expect(viewOf(container)).toBe(view);
	expect(view.dom.className).not.toBe(before);
});
