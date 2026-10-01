import { CompletionContext } from "@codemirror/autocomplete";
import { EditorState } from "@codemirror/state";
import { expect, test } from "vitest";
import { gcsimCompletionSource } from "./autocomplete";
import { gcsim } from "./language";

function complete(src: string, explicit = false) {
	const pos = src.indexOf("|");
	const doc = src.slice(0, pos) + src.slice(pos + 1);
	const state = EditorState.create({ doc, extensions: [gcsim()] });
	const result = gcsimCompletionSource(
		new CompletionContext(state, pos, explicit),
	);
	return (
		result && { from: result.from, labels: result.options.map((o) => o.label) }
	);
}

test("statement start offers characters and keywords", () => {
	const r = complete("ben|");
	expect(r?.from).toBe(0);
	expect(r?.labels).toContain("bennett");
	expect(r?.labels).toContain("active");
	expect(r?.labels).not.toContain("hp%");
});

test("after a character offers char/add and actions", () => {
	const r = complete("bennett |", true);
	expect(r?.labels).toEqual(expect.arrayContaining(["char", "add", "skill"]));
	expect(r?.labels).not.toContain("bennett");
});

test("after `add` offers what can be added", () => {
	expect(complete("bennett add |", true)?.labels).toEqual([
		"weapon",
		"set",
		"stats",
		"params",
	]);
});

test("in an add stats line offers stats", () => {
	const r = complete("xq add stats hp=4780 at|");
	expect(r?.labels).toContain("atk%");
	expect(r?.labels).not.toContain("skill");
});

test('weapon names inside weapon=""', () => {
	const r = complete('bennett add weapon="favon|');
	expect(r?.from).toBe('bennett add weapon="'.length);
	expect(r?.labels).toContain("favoniussword");
});

test('set names inside set=""', () => {
	expect(complete('bennett add set="nob|')?.labels).toContain("noblesseoblige");
});

test("an action list keeps offering actions", () => {
	expect(complete("xiangling attack, bu|")?.labels).toContain("burst");
});

test("active offers characters", () => {
	expect(complete("active zh|")?.labels).toContain("zhongli");
});

test("statements are split on ;", () => {
	expect(complete("active xq; ben|")?.labels).toContain("bennett");
});

test("a ; inside a comment does not end the statement", () => {
	expect(complete("bennett # a; b\nadd |", true)?.labels).toEqual(
		complete("bennett add |", true)?.labels,
	);
});

test("shortcuts show their key", () => {
	const state = EditorState.create({ doc: "x", extensions: [gcsim()] });
	const r = gcsimCompletionSource(new CompletionContext(state, 1, true));
	expect(r?.options.find((o) => o.label === "ht")?.detail).toBe("hutao");
});

test("nothing inside comments or other strings", () => {
	expect(complete("# ben|")).toBeNull();
	expect(complete('let x = "ben|')).toBeNull();
});

test("nothing without a word unless asked", () => {
	expect(complete("|")).toBeNull();
	expect(complete("|", true)?.labels.length).toBeGreaterThan(0);
});
