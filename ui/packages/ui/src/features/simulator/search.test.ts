import { describe, expect, it } from "vitest";
import { validateSimulatorSearch } from "./search";

describe("validateSimulatorSearch", () => {
	it("keeps a cfg string", () => {
		expect(validateSimulatorSearch({ cfg: "H4sI" })).toEqual({ cfg: "H4sI" });
	});

	it.each([
		["empty", ""],
		["non-string", 42],
	])("drops a %s cfg", (_, cfg) => {
		expect(validateSimulatorSearch({ cfg })).toEqual({});
	});
});
