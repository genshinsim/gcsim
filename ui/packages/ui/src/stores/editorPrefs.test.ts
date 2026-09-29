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

	it("loads prefs saved under the pre-existing keys", () => {
		const storage = fakeStorage({
			"gcsim-config-editor-theme": "github",
			"gcsim-config-editor-font-size": "18",
			"gcsim-config-editor-tools": JSON.stringify({ tips: false }),
		});
		expect(loadEditorPrefs(storage)).toEqual({
			theme: "github",
			fontSize: 18,
			toggles: { team: true, nameSearch: true, tips: false },
		});
	});

	it("ignores unparseable toggles", () => {
		const storage = fakeStorage({ "gcsim-config-editor-tools": "{nope" });
		expect(loadEditorPrefs(storage).toggles).toEqual(
			defaultEditorPrefs.toggles,
		);
	});

	it("round-trips through the same keys", () => {
		const storage = fakeStorage();
		const prefs = {
			theme: "xcode",
			fontSize: 12,
			toggles: { team: false, nameSearch: true, tips: false },
		};
		saveEditorPrefs(storage, prefs);
		expect(storage.data.get("gcsim-config-editor-theme")).toBe("xcode");
		expect(loadEditorPrefs(storage)).toEqual(prefs);
	});
});
