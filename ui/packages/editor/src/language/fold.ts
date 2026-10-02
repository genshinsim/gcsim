import { foldService } from "@codemirror/language";

export const gcsimFoldService = foldService.of((state, lineStart) => {
	const line = state.doc.lineAt(lineStart);
	if (!line.text.trimEnd().endsWith("{")) return null;

	const rest = state.doc.sliceString(line.to);
	let depth = 1;
	for (let i = 0; i < rest.length; i++) {
		if (rest[i] === "{") depth++;
		if (rest[i] === "}" && --depth === 0) {
			return { from: line.to, to: line.to + i };
		}
	}
	return null;
});
