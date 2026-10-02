import { indentUnit } from "@codemirror/language";
import type { ChangeSpec } from "@codemirror/state";
import type { EditorView } from "@codemirror/view";
import { formatGcsim, isSpace, sameTokens } from "./language/format";

function changesPerWhitespaceGap(src: string, out: string): ChangeSpec[] {
	const changes: ChangeSpec[] = [];
	let i = 0;
	let j = 0;
	while (i < src.length || j < out.length) {
		const from = i;
		const start = j;
		while (isSpace(src[i])) i++;
		while (isSpace(out[j])) j++;
		const insert = out.slice(start, j);
		if (src.slice(from, i) !== insert) changes.push({ from, to: i, insert });
		i++;
		j++;
	}
	return changes;
}

export function formatDocument(view: EditorView) {
	if (view.state.readOnly) return;
	const src = view.state.doc.toString();
	const out = formatGcsim(src, { indent: view.state.facet(indentUnit) });
	if (out === src || !sameTokens(src, out)) return;
	view.dispatch({
		changes: changesPerWhitespaceGap(src, out),
		userEvent: "format",
	});
}
