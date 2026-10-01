import { readFileSync } from "node:fs";
import { buildParserFile } from "@lezer/generator";
import { expect, test } from "vitest";

const read = (name: string) =>
	readFileSync(new URL(name, import.meta.url), "utf8");

test("committed parser matches gcsim.grammar (run `pnpm build:grammar`)", () => {
	const { parser, terms } = buildParserFile(read("./gcsim.grammar"), {
		fileName: "gcsim.grammar",
		typeScript: true,
	});
	expect(read("./parser.gen.ts")).toBe(parser);
	expect(read("./parser.gen.terms.ts")).toBe(terms);
});
