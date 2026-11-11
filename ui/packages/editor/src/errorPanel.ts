import { Facet, StateEffect, StateField } from "@codemirror/state";
import {
	type EditorView,
	type Panel,
	showPanel,
	type ViewUpdate,
} from "@codemirror/view";
import { errorRange, type GcsimError, parseErrorLine } from "./diagnostics";

export interface ErrorPanelValue {
	message: string | null;
	title: string;
}

export type ErrorPanelPosition = "top" | "bottom";

export const setErrorPanel = StateEffect.define<ErrorPanelValue>();

export const errorPanelPosition = Facet.define<
	ErrorPanelPosition,
	ErrorPanelPosition
>({ combine: (values) => values[0] ?? "bottom" });

export const errorPanelField = StateField.define<ErrorPanelValue>({
	create: () => ({ message: null, title: "" }),
	update(value, tr) {
		for (const e of tr.effects) if (e.is(setErrorPanel)) value = e.value;
		return value;
	},
	provide: (field) =>
		showPanel.compute([field, errorPanelPosition], (state) => {
			if (!state.field(field).message) return null;
			return state.facet(errorPanelPosition) === "top"
				? errorPanelTop
				: errorPanel;
		}),
});

function jumpTo(view: EditorView, err: GcsimError) {
	const { from } = errorRange(view.state.doc, err);
	view.dispatch({ selection: { anchor: from }, scrollIntoView: true });
	view.focus();
}

function render(dom: HTMLElement, view: EditorView, value: ErrorPanelValue) {
	const title = document.createElement("div");
	title.className = "cm-gcsim-errors-title";
	title.textContent = value.title;

	const list = document.createElement("ul");
	list.className = "cm-gcsim-errors-list";
	for (const text of (value.message ?? "").split("\n")) {
		if (text.trim() === "") continue;
		const item = document.createElement("li");
		const err = parseErrorLine(text);
		if (err) {
			const pos = document.createElement("button");
			pos.type = "button";
			pos.className = "cm-gcsim-errors-pos";
			pos.textContent = err.column
				? `${err.line}:${err.column}`
				: `${err.line}`;
			pos.onclick = () => jumpTo(view, err);
			item.append(pos, err.message);
		} else {
			item.textContent = text.trim();
		}
		list.append(item);
	}
	dom.replaceChildren(title, list);
}

function errorPanelTop(view: EditorView): Panel {
	return { ...errorPanel(view), top: true };
}

function errorPanel(view: EditorView): Panel {
	const dom = document.createElement("div");
	dom.className = "cm-gcsim-errors";
	dom.setAttribute("role", "alert");
	render(dom, view, view.state.field(errorPanelField));
	return {
		dom,
		update(update: ViewUpdate) {
			const value = update.state.field(errorPanelField);
			if (value !== update.startState.field(errorPanelField)) {
				render(dom, update.view, value);
			}
		},
	};
}
