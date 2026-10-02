import { expect, test } from "vitest";
import { formatGcsim, sameTokens } from "./format";

const lines = (...l: string[]) => `${l.join("\n")}\n`;

test("indents block bodies one tab and aligns } with the opening keyword", () => {
	expect(
		formatGcsim("while true\n{\nhutao attack;\n    if x { hutao skill; }\n}"),
	).toBe(
		lines(
			"while true {",
			"\thutao attack;",
			"\tif x {",
			"\t\thutao skill;",
			"\t}",
			"}",
		),
	);
});

test("keeps else on the closing brace line", () => {
	expect(formatGcsim("if x {\nhutao skill;\n}\nelse {\nhutao burst;\n}")).toBe(
		lines("if x {", "\thutao skill;", "} else {", "\thutao burst;", "}"),
	);
});

test("puts each case body on its own line, one level under the case", () => {
	expect(
		formatGcsim("switch n { case 1: hutao skill; default: hutao attack; }"),
	).toBe(
		lines(
			"switch n {",
			"\tcase 1:",
			"\t\thutao skill;",
			"\tdefault:",
			"\t\thutao attack;",
			"}",
		),
	);
});

test("lines a comment above a case label up with the label", () => {
	expect(
		formatGcsim(
			"switch n {\ncase 1:\nhutao skill;\n# fallback\ndefault:\nhutao attack;\n}",
		),
	).toBe(
		lines(
			"switch n {",
			"\tcase 1:",
			"\t\thutao skill;",
			"\t# fallback",
			"\tdefault:",
			"\t\thutao attack;",
			"}",
		),
	);
});

test("keeps a for header on one line", () => {
	expect(formatGcsim("for let i=0;i<3;i=i+1 {\nhutao attack;\n}")).toBe(
		lines("for let i = 0; i < 3; i = i + 1 {", "\thutao attack;", "}"),
	);
});

test("prints character, option and target lines as tight key=value pairs", () => {
	expect(
		formatGcsim(
			[
				"hutao char lvl = 90/90 cons = 1 talent = 9, 9, 9;",
				'hutao add weapon = "homa" refine = 1 +params = [stacks = 2];',
				"target lvl = 100 pos = 0, 2.4 ;",
			].join("\n"),
		),
	).toBe(
		lines(
			"hutao char lvl=90/90 cons=1 talent=9,9,9;",
			'hutao add weapon="homa" refine=1 +params=[stacks=2];',
			"target lvl=100 pos=0,2.4;",
		),
	);
});

test("spaces script operators but keeps calls, action params and repeats tight", () => {
	expect(
		formatGcsim(
			"let x=1+-2*f( y );\nhutao skill [ hold = 1 ] ;\nhutao attack : 3;\nlet n = i-1;",
		),
	).toBe(
		lines(
			"let x = 1 + -2 * f(y);",
			"hutao skill[hold=1];",
			"hutao attack:3;",
			"let n = i-1;",
		),
	);
});

test("keeps comments, with one space before a trailing comment", () => {
	expect(
		formatGcsim(
			"# team\nhutao attack;    # go\nwhile true {\n// loop\nhutao skill;\n}",
		),
	).toBe(
		lines(
			"# team",
			"hutao attack; # go",
			"while true {",
			"\t// loop",
			"\thutao skill;",
			"}",
		),
	);
});

test("keeps at most one blank line, and none just inside braces", () => {
	expect(
		formatGcsim("active hutao;\n\n\n\nwhile true {\n\nhutao attack;\n\n}\n"),
	).toBe(lines("active hutao;", "", "while true {", "\thutao attack;", "}"));
});

test("keeps the author's line breaks inside a statement, indented one level", () => {
	expect(formatGcsim("hutao add stats hp=4780\natk=311 em=187;")).toBe(
		lines("hutao add stats hp=4780", "\tatk=311 em=187;"),
	);
});

test("indents with the given unit", () => {
	expect(
		formatGcsim("while true {\nhutao attack;\n}", { indent: "    " }),
	).toBe(lines("while true {", "    hutao attack;", "}"));
});

test("only fixes stray spacing in an already tidy config", () => {
	const config = lines(
		"sucrose char lvl=90/90 cons=6 talent=9,9,9;",
		'sucrose add weapon="sacfragments" refine=3 lvl=90/90;',
		'sucrose add set="viridescentvenerer" count=4;',
		"sucrose add stats hp=4780 atk=311 em=187 em=187 em=187 ; #main",
		"",
		"options swap_delay=12 iteration=1;",
		"",
		"active sucrose;",
		"",
		"target lvl=100 resist=0.1 radius=2 pos=0,2.4 hp=999999999;",
		"",
		"sucrose burst, skill, attack;",
	);
	expect(formatGcsim(config)).toBe(
		config.replace("em=187 ; #main", "em=187; #main"),
	);
});

test("leaves the tokens of a half-typed config alone", () => {
	const config =
		'while true {\nif .x > 1 {\nhutao attack;\nlet s = "unterminated\nhutao skill & burst;\n}';
	const out = formatGcsim(config);
	expect(sameTokens(config, out)).toBe(true);
	expect(formatGcsim(out)).toBe(out);
});

test("sameTokens ignores whitespace between tokens but not inside them", () => {
	expect(sameTokens("x=1", "x = 1")).toBe(true);
	expect(sameTokens("a b", "ab")).toBe(false);
	expect(sameTokens('"a b"', '"a  b"')).toBe(false);
});
