import { defaultEditorPrefs } from "@gcsim/components";
import { describe, expect, it } from "vitest";
import { loadEditorPrefs, saveEditorPrefs } from "./editorPrefs";

function fakeStorage(init: Record<string, string> = {}) {
	const data = new Map(Object.entries(init));
	return {
		data,
		getItem: (k: string) => data.get(k) ?? null,
		setItem: (k: string, v: string) => {
			data.set(k, v);
		},
	};
}

describe("editorPrefs", () => {
	it("falls back to defaults when nothing is stored", () => {
		expect(loadEditorPrefs(fakeStorage())).toEqual(defaultEditorPrefs);
	});

	it("loads prefs saved under the pre-existing keys and ignores the old tools key", () => {
		const storage = fakeStorage({
			"gcsim-config-editor-color-theme": "dracula",
			"gcsim-config-editor-font-size": "18",
			"gcsim-config-editor-tools": JSON.stringify({ tips: false }),
		});
		expect(loadEditorPrefs(storage)).toEqual({
			theme: "dracula",
			fontSize: 18,
		});
	});

	it("ignores unknown themes, including old Ace theme names", () => {
		const storage = fakeStorage({
			"gcsim-config-editor-theme": "monokai",
			"gcsim-config-editor-color-theme": "xcode",
		});
		expect(loadEditorPrefs(storage).theme).toBe(defaultEditorPrefs.theme);
	});

	it.each([
		["abc", defaultEditorPrefs.fontSize],
		["146", 28],
		["2", 10],
	])("keeps a stored font size of %s within bounds", (stored, expected) => {
		const storage = fakeStorage({ "gcsim-config-editor-font-size": stored });
		expect(loadEditorPrefs(storage).fontSize).toBe(expected);
	});

	it("round-trips through the same keys", () => {
		const storage = fakeStorage();
		const prefs = {
			theme: "solarized_light" as const,
			fontSize: 12,
		};
		saveEditorPrefs(storage, prefs);
		expect(storage.data.get("gcsim-config-editor-font-size")).toBe("12");
		expect(loadEditorPrefs(storage)).toEqual(prefs);
	});
});
