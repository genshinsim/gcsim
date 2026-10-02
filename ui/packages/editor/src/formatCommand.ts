import { indentUnit } from "@codemirror/language";
import type { ChangeSpec } from "@codemirror/state";
import type { EditorView } from "@codemirror/view";
import { formatGcsim, sameTokens } from "./language/format";

const isSpace = (c: string | undefined) =>
	c === " " || c === "\t" || c === "\n" || c === "\r";

// One change per whitespace gap rather than one whole-document replace, so
// CodeMirror maps the cursor onto the same token. Null if the texts differ in
// anything but whitespace.
function whitespaceChanges(src: string, out: string): ChangeSpec[] | null {
	const changes: ChangeSpec[] = [];
	let i = 0;
	let j = 0;
	for (;;) {
		const from = i;
		const start = j;
		while (isSpace(src[i])) i++;
		while (isSpace(out[j])) j++;
		const insert = out.slice(start, j);
		if (src.slice(from, i) !== insert) changes.push({ from, to: i, insert });
		if (i >= src.length || j >= out.length)
			return i >= src.length && j >= out.length ? changes : null;
		if (src[i++] !== out[j++]) return null;
	}
}

export function formatDocument(view: EditorView): boolean {
	if (view.state.readOnly) return false;
	const src = view.state.doc.toString();
	const out = formatGcsim(src, { indent: view.state.facet(indentUnit) });
	if (out === src || !sameTokens(src, out)) return false;
	const changes = whitespaceChanges(src, out);
	if (!changes) return false;
	view.dispatch({ changes, userEvent: "format" });
	return true;
}
