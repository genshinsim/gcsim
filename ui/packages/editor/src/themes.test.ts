import { describe, expect, it } from "vitest";
import {
	DEFAULT_EDITOR_THEME,
	EDITOR_THEMES,
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
