import { type Diagnostic, setDiagnostics } from "@codemirror/lint";
import type { EditorState, Text } from "@codemirror/state";

export interface GcsimError {
	line: number;
	column?: number;
	message: string;
}

const LINE_PREFIX = /^\s*ln(\d+)(?::(\d+))?:\s*(.*)$/;

export function parseErrorLine(text: string): GcsimError | undefined {
	const m = LINE_PREFIX.exec(text);
	if (!m) return undefined;
	return {
		line: Number(m[1]),
		column: m[2] ? Number(m[2]) : undefined,
		message: m[3],
	};
}

export function parseErrors(message: string | null | undefined): GcsimError[] {
	if (!message) return [];
	return message
		.split("\n")
		.map(parseErrorLine)
		.filter((err) => err !== undefined);
}

export function errorRange(doc: Text, err: GcsimError) {
	const line = doc.line(Math.min(Math.max(err.line, 1), doc.lines));
	const offset = err.column ? err.column - 1 : 0;
	return { from: Math.min(line.from + offset, line.to), to: line.to };
}

export function toDiagnostics(doc: Text, errors: GcsimError[]): Diagnostic[] {
	return errors.map((err) => ({
		...errorRange(doc, err),
		severity: "error",
		message: err.message,
		source: "gcsim",
	}));
}

export function diagnosticsTransaction(
	state: EditorState,
	errors: GcsimError[],
) {
	return setDiagnostics(state, toDiagnostics(state.doc, errors));
}
