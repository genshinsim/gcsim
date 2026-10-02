import {
	autocompletion,
	closeBrackets,
	closeBracketsKeymap,
	completionKeymap,
} from "@codemirror/autocomplete";
import {
	defaultKeymap,
	history,
	historyKeymap,
	indentWithTab,
} from "@codemirror/commands";
import {
	bracketMatching,
	foldGutter,
	foldKeymap,
	indentOnInput,
	indentUnit,
} from "@codemirror/language";
import { lintGutter } from "@codemirror/lint";
import { highlightSelectionMatches, searchKeymap } from "@codemirror/search";
import {
	Annotation,
	Compartment,
	EditorState,
	type Extension,
} from "@codemirror/state";
import {
	drawSelection,
	EditorView,
	highlightActiveLineGutter,
	keymap,
	lineNumbers,
} from "@codemirror/view";
import { type Ref, useEffect, useImperativeHandle, useRef } from "react";
import { diagnosticsTransaction, parseErrors } from "./diagnostics";
import { errorPanelField, setErrorPanel } from "./errorPanel";
import { formatDocument } from "./formatCommand";
import { gcsimFoldService } from "./language/fold";
import { gcsim } from "./language/language";
import {
	DEFAULT_EDITOR_THEME,
	type EditorThemeId,
	editorTheme,
} from "./themes";

export interface EditorHandle {
	format: () => void;
}

export interface EditorProps {
	ref?: Ref<EditorHandle>;
	value: string;
	onChange?: (value: string) => void;
	readOnly?: boolean;
	fontSize?: number;
	theme?: EditorThemeId;
	maxLines?: number;
	error?: string | null;
	errorTitle?: string;
	id?: string;
	className?: string;
}

const LINE_HEIGHT = 1.5;

function sizing(fontSize: number, maxLines: number): Extension {
	return EditorView.theme({
		"&": { fontSize: `${fontSize}px` },
		".cm-scroller": {
			maxHeight: Number.isFinite(maxLines)
				? `${maxLines * LINE_HEIGHT}em`
				: "none",
		},
	});
}

const fromProps = Annotation.define<boolean>();

function changeBetween(current: string, next: string) {
	let from = 0;
	const max = Math.min(current.length, next.length);
	while (from < max && current[from] === next[from]) from++;
	let end = 0;
	while (
		end < max - from &&
		current[current.length - 1 - end] === next[next.length - 1 - end]
	)
		end++;
	return {
		from,
		to: current.length - end,
		insert: next.slice(from, next.length - end),
	};
}

function editable(readOnly: boolean): Extension {
	return [EditorState.readOnly.of(readOnly), EditorView.editable.of(!readOnly)];
}

export function Editor({
	ref,
	value,
	onChange,
	readOnly = false,
	fontSize = 14,
	theme = DEFAULT_EDITOR_THEME,
	maxLines = 35,
	error,
	errorTitle = "Invalid config",
	id,
	className,
}: EditorProps) {
	const parentRef = useRef<HTMLDivElement>(null);
	const viewRef = useRef<EditorView | null>(null);
	const onChangeRef = useRef(onChange);
	onChangeRef.current = onChange;
	const compartments = useRef({
		sizing: new Compartment(),
		editable: new Compartment(),
		theme: new Compartment(),
	}).current;
	const initial = useRef({ value, theme, fontSize, maxLines, readOnly });

	useImperativeHandle(
		ref,
		() => ({
			format: () => {
				if (viewRef.current) formatDocument(viewRef.current);
			},
		}),
		[],
	);

	useEffect(() => {
		const { value, theme, fontSize, maxLines, readOnly } = initial.current;
		const view = new EditorView({
			parent: parentRef.current ?? undefined,
			state: EditorState.create({
				doc: value,
				extensions: [
					lineNumbers(),
					highlightActiveLineGutter(),
					foldGutter(),
					gcsimFoldService,
					lintGutter(),
					history(),
					drawSelection(),
					indentOnInput(),
					indentUnit.of("\t"),
					EditorState.tabSize.of(4),
					bracketMatching(),
					closeBrackets(),
					highlightSelectionMatches(),
					autocompletion(),
					gcsim(),
					errorPanelField,
					compartments.theme.of(editorTheme(theme)),
					compartments.sizing.of(sizing(fontSize, maxLines)),
					compartments.editable.of(editable(readOnly)),
					keymap.of([
						...closeBracketsKeymap,
						...defaultKeymap,
						...searchKeymap,
						...historyKeymap,
						...foldKeymap,
						...completionKeymap,
						indentWithTab,
					]),
					EditorView.updateListener.of((update) => {
						if (
							update.docChanged &&
							!update.transactions.some((tr) => tr.annotation(fromProps))
						) {
							onChangeRef.current?.(update.state.doc.toString());
						}
					}),
				],
			}),
		});
		viewRef.current = view;
		return () => {
			view.destroy();
			viewRef.current = null;
		};
	}, [compartments]);

	useEffect(() => {
		const view = viewRef.current;
		if (!view) return;
		const current = view.state.doc.toString();
		if (current === value) return;
		view.dispatch({
			changes: changeBetween(current, value),
			annotations: fromProps.of(true),
		});
	}, [value]);

	useEffect(() => {
		viewRef.current?.dispatch({
			effects: compartments.sizing.reconfigure(sizing(fontSize, maxLines)),
		});
	}, [compartments, fontSize, maxLines]);

	useEffect(() => {
		viewRef.current?.dispatch({
			effects: compartments.theme.reconfigure(editorTheme(theme)),
		});
	}, [compartments, theme]);

	useEffect(() => {
		viewRef.current?.dispatch({
			effects: compartments.editable.reconfigure(editable(readOnly)),
		});
	}, [compartments, readOnly]);

	useEffect(() => {
		const view = viewRef.current;
		if (!view) return;
		const tr = diagnosticsTransaction(view.state, parseErrors(error));
		view.dispatch(tr, {
			effects: setErrorPanel.of({ message: error || null, title: errorTitle }),
		});
	}, [error, errorTitle]);

	return <div ref={parentRef} id={id} className={className} />;
}
