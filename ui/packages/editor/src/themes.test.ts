import { openSearchPanel, search } from "@codemirror/search";
import { EditorView } from "@codemirror/view";
import { afterEach, describe, expect, it } from "vitest";
import {
	DEFAULT_EDITOR_THEME,
	EDITOR_THEMES,
	type EditorThemeId,
	editorTheme,
	isEditorThemeId,
} from "./themes";

describe("editor themes", () => {
	it("have unique ids and include the default", () => {
		const ids = EDITOR_THEMES.map((th) => th.id);
		expect(new Set(ids).size).toBe(ids.length);
		expect(ids).toContain(DEFAULT_EDITOR_THEME);
	});

	it("reuses the built extension per theme", () => {
		expect(editorTheme("monokai")).toBe(editorTheme("monokai"));
		expect(editorTheme("monokai")).not.toBe(editorTheme("app"));
	});

	it("validates stored ids", () => {
		expect(isEditorThemeId("dracula")).toBe(true);
		expect(isEditorThemeId("xcode")).toBe(false);
		expect(isEditorThemeId(undefined)).toBe(false);
	});
});

describe("search panel controls", () => {
	let view: EditorView | undefined;
	afterEach(() => view?.destroy());

	function searchPanel(theme: EditorThemeId) {
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
	] as const)("follow the %s palette", (theme, bg, fg, raised) => {
		const { field, button } = searchPanel(theme);
		expect(field.backgroundColor).toBe(bg);
		expect(field.color).toBe(fg);
		expect(button.backgroundColor).toBe(raised);
		expect(button.backgroundImage).toBe("none");
		expect(button.color).toBe(fg);
	});
});
