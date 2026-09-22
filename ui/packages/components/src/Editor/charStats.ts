import type { model } from "@gcsim/types";
import type { TFunction } from "i18next";
import type { CharStatBlock } from "../Cards";

export const StatToIndexMap: Record<string, number> = {
	DEFP: 1,
	DEF: 2,
	HP: 3,
	HPP: 4,
	ATK: 5,
	ATKP: 6,
	ER: 7,
	EM: 8,
	CR: 9,
	CD: 10,
	Heal: 11,
	PyroP: 12,
	HydroP: 13,
	CryoP: 14,
	ElectroP: 15,
	AnemoP: 16,
	GeoP: 17,
	DendroP: 18,
	PhyP: 19,
};

type StatKind = "both" | "f" | "%";

interface StatRow {
	key: string;
	name: (t: TFunction) => string;
	flatIndex: number;
	percentIndex: number;
	t: StatKind;
}

// er/cr/cd/element%/heal render identically in the total-stats and snapshot views.
const percentRows: StatRow[] = [
	{ key: "er", name: (t) => t("stats.er"), flatIndex: -1, percentIndex: StatToIndexMap.ER, t: "%" },
	{ key: "cr", name: (t) => t("stats.cr"), flatIndex: -1, percentIndex: StatToIndexMap.CR, t: "%" },
	{ key: "cd", name: (t) => t("stats.cd"), flatIndex: -1, percentIndex: StatToIndexMap.CD, t: "%" },
	{ key: "electro", name: (t) => t("stats.electro%"), flatIndex: -1, percentIndex: StatToIndexMap.ElectroP, t: "%" },
	{ key: "pyro", name: (t) => t("stats.pyro%"), flatIndex: -1, percentIndex: StatToIndexMap.PyroP, t: "%" },
	{ key: "cryo", name: (t) => t("stats.cryo%"), flatIndex: -1, percentIndex: StatToIndexMap.CryoP, t: "%" },
	{ key: "hydro", name: (t) => t("stats.hydro%"), flatIndex: -1, percentIndex: StatToIndexMap.HydroP, t: "%" },
	{ key: "geo", name: (t) => t("stats.geo%"), flatIndex: -1, percentIndex: StatToIndexMap.GeoP, t: "%" },
	{ key: "anemo", name: (t) => t("stats.anemo%"), flatIndex: -1, percentIndex: StatToIndexMap.AnemoP, t: "%" },
	{ key: "phys", name: (t) => t("stats.phys%"), flatIndex: -1, percentIndex: StatToIndexMap.PhyP, t: "%" },
	{ key: "dendro", name: (t) => t("stats.dendro%"), flatIndex: -1, percentIndex: StatToIndexMap.DendroP, t: "%" },
	{ key: "heal", name: (t) => t("stats.heal"), flatIndex: -1, percentIndex: StatToIndexMap.Heal, t: "%" },
];

const totalRows: StatRow[] = [
	{ key: "hp", name: (t) => `${t("stats.hp")} / ${t("stats.hp%")}`, flatIndex: StatToIndexMap.HP, percentIndex: StatToIndexMap.HPP, t: "both" },
	{ key: "atk", name: (t) => `${t("stats.atk")} / ${t("stats.atk%")}`, flatIndex: StatToIndexMap.ATK, percentIndex: StatToIndexMap.ATKP, t: "both" },
	{ key: "def", name: (t) => `${t("stats.def")} / ${t("stats.def%")}`, flatIndex: StatToIndexMap.DEF, percentIndex: StatToIndexMap.DEFP, t: "both" },
	{ key: "em", name: (t) => t("stats.em"), flatIndex: StatToIndexMap.EM, percentIndex: -1, t: "f" },
	...percentRows,
];

const snapshotRows: StatRow[] = [
	{ key: "hp", name: (t) => t("stats.hp"), flatIndex: StatToIndexMap.HP, percentIndex: -1, t: "f" },
	{ key: "atk", name: (t) => t("stats.atk"), flatIndex: StatToIndexMap.ATK, percentIndex: -1, t: "f" },
	{ key: "def", name: (t) => t("stats.def"), flatIndex: StatToIndexMap.DEF, percentIndex: -1, t: "f" },
	{ key: "em", name: (t) => t("stats.em"), flatIndex: StatToIndexMap.EM, percentIndex: -1, t: "f" },
	...percentRows,
];

function buildBlocks(
	rows: StatRow[],
	chars: model.Character[],
	source: (c: model.Character) => number[],
	t: TFunction,
): { blocks: Record<string, CharStatBlock[]>; maxRows: number } {
	const values: Record<string, Record<string, { flat: number; per: number }>> =
		{};
	const counts: Record<string, number> = {};
	rows.forEach((row) => {
		values[row.key] = {};
		counts[row.key] = 0;
	});

	let maxRows = 0;
	chars.forEach((char) => {
		let rowCount = 0;
		const name = char.name ?? "";
		const stats = source(char);
		rows.forEach((row) => {
			if (!(name in values[row.key])) {
				values[row.key][name] = { flat: 0, per: 0 };
			}
			if ((stats[row.percentIndex] ?? 0) > 0 || (stats[row.flatIndex] ?? 0) > 0) {
				counts[row.key]++;
				rowCount++;
			}
			if (row.t === "both" || row.t === "f") {
				values[row.key][name].flat = stats[row.flatIndex];
			}
			if (row.t === "both" || row.t === "%") {
				values[row.key][name].per = stats[row.percentIndex];
			}
		});
		if (rowCount > maxRows) {
			maxRows = rowCount;
		}
	});

	const blocks: Record<string, CharStatBlock[]> = {};
	chars.forEach((char) => {
		blocks[char.name ?? ""] = [];
	});
	rows.forEach((row) => {
		if (counts[row.key] === 0) {
			return;
		}
		for (const name in values[row.key]) {
			blocks[name].push({
				key: row.key,
				name: row.name(t),
				t: row.t,
				flat: values[row.key][name].flat,
				percent: values[row.key][name].per,
			});
		}
	});
	return { blocks, maxRows };
}

export function ConsolidateCharStats(
	t: TFunction,
	chars: model.Character[],
): {
	stats: Record<string, CharStatBlock[]>;
	snapshot: Record<string, CharStatBlock[]>;
	maxRows: number;
} {
	const total = buildBlocks(totalRows, chars, (c) => c.stats ?? [], t);
	const snap = buildBlocks(snapshotRows, chars, (c) => c.snapshot ?? [], t);
	return {
		stats: total.blocks,
		snapshot: snap.blocks,
		maxRows: Math.max(total.maxRows, snap.maxRows),
	};
}
