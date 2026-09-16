import { describe, expect, it } from "vitest";
import { createKeyOrLabelPredicate } from "./utils";

describe("createKeyOrLabelPredicate", () => {
	const predicate = createKeyOrLabelPredicate((item: string) =>
		item === "xiangling" ? "Xiangling" : "Amber",
	);

	it("matches everything for an empty query", () => {
		expect(predicate("xiangling", "")).toBe(true);
		expect(predicate("amber", "   ")).toBe(true);
	});

	it("matches by raw key, case-insensitively", () => {
		expect(predicate("xiangling", "XIANG")).toBe(true);
		expect(predicate("amber", "xiang")).toBe(false);
	});

	it("matches by translated label, case-insensitively", () => {
		expect(predicate("xiangling", "ling")).toBe(true);
	});
});
