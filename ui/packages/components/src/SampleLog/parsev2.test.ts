import { describe, expect, it, vi } from "vitest";
import { parseLogV2 } from "./parsev2";
import { AllSampleOptions, DefaultSampleOptions } from "./SampleOptions";
import { sampleFixture as sample } from "./testdata";

const team = sample.character_details?.map((c) => c.name) ?? [];

function parse(selected: string[]) {
	return parseLogV2(sample.initial_character, team, sample.logs, selected);
}

function msgs(selected: string[]) {
	return parse(selected).flatMap((r) => r.slots.flat().map((e) => e.msg));
}

describe("parseLogV2", () => {
	it("emits one row per frame, in frame order, with a sim slot plus one slot per character", () => {
		const rows = parse(AllSampleOptions);
		expect(rows.length).toBeGreaterThan(0);
		for (const row of rows) {
			expect(row.slots).toHaveLength(team.length + 1);
		}
		const frames = rows.map((r) => r.f);
		expect(frames).toEqual([...new Set(frames)].sort((a, b) => a - b));
	});

	it("keeps only the selected events", () => {
		const events = new Set(
			parse(["damage"]).flatMap((r) => r.slots.flat().map((e) => e.event)),
		);
		expect([...events]).toEqual(["damage"]);
	});

	it("returns no rows when nothing is selected", () => {
		expect(parse([])).toEqual([]);
	});

	it("tracks the active character across swaps", () => {
		const rows = parse(["action"]);
		const activeAt = (f: number) => rows.find((r) => r.f === f)?.active;
		expect(activeAt(1)).toBe(team.indexOf("yaemiko") + 1);
		expect(activeAt(107)).toBe(team.indexOf("bennett") + 1);
		expect(activeAt(221)).toBe(team.indexOf("chevreuse") + 1);
	});

	it("formats damage with a grouped amount and crit tag in the character's slot", () => {
		const row = parse(["damage"]).find((r) => r.f === 122);
		const yae = team.indexOf("yaemiko") + 1;
		expect(row?.slots[yae].map((e) => e.msg)).toContain(
			"Sesshou Sakura Tick [8,815] (crit)",
		);
	});

	it("rewrites swaps and strips 'executed' from actions", () => {
		const out = msgs(["action"]);
		expect(out).toContain("swap to bennett");
		expect(out).toContain("burst");
	});

	it("adds status expiry rows only when status is selected", () => {
		const expired = "tf-4pc-icd expired [123 | 2.05s]";
		const row = parse(DefaultSampleOptions).find((r) => r.f === 180);
		expect(row?.slots[team.indexOf("bennett") + 1].map((e) => e.msg)).toContain(
			expired,
		);
		expect(msgs(["damage"])).not.toContain(expired);
	});

	it("drops ordering from the event details", () => {
		const item = parse(["damage"])[0].slots.flat()[0];
		expect(item.data).toBeDefined();
		expect(item.data).not.toHaveProperty("ordering");
		expect(JSON.parse(item.raw)).not.toHaveProperty("ordering");
	});

	it("returns no rows for an unparseable log", () => {
		vi.spyOn(console, "warn").mockImplementation(() => {});
		expect(parseLogV2("a", ["a"], "not json", AllSampleOptions)).toEqual([]);
	});
});
