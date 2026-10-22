// @vitest-environment jsdom
import { openSearchPanel, search } from "@codemirror/search";
import { EditorView } from "@codemirror/view";
import { afterEach, expect, it } from "vitest";
import { editorTheme } from "./themes";

let view: EditorView | undefined;
afterEach(() => view?.destroy());

function searchPanel(theme: "dracula" | "github_light") {
	view = new EditorView({
		parent: document.body,
		extensions: [search(), editorTheme(theme)],
	});
	openSearchPanel(view);
	const field = view.dom.querySelector<HTMLElement>(".cm-textfield");
	const button = view.dom.querySelector<HTMLElement>(".cm-button");
	if (!field || !button) throw new Error("no search panel");
	return { field: getComputedStyle(field), button: getComputedStyle(button) };
}

it.each([
	["dracula", "rgb(40, 42, 54)", "rgb(248, 248, 242)", "rgb(52, 55, 70)"],
	[
		"github_light",
		"rgb(255, 255, 255)",
		"rgb(36, 41, 47)",
		"rgb(246, 248, 250)",
	],
] as const)(
	"styles %s search controls from its palette",
	(theme, bg, fg, raised) => {
		const { field, button } = searchPanel(theme);
		expect(field.backgroundColor).toBe(bg);
		expect(field.color).toBe(fg);
		expect(button.backgroundColor).toBe(raised);
		expect(button.backgroundImage).toBe("none");
		expect(button.color).toBe(fg);
	},
);
