import type {
	Completion,
	CompletionContext,
	CompletionResult,
} from "@codemirror/autocomplete";
import { syntaxTree } from "@codemirror/language";
import {
	ACTIONS,
	CHARACTERS,
	ELEMENTS,
	KEYWORDS,
	SETS,
	STATS,
	WEAPONS,
} from "./keys";

const word = (type: string, detail?: string) => (label: string) => ({
	label,
	type,
	detail,
});

const named = (table: Readonly<Record<string, string>>, type: string) =>
	Object.entries(table).map(([label, key]) => ({
		label,
		type,
		detail: label === key ? undefined : key,
	}));

const keywords = KEYWORDS.map(word("keyword"));
const actions = ACTIONS.map(word("function", "action"));
const stats = STATS.map(word("property", "stat"));
const elements = ELEMENTS.map(word("constant", "element"));
const characters = named(CHARACTERS, "class");
const weapons = named(WEAPONS, "type");
const sets = named(SETS, "type");
const addTargets = ["weapon", "set", "stats", "params"].map(word("keyword"));
const charCommands = ["char", "add"].map(word("keyword"));
const everything = [
	...keywords,
	...characters,
	...actions,
	...stats,
	...elements,
];

const WORD = /[a-zA-Z_][\w\-%]*/;
const VALID_FOR = /^[\w\-%]*$/;

function statementBefore(context: CompletionContext): string {
	const start = Math.max(0, context.pos - 2000);
	const code = context.state
		.sliceDoc(start, context.pos)
		.replace(/(#|\/\/).*$/gm, "");
	return code.slice(
		Math.max(
			code.lastIndexOf(";"),
			code.lastIndexOf("{"),
			code.lastIndexOf("}"),
		) + 1,
	);
}

function hasUnclosedString(stmt: string) {
	return (stmt.match(/"/g)?.length ?? 0) % 2 === 1;
}

function optionsFor(stmt: string, partial: string): readonly Completion[] {
	const words = stmt
		.trim()
		.split(/[\s,]+/)
		.filter(Boolean);
	if (partial !== "") words.pop();

	if (words.length === 0) return [...characters, ...keywords];
	const [first, second, third] = words;
	if (first === "active") return characters;
	if (!Object.hasOwn(CHARACTERS, first)) return everything;

	if (words.length === 1) return [...charCommands, ...actions];
	if (second === "add") {
		if (words.length === 2) return addTargets;
		if (third === "stats") return stats;
		return [];
	}
	if (ACTIONS.includes(second)) return actions;
	return everything;
}

export function gcsimCompletionSource(
	context: CompletionContext,
): CompletionResult | null {
	const node = syntaxTree(context.state).resolveInner(context.pos, -1);
	if (node.name === "LineComment") return null;
	const stmt = statementBefore(context);

	const quoted = /\b(weapon|set)\s*=\s*"(\w*)$/.exec(stmt);
	if (quoted) {
		return {
			from: context.pos - quoted[2].length,
			options: quoted[1] === "weapon" ? weapons : sets,
			validFor: /^\w*$/,
		};
	}
	if (node.name === "StringLiteral" || hasUnclosedString(stmt)) return null;

	const before = context.matchBefore(WORD);
	if (!before && !context.explicit) return null;
	const partial = before?.text ?? "";
	const options = optionsFor(stmt, partial);
	if (options.length === 0) return null;
	return { from: before?.from ?? context.pos, options, validFor: VALID_FOR };
}
