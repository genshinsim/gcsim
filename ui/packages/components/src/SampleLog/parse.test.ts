import { describe, expect, it } from "vitest";
import { eventColor, strFrameWithSec } from "./parse";

describe("strFrameWithSec", () => {
	it("shows the frame and its seconds at 60 fps", () => {
		expect(strFrameWithSec(123)).toBe(" [123 | 2.05s]");
	});

	it("marks a permanent (-1) frame without seconds", () => {
		expect(strFrameWithSec(-1)).toBe(" [-1]");
	});
});

describe("eventColor", () => {
	it("gives known events their colour", () => {
		expect(eventColor("damage")).toBe("#2563EB");
	});

	it("leaves uncoloured known events blank", () => {
		expect(eventColor("heal")).toBe("");
	});

	it("falls back for unknown events", () => {
		expect(eventColor("nope")).toBe("gray-500");
	});
});
