import { foldService } from "@codemirror/language";

export const gcsimFoldService = foldService.of((state, lineStart) => {
	const line = state.doc.lineAt(lineStart);
	if (!line.text.trimEnd().endsWith("{")) return null;

	const doc = state.doc.sliceString(line.to);
	let depth = 0;
	for (let i = 0; i < doc.length; i++) {
		if (doc[i] === "{") depth++;
		if (doc[i] === "}" && --depth === 0) {
			return { from: line.to, to: line.to + i };
		}
	}
	return null;
});
