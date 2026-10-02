import { foldable } from "@codemirror/language";
import { EditorState } from "@codemirror/state";
import { expect, test } from "vitest";
import { gcsimFoldService } from "./fold";

const foldedText = (doc: string, lineNumber: number) => {
	const state = EditorState.create({ doc, extensions: [gcsimFoldService] });
	const line = state.doc.line(lineNumber);
	const range = foldable(state, line.from, line.to);
	return range && state.sliceDoc(range.from, range.to);
};

test("folds a block through its own closing brace, past nested blocks", () => {
	const doc = [
		"for let i = 0; i < 6; i = i + 1 {",
		"    if x {",
		"        a;",
		"    } else {",
		"        b;",
		"    }",
		"    c;",
		"}",
		"d;",
	].join("\n");
	expect(foldedText(doc, 1)).toBe(
		"\n    if x {\n        a;\n    } else {\n        b;\n    }\n    c;\n",
	);
});

test("folds a block with no nested blocks", () => {
	const doc = ["if x {", "    a;", "}", "while y {", "    b;", "}"].join("\n");
	expect(foldedText(doc, 1)).toBe("\n    a;\n");
});
