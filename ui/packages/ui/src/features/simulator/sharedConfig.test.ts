import Pako from "pako";
import { describe, expect, it } from "vitest";
import { decodeSharedConfig, encodeSharedConfig } from "./sharedConfig";

const config = `raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="engulfinglightning" refine=1 lvl=90/90;
active raiden;
raiden burst;
`;

function toBase64(bytes: Uint8Array): string {
	return btoa(String.fromCharCode(...bytes));
}

describe("decodeSharedConfig", () => {
	it("decodes gzipped base64 json with title and source", () => {
		const encoded = toBase64(
			Pako.gzip(JSON.stringify({ config, title: "Raiden", source: "calc" })),
		);
		expect(decodeSharedConfig(encoded)).toEqual({
			config,
			title: "Raiden",
			source: "calc",
		});
	});

	it("accepts zlib instead of gzip", () => {
		const encoded = toBase64(Pako.deflate(JSON.stringify({ config })));
		expect(decodeSharedConfig(encoded)).toEqual({ config });
	});

	it("accepts url-safe base64 without padding", () => {
		const encoded = toBase64(
			Pako.gzip(JSON.stringify({ config, title: "?>~" })),
		)
			.replaceAll("+", "-")
			.replaceAll("/", "_")
			.replace(/=+$/, "");
		expect(decodeSharedConfig(encoded)).toEqual({ config, title: "?>~" });
	});

	it("treats spaces as '+' left over from query-string decoding", () => {
		const encoded = encodeSharedConfig({ config: config.repeat(20) });
		expect(encoded).toContain("+");
		expect(decodeSharedConfig(encoded.replaceAll("+", " "))).toEqual({
			config: config.repeat(20),
		});
	});

	it("drops title and source that are not non-empty strings", () => {
		const encoded = toBase64(
			Pako.gzip(JSON.stringify({ config, title: 5, source: "" })),
		);
		expect(decodeSharedConfig(encoded)).toEqual({ config });
	});

	it.each([
		["not base64", "!!!"],
		["not compressed", btoa("hello")],
		["not json", toBase64(Pako.gzip("hello"))],
		["missing config", toBase64(Pako.gzip(JSON.stringify({ title: "x" })))],
		["empty config", toBase64(Pako.gzip(JSON.stringify({ config: "" })))],
		["non-object json", toBase64(Pako.gzip(JSON.stringify("cfg")))],
	])("throws when %s", (_, encoded) => {
		expect(() => decodeSharedConfig(encoded)).toThrow();
	});
});

describe("encodeSharedConfig", () => {
	it("round-trips through decodeSharedConfig", () => {
		const payload = { config, title: "Raiden", source: "calc" };
		expect(decodeSharedConfig(encodeSharedConfig(payload))).toEqual(payload);
	});
});
