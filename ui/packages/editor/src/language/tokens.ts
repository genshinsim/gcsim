import { ACTIONS, CHARACTERS, ELEMENTS, KEYWORDS, STATS } from "./keys";
import {
	ActionName,
	Bool,
	CharacterName,
	ElementName,
	Keyword,
	StatName,
} from "./parser.gen.terms";

const keywords = new Set(KEYWORDS);
const stats = new Set(STATS);
const elements = new Set(ELEMENTS);
const actions = new Set(ACTIONS);

// The stack argument is part of lezer's external specializer signature.
export function specializeIdentifier(value: string, _stack: unknown): number {
	if (keywords.has(value)) return Keyword;
	if (value === "true" || value === "false") return Bool;
	if (stats.has(value)) return StatName;
	if (elements.has(value)) return ElementName;
	if (Object.hasOwn(CHARACTERS, value)) return CharacterName;
	if (actions.has(value)) return ActionName;
	return -1;
}
