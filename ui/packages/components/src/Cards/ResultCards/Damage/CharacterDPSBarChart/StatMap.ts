import type { model } from "@gcsim/types";
import { useMemo } from "react";

export type StatMap = { [key: string]: model.DescriptiveStats };

export type StatMapDatum = {
	name: string;
	data: StatMap;
	total: number;
};

export type StatMapChartData = {
	data: StatMapDatum[];
	keys: string[];
	xMax: number;
};

export function useStatMapData(
	maps?: (StatMap | undefined)[],
	names?: string[],
): StatMapChartData {
	return useMemo(() => {
		if (maps == null || names == null) {
			return { data: [], keys: [], xMax: 0 };
		}

		const allKeys = new Set<string>();
		const data: StatMapDatum[] = [];

		let maxDPS = 0;
		for (let i = 0; i < maps.length; i++) {
			const map = maps[i];
			if (map == null) {
				continue;
			}

			let maxTotal = 0;
			let total = 0;
			for (const key in map) {
				allKeys.add(key);
				const mean = map[key].mean ?? 0;
				maxTotal += Math.max(map[key].max ?? 0, mean + (map[key].sd ?? 0));
				total += mean;
			}
			maxDPS = Math.max(maxDPS, maxTotal);
			data.push({ name: names[i], data: map, total: total });
		}

		return {
			data: data,
			keys: Array.from(allKeys),
			xMax: maxDPS,
		};
	}, [maps, names]);
}
