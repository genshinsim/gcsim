import { foldable } from "@codemirror/language";
import { EditorState } from "@codemirror/state";
import { expect, test } from "vitest";
import { gcsimFoldService } from "./fold";

const lines = (...l: string[]) => l.join("\n");

const foldedText = (doc: string, lineNumber: number) => {
	const state = EditorState.create({ doc, extensions: [gcsimFoldService] });
	const line = state.doc.line(lineNumber);
	const range = foldable(state, line.from, line.to);
	return range && state.sliceDoc(range.from, range.to);
};

const loop = lines(
	"for let i = 0; i < 6; i = i + 1 {",
	"\tif x {",
	"\t\ta;",
	"\t} else {",
	"\t\tb;",
	"\t}",
	"\tc;",
	"\tif y {",
	"\t\td;",
	"\t}",
	"}",
	"e;",
);

test("folds a block through its own closing brace, past nested blocks", () => {
	expect(foldedText(loop, 1)).toBe(
		lines(
			"",
			"\tif x {",
			"\t\ta;",
			"\t} else {",
			"\t\tb;",
			"\t}",
			"\tc;",
			"\tif y {",
			"\t\td;",
			"\t}",
			"",
		),
	);
});

test("folds nested and else blocks to their own closing brace", () => {
	expect(foldedText(loop, 2)).toBe(lines("", "\t\ta;", "\t"));
	expect(foldedText(loop, 4)).toBe(lines("", "\t\tb;", "\t"));
});

test("folds a block with no nested blocks", () => {
	const doc = lines("if x {", "\ta;", "}", "while y {", "\tb;", "}");
	expect(foldedText(doc, 1)).toBe(lines("", "\ta;", ""));
});
