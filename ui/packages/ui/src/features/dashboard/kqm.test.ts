import { model } from "@gcsim/types";
import { describe, expect, it } from "vitest";
import { parseKqmEntries } from "./kqm";

describe("parseKqmEntries", () => {
	it("reads the fields the dashboard shows", () => {
		const json = {
			data: [
				{
					_id: "abc",
					description: "Hu Tao vape",
					submitter: "someone",
					config: "ignored",
					summary: {
						team: [{ name: "hutao" }, { name: "xingqiu" }],
						mean_dps_per_target: 41234.5,
						mode: "TTK_MODE",
					},
				},
			],
		};
		expect(parseKqmEntries(json)).toEqual([
			{
				_id: "abc",
				description: "Hu Tao vape",
				submitter: "someone",
				summary: {
					team: [{ name: "hutao" }, { name: "xingqiu" }],
					mean_dps_per_target: 41234.5,
					mode: model.SimMode.TTK_MODE,
				},
			},
		]);
	});

	it("reads a numeric mode", () => {
		const [entry] = parseKqmEntries({
			data: [{ _id: "a", summary: { mode: 1 } }],
		});
		expect(entry.summary?.mode).toBe(model.SimMode.DURATION_MODE);
	});

	it("falls back to id when _id is missing", () => {
		expect(parseKqmEntries({ data: [{ id: "b" }] })).toEqual([{ _id: "b" }]);
	});

	it("returns no entries when data is missing or malformed", () => {
		expect(parseKqmEntries({})).toEqual([]);
		expect(parseKqmEntries(null)).toEqual([]);
		expect(parseKqmEntries({ data: "nope" })).toEqual([]);
	});

	it("skips entries that are not objects", () => {
		expect(parseKqmEntries({ data: [null, 3, { _id: "a" }] })).toEqual([
			{ _id: "a" },
		]);
	});
});
