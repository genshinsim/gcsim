import { indentUnit } from "@codemirror/language";
import type { ChangeSpec } from "@codemirror/state";
import type { EditorView } from "@codemirror/view";
import { formatGcsim, isSpace, sameTokens } from "./language/format";

function changedPart(gap: string, next: string, from: number): ChangeSpec[] {
	if (gap === next) return [];
	let head = 0;
	while (head < gap.length && gap[head] === next[head]) head++;
	let tail = 0;
	while (
		tail < gap.length - head &&
		tail < next.length - head &&
		gap[gap.length - 1 - tail] === next[next.length - 1 - tail]
	)
		tail++;
	return [
		{
			from: from + head,
			to: from + gap.length - tail,
			insert: next.slice(head, next.length - tail),
		},
	];
}

function changesPerWhitespaceGap(src: string, out: string): ChangeSpec[] {
	const changes: ChangeSpec[] = [];
	let i = 0;
	let j = 0;
	while (i < src.length || j < out.length) {
		const from = i;
		const start = j;
		while (isSpace(src[i])) i++;
		while (isSpace(out[j])) j++;
		changes.push(...changedPart(src.slice(from, i), out.slice(start, j), from));
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
