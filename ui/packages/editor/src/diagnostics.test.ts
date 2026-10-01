import { Text } from "@codemirror/state";
import { expect, test } from "vitest";
import { parseErrors, toDiagnostics } from "./diagnostics";

test("parses line and line:col prefixes, drops unpositioned errors", () => {
	expect(
		parseErrors(
			"ln3:5: unexpected token\n\tconfig does not contain any targets\nln7: bad action",
		),
	).toEqual([
		{ line: 3, column: 5, message: "unexpected token" },
		{ line: 7, column: undefined, message: "bad action" },
	]);
});

test("empty input", () => {
	expect(parseErrors(null)).toEqual([]);
	expect(parseErrors("")).toEqual([]);
});

const doc = Text.of(["first", "second line", "third"]);

test("marks from the column to the end of the line", () => {
	const [d] = toDiagnostics(doc, [{ line: 2, column: 8, message: "x" }]);
	expect(doc.sliceString(d.from, d.to)).toBe("line");
	expect(d.severity).toBe("error");
});

test("marks the whole line without a column", () => {
	const [d] = toDiagnostics(doc, [{ line: 1, message: "x" }]);
	expect(doc.sliceString(d.from, d.to)).toBe("first");
});

test("clamps positions past the document", () => {
	const [d] = toDiagnostics(doc, [{ line: 99, column: 99, message: "x" }]);
	expect(d.from).toBe(doc.length);
	expect(d.to).toBe(doc.length);
});
