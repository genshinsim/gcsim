import { expect, test } from "vitest";
import { parser } from "./parser.gen";

function tokens(src: string): [string, string][] {
	const out: [string, string][] = [];
	parser.parse(src).iterate({
		enter: (node) => {
			if (node.name !== "Program") {
				out.push([src.slice(node.from, node.to), node.name]);
			}
		},
	});
	return out;
}

const kind = (src: string) => tokens(src)[0]?.[1];

test.each([
	["bennett", "CharacterName"],
	["ht", "CharacterName"],
	["skill", "ActionName"],
	["low_plunge", "ActionName"],
	["hp%", "StatName"],
	["atkspd%", "StatName"],
	["pyro", "ElementName"],
	["let", "Keyword"],
	["set", "Keyword"],
	["add", "Keyword"],
	["true", "Bool"],
	["myvar", "Identifier"],
	["0.466", "Number"],
	['"favoniussword"', "StringLiteral"],
	["# subs", "LineComment"],
	["// subs", "LineComment"],
])("%s is %s", (src, name) => {
	expect(tokens(src)).toEqual([[src, name]]);
});

test("stat beats element, as in the Go lexer", () => {
	expect(kind("pyro%")).toBe("StatName");
});

test("a character line", () => {
	expect(
		tokens('bennett add weapon="thealleyflash" refine=1; # main').map(
			([, name]) => name,
		),
	).toEqual([
		"CharacterName",
		"Keyword",
		"Keyword",
		"Operator",
		"StringLiteral",
		"Keyword",
		"Operator",
		"Number",
		"Punctuation",
		"LineComment",
	]);
});
