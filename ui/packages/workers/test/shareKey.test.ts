import { describe, expect, it } from "vitest";
import { newShareKey } from "../src/share/shareKey";

describe("newShareKey", () => {
	it("is 12 characters from the no-lookalikes alphabet", () => {
		expect(newShareKey()).toMatch(
			/^[6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz]{12}$/,
		);
	});
});
