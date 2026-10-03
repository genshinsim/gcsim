import type { LogDetails, Sample } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { sampleFixture } from "../testdata";
import { fromLegacySample } from "./legacy";
import { isEvent, type SimEvent } from "./types";

const events = fromLegacySample(sampleFixture);

const find = (pred: (e: SimEvent) => boolean) => {
	const e = events.find(pred);
	if (e == null) {
		throw new Error("no matching event");
	}
	return e;
};

function line(over: Partial<LogDetails>): LogDetails {
	return {
		char_index: -1,
		ended: 0,
		event: "sim",
		frame: 0,
		msg: "",
		logs: {},
		...over,
	};
}

function sampleOf(logs: LogDetails[]): Sample {
	return {
		initial_character: "bennett",
		character_details: [{ name: "bennett" }, { name: "xiangling" }],
		logs,
	} as unknown as Sample;
}

describe("fromLegacySample", () => {
	it("returns events in frame order", () => {
		const frames = events.map((e) => e.frame);
		expect(frames).toEqual([...frames].sort((a, b) => a - b));
	});

	it("keeps the log category as the type and char_index as the character", () => {
		const e = find((e) => e.type === "damage");
		expect(e.frame).toBe(122);
		expect(e.characterIndex).toBe(3);
	});

	it("gives raw the line with its logs in ordering order and no ordering", () => {
		const out = fromLegacySample(
			sampleOf([
				line({
					event: "sim",
					msg: "x",
					logs: { b: 2, a: 1 },
					ordering: { a: 0, b: 1 },
				}),
			]),
		);
		const raw = out.find((e) => e.type === "sim")?.raw as LogDetails;
		expect(Object.keys(raw.logs)).toEqual(["a", "b"]);
		expect(raw).not.toHaveProperty("ordering");
	});

	describe("end", () => {
		it("is the end frame when the line ended after it started", () => {
			const added = find(
				(e) => isEvent(e, "status") && e.key === "bennettburst",
			);
			expect(added.frame).toBe(158);
			expect(added.end).toBe(912);
		});

		it("covers any category, not just statuses", () => {
			const anim = find((e) => e.type === "hitlag" && e.frame === 1);
			expect(anim.end).toBe(38);
		});

		it("is Infinity for a line that never ends", () => {
			const stam = find(
				(e) => isEvent(e, "status") && e.key === "utility-dash",
			);
			expect(stam.end).toBe(Number.POSITIVE_INFINITY);
		});

		it("is absent for an instant", () => {
			const skill = find((e) => e.type === "action" && e.frame === 1);
			expect(skill).not.toHaveProperty("end");
		});
	});

	describe("stint", () => {
		const field = events.filter((e) => isEvent(e, "stint"));

		it("follows the initial character then each executed swap", () => {
			expect(field.map((e) => [e.characterIndex, e.frame, e.end])).toEqual([
				[3, 0, 107],
				[0, 107, 221],
				[1, 221, 299],
			]);
		});

		it("names the character", () => {
			expect(field.map((e) => e.message)).toEqual([
				"yaemiko",
				"bennett",
				"chevreuse",
			]);
		});

		it("is one stint for the whole sample without swaps", () => {
			const out = fromLegacySample(
				sampleOf([line({ frame: 0 }), line({ frame: 50 })]),
			);
			expect(
				out
					.filter((e) => isEvent(e, "stint"))
					.map((e) => [e.characterIndex, e.frame, e.end]),
			).toEqual([[0, 0, 50]]);
		});
		it("skips a character outside the team", () => {
			const swap = line({
				event: "action",
				msg: "executed swap",
				char_index: 0,
				frame: 10,
				logs: { action: "swap" },
			});
			const out = fromLegacySample({
				...sampleOf([swap, line({ frame: 50 })]),
				character_details: undefined,
			});
			expect(out.filter((e) => isEvent(e, "stint"))).toEqual([]);
		});
	});

	describe("typed fields", () => {
		it("gives damage its amount", () => {
			const e = find((e) => e.type === "damage");
			expect(isEvent(e, "damage") && e.damage).toBeCloseTo(8815.17, 1);
		});

		it("gives actions their kind", () => {
			const bursts = events.filter(
				(e) => isEvent(e, "action") && e.action === "burst",
			);
			expect(bursts.map((e) => [e.frame, e.characterIndex])).toEqual([
				[158, 0],
				[221, 1],
			]);
		});

		it("gives an action without a kind an empty one", () => {
			const e = find(
				(e) => e.type === "action" && e.message.startsWith("swapping"),
			);
			expect(isEvent(e, "action") && e.action).toBe("");
		});

		it("gives statuses their key", () => {
			const e = find((e) => e.type === "status" && e.frame === 35);
			expect(isEvent(e, "status") && e.key).toBe("yae_oldest_totem_expiry");
		});

		it("fills every typed field", () => {
			for (const e of events) {
				if (isEvent(e, "action")) {
					expect(typeof e.action).toBe("string");
				}
				if (isEvent(e, "damage")) {
					expect(Number.isFinite(e.damage)).toBe(true);
				}
				if (isEvent(e, "status")) {
					expect(typeof e.key).toBe("string");
				}
				if (isEvent(e, "stint")) {
					expect(Number.isFinite(e.end)).toBe(true);
				}
			}
		});
	});

	describe("message", () => {
		it("adds the amount and modifiers to damage", () => {
			const out = fromLegacySample(
				sampleOf([
					line({
						event: "damage",
						msg: "Pyronado",
						logs: { damage: 12345.6, amp: "vaporize", crit: true },
					}),
				]),
			);
			expect(out.find((e) => e.type === "damage")?.message).toBe(
				"Pyronado [12,346] (vaporize crit)",
			);
		});

		it("drops 'executed' from actions and names the swap target", () => {
			const msgs = events
				.filter((e) => e.type === "action")
				.map((e) => e.message);
			expect(msgs).toContain("skill[hold=1]");
			expect(msgs).toContain("swap to bennett");
		});

		it("prefixes statuses with their key", () => {
			const e = find((e) => isEvent(e, "status") && e.key === "bennettburst");
			expect(e.message).toBe("bennettburst status added");
		});

		it("prefixes unknown categories with the category", () => {
			const out = fromLegacySample(
				sampleOf([line({ event: "queue", msg: "something" })]),
			);
			expect(out.find((e) => e.type === "queue")?.message).toBe(
				"queue: something",
			);
		});
	});

	it("accepts the log as a JSON string", () => {
		const s = {
			...sampleFixture,
			logs: JSON.stringify(sampleFixture.logs),
		} as unknown as Sample;
		expect(fromLegacySample(s)).toEqual(events);
	});

	it("returns nothing for a missing or broken log", () => {
		expect(fromLegacySample({ ...sampleFixture, logs: undefined })).toEqual([]);
		expect(
			fromLegacySample({ ...sampleFixture, logs: "{" } as unknown as Sample),
		).toEqual([]);
		expect(
			fromLegacySample({ ...sampleFixture, logs: "null" } as unknown as Sample),
		).toEqual([]);
	});
});
