import { CHARACTERS, SETS, WEAPONS } from "@gcsim/editor";
import {
	actionLabel,
	actions,
	artifactLabel,
	artifacts,
	characterLabel,
	characters,
	enemies,
	enemyLabel,
	statLabel,
	stats,
	weaponLabel,
	weapons,
} from "../common/gcsim";

export const SHORTCUT_KINDS = [
	"character",
	"weapon",
	"artifact",
	"enemy",
	"action",
	"stat",
] as const;

export type ShortcutKind = (typeof SHORTCUT_KINDS)[number];

export interface ShortcutEntry {
	kind: ShortcutKind;
	key: string;
	label: string;
	aliases: string[];
}

export interface ShortcutMatch {
	entry: ShortcutEntry;
	alias?: string;
}

function aliasesByKey(names: Readonly<Record<string, string>>) {
	const out = new Map<string, string[]>();
	for (const [name, key] of Object.entries(names)) {
		if (name === key) continue;
		const list = out.get(key);
		if (list) list.push(name);
		else out.set(key, [name]);
	}
	return out;
}

const characterAliases = aliasesByKey(CHARACTERS);
const weaponAliases = aliasesByKey(WEAPONS);
const artifactAliases = aliasesByKey(SETS);

export function buildShortcutEntries(): ShortcutEntry[] {
	const entries = <T extends string>(
		kind: ShortcutKind,
		keys: readonly T[],
		label: (key: T) => string,
		aliases?: Map<string, string[]>,
	) =>
		keys.map((key) => ({
			kind,
			key,
			label: label(key) || key,
			aliases: aliases?.get(key) ?? [],
		}));

	return [
		...entries("character", characters, characterLabel, characterAliases),
		...entries("weapon", weapons, weaponLabel, weaponAliases),
		...entries("artifact", artifacts, artifactLabel, artifactAliases),
		...entries("enemy", enemies, enemyLabel),
		...entries("action", actions, actionLabel),
		...entries("stat", stats, statLabel),
	];
}

function rank(text: string, query: string) {
	if (text === query) return 0;
	if (text.startsWith(query)) return 1;
	if (text.includes(query)) return 2;
	return undefined;
}

export function searchShortcuts(
	entries: readonly ShortcutEntry[],
	query: string,
): ShortcutMatch[] {
	const q = query.trim().toLowerCase();
	if (q.length === 0) return entries.map((entry) => ({ entry }));

	const scored: { match: ShortcutMatch; score: number; order: number }[] = [];
	entries.forEach((entry, order) => {
		let score = Math.min(
			rank(entry.key, q) ?? 3,
			rank(entry.label.toLowerCase(), q) ?? 3,
		);
		let alias: string | undefined;
		for (const a of entry.aliases) {
			const r = rank(a, q);
			if (r !== undefined && r < score) {
				score = r;
				alias = a;
			}
		}
		if (score < 3) scored.push({ match: { entry, alias }, score, order });
	});

	const best = new Map<ShortcutKind, number>();
	for (const { match, score } of scored) {
		const kind = match.entry.kind;
		best.set(kind, Math.min(best.get(kind) ?? score, score));
	}
	const kindRank = (m: ShortcutMatch) =>
		(best.get(m.entry.kind) ?? 0) * SHORTCUT_KINDS.length +
		SHORTCUT_KINDS.indexOf(m.entry.kind);
	return scored
		.sort(
			(a, b) =>
				kindRank(a.match) - kindRank(b.match) ||
				a.score - b.score ||
				a.order - b.order,
		)
		.map((s) => s.match);
}
