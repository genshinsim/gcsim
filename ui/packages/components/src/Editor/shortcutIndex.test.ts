import { describe, expect, it } from "vitest";
import { buildShortcutEntries, searchShortcuts } from "./shortcutIndex";

const entries = buildShortcutEntries();
const keys = (query: string) =>
	searchShortcuts(entries, query).map((m) => `${m.entry.kind}:${m.entry.key}`);

describe("searchShortcuts", () => {
	it("lists everything for an empty query", () => {
		expect(searchShortcuts(entries, "  ")).toHaveLength(entries.length);
	});

	it("finds a key by its shortcut and reports the alias", () => {
		const [first] = searchShortcuts(entries, "HT");
		expect(first.entry).toMatchObject({ kind: "character", key: "hutao" });
		expect(first.alias).toBe("ht");
	});

	it("does not report an alias when the key itself matches", () => {
		const [first] = searchShortcuts(entries, "hutao");
		expect(first.entry.key).toBe("hutao");
		expect(first.alias).toBeUndefined();
	});

	it("searches every kind", () => {
		expect(keys("burst")).toContain("action:burst");
		expect(keys("em")).toContain("stat:em");
		expect(keys("favoniuswarbow")).toContain("weapon:favoniuswarbow");
	});

	it("puts the kind with the best match first", () => {
		expect(keys("burst")[0]).toBe("action:burst");
	});
});
